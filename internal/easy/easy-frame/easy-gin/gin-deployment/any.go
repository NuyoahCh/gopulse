package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.SetTrustedProxies([]string{"192.168.1.2"})

	router.GET("/", func(c *gin.Context) {
		// 如果客户端是 192.168.1.2，则使用 X-Forwarded-For
		// 请求头中可信部分推断出原始客户端 IP。
		// 否则，直接返回客户端 IP
		fmt.Printf("ClientIP: %s\n", c.ClientIP())
	})
	router.Run()
}
