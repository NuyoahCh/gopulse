package chain

import (
	"context"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

// SequentialChain 顺序链
type SequentialChain struct {
	steps []func(context.Context, string) (string, error)
}

// NewSequentialChain 创建顺序链
func NewSequentialChain() *SequentialChain {
	return &SequentialChain{
		steps: make([]func(context.Context, string) (string, error), 0),
	}
}

// AddLambdaStep 添加 Lambda 步骤
func (c *SequentialChain) AddLambdaStep(step func(context.Context, string) (string, error)) *SequentialChain {
	c.steps = append(c.steps, step)
	return c
}

// AddStep 添加步骤（兼容旧接口）
func (c *SequentialChain) AddStep(step func(context.Context, interface{}) (interface{}, error)) *SequentialChain {
	// 包装为 string 类型
	wrapped := func(ctx context.Context, input string) (string, error) {
		result, err := step(ctx, input)
		if err != nil {
			return "", err
		}
		return result.(string), nil
	}
	c.steps = append(c.steps, wrapped)
	return c
}

// Run 执行链
func (c *SequentialChain) Run(ctx context.Context, input string) (string, error) {
	result := input
	var err error

	for _, step := range c.steps {
		result, err = step(ctx, result)
		if err != nil {
			return "", err
		}
	}

	return result, nil
}

// Example: 创建一个简单的翻译链
func CreateTranslationChain(chatModel model.ChatModel) *SequentialChain {
	chain := NewSequentialChain()

	// 步骤1: 翻译成英文
	translateToEnglish := func(ctx context.Context, input interface{}) (interface{}, error) {
		text := input.(string)
		messages := []*schema.Message{
			schema.UserMessage("Translate to English: " + text),
		}
		resp, err := chatModel.Generate(ctx, messages)
		if err != nil {
			return "", err
		}
		return resp.Content, nil
	}

	// 步骤2: 总结
	summarize := func(ctx context.Context, input interface{}) (interface{}, error) {
		text := input.(string)
		messages := []*schema.Message{
			schema.UserMessage("Summarize in one sentence: " + text),
		}
		resp, err := chatModel.Generate(ctx, messages)
		if err != nil {
			return "", err
		}
		return resp.Content, nil
	}

	chain.AddStep(translateToEnglish).AddStep(summarize)
	return chain
}
