/*
学习目标：
1. 掌握并行执行的使用
2. 理解并发处理的优势
3. 学会合并并行结果

核心概念：
- Parallel：并行执行节点
- 并发执行：多个任务同时执行
- 结果合并：将多个并行结果合并成一个输出

使用场景：
- 多个独立任务同时执行（提高性能）
- 从多个数据源获取数据
- 并行调用多个 AI 模型

运行方式：
go run 04_chain_parallel.go
*/

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cloudwego/eino/compose"
)

func main() {
	ctx := context.Background()

	// 创建并行节点
	parallel := compose.NewParallel()

	// 任务 1：计算平方
	parallel.AddLambda("square", compose.InvokableLambda(
		func(ctx context.Context, num int) (int, error) {
			fmt.Println("📊 任务1: 计算平方")
			time.Sleep(1 * time.Second) // 模拟耗时操作
			result := num * num
			fmt.Printf("📊 任务1 完成: %d² = %d\n", num, result)
			return result, nil
		},
	))

	// 任务 2：计算立方
	parallel.AddLambda("cube", compose.InvokableLambda(
		func(ctx context.Context, num int) (int, error) {
			fmt.Println("📈 任务2: 计算立方")
			time.Sleep(1 * time.Second) // 模拟耗时操作
			result := num * num * num
			fmt.Printf("📈 任务2 完成: %d³ = %d\n", num, result)
			return result, nil
		},
	))

	// 任务 3：计算阶乘
	parallel.AddLambda("factorial", compose.InvokableLambda(
		func(ctx context.Context, num int) (int, error) {
			fmt.Println("🔢 任务3: 计算阶乘")
			time.Sleep(1 * time.Second) // 模拟耗时操作
			result := 1
			for i := 2; i <= num; i++ {
				result *= i
			}
			fmt.Printf("🔢 任务3 完成: %d! = %d\n", num, result)
			return result, nil
		},
	))

	// 创建主 Chain
	chain := compose.NewChain[int, map[string]any]()

	chain.
		// 前置处理
		AppendLambda(compose.InvokableLambda(func(ctx context.Context, num int) (int, error) {
			fmt.Printf("\n=== 开始并行计算: 输入 = %d ===\n\n", num)
			return num, nil
		})).

		// 并行执行三个任务
		AppendParallel(parallel).

		// 后置处理：合并结果
		AppendLambda(compose.InvokableLambda(func(ctx context.Context, results map[string]any) (map[string]any, error) {
			fmt.Println("\n=== 所有任务完成 ===")

			// 计算总和
			total := results["square"].(int) + results["cube"].(int) + results["factorial"].(int)
			results["total"] = total

			return results, nil
		}))

	// 编译
	runnable, err := chain.Compile(ctx)
	if err != nil {
		log.Fatalf("编译失败: %v", err)
	}

	// 执行并计时
	startTime := time.Now()

	results, err := runnable.Invoke(ctx, 5)
	if err != nil {
		log.Fatalf("执行失败: %v", err)
	}

	elapsed := time.Since(startTime)

	// 输出结果
	fmt.Println("\n=== 计算结果 ===")
	fmt.Printf("平方: %d\n", results["square"])
	fmt.Printf("立方: %d\n", results["cube"])
	fmt.Printf("阶乘: %d\n", results["factorial"])
	fmt.Printf("总和: %d\n", results["total"])
	fmt.Printf("\n⏱️  总耗时: %v\n", elapsed)
	fmt.Println("💡 注意：三个任务并行执行，总耗时约为单个任务的时间（~1秒）而非三倍（~3秒）")

	// 思考题：
	// 1. 如果某个并行任务失败了，整个 Chain 会怎样？
	// 2. 并行任务之间能否共享数据？
	// 3. 如何控制并行任务的超时？
}
