package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func LatencyMonitor() gin.HandlerFunc {
	return func(c *gin.Context) {
		// [阶段 1: 入栈] 记录起始时间
		start := time.Now()
		path := c.Request.URL.Path

		// [阶段 2: 移交] 执行后续业务逻辑
		c.Next()

		// [阶段 3: 出栈] 业务逻辑已完成，此时可获取状态码和耗时
		latency := time.Since(start)
		status := c.Writer.Status()

		// 仅记录核心指标，避免在中间件中进行耗时 I/O
		log.Printf("[Metrics] | %d | %13v | %s", status, latency, path)
	}
}
