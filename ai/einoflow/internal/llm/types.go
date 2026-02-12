package llm

import (
	"context"
	"fmt"
	"io"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

// LLMProvider 定义 LLM 提供商接口
type LLMProvider interface {
	// GetChatModel 获取聊天模型
	GetChatModel(modelName string, opts ...model.Option) (model.ChatModel, error)

	// ListModels 列出可用模型
	ListModels() []string

	// Name 返回提供商名称
	Name() string
}

// ChatRequest 聊天请求
type ChatRequest struct {
	Model       string             `json:"model"`
	Messages    []*schema.Message  `json:"messages"`
	Stream      bool               `json:"stream"`
	Temperature *float64           `json:"temperature,omitempty"`
	MaxTokens   *int               `json:"max_tokens,omitempty"`
	Tools       []*schema.ToolInfo `json:"tools,omitempty"`
}

// ChatResponse 聊天响应
type ChatResponse struct {
	Message      *schema.Message `json:"message"`
	FinishReason string          `json:"finish_reason"`
	Usage        *Usage          `json:"usage,omitempty"`
}

// Usage token 使用统计
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// Manager LLM 管理器
type Manager struct {
	providers       map[string]LLMProvider
	defaultProvider string
}

// NewManager 创建 LLM 管理器
func NewManager() *Manager {
	return &Manager{
		providers:       make(map[string]LLMProvider),
		defaultProvider: "ark", // 默认使用豆包
	}
}

// SetDefaultProvider 设置默认提供商
func (m *Manager) SetDefaultProvider(name string) {
	m.defaultProvider = name
}

// RegisterProvider 注册提供商
func (m *Manager) RegisterProvider(provider LLMProvider) {
	m.providers[provider.Name()] = provider
}

// GetProvider 获取提供商
func (m *Manager) GetProvider(name string) (LLMProvider, bool) {
	provider, ok := m.providers[name]
	return provider, ok
}

// GetAllProviders 获取所有提供商
func (m *Manager) GetAllProviders() map[string]LLMProvider {
	return m.providers
}

// Chat 执行聊天
func (m *Manager) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	// 根据模型名称选择提供商
	provider := m.selectProvider(req.Model)
	if provider == nil {
		return nil, ErrProviderNotFound
	}

	// 获取聊天模型
	chatModel, err := provider.GetChatModel(req.Model)
	if err != nil {
		return nil, err
	}

	// 构建选项
	opts := []model.Option{}
	if req.Temperature != nil {
		opts = append(opts, model.WithTemperature(float32(*req.Temperature)))
	}
	if req.MaxTokens != nil {
		opts = append(opts, model.WithMaxTokens(*req.MaxTokens))
	}

	// 调用模型
	msg, err := chatModel.Generate(ctx, req.Messages, opts...)
	if err != nil {
		return nil, err
	}

	return &ChatResponse{
		Message:      msg,
		FinishReason: "",  // schema.Message 不包含 FinishReason
		Usage:        nil, // schema.Message 不包含 Usage 信息
	}, nil
}

func (m *Manager) selectProvider(modelName string) LLMProvider {
	// 简单的模型名称匹配逻辑
	for _, provider := range m.providers {
		for _, model := range provider.ListModels() {
			if model == modelName {
				return provider
			}
		}
	}
	// 如果没有匹配到，使用默认提供商
	if defaultProvider, ok := m.providers[m.defaultProvider]; ok {
		return defaultProvider
	}
	return nil
}

// ChatStream 执行流式聊天
func (m *Manager) ChatStream(ctx context.Context, req *ChatRequest) (*schema.StreamReader[*schema.Message], error) {
	// 根据模型名称选择提供商
	provider := m.selectProvider(req.Model)
	if provider == nil {
		return nil, ErrProviderNotFound
	}

	// 获取聊天模型
	chatModel, err := provider.GetChatModel(req.Model)
	if err != nil {
		return nil, err
	}

	// 构建选项
	opts := []model.Option{}
	if req.Temperature != nil {
		opts = append(opts, model.WithTemperature(float32(*req.Temperature)))
	}
	if req.MaxTokens != nil {
		opts = append(opts, model.WithMaxTokens(*req.MaxTokens))
	}

	// 调用流式模型
	return chatModel.Stream(ctx, req.Messages, opts...)
}

// StreamChunk 流式响应块
type StreamChunk struct {
	Content string `json:"content"`
	Done    bool   `json:"done"`
	Error   string `json:"error,omitempty"`
}

// ProcessStream 处理流式响应
func ProcessStream(ctx context.Context, stream *schema.StreamReader[*schema.Message], callback func(*StreamChunk) error) error {
	defer stream.Close()

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			// 流结束
			return callback(&StreamChunk{Done: true})
		}
		if err != nil {
			return callback(&StreamChunk{Error: err.Error(), Done: true})
		}

		// 发送内容块（即使为空也发送，因为流式响应可能有多个块）
		if msg.Content != "" || len(msg.ToolCalls) > 0 {
			if err := callback(&StreamChunk{Content: msg.Content, Done: false}); err != nil {
				return err
			}
		}
	}
}

var (
	ErrProviderNotFound = fmt.Errorf("provider not found")
)
