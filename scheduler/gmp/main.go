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

	mu    sync.Mutex // 实现互斥锁
	sched *Scheduler // 反向指向调度器，用于从阻塞中唤醒后重新入队
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
				g.sched.Submit(g)
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

// ---- P：逻辑处理器，持有本地 run queue ----

// P 队列
type P struct {
	id   int
	runq chan *G // 本地队列
}

// ---- Scheduler：调度器，管理 G/P/M ----

// Scheduler 调度器
type Scheduler struct {
	globalRunq chan *G // 全局队列
	ps         []*P    // 所有 P 的集合
	timeSlice  int     // 每次调度给 G 的“时间片”（最大步数）
}

func NewScheduler(numP int, timeSlice int) *Scheduler {
	s := &Scheduler{
		globalRunq: make(chan *G, 1024),
		ps:         make([]*P, numP),
		timeSlice:  timeSlice,
	}
	for i := 0; i < numP; i++ {
		s.ps[i] = &P{
			id:   i + 1,
			runq: make(chan *G, 256),
		}
	}
	return s
}

// Submit 把一个 runnable 的 G 加入调度器, 优先放某个 P 的本地队列，满了再放全局队列
func (s *Scheduler) Submit(g *G) {
	// 可能从 timer goroutine 里调用，所以不用锁，用 channel 就好
	idx := rand.Intn(len(s.ps))
	p := s.ps[idx]

	// 判断本地队列/全局队列
	select {
	case p.runq <- g:
		fmt.Printf("[Scheduler] 提交 %s 到 P%d 本地队列\n", g.String(), p.id)
	default:
		select {
		case s.globalRunq <- g:
			fmt.Printf("[Scheduler] 本地队列满了，提交 %s 到全局队列\n", g.String())
		default:
			// 极端情况下全局也满了，这里我们简单打印一下（真实 runtime 会有更复杂处理）
			fmt.Printf("[Scheduler] 警告：队列已满，丢弃 %s（示例代码，勿用于生产）\n", g.String())
		}
	}
}

// 把 G 优先重新丢回绑定的 P，本地队列满则丢全局
func (s *Scheduler) enqueueToP(p *P, g *G) {
	select {
	case p.runq <- g:
	default:
		select {
		case s.globalRunq <- g:
		default:
			fmt.Printf("[Scheduler] 警告：队列已满，无法重新入队 %s\n", g.String())
		}
	}
}

// 从当前 P、本全局队列、其他 P 依次尝试拿一个 G
func (s *Scheduler) getNextG(p *P) *G {
	// 1. 当前 P 的本地队列
	select {
	case g := <-p.runq:
		return g
	default:
	}

	// 2. 全局队列
	select {
	case g := <-s.globalRunq:
		return g
	default:
	}

	// 其他的 p 中进行获取
	for _, other := range s.ps {
		if other == p {
			continue
		}
		select {
		case g := <-other.runq:
			fmt.Printf("[Steal] P%d 从 P%d 偷到了 %s\n", p.id, other.id, g.String())
			return g
		default:
		}
	}
	return nil
}

// 在某个 M 上、绑定某个 P，执行一个 G，最多执行 timeSlice 步
// 如果时间片用完还没结束，就把 G 重新入队
func (s *Scheduler) runG(mID int, p *P, g *G, gWg *sync.WaitGroup) {
	fmt.Printf("M%d 在 P%d 上开始调度 %s\n", mID, p.id, g.String())

	for i := 0; i < s.timeSlice; i++ {
		res := g.doOneStep(mID, p.id)

		switch res {
		case StepRunning:
			// 继续执行本 G，直到时间片结束
		case StepBlocked:
			fmt.Printf("M%d：%s 在 P%d 上阻塞，切换到其他 G\n", mID, g.String(), p.id)
			return
		case StepDone:
			fmt.Printf("M%d：%s 在 P%d 上全部完成\n", mID, g.String(), p.id)
			gWg.Done()
			return
		}
	}

	// 时间片用完但 G 还没结束：模拟被抢占，重新入队
	fmt.Printf("M%d：%s 在 P%d 上时间片用完，重新入队\n", mID, g.String(), p.id)
	s.enqueueToP(p, g)
}

// worker：模拟一个 M，不断获取 G 执行
func (s *Scheduler) worker(mID int, p *P, gWg *sync.WaitGroup, stop <-chan struct{}) {
	for {
		select {
		case <-stop:
			fmt.Printf("M%d 退出（绑定 P%d）\n", mID, p.id)
			return
		default:
			g := s.getNextG(p)
			if g == nil {
				// 没有任务，稍微休息一下
				time.Sleep(10 * time.Millisecond)
				continue
			}
			s.runG(mID, p, g, gWg)
		}
	}
}
