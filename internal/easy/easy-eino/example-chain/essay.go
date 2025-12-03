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

type ArticleRequest struct {
	Topic    string
	Keywords []string
	Length   int // 目标字数
}

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

	// 构建文章生成流水线
	chain := compose.NewChain[ArticleRequest, string]()

	chain.
		// 步骤1: 生成大纲
		AppendLambda(compose.InvokableLambda(func(ctx context.Context, req ArticleRequest) (string, error) {
			fmt.Println("=== 步骤1: 生成文章大纲 ===")

			template := prompt.FromMessages(
				schema.FString,
				schema.SystemMessage("你是一个专业的内容策划师。请根据主题和关键词生成文章大纲。"),
				schema.UserMessage("主题: {topic}\\n关键词: {keywords}\\n\\n请生成一个3-5点的文章大纲。"),
			)

			messages, _ := template.Format(ctx, map[string]any{
				"topic":    req.Topic,
				"keywords": fmt.Sprintf("%v", req.Keywords),
			})

			response, err := chatModel.Generate(ctx, messages)
			if err != nil {
				return "", err
			}

			fmt.Printf("大纲:\\n%s\\n\\n", response.Content)
			return response.Content, nil
		})).

		// 步骤2: 扩写内容
		AppendLambda(compose.InvokableLambda(func(ctx context.Context, outline string) (string, error) {
			fmt.Println("=== 步骤2: 扩写文章内容 ===")

			template := prompt.FromMessages(
				schema.FString,
				schema.SystemMessage("你是一个专业的内容作者。请根据大纲扩写成完整文章，要求内容充实、逻辑清晰。"),
				schema.UserMessage("大纲:\\n{outline}\\n\\n请扩写成一篇800字左右的文章。"),
			)

			messages, _ := template.Format(ctx, map[string]any{
				"outline": outline,
			})

			response, err := chatModel.Generate(ctx, messages)
			if err != nil {
				return "", err
			}

			fmt.Printf("初稿完成，字数: %d\\n\\n", len(response.Content))
			return response.Content, nil
		})).

		// 步骤3: 润色优化
		AppendLambda(compose.InvokableLambda(func(ctx context.Context, draft string) (string, error) {
			fmt.Println("=== 步骤3: 润色优化 ===")

			template := prompt.FromMessages(
				schema.FString,
				schema.SystemMessage("你是一个专业的编辑。请优化文章的语言表达，使其更加流畅、生动。"),
				schema.UserMessage("文章:\\n{draft}\\n\\n请进行润色优化。"),
			)

			messages, _ := template.Format(ctx, map[string]any{
				"draft": draft,
			})

			response, err := chatModel.Generate(ctx, messages)
			if err != nil {
				return "", err
			}

			fmt.Println("润色完成\\n")
			return response.Content, nil
		})).

		// 步骤4: 格式化输出
		AppendLambda(compose.InvokableLambda(func(ctx context.Context, article string) (string, error) {
			fmt.Println("=== 步骤4: 格式化输出 ===")

			// 添加 Markdown 格式
			formatted := fmt.Sprintf("# 生成的文章\\n\\n%s\\n\\n---\\n*由 Eino AI 助手生成*", article)
			return formatted, nil
		}))

	runnable, err := chain.Compile(ctx)
	if err != nil {
		log.Fatalf("编译失败: %v", err)
	}

	// 生成文章
	request := ArticleRequest{
		Topic:    "AI 在软件开发中的应用",
		Keywords: []string{"AI", "编程", "效率", "未来"},
		Length:   800,
	}

	result, err := runnable.Invoke(ctx, request)
	if err != nil {
		log.Fatalf("执行失败: %v", err)
	}

	//fmt.Println("\\n" + "="*50)
	fmt.Println(result)
	//fmt.Println("=" * 50)

	// 保存到文件
	os.WriteFile("generated_article.md", []byte(result), 0644)
	fmt.Println("\\n文章已保存到 generated_article.md")
}
