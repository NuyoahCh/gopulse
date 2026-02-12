package main

import "github.com/gin-gonic/gin"

// 控制器函数
func doLogin(c *gin.Context) {
	// 获取post请求参数
	username := c.PostForm("username")
	password := c.PostForm("password")

	// 通过请求上下文对象Context, 直接往客户端返回一个字符串
	c.String(200, "username=%s,password=%s", username, password)
}

func main() {
	r := gin.Default()

	// 路由定义post请求, url路径为：/user/login, 绑定doLogin控制器函数
	r.POST("/user/login", doLogin)
}
