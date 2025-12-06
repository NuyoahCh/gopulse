package leave

import (
	"net/http"

	"github.com/Nuyoahch/einoverse/internal/domain/leave"
	leaveService "github.com/Nuyoahch/einoverse/internal/service/leave"
	"github.com/Nuyoahch/einoverse/pkg/errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handler 请假处理器
type Handler struct {
	service *leaveService.Service
	logger  *zap.Logger
}

// NewHandler 创建处理器
func NewHandler(service *leaveService.Service, logger *zap.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// CreateApplication 创建请假申请
func (h *Handler) CreateApplication(c *gin.Context) {
	var req leave.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.respondError(c, http.StatusBadRequest, errors.ErrInvalidInput)
		return
	}

	app, err := h.service.CreateApplication(&req)
	if err != nil {
		h.respondError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": app,
	})
}

// GetApplication 获取申请
func (h *Handler) GetApplication(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		h.respondError(c, http.StatusBadRequest, errors.ErrInvalidInput)
		return
	}

	app, err := h.service.GetApplication(id)
	if err != nil {
		if err == errors.ErrApplicationNotFound {
			h.respondError(c, http.StatusNotFound, err)
		} else {
			h.respondError(c, http.StatusInternalServerError, err)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": app,
	})
}

// ApproveApplication 审批申请
func (h *Handler) ApproveApplication(c *gin.Context) {
	// 从URL路径参数获取申请ID
	id := c.Param("id")
	if id == "" {
		h.respondError(c, http.StatusBadRequest, errors.ErrInvalidInput)
		return
	}

	var req leave.ApprovalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.respondError(c, http.StatusBadRequest, errors.ErrInvalidInput)
		return
	}

	// 使用URL路径中的ID，忽略body中的application_id（如果提供）
	recommendation, err := h.service.ApproveApplication(id, req.Approver)
	if err != nil {
		h.respondError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": recommendation,
	})
}

// respondError 统一错误响应处理
func (h *Handler) respondError(c *gin.Context, status int, err error) {
	if bizErr, ok := err.(*errors.BusinessError); ok {
		c.JSON(status, gin.H{
			"code": bizErr.Code,
			"msg":  bizErr.Message,
		})
	} else {
		c.JSON(status, gin.H{
			"code": "INTERNAL_ERROR",
			"msg":  err.Error(),
		})
	}
}
