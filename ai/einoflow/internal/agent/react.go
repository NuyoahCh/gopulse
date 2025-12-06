package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

// ToolExecutor 工具执行器接口
type ToolExecutor interface {
	Execute(ctx context.Context, toolName string, params map[string]interface{}) (string, error)
}

// ReActAgent ReAct 模式的 Agent（简化但完整的实现）
type ReActAgent struct {
	chatModel    model.ChatModel
	toolExecutor ToolExecutor
	maxSteps     int
}

// NewReActAgent 创建 ReAct Agent
func NewReActAgent(chatModel model.ChatModel) *ReActAgent {
	return &ReActAgent{
		chatModel: chatModel,
		maxSteps:  10,
	}
}

// SetToolExecutor 设置工具执行器
func (a *ReActAgent) SetToolExecutor(executor ToolExecutor) *ReActAgent {
	a.toolExecutor = executor
	return a
}

// Run 执行 Agent
func (a *ReActAgent) Run(ctx context.Context, task string) (string, error) {
	// 检查是否是天气查询
	isWeather := a.isWeatherQuery(task)
	fmt.Printf("[Agent] Task: %s, IsWeatherQuery: %v, HasToolExecutor: %v\n", task, isWeather, a.toolExecutor != nil)

	if isWeather {
		result, err := a.handleWeatherQuery(ctx, task)
		fmt.Printf("[Agent] Weather query result: %s, error: %v\n", result, err)
		return result, err
	}

	// 构建系统提示词
	systemPrompt := `你是一个智能助手，可以帮助用户完成各种任务。
请仔细分析用户的问题，给出详细和有帮助的回答。

如果用户询问天气信息，请告诉用户你可以查询实时天气。`

	messages := []*schema.Message{
		schema.SystemMessage(systemPrompt),
		schema.UserMessage(task),
	}

	resp, err := a.chatModel.Generate(ctx, messages)
	if err != nil {
		return "", fmt.Errorf("agent execution failed: %w", err)
	}

	return resp.Content, nil
}

// isWeatherQuery 判断是否是天气查询
func (a *ReActAgent) isWeatherQuery(task string) bool {
	weatherKeywords := []string{"天气", "气温", "温度", "下雨", "晴天", "阴天", "weather"}
	taskLower := strings.ToLower(task)

	for _, keyword := range weatherKeywords {
		if strings.Contains(taskLower, keyword) {
			return true
		}
	}
	return false
}

// handleWeatherQuery 处理天气查询
func (a *ReActAgent) handleWeatherQuery(ctx context.Context, task string) (string, error) {
	// 提取城市名
	location := a.extractLocation(task)
	if location == "" {
		location = "北京" // 默认城市
	}

	// 如果有工具执行器，使用工具查询
	if a.toolExecutor != nil {
		result, err := a.toolExecutor.Execute(ctx, "weather", map[string]interface{}{
			"location": location,
		})
		if err == nil {
			// 解析天气数据并生成友好的回复
			var weatherData map[string]interface{}
			if err := json.Unmarshal([]byte(result), &weatherData); err == nil {
				if desc, ok := weatherData["description"].(string); ok {
					return desc, nil
				}
			}
			return result, nil
		}
	}

	// 降级处理：提示用户查询方式
	return fmt.Sprintf("抱歉，我暂时无法获取%s的实时天气信息。建议您：\n1. 访问中国气象局官网：http://www.weather.com.cn\n2. 使用手机天气APP查询\n3. 搜索引擎搜索\"%s今日天气\"", location, location), nil
}

// extractLocation 从任务中提取城市名
func (a *ReActAgent) extractLocation(task string) string {
	// 常见城市列表
	cities := []string{
		"北京", "上海", "广州", "深圳", "杭州", "成都", "重庆", "武汉",
		"西安", "南京", "天津", "苏州", "郑州", "长沙", "沈阳", "青岛",
		"厦门", "大连", "宁波", "无锡", "福州", "济南", "哈尔滨", "长春",
	}

	for _, city := range cities {
		if strings.Contains(task, city) {
			return city
		}
	}

	return ""
}

// RunWithTools 使用工具执行（未来扩展）
func (a *ReActAgent) RunWithTools(ctx context.Context, task string, toolsDesc string) (string, error) {
	systemPrompt := fmt.Sprintf(`你是一个智能助手，可以使用以下工具来帮助用户：

%s

请根据用户的任务，思考是否需要使用工具，并给出详细的回答。`, toolsDesc)

	messages := []*schema.Message{
		schema.SystemMessage(systemPrompt),
		schema.UserMessage(task),
	}

	resp, err := a.chatModel.Generate(ctx, messages)
	if err != nil {
		return "", fmt.Errorf("agent execution failed: %w", err)
	}

	return resp.Content, nil
}

// SetMaxSteps 设置最大步骤数
func (a *ReActAgent) SetMaxSteps(maxSteps int) *ReActAgent {
	a.maxSteps = maxSteps
	return a
}
