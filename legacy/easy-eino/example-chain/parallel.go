package main

import (
	"context"
	"fmt"
	"log"
	"os"

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

	// 创建并行节点
	parallel := compose.NewParallel()

	// 任务1: 提取关键词
	parallel.AddLambda("keywords", compose.InvokableLambda(
		func(ctx context.Context, input map[string]any) (string, error) {
			fmt.Println("并行任务1: 提取关键词")

			template := prompt.FromMessages(
				schema.FString,
				schema.SystemMessage("请提取文本中的关键词，用逗号分隔。"),
				schema.UserMessage("{text}"),
			)

			messages, _ := template.Format(ctx, input)
			response, err := chatModel.Generate(ctx, messages)
			if err != nil {
				return "", err
			}
			return response.Content, nil
		},
	))

	// 任务2: 情感分析
	parallel.AddLambda("sentiment", compose.InvokableLambda(
		func(ctx context.Context, input map[string]any) (string, error) {
			fmt.Println("并行任务2: 情感分析")

			template := prompt.FromMessages(
				schema.FString,
				schema.SystemMessage("请分析文本的情感倾向（正面/负面/中性）。"),
				schema.UserMessage("{text}"),
			)

			messages, _ := template.Format(ctx, input)
			response, err := chatModel.Generate(ctx, messages)
			if err != nil {
				return "", err
			}
			return response.Content, nil
		},
	))

	// 任务3: 摘要生成
	parallel.AddLambda("summary", compose.InvokableLambda(
		func(ctx context.Context, input map[string]any) (string, error) {
			fmt.Println("并行任务3: 生成摘要")

			template := prompt.FromMessages(
				schema.FString,
				schema.SystemMessage("请用一句话总结文本内容。"),
				schema.UserMessage("{text}"),
			)

			messages, _ := template.Format(ctx, input)
			response, err := chatModel.Generate(ctx, messages)
			if err != nil {
				return "", err
			}
			return response.Content, nil
		},
	))

	// 创建主 Chain
	chain := compose.NewChain[string, map[string]any]()

	chain.
		// 准备输入
		AppendLambda(compose.InvokableLambda(func(ctx context.Context, text string) (map[string]any, error) {
			return map[string]any{"text": text}, nil
		})).
		// 并行执行三个任务
		AppendParallel(parallel).
		// 处理并行结果
		AppendLambda(compose.InvokableLambda(func(ctx context.Context, results map[string]any) (map[string]any, error) {
			fmt.Println("\\n=== 所有任务完成 ===")
			return results, nil
		}))

	runnable, err := chain.Compile(ctx)
	if err != nil {
		log.Fatalf("编译失败: %v", err)
	}

	text := "Eino 是一个优秀的 AI 开发框架，它让开发者能够轻松构建智能应用。框架设计精良，文档完善，社区活跃。"

	results, err := runnable.Invoke(ctx, text)
	if err != nil {
		log.Fatalf("执行失败: %v", err)
	}

	fmt.Printf("\\n关键词: %s\\n", results["keywords"])
	fmt.Printf("情感: %s\\n", results["sentiment"])
	fmt.Printf("摘要: %s\\n", results["summary"])
}
