package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LoggerMiddleware 日志中间件
func LoggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	// 日志中间件，记录请求信息
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// 继续执行下一个中间件
		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		logger.Info("HTTP request",
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", status),
			zap.Duration("latency", latency),
			zap.String("client_ip", c.ClientIP()),
		)
	}
}

// RecoveryMiddleware 恢复中间件
func RecoveryMiddleware(logger *zap.Logger) gin.HandlerFunc {
	// 恢复中间件，捕获panic并记录错误
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		// 记录panic信息
		logger.Error("panic recovered",
			zap.Any("error", recovered),
			zap.String("path", c.Request.URL.Path),
		)
		// 返回500错误响应
		c.JSON(500, gin.H{
			"code": "INTERNAL_ERROR",
			"msg":  "Internal server error",
		})
		// 终止请求
		c.Abort()
	})
}
