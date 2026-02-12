package rag

import (
	"context"
	"sort"
	"sync"

	"github.com/cloudwego/eino/schema"
)

// MemoryVectorStore 内存向量存储
type MemoryVectorStore struct {
	mu         sync.RWMutex
	documents  []*schema.Document
	embeddings [][]float64
}

// NewMemoryVectorStore 创建内存向量存储
func NewMemoryVectorStore() *MemoryVectorStore {
	return &MemoryVectorStore{
		documents:  make([]*schema.Document, 0),
		embeddings: make([][]float64, 0),
	}
}

// Add 添加文档和向量
func (s *MemoryVectorStore) Add(ctx context.Context, docs []*schema.Document, embeddings [][]float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.documents = append(s.documents, docs...)
	s.embeddings = append(s.embeddings, embeddings...)

	return nil
}

// Search 搜索相似文档
func (s *MemoryVectorStore) Search(ctx context.Context, queryEmbedding []float64, topK int) ([]*schema.Document, []float64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 计算相似度
	type result struct {
		doc   *schema.Document
		score float64
		index int
	}

	results := make([]result, len(s.documents))
	for i, emb := range s.embeddings {
		score := CosineSimilarity(queryEmbedding, emb)
		results[i] = result{
			doc:   s.documents[i],
			score: score,
			index: i,
		}
	}

	// 排序
	sort.Slice(results, func(i, j int) bool {
		return results[i].score > results[j].score
	})

	// 返回 topK
	if topK > len(results) {
		topK = len(results)
	}

	docs := make([]*schema.Document, topK)
	scores := make([]float64, topK)
	for i := 0; i < topK; i++ {
		docs[i] = results[i].doc
		scores[i] = results[i].score
	}

	return docs, scores, nil
}

// GetAllDocuments 获取所有存储的文档
func (s *MemoryVectorStore) GetAllDocuments() []*schema.Document {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 返回文档的副本
	docs := make([]*schema.Document, len(s.documents))
	copy(docs, s.documents)
	return docs
}

// Count 获取存储的文档数量
func (s *MemoryVectorStore) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.documents)
}

// Clear 清空所有文档
func (s *MemoryVectorStore) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.documents = make([]*schema.Document, 0)
	s.embeddings = make([][]float64, 0)
}
