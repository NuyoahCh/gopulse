/*
学习目标：
1. 理解多轮对话的上下文管理
2. 掌握消息历史的维护方法
3. 了解 LLM 的无状态特性

核心概念：
- 对话上下文：LLM 是无状态的，需要每次传入完整历史
- 消息追加：使用 append 维护对话历史
- AssistantMessage：保存 AI 的回复到历史中

运行方式：
go run 02_multi_turn.go
*/

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino/schema"
)

func main() {
	ctx := context.Background()

	// 创建 ChatModel
	chatModel, err := deepseek.NewChatModel(ctx, &deepseek.ChatModelConfig{
		APIKey:  os.Getenv("DEEPSEEK_API_KEY"),
		Model:   "deepseek-chat",
		BaseURL: "https://api.deepseek.com",
	})
	if err != nil {
		log.Fatalf("创建失败: %v", err)
	}

	// 初始化对话历史（包含系统消息）
	messages := []*schema.Message{
		schema.SystemMessage("你是一个 Go 语言专家，擅长用简洁的语言解释复杂概念。"),
	}

	// 第一轮对话
	fmt.Println("=== 第一轮对话 ===")
	messages = append(messages, schema.UserMessage("请用一句话解释 goroutine"))

	response1, err := chatModel.Generate(ctx, messages)
	if err != nil {
		log.Fatalf("第一轮失败: %v", err)
	}
	fmt.Printf("AI: %s\n\n", response1.Content)

	// 关键：将 AI 的回复添加到历史中
	messages = append(messages, response1)

	// 第二轮对话（追问细节）
	fmt.Println("=== 第二轮对话 ===")
	messages = append(messages, schema.UserMessage("那它和操作系统线程有什么区别？"))

	response2, err := chatModel.Generate(ctx, messages)
	if err != nil {
		log.Fatalf("第二轮失败: %v", err)
	}
	fmt.Printf("AI: %s\n\n", response2.Content)

	// 第三轮对话（继续追问）
	messages = append(messages, response2)
	fmt.Println("=== 第三轮对话 ===")
	messages = append(messages, schema.UserMessage("请给一个简单的代码示例"))

	response3, err := chatModel.Generate(ctx, messages)
	if err != nil {
		log.Fatalf("第三轮失败: %v", err)
	}
	fmt.Printf("AI: %s\n\n", response3.Content)

	// 打印完整对话历史
	fmt.Println("=== 完整对话历史 ===")
	for i, msg := range messages {
		fmt.Printf("%d. [%s] %s\n", i+1, msg.Role, truncate(msg.Content, 50))
	}

	// 思考题：
	// 1. 如果不保存 AI 的回复，第二轮对话会发生什么？
	// 2. 对话历史越来越长，会有什么问题？如何优化？
	// 3. 如何实现一个交互式的命令行对话程序？
}

// 辅助函数：截断长文本
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
