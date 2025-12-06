package embedding

import (
	"context"
	"fmt"

	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
)

// ArkEmbedder 字节豆包 Embedding 实现
type ArkEmbedder struct {
	client *arkruntime.Client
	model  string
}

// NewArkEmbedder 创建字节豆包 Embedder
func NewArkEmbedder(apiKey, baseURL, modelID string) (*ArkEmbedder, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("ARK API key is required")
	}
	if modelID == "" {
		modelID = "doubao-embedding-large-text-250515" // 默认模型
	}

	client := arkruntime.NewClientWithApiKey(
		apiKey,
		arkruntime.WithBaseUrl(baseURL),
		arkruntime.WithRegion("cn-beijing"),
	)

	return &ArkEmbedder{
		client: client,
		model:  modelID,
	}, nil
}

// EmbedText 对单个文本进行向量化
func (e *ArkEmbedder) EmbedText(ctx context.Context, text string) ([]float64, error) {
	req := &model.EmbeddingRequest{
		Model: e.model,
		Input: text,
	}

	resp, err := e.client.CreateEmbeddings(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create embedding: %w", err)
	}

	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("no embedding data returned")
	}

	// 转换为 []float64
	embedding := make([]float64, len(resp.Data[0].Embedding))
	for i, v := range resp.Data[0].Embedding {
		embedding[i] = float64(v)
	}

	return embedding, nil
}

// EmbedTexts 批量对文本进行向量化
func (e *ArkEmbedder) EmbedTexts(ctx context.Context, texts []string) ([][]float64, error) {
	embeddings := make([][]float64, len(texts))

	// 批量处理（ARK 支持批量）
	for i := 0; i < len(texts); i += 10 { // 每次处理10个
		end := i + 10
		if end > len(texts) {
			end = len(texts)
		}

		batch := texts[i:end]
		for j, text := range batch {
			embedding, err := e.EmbedText(ctx, text)
			if err != nil {
				return nil, fmt.Errorf("failed to embed text %d: %w", i+j, err)
			}
			embeddings[i+j] = embedding
		}
	}

	return embeddings, nil
}

// GetDimension 获取向量维度
func (e *ArkEmbedder) GetDimension() int {
	// doubao-embedding-large-text-250515 的维度是 1024
	return 1024
}

// GetModel 获取模型名称
func (e *ArkEmbedder) GetModel() string {
	return e.model
}
