package main

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/cloudwego/eino-ext/components/indexer/milvus"
	"github.com/cloudwego/eino/schema"
	"github.com/joho/godotenv"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Failed to load .env file")
	}

	ctx := context.Background()
	// model, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
	// 	APIKey: os.Getenv("ARK_API_KEY"),
	// 	Model:  os.Getenv("ARK_CHAT_MODEL"),
	// })

	// template := prompt.FromMessages(schema.FString,
	// 	schema.SystemMessage("你是一个{role}"),
	// 	&schema.Message{
	// 		Role:    schema.User,
	// 		Content: "提醒我每天{task}",
	// 	},
	// )

	// params := map[string]any{
	// 	"role": "爱睡觉的老爸",
	// 	"task": "好好学习，天天向上",
	// }

	// input := []*schema.Message{
	// 	schema.SystemMessage("你是一个爱唠叨的老妈"),
	// 	schema.UserMessage("提醒我每天早睡早起"),
	// }

	// message, err := template.Format(ctx, params)

	// response, err := model.Generate(ctx, message)
	// if err != nil {
	// 	panic(err)
	// }
	// print(response.Content)

	// respose, err := model.Generate(ctx, input)
	// if err != nil {
	// 	panic(err)
	// }
	// print(respose.Content)

	// reader, err := model.Stream(ctx, input)
	// if err != nil {
	// 	panic(err)
	// }
	// defer reader.Close()

	// for {
	// 	chunk, err := reader.Recv()
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	print(chunk.Content)
	// }

	// ctx := context.Background()

	embedder, err := ark.NewEmbedder(ctx, &ark.EmbeddingConfig{
		APIKey: os.Getenv("ARK_API_KEY"),
		Model:  os.Getenv("ARK_EMBEDDING_MODEL"),
	})

	if err != nil {
		panic(err)
	}

	// // 生成文本向量
	// texts := []string{
	// 	"这是第一段示例文本",
	// 	"这是第二段示例文本",
	// }

	// embeddings, err := embedder.EmbedStrings(ctx, texts)
	// if err != nil {
	// 	panic(err)
	// }

	// // 使用生成的向量
	// for i, embedding := range embeddings {
	// 	println("文本", i+1, "的向量维度:", len(embedding))
	// }
	// println("生成的向量数量:", len(embeddings))

	InitClient()

	var collection = "test"

	var fields = []*entity.Field{
		{
			Name:     "id",
			DataType: entity.FieldTypeVarChar,
			TypeParams: map[string]string{
				"max_length": "256",
			},
			PrimaryKey: true,
		},
		{
			Name:     "vector", // 确保字段名匹配
			DataType: entity.FieldTypeBinaryVector,
			TypeParams: map[string]string{
				"dim": "81920",
			},
		},
		{
			Name:     "content",
			DataType: entity.FieldTypeVarChar,
			TypeParams: map[string]string{
				"max_length": "8192",
			},
		},
		{
			Name:     "metadata",
			DataType: entity.FieldTypeJSON,
		},
	}

	indexer, err := milvus.NewIndexer(ctx, &milvus.IndexerConfig{
		Client:     MilvusCli,
		Collection: collection,
		Fields:     fields,
		Embedding:  embedder,
	})

	if err != nil {
		panic(err)
	}

	docs := []*schema.Document{
		{
			ID:      "1",
			Content: "这是一个关于人工智能的文档。",
			MetaData: map[string]any{
				"category": "technology",
			},
		},

		{
			ID:      "2",
			Content: "今天的天气非常好，适合户外活动。",
			MetaData: map[string]any{
				"category": "weather",
			},
		},
	}

	ids, err := indexer.Store(ctx, docs)
	if err != nil {
		panic(err)
	}

	fmt.Println("Stored document IDs:", ids)

}
