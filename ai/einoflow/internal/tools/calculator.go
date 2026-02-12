package tools

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

// CalculatorTool 计算器工具
type CalculatorTool struct{}

// NewCalculatorTool 创建计算器工具
func NewCalculatorTool() *CalculatorTool {
	return &CalculatorTool{}
}

// Calculate 执行计算
// @tool Perform basic arithmetic calculations
func (t *CalculatorTool) Calculate(ctx context.Context, expression string) (string, error) {
	// 简单的计算实现
	parts := strings.Fields(expression)
	if len(parts) != 3 {
		return "", fmt.Errorf("invalid expression format")
	}

	a, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return "", err
	}

	b, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return "", err
	}

	var result float64
	switch parts[1] {
	case "+":
		result = a + b
	case "-":
		result = a - b
	case "*":
		result = a * b
	case "/":
		if b == 0 {
			return "", fmt.Errorf("division by zero")
		}
		result = a / b
	default:
		return "", fmt.Errorf("unsupported operator: %s", parts[1])
	}

	return fmt.Sprintf("%.2f", result), nil
}
