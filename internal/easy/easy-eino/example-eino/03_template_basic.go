/*
学习目标：
1. 理解 ChatTemplate 的作用：分离提示词和业务逻辑
2. 掌握变量插值：使用 {variable} 语法
3. 了解 schema.FString 格式化器

核心概念：
- ChatTemplate：提示词模板，类似于 HTML 模板
- 变量插值：动态替换模板中的占位符
- 可复用性：同一个模板可以用于不同的输入

为什么需要 ChatTemplate？
- 问题：硬编码提示词难以维护和复用
- 解决：将提示词模板化，参数化

运行方式：
go run 03_template_basic.go
*/

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

func main() {
	ctx := context.Background()

	// 步骤 1：创建 ChatTemplate
	// prompt.FromMessages 创建一个消息模板
	// schema.FString 是格式化器，支持 {variable} 语法
	template := prompt.FromMessages(
		schema.FString, // 格式化器类型
		schema.SystemMessage("你是一个{role}，擅长{skill}"),
		schema.UserMessage("{question}"),
	)

	fmt.Println("=== 示例 1：基础用法 ===\n")

	// 步骤 2：准备变量
	variables := map[string]any{
		"role":     "专业的 Go 语言工程师",
		"skill":    "用简洁的语言解释复杂概念",
		"question": "什么是 interface？",
	}

	// 步骤 3：格式化消息
	// Format 方法会将变量替换到模板中
	messages, err := template.Format(ctx, variables)
	if err != nil {
		log.Fatalf("格式化失败: %v", err)
	}

	// 步骤 4：查看生成的消息
	fmt.Println("生成的消息:")
	for i, msg := range messages {
		fmt.Printf("%d. [%s] %s\n", i+1, msg.Role, msg.Content)
	}

	// 步骤 5：使用生成的消息调用模型
	chatModel, err := deepseek.NewChatModel(ctx, &deepseek.ChatModelConfig{
		APIKey:  os.Getenv("DEEPSEEK_API_KEY"),
		Model:   "deepseek-chat",
		BaseURL: "https://api.deepseek.com",
	})
	if err != nil {
		log.Fatalf("创建模型失败: %v", err)
	}

	response, err := chatModel.Generate(ctx, messages)
	if err != nil {
		log.Fatalf("生成失败: %v", err)
	}

	fmt.Printf("\nAI 回答:\n%s\n\n", response.Content)

	// 示例 2：复用同一个模板，不同的变量
	fmt.Println("=== 示例 2：模板复用 ===\n")

	variables2 := map[string]any{
		"role":     "资深的 Python 开发者",
		"skill":    "对比不同编程语言的特性",
		"question": "Python 的装饰器和 Go 的中间件有什么相似之处？",
	}

	messages2, err := template.Format(ctx, variables2)
	if err != nil {
		log.Fatalf("格式化失败: %v", err)
	}

	response2, err := chatModel.Generate(ctx, messages2)
	if err != nil {
		log.Fatalf("生成失败: %v", err)
	}

	fmt.Printf("AI 回答:\n%s\n", response2.Content)

	// 思考题：
	// 1. 如果变量名写错了（如 {rol} 而不是 {role}），会发生什么？
	// 2. 能否在模板中使用多次同一个变量？
	// 3. 除了 FString，还有哪些格式化器？（提示：查看 schema 包）
}
