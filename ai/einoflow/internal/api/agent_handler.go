package api

import (
	"net/http"

	"einoflow/internal/agent"
	"einoflow/internal/tools"
	"einoflow/pkg/logger"

	"github.com/cloudwego/eino/components/model"
	"github.com/gin-gonic/gin"
)

type AgentHandler struct {
	chatModel    model.ChatModel
	toolRegistry *tools.Registry
	toolExecutor *tools.Executor
}

func NewAgentHandler(chatModel model.ChatModel) *AgentHandler {
	// 创建工具注册表
	toolRegistry := tools.NewRegistry("./data/einoflow.db", "./data/files")
	toolExecutor := tools.NewExecutor(toolRegistry)

	return &AgentHandler{
		chatModel:    chatModel,
		toolRegistry: toolRegistry,
		toolExecutor: toolExecutor,
	}
}

// AgentRequest Agent 请求
type AgentRequest struct {
	Task string `json:"task" binding:"required"`
}

// AgentResponse Agent 响应
type AgentResponse struct {
	Answer string `json:"answer"`
}

func (h *AgentHandler) Run(c *gin.Context) {
	var req AgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 创建 Agent 并设置工具执行器
	reactAgent := agent.NewReActAgent(h.chatModel).SetToolExecutor(h.toolExecutor)

	// 执行任务
	result, err := reactAgent.Run(c.Request.Context(), req.Task)
	if err != nil {
		logger.Error("Agent execution failed: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, &AgentResponse{
		Answer: result,
	})
}
