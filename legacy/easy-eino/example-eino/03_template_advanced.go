/*
学习目标：
1. 掌握多消息模板（包含历史对话）
2. 学会构建可复用的模板库
3. 理解模板的最佳实践

实战场景：
构建一个代码审查助手，可以审查不同编程语言的代码

运行方式：
go run 03_template_advanced.go
*/

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

// 定义可复用的模板库
var (
	// 代码审查模板
	CodeReviewTemplate = prompt.FromMessages(
		schema.FString,
		schema.SystemMessage(`你是一个专业的{language}代码审查专家。
审查标准：
1. 代码规范性
2. 性能优化建议
3. 潜在的 bug
4. 最佳实践建议

请用简洁的语言给出审查意见。`),
		schema.UserMessage("请审查以下代码：\n\n```{language}\n{code}\n```"),
	)

	// 代码解释模板
	CodeExplainTemplate = prompt.FromMessages(
		schema.FString,
		schema.SystemMessage("你是一个{language}专家，擅长用{style}的方式解释代码。"),
		schema.UserMessage("请解释以下代码的作用：\n\n```{language}\n{code}\n```"),
	)

	// 技术对比模板
	TechCompareTemplate = prompt.FromMessages(
		schema.FString,
		schema.SystemMessage("你是一个技术专家，擅长对比不同技术的优劣。"),
		schema.UserMessage("请对比 {tech1} 和 {tech2} 在 {aspect} 方面的差异。"),
	)
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

	// 场景 1：审查 Go 代码
	fmt.Println("=== 场景 1：代码审查 ===\n")
	reviewGoCode(ctx, chatModel)

	// 场景 2：解释 Python 代码
	fmt.Println("\n=== 场景 2：代码解释 ===\n")
	explainPythonCode(ctx, chatModel)

	// 场景 3：技术对比
	fmt.Println("\n=== 场景 3：技术对比 ===\n")
	compareTech(ctx, chatModel)
}

func reviewGoCode(ctx context.Context, chatModel *deepseek.ChatModel) {
	variables := map[string]any{
		"language": "Go",
		"code": `func processUsers(users []User) {
    for i := 0; i < len(users); i++ {
        user := users[i]
        go func() {
            fmt.Println(user.Name)
        }()
    }
}`,
	}

	messages, err := CodeReviewTemplate.Format(ctx, variables)
	if err != nil {
		log.Fatalf("格式化失败: %v", err)
	}

	response, err := chatModel.Generate(ctx, messages)
	if err != nil {
		log.Fatalf("生成失败: %v", err)
	}

	fmt.Printf("审查意见:\n%s\n", response.Content)
}

func explainPythonCode(ctx context.Context, chatModel *deepseek.ChatModel) {
	variables := map[string]any{
		"language": "Python",
		"style":    "通俗易懂",
		"code": `@decorator
def hello(name):
    return f"Hello, {name}!"`,
	}

	messages, err := CodeExplainTemplate.Format(ctx, variables)
	if err != nil {
		log.Fatalf("格式化失败: %v", err)
	}

	response, err := chatModel.Generate(ctx, messages)
	if err != nil {
		log.Fatalf("生成失败: %v", err)
	}

	fmt.Printf("代码解释:\n%s\n", response.Content)
}

func compareTech(ctx context.Context, chatModel *deepseek.ChatModel) {
	variables := map[string]any{
		"tech1":  "Eino",
		"tech2":  "LangChain",
		"aspect": "类型安全性",
	}

	messages, err := TechCompareTemplate.Format(ctx, variables)
	if err != nil {
		log.Fatalf("格式化失败: %v", err)
	}

	response, err := chatModel.Generate(ctx, messages)
	if err != nil {
		log.Fatalf("生成失败: %v", err)
	}

	fmt.Printf("对比分析:\n%s\n", response.Content)
}
