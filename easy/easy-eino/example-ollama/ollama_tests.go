package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino/schema"
)

func main() {
	ctx := context.Background()

	// Ollama 默认运行在 http://localhost:11434
	chatModel, err := ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
		BaseURL: "http://localhost:11434",
		Model:   "llama2", // 或其他已下载的模型
	})
	if err != nil {
		log.Fatalf("创建失败: %v", err)
	}

	messages := []*schema.Message{
		schema.UserMessage("golang是什么?"),
	}

	stream, err := chatModel.Stream(ctx, messages)
	if err != nil {
		log.Fatalf("流式生成失败: %v", err)
	}
	defer stream.Close()

	fmt.Print("AI 回复: ")

	for {
		chunk, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			log.Fatalf("接收失败: %v", err)
		}

		// 实时打印
		fmt.Print(chunk.Content)
	}
}
