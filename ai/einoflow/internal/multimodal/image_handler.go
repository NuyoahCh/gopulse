package multimodal

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/cloudwego/eino/schema"
)

// ImageHandler 图像处理器
type ImageHandler struct {
	maxImageSize int64 // 最大图像大小（字节）
}

// NewImageHandler 创建图像处理器
func NewImageHandler() *ImageHandler {
	return &ImageHandler{
		maxImageSize: 20 * 1024 * 1024, // 20MB
	}
}

// LoadImageFromURL 从 URL 加载图像
func (h *ImageHandler) LoadImageFromURL(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch image: status %d", resp.StatusCode)
	}

	// 读取图像数据
	data, err := io.ReadAll(io.LimitReader(resp.Body, h.maxImageSize))
	if err != nil {
		return "", fmt.Errorf("failed to read image: %w", err)
	}

	// 转换为 base64
	return base64.StdEncoding.EncodeToString(data), nil
}

// LoadImageFromFile 从文件加载图像
func (h *ImageHandler) LoadImageFromFile(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	if int64(len(data)) > h.maxImageSize {
		return "", fmt.Errorf("image too large: %d bytes (max %d)", len(data), h.maxImageSize)
	}

	return base64.StdEncoding.EncodeToString(data), nil
}

// CreateImageMessage 创建包含图像的消息
// 注意：当前版本使用简化实现，将图像 URL 附加到文本内容中
func (h *ImageHandler) CreateImageMessage(role, text, imageURL string) *schema.Message {
	// 简化实现：将图像信息附加到文本中
	content := fmt.Sprintf("%s\n[Image: %s]", text, imageURL)
	return &schema.Message{
		Role:    schema.RoleType(role),
		Content: content,
	}
}

// CreateImageMessageFromBase64 创建包含 base64 图像的消息
func (h *ImageHandler) CreateImageMessageFromBase64(role, text, base64Data, mimeType string) *schema.Message {
	// 简化实现：创建 data URL 并附加到文本中
	dataURL := fmt.Sprintf("data:%s;base64,%s", mimeType, base64Data[:min(100, len(base64Data))])
	content := fmt.Sprintf("%s\n[Image Data: %s...]", text, dataURL)
	return &schema.Message{
		Role:    schema.RoleType(role),
		Content: content,
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// SetMaxImageSize 设置最大图像大小
func (h *ImageHandler) SetMaxImageSize(size int64) {
	h.maxImageSize = size
}

// GetMaxImageSize 获取最大图像大小
func (h *ImageHandler) GetMaxImageSize() int64 {
	return h.maxImageSize
}
