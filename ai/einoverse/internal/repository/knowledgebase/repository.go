package knowledgebase

import (
	"errors"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/Nuyoahch/einoverse/internal/domain/knowledgebase"
	"github.com/google/uuid"
)

// Repository 知识库仓储接口
type Repository interface {
	Create(doc *knowledgebase.Document) error                             // 创建文档
	GetByID(id string) (*knowledgebase.Document, error)                   // 根据ID获取文档
	Search(query string, limit int) ([]knowledgebase.SearchResult, error) // 搜索文档
	Update(id string, doc *knowledgebase.Document) error                  // 更新文档
	Delete(id string) error                                               // 删除文档
	List(offset, limit int) ([]knowledgebase.Document, int, error)        // 列表文档
}

// InMemoryRepository 内存实现的仓储
type InMemoryRepository struct {
	docs  map[string]*knowledgebase.Document // 文档
	mutex sync.RWMutex                       // 互斥锁
}

// NewInMemoryRepository 创建内存仓储
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		docs: make(map[string]*knowledgebase.Document),
	}
}

// Create 创建文档
func (r *InMemoryRepository) Create(doc *knowledgebase.Document) error {
	// 控制并发操作
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// 如果文档ID为空，则生成一个UUID
	if doc.ID == "" {
		doc.ID = uuid.New().String()
	}

	// 如果创建时间为空，则设置为当前时间
	now := time.Now()
	if doc.CreatedAt.IsZero() {
		doc.CreatedAt = now
	}
	// 如果更新时间为空，则设置为当前时间
	doc.UpdatedAt = now
	r.docs[doc.ID] = doc
	return nil
}

// GetByID 根据ID获取文档
func (r *InMemoryRepository) GetByID(id string) (*knowledgebase.Document, error) {
	// 控制并发操作
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// 根据ID获取文档
	doc, exists := r.docs[id]
	if !exists {
		return nil, nil
	}

	// 返回副本
	docCopy := *doc
	return &docCopy, nil
}

// Search 搜索文档
func (r *InMemoryRepository) Search(query string, limit int) ([]knowledgebase.SearchResult, error) {
	// 控制并发操作
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// 初始化搜索结果
	results := make([]knowledgebase.SearchResult, 0)
	queryLower := strings.ToLower(query)

	// 遍历文档
	for _, doc := range r.docs {
		// 计算分数
		score := calculateScore(doc, queryLower)
		// 如果分数大于0，则添加到搜索结果
		if score > 0 {
			results = append(results, knowledgebase.SearchResult{
				Document: *doc,
				Score:    score,
			})
		}
	}

	// 按分数排序
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	// 如果搜索结果超过限制，则截取
	if len(results) > limit {
		results = results[:limit]
	}

	return results, nil
}

// Update 更新文档
func (r *InMemoryRepository) Update(id string, doc *knowledgebase.Document) error {
	// 控制并发操作
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// 如果文档不存在，则返回错误
	if _, exists := r.docs[id]; !exists {
		return errors.New("document not found")
	}

	// 如果更新时间为空，则设置为当前时间
	doc.UpdatedAt = time.Now()
	r.docs[id] = doc
	return nil
}

// Delete 删除文档
func (r *InMemoryRepository) Delete(id string) error {
	// 控制并发操作
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// 如果文档不存在，则返回错误
	if _, exists := r.docs[id]; !exists {
		return errors.New("document not found")
	}

	// 删除文档
	delete(r.docs, id)
	return nil
}

// List 列表文档
func (r *InMemoryRepository) List(offset, limit int) ([]knowledgebase.Document, int, error) {
	// 控制并发操作
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// 计算总数
	total := len(r.docs)
	docs := make([]knowledgebase.Document, 0, limit)

	// 遍历文档
	count := 0
	for _, doc := range r.docs {
		// 如果偏移量大于等于总数，则跳过
		if count < offset {
			count++
			continue
		}
		// 如果文档数量超过限制，则跳出循环
		if len(docs) >= limit {
			break
		}
		// 追加文档内容
		docs = append(docs, *doc)
		count++
	}

	return docs, total, nil
}

// 辅助函数
func calculateScore(doc *knowledgebase.Document, query string) float64 {
	// 初始化分数
	score := 0.0
	// 转化为小写操作
	titleLower := strings.ToLower(doc.Title)
	contentLower := strings.ToLower(doc.Content)
	queryLower := strings.ToLower(query)

	// 支持多关键词匹配：将查询词按空格分割
	queryWords := strings.Fields(queryLower)
	
	// 如果只有一个词或整个查询作为短语匹配
	if len(queryWords) <= 1 {
		// 完整查询匹配（精确匹配）
		if strings.Contains(titleLower, queryLower) {
			score += 10.0
		}
		if strings.Contains(contentLower, queryLower) {
			score += 5.0
		}
		for _, tag := range doc.Tags {
			if strings.Contains(strings.ToLower(tag), queryLower) {
				score += 3.0
			}
		}
	} else {
		// 多关键词匹配：每个关键词单独匹配，累加分数
		titleMatches := 0
		contentMatches := 0
		tagMatches := 0
		
		for _, word := range queryWords {
			if word == "" {
				continue
			}
			if strings.Contains(titleLower, word) {
				titleMatches++
			}
			if strings.Contains(contentLower, word) {
				contentMatches++
			}
			for _, tag := range doc.Tags {
				if strings.Contains(strings.ToLower(tag), word) {
					tagMatches++
					break // 每个标签只匹配一次
				}
			}
		}
		
		// 标题匹配：匹配的关键词越多，分数越高
		if titleMatches > 0 {
			score += 10.0 * float64(titleMatches) / float64(len(queryWords))
		}
		// 内容匹配
		if contentMatches > 0 {
			score += 5.0 * float64(contentMatches) / float64(len(queryWords))
		}
		// 标签匹配
		if tagMatches > 0 {
			score += 3.0 * float64(tagMatches) / float64(len(queryWords))
		}
		
		// 如果所有关键词都匹配，额外加分（短语匹配奖励）
		if titleMatches == len(queryWords) {
			score += 5.0
		}
		if contentMatches == len(queryWords) {
			score += 3.0
		}
	}

	return score
}
