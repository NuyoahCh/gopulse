package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	// 加载templates目录下面的所有模版文件，包括子目录
	// **/* 代表所有子目录下的所有文件
	router.LoadHTMLGlob("templates/**/*")

	router.GET("/posts/index", func(c *gin.Context) {
		// 子目录的模版文件，需要加上目录名，例如：posts/index.tmpl
		c.HTML(http.StatusOK, "posts/index.tmpl", gin.H{
			"title": "Posts",
		})
	})

	router.GET("/users/index", func(c *gin.Context) {
		// 子目录的模版文件，需要加上目录名，例如：users/index.tmpl
		c.HTML(http.StatusOK, "users/index.tmpl", gin.H{
			"title": "Users",
		})
	})
	router.Run(":8080")
}
