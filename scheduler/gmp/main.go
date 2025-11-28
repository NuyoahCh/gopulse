package gmp

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// ---- G 状态 & 步骤执行结果 ----

// GState Goroutine 调度状态
type GState int

// Goroutine 相关参数
const (
	GIdle GState = iota
	GRunnable
	GRunning
	GWaiting
	GDead
)

// StepResult 步骤结果
type StepResult int

// 步骤执行流程
const (
	StepRunning StepResult = iota // 还有后续步骤
	StepBlocked                   // 本步进入阻塞（模拟系统调用）
	StepDone                      // 所有步骤完成
)

// ---- G：被调度的“goroutine” ----

// G 协程
type G struct {
	id            int
	name          string
	totalSteps    int           // 总共要执行多少步
	doneSteps     int           // 已执行的步数
	blockAtStep   int           // 在第几步模拟阻塞（0 表示不阻塞）
	blockDuration time.Duration // 阻塞多久
	state         GState        // 表示状态

	mu sync.Mutex // 实现互斥锁

	// todo 实现调度器
}

// String helper func
func (g *G) String() string {
	return fmt.Sprintf("G%d(%s)", g.id, g.name)
}

// 执行 G 的“一步”逻辑，返回这一步之后的状态
func (g *G) doOneStep(mID, pID int) StepResult {
	g.mu.Lock()
	defer g.mu.Unlock()

	// Dead 就返回 Done 结束
	if g.state == GDead {
		return StepDone
	}
	// Wait 就返回 Blocked 阻塞
	if g.state == GWaiting {
		return StepBlocked
	}

	g.state = GRunning
	// 执行的步数
	g.doneSteps++

	step := g.doneSteps
	fmt.Printf("    %s 在 M%d/P%d 上执行第 %d/%d 步\n",
		g.String(), mID, pID, step, g.totalSteps)

	// 模拟这一步的计算耗时
	time.Sleep(time.Duration(rand.Intn(50)+30) * time.Millisecond)

	// 检查是否需要在这一“步”模拟阻塞系统调用（比如 read/write/net 调用）
	if g.blockAtStep > 0 && step == g.blockAtStep && g.blockDuration > 0 {
		fmt.Printf("    %s 在第 %d 步执行阻塞系统调用，阻塞 %v\n",
			g.String(), step, g.blockDuration)

		g.state = GWaiting
		// 保存阻塞时间并清零，避免多次阻塞
		d := g.blockDuration
		g.blockDuration = 0

		// 用一个额外的 goroutine + timer 来模拟 netpoll/timer 唤醒这个 G
		go func(g *G, d time.Duration) {
			time.Sleep(d)
			g.mu.Lock()
			defer g.mu.Unlock()

			if g.state == GWaiting {
				g.state = GRunnable
				fmt.Printf("    [Timer] %s 阻塞结束，重新变为 runnable，重新入队\n", g.String())
				// todo 实现调度器
			}
		}(g, d)

		return StepBlocked
	}

	if g.doneSteps >= g.totalSteps {
		g.state = GDead
		return StepDone
	}

	// 这一步完成，但 G 还有后续步骤
	g.state = GRunnable
	return StepRunning
}
