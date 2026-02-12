package rag

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sort"
	"sync"

	"github.com/cloudwego/eino/schema"
	_ "github.com/mattn/go-sqlite3"
)

// PersistentVectorStore 持久化向量存储（基于 SQLite）
type PersistentVectorStore struct {
	mu sync.RWMutex
	db *sql.DB
}

// NewPersistentVectorStore 创建持久化向量存储
func NewPersistentVectorStore(dbPath string) (*PersistentVectorStore, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// 创建表
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS documents (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		content TEXT NOT NULL,
		metadata TEXT,
		embedding TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_created_at ON documents(created_at);
	`

	if _, err := db.Exec(createTableSQL); err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return &PersistentVectorStore{db: db}, nil
}

// Add 添加文档和向量
func (s *PersistentVectorStore) Add(ctx context.Context, docs []*schema.Document, embeddings [][]float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, "INSERT INTO documents (content, metadata, embedding) VALUES (?, ?, ?)")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for i, doc := range docs {
		// 序列化 metadata
		metadataJSON, err := json.Marshal(doc.MetaData)
		if err != nil {
			return fmt.Errorf("failed to marshal metadata: %w", err)
		}

		// 序列化 embedding
		embeddingJSON, err := json.Marshal(embeddings[i])
		if err != nil {
			return fmt.Errorf("failed to marshal embedding: %w", err)
		}

		if _, err := stmt.ExecContext(ctx, doc.Content, string(metadataJSON), string(embeddingJSON)); err != nil {
			return fmt.Errorf("failed to insert document: %w", err)
		}
	}

	return tx.Commit()
}

// Search 搜索相似文档
func (s *PersistentVectorStore) Search(ctx context.Context, queryEmbedding []float64, topK int) ([]*schema.Document, []float64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 查询所有文档
	rows, err := s.db.QueryContext(ctx, "SELECT id, content, metadata, embedding FROM documents")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to query documents: %w", err)
	}
	defer rows.Close()

	type result struct {
		doc   *schema.Document
		score float64
	}

	var results []result

	for rows.Next() {
		var id int
		var content, metadataJSON, embeddingJSON string

		if err := rows.Scan(&id, &content, &metadataJSON, &embeddingJSON); err != nil {
			return nil, nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// 反序列化 metadata
		var metadata map[string]any
		if err := json.Unmarshal([]byte(metadataJSON), &metadata); err != nil {
			metadata = make(map[string]any)
		}

		// 反序列化 embedding
		var embedding []float64
		if err := json.Unmarshal([]byte(embeddingJSON), &embedding); err != nil {
			continue // 跳过无效的 embedding
		}

		// 计算相似度
		score := CosineSimilarity(queryEmbedding, embedding)

		results = append(results, result{
			doc: &schema.Document{
				Content:  content,
				MetaData: metadata,
			},
			score: score,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("error iterating rows: %w", err)
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
func (s *PersistentVectorStore) GetAllDocuments() ([]*schema.Document, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rows, err := s.db.Query("SELECT content, metadata FROM documents ORDER BY created_at DESC")
	if err != nil {
		return nil, fmt.Errorf("failed to query documents: %w", err)
	}
	defer rows.Close()

	var docs []*schema.Document

	for rows.Next() {
		var content, metadataJSON string
		if err := rows.Scan(&content, &metadataJSON); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		var metadata map[string]any
		if err := json.Unmarshal([]byte(metadataJSON), &metadata); err != nil {
			metadata = make(map[string]any)
		}

		docs = append(docs, &schema.Document{
			Content:  content,
			MetaData: metadata,
		})
	}

	return docs, nil
}

// Count 获取存储的文档数量
func (s *PersistentVectorStore) Count() (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM documents").Scan(&count)
	return count, err
}

// Clear 清空所有文档
func (s *PersistentVectorStore) Clear() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec("DELETE FROM documents")
	return err
}

// Close 关闭数据库连接
func (s *PersistentVectorStore) Close() error {
	return s.db.Close()
}
