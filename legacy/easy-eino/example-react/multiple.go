package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
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

	// 创建一些实用工具
	tools := []tool.BaseTool{
		// 这里可以添加各种工具
	}

	agent, err := react.NewAgent(ctx, &react.AgentConfig{
		ToolCallingModel: chatModel,
		ToolsConfig: compose.ToolsNodeConfig{
			Tools: tools,
		},
	})
	if err != nil {
		log.Fatalf("创建 Agent 失败: %v", err)
	}

	// 对话历史
	messages := []*schema.Message{
		schema.SystemMessage("你是一个友好的 AI 助手，可以帮助用户解决问题。"),
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("=== 多轮对话 Agent（输入 'exit' 退出）===\\n")

	for {
		fmt.Print("你: ")
		if !scanner.Scan() {
			break
		}

		userInput := strings.TrimSpace(scanner.Text())
		if userInput == "exit" {
			fmt.Println("再见！")
			break
		}

		if userInput == "" {
			continue
		}

		// 添加用户消息
		messages = append(messages, schema.UserMessage(userInput))

		// Agent 生成回复
		response, err := agent.Generate(ctx, messages)
		if err != nil {
			fmt.Printf("生成失败: %v\\n", err)
			continue
		}

		// 添加 AI 回复到历史
		messages = append(messages, response)

		fmt.Printf("\\nAgent: %s\\n\\n", response.Content)
	}
}
