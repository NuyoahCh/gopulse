go run easy.go


\n\n=== 最终结果 ===\nGoroutine 是 Go 语言中的轻量级并发执行单元，由 Go 运行时（runtime）管理，可以看作是 Go 语言实现并发编程的核心机制。以下是关于 goroutine 的关键特性：

1. **轻量级**
    - 初始栈大小仅 2KB（可动态扩容/缩容）
    - 创建和销毁开销远小于操作系统线程
    - 单进程可轻松创建数十万个 goroutine

2. **调度模型**
    - 采用 M:P:G 调度模型（MPG）
    - 由 Go 运行时调度，非操作系统直接调度
    - M（Machine）对应操作系统线程
    - P（Processor）维护运行队列的上下文
    - G（Goroutine）即并发执行单元

3. **创建方式**
   ```go
   // 使用 go 关键字快速创建
   go func() {
       // 并发执行的任务
   }()
   ```

4. **通信机制**
    - 通过 channel 进行 goroutine 间通信
    - 遵循 CSP（Communicating Sequential Processes）模型
    - 支持数据传递和同步控制

5. **实践特性**
    - 自动在多个系统线程间复用
    - 支持抢占式调度（since Go 1.14）
    - 当 goroutine 阻塞时，调度器会自动切换其他 goroutine 执行

示例场景：
```go
func main() {
    go worker("A") // 启动并发任务
    go worker("B")
    time.Sleep(time.Second)
}

func worker(name string) {
    for i := 0; i < 3; i++ {
        fmt.Printf("%s: %d\n", name, i)
    }
}
```

Goroutine 的设计使 Go 程序能够以极低资源消耗实现高并发，成为构建现代分布式系统和网络服务的理想选择。\n%                                                      