package chain

import (
	"context"
	"sync"
)

// ParallelChain 并行链（简化版本，用于演示）
// 注意：这是一个简化实现，实际使用应该用 eino 的 Graph API
type ParallelChain struct {
	branches []func(context.Context, interface{}) (interface{}, error)
}

// NewParallelChain 创建并行链
func NewParallelChain() *ParallelChain {
	return &ParallelChain{
		branches: make([]func(context.Context, interface{}) (interface{}, error), 0),
	}
}

// AddBranch 添加分支
func (c *ParallelChain) AddBranch(branch func(context.Context, interface{}) (interface{}, error)) *ParallelChain {
	c.branches = append(c.branches, branch)
	return c
}

// Run 并行执行所有分支
func (c *ParallelChain) Run(ctx context.Context, input interface{}) ([]interface{}, error) {
	results := make([]interface{}, len(c.branches))
	errors := make([]error, len(c.branches))

	var wg sync.WaitGroup
	wg.Add(len(c.branches))

	for i, branch := range c.branches {
		go func(idx int, b func(context.Context, interface{}) (interface{}, error)) {
			defer wg.Done()
			result, err := b(ctx, input)
			results[idx] = result
			errors[idx] = err
		}(i, branch)
	}

	wg.Wait()

	// 检查错误
	for _, err := range errors {
		if err != nil {
			return nil, err
		}
	}

	return results, nil
}
