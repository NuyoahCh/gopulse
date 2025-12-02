package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/schema"
)

func main() {
	ctx := context.Background()

	// 使用 ARK ChatModel
	chatModel, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey: os.Getenv("ARK_API_KEY"),
		Model:  os.Getenv("ARK_MODEL_NAME"), // 例如: ep-xxxxx
	})
	if err != nil {
		log.Fatalf("创建 ChatModel 失败: %v", err)
	}

	messages := []*schema.Message{
		schema.SystemMessage("你是一个友好的 AI 助手"),
		schema.UserMessage("你好，请用介绍一下字节跳动大模型Eino"),
	}

	response, err := chatModel.Generate(ctx, messages)
	if err != nil {
		log.Fatalf("生成响应失败: %v", err)
	}

	fmt.Printf("AI 响应: %s\\n", response.Content)
}
