package api

import (
	"einoflow/internal/config"
	"einoflow/internal/llm"
	"einoflow/internal/llm/providers"

	"github.com/cloudwego/eino/components/model"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, cfg *config.Config) {
	// 初始化 LLM 管理器
	llmManager := llm.NewManager()

	// 注册提供商（优先注册豆包）
	var defaultChatModel model.ChatModel
	if cfg.ArkAPIKey != "" {
		provider := providers.NewArkProvider(cfg.ArkAPIKey, cfg.ArkBaseURL)
		llmManager.RegisterProvider(provider)
		// 获取默认模型
		defaultChatModel, _ = provider.GetChatModel("doubao-seed-1-6-lite-251015")
	}

	if cfg.OpenAIKey != "" {
		llmManager.RegisterProvider(providers.NewOpenAIProvider(cfg.OpenAIKey, cfg.OpenAIBaseURL))
	}

	// API v1 路由组
	v1 := router.Group("/api/v1")
	{
		// LLM 相关
		llmGroup := v1.Group("/llm")
		{
			llmHandler := NewLLMHandler(llmManager)
			llmGroup.POST("/chat", llmHandler.Chat)
			llmGroup.POST("/chat/stream", llmHandler.ChatStream)
			llmGroup.GET("/models", llmHandler.ListModels)
		}

		// Chain 相关
		chainGroup := v1.Group("/chain")
		{
			chainHandler := NewChainHandler(defaultChatModel)
			chainGroup.POST("/run", chainHandler.Run)
		}

		// Agent 相关
		agentGroup := v1.Group("/agent")
		{
			agentHandler := NewAgentHandler(defaultChatModel)
			agentGroup.POST("/run", agentHandler.Run)
		}

		// RAG 相关
		ragGroup := v1.Group("/rag")
		{
			ragHandler := NewRAGHandler(
				defaultChatModel,
				cfg.ArkAPIKey,
				cfg.ArkBaseURL,
				cfg.ArkEmbeddingModel,
			)
			ragGroup.POST("/index", ragHandler.Index)
			ragGroup.POST("/upload", ragHandler.UploadFile) // 文件上传
			ragGroup.POST("/query", ragHandler.Query)
			ragGroup.GET("/stats", ragHandler.GetStats) // 查看存储的文档
			ragGroup.DELETE("/clear", ragHandler.Clear) // 清空文档
		}

		// Graph 相关
		graphGroup := v1.Group("/graph")
		{
			graphHandler := NewGraphHandler(defaultChatModel)
			graphGroup.POST("/run", graphHandler.Run)
		}

		// 多模态相关
		multimodalGroup := v1.Group("/multimodal")
		{
			multimodalHandler := NewMultimodalHandler(defaultChatModel)
			multimodalGroup.POST("/chat", multimodalHandler.ChatWithImage)
		}
	}
}
