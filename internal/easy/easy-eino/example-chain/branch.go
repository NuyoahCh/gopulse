package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/cloudwego/eino/compose"
)

func main() {
	ctx := context.Background()

	// 定义分支条件
	branchCondition := func(ctx context.Context, input map[string]any) (string, error) {
		language := input["language"].(string)
		language = strings.ToLower(language)

		fmt.Printf("检测到语言: %s\\n", language)

		// 根据语言选择分支
		if language == "go" || language == "golang" {
			return "go_branch", nil
		} else if language == "python" {
			return "python_branch", nil
		}
		return "other_branch", nil
	}

	// Go 分支处理
	goBranch := compose.InvokableLambda(func(ctx context.Context, input map[string]any) (map[string]any, error) {
		fmt.Println("→ 执行 Go 分支")
		input["advice"] = "推荐使用 Eino 框架进行 AI 开发"
		input["features"] = []string{"高性能", "并发支持", "类型安全"}
		return input, nil
	})

	// Python 分支处理
	pythonBranch := compose.InvokableLambda(func(ctx context.Context, input map[string]any) (map[string]any, error) {
		fmt.Println("→ 执行 Python 分支")
		input["advice"] = "推荐使用 LangChain 框架进行 AI 开发"
		input["features"] = []string{"生态丰富", "易上手", "社区活跃"}
		return input, nil
	})

	// 其他语言分支
	otherBranch := compose.InvokableLambda(func(ctx context.Context, input map[string]any) (map[string]any, error) {
		fmt.Println("→ 执行其他语言分支")
		input["advice"] = "建议参考该语言的 AI 开发库"
		input["features"] = []string{"待探索"}
		return input, nil
	})

	// 创建 Chain
	chain := compose.NewChain[map[string]any, map[string]any]()

	chain.
		AppendLambda(compose.InvokableLambda(func(ctx context.Context, input map[string]any) (map[string]any, error) {
			fmt.Println("=== 开始处理 ===")
			return input, nil
		})).
		// 添加条件分支
		AppendBranch(compose.NewChainBranch(branchCondition).
			AddLambda("go_branch", goBranch).
			AddLambda("python_branch", pythonBranch).
			AddLambda("other_branch", otherBranch),
		).
		AppendLambda(compose.InvokableLambda(func(ctx context.Context, input map[string]any) (map[string]any, error) {
			fmt.Println("=== 处理完成 ===")
			return input, nil
		}))

	runnable, err := chain.Compile(ctx)
	if err != nil {
		log.Fatalf("编译失败: %v", err)
	}

	// 测试不同的输入
	testCases := []map[string]any{
		{"language": "Go", "task": "开发 AI 应用"},
		{"language": "Python", "task": "机器学习"},
		{"language": "Java", "task": "企业应用"},
	}

	for i, testCase := range testCases {
		fmt.Printf("\\n========== 测试 %d ==========\\n", i+1)
		result, err := runnable.Invoke(ctx, testCase)
		if err != nil {
			log.Printf("执行失败: %v", err)
			continue
		}

		fmt.Printf("建议: %s\\n", result["advice"])
		fmt.Printf("特性: %v\\n", result["features"])
	}
}
