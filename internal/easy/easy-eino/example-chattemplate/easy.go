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

	// 1. 创建 ChatTemplate
	template := prompt.FromMessages(
		schema.FString, // 使用 FString 格式化（类似 Python 的 f-string）
		schema.SystemMessage("你是一个{role}"),
		schema.UserMessage("{question}"),
	)

	// 2. 准备变量
	variables := map[string]any{
		"role":     "专业的 Go 语言工程师",
		"question": "什么是 interface？",
	}

	// 3. 格式化消息
	messages, err := template.Format(ctx, variables)
	if err != nil {
		log.Fatalf("格式化失败: %v", err)
	}

	// 4. 查看生成的消息
	fmt.Println("生成的消息:")
	for i, msg := range messages {
		fmt.Printf("%d. [%s] %s\\n", i+1, msg.Role, msg.Content)
	}

	// 5. 使用生成的消息调用模型
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

	fmt.Printf("\\nAI 回答:\\n%s\\n", response.Content)
}
