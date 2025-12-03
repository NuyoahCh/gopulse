/*
学习目标：
1. 掌握 Lambda 节点的使用
2. 理解如何在 Chain 中插入自定义逻辑
3. 学会数据转换和处理

核心概念：
- Lambda：自定义逻辑节点，可以执行任意 Go 代码
- InvokableLambda：将函数包装成可调用的节点
- 数据转换：在节点之间进行数据格式转换

使用场景：
- 数据预处理（清洗、验证）
- 数据后处理（格式化、过滤）
- 业务逻辑（计算、判断）

运行方式：
go run 04_chain_lambda.go
*/

package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/cloudwego/eino/compose"
)

func main() {
	ctx := context.Background()

	// 场景：文本处理管道
	// 输入字符串 → 清洗 → 转大写 → 添加时间戳 → 输出

	chain := compose.NewChain[string, string]()

	chain.
		// Lambda 1: 数据清洗（去除空格和特殊字符）
		AppendLambda(compose.InvokableLambda(func(ctx context.Context, input string) (string, error) {
			fmt.Printf("步骤1 [清洗] 输入: %q\n", input)

			// 去除首尾空格
			cleaned := strings.TrimSpace(input)
			// 去除多余空格
			cleaned = strings.Join(strings.Fields(cleaned), " ")

			fmt.Printf("步骤1 [清洗] 输出: %q\n\n", cleaned)
			return cleaned, nil
		})).

		// Lambda 2: 转大写
		AppendLambda(compose.InvokableLambda(func(ctx context.Context, input string) (string, error) {
			fmt.Printf("步骤2 [转大写] 输入: %q\n", input)

			result := strings.ToUpper(input)

			fmt.Printf("步骤2 [转大写] 输出: %q\n\n", result)
			return result, nil
		})).

		// Lambda 3: 添加时间戳
		AppendLambda(compose.InvokableLambda(func(ctx context.Context, input string) (string, error) {
			fmt.Printf("步骤3 [添加时间戳] 输入: %q\n", input)

			timestamp := time.Now().Format("2006-01-02 15:04:05")
			result := fmt.Sprintf("[%s] %s", timestamp, input)

			fmt.Printf("步骤3 [添加时间戳] 输出: %q\n\n", result)
			return result, nil
		}))

	// 编译 Chain
	runnable, err := chain.Compile(ctx)
	if err != nil {
		log.Fatalf("编译失败: %v", err)
	}

	// 执行 Chain
	input := "  hello   eino   framework  "
	fmt.Printf("原始输入: %q\n\n", input)
	fmt.Println("=== 开始处理 ===\n")

	output, err := runnable.Invoke(ctx, input)
	if err != nil {
		log.Fatalf("执行失败: %v", err)
	}

	fmt.Println("=== 处理完成 ===")
	fmt.Printf("\n最终输出: %q\n", output)

	// 示例 2：类型转换 Chain
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("示例 2：类型转换 Chain")
	fmt.Println(strings.Repeat("=", 60) + "\n")

	typeConversionDemo()
}

// 示例 2：演示不同类型之间的转换
func typeConversionDemo() {
	ctx := context.Background()

	// Chain: int → string → map[string]any
	chain := compose.NewChain[int, map[string]any]()

	chain.
		// Lambda 1: int → string
		AppendLambda(compose.InvokableLambda(func(ctx context.Context, num int) (string, error) {
			result := fmt.Sprintf("数字 %d", num)
			fmt.Printf("转换1: %d → %q\n", num, result)
			return result, nil
		})).

		// Lambda 2: string → map[string]any
		AppendLambda(compose.InvokableLambda(func(ctx context.Context, str string) (map[string]any, error) {
			result := map[string]any{
				"text":      str,
				"length":    len(str),
				"timestamp": time.Now().Unix(),
			}
			fmt.Printf("转换2: %q → %v\n", str, result)
			return result, nil
		}))

	runnable, _ := chain.Compile(ctx)
	output, _ := runnable.Invoke(ctx, 42)

	fmt.Printf("\n最终结果: %v\n", output)

	// 思考题：
	// 1. Lambda 中可以访问外部变量吗？
	// 2. 如果 Lambda 返回 error，Chain 会如何处理？
	// 3. Lambda 可以是异步的吗？
}
