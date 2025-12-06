package knowledgebase

import (
	"errors"
	"net/http"
	"strconv"

	kbDomain "github.com/Nuyoahch/einoverse/internal/domain/knowledgebase"
	kbService "github.com/Nuyoahch/einoverse/internal/service/knowledgebase"
	appErrors "github.com/Nuyoahch/einoverse/pkg/errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handler 知识库处理器
type Handler struct {
	service *kbService.Service
	logger  *zap.Logger
}

// NewHandler 创建处理器
func NewHandler(service *kbService.Service, logger *zap.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// CreateDocument 创建文档
func (h *Handler) CreateDocument(c *gin.Context) {
	// 绑定请求体
	var doc kbDomain.Document
	if err := c.ShouldBindJSON(&doc); err != nil {
		// 记录请求体错误
		h.logger.Warn("invalid request body for create document", zap.Error(err))
		// 返回错误响应
		h.respondError(c, http.StatusBadRequest, appErrors.ErrInvalidInput, err)
		return
	}

	// 验证必填字段
	if doc.Title == "" || doc.Content == "" {
		// 记录必填字段错误
		h.logger.Warn("missing required fields for create document",
			zap.Bool("has_title", doc.Title != ""),
			zap.Bool("has_content", doc.Content != ""))
		// 返回错误响应
		h.respondError(c, http.StatusBadRequest, appErrors.ErrInvalidInput, nil)
		return
	}

	// 创建文档
	if err := h.service.CreateDocument(&doc); err != nil {
		h.logger.Error("failed to create document", zap.Error(err))
		h.respondError(c, http.StatusInternalServerError, err, err)
		return
	}

	// 返回成功响应
	h.logger.Info("document created successfully", zap.String("id", doc.ID), zap.String("title", doc.Title))
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"id": doc.ID,
		},
	})
}

// GetDocument 获取文档
func (h *Handler) GetDocument(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		h.logger.Warn("empty document id in request")
		h.respondError(c, http.StatusBadRequest, appErrors.ErrInvalidInput, nil)
		return
	}

	doc, err := h.service.GetDocument(id)
	if err != nil {
		if errors.Is(err, appErrors.ErrDocumentNotFound) {
			h.logger.Info("document not found", zap.String("id", id))
			h.respondError(c, http.StatusNotFound, appErrors.ErrDocumentNotFound, nil)
		} else {
			h.logger.Error("failed to get document", zap.String("id", id), zap.Error(err))
			h.respondError(c, http.StatusInternalServerError, err, err)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": doc,
	})
}

// SearchDocuments 搜索文档
func (h *Handler) SearchDocuments(c *gin.Context) {
	query := c.Query("q")
	limitStr := c.DefaultQuery("limit", "5")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		h.logger.Warn("invalid limit parameter", zap.String("limit", limitStr), zap.Error(err))
		limit = 5
	}

	// 限制最大返回数量，防止资源消耗过大
	if limit > 100 {
		limit = 100
	}

	// 如果查询为空，返回空结果
	if query == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "success",
			"data": []kbDomain.SearchResult{},
		})
		return
	}

	// 搜索文档
	results, err := h.service.SearchDocuments(query, limit)
	if err != nil {
		// 记录搜索文档错误
		h.logger.Error("failed to search documents", zap.String("query", query), zap.Error(err))
		// 返回错误响应
		h.respondError(c, http.StatusInternalServerError, err, err)
		return
	}

	h.logger.Info("documents searched", zap.String("query", query), zap.Int("count", len(results)))
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": results,
	})
}

// AskQuestion 问答
func (h *Handler) AskQuestion(c *gin.Context) {
	var req struct {
		Question string `json:"question" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request body for ask question", zap.Error(err))
		h.respondError(c, http.StatusBadRequest, appErrors.ErrInvalidInput, err)
		return
	}

	// 验证问题不为空
	if req.Question == "" {
		h.logger.Warn("empty question in request")
		h.respondError(c, http.StatusBadRequest, appErrors.ErrInvalidInput, nil)
		return
	}

	response, err := h.service.AskQuestion(req.Question)
	if err != nil {
		if errors.Is(err, appErrors.ErrLLMServiceError) {
			h.logger.Error("LLM service error", zap.String("question", req.Question), zap.Error(err))
		} else {
			h.logger.Error("failed to answer question", zap.String("question", req.Question), zap.Error(err))
		}
		h.respondError(c, http.StatusInternalServerError, err, err)
		return
	}

	h.logger.Info("question answered", zap.String("question", req.Question), zap.Float64("confidence", response.Confidence))
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": response,
	})
}

// respondError 统一错误响应处理
func (h *Handler) respondError(c *gin.Context, status int, err error, originalErr error) {
	// 记录原始错误（如果不是业务错误）
	if originalErr != nil {
		var bizErr *appErrors.BusinessError
		if !errors.As(originalErr, &bizErr) {
			h.logger.Error("internal error occurred",
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
				zap.Error(originalErr))
		}
	}

	var bizErr *appErrors.BusinessError
	if errors.As(err, &bizErr) {
		c.JSON(status, gin.H{
			"code": bizErr.Code,
			"msg":  bizErr.Message,
		})
		return
	}

	// 对于非业务错误，不暴露详细错误信息给客户端
	msg := "Internal server error"
	if status == http.StatusBadRequest {
		msg = err.Error()
	}

	c.JSON(status, gin.H{
		"code": "INTERNAL_ERROR",
		"msg":  msg,
	})
}
