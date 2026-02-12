package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// LoginRequest 定义登录请求的结构体
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func main() {
	g := gin.New()
	g.POST("/login", func(ctx *gin.Context) {
		r := &LoginRequest{}
		// 绑定请求体到结构体
		ctx.ShouldBind(r)
		fmt.Printf("login-request:%+v\n", r)
	})

	g.Run(":9090")
}
