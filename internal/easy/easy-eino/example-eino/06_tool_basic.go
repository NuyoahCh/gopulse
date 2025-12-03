/*
学习目标：
1. 理解什么是工具（Tool）
2. 掌握工具的三要素：Info、Params、Run
3. 学会定义和使用简单工具

核心概念：
- Tool：扩展 AI 能力的外部函数
- ToolInfo：工具的元信息（名称、描述、参数）
- InvokableRun：工具的执行逻辑

为什么需要工具？
- LLM 的局限：无法访问实时数据、执行计算、调用 API
- 工具的价值：让 AI 能够"行动"，而不仅仅是"对话"

运行方式：
go run 06_tool_basic.go
*/

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

// 示例 1：获取当前时间的工具
type CurrentTimeTool struct{}

// Info 返回工具的元信息
// 这是 AI 用来理解工具功能的关键信息
func (t *CurrentTimeTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "get_current_time",
		Desc: "获取当前的日期和时间",
		// 这个工具不需要参数
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{}),
	}, nil
}

// InvokableRun 执行工具逻辑
// argumentsInJSON：AI 传入的参数（JSON 格式）
func (t *CurrentTimeTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	// 获取当前时间
	now := time.Now()

	// 格式化结果
	result := map[string]string{
		"datetime": now.Format("2006-01-02 15:04:05"),
		"date":     now.Format("2006-01-02"),
		"time":     now.Format("15:04:05"),
		"weekday":  now.Weekday().String(),
	}

	// 返回 JSON 格式的结果
	resultJSON, err := json.Marshal(result)
	if err != nil {
		return "", err
	}

	return string(resultJSON), nil
}

// 示例 2：天气查询工具（模拟）
type WeatherTool struct{}

func (t *WeatherTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "get_weather",
		Desc: "查询指定城市的天气信息",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"city": {
				Type:     "string",
				Desc:     "城市名称，例如：北京、上海",
				Required: true,
			},
		}),
	}, nil
}

type WeatherParams struct {
	City string `json:"city"`
}

type WeatherResult struct {
	City        string `json:"city"`
	Temperature int    `json:"temperature"`
	Weather     string `json:"weather"`
	Humidity    int    `json:"humidity"`
}

func (t *WeatherTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	// 1. 解析参数
	var params WeatherParams
	if err := json.Unmarshal([]byte(argumentsInJSON), &params); err != nil {
		return "", fmt.Errorf("解析参数失败: %w", err)
	}

	fmt.Printf("🔧 [工具执行] 查询 %s 的天气...\n", params.City)

	// 2. 模拟查询天气（实际应该调用天气 API）
	result := WeatherResult{
		City:        params.City,
		Temperature: 22,
		Weather:     "晴天",
		Humidity:    60,
	}

	// 3. 返回结果
	resultJSON, err := json.Marshal(result)
	if err != nil {
		return "", err
	}

	return string(resultJSON), nil
}

// 示例 3：计算器工具
type CalculatorTool struct{}

func (t *CalculatorTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "calculator",
		Desc: "执行数学计算（加、减、乘、除）",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"operation": {
				Type:     "string",
				Desc:     "运算类型: add(加), subtract(减), multiply(乘), divide(除)",
				Required: true,
			},
			"a": {
				Type:     "number",
				Desc:     "第一个数字",
				Required: true,
			},
			"b": {
				Type:     "number",
				Desc:     "第二个数字",
				Required: true,
			},
		}),
	}, nil
}

type CalculatorParams struct {
	Operation string  `json:"operation"`
	A         float64 `json:"a"`
	B         float64 `json:"b"`
}

type CalculatorResult struct {
	Result float64 `json:"result"`
	Error  string  `json:"error,omitempty"`
}

func (t *CalculatorTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	var params CalculatorParams
	if err := json.Unmarshal([]byte(argumentsInJSON), &params); err != nil {
		return "", fmt.Errorf("解析参数失败: %w", err)
	}

	fmt.Printf("🔧 [工具执行] 计算 %g %s %g\n", params.A, params.Operation, params.B)

	var result float64
	switch params.Operation {
	case "add":
		result = params.A + params.B
	case "subtract":
		result = params.A - params.B
	case "multiply":
		result = params.A * params.B
	case "divide":
		if params.B == 0 {
			resultJSON, _ := json.Marshal(CalculatorResult{Error: "除数不能为0"})
			return string(resultJSON), nil
		}
		result = params.A / params.B
	default:
		resultJSON, _ := json.Marshal(CalculatorResult{Error: "不支持的运算"})
		return string(resultJSON), nil
	}

	resultJSON, _ := json.Marshal(CalculatorResult{Result: result})
	return string(resultJSON), nil
}

func main() {
	ctx := context.Background()

	fmt.Println("=== 工具测试 ===\n")

	// 测试 1：时间工具
	fmt.Println("【测试 1：获取当前时间】")
	timeTool := &CurrentTimeTool{}

	info, _ := timeTool.Info(ctx)
	fmt.Printf("工具名称: %s\n", info.Name)
	fmt.Printf("工具描述: %s\n", info.Desc)

	result, err := timeTool.InvokableRun(ctx, "{}")
	if err != nil {
		log.Fatalf("执行失败: %v", err)
	}
	fmt.Printf("执行结果: %s\n\n", result)

	// 测试 2：天气工具
	fmt.Println("【测试 2：查询天气】")
	weatherTool := &WeatherTool{}

	info, _ = weatherTool.Info(ctx)
	fmt.Printf("工具名称: %s\n", info.Name)
	fmt.Printf("工具描述: %s\n", info.Desc)

	params := `{"city": "北京"}`
	result, err = weatherTool.InvokableRun(ctx, params)
	if err != nil {
		log.Fatalf("执行失败: %v", err)
	}
	fmt.Printf("执行结果: %s\n\n", result)

	// 测试 3：计算器工具
	fmt.Println("【测试 3：计算器】")
	calculator := &CalculatorTool{}

	testCases := []struct {
		operation string
		a, b      float64
	}{
		{"add", 10, 5},
		{"multiply", 7, 8},
		{"divide", 100, 4},
	}

	for _, tc := range testCases {
		params := fmt.Sprintf(`{"operation":"%s","a":%g,"b":%g}`, tc.operation, tc.a, tc.b)
		result, _ := calculator.InvokableRun(ctx, params)
		fmt.Printf("结果: %s\n", result)
	}

	// 思考题：
	// 1. 工具的 Info 和 InvokableRun 分别有什么作用？
	// 2. 为什么参数和结果都要用 JSON 格式？
	// 3. 如何实现一个查询数据库的工具？
}
