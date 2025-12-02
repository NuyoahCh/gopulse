/*
🎯 综合实战项目：智能助手系统

整合所有学习内容：
1. ChatModel - AI 对话能力
2. ChatTemplate - 提示词模板
3. Chain - 流程编排
4. Stream - 流式输出
5. Tool - 工具调用
6. Agent - 智能决策
7. Callback - 监控统计

功能：
- 支持多轮对话
- 自动调用工具
- 流式输出
- 性能监控
- Token 统计

运行方式：
go run 08_final_project.go
*/

package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"time"

	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
)

// ========== Callback：性能监控 ==========
type PerformanceMonitor struct {
	callbacks.HandlerBuilder
	totalCalls  int
	totalTokens int
	totalTime   time.Duration
	startTime   time.Time
}

func (p *PerformanceMonitor) OnStart(ctx context.Context, info *callbacks.RunInfo, input callbacks.CallbackInput) context.Context {
	p.totalCalls++
	p.startTime = time.Now()
	return ctx
}

func (p *PerformanceMonitor) OnEnd(ctx context.Context, info *callbacks.RunInfo, output callbacks.CallbackOutput) context.Context {
	elapsed := time.Since(p.startTime)
	p.totalTime += elapsed

	if msg, ok := output.(*schema.Message); ok {
		if msg.ResponseMeta != nil && msg.ResponseMeta.Usage != nil {
			p.totalTokens += msg.ResponseMeta.Usage.PromptTokens + msg.ResponseMeta.Usage.CompletionTokens
		}
	}
	return ctx
}

func (p *PerformanceMonitor) OnStartWithStreamInput(ctx context.Context, info *callbacks.RunInfo, input *schema.StreamReader[callbacks.CallbackInput]) context.Context {
	return ctx
}

func (p *PerformanceMonitor) OnEndWithStreamOutput(ctx context.Context, info *callbacks.RunInfo, output *schema.StreamReader[callbacks.CallbackOutput]) context.Context {
	return ctx
}

func (p *PerformanceMonitor) OnError(ctx context.Context, info *callbacks.RunInfo, err error) context.Context {
	return ctx
}

func (p *PerformanceMonitor) GetStats() string {
	avgTime := time.Duration(0)
	if p.totalCalls > 0 {
		avgTime = p.totalTime / time.Duration(p.totalCalls)
	}
	return fmt.Sprintf("调用次数: %d | Token: %d | 总耗时: %v | 平均耗时: %v",
		p.totalCalls, p.totalTokens, p.totalTime, avgTime)
}

// ========== 工具集 ==========

// 创建所有工具
func createAllTools() []tool.BaseTool {
	// 工具 1：获取当前时间
	timeTool := utils.NewTool(
		&schema.ToolInfo{
			Name:        "get_time",
			Desc:        "获取当前时间",
			ParamsOneOf: nil,
		},
		func(ctx context.Context, params map[string]any) (string, error) {
			now := time.Now().Format("2006-01-02 15:04:05")
			fmt.Printf("  🔧 [工具] 获取时间 → %s\n", now)
			return now, nil
		},
	)

	// 工具 2：计算器
	calculator := utils.NewTool(
		&schema.ToolInfo{
			Name: "calculator",
			Desc: "执行数学计算，支持加减乘除和幂运算",
			ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
				"operation": {
					Type:     "string",
					Desc:     "运算类型: add(加), sub(减), mul(乘), div(除), pow(幂)",
					Required: true,
				},
				"a": {
					Type:     "number",
					Desc:     "第一个数字",
					Required: true,
				},
				"b": {
					Type:     "number",
					Desc:     "第二个数字",
					Required: true,
				},
			}),
		},
		func(ctx context.Context, params map[string]any) (string, error) {
			op := params["operation"].(string)
			a := params["a"].(float64)
			b := params["b"].(float64)

			var result float64
			switch op {
			case "add":
				result = a + b
			case "sub":
				result = a - b
			case "mul":
				result = a * b
			case "div":
				if b == 0 {
					return "错误：除数不能为0", nil
				}
				result = a / b
			case "pow":
				result = math.Pow(a, b)
			default:
				return "错误：不支持的运算", nil
			}

			resultStr := fmt.Sprintf("%.2f", result)
			fmt.Printf("  🔧 [工具] 计算 %g %s %g = %s\n", a, op, b, resultStr)
			return resultStr, nil
		},
	)

	// 工具 3：天气查询
	weatherTool := utils.NewTool(
		&schema.ToolInfo{
			Name: "get_weather",
			Desc: "查询指定城市的天气",
			ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
				"city": {
					Type:     "string",
					Desc:     "城市名称",
					Required: true,
				},
			}),
		},
		func(ctx context.Context, params map[string]any) (string, error) {
			city := params["city"].(string)
			// 模拟天气数据
			result := fmt.Sprintf("%s：晴天，温度 22°C，湿度 60%%，空气质量优", city)
			fmt.Printf("  🔧 [工具] 查询天气 %s → %s\n", city, result)
			return result, nil
		},
	)

	return []tool.BaseTool{timeTool, calculator, weatherTool}
}

// ========== 智能助手 ==========

type SmartAssistant struct {
	agent   *react.Agent
	monitor *PerformanceMonitor
	history []*schema.Message
}

func NewSmartAssistant(ctx context.Context) (*SmartAssistant, error) {
	// 创建 ChatModel
	chatModel, err := deepseek.NewChatModel(ctx, &deepseek.ChatModelConfig{
		APIKey:  os.Getenv("DEEPSEEK_API_KEY"),
		Model:   "deepseek-chat",
		BaseURL: "https://api.deepseek.com",
	})
	if err != nil {
		return nil, err
	}

	// 创建工具
	tools := createAllTools()

	// 创建 Agent
	ragent, err := react.NewAgent(ctx, &react.AgentConfig{
		ToolCallingModel: chatModel,
		ToolsConfig: compose.ToolsNodeConfig{
			Tools: tools,
		},
		MaxStep: 5,
	})
	if err != nil {
		return nil, err
	}

	// 创建监控器
	monitor := &PerformanceMonitor{}

	return &SmartAssistant{
		agent:   ragent,
		monitor: monitor,
		history: []*schema.Message{
			schema.SystemMessage("你是一个智能助手，可以帮助用户查询时间、计算数学问题、查询天气等。请简洁明了地回答问题。"),
		},
	}, nil
}

// 对话
func (sa *SmartAssistant) Chat(ctx context.Context, userInput string) (string, error) {
	// 添加用户消息
	sa.history = append(sa.history, schema.UserMessage(userInput))

	// 调用 Agent
	response, err := sa.agent.Generate(ctx, sa.history,
		agent.WithComposeOptions(compose.WithCallbacks(sa.monitor)))
	if err != nil {
		return "", err
	}

	// 保存 AI 回复
	sa.history = append(sa.history, response)

	return response.Content, nil
}

// 获取统计信息
func (sa *SmartAssistant) GetStats() string {
	return sa.monitor.GetStats()
}

// ========== 主程序 ==========

func main() {
	ctx := context.Background()

	// 创建智能助手
	assistant, err := NewSmartAssistant(ctx)
	if err != nil {
		log.Fatalf("创建助手失败: %v", err)
	}

	fmt.Println("╔═══════════════════════════════════════════════════════╗")
	fmt.Println("║          🤖 智能助手系统 - Eino 综合实战项目          ║")
	fmt.Println("╚═══════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("功能：")
	fmt.Println("  ✅ 多轮对话")
	fmt.Println("  ✅ 自动调用工具（时间、计算器、天气）")
	fmt.Println("  ✅ 性能监控")
	fmt.Println()
	fmt.Println(strings.Repeat("=", 60))

	// 测试场景
	testCases := []struct {
		name  string
		input string
	}{
		{"时间查询", "现在几点了？"},
		{"数学计算", "帮我计算 15 的平方"},
		{"天气查询", "北京今天天气怎么样？"},
		{"复杂问题", "如果我现在出发去北京，天气适合吗？"},
	}

	for i, tc := range testCases {
		fmt.Printf("\n【场景 %d：%s】\n", i+1, tc.name)
		fmt.Printf("👤 用户: %s\n\n", tc.input)

		response, err := assistant.Chat(ctx, tc.input)
		if err != nil {
			log.Printf("对话失败: %v", err)
			continue
		}

		fmt.Printf("🤖 助手: %s\n", response)
		fmt.Println(strings.Repeat("-", 60))
	}

	// 显示统计信息
	fmt.Println()
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("\n📊 性能统计: %s\n", assistant.GetStats())
	fmt.Println()
	fmt.Println("✨ 恭喜！你已经掌握了 Eino 框架的所有核心知识！")
	fmt.Println()
	fmt.Println("下一步建议：")
	fmt.Println("  1. 阅读 Eino 官方文档深入学习")
	fmt.Println("  2. 尝试构建自己的 AI 应用")
	fmt.Println("  3. 探索更多高级特性（RAG、多 Agent 等）")
	fmt.Println()
}
