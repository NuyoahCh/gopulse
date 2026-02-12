package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 1. 初始化 Redis 存储后端
	store, _ := redis.NewStore(10, "tcp", "localhost:6379", "", "secret-key")

	// 2. 注入中间件（全局应用）
	r.Use(sessions.Sessions("GOSESSID", store))

	r.GET("/login", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set("uid", 10086) // 写入会话
		session.Save()            // 必须手动保存
		c.JSON(200, gin.H{"status": "logged in"})
	})
}
