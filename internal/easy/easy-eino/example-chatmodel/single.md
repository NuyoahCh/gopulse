go run single.go

回答:\n在 Go 语言中，**goroutine** 是一种轻量级的并发执行单元，由 Go 运行时（runtime）管理。它类似于传统编程语言中的线程，但比线程更加高效和轻量。

### 主要特点：
1. **轻量级**：
    - 一个 goroutine 的初始栈大小通常只有几 KB，并且可以根据需要动态增长或缩减。
    - 相比之下，一个操作系统线程通常需要几 MB 的栈空间。

2. **由 Go 运行时管理**：
    - goroutine 的调度由 Go 运行时负责，而不是操作系统内核。Go 运行时使用 M:N 调度模型，将多个 goroutine 映射到少量的操作系统线程上执行。

3. **创建简单**：
    - 使用 `go` 关键字即可启动一个 goroutine，例如：
      ```go
      go func() {
          // 并发执行的代码
      }()
      ```

4. **通信通过 Channel**：
    - goroutine 之间通常通过 **channel** 进行通信，这是 Go 语言并发模型的核心思想：“不要通过共享内存来通信，而应该通过通信来共享内存”。

### 示例：
```go
package main

import (
    "fmt"
    "time"
)

func sayHello() {
    fmt.Println("Hello from goroutine!")
}

func main() {
    // 启动一个 goroutine
    go sayHello()

    // 主 goroutine 等待一段时间，确保子 goroutine 有机会执行
    time.Sleep(100 * time.Millisecond)
    fmt.Println("Hello from main!")
}
```

### 优势：
- **高并发**：可以轻松创建成千上万个 goroutine 而不会导致系统资源耗尽。
- **高效调度**：Go 运行时的调度器能够充分利用多核 CPU 的性能。
- **简化并发编程**：通过 goroutine 和 channel 的组合，可以编写出清晰、安全的并发代码。

### 注意事项：
- goroutine 之间没有父子关系，一个 goroutine 的退出不会影响其他 goroutine。
- 需要小心处理 goroutine 之间的同步和数据竞争问题，可以使用 `sync` 包或 channel 来协调。

goroutine 是 Go 语言并发编程的核心，它让编写高并发程序变得简单而高效。\n%                                                                                                                                          