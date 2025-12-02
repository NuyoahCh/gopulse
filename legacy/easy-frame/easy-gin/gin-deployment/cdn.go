package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	// 使用预定义的 gin.PlatformXXX 头
	// Google App Engine
	router.TrustedPlatform = gin.PlatformGoogleAppEngine
	// Cloudflare
	router.TrustedPlatform = gin.PlatformCloudflare
	// Fly.io
	router.TrustedPlatform = gin.PlatformFlyIO
	// 或者，你可以设置自己的可信请求头。但要确保你的 CDN
	// 能防止用户传递此请求头！例如，如果你的 CDN 将客户端
	// IP 放在 X-CDN-Client-IP 中：
	router.TrustedPlatform = "X-CDN-Client-IP"

	router.GET("/", func(c *gin.Context) {
		// 如果设置了 TrustedPlatform，ClientIP() 将解析
		// 对应的请求头并直接返回 IP
		fmt.Printf("ClientIP: %s\n", c.ClientIP())
	})
	router.Run()
}
