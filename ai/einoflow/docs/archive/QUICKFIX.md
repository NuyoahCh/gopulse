# 快速修复指南

## 问题 1: 流式对话返回空内容

### 原因
流式响应中每个 chunk 的 `Content` 可能为空，这是正常的流式行为。某些模型会分多个块发送内容。

### 已修复
- ✅ 更新了 `ProcessStream` 函数，只在有内容时才发送
- ✅ 过滤掉空内容块

### 测试方法
```bash
# 1. 启动服务器
make run

# 2. 在另一个终端测试流式响应
./scripts/test_stream.sh
```

## 问题 2: Demo 程序模型 404 错误

### 原因
使用了错误的模型 ID: `ep-20241116153014-gfmhp`

### 已修复
- ✅ 将模型 ID 改为: `doubao-seed-1-6-lite-251015`
- ✅ 这是豆包支持的正确模型名称

### 测试方法
```bash
# 运行 demo 程序
make demo

# 选择任意功能测试
# 1. 基础对话
# 2. 流式对话
# 3. Agent 工具调用
# 4. Graph 多步骤处理
```

## 验证修复

### 1. 编译检查
```bash
go build ./...
```

### 2. 运行服务器
```bash
make run
```

### 3. 测试基础对话
```bash
curl -X POST http://localhost:8080/api/v1/llm/chat \
  -H "Content-Type: application/json" \
  -d '{
    "model": "doubao-seed-1-6-lite-251015",
    "messages": [{"role": "user", "content": "你好"}]
  }'
```

### 4. 测试流式对话
```bash
curl -N -X POST http://localhost:8080/api/v1/llm/chat/stream \
  -H "Content-Type: application/json" \
  -d '{
    "model": "doubao-seed-1-6-lite-251015",
    "messages": [{"role": "user", "content": "讲个笑话"}]
  }'
```

### 5. 查看可用模型
```bash
curl http://localhost:8080/api/v1/llm/models
```

## 重要提示

### 使用正确的模型 ID
豆包支持的模型：
- ✅ `doubao-seed-1-6-lite-251015` - Lite 版本（推荐）
- ✅ `doubao-seed-1-6-vision-250815` - Vision 版本
- ❌ `ep-20241116153014-gfmhp` - 这是 endpoint ID，不是模型名称

### 检查 API 配置
确保 `.env` 文件配置正确：
```bash
ARK_API_KEY="your-api-key"
ARK_BASE_URL=https://ark.cn-beijing.volces.com/api/v3
```

### 如果仍有问题

1. **查看服务器日志**
   ```bash
   make run
   # 观察日志输出
   ```

2. **测试健康检查**
   ```bash
   curl http://localhost:8080/health
   ```

3. **检查网络连接**
   ```bash
   curl -I https://ark.cn-beijing.volces.com
   ```

4. **查看详细错误**
   修改 `.env`:
   ```bash
   LOG_LEVEL=debug
   LOG_FORMAT=text
   ```

## 下一步

所有修复已完成，你现在可以：
1. ✅ 运行 `make demo` 测试所有功能
2. ✅ 运行 `make run` 启动 Web 服务
3. ✅ 使用 API 进行对话和流式对话
4. ✅ 查看 `docs/TROUBLESHOOTING.md` 了解更多故障排查信息
