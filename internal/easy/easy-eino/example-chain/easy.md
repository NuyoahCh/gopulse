go run easy.go

回答: Goroutine 是 Go 语言中并发编程的核心概念，它是一种轻量级的线程，由 Go 运行时（runtime）管理。以下是关于 goroutine 的主要特点：

1. **轻量级**
    - 初始栈大小仅 2KB（可动态扩容/缩容）
    - 创建和销毁开销极小，可轻松创建数十万个

2. **创建方式**
   ```go
   go funcName()  // 在函数调用前添加 go 关键字
   ```

3. **调度机制**
    - 采用 M:N 调度模型（M 个 goroutine 映射到 N 个 OS 线程）
    - 由 Go 运行时自主调度，非操作系统调度
    - 基于工作窃取的调度算法

4. **通信机制**
    - 通过 channel 进行 goroutine 间通信
    - 遵循 CSP（Communicating Sequential Processes）模型

5. **优势特点**
    - 创建成本远低于线程（线程通常 MB 级栈空间）
    - 快速启动时间（无需操作系统介入）
    - 智能调度可自动利用多核处理器

示例：
```go
func main() {
    go func() {
        fmt.Println("Hello from goroutine")
    }()
    time.Sleep(time.Millisecond) // 等待goroutine执行
}
```


注意：实际开发中应使用 sync.WaitGroup 或 channel 来协调 goroutine 执行，而非 time.Sleep。\n%                       