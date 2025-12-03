package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/compose"
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
		log.Fatalf("创建模型失败: %v", err)
	}

	// 创建支持流式的 Chain
	chain := compose.NewChain[map[string]any, *schema.Message]()

	template := prompt.FromMessages(
		schema.FString,
		schema.SystemMessage("你是一个故事作家"),
		schema.UserMessage("请写一个关于{topic}的短故事"),
	)

	chain.
		AppendChatTemplate(template).
		AppendChatModel(chatModel)

	runnable, err := chain.Compile(ctx)
	if err != nil {
		log.Fatalf("编译失败: %v", err)
	}

	input := map[string]any{
		"topic": "勇敢的程序员",
	}

	// 流式执行
	stream, err := runnable.Stream(ctx, input)
	if err != nil {
		log.Fatalf("流式执行失败: %v", err)
	}
	defer stream.Close()

	fmt.Print("AI 正在创作: ")

	for {
		chunk, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			log.Fatalf("接收失败: %v", err)
		}

		fmt.Print(chunk.Content)
	}

	fmt.Println("\\n\\n创作完成！")
}
