/*
对比实验：保存 vs 不保存 AI 回复

场景 A：保存 AI 回复（正确）
场景 B：不保存 AI 回复（错误）
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

	chatModel, err := deepseek.NewChatModel(ctx, &deepseek.ChatModelConfig{
		APIKey:  os.Getenv("DEEPSEEK_API_KEY"),
		Model:   "deepseek-chat",
		BaseURL: "https://api.deepseek.com",
	})
	if err != nil {
		log.Fatalf("创建失败: %v", err)
	}

	// fmt.Println("\n" + "="*60)
	fmt.Println("场景 A：保存 AI 回复（正确做法）")
	// fmt.Println("="*60 + "\n")
	scenarioA(ctx, chatModel)

	// fmt.Println("\n" + "="*60)
	fmt.Println("场景 B：不保存 AI 回复（错误做法）")
	// fmt.Println("="*60 + "\n")
	scenarioB(ctx, chatModel)
}

// 场景 A：正确做法 - 保存 AI 回复
func scenarioA(ctx context.Context, chatModel *deepseek.ChatModel) {
	messages := []*schema.Message{
		schema.SystemMessage("你是一个简洁的 AI 助手"),
	}

	// 第一轮
	fmt.Println("[第一轮] 用户: 我的名字是小明")
	messages = append(messages, schema.UserMessage("我的名字是小明"))

	response1, _ := chatModel.Generate(ctx, messages)
	fmt.Printf("[第一轮] AI: %s\n\n", response1.Content)

	// 👇 关键：保存 AI 的回复
	messages = append(messages, response1)

	// 第二轮：追问（需要用到第一轮的信息）
	fmt.Println("[第二轮] 用户: 我刚才说我叫什么名字？")
	messages = append(messages, schema.UserMessage("我刚才说我叫什么名字？"))

	response2, _ := chatModel.Generate(ctx, messages)
	fmt.Printf("[第二轮] AI: %s\n", response2.Content)

	// 打印完整对话历史
	fmt.Println("\n--- 内存中的对话历史 ---")
	for i, msg := range messages {
		fmt.Printf("%d. [%s] %s\n", i+1, msg.Role, truncate(msg.Content, 40))
	}
}

// 场景 B：错误做法 - 不保存 AI 回复
func scenarioB(ctx context.Context, chatModel *deepseek.ChatModel) {
	messages := []*schema.Message{
		schema.SystemMessage("你是一个简洁的 AI 助手"),
	}

	// 第一轮
	fmt.Println("[第一轮] 用户: 我的名字是小明")
	messages = append(messages, schema.UserMessage("我的名字是小明"))

	response1, _ := chatModel.Generate(ctx, messages)
	fmt.Printf("[第一轮] AI: %s\n\n", response1.Content)

	// ❌ 错误：没有保存 AI 的回复！
	// messages = append(messages, response1)  // 这行被注释掉了

	// 第二轮：追问
	fmt.Println("[第二轮] 用户: 我刚才说我叫什么名字？")
	messages = append(messages, schema.UserMessage("我刚才说我叫什么名字？"))

	response2, _ := chatModel.Generate(ctx, messages)
	fmt.Printf("[第二轮] AI: %s\n", response2.Content)

	// 打印对话历史（缺少 AI 的第一轮回复）
	fmt.Println("\n--- 内存中的对话历史 ---")
	for i, msg := range messages {
		fmt.Printf("%d. [%s] %s\n", i+1, msg.Role, truncate(msg.Content, 40))
	}
	fmt.Println("\n⚠️  注意：AI 的第一轮回复丢失了！")
}
