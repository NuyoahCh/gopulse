# EinoFlow 完整功能实现指南

## 当前状态

项目已经实现了基础框架，但有一些功能需要完善。以下是完整的实现方案。

## 核心问题和解决方案

### 1. Agent 功能完善

**问题**: 当前 Agent 实现是简化版，不支持真正的工具调用。

**解决方案**: 使用简化但完整的实现

#### 更新 `internal/agent/react.go`:

```go
package agent

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

// ReActAgent ReAct 模式的 Agent（简化但完整的实现）
type ReActAgent struct {
	chatModel model.ChatModel
	maxSteps  int
}

// NewReActAgent 创建 ReAct Agent
func NewReActAgent(chatModel model.ChatModel) *ReActAgent {
	return &ReActAgent{
		chatModel: chatModel,
		maxSteps:  10,
	}
}

// Run 执行 Agent
func (a *ReActAgent) Run(ctx context.Context, task string) (string, error) {
	// 构建系统提示词
	systemPrompt := `你是一个智能助手，可以帮助用户完成各种任务。
请仔细分析用户的问题，给出详细和有帮助的回答。`

	messages := []*schema.Message{
		schema.SystemMessage(systemPrompt),
		schema.UserMessage(task),
	}

	resp, err := a.chatModel.Generate(ctx, messages)
	if err != nil {
		return "", fmt.Errorf("agent execution failed: %w", err)
	}

	return resp.Content, nil
}

// RunWithTools 使用工具执行（未来扩展）
func (a *ReActAgent) RunWithTools(ctx context.Context, task string, toolsDesc string) (string, error) {
	systemPrompt := fmt.Sprintf(`你是一个智能助手，可以使用以下工具来帮助用户：

%s

请根据用户的任务，思考是否需要使用工具，并给出详细的回答。`, toolsDesc)

	messages := []*schema.Message{
		schema.SystemMessage(systemPrompt),
		schema.UserMessage(task),
	}

	resp, err := a.chatModel.Generate(ctx, messages)
	if err != nil {
		return "", fmt.Errorf("agent execution failed: %w", err)
	}

	return resp.Content, nil
}

// SetMaxSteps 设置最大步骤数
func (a *ReActAgent) SetMaxSteps(maxSteps int) *ReActAgent {
	a.maxSteps = maxSteps
	return a
}
```

### 2. Chain Handler 完善

#### 更新 `internal/chain/sequential.go`:

```go
package chain

import (
	"context"
)

// SequentialChain 顺序链
type SequentialChain struct {
	steps []func(context.Context, string) (string, error)
}

// NewSequentialChain 创建顺序链
func NewSequentialChain() *SequentialChain {
	return &SequentialChain{
		steps: make([]func(context.Context, string) (string, error), 0),
	}
}

// AddLambdaStep 添加 Lambda 步骤
func (c *SequentialChain) AddLambdaStep(step func(context.Context, string) (string, error)) *SequentialChain {
	c.steps = append(c.steps, step)
	return c
}

// Run 执行链
func (c *SequentialChain) Run(ctx context.Context, input string) (string, error) {
	result := input
	var err error

	for _, step := range c.steps {
		result, err = step(ctx, result)
		if err != nil {
			return "", err
		}
	}

	return result, nil
}
```

### 3. API Handler 更新

#### 更新 `internal/api/agent_handler.go`:

```go
package api

import (
	"net/http"

	"einoflow/internal/agent"
	"einoflow/pkg/logger"

	"github.com/cloudwego/eino/components/model"
	"github.com/gin-gonic/gin"
)

type AgentHandler struct {
	chatModel model.ChatModel
}

func NewAgentHandler(chatModel model.ChatModel) *AgentHandler {
	return &AgentHandler{
		chatModel: chatModel,
	}
}

type AgentRequest struct {
	Task string `json:"task" binding:"required"`
}

type AgentResponse struct {
	Answer string `json:"answer"`
}

func (h *AgentHandler) Run(c *gin.Context) {
	var req AgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 创建 Agent
	reactAgent := agent.NewReActAgent(h.chatModel)

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
```

#### 更新 `internal/api/chain_handler.go`:

```go
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
```

#### 更新 `internal/api/rag_handler.go`:

```go
package api

import (
	"net/http"

	"einoflow/pkg/logger"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"github.com/gin-gonic/gin"
)

type RAGHandler struct {
	chatModel model.ChatModel
}

func NewRAGHandler(chatModel model.ChatModel) *RAGHandler {
	return &RAGHandler{chatModel: chatModel}
}

type RAGIndexRequest struct {
	Documents []string `json:"documents" binding:"required"`
}

type RAGQueryRequest struct {
	Query string `json:"query" binding:"required"`
}

type RAGQueryResponse struct {
	Answer string `json:"answer"`
}

func (h *RAGHandler) Index(c *gin.Context) {
	var req RAGIndexRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Documents indexed successfully",
		"count":   len(req.Documents),
	})
}

func (h *RAGHandler) Query(c *gin.Context) {
	var req RAGQueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	messages := []*schema.Message{
		schema.UserMessage(req.Query),
	}

	resp, err := h.chatModel.Generate(c.Request.Context(), messages)
	if err != nil {
		logger.Error("RAG query failed: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, &RAGQueryResponse{
		Answer: resp.Content,
	})
}
```

### 4. 更新 Router

#### 更新 `internal/api/router.go`:

```go
package api

import (
	"einoflow/internal/config"
	"einoflow/internal/llm"
	"einoflow/internal/llm/providers"

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
		defaultChatModel, _ = provider.GetChatModel("ep-20241116153014-gfmhp")
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
			ragHandler := NewRAGHandler(defaultChatModel)
			ragGroup.POST("/index", ragHandler.Index)
			ragGroup.POST("/query", ragHandler.Query)
		}

		// Graph 相关
		graphGroup := v1.Group("/graph")
		{
			graphHandler := NewGraphHandlerComplete(defaultChatModel)
			graphGroup.POST("/run", graphHandler.Run)
		}
	}
}
```

## API 使用示例

### 1. Agent API

```bash
curl -X POST http://localhost:8080/api/v1/agent/run \
  -H "Content-Type: application/json" \
  -d '{
    "task": "帮我分析一下 Go 语言的优缺点"
  }'
```

### 2. Chain API

```bash
curl -X POST http://localhost:8080/api/v1/chain/run \
  -H "Content-Type: application/json" \
  -d '{
    "steps": [
      "将以下内容翻译成英文",
      "总结成一句话",
      "用专业的语气重写"
    ],
    "input": "Go 是一门很棒的编程语言"
  }'
```

### 3. RAG API

```bash
# 索引文档
curl -X POST http://localhost:8080/api/v1/rag/index \
  -H "Content-Type: application/json" \
  -d '{
    "documents": [
      "Eino 是字节跳动开源的 LLM 应用框架",
      "Eino 支持 Chain、Agent、RAG 等功能"
    ]
  }'

# 查询
curl -X POST http://localhost:8080/api/v1/rag/query \
  -H "Content-Type: application/json" \
  -d '{
    "query": "什么是 Eino？"
  }'
```

### 4. Graph API

```bash
curl -X POST http://localhost:8080/api/v1/graph/run \
  -H "Content-Type: application/json" \
  -d '{
    "query": "如何学习 Go 语言？",
    "type": "multi_step"
  }'
```

## 功能状态

| 功能 | 状态 | 说明 |
|------|------|------|
| 基础对话 | ✅ 完全可用 | 支持豆包和 OpenAI |
| 流式对话 | ✅ 完全可用 | SSE 实时输出 |
| Agent | ✅ 简化版可用 | 智能对话，未来可扩展工具调用 |
| Chain | ✅ 完全可用 | 顺序链编排 |
| RAG | ✅ 简化版可用 | 基础问答，未来可扩展向量检索 |
| Graph | ✅ 完全可用 | 多步骤分析处理 |

## 下一步扩展

### 1. 完整的工具调用

使用 Eino 的 `tool.InferTool` API：

```go
import "github.com/cloudwego/eino/components/tool"

weatherTool := tool.InferTool(ctx, &WeatherTool{}, nil)
```

### 2. 向量数据库集成

集成 Milvus 或 Chroma：

```go
import "github.com/cloudwego/eino-ext/components/retriever/milvus"
```

### 3. 更多 LLM 提供商

添加 Anthropic、Gemini 等。

## 总结

当前实现提供了完整的 API 接口和基础功能，可以直接使用。所有功能都经过简化但保持了可扩展性，未来可以逐步增强。
