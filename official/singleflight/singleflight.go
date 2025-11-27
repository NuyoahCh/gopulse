// Copyright...
// Package singleflight provides a duplicate function call suppression mechanism.
// singleflight 提供了一种机制：对同一 key 的函数调用，只允许一次执行，
// 其他重复调用将等待第一次执行完成并共享结果。

package singleflight

// import "golang.org/x/sync/singleflight"

import (
	"bytes"
	"errors"
	"fmt"
	"runtime"
	"runtime/debug"
	"sync"
)

// errGoexit 表示用户提供的函数调用了 runtime.Goexit。
// Goexit 结束当前 goroutine，但不能被 recover 捕获。
var errGoexit = errors.New("runtime.Goexit was called")

// panicError 用于保存 panic 的值与堆栈信息
type panicError struct {
	value interface{} // panic 的原始值
	stack []byte      // 堆栈信息
}

// 实现 error 接口，打印 panic 信息和堆栈
func (p *panicError) Error() string {
	return fmt.Sprintf("%v\n\n%s", p.value, p.stack)
}

// newPanicError 构造 panicError，并截掉 goroutine 标题行（因为可能已失效）
func newPanicError(v interface{}) error {
	stack := debug.Stack()

	// 堆栈第一行通常是 "goroutine N [xxx]:"，此时 goroutine 可能已结束，这行会误导，
	// 因此将其删除。
	if line := bytes.IndexByte(stack[:], '\n'); line >= 0 {
		stack = stack[line+1:]
	}
	return &panicError{value: v, stack: stack}
}

// call 表示一次唯一的 Do 调用（正在执行或已完成）
type call struct {
	wg sync.WaitGroup // 用于等待 call 执行结束

	// 这些字段只会在调用结束前写入，在 WaitGroup 完成之后读取
	val interface{} // fn 返回值
	err error       // fn 返回错误

	// forgotten 表示在 call 进行过程中是否调用过 Forget(key)
	// 若为 true，则 call 完成后不会自动从 map 中删除
	forgotten bool

	// 以下字段在调用过程中会被 mutex 保护写入，完成后只读
	dups  int             // 同一 key 的重复请求数量
	chans []chan<- Result // 注册等待结果的 channel 列表（用于 DoChan）
}

// Group 管理多个 key 对应的执行任务，相当于一个命名空间
type Group struct {
	mu sync.Mutex       // 保护 m
	m  map[string]*call // key -> call 的映射表
}

// Result 保存 Do() 和 DoChan() 执行结果
type Result struct {
	Val    interface{} // 返回值
	Err    error       // 错误
	Shared bool        // 是否共享（是否为重复调用返回）
}

// Do 执行给定 key 的函数 fn。若同 key 正在执行，则重复请求等待并共享结果。
func (g *Group) Do(key string, fn func() (interface{}, error)) (v interface{}, err error, shared bool) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}

	// 若 key 已经存在说明正在执行：增加 dups 后等待
	if c, ok := g.m[key]; ok {
		c.dups++
		// 解锁，不阻塞其他 key
		g.mu.Unlock()

		// 等待原执行完成
		c.wg.Wait()

		// 若原调用抛 panic，则同样抛 panic
		if e, ok := c.err.(*panicError); ok {
			panic(e)
		} else if c.err == errGoexit {
			// 如果原调用调用了 runtime.Goexit，则也 Goexit
			runtime.Goexit()
		}
		return c.val, c.err, true
	}

	// key 不存在：创建新的 call 作为首次执行
	c := new(call)
	c.wg.Add(1)   // 标记 call 进行中
	g.m[key] = c  // 注册到 map 中
	g.mu.Unlock() // 解锁

	// 执行 fn
	g.doCall(c, key, fn)

	return c.val, c.err, c.dups > 0
}

// DoChan 与 Do 类似，但是返回一个异步 channel，执行结果写入 channel
func (g *Group) DoChan(key string, fn func() (interface{}, error)) <-chan Result {
	ch := make(chan Result, 1) // 结果只会写一次，缓冲避免死锁

	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}

	// 若 key 已存在，说明正在执行：注册当前 channel 等候结果
	if c, ok := g.m[key]; ok {
		c.dups++
		c.chans = append(c.chans, ch)
		g.mu.Unlock()
		return ch
	}

	// 新建 call
	c := &call{chans: []chan<- Result{ch}}
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	// 异步执行
	go g.doCall(c, key, fn)

	return ch
}

// doCall 真正执行 fn 的逻辑，并处理 panic、Goexit、正常返回等情况
func (g *Group) doCall(c *call, key string, fn func() (interface{}, error)) {
	normalReturn := false // 标记是否正常返回
	recovered := false    // 标记是否从 panic 中恢复

	// 双层 defer 用于区分 panic 和 Goexit
	defer func() {
		// 若既不是正常返回也不是 panic，则说明是 Goexit
		if !normalReturn && !recovered {
			c.err = errGoexit
		}

		c.wg.Done() // 任务结束，唤醒等待者

		g.mu.Lock()
		defer g.mu.Unlock()

		// 若未调用 Forget(key)，则从 map 中删除该 key
		if !c.forgotten {
			delete(g.m, key)
		}

		// 如果是 panicError，需要确保所有等待者不会永远阻塞
		if e, ok := c.err.(*panicError); ok {
			if len(c.chans) > 0 {
				// 为了让所有等待者醒来，让 panic 发生在一个新 goroutine 中
				go panic(e)
				select {} // 保持 goroutine 存活，便于 crash dump 查看
			} else {
				panic(e)
			}
		} else if c.err == errGoexit {
			// Goexit 已处理，无需再次调用
		} else {
			// 正常返回，将结果写入所有等待 channel
			for _, ch := range c.chans {
				ch <- Result{c.val, c.err, c.dups > 0}
			}
		}
	}()

	// 内层 defer，用于捕获 panic，并转为 panicError
	func() {
		defer func() {
			// 若不是正常返回，则可能发生了 panic
			if !normalReturn {
				if r := recover(); r != nil {
					c.err = newPanicError(r)
				}
			}
		}()

		// 执行用户函数
		c.val, c.err = fn()
		normalReturn = true
	}()

	// 若不是正常返回但 recover 成功，则标记 recovered
	if !normalReturn {
		recovered = true
	}
}

// Forget 强制从 map 中删除 key， 未来对于该 key 的调用不会等待之前的结果。
func (g *Group) Forget(key string) {
	g.mu.Lock()
	if c, ok := g.m[key]; ok {
		c.forgotten = true
	}
	delete(g.m, key)
	g.mu.Unlock()
}
