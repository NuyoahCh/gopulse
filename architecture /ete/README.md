代码包含两种模式：

* linkReliable=true：链路层做“每包校验+重传”，但磁盘仍可能腐化 → 只有端到端校验才能兜底

* linkReliable=false：链路可能翻转比特 → 端到端校验与重试把整体传输拉回正确

```go
go run main.go \
  -sizeMB=4 \
  -packet=1024 \
  -flipProb=0.02 \
  -diskProb=0.05 \
  -linkReliable=false \
  -retries=5
```

将 -linkReliable 切换为 false，感受“无链路可靠性”时端到端校验与重试如何仍然保障整体正确性。

调大 -diskProb，你会看到：即使链路层很可靠，写盘后也可能腐化，只有端到端校验+重试才能保证最终一致。

观察日志中：

[LINK]：链路层的“每包校验+重传”（性能增强）

[DISK]：存储层模拟腐化

[APP]：端到端哈希校验与重试（正确性保证）

