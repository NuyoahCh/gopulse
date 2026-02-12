package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// WeatherTool 天气查询工具
type WeatherTool struct {
	client *http.Client
}

// NewWeatherTool 创建天气工具
func NewWeatherTool() *WeatherTool {
	return &WeatherTool{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// WeatherResponse 天气响应结构
type WeatherResponse struct {
	Location    string `json:"location"`
	Temperature string `json:"temperature"`
	Condition   string `json:"condition"`
	Humidity    string `json:"humidity"`
	WindSpeed   string `json:"wind_speed"`
	FeelsLike   string `json:"feels_like"`
	Description string `json:"description"`
}

// GetWeather 获取天气信息
// @tool Get current weather for a location
func (t *WeatherTool) GetWeather(ctx context.Context, location string) (string, error) {
	// 使用 wttr.in 免费天气服务（支持中文城市名）
	// 格式：wttr.in/城市名?format=j1（返回 JSON）
	url := fmt.Sprintf("https://wttr.in/%s?format=j1&lang=zh", location)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return t.getFallbackWeather(location), nil
	}

	resp, err := t.client.Do(req)
	if err != nil {
		return t.getFallbackWeather(location), nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return t.getFallbackWeather(location), nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return t.getFallbackWeather(location), nil
	}

	// 解析 wttr.in 的 JSON 响应
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return t.getFallbackWeather(location), nil
	}

	// 提取天气信息
	weather := t.parseWeatherData(location, result)

	data, err := json.Marshal(weather)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// parseWeatherData 解析天气数据
func (t *WeatherTool) parseWeatherData(location string, data map[string]interface{}) WeatherResponse {
	weather := WeatherResponse{
		Location: location,
	}

	// 获取当前天气
	if currentCondition, ok := data["current_condition"].([]interface{}); ok && len(currentCondition) > 0 {
		if current, ok := currentCondition[0].(map[string]interface{}); ok {
			if temp, ok := current["temp_C"].(string); ok {
				weather.Temperature = temp + "°C"
			}
			if feelsLike, ok := current["FeelsLikeC"].(string); ok {
				weather.FeelsLike = feelsLike + "°C"
			}
			if humidity, ok := current["humidity"].(string); ok {
				weather.Humidity = humidity + "%"
			}
			if windSpeed, ok := current["windspeedKmph"].(string); ok {
				weather.WindSpeed = windSpeed + " km/h"
			}

			// 获取天气描述（中文）
			if weatherDesc, ok := current["lang_zh"].([]interface{}); ok && len(weatherDesc) > 0 {
				if desc, ok := weatherDesc[0].(map[string]interface{}); ok {
					if value, ok := desc["value"].(string); ok {
						weather.Condition = value
					}
				}
			}
		}
	}

	// 生成友好的描述
	weather.Description = fmt.Sprintf("%s当前天气：%s，温度%s，体感温度%s，湿度%s，风速%s",
		location, weather.Condition, weather.Temperature, weather.FeelsLike, weather.Humidity, weather.WindSpeed)

	return weather
}

// getFallbackWeather 获取备用天气信息（当 API 调用失败时）
func (t *WeatherTool) getFallbackWeather(location string) string {
	// 提供友好的错误提示
	weather := WeatherResponse{
		Location:    location,
		Temperature: "无法获取",
		Condition:   "数据不可用",
		Description: fmt.Sprintf("抱歉，暂时无法获取%s的实时天气信息。建议您：\n1. 访问中国气象局官网：http://www.weather.com.cn\n2. 使用手机天气APP查询\n3. 搜索引擎搜索\"%s今日天气\"", location, location),
	}

	data, _ := json.Marshal(weather)
	return string(data)
}

// GetWeatherDescription 获取天气描述（用于 Agent 工具调用）
func (t *WeatherTool) GetWeatherDescription() string {
	return "获取指定城市的实时天气信息，包括温度、天气状况、湿度、风速等。支持中文城市名，如：北京、上海、深圳等。"
}

// GetWeatherParameters 获取工具参数定义
func (t *WeatherTool) GetWeatherParameters() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"location": map[string]interface{}{
				"type":        "string",
				"description": "城市名称，如：北京、上海、深圳",
			},
		},
		"required": []string{"location"},
	}
}
