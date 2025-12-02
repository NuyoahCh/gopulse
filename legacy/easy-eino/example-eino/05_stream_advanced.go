/*
学习目标：
1. 掌握流式数据的高级处理技巧
2. 学会实现流式数据的过滤和转换
3. 理解流式处理的错误处理

高级技巧：
- 流式过滤：只输出符合条件的内容
- 流式转换：实时修改数据格式
- 流式聚合：累积计算统计信息
- 错误恢复：处理流中断和重试

运行方式：
go run 05_stream_advanced.go
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
		schema.SystemMessage("你是一个代码生成助手"),
		schema.UserMessage("请生成一个 Go 语言的 HTTP 服务器示例代码"),
	}

	// 技巧 1：流式过滤（只输出代码块）
	fmt.Println("=== 技巧 1：流式过滤 ===\n")
	streamWithFilter(ctx, chatModel, messages)

	// 技巧 2：流式转换（添加行号）
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("=== 技巧 2：流式转换 ===\n")
	streamWithTransform(ctx, chatModel, messages)

	// 技巧 3：流式聚合（实时统计）
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("=== 技巧 3：流式聚合 ===\n")
	streamWithAggregation(ctx, chatModel, messages)
}

// 技巧 1：流式过滤
func streamWithFilter(ctx context.Context, chatModel *deepseek.ChatModel, messages []*schema.Message) {
	stream, err := chatModel.Stream(ctx, messages)
	if err != nil {
		log.Fatalf("流式生成失败: %v", err)
	}
	defer stream.Close()

	inCodeBlock := false
	fmt.Println("📄 提取的代码:")
	fmt.Println(strings.Repeat("-", 60))

	for {
		chunk, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			log.Fatalf("接收失败: %v", err)
		}

		content := chunk.Content

		// 检测代码块标记
		if strings.Contains(content, "```") {
			inCodeBlock = !inCodeBlock
			continue
		}

		// 只输出代码块内的内容
		if inCodeBlock {
			fmt.Print(content)
		}
	}

	fmt.Println("\n" + strings.Repeat("-", 60))
}

// 技巧 2：流式转换（添加行号）
func streamWithTransform(ctx context.Context, chatModel *deepseek.ChatModel, messages []*schema.Message) {
	stream, err := chatModel.Stream(ctx, messages)
	if err != nil {
		log.Fatalf("流式生成失败: %v", err)
	}
	defer stream.Close()

	lineNumber := 1
	isNewLine := true

	fmt.Println("📝 带行号的输出:")
	fmt.Println(strings.Repeat("-", 60))

	for {
		chunk, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			log.Fatalf("接收失败: %v", err)
		}

		content := chunk.Content

		// 逐字符处理，在每行开头添加行号
		for _, char := range content {
			if isNewLine {
				fmt.Printf("%3d | ", lineNumber)
				isNewLine = false
			}

			fmt.Print(string(char))

			if char == '\n' {
				lineNumber++
				isNewLine = true
			}
		}
	}

	fmt.Println("\n" + strings.Repeat("-", 60))
}

// 技巧 3：流式聚合（实时统计）
func streamWithAggregation(ctx context.Context, chatModel *deepseek.ChatModel, messages []*schema.Message) {
	stream, err := chatModel.Stream(ctx, messages)
	if err != nil {
		log.Fatalf("流式生成失败: %v", err)
	}
	defer stream.Close()

	// 统计信息
	stats := struct {
		TotalChars  int
		TotalChunks int
		CodeLines   int
		TextLines   int
		StartTime   time.Time
		LastUpdate  time.Time
	}{
		StartTime:  time.Now(),
		LastUpdate: time.Now(),
	}

	inCodeBlock := false

	fmt.Println("📊 实时统计:")
	fmt.Println(strings.Repeat("-", 60))

	for {
		chunk, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			log.Fatalf("接收失败: %v", err)
		}

		content := chunk.Content
		stats.TotalChars += len(content)
		stats.TotalChunks++

		// 统计代码行和文本行
		if strings.Contains(content, "```") {
			inCodeBlock = !inCodeBlock
		}

		lines := strings.Count(content, "\n")
		if inCodeBlock {
			stats.CodeLines += lines
		} else {
			stats.TextLines += lines
		}

		// 每 500ms 更新一次统计信息
		if time.Since(stats.LastUpdate) > 500*time.Millisecond {
			elapsed := time.Since(stats.StartTime)
			charsPerSec := float64(stats.TotalChars) / elapsed.Seconds()

			fmt.Printf("\r⏱️  %.1fs | 📝 %d 字符 | 📦 %d 块 | 💻 %d 代码行 | 📄 %d 文本行 | ⚡ %.0f 字符/秒",
				elapsed.Seconds(),
				stats.TotalChars,
				stats.TotalChunks,
				stats.CodeLines,
				stats.TextLines,
				charsPerSec,
			)
			stats.LastUpdate = time.Now()
		}
	}

	// 最终统计
	elapsed := time.Since(stats.StartTime)
	fmt.Printf("\n\n✅ 完成！总耗时: %.2fs\n", elapsed.Seconds())
	fmt.Printf("📊 最终统计:\n")
	fmt.Printf("  - 总字符: %d\n", stats.TotalChars)
	fmt.Printf("  - 数据块: %d\n", stats.TotalChunks)
	fmt.Printf("  - 代码行: %d\n", stats.CodeLines)
	fmt.Printf("  - 文本行: %d\n", stats.TextLines)
	fmt.Printf("  - 平均速度: %.0f 字符/秒\n", float64(stats.TotalChars)/elapsed.Seconds())
}