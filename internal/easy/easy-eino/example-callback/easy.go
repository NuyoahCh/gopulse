package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

// LoggerCallback 日志回调
type LoggerCallback struct {
	callbacks.HandlerBuilder // 辅助实现，提供默认空实现
}

func (cb *LoggerCallback) OnStartWithStreamInput(ctx context.Context, info *callbacks.RunInfo, input *schema.StreamReader[callbacks.CallbackInput]) context.Context {
	fmt.Printf("\\n[开始] %s (流式输入)\\n", info.Name)
	return ctx
}

func (cb *LoggerCallback) OnEndWithStreamOutput(ctx context.Context, info *callbacks.RunInfo, output *schema.StreamReader[callbacks.CallbackOutput]) context.Context {
	fmt.Printf("\\n[完成] %s (流式输出)\\n", info.Name)
	return ctx
}

func (cb *LoggerCallback) OnStart(ctx context.Context, info *callbacks.RunInfo, input callbacks.CallbackInput) context.Context {
	fmt.Printf("\\n[开始] %s\\n", info.Name)
	fmt.Printf("  类型: %s\\n", info.Type)

	// 格式化输入
	inputJSON, _ := json.MarshalIndent(input, "  ", "  ")
	fmt.Printf("  输入: %s\\n", string(inputJSON))

	// 记录开始时间
	return context.WithValue(ctx, "start_time", time.Now())
}

func (cb *LoggerCallback) OnEnd(ctx context.Context, info *callbacks.RunInfo, output callbacks.CallbackOutput) context.Context {
	fmt.Printf("\\n[完成] %s\\n", info.Name)

	// 计算耗时
	if startTime, ok := ctx.Value("start_time").(time.Time); ok {
		duration := time.Since(startTime)
		fmt.Printf("  耗时: %v\\n", duration)
	}

	// 格式化输出（只显示部分，避免太长）
	outputJSON, _ := json.Marshal(output)
	if len(outputJSON) > 200 {
		fmt.Printf("  输出: %s...\\n", string(outputJSON[:200]))
	} else {
		fmt.Printf("  输出: %s\\n", string(outputJSON))
	}

	return ctx
}

func (cb *LoggerCallback) OnError(ctx context.Context, info *callbacks.RunInfo, err error) context.Context {
	fmt.Printf("\\n[错误] %s\\n", info.Name)
	fmt.Printf("  错误: %v\\n", err)
	return ctx
}

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

	// 编译时注册 Callback
	runnable, err := chain.Compile(ctx, compose.WithEagerExecution())
	if err != nil {
		log.Fatalf("编译失败: %v", err)
	}

	// 执行
	input := map[string]any{
		"role":     "Go 专家",
		"question": "什么是 goroutine？",
	}

	response, err := runnable.Invoke(ctx, input)
	if err != nil {
		log.Fatalf("执行失败: %v", err)
	}

	fmt.Printf("\\n\\n=== 最终结果 ===\\n%s\\n", response.Content)
}
