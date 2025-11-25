package singleflight

import (
	"context"
	"sync"
)

// Result 是一次调用返回的结果结构体。
type Result struct {
	Val    any   // 业务返回值
	Err    error // 业务错误
	Shared bool  // 是否为“复用”的结果
}

// Group 用于管理同一 key 的并发调用，使其在同一时间只会执行一次。
type Group struct {
	mu sync.Mutex
	m  map[string]*call
}

// call 代表一次正在进行中的调用。
type call struct {
	wg        sync.WaitGroup  // 等待调用结束
	val       any             // 业务返回值
	err       error           // 业务错误
	dups      int             // 复用次数
	chans     []chan<- Result // DoChan/DoContext 这种通过 channel 等待结果的调用者列表
	forgotten bool            // 是否被 Forget 标记（避免结束时重复 delete）
}

// Do 确保同一个 key 的 fn 在同一时间只会被执行一次。
func (g *Group) Do(key string, fn func() (any, error)) (any, error, bool) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok {
		// 已有同在 key 调用在执行，复用
		c.dups++
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err, true
	}

	// 当前 goroutine 成为“leader”
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	// 同步执行用户逻辑
	c.run(g, key, fn)

	return c.val, c.err, c.dups > 0
}

// DoChan 和 Do 类似，但返回一个 channel，调用方可以 select 等待结果。
func (g *Group) DoChan(key string, fn func() (any, error)) <-chan Result {
	ch := make(chan Result, 1)

	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}

	if c, ok := g.m[key]; ok {
		// 已有同 key 调用在执行，复用，只需把自己的 channel 加进去
		c.dups++
		c.chans = append(c.chans, ch)
		g.mu.Unlock()
		return ch
	}

	// 当前 goroutine 成为“leader”，但这里我们用 goroutine 异步执行 fn
	c := new(call)
	c.wg.Add(1)
	c.chans = append(c.chans, ch)
	g.m[key] = c
	g.mu.Unlock()

	go c.run(g, key, fn)
	return ch
}

// DoContext 在 DoChan 的基础上加了 context 支持，
func (g *Group) DoContext(ctx context.Context, key string, fn func() (any, error)) (any, error, bool) {
	// 当 ctx 先结束时，会返回 ctx.Err()，但内部 fn 依然会继续执行并可被其他调用复用。
	ch := g.DoChan(key, fn)

	select {
	case <-ctx.Done():
		// 自己这次调用不再关心结果，返回 ctx.Err
		return nil, ctx.Err(), false
	case res := <-ch:
		return res.Val, res.Err, res.Shared
	}
}

// Forget 让 group 忘记某个 key，
func (g *Group) Forget(key string) {
	// 这样即便该 key 对应的调用还在执行，之后对同一 key 的 Do/DoChan 将会触发新的调用。
	g.mu.Lock()
	if g.m == nil {
		g.mu.Unlock()
		return
	}
	if c, ok := g.m[key]; ok {
		c.forgotten = true
		delete(g.m, key)
	}
	g.mu.Unlock()
}

// run 真正执行 fn 的逻辑，统一处理结束后的收尾工作：
func (c *call) run(g *Group, key string, fn func() (any, error)) {
	// 执行用户逻辑（同步或异步都要保证 runs before cleanup）
	c.val, c.err = fn()

	// cleanup，需要在 fn 完全结束后进行
	g.mu.Lock()
	if !c.forgotten {
		delete(g.m, key)
	}
	chans := c.chans
	res := Result{
		Val:    c.val,
		Err:    c.err,
		Shared: c.dups > 0,
	}
	g.mu.Unlock()

	// 唤醒等待者
	for _, ch := range chans {
		ch <- res
		close(ch)
	}

	// 标记 fn 已完成
	c.wg.Done()
}
