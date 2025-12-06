package providers

import (
	"context"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
)

type OpenAIProvider struct {
	apiKey  string
	baseURL string
}

func NewOpenAIProvider(apiKey, baseURL string) *OpenAIProvider {
	return &OpenAIProvider{
		apiKey:  apiKey,
		baseURL: baseURL,
	}
}

func (p *OpenAIProvider) GetChatModel(modelName string, opts ...model.Option) (model.ChatModel, error) {
	config := &openai.ChatModelConfig{
		APIKey:  p.apiKey,
		BaseURL: p.baseURL,
		Model:   modelName,
	}
	
	return openai.NewChatModel(context.Background(), config)
}

func (p *OpenAIProvider) ListModels() []string {
	return []string{
		"gpt-4",
		"gpt-4-turbo",
		"gpt-4o",
		"gpt-4o-mini",
		"gpt-3.5-turbo",
	}
}

func (p *OpenAIProvider) Name() string {
	return "openai"
}