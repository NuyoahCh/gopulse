package middleware

import (
	"time"

	"einoflow/pkg/logger"

	"github.com/gin-gonic/gin"
)

// Logger 结构化日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// 处理请求
		c.Next()

		// 计算延迟
		latency := time.Since(start)

		// 获取状态码
		statusCode := c.Writer.Status()

		// 获取 Request ID
		requestID := GetRequestID(c)

		// 构建日志字段
		fields := map[string]interface{}{
			"request_id": requestID,
			"method":     c.Request.Method,
			"path":       path,
			"query":      query,
			"status":     statusCode,
			"latency_ms": latency.Milliseconds(),
			"client_ip":  c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		}

		// 如果有错误，添加错误信息
		if len(c.Errors) > 0 {
			fields["errors"] = c.Errors.String()
		}

		// 根据状态码选择日志级别
		if statusCode >= 500 {
			logger.WithFields(fields).Error("Server error")
		} else if statusCode >= 400 {
			logger.WithFields(fields).Warn("Client error")
		} else {
			logger.WithFields(fields).Info("Request completed")
		}
	}
}
