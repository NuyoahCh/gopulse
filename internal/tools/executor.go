package tools

import (
	"context"
	"fmt"
)

// Executor 工具执行器
type Executor struct {
	registry *Registry
}

// NewExecutor 创建工具执行器
func NewExecutor(registry *Registry) *Executor {
	return &Executor{
		registry: registry,
	}
}

// Execute 执行工具
func (e *Executor) Execute(ctx context.Context, toolName string, params map[string]interface{}) (string, error) {
	fmt.Printf("[ToolExecutor] Executing tool: %s, params: %v\n", toolName, params)

	tool, ok := e.registry.Get(toolName)
	if !ok {
		err := fmt.Errorf("tool not found: %s", toolName)
		fmt.Printf("[ToolExecutor] Error: %v\n", err)
		return "", err
	}

	switch toolName {
	case "weather":
		if weatherTool, ok := tool.(*WeatherTool); ok {
			location, _ := params["location"].(string)
			if location == "" {
				location = "北京"
			}
			fmt.Printf("[ToolExecutor] Calling weather tool for location: %s\n", location)
			result, err := weatherTool.GetWeather(ctx, location)
			fmt.Printf("[ToolExecutor] Weather result: %s, error: %v\n", result, err)
			return result, err
		}
	case "calculator":
		if calcTool, ok := tool.(*CalculatorTool); ok {
			expression, _ := params["expression"].(string)
			return calcTool.Calculate(ctx, expression)
		}
	case "search":
		if searchTool, ok := tool.(*SearchTool); ok {
			query, _ := params["query"].(string)
			return searchTool.Search(ctx, query)
		}
	}

	err := fmt.Errorf("tool execution not implemented: %s", toolName)
	fmt.Printf("[ToolExecutor] Error: %v\n", err)
	return "", err
}
