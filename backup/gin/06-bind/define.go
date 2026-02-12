package main

import "github.com/gin-gonic/gin"

// StructA 定义嵌套结构体
type StructA struct {
	FieldA string `form:"field_a"`
}

// StructB 定义包含 StructA 的结构体
type StructB struct {
	NestedStruct StructA
	FieldB       string `form:"field_b"`
}

// StructC 定义包含 StructA 指针的结构体
type StructC struct {
	NestedStructPointer *StructA
	FieldC              string `form:"field_c"`
}

// StructD 定义包含匿名结构体的结构体
type StructD struct {
	NestedAnonyStruct struct {
		FieldX string `form:"field_x"`
	}
	FieldD string `form:"field_d"`
}

// GetDataB 处理包含嵌套结构体的请求
func GetDataB(c *gin.Context) {
	var b StructB
	c.Bind(&b)
	c.JSON(200, gin.H{
		"a": b.NestedStruct,
		"b": b.FieldB,
	})
}

// GetDataC 处理包含结构体指针的请求
func GetDataC(c *gin.Context) {
	var b StructC
	c.Bind(&b)
	c.JSON(200, gin.H{
		"a": b.NestedStructPointer,
		"c": b.FieldC,
	})
}

// GetDataD 处理包含匿名结构体的请求
func GetDataD(c *gin.Context) {
	var b StructD
	c.Bind(&b)
	c.JSON(200, gin.H{
		"x": b.NestedAnonyStruct,
		"d": b.FieldD,
	})
}

// func main() {
// 	router := gin.Default()
// 	router.GET("/getb", GetDataB)
// 	router.GET("/getc", GetDataC)
// 	router.GET("/getd", GetDataD)

// 	router.Run()
// }
