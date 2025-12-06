package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"einoflow/internal/agent"
	"einoflow/internal/config"
	"einoflow/internal/graph"
	"einoflow/internal/tools"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"github.com/joho/godotenv"
)

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// 创建豆包聊天模型
	chatModel, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey:  cfg.ArkAPIKey,
		BaseURL: cfg.ArkBaseURL,
		Model:   "doubao-seed-1-6-lite-251015", // 豆包 Lite 模型
	})
	if err != nil {
		log.Fatal(err)
	}

	// 创建输入扫描器
	scanner := bufio.NewScanner(os.Stdin)

	// 演示菜单
	for {
		fmt.Println("\n=== EinoFlow 功能演示 ===")
		fmt.Println("1. 基础对话")
		fmt.Println("2. 流式对话")
		fmt.Println("3. Agent 工具调用")
		fmt.Println("4. Graph 多步骤处理")
		fmt.Println("5. 退出")
		fmt.Print("\n请选择功能 (1-5): ")

		if !scanner.Scan() {
			break
		}
		choiceStr := strings.TrimSpace(scanner.Text())
		choice, err := strconv.Atoi(choiceStr)
		if err != nil {
			fmt.Println("无效选择，请输入数字 1-5")
			continue
		}

		switch choice {
		case 1:
			demoBasicChat(ctx, chatModel)
		case 2:
			demoStreamChat(ctx, chatModel)
		case 3:
			demoAgent(ctx, chatModel, cfg)
		case 4:
			demoGraph(ctx, chatModel)
		case 5:
			fmt.Println("再见！")
			os.Exit(0)
		default:
			fmt.Println("无效选择，请重试")
		}
	}
}

func demoBasicChat(ctx context.Context, chatModel model.ChatModel) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("\n请输入问题: ")
	if !scanner.Scan() {
		return
	}
	question := strings.TrimSpace(scanner.Text())

	messages := []*schema.Message{
		schema.UserMessage(question),
	}

	fmt.Println("\n思考中...")
	resp, err := chatModel.Generate(ctx, messages)
	if err != nil {
		log.Printf("错误: %v\n", err)
		return
	}

	fmt.Printf("\n回答: %s\n", resp.Content)
}

func demoStreamChat(ctx context.Context, chatModel model.ChatModel) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("\n请输入问题: ")
	if !scanner.Scan() {
		return
	}
	question := strings.TrimSpace(scanner.Text())

	messages := []*schema.Message{
		schema.UserMessage(question),
	}

	fmt.Println("\n回答: ")
	stream, err := chatModel.Stream(ctx, messages)
	if err != nil {
		log.Printf("错误: %v\n", err)
		return
	}
	defer stream.Close()

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("错误: %v\n", err)
			break
		}
		fmt.Print(msg.Content)
	}
	fmt.Println()
}

func demoAgent(ctx context.Context, chatModel model.ChatModel, cfg *config.Config) {
	// 创建工具注册表
	registry := tools.NewRegistry(cfg.DBPath, "./data/documents")
	_ = registry.GetAll() // 工具列表（当前简化版本不使用）

	// 创建 Agent
	reactAgent := agent.NewReActAgent(chatModel)

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("\n请输入任务: ")
	if !scanner.Scan() {
		return
	}
	task := strings.TrimSpace(scanner.Text())

	fmt.Println("\nAgent 执行中...")
	fmt.Println("⏳ 正在调用模型生成回答，请稍候...")

	result, err := reactAgent.Run(ctx, task)
	if err != nil {
		log.Printf("错误: %v\n", err)
		return
	}

	fmt.Printf("\nAgent 响应:\n%s\n", result)
}

func demoGraph(ctx context.Context, chatModel model.ChatModel) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("\n请输入复杂问题: ")
	if !scanner.Scan() {
		return
	}
	question := strings.TrimSpace(scanner.Text())

	// 创建多步骤处理图
	g := graph.CreateMultiStepGraph(chatModel)

	// 设置进度回调
	g.SetProgressCallback(func(nodeID string, input, output interface{}) {
		switch nodeID {
		case "analyze":
			fmt.Println("✅ 步骤 1/3: 问题分析完成")
			fmt.Println("步骤 2/3: 制定计划...")
		case "plan":
			fmt.Println("✅ 步骤 2/3: 计划制定完成")
			fmt.Println("步骤 3/3: 执行并总结...")
		case "execute":
			fmt.Println("✅ 步骤 3/3: 执行完成")
		}
	})

	fmt.Println("\n执行多步骤分析...")
	fmt.Println("步骤 1/3: 分析问题...")

	// 执行图
	result, err := g.Run(ctx, question)
	if err != nil {
		log.Printf("错误: %v\n", err)
		return
	}

	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Printf("\n最终答案:\n%s\n", result)
	fmt.Println("\n" + strings.Repeat("=", 50))
}
