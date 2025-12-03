/*
学习目标：
1. 理解 AI 如何调用工具（Function Calling）
2. 掌握 ToolsNode 的使用
3. 学会构建工具调用链路

核心概念：
- Function Calling：AI 自动选择和调用工具
- ToolsNode：工具执行节点
- 工具调用流程：问题 → AI 分析 → 选择工具 → 执行 → 返回结果

运行方式：
go run 06_tool_with_model.go
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

	// 1. 创建支持 Function Calling 的 ChatModel
	chatModel, err := deepseek.NewChatModel(ctx, &deepseek.ChatModelConfig{
		APIKey:  os.Getenv("DEEPSEEK_API_KEY"),
		Model:   "deepseek-chat",
		BaseURL: "https://api.deepseek.com",
	})
	if err != nil {
		log.Fatalf("创建模型失败: %v", err)
	}

	// 2. 创建工具集
	tools := createTools()

	// 3. 使用 ReAct Agent（这是 Eino 推荐的工具调用方式）
	// ReAct Agent 会自动处理：Model → 工具选择 → 工具执行 → Model 处理结果
	agent, err := react.NewAgent(ctx, &react.AgentConfig{
		ToolCallingModel: chatModel,
		ToolsConfig: compose.ToolsNodeConfig{
			Tools: tools,
		},
		MaxStep: 3, // 最多推理 3 步
	})
	if err != nil {
		log.Fatalf("创建 Agent 失败: %v", err)
	}

	// 4. 测试不同的问题
	testQuestions := []string{
		"现在几点了？",
		"帮我计算 123 + 456",
		"北京今天天气怎么样？",
		"你好，介绍一下自己", // 不需要工具
	}

	for i, question := range testQuestions {
		fmt.Printf("\n========== 测试 %d ==========\n", i+1)
		fmt.Printf("❓ 问题: %s\n\n", question)

		messages := []*schema.Message{
			schema.UserMessage(question),
		}

		response, err := agent.Generate(ctx, messages)
		if err != nil {
			log.Printf("执行失败: %v", err)
			continue
		}

		fmt.Printf("🤖 回答: %s\n", response.Content)
	}

	// 思考题：
	// 1. AI 如何知道该调用哪个工具？
	// 2. 如果没有合适的工具，AI 会怎么做？
	// 3. 工具调用失败了怎么办？
}

// 创建工具集
func createTools() []tool.BaseTool {
	// 工具 1：获取当前时间
	// 注意：对于无参数的工具，ParamsOneOf 应该设置为 nil 或者使用空对象 schema.NewParamsOneOfByParams(nil)
	timeTool := utils.NewTool(
		&schema.ToolInfo{
			Name:        "get_current_time",
			Desc:        "获取当前的日期和时间",
			ParamsOneOf: nil, // 无参数工具
		},
		func(ctx context.Context, params map[string]any) (string, error) {
			now := time.Now().Format("2006-01-02 15:04:05")
			fmt.Printf("  🔧 [工具执行] get_current_time → %s\n", now)
			return now, nil
		},
	)

	// 工具 2：计算器
	calculator := utils.NewTool(
		&schema.ToolInfo{
			Name: "calculator",
			Desc: "执行数学计算",
			ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
				"expression": {
					Type:     "string",
					Desc:     "数学表达式，例如: 123 + 456",
					Required: true,
				},
			}),
		},
		func(ctx context.Context, params map[string]any) (string, error) {
			expr := params["expression"].(string)
			// 简化：实际应该用表达式解析库
			result := "579" // 假设结果
			fmt.Printf("  🔧 [工具执行] calculator(%s) → %s\n", expr, result)
			return result, nil
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
			result := fmt.Sprintf("%s：晴天，22°C，湿度60%%", city)
			fmt.Printf("  🔧 [工具执行] get_weather(%s) → %s\n", city, result)
			return result, nil
		},
	)

	return []tool.BaseTool{timeTool, calculator, weatherTool}
}
