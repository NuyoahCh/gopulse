package main

import (
	"context"
	"fmt"
	"log"

	"einoflow/internal/rag"

	embeddingopenai "github.com/cloudwego/eino-ext/components/embedding/openai"
	modelopenai "github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/schema"
)

func main() {
	ctx := context.Background()

	// 1. 创建 Embedding 模型
	embedder, err := embeddingopenai.NewEmbedder(ctx, &embeddingopenai.EmbeddingConfig{
		APIKey: "your-api-key",
		Model:  "text-embedding-3-small",
	})
	if err != nil {
		log.Fatal(err)
	}

	// 2. 创建向量存储
	vectorStore := rag.NewMemoryVectorStore()

	// 3. 加载文档
	loader := rag.NewDocumentLoader()
	doc := loader.LoadFromText(ctx, "Eino 是字节跳动开源的 LLM 应用开发框架，支持 Chain、Agent、RAG 等功能。", map[string]interface{}{
		"source": "intro",
	})

	// 4. 文本分割
	splitter := rag.NewTextSplitter(500, 50)
	chunks := splitter.Split(doc)

	// 5. 向量化并存储
	for _, chunk := range chunks {
		embeddings, err := embedder.EmbedStrings(ctx, []string{chunk.Content})
		if err != nil {
			log.Fatal(err)
		}
		vectorStore.Add(ctx, []*schema.Document{chunk}, embeddings)
	}

	// 6. 创建检索器
	retriever := rag.NewVectorRetriever(embedder, vectorStore, 3)

	// 7. 检索相关文档
	query := "Eino 支持哪些功能？"
	docs, err := retriever.Retrieve(ctx, query)
	if err != nil {
		log.Fatal(err)
	}

	// 8. 使用 LLM 生成答案
	chatModel, err := modelopenai.NewChatModel(ctx, &modelopenai.ChatModelConfig{
		APIKey: "your-api-key",
		Model:  "gpt-4",
	})
	if err != nil {
		log.Fatal(err)
	}

	// 构建上下文
	contextStr := ""
	for _, doc := range docs {
		contextStr += doc.Content + "\n"
	}

	messages := []*schema.Message{
		schema.SystemMessage("根据以下上下文回答问题：\n" + contextStr),
		schema.UserMessage(query),
	}

	resp, err := chatModel.Generate(ctx, messages)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("问题: %s\n", query)
	fmt.Printf("答案: %s\n", resp.Content)
}
