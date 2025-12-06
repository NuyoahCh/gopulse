package api

import (
	"context"
	"net/http"

	"einoflow/internal/chain"
	"einoflow/pkg/logger"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"github.com/gin-gonic/gin"
)

type ChainHandler struct {
	chatModel model.ChatModel
}

func NewChainHandler(chatModel model.ChatModel) *ChainHandler {
	return &ChainHandler{chatModel: chatModel}
}

type ChainRequest struct {
	Steps []string `json:"steps" binding:"required"`
	Input string   `json:"input" binding:"required"`
}

type ChainResponse struct {
	Result string `json:"result"`
	Steps  int    `json:"steps"`
}

func (h *ChainHandler) Run(c *gin.Context) {
	var req ChainRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 创建顺序链
	seqChain := chain.NewSequentialChain()
	for _, stepPrompt := range req.Steps {
		prompt := stepPrompt // 捕获变量
		lambda := func(ctx context.Context, input string) (string, error) {
			messages := []*schema.Message{
				schema.SystemMessage(prompt),
				schema.UserMessage(input),
			}
			resp, err := h.chatModel.Generate(ctx, messages)
			if err != nil {
				return "", err
			}
			return resp.Content, nil
		}
		seqChain.AddLambdaStep(lambda)
	}

	result, err := seqChain.Run(c.Request.Context(), req.Input)
	if err != nil {
		logger.Error("Chain execution failed: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, &ChainResponse{
		Result: result,
		Steps:  len(req.Steps),
	})
}
