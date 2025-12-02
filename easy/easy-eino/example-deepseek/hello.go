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
	// 1. 创建上下文
	ctx := context.Background()

	// 2. 创建 ChatModel (使用 DeepSeek)
	chatModel, err := deepseek.NewChatModel(ctx, &deepseek.ChatModelConfig{
		APIKey:  os.Getenv("DEEPSEEK_API_KEY"),
		Model:   "deepseek-chat",
		BaseURL: "https://api.deepseek.com",
	})
	if err != nil {
		log.Fatalf("创建 ChatModel 失败: %v", err)
	}

	// 3. 准备消息
	messages := []*schema.Message{
		schema.SystemMessage("你是一个友好的 AI 助手"),
		schema.UserMessage("你好，请介绍一下 Eino 框架"),
	}

	// 4. 调用模型生成响应
	response, err := chatModel.Generate(ctx, messages)
	if err != nil {
		log.Fatalf("生成响应失败: %v", err)
	}

	// 5. 输出结果
	fmt.Printf("AI 响应: %s\\n", response.Content)

	// 6. 输出 token 使用情况
	if response.ResponseMeta != nil && response.ResponseMeta.Usage != nil {
		fmt.Printf("\\nToken 使用统计:\\n")
		fmt.Printf("  输入 Token: %d\\n", response.ResponseMeta.Usage.PromptTokens)
		fmt.Printf("  输出 Token: %d\\n", response.ResponseMeta.Usage.CompletionTokens)
		fmt.Printf("  总计 Token: %d\\n", response.ResponseMeta.Usage.TotalTokens)
	}
}
