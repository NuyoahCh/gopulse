/*
学习目标：
1. 理解 Callback 机制的作用
2. 掌握如何监控 AI 调用过程
3. 学会实现自定义 Callback

核心概念：
- Callback：在特定事件发生时被调用的函数
- 监控：跟踪 AI 的输入、输出、耗时
- 调试：定位问题、优化性能

运行方式：
go run 07_callback_simple.go
*/

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
)

// 自定义 Callback：监控 AI 调用
type MonitorCallback struct {
	callbacks.HandlerBuilder // 提供默认空实现
	callCount                int
}

// OnStart：AI 调用开始时触发
func (m *MonitorCallback) OnStart(ctx context.Context, info *callbacks.RunInfo, input callbacks.CallbackInput) context.Context {
	m.callCount++

	fmt.Printf("\n🚀 [开始] 第 %d 次调用\n", m.callCount)
	fmt.Printf("   组件: %s\n", info.Name)
	fmt.Printf("   类型: %s\n", info.Type)

	// 记录开始时间
	return context.WithValue(ctx, "start_time_"+info.Name, time.Now())
}

// OnEnd：AI 调用结束时触发
func (m *MonitorCallback) OnEnd(ctx context.Context, info *callbacks.RunInfo, output callbacks.CallbackOutput) context.Context {
	// 计算耗时
	if startTime, ok := ctx.Value("start_time_" + info.Name).(time.Time); ok {
		elapsed := time.Since(startTime)
		fmt.Printf("\n✅ [完成] %s - 耗时: %v\n", info.Name, elapsed)
	}

	// 如果是 ChatModel 输出，显示内容
	if msg, ok := output.(*schema.Message); ok {
		content := msg.Content
		if len(content) > 100 {
			content = content[:100] + "..."
		}
		fmt.Printf("   输出: %s\n", content)

		// 显示 Token 统计
		if msg.ResponseMeta != nil && msg.ResponseMeta.Usage != nil {
			usage := msg.ResponseMeta.Usage
			fmt.Printf("   Token: 输入=%d, 输出=%d, 总计=%d\n",
				usage.PromptTokens,
				usage.CompletionTokens,
				usage.PromptTokens+usage.CompletionTokens)
		}
	}

	return ctx
}

// OnStartWithStreamInput：流式输入开始时触发
func (m *MonitorCallback) OnStartWithStreamInput(ctx context.Context, info *callbacks.RunInfo, input *schema.StreamReader[callbacks.CallbackInput]) context.Context {
	m.callCount++
	fmt.Printf("\n🚀 [开始] 第 %d 次调用 (流式输入)\n", m.callCount)
	fmt.Printf("   组件: %s\n", info.Name)
	return context.WithValue(ctx, "start_time_"+info.Name, time.Now())
}

// OnEndWithStreamOutput：流式输出结束时触发
func (m *MonitorCallback) OnEndWithStreamOutput(ctx context.Context, info *callbacks.RunInfo, output *schema.StreamReader[callbacks.CallbackOutput]) context.Context {
	if startTime, ok := ctx.Value("start_time_" + info.Name).(time.Time); ok {
		elapsed := time.Since(startTime)
		fmt.Printf("\n✅ [完成] %s (流式输出) - 耗时: %v\n", info.Name, elapsed)
	}
	return ctx
}

// OnError：发生错误时触发
func (m *MonitorCallback) OnError(ctx context.Context, info *callbacks.RunInfo, err error) context.Context {
	fmt.Printf("\n❌ [错误] %s\n", info.Name)
	fmt.Printf("   错误: %v\n", err)
	return ctx
}

func main() {
	ctx := context.Background()

	// 创建监控 Callback
	monitor := &MonitorCallback{}

	// 创建 ChatModel
	chatModel, err := deepseek.NewChatModel(ctx, &deepseek.ChatModelConfig{
		APIKey:  os.Getenv("DEEPSEEK_API_KEY"),
		Model:   "deepseek-chat",
		BaseURL: "https://api.deepseek.com",
	})
	if err != nil {
		log.Fatalf("创建模型失败: %v", err)
	}

	// 创建 Agent（Agent 支持 Callback）
	ragent, err := react.NewAgent(ctx, &react.AgentConfig{
		ToolCallingModel: chatModel,
		ToolsConfig:      compose.ToolsNodeConfig{},
	})
	if err != nil {
		log.Fatalf("创建 Agent 失败: %v", err)
	}

	fmt.Println("=== Callback 监控演示 ===")
	fmt.Println("观察每次 AI 调用的详细信息\n")
	fmt.Println(strings.Repeat("=", 60))

	// 测试 1：简单问答
	fmt.Println("\n【测试 1：简单问答】")
	messages1 := []*schema.Message{
		schema.UserMessage("什么是 Go 语言？用一句话回答"),
	}

	// 使用 Callback（通过 agent.WithComposeOptions）
	response1, err := ragent.Generate(ctx, messages1,
		agent.WithComposeOptions(compose.WithCallbacks(monitor)))
	if err != nil {
		log.Printf("生成失败: %v", err)
	} else {
		fmt.Printf("\n💬 最终回答: %s\n", response1.Content)
	}

	// 测试 2：复杂问答
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("\n【测试 2：复杂问答】")
	messages2 := []*schema.Message{
		schema.SystemMessage("你是一个技术专家"),
		schema.UserMessage("解释一下 Eino 框架的核心优势（50字以内）"),
	}

	response2, err := ragent.Generate(ctx, messages2,
		agent.WithComposeOptions(compose.WithCallbacks(monitor)))
	if err != nil {
		log.Printf("生成失败: %v", err)
	} else {
		fmt.Printf("\n💬 最终回答: %s\n", response2.Content)
	}

	// 打印统计信息
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Printf("\n📊 统计信息: 总调用次数 = %d\n", monitor.callCount)

	// 思考题：
	// 1. Callback 在什么场景下最有用？
	//    - 日志记录、性能监控、成本统计、调试分析
	// 2. 如何用 Callback 实现成本控制？
	//    - 在 OnEnd 中累计 Token，超过阈值时抛出错误
	// 3. 如何用 Callback 实现请求限流？
	//    - 在 OnStart 中检查请求频率，超过限制时延迟或拒绝
}
