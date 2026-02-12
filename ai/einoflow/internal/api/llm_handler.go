package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"einoflow/internal/llm"
	"einoflow/internal/memory"
	"einoflow/internal/middleware"
	"einoflow/pkg/logger"

	"github.com/gin-gonic/gin"
)

type LLMHandler struct {
	manager        *llm.Manager
	contextManager *memory.ContextManager
}

func NewLLMHandler(manager *llm.Manager) *LLMHandler {
	// 创建上下文管理器（4096 tokens，适合大多数模型）
	ctxMgr, err := memory.NewContextManager(4096)
	if err != nil {
		logger.Warn(fmt.Sprintf("Failed to create context manager: %v", err))
	}

	return &LLMHandler{
		manager:        manager,
		contextManager: ctxMgr,
	}
}

// Chat godoc
// @Summary      LLM 聊天接口
// @Description  与 LLM 进行对话，支持多个模型提供商
// @Tags         LLM
// @Accept       json
// @Produce      json
// @Param        request  body      object  true  "聊天请求 {model: string, messages: array}"
// @Success      200      {object}  object  "聊天响应"
// @Failure      400      {object}  map[string]interface{}  "请求格式错误"
// @Failure      500      {object}  map[string]interface{}  "服务器错误"
// @Router       /llm/chat [post]
func (h *LLMHandler) Chat(c *gin.Context) {
	requestID := middleware.GetRequestID(c)

	var req llm.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.WithFields(map[string]interface{}{
			"request_id": requestID,
			"error":      err.Error(),
		}).Warn("Invalid chat request")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":      "Invalid request format",
			"request_id": requestID,
		})
		return
	}

	// 截断消息以适应上下文窗口
	if h.contextManager != nil {
		originalCount := len(req.Messages)
		req.Messages = h.contextManager.TruncateMessages(req.Messages)
		if len(req.Messages) < originalCount {
			logger.WithFields(map[string]interface{}{
				"request_id":       requestID,
				"model":            req.Model,
				"original_count":   originalCount,
				"truncated_count":  len(req.Messages),
				"tokens":           h.contextManager.CountTokens(req.Messages),
				"available_tokens": h.contextManager.GetAvailableTokens(req.Messages),
			}).Info("Context truncated")
		}
	}

	resp, err := h.manager.Chat(c.Request.Context(), &req)
	if err != nil {
		logger.WithFields(map[string]interface{}{
			"request_id": requestID,
			"model":      req.Model,
			"error":      err.Error(),
		}).Error("Chat request failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "Failed to process chat request",
			"request_id": requestID,
		})
		return
	}

	logger.WithFields(map[string]interface{}{
		"request_id": requestID,
		"model":      req.Model,
	}).Debug("Chat request completed")

	c.JSON(http.StatusOK, resp)
}

// ChatStream godoc
// @Summary      LLM 流式聊天接口
// @Description  与 LLM 进行流式对话，实时返回响应
// @Tags         LLM
// @Accept       json
// @Produce      text/event-stream
// @Param        request  body      object  true  "聊天请求 {model: string, messages: array}"
// @Success      200      {string}  string  "SSE 流式响应"
// @Failure      400      {object}  map[string]interface{}  "请求格式错误"
// @Failure      500      {object}  map[string]interface{}  "服务器错误"
// @Router       /llm/chat/stream [post]
func (h *LLMHandler) ChatStream(c *gin.Context) {
	requestID := middleware.GetRequestID(c)

	var req llm.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.WithFields(map[string]interface{}{
			"request_id": requestID,
			"error":      err.Error(),
		}).Warn("Invalid stream chat request")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":      "Invalid request format",
			"request_id": requestID,
		})
		return
	}

	// 截断消息以适应上下文窗口
	if h.contextManager != nil {
		originalCount := len(req.Messages)
		req.Messages = h.contextManager.TruncateMessages(req.Messages)
		if len(req.Messages) < originalCount {
			logger.WithFields(map[string]interface{}{
				"request_id":      requestID,
				"model":           req.Model,
				"original_count":  originalCount,
				"truncated_count": len(req.Messages),
			}).Info("Stream context truncated")
		}
	}

	req.Stream = true

	// 设置 SSE 头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Transfer-Encoding", "chunked")

	// 获取流式响应
	stream, err := h.manager.ChatStream(c.Request.Context(), &req)
	if err != nil {
		logger.WithFields(map[string]interface{}{
			"request_id": requestID,
			"model":      req.Model,
			"error":      err.Error(),
		}).Error("Stream chat request failed")
		c.SSEvent("error", gin.H{
			"error":      "Failed to process stream chat request",
			"request_id": requestID,
		})
		return
	}

	// 处理流式响应
	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "streaming not supported"})
		return
	}

	err = llm.ProcessStream(c.Request.Context(), stream, func(chunk *llm.StreamChunk) error {
		data, _ := json.Marshal(chunk)
		fmt.Fprintf(c.Writer, "data: %s\n\n", data)
		flusher.Flush()
		return nil
	})

	if err != nil {
		logger.Error("Stream processing failed: " + err.Error())
	}
}

// ListModels godoc
// @Summary      获取可用模型列表
// @Description  获取所有已配置的 LLM 模型列表
// @Tags         LLM
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "模型列表"
// @Router       /llm/models [get]
func (h *LLMHandler) ListModels(c *gin.Context) {
	providerName := c.Query("provider")

	if providerName != "" {
		// 获取特定提供商的模型
		provider, ok := h.manager.GetProvider(providerName)
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "provider not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"provider": providerName,
			"models":   provider.ListModels(),
		})
		return
	}

	// 返回所有提供商的模型，转换为前端期望的格式
	type ModelInfo struct {
		ID       string `json:"id"`
		Provider string `json:"provider"`
		Name     string `json:"name"`
	}

	var models []ModelInfo
	for providerName, provider := range h.manager.GetAllProviders() {
		modelIDs := provider.ListModels()
		for _, modelID := range modelIDs {
			models = append(models, ModelInfo{
				ID:       modelID,
				Provider: providerName,
				Name:     modelID, // 可以后续优化为更友好的名称
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"models": models,
	})
}
