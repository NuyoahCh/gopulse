# 故障排查指南

## 常见问题

### 1. 流式响应返回空内容

**症状**: SSE 流返回 `data: {"content":"","done":false}`

**可能原因**:
1. 模型 ID 不正确
2. 流式响应格式问题
3. API 配置错误

**解决方案**:

#### 检查模型 ID
确保使用正确的模型 ID。豆包支持的模型：
- `doubao-seed-1-6-lite-251015` (推荐，Lite 版本)
- `doubao-seed-1-6-vision-250815` (Vision 版本)

#### 测试流式响应
```bash
# 启动服务器
make run

# 在另一个终端运行测试
./scripts/test_stream.sh
```

#### 检查 API 配置
确保 `.env` 文件中的配置正确：
```bash
ARK_API_KEY="your-api-key"
ARK_BASE_URL=https://ark.cn-beijing.volces.com/api/v3
```

### 2. 模型 404 错误

**错误信息**: 
```
Error code: 404 - {"code":"InvalidEndpointOrModel.NotFound","message":"The model or endpoint ep-xxx does not exist..."}
```

**原因**: 使用了不存在的模型 ID 或 endpoint ID

**解决方案**:

1. **使用模型名称而不是 endpoint ID**
   ```go
   // ❌ 错误 - 使用 endpoint ID
   Model: "ep-20241116153014-gfmhp"
   
   // ✅ 正确 - 使用模型名称
   Model: "doubao-seed-1-6-lite-251015"
   ```

2. **查看可用模型列表**
   ```bash
   curl http://localhost:8080/api/v1/llm/models
   ```

3. **更新配置文件**
   - 修改 `examples/complete_demo.go` 中的模型 ID
   - 修改 API 请求中的 `model` 字段

### 3. API Key 无效

**错误信息**: `401 Unauthorized` 或 `403 Forbidden`

**解决方案**:
1. 检查 `.env` 文件中的 API Key 是否正确
2. 确认 API Key 有访问权限
3. 检查 BaseURL 是否正确

### 4. 流式响应不显示

**症状**: 流式请求没有输出或输出延迟

**可能原因**:
1. 缓冲问题
2. 客户端不支持 SSE
3. 代理或防火墙拦截

**解决方案**:
```bash
# 使用 curl 测试（-N 禁用缓冲）
curl -N -X POST http://localhost:8080/api/v1/llm/chat/stream \
  -H "Content-Type: application/json" \
  -d '{"model":"doubao-seed-1-6-lite-251015","messages":[{"role":"user","content":"你好"}]}'
```

## 调试技巧

### 1. 启用详细日志
修改 `.env`:
```bash
LOG_LEVEL=debug
LOG_FORMAT=text
```

### 2. 查看服务器日志
```bash
# 启动服务器并查看日志
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

### 4. 检查健康状态
```bash
curl http://localhost:8080/health
```

## 配置示例

### 正确的 .env 配置
```bash
# 字节豆包配置
ARK_API_KEY="your-actual-api-key"
ARK_BASE_URL=https://ark.cn-beijing.volces.com/api/v3

# OpenAI 配置（可选）
OPENAI_API_KEY="sk-xxx"
OPENAI_BASE_URL=https://api.openai.com/v1

# 服务配置
SERVER_PORT=8080
SERVER_HOST=0.0.0.0
LOG_LEVEL=info
LOG_FORMAT=json
```

### 正确的 API 请求
```json
{
  "model": "doubao-seed-1-6-lite-251015",
  "messages": [
    {
      "role": "user",
      "content": "你好"
    }
  ],
  "temperature": 0.7,
  "max_tokens": 1000
}
```

## 获取帮助

如果问题仍未解决：
1. 查看服务器日志获取详细错误信息
2. 确认网络连接正常
3. 检查 API 配额是否用尽
4. 联系 API 提供商支持
