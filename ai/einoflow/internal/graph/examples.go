package graph

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

// CreateRAGGraph 创建 RAG 图
func CreateRAGGraph(chatModel model.ChatModel, retriever NodeHandler) *Graph {
	graph := NewGraph()

	// 节点1: 检索相关文档
	retrieveNode := func(ctx context.Context, input interface{}) (interface{}, error) {
		query := input.(string)
		docs, err := retriever(ctx, query)
		if err != nil {
			return nil, err
		}
		return map[string]interface{}{
			"query": query,
			"docs":  docs,
		}, nil
	}

	// 节点2: 生成答案
	generateNode := func(ctx context.Context, input interface{}) (interface{}, error) {
		data := input.(map[string]interface{})
		query := data["query"].(string)
		docs := data["docs"]

		// 构建提示词
		prompt := fmt.Sprintf("基于以下文档回答问题：\n\n文档：%v\n\n问题：%s", docs, query)

		messages := []*schema.Message{
			schema.UserMessage(prompt),
		}

		resp, err := chatModel.Generate(ctx, messages)
		if err != nil {
			return "", err
		}

		return resp.Content, nil
	}

	graph.AddNode("retrieve", retrieveNode).
		AddNode("generate", generateNode).
		AddEdge("retrieve", "generate").
		SetStart("retrieve").
		SetEnd("generate")

	return graph
}

// CreateMultiStepGraph 创建多步骤处理图
func CreateMultiStepGraph(chatModel model.ChatModel) *Graph {
	graph := NewGraph()

	// 步骤1: 分析问题
	analyzeNode := func(ctx context.Context, input interface{}) (interface{}, error) {
		question := input.(string)
		messages := []*schema.Message{
			schema.SystemMessage("你是一个问题分析专家，请分析用户问题的意图和关键点。"),
			schema.UserMessage(question),
		}
		resp, err := chatModel.Generate(ctx, messages)
		if err != nil {
			return "", err
		}
		return resp.Content, nil
	}

	// 步骤2: 制定计划
	planNode := func(ctx context.Context, input interface{}) (interface{}, error) {
		analysis := input.(string)
		messages := []*schema.Message{
			schema.SystemMessage("根据问题分析，制定解决方案的步骤。"),
			schema.UserMessage(analysis),
		}
		resp, err := chatModel.Generate(ctx, messages)
		if err != nil {
			return "", err
		}
		return resp.Content, nil
	}

	// 步骤3: 执行并总结
	executeNode := func(ctx context.Context, input interface{}) (interface{}, error) {
		plan := input.(string)
		messages := []*schema.Message{
			schema.SystemMessage("根据计划执行并给出最终答案。"),
			schema.UserMessage(plan),
		}
		resp, err := chatModel.Generate(ctx, messages)
		if err != nil {
			return "", err
		}
		return resp.Content, nil
	}

	graph.AddNode("analyze", analyzeNode).
		AddNode("plan", planNode).
		AddNode("execute", executeNode).
		AddEdge("analyze", "plan").
		AddEdge("plan", "execute").
		SetStart("analyze").
		SetEnd("execute")

	return graph
}
