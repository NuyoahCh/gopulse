package rag

import (
	"strings"

	"github.com/cloudwego/eino/schema"
)

// TextSplitter 文本分割器
type TextSplitter struct {
	chunkSize    int
	chunkOverlap int
}

// NewTextSplitter 创建文本分割器
func NewTextSplitter(chunkSize, chunkOverlap int) *TextSplitter {
	return &TextSplitter{
		chunkSize:    chunkSize,
		chunkOverlap: chunkOverlap,
	}
}

// Split 分割文档
func (s *TextSplitter) Split(doc *schema.Document) []*schema.Document {
	text := doc.Content
	chunks := make([]*schema.Document, 0)

	// 简单的按字符分割
	for i := 0; i < len(text); i += s.chunkSize - s.chunkOverlap {
		end := i + s.chunkSize
		if end > len(text) {
			end = len(text)
		}

		chunk := &schema.Document{
			Content:  text[i:end],
			MetaData: doc.MetaData,
		}
		chunks = append(chunks, chunk)

		if end >= len(text) {
			break
		}
	}

	return chunks
}

// SplitByParagraph 按段落分割
func (s *TextSplitter) SplitByParagraph(doc *schema.Document) []*schema.Document {
	paragraphs := strings.Split(doc.Content, "\n\n")
	chunks := make([]*schema.Document, 0, len(paragraphs))

	for _, para := range paragraphs {
		if strings.TrimSpace(para) == "" {
			continue
		}

		chunk := &schema.Document{
			Content:  para,
			MetaData: doc.MetaData,
		}
		chunks = append(chunks, chunk)
	}

	return chunks
}
