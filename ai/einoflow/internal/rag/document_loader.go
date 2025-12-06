package rag

import (
	"context"
	"os"

	"github.com/cloudwego/eino/schema"
)

// DocumentLoader 文档加载器
type DocumentLoader struct{}

// NewDocumentLoader 创建文档加载器
func NewDocumentLoader() *DocumentLoader {
	return &DocumentLoader{}
}

// LoadFromFile 从文件加载文档
func (l *DocumentLoader) LoadFromFile(ctx context.Context, filepath string) ([]*schema.Document, error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	doc := &schema.Document{
		Content: string(content),
		MetaData: map[string]interface{}{
			"source": filepath,
		},
	}

	return []*schema.Document{doc}, nil
}

// LoadFromText 从文本加载文档
func (l *DocumentLoader) LoadFromText(ctx context.Context, text string, metadata map[string]interface{}) *schema.Document {
	return &schema.Document{
		Content:  text,
		MetaData: metadata,
	}
}
