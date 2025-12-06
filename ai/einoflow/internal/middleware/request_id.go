package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const RequestIDKey = "X-Request-ID"

// RequestID 中间件为每个请求生成唯一 ID
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 尝试从请求头获取 Request ID
		requestID := c.GetHeader(RequestIDKey)

		// 如果没有，生成一个新的
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// 设置到上下文和响应头
		c.Set(RequestIDKey, requestID)
		c.Header(RequestIDKey, requestID)

		c.Next()
	}
}

// GetRequestID 从上下文获取 Request ID
func GetRequestID(c *gin.Context) string {
	if requestID, exists := c.Get(RequestIDKey); exists {
		return requestID.(string)
	}
	return ""
}
