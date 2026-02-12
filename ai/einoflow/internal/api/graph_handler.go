package api

import (
	"net/http"

	"einoflow/internal/graph"
	"einoflow/pkg/logger"

	"github.com/cloudwego/eino/components/model"
	"github.com/gin-gonic/gin"
)

// GraphHandler Graph 处理器
type GraphHandler struct {
	chatModel model.ChatModel
}

// NewGraphHandler 创建 Graph 处理器
func NewGraphHandler(chatModel model.ChatModel) *GraphHandler {
	return &GraphHandler{chatModel: chatModel}
}

// GraphRequest Graph 请求
type GraphRequest struct {
	Query string `json:"query" binding:"required"`
	Type  string `json:"type"` // multi_step, rag
}

// GraphResponse Graph 响应
type GraphResponse struct {
	Result string   `json:"result"`
	Steps  []string `json:"steps"`
}

// Run 执行 Graph
func (h *GraphHandler) Run(c *gin.Context) {
	var req GraphRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 创建多步骤图
	g := graph.CreateMultiStepGraph(h.chatModel)

	// 执行图
	result, err := g.Run(c.Request.Context(), req.Query)
	if err != nil {
		logger.Error("Graph execution failed: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, &GraphResponse{
		Result: result.(string),
		Steps:  []string{"analyze", "plan", "execute"},
	})
}
