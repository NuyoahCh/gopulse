package memory

import (
	"context"
	"sync"

	"github.com/cloudwego/eino/schema"
)

// ChatHistory 聊天历史
type ChatHistory struct {
	mu       sync.RWMutex
	messages []*schema.Message
	maxSize  int
}

// NewChatHistory 创建聊天历史
func NewChatHistory(maxSize int) *ChatHistory {
	return &ChatHistory{
		messages: make([]*schema.Message, 0),
		maxSize:  maxSize,
	}
}

// Add 添加消息
func (h *ChatHistory) Add(ctx context.Context, message *schema.Message) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.messages = append(h.messages, message)

	// 保持最大大小
	if len(h.messages) > h.maxSize {
		h.messages = h.messages[len(h.messages)-h.maxSize:]
	}
}

// GetAll 获取所有消息
func (h *ChatHistory) GetAll(ctx context.Context) []*schema.Message {
	h.mu.RLock()
	defer h.mu.RUnlock()

	// 返回副本
	messages := make([]*schema.Message, len(h.messages))
	copy(messages, h.messages)
	return messages
}

// GetLast 获取最后 n 条消息
func (h *ChatHistory) GetLast(ctx context.Context, n int) []*schema.Message {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if n > len(h.messages) {
		n = len(h.messages)
	}

	messages := make([]*schema.Message, n)
	copy(messages, h.messages[len(h.messages)-n:])
	return messages
}

// Clear 清空历史
func (h *ChatHistory) Clear(ctx context.Context) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.messages = make([]*schema.Message, 0)
}

// Size 获取消息数量
func (h *ChatHistory) Size() int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return len(h.messages)
}
