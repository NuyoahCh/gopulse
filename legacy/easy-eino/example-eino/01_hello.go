/*
学习目标：
1. 理解 Eino 的核心概念：Component（组件）
2. 掌握 ChatModel 的基本用法
3. 了解 schema.Message 的作用

核心概念：
- Component：Eino 的基本构建块，每个组件都有明确的输入/输出类型
- ChatModel：对话模型的抽象接口，屏蔽了底层实现细节
- schema.Message：Eino 的统一消息格式

运行方式：
export DEEPSEEK_API_KEY="your-key"
go run 01_hello.go
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

	// 步骤 1：创建 ChatModel 组件
	// ChatModel 是 Eino 的核心组件之一，代表一个对话模型
	// 这里使用 DeepSeek 作为实现，但你也可以替换为 Ollama、ARK 等
	chatModel, err := deepseek.NewChatModel(ctx, &deepseek.ChatModelConfig{
		APIKey:  os.Getenv("DEEPSEEK_API_KEY"), // 从环境变量读取 API Key
		Model:   "deepseek-chat",               // 指定模型名称
		BaseURL: "https://api.deepseek.com",    // API 基础 URL
	})
	if err != nil {
		log.Fatalf("创建 ChatModel 失败: %v", err)
	}

	// 步骤 2：构建消息
	// schema.Message 是 Eino 的统一消息格式，有三种类型：
	// - SystemMessage：系统角色设定，定义 AI 的行为模式
	// - UserMessage：用户输入
	// - AssistantMessage：AI 的回复
	messages := []*schema.Message{
		schema.SystemMessage("你是一个专业的 Go 语言工程师，擅长用简洁的语言解释技术概念。"),
		schema.UserMessage("请用一句话解释什么是 Eino 框架？"),
	}

	// 步骤 3：调用 ChatModel 生成回复
	// Generate 方法是同步调用，会等待模型生成完整的回复
	response, err := chatModel.Generate(ctx, messages)
	if err != nil {
		log.Fatalf("生成回复失败: %v", err)
	}

	// 步骤 4：输出结果
	fmt.Printf("\n=== AI 回答 ===\n%s\n", response.Content)

	// 思考题：
	// 1. 为什么需要 SystemMessage？如果去掉会有什么影响？
	// 2. ChatModel 还有哪些方法？（提示：查看 Stream 方法）
	// 3. 如果要替换为 Ollama 模型，需要改动哪些代码？

	// 练习：
	// 1. 修改 SystemMessage，让 AI 扮演不同的角色（如数学老师、诗人）
	// 2. 尝试问不同的问题，观察回答的变化
	// 3. 查看 response 对象还有哪些字段（提示：response.Role）
}
