package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Response 定义响应结构体
type Response struct {
	Code int    `json:"code"`
	Data any    `json:"data"`
	Msg  string `json:"msg"`
}

// IndexHandler 处理路径的请求
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	// 输出请求的方法和 URL
	fmt.Println(r.Method, r.URL.String())
	// 仅在 GET 请求时读取并打印请求体内容
	if r.Method == "GET" {
		byteData, _ := io.ReadAll(r.Body)
		fmt.Println(string(byteData))
	}
	fmt.Println(r.Header)
	w.Write([]byte("Hello, World!"))
}

// GET 处理 GET 请求
func GET(res http.ResponseWriter, req *http.Request) {
	// 获取参数
	fmt.Println(req.URL.String())

	byteData, _ := json.Marshal(Response{
		Code: 0,
		Data: map[string]any{},
		Msg:  "成功",
	})
	res.Write(byteData)

}

// POST 处理 POST 请求
func POST(res http.ResponseWriter, req *http.Request) {
	// 获取参数
	byteData, _ := io.ReadAll(req.Body)
	fmt.Println(string(byteData))
	byteData, _ = json.Marshal(Response{
		Code: 0,
		Data: map[string]any{},
		Msg:  "成功",
	})
	res.Write(byteData)
}

func main() {
	// 使用 net/http 包创建一个简单的 HTTP 服务器--基础实现
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("Hello, World!"))
	// })
	http.HandleFunc("/index", IndexHandler)

	// 注册 GET 和 POST 处理函数
	http.HandleFunc("/get", GET)
	http.HandleFunc("/post", POST)

	// 输出服务器启动信息
	fmt.Println("Starting server on :8080")

	// 启动服务器，监听端口 8080
	http.ListenAndServe(":8080", nil)
}
