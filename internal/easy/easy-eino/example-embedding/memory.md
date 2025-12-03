go run memory.go


=== 向量化文档 ===
✅ 成功向量化 3 个文档\n\n=== 检索: Go 语言的并发是如何实现的？ ===\n✅ 检索到 2 个相关文档:\n\n文档 1:\n  内容: Eino 是用 Go 开发的 AI 应用框架，提供了丰富的组件和灵活的编排能力。\n  元数据: map[source:eino_intro type:framework]\n\n文档 2:\n  内容: Go 的并发模型基于 CSP（通信顺序进程）。主要通过 goroutine 和 channel 实现。\n  元数据: map[source:go_concurrency type:advanced]\n\n=== AI 回答（基于检索的知识）===
Go 语言基于 CSP（通信顺序进程）理论实现了并发模型，主要通过以下两个核心机制：

1. **Goroutine**
    - 轻量级线程，由 Go 运行时管理
    - 创建成本极低（初始栈约 2KB）
    - 语法简单，使用 `go` 关键字即可启动

2. **Channel**
    - 类型安全的通信管道
    - 提供同步和数据传输功能
    - 支持带缓冲和无缓冲两种模式

示例：
```go
// 创建通道
ch := make(chan int)

// 启动 goroutine
go func() {
    ch <- 42  // 发送数据
}()

value := <-ch  // 接收数据
```

这种设计实现了：
- 通过通信共享内存（而非通过共享内存通信）
- 天然的并发安全
- 优雅的并发程序结构

这正是 Eino 框架能够高效处理 AI 应用并发任务的基础架构。