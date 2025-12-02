/*
学习目标：
1. 理解 ReAct Agent 的工作原理
2. 掌握 Agent 的创建和使用
3. 学会观察 Agent 的推理过程

核心概念：
- ReAct：Reasoning（推理）+ Acting（行动）
- Agent：能够自主决策和行动的 AI 系统
- 推理循环：思考 → 行动 → 观察 → 思考...

ReAct 流程：
1. 用户提问
2. Agent 推理：需要什么信息？
3. Agent 行动：调用工具获取信息
4. Agent 观察：分析工具返回的结果
5. Agent 回答：基于信息给出答案

运行方式：
go run 06_react_agent.go
*/

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
)

func main() {
	ctx := context.Background()

	// 1. 创建 ChatModel
	chatModel, err := deepseek.NewChatModel(ctx, &deepseek.ChatModelConfig{
		APIKey:  os.Getenv("DEEPSEEK_API_KEY"),
		Model:   "deepseek-chat",
		BaseURL: "https://api.deepseek.com",
	})
	if err != nil {
		log.Fatalf("创建模型失败: %v", err)
	}

	// 2. 创建工具集
	tools := createToolsForAgent()

	// 3. 创建 ReAct Agent
	agent, err := react.NewAgent(ctx, &react.AgentConfig{
		ToolCallingModel: chatModel,
		ToolsConfig: compose.ToolsNodeConfig{
			Tools: tools,
		},
		MaxStep: 5, // 最多推理 5 步
	})
	if err != nil {
		log.Fatalf("创建 Agent 失败: %v", err)
	}

	// 4. 测试 Agent
	testCases := []string{
		"现在几点了？明天这个时候是几点？",
		"帮我计算 (100 + 50) * 2",
		"北京和上海今天的天气对比",
	}

	for i, question := range testCases {
		// fmt.Printf("\n%s\n", "="*60)
		fmt.Printf("测试 %d: %s\n", i+1, question)
		// fmt.Printf("%s\n\n", "="*60)

		messages := []*schema.Message{
			schema.UserMessage(question),
		}

		response, err := agent.Generate(ctx, messages)
		if err != nil {
			log.Printf("执行失败: %v", err)
			continue
		}

		fmt.Printf("\n✅ Agent 最终回答:\n%s\n", response.Content)
	}

	// 思考题：
	// 1. Agent 如何决定调用哪些工具？
	// 2. MaxStep 的作用是什么？
	// 3. 如何让 Agent 更智能地推理？
}

func createToolsForAgent() []tool.BaseTool {
	// 工具 1：获取当前时间
	timeTool := utils.NewTool(
		&schema.ToolInfo{
			Name:        "get_current_time",
			Desc:        "获取当前的日期和时间，返回格式：YYYY-MM-DD HH:MM:SS",
			ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{}),
		},
		func(ctx context.Context, params map[string]any) (string, error) {
			now := time.Now()
			result := now.Format("2006-01-02 15:04:05")
			fmt.Printf("\n🔧 [工具] get_current_time\n")
			fmt.Printf("   返回: %s\n", result)
			return result, nil
		},
	)

	// 工具 2：时间计算
	timeCalculator := utils.NewTool(
		&schema.ToolInfo{
			Name: "calculate_time",
			Desc: "计算未来或过去的时间",
			ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
				"hours": {
					Type:     "number",
					Desc:     "小时数（正数表示未来，负数表示过去）",
					Required: true,
				},
			}),
		},
		func(ctx context.Context, params map[string]any) (string, error) {
			hours := int(params["hours"].(float64))
			future := time.Now().Add(time.Duration(hours) * time.Hour)
			result := future.Format("2006-01-02 15:04:05")
			fmt.Printf("\n🔧 [工具] calculate_time\n")
			fmt.Printf("   参数: %d 小时后\n", hours)
			fmt.Printf("   返回: %s\n", result)
			return result, nil
		},
	)

	// 工具 3：计算器
	calculator := utils.NewTool(
		&schema.ToolInfo{
			Name: "calculator",
			Desc: "执行数学计算，支持加减乘除和括号",
			ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
				"expression": {
					Type:     "string",
					Desc:     "数学表达式",
					Required: true,
				},
			}),
		},
		func(ctx context.Context, params map[string]any) (string, error) {
			expr := params["expression"].(string)
			// 简化：实际应该用表达式解析
			result := "300" // 假设结果
			fmt.Printf("\n🔧 [工具] calculator\n")
			fmt.Printf("   表达式: %s\n", expr)
			fmt.Printf("   返回: %s\n", result)
			return result, nil
		},
	)

	// 工具 4：天气查询
	weatherTool := utils.NewTool(
		&schema.ToolInfo{
			Name: "get_weather",
			Desc: "查询指定城市的实时天气",
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
			result := fmt.Sprintf("%s：晴天，温度22°C，湿度60%%，空气质量优", city)
			fmt.Printf("\n🔧 [工具] get_weather\n")
			fmt.Printf("   城市: %s\n", city)
			fmt.Printf("   返回: %s\n", result)
			return result, nil
		},
	)

	return []tool.BaseTool{timeTool, timeCalculator, calculator, weatherTool}
}
