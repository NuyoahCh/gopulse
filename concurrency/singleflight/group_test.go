package singleflight

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// 基础测试：单次调用
func TestDoSingle(t *testing.T) {
	var g Group

	v, err, shared := g.Do("k", func() (any, error) {
		return "hello", nil
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.(string) != "hello" {
		t.Fatalf("unexpected value: %v", v)
	}
	if shared {
		t.Fatalf("first call should not be shared")
	}
}

// 并发测试：多个 goroutine 对同一 key 并发调用，只应执行一次 fn
// TODO 并发测试还是存在问题，待实现
func TestDoDuplicateSuppression(t *testing.T) {
	var g Group
	var counter int32

	started := make(chan struct{})
	block := make(chan struct{})

	// 第一个 goroutine，负责真正执行 fn
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, _, shared := g.Do("k", func() (any, error) {
			atomic.AddInt32(&counter, 1)
			// 告诉外面：我已经开始执行 fn 了
			close(started)
			// 阻塞一段时间，让其他 goroutine 有机会进来复用
			<-block
			return "ok", nil
		})
		if shared {
			t.Errorf("first caller should not be shared")
		}
	}()

	// 确保第一个 fn 已经开始执行
	<-started

	const N = 10
	var wg2 sync.WaitGroup
	wg2.Add(N)

	// 其它 goroutine 并发调用
	for i := 0; i < N; i++ {
		go func(i int) {
			defer wg2.Done()
			v, err, shared := g.Do("k", func() (any, error) {
				t.Errorf("duplicate fn should not be called")
				return nil, nil
			})
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if v.(string) != "ok" {
				t.Errorf("unexpected value: %v", v)
			}
			if i > 0 && !shared {
				t.Errorf("duplicate caller should see shared=true")
			}
		}(i)
	}

	// 释放阻塞，让真正的 fn 结束
	close(block)

	wg2.Wait()
	wg.Wait()

	if got := atomic.LoadInt32(&counter); got != 1 {
		t.Fatalf("fn should be called once, got %d", got)
	}
}

// 测试不同 key 互不影响
func TestDoDifferentKeys(t *testing.T) {
	var g Group
	var c1, c2 int32

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		v, _, _ := g.Do("k1", func() (any, error) {
			atomic.AddInt32(&c1, 1)
			time.Sleep(10 * time.Millisecond)
			return "v1", nil
		})
		if v.(string) != "v1" {
			t.Errorf("unexpected value for k1: %v", v)
		}
	}()

	go func() {
		defer wg.Done()
		v, _, _ := g.Do("k2", func() (any, error) {
			atomic.AddInt32(&c2, 1)
			time.Sleep(10 * time.Millisecond)
			return "v2", nil
		})
		if v.(string) != "v2" {
			t.Errorf("unexpected value for k2: %v", v)
		}
	}()

	wg.Wait()

	if atomic.LoadInt32(&c1) != 1 || atomic.LoadInt32(&c2) != 1 {
		t.Fatalf("each key should be called once, got c1=%d c2=%d",
			c1, c2)
	}
}

// 测试 DoChan 的基本行为
func TestDoChan(t *testing.T) {
	var g Group
	var counter int32

	resCh := g.DoChan("k", func() (any, error) {
		atomic.AddInt32(&counter, 1)
		time.Sleep(20 * time.Millisecond)
		return "x", nil
	})

	res := <-resCh
	if res.Err != nil {
		t.Fatalf("unexpected error: %v", res.Err)
	}
	if res.Val.(string) != "x" {
		t.Fatalf("unexpected value: %v", res.Val)
	}
	if res.Shared {
		t.Fatalf("first caller via DoChan should not be shared")
	}
	if atomic.LoadInt32(&counter) != 1 {
		t.Fatalf("fn should be called once, got %d", counter)
	}
}

// 测试 DoContext：ctx 先取消
func TestDoContextCanceled(t *testing.T) {
	var g Group
	var counter int32

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	started := make(chan struct{})

	go func() {
		// 稍后一点再取消，确保 fn 已经在执行
		<-started
		time.Sleep(20 * time.Millisecond)
		cancel()
	}()

	// DoContext 调用
	v, err, shared := g.DoContext(ctx, "k", func() (any, error) {
		atomic.AddInt32(&counter, 1)
		close(started)
		time.Sleep(100 * time.Millisecond) // 模拟很慢的操作
		return "done", nil
	})

	if err == nil {
		t.Fatalf("expected context error, got nil")
	}
	if v != nil {
		t.Fatalf("expected nil value when canceled, got %v", v)
	}
	if shared {
		t.Fatalf("when ctx canceled, shared should be false for this caller")
	}

	// 再发起一个正常的 Do，应该可以复用之前那次还在执行的 fn
	v2, err2, shared2 := g.Do("k", func() (any, error) {
		t.Fatalf("second fn should not be called; should reuse existing call")
		return nil, nil
	})

	if err2 != nil {
		t.Fatalf("unexpected error: %v", err2)
	}
	if v2.(string) != "done" {
		t.Fatalf("unexpected value: %v", v2)
	}
	if !shared2 {
		t.Fatalf("second caller should see shared=true")
	}

	if atomic.LoadInt32(&counter) != 1 {
		t.Fatalf("fn should be called once, got %d", counter)
	}
}

// 测试 Forget：在调用进行中 Forget，然后再次调用同 key，会触发第二次执行
func TestForget(t *testing.T) {
	var g Group
	var counter int32

	var wg sync.WaitGroup
	wg.Add(2)

	fn := func() (any, error) {
		atomic.AddInt32(&counter, 1)
		time.Sleep(50 * time.Millisecond)
		return "ok", nil
	}

	// 第一次调用
	go func() {
		defer wg.Done()
		v, err, shared := g.Do("k", fn)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if v.(string) != "ok" {
			t.Errorf("unexpected value: %v", v)
		}
		if shared {
			t.Errorf("first call should not be shared")
		}
	}()

	// 确保第一次已经在执行
	time.Sleep(10 * time.Millisecond)
	g.Forget("k")

	// 第二次调用，同 key，但因为 Forget 过，应该再执行一次 fn
	go func() {
		defer wg.Done()
		v, err, shared := g.Do("k", fn)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if v.(string) != "ok" {
			t.Errorf("unexpected value: %v", v)
		}
		// shared 这里可能 true/false，取决于两个 fn 的先后完成顺序，这里不强行断言
		_ = shared
	}()

	wg.Wait()

	if got := atomic.LoadInt32(&counter); got != 2 {
		t.Fatalf("fn should be called twice due to Forget, got %d", got)
	}
}
