package api

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"einoflow/internal/embedding"
	"einoflow/internal/rag"
	"einoflow/pkg/logger"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"github.com/gin-gonic/gin"
)

type RAGHandler struct {
	chatModel        model.ChatModel
	vectorStore      *rag.MemoryVectorStore     // 内存存储（临时）
	persistentStore  *rag.PersistentVectorStore // 持久化存储
	embedder         *embedding.ArkEmbedder     // Embedding 模型
	usePersistent    bool
	useRealEmbedding bool
}

func NewRAGHandler(chatModel model.ChatModel, arkAPIKey, arkBaseURL, embeddingModel string) *RAGHandler {
	// 尝试创建持久化存储
	persistentStore, err := rag.NewPersistentVectorStore("./data/vector_store.db")
	usePersistent := err == nil

	if usePersistent {
		logger.Info("Using persistent vector store (SQLite)")
	} else {
		logger.Info("Using memory vector store (data will be lost on restart)")
	}

	// 尝试创建 Embedding 模型
	var embedder *embedding.ArkEmbedder
	useRealEmbedding := false
	if arkAPIKey != "" && embeddingModel != "" {
		embedder, err = embedding.NewArkEmbedder(arkAPIKey, arkBaseURL, embeddingModel)
		if err == nil {
			useRealEmbedding = true
			logger.Info(fmt.Sprintf("Using ARK Embedding model: %s", embeddingModel))
		} else {
			logger.Warn(fmt.Sprintf("Failed to create ARK embedder: %v, using simple embedding", err))
		}
	} else {
		logger.Info("Using simple character-based embedding")
	}

	return &RAGHandler{
		chatModel:        chatModel,
		vectorStore:      rag.NewMemoryVectorStore(),
		persistentStore:  persistentStore,
		embedder:         embedder,
		usePersistent:    usePersistent,
		useRealEmbedding: useRealEmbedding,
	}
}

type RAGIndexRequest struct {
	Documents []string `json:"documents" binding:"required"`
}

type RAGQueryRequest struct {
	Query string `json:"query" binding:"required"`
}

type RAGQueryResponse struct {
	Answer    string   `json:"answer"`
	Documents []string `json:"documents,omitempty"` // 相关文档
}

type RAGStatsResponse struct {
	Count     int      `json:"count"`
	Documents []string `json:"documents"`
}

type RAGUploadResponse struct {
	Message       string `json:"message"`
	Filename      string `json:"filename"`
	DocumentCount int    `json:"document_count"`
	TotalCount    int    `json:"total_count"`
}

func (h *RAGHandler) Index(c *gin.Context) {
	var req RAGIndexRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 将文档转换为 schema.Document
	docs := make([]*schema.Document, len(req.Documents))
	embeddings := make([][]float64, len(req.Documents))

	for i, content := range req.Documents {
		docs[i] = &schema.Document{
			Content: content,
			MetaData: map[string]any{
				"index": i,
				"id":    fmt.Sprintf("doc_%d", i),
			},
		}

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
	}

	// 添加到向量存储
	var err error
	var totalCount int

	if h.usePersistent {
		err = h.persistentStore.Add(c.Request.Context(), docs, embeddings)
		if err == nil {
			totalCount, _ = h.persistentStore.Count()
		}
	} else {
		err = h.vectorStore.Add(c.Request.Context(), docs, embeddings)
		totalCount = h.vectorStore.Count()
	}

	if err != nil {
		logger.Error("Failed to add documents: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to index documents"})
		return
	}

	logger.Info(fmt.Sprintf("Indexed %d documents, total: %d", len(req.Documents), totalCount))

	c.JSON(http.StatusOK, gin.H{
		"message": "Documents indexed successfully",
		"count":   len(req.Documents),
		"total":   h.vectorStore.Count(),
	})
}

func (h *RAGHandler) Query(c *gin.Context) {
	var req RAGQueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检索相关文档
	var contextDocs []string
	var hasDocuments bool

	if h.usePersistent {
		count, _ := h.persistentStore.Count()
		hasDocuments = count > 0
	} else {
		hasDocuments = h.vectorStore.Count() > 0
	}

	if hasDocuments {
		// 使用真实 Embedding 或简单 Embedding
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

		var docs []*schema.Document
		var scores []float64
		var err error

		if h.usePersistent {
			docs, scores, err = h.persistentStore.Search(c.Request.Context(), queryEmbedding, 3)
		} else {
			docs, scores, err = h.vectorStore.Search(c.Request.Context(), queryEmbedding, 3)
		}

		if err == nil && len(docs) > 0 {
			contextDocs = make([]string, len(docs))
			for i, doc := range docs {
				contextDocs[i] = doc.Content
				logger.Info(fmt.Sprintf("Retrieved doc %d (score: %.3f): %s", i, scores[i], doc.Content[:min(50, len(doc.Content))]))
			}
		}
	}

	// 构建上下文
	var systemPrompt string
	if len(contextDocs) > 0 {
		systemPrompt = fmt.Sprintf("以下是相关的背景信息：\n\n%s\n\n请基于以上信息回答问题。",
			joinStrings(contextDocs, "\n\n"))
	} else {
		systemPrompt = "请回答以下问题："
	}

	messages := []*schema.Message{
		schema.SystemMessage(systemPrompt),
		schema.UserMessage(req.Query),
	}

	resp, err := h.chatModel.Generate(c.Request.Context(), messages)
	if err != nil {
		logger.Error("RAG query failed: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, &RAGQueryResponse{
		Answer:    resp.Content,
		Documents: contextDocs,
	})
}

// GetStats 获取 RAG 统计信息
func (h *RAGHandler) GetStats(c *gin.Context) {
	var docs []*schema.Document
	var err error

	if h.usePersistent {
		docs, err = h.persistentStore.GetAllDocuments()
	} else {
		docs = h.vectorStore.GetAllDocuments()
	}

	if err != nil {
		logger.Error("Failed to get documents: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get documents"})
		return
	}

	docContents := make([]string, len(docs))
	for i, doc := range docs {
		docContents[i] = doc.Content
	}

	c.JSON(http.StatusOK, &RAGStatsResponse{
		Count:     len(docs),
		Documents: docContents,
	})
}

// Clear 清空所有文档
func (h *RAGHandler) Clear(c *gin.Context) {
	var err error

	if h.usePersistent {
		err = h.persistentStore.Clear()
	} else {
		h.vectorStore.Clear()
	}

	if err != nil {
		logger.Error("Failed to clear documents: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to clear documents"})
		return
	}

	logger.Info("Cleared all documents from vector store")
	c.JSON(http.StatusOK, gin.H{
		"message": "All documents cleared",
	})
}

// simpleEmbedding 简单的向量化方法（基于字符特征）
func simpleEmbedding(text string) []float64 {
	// 创建一个 128 维的向量
	vec := make([]float64, 128)

	// 基于字符特征生成向量
	for i, char := range text {
		idx := i % 128
		vec[idx] += float64(char) / 1000.0
	}

	// 添加长度特征
	vec[0] += float64(len(text)) / 100.0

	return vec
}

func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// UploadFile 上传文件并索引
func (h *RAGHandler) UploadFile(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no file uploaded"})
		return
	}

	// 检查文件大小（限制 10MB）
	if file.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file size exceeds 10MB limit"})
		return
	}

	// 打开文件
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open file"})
		return
	}
	defer src.Close()

	// 读取文件内容
	content, err := io.ReadAll(src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})
		return
	}

	// 将文件内容转换为文本
	text := string(content)

	// 分块处理文本（每 500 字符一块）
	chunks := h.splitText(text, 500)

	if len(chunks) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is empty or cannot be processed"})
		return
	}

	// 创建文档和 embeddings
	docs := make([]*schema.Document, len(chunks))
	embeddings := make([][]float64, len(chunks))

	for i, chunk := range chunks {
		docs[i] = &schema.Document{
			Content: chunk,
			MetaData: map[string]any{
				"index":    i,
				"filename": file.Filename,
				"source":   "upload",
			},
		}

		// 使用真实 Embedding 或简单 Embedding
		if h.useRealEmbedding {
			embedding, err := h.embedder.EmbedText(c.Request.Context(), chunk)
			if err != nil {
				logger.Warn(fmt.Sprintf("Failed to embed text, using simple embedding: %v", err))
				embeddings[i] = simpleEmbedding(chunk)
			} else {
				embeddings[i] = embedding
			}
		} else {
			embeddings[i] = simpleEmbedding(chunk)
		}
	}

	// 添加到向量存储
	var totalCount int
	if h.usePersistent {
		err = h.persistentStore.Add(c.Request.Context(), docs, embeddings)
		if err == nil {
			totalCount, _ = h.persistentStore.Count()
		}
	} else {
		err = h.vectorStore.Add(c.Request.Context(), docs, embeddings)
		totalCount = h.vectorStore.Count()
	}

	if err != nil {
		logger.Error("Failed to add documents: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to index file"})
		return
	}

	logger.Info(fmt.Sprintf("Indexed file %s with %d chunks, total: %d", file.Filename, len(chunks), totalCount))

	c.JSON(http.StatusOK, &RAGUploadResponse{
		Message:       "File uploaded and indexed successfully",
		Filename:      file.Filename,
		DocumentCount: len(chunks),
		TotalCount:    totalCount,
	})
}

// splitText 分割文本为块
func (h *RAGHandler) splitText(text string, chunkSize int) []string {
	if len(text) == 0 {
		return nil
	}

	var chunks []string
	runes := []rune(text)

	for i := 0; i < len(runes); i += chunkSize {
		end := i + chunkSize
		if end > len(runes) {
			end = len(runes)
		}

		chunk := string(runes[i:end])
		// 去除空白块
		if len(strings.TrimSpace(chunk)) > 0 {
			chunks = append(chunks, chunk)
		}
	}

	return chunks
}
