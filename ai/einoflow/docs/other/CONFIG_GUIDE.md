# 配置指南

## 📋 快速配置

### 1. 复制配置文件

```bash
cp .env.example .env
```

### 2. 编辑 .env 文件

使用你喜欢的编辑器打开 `.env` 文件：

```bash
vim .env
# 或
nano .env
```

### 3. 填入真实的 API Key

至少配置一个 LLM API Key（推荐使用豆包）：

```bash
# 字节豆包配置（推荐）
ARK_API_KEY=你的真实密钥
ARK_BASE_URL=https://ark.cn-beijing.volces.com/api/v3

# 或者使用 OpenAI
OPENAI_API_KEY=你的真实密钥
OPENAI_BASE_URL=https://api.openai.com/v1
```

### 4. 验证配置

```bash
./scripts/check-config.sh
```

## 🔑 获取 API Key

### 字节豆包（推荐）

1. 访问 [火山引擎控制台](https://console.volcengine.com/ark)
2. 注册/登录账号
3. 创建推理接入点
4. 复制 API Key

### OpenAI

1. 访问 [OpenAI Platform](https://platform.openai.com/)
2. 注册/登录账号
3. 进入 API Keys 页面
4. 创建新的 API Key

## 🛡️ 安全最佳实践

### ✅ 正确做法

1. **所有密钥都在 .env 文件中**
   - `.env` 文件已被 `.gitignore` 保护
   - 不会被提交到 Git 仓库

2. **代码中没有硬编码**
   - 所有配置都通过 `getEnv()` 从环境变量读取
   - 默认值为空字符串，强制用户配置

3. **团队协作**
   - 每个开发者维护自己的 `.env` 文件
   - 通过 `.env.example` 共享配置模板

### ❌ 错误做法

1. **不要将 .env 提交到 Git**
   ```bash
   # 如果不小心提交了
   git rm --cached .env
   git commit -m "Remove .env file"
   ```

2. **不要在代码中硬编码密钥**
   ```go
   // ❌ 错误
   ArkAPIKey: "feabe6d9-8244-4e30-aff4-e7ad167a2ae9"
   
   // ✅ 正确
   ArkAPIKey: getEnv("ARK_API_KEY", "")
   ```

3. **不要在日志中打印密钥**
   ```go
   // ❌ 错误
   logger.Info("API Key: " + cfg.ArkAPIKey)
   
   // ✅ 正确
   logger.Info("API Key configured: " + strconv.FormatBool(cfg.ArkAPIKey != ""))
   ```

## 📝 配置项说明

### LLM API 配置

| 配置项 | 说明 | 必填 | 默认值 |
|--------|------|------|--------|
| `ARK_API_KEY` | 字节豆包 API Key | 至少一个 | - |
| `ARK_BASE_URL` | 字节豆包 API 地址 | 否 | `https://ark.cn-beijing.volces.com/api/v3` |
| `ARK_EMBEDDING_MODEL` | 豆包 Embedding 模型 | 否 | `doubao-embedding-large-text-250515` |
| `OPENAI_API_KEY` | OpenAI API Key | 至少一个 | - |
| `OPENAI_BASE_URL` | OpenAI API 地址 | 否 | `https://api.openai.com/v1` |
| `ANTHROPIC_API_KEY` | Anthropic API Key | 至少一个 | - |

### 服务配置

| 配置项 | 说明 | 默认值 |
|--------|------|--------|
| `SERVER_HOST` | 服务监听地址 | `0.0.0.0` |
| `SERVER_PORT` | 服务监听端口 | `8080` |
| `LOG_LEVEL` | 日志级别 | `info` |
| `LOG_FORMAT` | 日志格式 | `json` |

### 数据库配置

| 配置项 | 说明 | 默认值 |
|--------|------|--------|
| `DB_PATH` | SQLite 数据库路径 | `./data/einoflow.db` |
| `VECTOR_STORE_TYPE` | 向量存储类型 | `memory` |
| `VECTOR_DIM` | 向量维度 | `1536` |

## 🔧 环境变量优先级

配置加载顺序（优先级从高到低）：

1. **系统环境变量**（最高优先级）
2. **.env 文件**
3. **代码默认值**（最低优先级）

示例：

```bash
# 临时覆盖配置
SERVER_PORT=9000 go run cmd/server/main.go

# 或在 .env 中配置
SERVER_PORT=9000
```

## 🚀 不同环境配置

### 开发环境

```bash
# .env
LOG_LEVEL=debug
LOG_FORMAT=text
VECTOR_STORE_TYPE=memory
```

### 生产环境

```bash
# .env
LOG_LEVEL=info
LOG_FORMAT=json
VECTOR_STORE_TYPE=persistent
```

## 📚 相关文档

- [安全配置指南](./SECURITY.md) - 详细的安全最佳实践
- [快速开始](../README.md#快速开始) - 项目启动指南
- [API 文档](http://localhost:8080/swagger/index.html) - API 接口文档

## ❓ 常见问题

### Q: 为什么启动失败提示 "at least one LLM API key must be configured"？

A: 你需要在 `.env` 文件中至少配置一个 LLM API Key（`ARK_API_KEY`、`OPENAI_API_KEY` 或 `ANTHROPIC_API_KEY`）。

### Q: 可以同时配置多个 LLM 提供商吗？

A: 可以！你可以同时配置豆包、OpenAI 和 Anthropic，系统会根据模型名称自动选择对应的提供商。

### Q: 如何切换不同的模型？

A: 在 API 请求中指定 `model` 参数即可，例如：
```json
{
  "model": "doubao-seed-1-6-lite-251015",
  "messages": [...]
}
```

### Q: .env 文件会被提交到 Git 吗？

A: 不会！`.env` 文件已经在 `.gitignore` 中，不会被 Git 跟踪。你可以运行 `./scripts/check-config.sh` 验证。
