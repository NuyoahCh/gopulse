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

	// 创建一个简单的 Chain：输入字符串 → 转大写 → 添加前缀 → 输出
	chain := compose.NewChain[string, string]()

	chain.
		// Lambda 1: 转大写
		AppendLambda(compose.InvokableLambda(func(ctx context.Context, input string) (string, error) {
			fmt.Printf("步骤1: 输入 = %s\\n", input)
			result := strings.ToUpper(input)
			fmt.Printf("步骤1: 输出 = %s\\n", result)
			return result, nil
		})).
		// Lambda 2: 添加前缀
		AppendLambda(compose.InvokableLambda(func(ctx context.Context, input string) (string, error) {
			fmt.Printf("步骤2: 输入 = %s\\n", input)
			result := "处理结果: " + input
			fmt.Printf("步骤2: 输出 = %s\\n", result)
			return result, nil
		}))

	runnable, err := chain.Compile(ctx)
	if err != nil {
		log.Fatalf("编译失败: %v", err)
	}

	output, err := runnable.Invoke(ctx, "hello eino")
	if err != nil {
		log.Fatalf("执行失败: %v", err)
	}

	fmt.Printf("\\n最终输出: %s\\n", output)
}
