# EinoFlow - 基于 Eino 的 AI 应用平台

EinoFlow 构建于字节跳动开源框架 [Eino](https://github.com/cloudwego/eino) 之上，开箱支持对话、Agent、RAG、Chain、Graph 等核心 AI 能力，并提供后端 API 与 React/Vite 前端示例，适合快速验证原型或作为二次开发的基座。

## ✨ 功能特性
- **多模型 LLM**：内置豆包（默认）、OpenAI 等提供商，支持文本/多模态与流式输出。
- **RAG 能力**：文档加载、分块、向量化与检索，**支持文件上传**（TXT/MD），默认内存向量库，便于扩展。
- **Chain & Graph 编排**：顺序链、并行链、多步骤 Graph，支持复杂流程拆解。
- **Agent 系统**：ReAct Agent + Function Calling，可插入自定义工具，**支持实时天气查询**，多轮对话。
- **上下文与可观测性**：内置对话记忆、上下文截断、结构化日志与调用追踪。
- **前端演示**：React + Vite UI，涵盖 Chat/Agent/RAG/Graph 页面，**完整 Markdown 渲染**，便于集成。

### 🆕 最新更新（v1.1.0）
- 🌤️ **天气查询工具**：Agent 可查询实时天气信息，支持中国主要城市
- 📁 **RAG 文件上传**：支持直接上传 TXT/MD 文件进行索引，最大 10MB
- 🎨 **Markdown 渲染**：前端完整支持 Markdown 格式，包括代码高亮、表格、列表等
- 🔒 **安全配置优化**：移除硬编码密钥，所有配置通过环境变量管理

详细说明请查看 [新功能文档](docs/NEW_FEATURES.md) 和 [测试指南](docs/TESTING_GUIDE.md)。

## 🧭 仓库结构
EinoFlow 是一个构建在字节跳动 [Eino](https://github.com/cloudwego/eino) 框架之上的完整 AI 应用平台，开箱即可体验 LLM 对话、Agent、RAG、Chain、Graph 等能力，并提供前后端示例、可观测性与配置管理，适合快速验证 AI 产品原型或作为二次开发的基础。

## ✨ 功能特性
- **多模型 LLM**：内置字节豆包（默认）、OpenAI 等提供商，支持文本/多模态、流式与非流式输出。
- **RAG 体系**：文档加载、分块与向量化，内置检索与基础问答能力，便于扩展向量数据库。
- **Chain & Graph 编排**：顺序链、并行链、分支链及多步骤 Graph，支持复杂业务流程组装。
- **Agent 系统**：ReAct Agent + Function Calling，可接入自定义工具，支持多轮对话。
- **Memory 与可观测性**：对话记忆、上下文窗口控制，结构化日志、请求追踪、性能与成本统计。
- **前后端示例**：后端 RESTful API，前端基于 React + Vite 的演示界面，便于二次集成。

## 🧭 项目结构
```
.
├── cmd/server/main.go         # 服务入口
├── internal/                  # 核心业务逻辑
│   ├── api/                   # 路由与 Handler
│   ├── llm/                   # LLM 抽象与 Provider
│   ├── agent/                 # ReAct Agent
│   ├── chain/                 # Chain 编排
│   ├── graph/                 # Graph 任务编排
│   ├── rag/                   # 文档加载、分块、检索
│   ├── memory/                # 对话上下文管理
│   └── config/                # 配置加载
├── web/                       # React/Vite 前端示例
├── scripts/                   # 启动与测试脚本
├── docs/                      # 说明文档与指南
└── examples/                  # Go 端到端示例
```

## 🚀 快速开始
### 1) 环境准备
- Go **1.21+**（建议与 `go.mod` 保持一致）
- Node.js **18+**（使用前端示例时）

### 2) 拉取代码与依赖
```bash
git clone https://github.com/your-org/einoflow.git
cd einoflow

# 安装后端依赖
make install   # 等同于 go mod download && go mod tidy

# （可选）安装前端依赖
cd web && npm install && cd ..
```

### 3) 配置环境变量
复制示例文件并填入至少一个可用的模型 Key（推荐豆包 `ARK_API_KEY`）：
```bash
cp .env.example .env
# 编辑 .env 文件，填入真实的 API Key
```

**⚠️ 安全提示**：
- `.env` 文件包含敏感信息，已被 `.gitignore` 保护，**不会提交到 Git**
- 所有 API Key 都从 `.env` 读取，代码中**没有硬编码**的密钥
- 详细安全配置请参考 [docs/SECURITY.md](docs/SECURITY.md)

常用字段说明：
- `ARK_API_KEY` / `ARK_BASE_URL`：豆包配置，默认模型 `doubao-seed-1-6-lite-251015`、`doubao-seed-1-6-vision-250815`。
- `OPENAI_API_KEY` / `OPENAI_BASE_URL`：OpenAI 备用配置（`gpt-4o`、`gpt-4o-mini` 等）。
- `SERVER_HOST` / `SERVER_PORT`：监听地址，默认 `0.0.0.0:8080`。
- `DB_PATH`：SQLite 路径，默认 `./data/einoflow.db`。
- `VECTOR_STORE_TYPE`：向量存储类型，默认 `memory`。

### 4) 启动服务
**方式 A：一键启动（推荐）**
```bash
./scripts/start-dev.sh   # 启动后端与前端，自动检查依赖
```

**方式 B：手动启动**
```bash
# 终端 1 - 启动后端
go run cmd/server/main.go

# 终端 2 - 启动前端（可选）
cd web
npm run dev
```

访问前端：`http://localhost:5173`，后端健康检查：
```bash
curl http://localhost:8080/health
```
返回 `"status":"ok"` 即表示后端正常。

## 📡 API 快速体验
以下示例使用默认豆包模型，依据需要可替换为 OpenAI 模型。

- **获取模型列表**
```bash
curl http://localhost:8080/api/v1/llm/models
```

- **基础对话**
```bash
curl -X POST http://localhost:8080/api/v1/llm/chat \
  -H "Content-Type: application/json" \
  -d '{
    "model": "doubao-seed-1-6-lite-251015",
    "messages": [{"role": "user", "content": "你好，介绍一下 Eino"}]
  }'
```

- **流式对话**
```bash
curl -N -X POST http://localhost:8080/api/v1/llm/chat/stream \
  -H "Content-Type: application/json" \
  -d '{
    "model": "doubao-seed-1-6-lite-251015",
    "messages": [{"role": "user", "content": "写一首关于编程的诗"}]
  }'
```

- **Agent 任务**
```bash
curl -X POST http://localhost:8080/api/v1/agent/run \
  -H "Content-Type: application/json" \
  -d '{"task": "分析 Go 语言与 Python 的优缺点并给出学习建议"}'
```

- **Chain 多步骤处理**
```bash
curl -X POST http://localhost:8080/api/v1/chain/run \
  -H "Content-Type: application/json" \
  -d '{
    "steps": ["翻译成英文", "总结成一句话", "用专业语气重写"],
    "input": "Go 是一门简洁高效的编程语言"
  }'
```

- **RAG 索引与查询**
```bash
# 索引文档
curl -X POST http://localhost:8080/api/v1/rag/index \
  -H "Content-Type: application/json" \
  -d '{
    "documents": [
      "Eino 是字节跳动开源的 LLM 应用开发框架",
      "Eino 支持 Chain、Agent、RAG、Graph 等功能"
    ]
  }'

# 查询
curl -X POST http://localhost:8080/api/v1/rag/query \
  -H "Content-Type: application/json" \
  -d '{"query": "Eino 有哪些主要功能？"}'
```

- **Graph 多步骤分析**（执行约 3–4 分钟）
```bash
│   ├── api/                   # 路由与各功能 Handler
│   ├── llm/                   # LLM 抽象与模型提供商
│   ├── agent/                 # ReAct Agent
│   ├── chain/                 # Chain 编排实现
│   ├── graph/                 # Graph 任务编排
│   ├── rag/                   # 文档加载、分块、检索
│   ├── memory/                # 对话历史与上下文
│   └── config/                # 配置加载
├── pkg/logger/                # 日志工具
├── web/                       # React/Vite 前端示例
├── docs/                      # 详细文档与指南
└── scripts/                   # 开发/启动脚本
```

## 🚀 快速开始
### 1. 环境准备
- Go **1.24** 及以上
- Node.js **18+**（运行前端示例）

### 2. 获取代码与依赖
```bash
# 克隆项目
git clone https://github.com/your-org/einoflow.git
cd einoflow

# 安装 Go 依赖
go mod download

# （可选）安装前端依赖
cd web && npm install
```

### 3. 配置环境变量
复制示例文件并填写密钥，至少保证字节豆包或 OpenAI 的 Key 有效：
```bash
cp .env.example .env
```
关键字段：
- `ARK_API_KEY` / `ARK_BASE_URL`：字节豆包配置（推荐）
- `OPENAI_API_KEY` / `OPENAI_BASE_URL`：OpenAI 备用配置
- `SERVER_PORT`、`SERVER_HOST`：服务监听地址
- `DB_PATH`、`VECTOR_STORE_TYPE`：存储配置

### 4. 启动服务
```bash
# 启动后端
go run cmd/server/main.go

# 可选：启动前端（另开终端）
cd web
npm run dev
```
访问前端：`http://localhost:5173`

### 5. 健康检查
```bash
curl http://localhost:8080/health
```
看到 `"status":"ok"` 即表示后端正常运行。

## 📡 API 快速体验
以下示例默认使用豆包 `ep-20241116153014-gfmhp` 模型，可根据需要替换。

- **获取模型列表**
```bash
curl http://localhost:8080/api/v1/llm/models
```

- **基础对话**
```bash
curl -X POST http://localhost:8080/api/v1/llm/chat \
  -H "Content-Type: application/json" \
  -d '{
    "model": "ep-20241116153014-gfmhp",
    "messages": [{"role": "user", "content": "你好，介绍一下 Eino"}]
  }'
```

- **流式对话**
```bash
curl -N -X POST http://localhost:8080/api/v1/llm/chat/stream \
  -H "Content-Type: application/json" \
  -d '{
    "model": "ep-20241116153014-gfmhp",
    "messages": [{"role": "user", "content": "写一首关于编程的诗"}]
  }'
```

- **Agent 任务**
```bash
curl -X POST http://localhost:8080/api/v1/agent/run \
  -H "Content-Type: application/json" \
  -d '{"task": "分析 Go 语言与 Python 的优缺点并给出学习建议"}'
```

- **Chain 多步骤处理**
```bash
curl -X POST http://localhost:8080/api/v1/chain/run \
  -H "Content-Type: application/json" \
  -d '{
    "steps": ["翻译成英文", "总结成一句话", "用专业语气重写"],
    "input": "Go 是一门简洁高效的编程语言"
  }'
```

- **RAG 索引与查询**
```bash
# 索引文档
curl -X POST http://localhost:8080/api/v1/rag/index \
  -H "Content-Type: application/json" \
  -d '{
    "documents": [
      "Eino 是字节跳动开源的 LLM 应用开发框架",
      "Eino 支持 Chain、Agent、RAG、Graph 等功能"
    ]
  }'

# 查询
curl -X POST http://localhost:8080/api/v1/rag/query \
  -H "Content-Type: application/json" \
  -d '{"query": "Eino 有哪些主要功能？"}'
```

- **Graph 多步骤分析**
```bash
curl -X POST http://localhost:8080/api/v1/graph/run \
  -H "Content-Type: application/json" \
  -d '{"query": "如何成为优秀的 Go 开发者？", "type": "multi_step"}'
```

更多调用示例、超时时间与响应说明可参考 `QUICK_REFERENCE.md` 与 `docs/DEMO_GUIDE.md`。

## 🧪 开发与测试
```bash
# 后端测试
go test ./...

# 前端类型与构建检查（如启用前端）
cd web && npm run build
```

常用脚本：`make run`/`make test`、`scripts/test-api.sh`（API 冒烟）、`scripts/test_stream.sh`（流式检查）。

## 🛠️ 常见问题速查
- **启动失败**：检查 `.env` 是否配置了有效的 `ARK_API_KEY` 或 `OPENAI_API_KEY`，并确认 8080/5173 端口未被占用。
- **流式响应无输出**：使用 `curl -N` 或前端 SSE 客户端；确保后端终端无错误日志。
- **Graph/Agent 耗时较长**：Agent 约 1–2 分钟，Graph 约 3–4 分钟属正常，请勿频繁刷新。
- **RAG 查询为空**：先调用 `/api/v1/rag/index` 完成索引，再发起查询；默认向量库为内存，重启会清空数据。
更多排查方法见 `QUICK_REFERENCE.md` 与 `DEBUG_GUIDE.md`。

## 📚 更多资源
- `QUICK_REFERENCE.md`：快速启动、FAQ、端点与耗时参考
- `DEBUG_GUIDE.md`：调试与日志查看技巧
- `TEST_CHECKLIST.md`：回归与端到端测试清单
- `FRONTEND_IMPLEMENTATION.md` / `web/SETUP.md`：前端架构与使用说明
- `FIXES_APPLIED.md` / `FIXES_ROUND*.md`：阶段性修复与现状
更多调用示例与响应格式参见 `docs/QUICKSTART.md` 与 `docs/DEMO_GUIDE.md`。

## 🧪 开发与测试
```bash
# 运行单元测试
go test ./...

# 前端构建检查（如启用前端）
cd web
npm run build
```

## 🛠️ 常见问题速查
- 启动失败：确认 `.env` 已填写有效的 `ARK_API_KEY` 或 `OPENAI_API_KEY`。
- 流式响应无输出：使用 `curl -N`，或在前端使用 SSE 客户端。
- Graph/Agent 任务耗时：Graph 约 3-4 分钟，Agent 约 1-2 分钟，属正常现象。
- RAG 查询为空：先调用 `/api/v1/rag/index` 索引文档，再发起查询。
更多排查方法见 `docs/TROUBLESHOOTING.md` 与 `QUICK_REFERENCE.md`。

## 📚 更多资源
- `docs/PROJECT_SUMMARY.md`：功能与结构概览
- `docs/COMPLETE_IMPLEMENTATION.md`：实现细节
- `docs/ADVANCED_FEATURES_IMPLEMENTATION.md`：高级能力说明
- `docs/FINAL_STATUS.md`：当前完成度与限制

欢迎提交 Issue 或 PR，一起完善 EinoFlow！
