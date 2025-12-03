package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
	"io"
	"log"
	"os"
)

func main() {
	ctx := context.Background()

	// 创建模型和工具（省略，与上一个示例相同）
	chatModel, err := deepseek.NewChatModel(ctx, &deepseek.ChatModelConfig{
		APIKey:  os.Getenv("DEEPSEEK_API_KEY"),
		Model:   "deepseek-chat",
		BaseURL: "https://api.deepseek.com",
	})
	if err != nil {
		log.Fatalf("创建模型失败: %v", err)
	}

	searchTool := utils.NewTool(
		&schema.ToolInfo{
			Name: "search",
			Desc: "搜索信息",
			ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
				"query": {
					Type:     "string",
					Desc:     "搜索关键词",
					Required: true,
				},
			}),
		},
		func(ctx context.Context, params map[string]any) (string, error) {
			query := params["query"].(string)
			// 模拟搜索结果
			result := fmt.Sprintf("找到关于 '%s' 的信息: Go 是 Google 开发的编程语言，以并发和简洁著称。", query)
			fmt.Printf("\\n[工具执行] search(%s)\\n", query)
			return result, nil
		},
	)

	agent, err := react.NewAgent(ctx, &react.AgentConfig{
		ToolCallingModel: chatModel,
		ToolsConfig: compose.ToolsNodeConfig{
			Tools: []tool.BaseTool{searchTool},
		},
	})
	if err != nil {
		log.Fatalf("创建 Agent 失败: %v", err)
	}

	// 流式调用
	messages := []*schema.Message{
		schema.UserMessage("请告诉我 Go 语言的特点"),
	}

	fmt.Println("=== 用户: 请告诉我 Go 语言的特点 ===\\n")
	fmt.Print("Agent: ")

	stream, err := agent.Stream(ctx, messages)
	if err != nil {
		log.Fatalf("流式生成失败: %v", err)
	}
	defer stream.Close()

	// 逐块接收并打印
	for {
		chunk, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			log.Fatalf("接收失败: %v", err)
		}

		// 打印内容（打字机效果）
		fmt.Print(chunk.Content)
	}

	fmt.Println("\\n\\n=== 完成 ===")

}
