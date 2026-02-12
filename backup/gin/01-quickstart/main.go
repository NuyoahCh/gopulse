package main

import "github.com/gin-gonic/gin"

func main() {
	// 1 初始化
	r := gin.Default()
	// 2 路由
	r.GET("/", func(c *gin.Context) {
		c.String(200, "hello world")
	})
	// 3 监听运行
	r.Run(":8080")
}
