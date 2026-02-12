package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/query", Handler)

	_ = r.Run(":8080")
}

func Handler(c *gin.Context) {
	// 获取name参数, 通过Query获取的参数值是String类型。
	name := c.Query("name")

	// 获取name参数, 跟Query函数的区别是，可以通过第二个参数设置默认值。
	nameWithDefault := c.DefaultQuery("name", "shanyangsuanfa")

	// 获取id参数, 通过GetQuery获取的参数值也是String类型,
	// 区别是GetQuery返回两个参数，第一个是参数值，第二个参数是参数是否存在的bool值，可以用来判断参数是否存在。
	id, ok := c.GetQuery("id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "missing id",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"name":              name,
		"name_with_default": nameWithDefault,
		"id":                id,
	})
}
