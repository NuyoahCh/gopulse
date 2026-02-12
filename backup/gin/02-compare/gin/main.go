package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 定义响应结构体
type Response struct {
	Code int    `json:"code"`
	Data any    `json:"data"`
	Msg  string `json:"msg"`
}

// IndexHandler 处理路径的请求
func IndexHandler(c *gin.Context) {
	// Gin 优势：直接返回字符串，无需手动处理 ResponseWriter
	c.String(http.StatusOK, "Hello, World!")
}

// GET 处理 GET 请求
func GET(c *gin.Context) {
	// Gin 优势：c.JSON() 自动序列化并设置正确的 Content-Type
	c.JSON(http.StatusOK, Response{
		Code: 0,
		Data: map[string]any{},
		Msg:  "成功",
	})
}

// POST 处理 POST 请求
func POST(c *gin.Context) {
	// Gin 优势：自动 JSON 绑定和验证
	var body map[string]any
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code: -1,
			Data: nil,
			Msg:  "参数错误",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code: 0,
		Data: body,
		Msg:  "成功",
	})
}

func main() {
	// Gin 优势：链式路由注册，代码更简洁
	r := gin.Default()

	r.GET("/index", IndexHandler)
	r.GET("/get", GET)
	r.POST("/post", POST)

	// Gin 优势：Run() 自动启动服务并打印日志
	r.Run(":8080")
}
