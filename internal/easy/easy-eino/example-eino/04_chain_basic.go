/*
学习目标：
1. 理解 Chain 的本质：数据流管道
2. 掌握 Chain 的基本用法：Template → Model
3. 理解编译和执行的两阶段模式

核心概念：
- Chain：将多个组件串联成一个数据处理管道
- Compile：编译阶段，进行类型检查和优化
- Invoke：执行阶段，传入数据并获得结果

数据流：
Input (map[string]any)
  → ChatTemplate (格式化)
  → []*schema.Message
  → ChatModel (生成)
  → *schema.Message (Output)

运行方式：
go run 04_chain_basic.go
*/

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

	// 步骤 1：创建 ChatTemplate（第一个节点）
	chatTemplate := prompt.FromMessages(
		schema.FString,
		schema.SystemMessage("你是一个{role}"),
		schema.UserMessage("{question}"),
	)

	// 步骤 2：创建 ChatModel（第二个节点）
	chatModel, err := deepseek.NewChatModel(ctx, &deepseek.ChatModelConfig{
		APIKey:  os.Getenv("DEEPSEEK_API_KEY"),
		Model:   "deepseek-chat",
		BaseURL: "https://api.deepseek.com",
	})
	if err != nil {
		log.Fatalf("创建模型失败: %v", err)
	}

	// 步骤 3：创建 Chain 并串联组件
	// NewChain[InputType, OutputType]() 定义了整个 Chain 的输入输出类型
	chain := compose.NewChain[map[string]any, *schema.Message]()

	chain.
		AppendChatTemplate(chatTemplate). // 节点1：格式化模板
		AppendChatModel(chatModel)        // 节点2：调用模型

	fmt.Println("=== Chain 结构 ===")
	fmt.Println("Input: map[string]any")
	fmt.Println("  ↓")
	fmt.Println("ChatTemplate (格式化)")
	fmt.Println("  ↓")
	fmt.Println("[]*schema.Message")
	fmt.Println("  ↓")
	fmt.Println("ChatModel (生成)")
	fmt.Println("  ↓")
	fmt.Println("Output: *schema.Message")
	fmt.Println()

	// 步骤 4：编译 Chain
	// Compile 会进行类型检查，确保节点之间的类型匹配
	runnable, err := chain.Compile(ctx)
	if err != nil {
		log.Fatalf("编译失败: %v", err)
	}
	fmt.Println("✅ Chain 编译成功")

	// 步骤 5：执行 Chain
	input := map[string]any{
		"role":     "Go 语言专家",
		"question": "请用一句话解释 Eino 的 Chain 是什么？",
	}

	fmt.Printf("\n输入: %v\n\n", input)

	output, err := runnable.Invoke(ctx, input)
	if err != nil {
		log.Fatalf("执行失败: %v", err)
	}

	fmt.Printf("输出: %s\n", output.Content)

	// 思考题：
	// 1. 如果 ChatTemplate 的输出类型和 ChatModel 的输入类型不匹配会怎样？
	// 2. 为什么需要 Compile 这一步？能否直接执行？
	// 3. Chain 和直接调用组件有什么区别？
}
