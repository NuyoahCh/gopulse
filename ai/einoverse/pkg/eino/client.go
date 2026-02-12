package eino

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Client Eino 客户端
type Client struct {
	apiKey  string
	apiBase string
	model   string
	client  *http.Client
}

// Message 消息结构
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// NewClient 创建 Eino 客户端
func NewClient() *Client {
	// 获取环境变量
	apiKey := os.Getenv("ARK_API_KEY")
	apiBase := os.Getenv("ARK_API_BASE_URL")
	// 如果API基础URL为空，则使用默认URL
	if apiBase == "" {
		apiBase = "https://ark.cn-beijing.volces.com/api/v3"
	}

	model := os.Getenv("ARK_MODEL_NAME")
	// 如果模型名称为空，则使用默认模型
	if model == "" {
		model = "doubao-seed-1-6-lite-251015"
	}

	// 创建客户端
	return &Client{
		apiKey:  apiKey,
		apiBase: apiBase,
		model:   model,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// ChatCompletion 调用聊天补全接口
func (c *Client) ChatCompletion(messages []Message) (string, error) {
	if c.apiKey == "" {
		return "", fmt.Errorf("ARK_API_KEY not configured")
	}

	// 构建请求体
	reqBody := map[string]interface{}{
		"model":    c.model,
		"messages": messages,
	}

	// 将请求体转换为JSON
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request failed: %w", err)
	}

	// 创建请求
	req, err := http.NewRequest("POST", c.apiBase+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("create request failed: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	// 发送请求
	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// 如果响应状态码不为200，则返回错误
	if resp.StatusCode != http.StatusOK {
		// 读取响应体
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API error [%d]: %s", resp.StatusCode, string(body))
	}

	// 解析响应体
	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}

	// 解析响应体
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode response failed: %w", err)
	}

	// 如果响应体有错误，则返回错误
	if result.Error.Message != "" {
		return "", fmt.Errorf("API error: %s", result.Error.Message)
	}

	// 如果响应体没有选择，则返回错误
	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no response from model")
	}

	// 返回响应体
	return result.Choices[0].Message.Content, nil
}

// IsAvailable 检查服务是否可用
func (c *Client) IsAvailable() bool {
	return c.apiKey != ""
}
