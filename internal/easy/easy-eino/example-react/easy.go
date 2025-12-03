package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino/components/tool"
	"log"
	"os"
	"time"

	"github.com/cloudwego/eino-ext/components/model/deepseek"
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

	// 2. 创建工具
	// 获取当前时间的工具
	timeTool := utils.NewTool(
		&schema.ToolInfo{
			Name:        "get_current_time",
			Desc:        "获取当前时间",
			ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{}),
		},
		func(ctx context.Context, params map[string]any) (string, error) {
			now := time.Now().Format("2006-01-02 15:04:05")
			fmt.Printf("[工具执行] get_current_time -> %s\\n", now)
			return now, nil
		},
	)

	// 简单计算器工具
	calculator := utils.NewTool(
		&schema.ToolInfo{
			Name: "calculator",
			Desc: "执行简单的数学计算（加减乘除）",
			ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
				"expression": {
					Type:     "string",
					Desc:     "数学表达式，例如: 10 + 5",
					Required: true,
				},
			}),
		},
		func(ctx context.Context, params map[string]any) (string, error) {
			expr := params["expression"].(string)
			// 简化：这里只演示，实际应该用真正的表达式解析
			result := "15" // 假设结果
			fmt.Printf("[工具执行] calculator(%s) -> %s\\n", expr, result)
			return result, nil
		},
	)

	// 3. 创建 React Agent
	agent, err := react.NewAgent(ctx, &react.AgentConfig{
		ToolCallingModel: chatModel,
		ToolsConfig: compose.ToolsNodeConfig{
			Tools: []tool.BaseTool{timeTool, calculator},
		},
	})
	if err != nil {
		log.Fatalf("创建 Agent 失败: %v", err)
	}

	// 4. 使用 Agent
	messages := []*schema.Message{
		schema.UserMessage("现在几点了？"),
	}

	fmt.Println("=== 用户: 现在几点了？ ===\\n")

	response, err := agent.Generate(ctx, messages)
	if err != nil {
		log.Fatalf("生成失败: %v", err)
	}

	fmt.Printf("\\n=== Agent 回答 ===\\n%s\\n", response.Content)
}
