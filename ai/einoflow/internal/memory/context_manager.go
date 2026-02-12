package memory

import (
	"fmt"

	"github.com/cloudwego/eino/schema"
	"github.com/pkoukk/tiktoken-go"
)

// ContextManager 上下文窗口管理器
type ContextManager struct {
	maxTokens     int
	tokenizer     *tiktoken.Tiktoken
	reserveTokens int // 为响应预留的 token 数
}

// NewContextManager 创建上下文管理器
func NewContextManager(maxTokens int) (*ContextManager, error) {
	// 使用 cl100k_base 编码器（适用于 GPT-4 和豆包）
	tokenizer, err := tiktoken.GetEncoding("cl100k_base")
	if err != nil {
		return nil, fmt.Errorf("failed to get tokenizer: %w", err)
	}

	return &ContextManager{
		maxTokens:     maxTokens,
		tokenizer:     tokenizer,
		reserveTokens: 1000, // 为响应预留 1000 tokens
	}, nil
}

// CountTokens 计算消息的 token 数
func (cm *ContextManager) CountTokens(messages []*schema.Message) int {
	totalTokens := 0
	for _, msg := range messages {
		// 每条消息的开销：role + content + 格式化
		tokens := cm.tokenizer.Encode(msg.Content, nil, nil)
		totalTokens += len(tokens) + 4 // +4 for message formatting
	}
	return totalTokens + 3 // +3 for reply priming
}

// CountTextTokens 计算文本的 token 数
func (cm *ContextManager) CountTextTokens(text string) int {
	tokens := cm.tokenizer.Encode(text, nil, nil)
	return len(tokens)
}

// TruncateMessages 截断消息以适应上下文窗口
func (cm *ContextManager) TruncateMessages(messages []*schema.Message) []*schema.Message {
	if len(messages) == 0 {
		return messages
	}

	// 计算当前 token 数
	currentTokens := cm.CountTokens(messages)
	maxAllowed := cm.maxTokens - cm.reserveTokens

	if currentTokens <= maxAllowed {
		return messages // 不需要截断
	}

	// 保留系统消息（如果有）
	var systemMsg *schema.Message
	startIdx := 0
	if len(messages) > 0 && messages[0].Role == "system" {
		systemMsg = messages[0]
		startIdx = 1
	}

	// 从最新的消息开始保留
	result := make([]*schema.Message, 0)
	if systemMsg != nil {
		result = append(result, systemMsg)
	}

	// 从后往前添加消息，直到达到 token 限制
	tokens := 0
	if systemMsg != nil {
		tokens = cm.CountTokens([]*schema.Message{systemMsg})
	}

	for i := len(messages) - 1; i >= startIdx; i-- {
		msgTokens := cm.CountTokens([]*schema.Message{messages[i]})
		if tokens+msgTokens > maxAllowed {
			break
		}
		result = append([]*schema.Message{messages[i]}, result...)
		tokens += msgTokens
	}

	return result
}

// TruncateText 截断文本以适应 token 限制
func (cm *ContextManager) TruncateText(text string, maxTokens int) string {
	tokens := cm.tokenizer.Encode(text, nil, nil)
	if len(tokens) <= maxTokens {
		return text
	}

	// 截断 tokens 并解码
	truncatedTokens := tokens[:maxTokens]
	return cm.tokenizer.Decode(truncatedTokens)
}

// GetMaxTokens 获取最大 token 数
func (cm *ContextManager) GetMaxTokens() int {
	return cm.maxTokens
}

// GetAvailableTokens 获取可用的 token 数
func (cm *ContextManager) GetAvailableTokens(messages []*schema.Message) int {
	used := cm.CountTokens(messages)
	available := cm.maxTokens - used - cm.reserveTokens
	if available < 0 {
		return 0
	}
	return available
}

// SetReserveTokens 设置预留的 token 数
func (cm *ContextManager) SetReserveTokens(tokens int) {
	cm.reserveTokens = tokens
}

// GetReserveTokens 获取预留的 token 数
func (cm *ContextManager) GetReserveTokens() int {
	return cm.reserveTokens
}
