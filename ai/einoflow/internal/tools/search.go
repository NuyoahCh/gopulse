package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// SearchTool 搜索工具
type SearchTool struct {
	apiKey string
}

// NewSearchTool 创建搜索工具
func NewSearchTool(apiKey string) *SearchTool {
	return &SearchTool{apiKey: apiKey}
}

// Search 执行网络搜索
// @tool Search the web for information
func (t *SearchTool) Search(ctx context.Context, query string) (string, error) {
	// 使用 DuckDuckGo 即时答案 API（无需 API Key）
	baseURL := "https://api.duckduckgo.com/"
	params := url.Values{}
	params.Add("q", query)
	params.Add("format", "json")
	params.Add("no_html", "1")
	params.Add("skip_disambig", "1")

	resp, err := http.Get(baseURL + "?" + params.Encode())
	if err != nil {
		return "", fmt.Errorf("search failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	// 提取摘要
	if abstract, ok := result["AbstractText"].(string); ok && abstract != "" {
		return abstract, nil
	}

	// 如果没有摘要，返回相关主题
	if relatedTopics, ok := result["RelatedTopics"].([]interface{}); ok && len(relatedTopics) > 0 {
		if topic, ok := relatedTopics[0].(map[string]interface{}); ok {
			if text, ok := topic["Text"].(string); ok {
				return text, nil
			}
		}
	}

	return "No results found", nil
}
