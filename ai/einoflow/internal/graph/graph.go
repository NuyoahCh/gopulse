package graph

import (
	"context"
	"fmt"
)

// NodeHandler 节点处理函数
type NodeHandler func(ctx context.Context, input interface{}) (interface{}, error)

// Node 图节点
type Node struct {
	ID      string
	Handler NodeHandler
}

// Edge 图边
type Edge struct {
	From string
	To   string
}

// ProgressCallback 进度回调函数
type ProgressCallback func(nodeID string, input, output interface{})

// Graph 图编排
type Graph struct {
	nodes            map[string]*Node
	edges            []*Edge
	start            string
	end              string
	progressCallback ProgressCallback
}

// NewGraph 创建图
func NewGraph() *Graph {
	return &Graph{
		nodes:            make(map[string]*Node),
		edges:            make([]*Edge, 0),
		progressCallback: nil,
	}
}

// SetProgressCallback 设置进度回调
func (g *Graph) SetProgressCallback(callback ProgressCallback) *Graph {
	g.progressCallback = callback
	return g
}

// AddNode 添加节点
func (g *Graph) AddNode(id string, handler NodeHandler) *Graph {
	g.nodes[id] = &Node{
		ID:      id,
		Handler: handler,
	}
	return g
}

// AddEdge 添加边
func (g *Graph) AddEdge(from, to string) *Graph {
	g.edges = append(g.edges, &Edge{
		From: from,
		To:   to,
	})
	return g
}

// SetStart 设置起始节点
func (g *Graph) SetStart(nodeID string) *Graph {
	g.start = nodeID
	return g
}

// SetEnd 设置结束节点
func (g *Graph) SetEnd(nodeID string) *Graph {
	g.end = nodeID
	return g
}

// Run 执行图
func (g *Graph) Run(ctx context.Context, input interface{}) (interface{}, error) {
	if g.start == "" {
		return nil, fmt.Errorf("start node not set")
	}

	// 简单的顺序执行（实际应该根据边的关系执行）
	visited := make(map[string]bool)
	return g.executeNode(ctx, g.start, input, visited)
}

func (g *Graph) executeNode(ctx context.Context, nodeID string, input interface{}, visited map[string]bool) (interface{}, error) {
	if visited[nodeID] {
		return input, nil // 避免循环
	}
	visited[nodeID] = true

	node, ok := g.nodes[nodeID]
	if !ok {
		return nil, fmt.Errorf("node %s not found", nodeID)
	}

	// 执行节点
	output, err := node.Handler(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("node %s execution failed: %w", nodeID, err)
	}

	// 调用进度回调
	if g.progressCallback != nil {
		g.progressCallback(nodeID, input, output)
	}

	// 如果是结束节点，返回结果
	if nodeID == g.end {
		return output, nil
	}

	// 找到下一个节点
	for _, edge := range g.edges {
		if edge.From == nodeID {
			return g.executeNode(ctx, edge.To, output, visited)
		}
	}

	return output, nil
}

// Validate 验证图结构
func (g *Graph) Validate() error {
	if g.start == "" {
		return fmt.Errorf("start node not set")
	}
	if _, ok := g.nodes[g.start]; !ok {
		return fmt.Errorf("start node %s not found", g.start)
	}
	if g.end != "" {
		if _, ok := g.nodes[g.end]; !ok {
			return fmt.Errorf("end node %s not found", g.end)
		}
	}
	return nil
}
