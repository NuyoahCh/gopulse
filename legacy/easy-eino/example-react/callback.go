package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
)

// 自定义 Callback
type MyCallback struct {
	callbacks.HandlerBuilder
}

func (cb *MyCallback) OnStart(ctx context.Context, info *callbacks.RunInfo, input callbacks.CallbackInput) context.Context {
	fmt.Printf("\\n[Callback] 开始执行: %s\\n", info.Name)
	inputJSON, _ := json.MarshalIndent(input, "", "  ")
	fmt.Printf("输入: %s\\n", string(inputJSON))
	return ctx
}

func (cb *MyCallback) OnEnd(ctx context.Context, info *callbacks.RunInfo, output callbacks.CallbackOutput) context.Context {
	fmt.Printf("\\n[Callback] 完成执行: %s\\n", info.Name)
	return ctx
}

func (cb *MyCallback) OnEndWithStreamOutput(ctx context.Context, info *callbacks.RunInfo, output *schema.StreamReader[callbacks.CallbackOutput]) context.Context {
	fmt.Printf("\\n[Callback] 完成流式执行: %s\\n", info.Name)
	return ctx
}

func (cb *MyCallback) OnError(ctx context.Context, info *callbacks.RunInfo, err error) context.Context {
	fmt.Printf("\\n[Callback] 执行出错: %s, 错误: %v\\n", info.Name, err)
	return ctx
}

func (cb *MyCallback) OnStartWithStreamInput(ctx context.Context, info *callbacks.RunInfo, input *schema.StreamReader[callbacks.CallbackInput]) context.Context {
	fmt.Printf("\\n[Callback] 开始流式执行: %s\\n", info.Name)
	return ctx
}

func main() {
	ctx := context.Background()

	// 创建模型
	chatModel, err := deepseek.NewChatModel(ctx, &deepseek.ChatModelConfig{
		APIKey:  os.Getenv("DEEPSEEK_API_KEY"),
		Model:   "deepseek-chat",
		BaseURL: "https://api.deepseek.com",
	})
	if err != nil {
		log.Fatalf("创建模型失败: %v", err)
	}

	testTool := utils.NewTool(
		&schema.ToolInfo{
			Name:        "test_tool",
			Desc:        "测试工具",
			ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{}),
		},
		func(ctx context.Context, params map[string]any) (string, error) {
			return "工具执行成功", nil
		},
	)

	ragent, _ := react.NewAgent(ctx, &react.AgentConfig{
		ToolCallingModel: chatModel,
		ToolsConfig: compose.ToolsNodeConfig{
			Tools: []tool.BaseTool{testTool},
		},
	})

	// 使用 Callback
	response, _ := ragent.Generate(ctx,
		[]*schema.Message{
			schema.UserMessage("测试一下"),
		},
		agent.WithComposeOptions(compose.WithCallbacks(&MyCallback{})),
	)

	fmt.Printf("\\n最终回答: %s\\n", response.Content)
}
