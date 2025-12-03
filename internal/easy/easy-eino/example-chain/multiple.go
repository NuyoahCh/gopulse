package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/compose"
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

	// 构建处理链
	chain := compose.NewChain[string, string]()

	chain.
		// 步骤1: 数据清洗
		AppendLambda(compose.InvokableLambda(func(ctx context.Context, rawText string) (string, error) {
			fmt.Println("=== 步骤1: 数据清洗 ===")
			// 去除多余空格、换行等
			cleaned := strings.TrimSpace(rawText)
			cleaned = strings.ReplaceAll(cleaned, "\\n\\n", "\\n")
			fmt.Printf("清洗后: %s\\n\\n", cleaned)
			return cleaned, nil
		})).

		// 步骤2: 转换为 AI 分析输入
		AppendLambda(compose.InvokableLambda(func(ctx context.Context, text string) (map[string]any, error) {
			fmt.Println("=== 步骤2: 准备分析 ===")
			return map[string]any{
				"text": text,
			}, nil
		})).

		// 步骤3: AI 分析
		AppendGraph(func() *compose.Chain[map[string]any, *schema.Message] {
			analysisChain := compose.NewChain[map[string]any, *schema.Message]()

			template := prompt.FromMessages(
				schema.FString,
				schema.SystemMessage("你是一个文本分析专家。请分析以下文本的关键信息、主题和情感。"),
				schema.UserMessage("{text}"),
			)

			analysisChain.
				AppendChatTemplate(template).
				AppendChatModel(chatModel)

			return analysisChain
		}()).

		// 步骤4: 提取分析结果
		AppendLambda(compose.InvokableLambda(func(ctx context.Context, msg *schema.Message) (string, error) {
			fmt.Println("\\n=== 步骤3: 提取结果 ===")
			return msg.Content, nil
		}))

	runnable, err := chain.Compile(ctx)
	if err != nil {
		log.Fatalf("编译失败: %v", err)
	}

	// 测试输入
	rawInput := `
		Eino   是一个强大的  AI  开发框架。
		
		它提供了丰富的组件和灵活的编排能力。
		
		
		开发者可以快速构建   AI   应用。
	`

	result, err := runnable.Invoke(ctx, rawInput)
	if err != nil {
		log.Fatalf("执行失败: %v", err)
	}

	fmt.Printf("\\n=== 最终分析结果 ===\\n%s\\n", result)
}
