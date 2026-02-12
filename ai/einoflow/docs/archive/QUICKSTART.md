# EinoFlow 快速开始指南

## 🎯 项目概述

EinoFlow 是一个基于字节跳动 Eino 框架的完整 AI 应用平台，主要使用**字节豆包**作为 LLM 提供商，支持：

- ✅ **基础对话** - 单轮问答
- ✅ **流式对话** - 实时输出
- ✅ **Agent 系统** - 智能任务处理
- ✅ **Chain 编排** - 多步骤处理
- ✅ **RAG 问答** - 知识库检索
- ✅ **Graph 编排** - 复杂流程处理

## 🚀 快速启动（3 步）

### 1. 配置环境变量

确保 `.env` 文件已配置好字节豆包的 API Key：

```env
# 字节豆包配置（主要使用）
ARK_API_KEY=your_ark_api_key_here
ARK_BASE_URL=https://ark.cn-beijing.volces.com/api/v3

# OpenAI 配置（备用）
OPENAI_API_KEY=your_openai_key_here
OPENAI_BASE_URL=https://api.openai.com/v1
```

### 2. 启动服务

```bash
go run cmd/server/main.go
```

看到以下输出表示启动成功：
```
{"level":"info","msg":"Starting EinoFlow server...","time":"..."}
{"level":"info","msg":"Server listening on 0.0.0.0:8080","time":"..."}
```

### 3. 测试 API

```bash
# 测试基础对话
curl -X POST http://localhost:8080/api/v1/llm/chat \
  -H "Content-Type: application/json" \
  -d '{
    "model": "ep-20241116153014-gfmhp",
    "messages": [{"role": "user", "content": "你好"}]
  }'
```

## 📋 完整 API 列表

### 1. LLM 对话

#### 基础对话
```bash
curl -X POST http://localhost:8080/api/v1/llm/chat \
  -H "Content-Type: application/json" \
  -d '{
    "model": "ep-20241116153014-gfmhp",
    "messages": [
      {"role": "user", "content": "解释一下什么是 Eino 框架"}
    ]
  }'
```

#### 流式对话
```bash
curl -N -X POST http://localhost:8080/api/v1/llm/chat/stream \
  -H "Content-Type: application/json" \
  -d '{
    "model": "ep-20241116153014-gfmhp",
    "messages": [
      {"role": "user", "content": "写一首关于编程的诗"}
    ]
  }'
```

#### 获取模型列表
```bash
# 所有提供商的模型
curl http://localhost:8080/api/v1/llm/models

# 指定提供商
curl http://localhost:8080/api/v1/llm/models?provider=ark
```

### 2. Agent 智能任务

```bash
curl -X POST http://localhost:8080/api/v1/agent/run \
  -H "Content-Type: application/json" \
  -d '{
    "task": "分析 Go 语言和 Python 的优缺点，给出学习建议"
  }'
```

**响应示例：**
```json
{
  "answer": "Go 语言优点：1. 并发性能强... Python 优点：1. 语法简洁..."
}
```

### 3. Chain 多步骤处理

```bash
curl -X POST http://localhost:8080/api/v1/chain/run \
  -H "Content-Type: application/json" \
  -d '{
    "steps": [
      "将以下内容翻译成英文",
      "总结成一句话",
      "用专业的语气重写"
    ],
    "input": "Go 是一门很棒的编程语言，它简洁、高效、并发性能强"
  }'
```

**响应示例：**
```json
{
  "result": "Go is a professionally acclaimed programming language...",
  "steps": 3
}
```

### 4. RAG 知识库问答

#### 索引文档
```bash
curl -X POST http://localhost:8080/api/v1/rag/index \
  -H "Content-Type: application/json" \
  -d '{
    "documents": [
      "Eino 是字节跳动开源的 LLM 应用开发框架",
      "Eino 支持 Chain、Agent、RAG、Graph 等功能",
      "Eino 使用 Go 语言编写，性能优秀"
    ]
  }'
```

#### 查询
```bash
curl -X POST http://localhost:8080/api/v1/rag/query \
  -H "Content-Type: application/json" \
  -d '{
    "query": "Eino 有哪些主要功能？"
  }'
```

### 5. Graph 复杂编排

```bash
curl -X POST http://localhost:8080/api/v1/graph/run \
  -H "Content-Type: application/json" \
  -d '{
    "query": "如何成为一名优秀的 Go 语言开发者？",
    "type": "multi_step"
  }'
```

**处理流程：**
1. 分析问题 → 2. 制定计划 → 3. 执行总结

## 🎨 使用场景示例

### 场景 1：智能客服

```bash
# 使用 Agent 处理客户问题
curl -X POST http://localhost:8080/api/v1/agent/run \
  -H "Content-Type: application/json" \
  -d '{
    "task": "用户询问：你们的产品支持哪些功能？请详细介绍"
  }'
```

### 场景 2：内容创作

```bash
# 使用 Chain 进行内容创作流程
curl -X POST http://localhost:8080/api/v1/chain/run \
  -H "Content-Type: application/json" \
  -d '{
    "steps": [
      "根据主题生成大纲",
      "扩展每个要点",
      "润色文字，使其更专业"
    ],
    "input": "主题：人工智能的未来发展"
  }'
```

### 场景 3：知识问答

```bash
# 先索引知识库
curl -X POST http://localhost:8080/api/v1/rag/index \
  -H "Content-Type: application/json" \
  -d '{
    "documents": [
      "公司成立于2020年，专注于AI技术",
      "主要产品包括智能对话、文本生成、图像识别",
      "服务客户超过1000家企业"
    ]
  }'

# 然后查询
curl -X POST http://localhost:8080/api/v1/rag/query \
  -H "Content-Type: application/json" \
  -d '{
    "query": "公司什么时候成立的？有哪些产品？"
  }'
```

### 场景 4：复杂分析

```bash
# 使用 Graph 进行多步骤分析
curl -X POST http://localhost:8080/api/v1/graph/run \
  -H "Content-Type: application/json" \
  -d '{
    "query": "分析当前 AI 行业的发展趋势，并给出投资建议",
    "type": "multi_step"
  }'
```

## 🔧 可用的豆包模型

| 模型 ID | 名称 | 上下文长度 | 适用场景 |
|---------|------|-----------|----------|
| `ep-20241116152913-xdvqz` | 豆包-lite-4k | 4K | 简单对话 |
| `ep-20241116153014-gfmhp` | 豆包-pro-4k | 4K | **推荐使用** |
| `ep-20241116153056-8nqkl` | 豆包-turbo-4k | 4K | 快速响应 |
| `ep-20241116153137-jzlgr` | 豆包-lite-32k | 32K | 长文本 |
| `ep-20241116153211-lnmwz` | 豆包-pro-32k | 32K | 复杂任务 |

## 📊 性能参考

| API | 平均响应时间 | 说明 |
|-----|-------------|------|
| `/llm/chat` | 3-8 秒 | 单次调用 |
| `/llm/chat/stream` | 5-15 秒 | 流式输出 |
| `/agent/run` | 5-10 秒 | 智能处理 |
| `/chain/run` (3步) | 15-25 秒 | 3次模型调用 |
| `/graph/run` | 20-40 秒 | 多步骤分析 |

## 🐛 常见问题

### Q1: 服务启动失败？
**A:** 检查 `.env` 文件是否配置了 `ARK_API_KEY`

### Q2: API 返回 500 错误？
**A:** 查看服务器日志，通常是 API Key 无效或网络问题

### Q3: 流式响应不工作？
**A:** 使用 `curl -N` 参数，或在浏览器中使用 EventSource

### Q4: 响应太慢？
**A:** 
- 使用 `豆包-turbo-4k` 模型
- 减少 Chain/Graph 的步骤数
- 检查网络连接

## 📚 更多文档

- `README.md` - 项目概览
- `docs/DEMO_GUIDE.md` - 演示指南
- `docs/FINAL_STATUS.md` - 完成状态
- `docs/COMPLETE_IMPLEMENTATION.md` - 实现细节

## 🎉 开始使用

现在你已经了解了所有功能，开始构建你的 AI 应用吧！

```bash
# 启动服务
go run cmd/server/main.go

# 在另一个终端测试
curl http://localhost:8080/health
```

祝你使用愉快！🚀
