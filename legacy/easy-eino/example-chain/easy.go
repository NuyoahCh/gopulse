package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

func main() {
	ctx := context.Background()

	// 1. 创建 ChatTemplate
	chatTemplate := prompt.FromMessages(
		schema.FString,
		schema.SystemMessage("你是一个{role}"),
		schema.UserMessage("{question}"),
	)

	// 2. 创建 ChatModel
	chatModel, err := deepseek.NewChatModel(ctx, &deepseek.ChatModelConfig{
		APIKey:  os.Getenv("DEEPSEEK_API_KEY"),
		Model:   "deepseek-chat",
		BaseURL: "https://api.deepseek.com",
	})
	if err != nil {
		log.Fatalf("创建模型失败: %v", err)
	}

	// 3. 创建 Chain：ChatTemplate → ChatModel
	// 输入：map[string]any → 输出：*schema.Message
	chain := compose.NewChain[map[string]any, *schema.Message]()
	chain.
		AppendChatTemplate(chatTemplate). // 第一步：格式化模板
		AppendChatModel(chatModel)        // 第二步：调用模型

	// 4. 编译 Chain
	runnable, err := chain.Compile(ctx)
	if err != nil {
		log.Fatalf("编译失败: %v", err)
	}

	// 5. 执行 Chain
	input := map[string]any{
		"role":     "Go 语言专家",
		"question": "什么是 goroutine？",
	}

	output, err := runnable.Invoke(ctx, input)
	if err != nil {
		log.Fatalf("执行失败: %v", err)
	}

	fmt.Printf("回答: %s\\n", output.Content)
}
