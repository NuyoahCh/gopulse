package providers

import (
	"context"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/components/model"
)

// ArkProvider 字节豆包 Provider
type ArkProvider struct {
	apiKey  string
	baseURL string
}

// NewArkProvider 创建字节豆包 Provider
func NewArkProvider(apiKey, baseURL string) *ArkProvider {
	return &ArkProvider{
		apiKey:  apiKey,
		baseURL: baseURL,
	}
}

// GetChatModel 获取聊天模型
func (p *ArkProvider) GetChatModel(modelName string, opts ...model.Option) (model.ChatModel, error) {
	config := &ark.ChatModelConfig{
		APIKey: p.apiKey,
		Model:  modelName,
	}

	if p.baseURL != "" {
		config.BaseURL = p.baseURL
	}

	return ark.NewChatModel(context.Background(), config)
}

// ListModels 列出可用模型
func (p *ArkProvider) ListModels() []string {
	return []string{
		"doubao-seed-1-6-lite-251015",
		"doubao-seed-1-6-vision-250815",
	}
}

// Name 返回提供商名称
func (p *ArkProvider) Name() string {
	return "ark"
}
