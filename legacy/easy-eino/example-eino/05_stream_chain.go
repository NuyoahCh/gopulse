/*
学习目标：
1. 掌握在 Chain 中使用流式处理
2. 理解 Eino 的自动流处理机制
3. 学会流式数据的转换和处理

核心概念：
- Chain 的流式支持：runnable.Stream()
- 自动流处理：Eino 自动处理流的传递和转换
- 流式 Lambda：在流中插入自定义处理逻辑

Eino 的流处理魔法：
- 自动装箱：非流 → 流
- 自动拆箱：流 → 非流（当下游需要完整数据时）
- 自动合并：多个流 → 单个流
- 自动复制：单个流 → 多个流（分支时）

运行方式：
go run 05_stream_chain.go
*/

package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

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

	// 场景 1：简单流式 Chain
	fmt.Println("=== 场景 1：Template → Model 流式输出 ===\n")
	simpleStreamChain(ctx, chatModel)

	// 场景 2：带数据处理的流式 Chain
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("=== 场景 2：流式输出 + 实时统计 ===\n")
	streamWithProcessing(ctx, chatModel)
}

// 场景 1：简单的流式 Chain
func simpleStreamChain(ctx context.Context, chatModel *deepseek.ChatModel) {
	// 创建 Chain
	chain := compose.NewChain[map[string]any, *schema.Message]()

	template := prompt.FromMessages(
		schema.FString,
		schema.SystemMessage("你是一个{role}"),
		schema.UserMessage("{question}"),
	)

	chain.
		AppendChatTemplate(template).
		AppendChatModel(chatModel)

	runnable, err := chain.Compile(ctx)
	if err != nil {
		log.Fatalf("编译失败: %v", err)
	}

	input := map[string]any{
		"role":     "诗人",
		"question": "请写一首关于代码之美的小诗",
	}

	// 流式执行 Chain
	stream, err := runnable.Stream(ctx, input)
	if err != nil {
		log.Fatalf("流式执行失败: %v", err)
	}
	defer stream.Close()

	fmt.Print("🎭 诗人创作中: ")

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

	fmt.Println("\n")
}

// 场景 2：流式输出 + 实时统计
func streamWithProcessing(ctx context.Context, chatModel *deepseek.ChatModel) {
	chain := compose.NewChain[string, *schema.Message]()

	// Lambda 1: 构建消息
	chain.AppendLambda(compose.InvokableLambda(func(ctx context.Context, topic string) ([]*schema.Message, error) {
		return []*schema.Message{
			schema.SystemMessage("你是一个技术博主"),
			schema.UserMessage(fmt.Sprintf("请写一篇关于 %s 的技术文章（150字左右）", topic)),
		}, nil
	}))

	// Lambda 2: 调用模型
	chain.AppendChatModel(chatModel)

	runnable, err := chain.Compile(ctx)
	if err != nil {
		log.Fatalf("编译失败: %v", err)
	}

	// 流式执行
	stream, err := runnable.Stream(ctx, "Go 语言的并发模型")
	if err != nil {
		log.Fatalf("流式执行失败: %v", err)
	}
	defer stream.Close()

	fmt.Println("📝 文章生成中...")
	fmt.Println(strings.Repeat("-", 60))

	var (
		fullContent string
		charCount   int
		chunkCount  int
		lineCount   = 1
	)

	for {
		chunk, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			log.Fatalf("接收失败: %v", err)
		}

		// 实时输出
		fmt.Print(chunk.Content)

		// 实时统计
		fullContent += chunk.Content
		charCount += len(chunk.Content)
		chunkCount++
		lineCount += strings.Count(chunk.Content, "\n")
	}

	fmt.Println("\n" + strings.Repeat("-", 60))
	fmt.Printf("📊 统计信息:\n")
	fmt.Printf("  - 总字符数: %d\n", charCount)
	fmt.Printf("  - 数据块数: %d\n", chunkCount)
	fmt.Printf("  - 行数: %d\n", lineCount)
	fmt.Printf("  - 平均块大小: %.1f 字符\n", float64(charCount)/float64(chunkCount))
}
