package api

import (
	"net/http"

	"einoflow/internal/multimodal"
	"einoflow/pkg/logger"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"github.com/gin-gonic/gin"
)

type MultimodalHandler struct {
	chatModel    model.ChatModel
	imageHandler *multimodal.ImageHandler
}

func NewMultimodalHandler(chatModel model.ChatModel) *MultimodalHandler {
	return &MultimodalHandler{
		chatModel:    chatModel,
		imageHandler: multimodal.NewImageHandler(),
	}
}

type ImageChatRequest struct {
	Text     string `json:"text" binding:"required"`
	ImageURL string `json:"image_url,omitempty"`
	ImageB64 string `json:"image_b64,omitempty"`
	MimeType string `json:"mime_type,omitempty"` // image/jpeg, image/png, etc.
}

type ImageChatResponse struct {
	Answer string `json:"answer"`
}

func (h *MultimodalHandler) ChatWithImage(c *gin.Context) {
	var req ImageChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var message *schema.Message

	if req.ImageURL != "" {
		// 使用图像 URL
		message = h.imageHandler.CreateImageMessage("user", req.Text, req.ImageURL)
	} else if req.ImageB64 != "" {
		// 使用 base64 图像
		mimeType := req.MimeType
		if mimeType == "" {
			mimeType = "image/jpeg" // 默认
		}
		message = h.imageHandler.CreateImageMessageFromBase64("user", req.Text, req.ImageB64, mimeType)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "either image_url or image_b64 is required"})
		return
	}

	// 调用模型
	resp, err := h.chatModel.Generate(c.Request.Context(), []*schema.Message{message})
	if err != nil {
		logger.Error("Multimodal chat failed: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, &ImageChatResponse{
		Answer: resp.Content,
	})
}
