package rag

import (
	"context"

	"github.com/cloudwego/eino/components/embedding"
	"github.com/cloudwego/eino/schema"
)

// VectorRetriever 向量检索器
type VectorRetriever struct {
	embedder    embedding.Embedder
	vectorStore VectorStore
	topK        int
}

// VectorStore 向量存储接口
type VectorStore interface {
	Add(ctx context.Context, docs []*schema.Document, embeddings [][]float64) error
	Search(ctx context.Context, queryEmbedding []float64, topK int) ([]*schema.Document, []float64, error)
}

// NewVectorRetriever 创建向量检索器
func NewVectorRetriever(embedder embedding.Embedder, vectorStore VectorStore, topK int) *VectorRetriever {
	return &VectorRetriever{
		embedder:    embedder,
		vectorStore: vectorStore,
		topK:        topK,
	}
}

// Retrieve 检索文档
func (r *VectorRetriever) Retrieve(ctx context.Context, query string) ([]*schema.Document, error) {
	// 1. 对查询进行向量化
	embeddings, err := r.embedder.EmbedStrings(ctx, []string{query})
	if err != nil {
		return nil, err
	}
	queryEmbedding := embeddings[0]

	// 2. 在向量存储中搜索
	docs, _, err := r.vectorStore.Search(ctx, queryEmbedding, r.topK)
	if err != nil {
		return nil, err
	}

	return docs, nil
}
