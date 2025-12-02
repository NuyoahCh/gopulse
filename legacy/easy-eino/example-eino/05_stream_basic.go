/*
学习目标：
1. 理解流式处理的概念和价值
2. 掌握 Stream 方法的基本用法
3. 学会正确处理流的生命周期

核心概念：
- 流式处理：数据逐块返回，而非等待全部完成
- StreamReader：流读取器接口
- Recv()：接收下一个数据块
- io.EOF：流结束标志

为什么需要流式处理？
- 用户体验：实时看到输出（打字机效果）
- 性能优化：边生成边处理，降低首字延迟
- 长文本生成：避免超时，及时反馈进度

对比：
- Generate()：等待完整结果（适合短文本）
- Stream()：逐块返回（适合长文本、实时交互）

运行方式：
go run 05_stream_basic.go
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
	"time"

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

	messages := []*schema.Message{
		schema.SystemMessage("你是一个专业的技术作家"),
		schema.UserMessage("请写一篇关于 Eino 框架流式处理的文章（200字左右）"),
	}

	fmt.Println("=== 对比演示 ===\n")

	// 演示 1：非流式（Generate）
	fmt.Println("【方式1：非流式 Generate】")
	fmt.Println("等待中...")
	startTime := time.Now()

	response, err := chatModel.Generate(ctx, messages)
	if err != nil {
		log.Fatalf("生成失败: %v", err)
	}

	elapsed := time.Since(startTime)
	fmt.Printf("\n⏱️  耗时: %v\n", elapsed)
	fmt.Printf("📝 结果:\n%s\n\n", response.Content)

	// 演示 2：流式（Stream）
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("【方式2：流式 Stream】")
	fmt.Print("AI 正在输出: ")

	startTime = time.Now()

	// 步骤 1：调用 Stream 方法获取流
	stream, err := chatModel.Stream(ctx, messages)
	if err != nil {
		log.Fatalf("流式生成失败: %v", err)
	}
	// 步骤 2：确保关闭流（重要！）
	defer stream.Close()

	var fullContent string
	chunkCount := 0

	// 步骤 3：循环接收数据块
	for {
		// 步骤 4：接收下一个块
		chunk, err := stream.Recv()
		if err != nil {
			// 步骤 5：检查是否流结束
			if errors.Is(err, io.EOF) {
				break // 正常结束
			}
			log.Fatalf("接收失败: %v", err)
		}

		// 步骤 6：处理数据块
		fmt.Print(chunk.Content) // 实时打印（打字机效果）
		fullContent += chunk.Content
		chunkCount++

		// 模拟打字机效果（可选）
		time.Sleep(10 * time.Millisecond)
	}

	elapsed = time.Since(startTime)
	fmt.Printf("\n\n⏱️  耗时: %v\n", elapsed)
	fmt.Printf("📊 统计: 共接收 %d 个数据块\n", chunkCount)
	fmt.Printf("📝 完整内容长度: %d 字符\n", len(fullContent))

	// 思考题：
	// 1. 如果忘记 defer stream.Close() 会有什么后果？
	// 2. 流式和非流式的总耗时差别大吗？为什么？
	// 3. 如何在流式处理中实现进度条？
}
