package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/schema"
)

func main() {
	ctx := context.Background()

	// 创建 OpenAI 聊天模型
	// 从环境变量读取 API Key
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is required")
	}

	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		APIKey: apiKey,
		Model:  "gpt-4o-mini",
	})
	if err != nil {
		log.Fatal(err)
	}

	// 准备消息
	messages := []*schema.Message{
		schema.UserMessage("你好，请介绍一下 Eino 框架"),
	}

	// 调用模型
	resp, err := chatModel.Generate(ctx, messages)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("回复: %s\n", resp.Content)
	fmt.Println("调用成功")
}
