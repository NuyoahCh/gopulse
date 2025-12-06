package knowledgebase

import (
	"time"
)

// Document 知识库文档
type Document struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Tags      []string  `json:"tags"`
	Author    string    `json:"author,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SearchResult 搜索结果
type SearchResult struct {
	Document Document `json:"document"`
	Score    float64  `json:"score"`
}

// QAResponse 问答响应
type QAResponse struct {
	Answer     string   `json:"answer"`
	SourceDocs []string `json:"source_docs,omitempty"`
	Confidence float64  `json:"confidence,omitempty"`
}
