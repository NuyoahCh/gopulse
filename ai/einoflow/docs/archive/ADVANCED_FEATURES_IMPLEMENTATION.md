# 高级功能实现指南

## 🎯 待实现的三个功能

### 1. 真正的 Embedding 模型（字节豆包）✅ 已创建
### 2. 上下文窗口管理
### 3. 多模态支持

---

## 1️⃣ 真正的 Embedding 模型

### 已创建文件
- `internal/embedding/ark_embedder.go` ✅

### 需要完成的步骤

#### 步骤 1: 更新配置
编辑 `.env` 文件，添加：
```env
ARK_EMBEDDING_MODEL="doubao-embedding-large-text-250515"
```

编辑 `internal/config/config.go`，添加字段：
```go
type Config struct {
    // ... 现有字段
    ArkEmbeddingModel string `mapstructure:"ARK_EMBEDDING_MODEL"`
}
```

#### 步骤 2: 更新 RAG Handler 的 Index 方法

在 `internal/api/rag_handler.go` 的 `Index` 方法中，替换：
```go
// 旧代码（第98-99行）
// 简单的向量化：使用字符串长度和内容特征
embeddings[i] = simpleEmbedding(content)
```

为：
```go
// 使用真实 Embedding 或简单 Embedding
if h.useRealEmbedding {
    embedding, err := h.embedder.EmbedText(c.Request.Context(), content)
    if err != nil {
        logger.Warn(fmt.Sprintf("Failed to embed text, using simple embedding: %v", err))
        embeddings[i] = simpleEmbedding(content)
    } else {
        embeddings[i] = embedding
    }
} else {
    embeddings[i] = simpleEmbedding(content)
}
```

#### 步骤 3: 更新 Query 方法

在 `Query` 方法中（约第147行），替换：
```go
queryEmbedding := simpleEmbedding(req.Query)
```

为：
```go
var queryEmbedding []float64
if h.useRealEmbedding {
    var err error
    queryEmbedding, err = h.embedder.EmbedText(c.Request.Context(), req.Query)
    if err != nil {
        logger.Warn(fmt.Sprintf("Failed to embed query, using simple embedding: %v", err))
        queryEmbedding = simpleEmbedding(req.Query)
    }
} else {
    queryEmbedding = simpleEmbedding(req.Query)
}
```

#### 步骤 4: 更新 Router

在 `internal/api/router.go` 中（约第58行），替换：
```go
ragHandler := NewRAGHandler(defaultChatModel)
```

为：
```go
ragHandler := NewRAGHandler(
    defaultChatModel,
    cfg.ArkAPIKey,
    cfg.ArkBaseURL,
    cfg.ArkEmbeddingModel,
)
```

#### 步骤 5: 安装依赖
```bash
go get github.com/volcengine/volcengine-go-sdk
```

---

## 2️⃣ 上下文窗口管理

### 创建文件：`internal/memory/context_manager.go`

```go
package memory

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/schema"
	"github.com/pkoukk/tiktoken-go"
)

// ContextManager 上下文窗口管理器
type ContextManager struct {
	maxTokens     int
	tokenizer     *tiktoken.Tiktoken
	reserveTokens int // 为响应预留的 token 数
}

// NewContextManager 创建上下文管理器
func NewContextManager(maxTokens int) (*ContextManager, error) {
	// 使用 cl100k_base 编码器（适用于 GPT-4 和豆包）
	tokenizer, err := tiktoken.GetEncoding("cl100k_base")
	if err != nil {
		return nil, fmt.Errorf("failed to get tokenizer: %w", err)
	}

	return &ContextManager{
		maxTokens:     maxTokens,
		tokenizer:     tokenizer,
		reserveTokens: 1000, // 为响应预留 1000 tokens
	}, nil
}

// CountTokens 计算消息的 token 数
func (cm *ContextManager) CountTokens(messages []*schema.Message) int {
	totalTokens := 0
	for _, msg := range messages {
		// 每条消息的开销：role + content + 格式化
		tokens := cm.tokenizer.Encode(msg.Content, nil, nil)
		totalTokens += len(tokens) + 4 // +4 for message formatting
	}
	return totalTokens + 3 // +3 for reply priming
}

// TruncateMessages 截断消息以适应上下文窗口
func (cm *ContextManager) TruncateMessages(messages []*schema.Message) []*schema.Message {
	if len(messages) == 0 {
		return messages
	}

	// 计算当前 token 数
	currentTokens := cm.CountTokens(messages)
	maxAllowed := cm.maxTokens - cm.reserveTokens

	if currentTokens <= maxAllowed {
		return messages // 不需要截断
	}

	// 保留系统消息（如果有）
	var systemMsg *schema.Message
	startIdx := 0
	if len(messages) > 0 && messages[0].Role == "system" {
		systemMsg = messages[0]
		startIdx = 1
	}

	// 从最新的消息开始保留
	result := make([]*schema.Message, 0)
	if systemMsg != nil {
		result = append(result, systemMsg)
	}

	// 从后往前添加消息，直到达到 token 限制
	tokens := 0
	if systemMsg != nil {
		tokens = cm.CountTokens([]*schema.Message{systemMsg})
	}

	for i := len(messages) - 1; i >= startIdx; i-- {
		msgTokens := cm.CountTokens([]*schema.Message{messages[i]})
		if tokens+msgTokens > maxAllowed {
			break
		}
		result = append([]*schema.Message{messages[i]}, result...)
		tokens += msgTokens
	}

	return result
}

// GetMaxTokens 获取最大 token 数
func (cm *ContextManager) GetMaxTokens() int {
	return cm.maxTokens
}

// GetAvailableTokens 获取可用的 token 数
func (cm *ContextManager) GetAvailableTokens(messages []*schema.Message) int {
	used := cm.CountTokens(messages)
	return cm.maxTokens - used - cm.reserveTokens
}
```

### 使用示例

在 `internal/api/llm_handler.go` 中使用：

```go
import "einoflow/internal/memory"

type LLMHandler struct {
	manager        *llm.Manager
	contextManager *memory.ContextManager
}

func NewLLMHandler(manager *llm.Manager) *LLMHandler {
	// 创建上下文管理器（4096 tokens）
	ctxMgr, _ := memory.NewContextManager(4096)
	
	return &LLMHandler{
		manager:        manager,
		contextManager: ctxMgr,
	}
}

func (h *LLMHandler) Chat(c *gin.Context) {
	var req ChatRequest
	// ... 解析请求

	// 截断消息以适应上下文窗口
	if h.contextManager != nil {
		req.Messages = h.contextManager.TruncateMessages(req.Messages)
		logger.Info(fmt.Sprintf("Context tokens: %d, available: %d",
			h.contextManager.CountTokens(req.Messages),
			h.contextManager.GetAvailableTokens(req.Messages)))
	}

	// ... 继续处理
}
```

### 安装依赖
```bash
go get github.com/pkoukk/tiktoken-go
```

---

## 3️⃣ 多模态支持

### 创建文件：`internal/multimodal/image_handler.go`

```go
package multimodal

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/cloudwego/eino/schema"
)

// ImageHandler 图像处理器
type ImageHandler struct {
	maxImageSize int64 // 最大图像大小（字节）
}

// NewImageHandler 创建图像处理器
func NewImageHandler() *ImageHandler {
	return &ImageHandler{
		maxImageSize: 20 * 1024 * 1024, // 20MB
	}
}

// LoadImageFromURL 从 URL 加载图像
func (h *ImageHandler) LoadImageFromURL(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch image: status %d", resp.StatusCode)
	}

	// 读取图像数据
	data, err := io.ReadAll(io.LimitReader(resp.Body, h.maxImageSize))
	if err != nil {
		return "", fmt.Errorf("failed to read image: %w", err)
	}

	// 转换为 base64
	return base64.StdEncoding.EncodeToString(data), nil
}

// LoadImageFromFile 从文件加载图像
func (h *ImageHandler) LoadImageFromFile(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	if int64(len(data)) > h.maxImageSize {
		return "", fmt.Errorf("image too large: %d bytes (max %d)", len(data), h.maxImageSize)
	}

	return base64.StdEncoding.EncodeToString(data), nil
}

// CreateImageMessage 创建包含图像的消息
func (h *ImageHandler) CreateImageMessage(role, text, imageURL string) *schema.Message {
	return &schema.Message{
		Role: role,
		Content: text,
		MultiContent: []*schema.MultiContent{
			{
				Type: "text",
				Text: text,
			},
			{
				Type:     "image_url",
				ImageURL: &schema.ImageURL{URL: imageURL},
			},
		},
	}
}

// CreateImageMessageFromBase64 创建包含 base64 图像的消息
func (h *ImageHandler) CreateImageMessageFromBase64(role, text, base64Data, mimeType string) *schema.Message {
	dataURL := fmt.Sprintf("data:%s;base64,%s", mimeType, base64Data)
	return h.CreateImageMessage(role, text, dataURL)
}
```

### 创建 API Handler：`internal/api/multimodal_handler.go`

```go
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
```

### 在 Router 中注册

在 `internal/api/router.go` 中添加：

```go
// 多模态相关
multimodalGroup := v1.Group("/multimodal")
{
    multimodalHandler := NewMultimodalHandler(defaultChatModel)
    multimodalGroup.POST("/chat", multimodalHandler.ChatWithImage)
}
```

### API 使用示例

```bash
# 使用图像 URL
curl -X POST http://localhost:8080/api/v1/multimodal/chat \
  -H "Content-Type: application/json" \
  -d '{
    "text": "这张图片里有什么？",
    "image_url": "https://example.com/image.jpg"
  }'

# 使用 base64 图像
curl -X POST http://localhost:8080/api/v1/multimodal/chat \
  -H "Content-Type: application/json" \
  -d '{
    "text": "描述这张图片",
    "image_b64": "iVBORw0KGgoAAAANSUhEUgAA...",
    "mime_type": "image/png"
  }'
```

---

## 📝 实施步骤总结

### 1. Embedding 模型（优先级最高）
```bash
# 1. 已创建 ark_embedder.go ✅
# 2. 更新 .env 添加 ARK_EMBEDDING_MODEL
# 3. 更新 config.go 添加配置字段
# 4. 更新 rag_handler.go 的 Index 和 Query 方法
# 5. 更新 router.go 传递参数
# 6. 安装依赖
go get github.com/volcengine/volcengine-go-sdk
```

### 2. 上下文窗口管理
```bash
# 1. 创建 context_manager.go
# 2. 更新 llm_handler.go 使用上下文管理
# 3. 安装依赖
go get github.com/pkoukk/tiktoken-go
```

### 3. 多模态支持
```bash
# 1. 创建 image_handler.go
# 2. 创建 multimodal_handler.go
# 3. 更新 router.go 注册路由
```

---

## 🎯 预期效果

### Embedding 模型
- ✅ RAG 检索准确度大幅提升
- ✅ 使用专业的 1024 维向量
- ✅ 自动降级到简单 embedding

### 上下文窗口管理
- ✅ 自动截断超长对话
- ✅ 避免 token 超限错误
- ✅ 智能保留最新消息

### 多模态支持
- ✅ 支持图像理解
- ✅ 支持 URL 和 base64
- ✅ 统一的 API 接口

---

## 📚 相关文档

- ARK Embedding API: https://www.volcengine.com/docs/82379
- Tiktoken: https://github.com/pkoukk/tiktoken-go
- Eino Multimodal: https://github.com/cloudwego/eino

---

**实施建议：** 按优先级逐个实现，每个功能实现后进行测试，确保稳定后再继续下一个。
