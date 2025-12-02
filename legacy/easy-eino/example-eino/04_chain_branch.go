/*
学习目标：
1. 掌握条件分支的使用
2. 理解动态路由的概念
3. 学会构建复杂的业务逻辑

核心概念：
- Branch：条件分支节点
- 路由函数：根据输入决定执行哪个分支
- 分支处理：不同分支执行不同的逻辑

使用场景：
- 根据用户类型提供不同服务
- 根据数据类型选择不同处理方式
- 实现复杂的业务规则

运行方式：
go run 04_chain_branch.go
*/

package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/cloudwego/eino/compose"
)

// 定义用户请求结构
type UserRequest struct {
	UserType string // "vip", "normal", "guest"
	Action   string
	Data     string
}

func main() {
	ctx := context.Background()

	// 定义路由函数：根据用户类型选择分支
	routeFunc := func(ctx context.Context, req UserRequest) (string, error) {
		userType := strings.ToLower(req.UserType)
		fmt.Printf("🔀 路由判断: 用户类型 = %s\n", userType)

		switch userType {
		case "vip":
			return "vip_branch", nil
		case "normal":
			return "normal_branch", nil
		default:
			return "guest_branch", nil
		}
	}

	// VIP 用户分支
	vipBranch := compose.InvokableLambda(func(ctx context.Context, req UserRequest) (string, error) {
		fmt.Println("✨ 执行 VIP 分支")
		result := fmt.Sprintf("VIP 用户 %s，享受优先服务！处理结果：%s（已加速）",
			req.Action, strings.ToUpper(req.Data))
		return result, nil
	})

	// 普通用户分支
	normalBranch := compose.InvokableLambda(func(ctx context.Context, req UserRequest) (string, error) {
		fmt.Println("👤 执行普通用户分支")
		result := fmt.Sprintf("普通用户 %s，标准服务。处理结果：%s",
			req.Action, req.Data)
		return result, nil
	})

	// 访客分支
	guestBranch := compose.InvokableLambda(func(ctx context.Context, req UserRequest) (string, error) {
		fmt.Println("🚶 执行访客分支")
		result := fmt.Sprintf("访客 %s，功能受限。请注册以获得更多服务。", req.Action)
		return result, nil
	})

	// 创建 Chain
	chain := compose.NewChain[UserRequest, string]()

	chain.
		// 前置处理
		AppendLambda(compose.InvokableLambda(func(ctx context.Context, req UserRequest) (UserRequest, error) {
			fmt.Println("=== 开始处理请求 ===")
			fmt.Printf("用户类型: %s, 操作: %s\n", req.UserType, req.Action)
			return req, nil
		})).

		// 条件分支
		AppendBranch(compose.NewChainBranch(routeFunc).
			AddLambda("vip_branch", vipBranch).
			AddLambda("normal_branch", normalBranch).
			AddLambda("guest_branch", guestBranch),
		).

		// 后置处理
		AppendLambda(compose.InvokableLambda(func(ctx context.Context, result string) (string, error) {
			fmt.Println("=== 请求处理完成 ===")
			return result, nil
		}))

	// 编译
	runnable, err := chain.Compile(ctx)
	if err != nil {
		log.Fatalf("编译失败: %v", err)
	}

	// 测试不同类型的用户
	testCases := []UserRequest{
		{UserType: "VIP", Action: "查询数据", Data: "eino framework"},
		{UserType: "normal", Action: "查询数据", Data: "hello world"},
		{UserType: "guest", Action: "查询数据", Data: "test"},
	}

	// 测试
	for i, req := range testCases {
		fmt.Printf("\n========== 测试 %d ==========\n", i+1)
		result, err := runnable.Invoke(ctx, req)
		if err != nil {
			log.Printf("执行失败: %v", err)
			continue
		}
		fmt.Printf("\n📋 结果: %s\n", result)
	}

	// 思考题：
	// 1. 如果路由函数返回一个不存在的分支名会怎样？
	// 2. 能否在分支内部再嵌套分支？
	// 3. 如何实现默认分支（类似 switch 的 default）？
}