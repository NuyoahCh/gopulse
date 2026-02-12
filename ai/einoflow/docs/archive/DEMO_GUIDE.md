# Demo 演示指南

## 运行演示程序

```bash
make demo
```

## 功能说明

### 1. 基础对话 ✅ 完全可用

**功能**: 与 AI 进行简单的问答对话

**示例**:
```
请输入问题: hello
回答: Hello! How can I assist you today? 😊
```

**特点**:
- 单轮对话
- 快速响应
- 适合简单问答

---

### 2. 流式对话 ✅ 完全可用

**功能**: 实时流式输出 AI 的回答

**示例**:
```
请输入问题: 今天天气怎么样
回答: 
天气会因具体城市/地区和实时时间而变化...
（内容会逐字显示）
```

**特点**:
- 实时输出
- 更好的用户体验
- 适合长文本生成

---

### 3. Agent 工具调用 ⚠️ 演示版本

**功能**: 展示 Agent 如何理解和描述工具使用

**当前状态**: 
- ✅ 模型会识别需要使用的工具
- ✅ 模型会描述如何使用工具
- ❌ 不会真正执行工具调用（简化实现）

**示例输出**:
```
请输入任务: 给我找到今天北京的天气，进行记录

Agent 执行中...
💡 提示: 当前是简化版 Agent，模型会描述如何使用工具，但不会真正执行
   如需完整的工具调用功能，请使用 Eino 的 react.NewAgent API

Agent 响应:
<|FunctionCallBegin|>[{"name":"tools.WeatherTool","parameters":{"city":"北京","date":"今天"}}]<|FunctionCallEnd|>
```

**说明**:
- 模型返回的 `<|FunctionCallBegin|>` 是豆包模型的工具调用格式
- 这表明模型正确识别了需要调用天气工具
- 要实现真正的工具执行，需要：
  1. 解析这个 JSON 格式
  2. 调用对应的工具函数
  3. 将结果返回给模型
  4. 让模型生成最终答案

**完整实现参考**:
```go
// 使用 Eino 的完整 ReAct Agent
import "github.com/cloudwego/eino/flow/agent/react"

agent, err := react.NewAgent(ctx, &react.AgentConfig{
    ToolCallingModel: toolCallingModel,
    ToolsConfig: compose.ToolsNodeConfig{
        Tools: toolsList,
    },
})
```

---

### 4. Graph 多步骤处理 ✅ 完全可用

**功能**: 将复杂问题分解为多个步骤处理

**处理流程**:
1. **分析问题** - 理解用户意图和关键点
2. **制定计划** - 设计解决方案的步骤
3. **执行总结** - 给出最终答案

**示例**:
```
请输入复杂问题: Go和Java哪个好

执行多步骤分析...
步骤 1/3: 分析问题...
✅ 步骤 1/3: 问题分析完成
步骤 2/3: 制定计划...
✅ 步骤 2/3: 计划制定完成
步骤 3/3: 执行并总结...
✅ 步骤 3/3: 执行完成

==================================================

最终答案:
（详细的对比分析）

==================================================
```

**特点**:
- 结构化思考
- 多步骤推理
- 适合复杂问题分析

**注意**: 
- 每个步骤都会调用模型，所以总耗时较长（通常 20-40 秒）
- 如果某个步骤超时，可能是网络问题或模型响应慢

---

## 常见问题

### Q1: Agent 为什么不执行工具？

**A**: 当前是简化版实现，主要用于演示。完整的工具调用需要：
1. 使用 Eino 的 `react.NewAgent` API
2. 配置 `ToolCallingChatModel` 而不是普通的 `ChatModel`
3. 实现工具调用的解析和执行逻辑

### Q2: Graph 执行太慢怎么办？

**A**: Graph 需要调用模型 3 次（每个步骤一次），建议：
1. 使用更快的模型（如 lite 版本）
2. 简化问题描述
3. 检查网络连接

### Q3: 如何实现真正的工具调用？

**A**: 参考 Eino 官方文档：
```go
// 1. 使用 ToolCallingChatModel
toolCallingModel, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
    APIKey: apiKey,
    Model:  "doubao-seed-1-6-lite-251015",
})

// 2. 创建 ReAct Agent
agent, err := react.NewAgent(ctx, &react.AgentConfig{
    ToolCallingModel: toolCallingModel,
    ToolsConfig: compose.ToolsNodeConfig{
        Tools: []tool.InvokableTool{
            weatherTool,
            calculatorTool,
        },
    },
})

// 3. 执行
result, err := agent.Invoke(ctx, &react.AgentInput{
    Query: "查询北京天气",
})
```

---

## 性能参考

| 功能 | 平均响应时间 | 模型调用次数 |
|------|-------------|-------------|
| 基础对话 | 3-8 秒 | 1 次 |
| 流式对话 | 5-15 秒 | 1 次（流式）|
| Agent | 5-10 秒 | 1 次 |
| Graph | 20-40 秒 | 3 次 |

---

## 下一步

### 使用 Web API
```bash
# 启动服务器
make run

# 测试基础对话
curl -X POST http://localhost:8080/api/v1/llm/chat \
  -H "Content-Type: application/json" \
  -d '{
    "model": "doubao-seed-1-6-lite-251015",
    "messages": [{"role": "user", "content": "你好"}]
  }'
```

### 查看更多示例
- `examples/llm/basic_chat.go` - LLM 基础使用
- `examples/agent/weather_agent.go` - Agent 示例
- `examples/rag/simple_rag.go` - RAG 示例

### 阅读文档
- `docs/TROUBLESHOOTING.md` - 故障排查
- `QUICKFIX.md` - 快速修复指南
- `README.md` - 项目概览
