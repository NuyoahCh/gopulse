# 🔍 调试指南

## 已修复的问题

### 1. ✅ Vite 代理配置
**问题**: 前端请求 `/api/*` 返回 404  
**原因**: `vite.config.ts` 缺少代理配置  
**修复**: 添加了 proxy 配置，将 `/api` 请求转发到 `http://localhost:8080`

```typescript
// vite.config.ts
proxy: {
  '/api': {
    target: 'http://localhost:8080',
    changeOrigin: true,
    secure: false,
  }
}
```

### 2. ✅ Agent API 字段不匹配
**问题**: Agent 返回 404 或数据显示错误  
**原因**: 后端返回 `answer` 字段，前端期望 `result` 字段  
**修复**: 统一使用 `answer` 字段

---

## 🚀 重启服务

修改配置后需要重启前端服务：

```bash
# 1. 停止当前服务（Ctrl+C）

# 2. 重新启动
cd web
npm run dev
```

或者使用启动脚本：
```bash
./scripts/start-dev.sh
```

---

## 🧪 测试步骤

### 1. 测试后端 API

```bash
# 测试模型列表
curl http://localhost:8080/api/v1/llm/models

# 测试对话
curl -X POST http://localhost:8080/api/v1/llm/chat \
  -H "Content-Type: application/json" \
  -d '{
    "provider": "ark",
    "model": "doubao-seed-1-6-lite-251015",
    "messages": [{"role": "user", "content": "Hello"}]
  }'

# 测试 Agent
curl -X POST http://localhost:8080/api/v1/agent/run \
  -H "Content-Type: application/json" \
  -d '{"task": "写一篇关于 Go 的文章"}'

# 测试 RAG 索引
curl -X POST http://localhost:8080/api/v1/rag/index \
  -H "Content-Type: application/json" \
  -d '{"documents": ["这是一个测试文档"]}'

# 测试 RAG 查询
curl -X POST http://localhost:8080/api/v1/rag/query \
  -H "Content-Type: application/json" \
  -d '{"query": "测试"}'

# 测试 Graph
curl -X POST http://localhost:8080/api/v1/graph/run \
  -H "Content-Type: application/json" \
  -d '{"input": "如何设计微服务架构？"}'
```

### 2. 测试前端页面

访问 `http://localhost:5173` 并测试：

#### AI 对话 (`/chat`)
1. 选择模型
2. 输入问题："你好"
3. 查看回答
4. 测试流式输出

#### AI Agent (`/agent`)
1. 输入任务："给我写一篇 Go 语言的文章"
2. 点击"运行 Agent"
3. 查看结果

#### RAG 检索 (`/rag`)
1. 切换到"索引文档"
2. 输入文档内容
3. 点击"开始索引"
4. 切换到"查询问答"
5. 输入问题
6. 查看答案

#### Graph 编排 (`/graph`)
1. 输入问题："如何设计微服务架构？"
2. 点击"运行 Graph"
3. 查看多步骤执行

---

## 🐛 常见问题排查

### 问题 1: 页面显示黑框

**可能原因**:
1. API 请求失败
2. 错误未正确处理
3. 样式加载失败

**排查步骤**:
```bash
# 1. 打开浏览器开发者工具（F12）
# 2. 查看 Console 标签页的错误信息
# 3. 查看 Network 标签页的请求状态
```

**解决方案**:
- 检查后端是否运行：`curl http://localhost:8080/api/v1/llm/models`
- 检查前端代理配置：`vite.config.ts`
- 查看浏览器控制台错误

### 问题 2: API 返回 404

**检查清单**:
- [ ] 后端服务是否运行
- [ ] 端口是否正确（8080）
- [ ] Vite 代理是否配置
- [ ] API 路径是否正确

**验证代理**:
```bash
# 在浏览器控制台执行
fetch('/api/v1/llm/models')
  .then(r => r.json())
  .then(console.log)
```

### 问题 3: 命令行无显示

**原因**: Gin 默认不输出请求日志到标准输出

**查看日志**:
```bash
# 后端日志会显示在启动后端的终端
# 查找类似这样的输出：
[GIN] 2025/11/17 - 08:51:38 | 200 | 412.292µs | ::1 | GET "/api/v1/llm/models"
```

### 问题 4: CORS 错误

**错误信息**:
```
Access to fetch at 'http://localhost:8080/api/v1/...' from origin 'http://localhost:5173' has been blocked by CORS policy
```

**解决**: 使用 Vite 代理（已配置）而不是直接请求后端

---

## 📊 API 字段对照表

### Agent API
| 前端字段 | 后端字段 | 类型 |
|---------|---------|------|
| task | task | string |
| answer | answer | string |

### RAG API
| 前端字段 | 后端字段 | 类型 |
|---------|---------|------|
| documents | documents | string[] |
| query | query | string |
| answer | answer | string |
| sources | sources | string[] |

### Graph API
| 前端字段 | 后端字段 | 类型 |
|---------|---------|------|
| input | input | string |
| result | result | string |
| steps | steps | object[] |

---

## 🔧 开发者工具使用

### Chrome DevTools

1. **Console 标签页**
   - 查看 JavaScript 错误
   - 查看 API 响应
   - 执行调试代码

2. **Network 标签页**
   - 查看所有网络请求
   - 检查请求/响应头
   - 查看响应数据

3. **Sources 标签页**
   - 设置断点
   - 单步调试
   - 查看变量值

### 快速调试命令

```javascript
// 在浏览器 Console 中执行

// 1. 测试 API 连接
fetch('/api/v1/llm/models')
  .then(r => r.json())
  .then(console.log)
  .catch(console.error);

// 2. 测试对话 API
fetch('/api/v1/llm/chat', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    provider: 'ark',
    model: 'doubao-seed-1-6-lite-251015',
    messages: [{ role: 'user', content: 'Hello' }]
  })
})
  .then(r => r.json())
  .then(console.log)
  .catch(console.error);

// 3. 查看当前页面状态
console.log('Location:', window.location.href);
console.log('React DevTools:', window.__REACT_DEVTOOLS_GLOBAL_HOOK__);
```

---

## ✅ 验证修复

### 1. 重启服务
```bash
# 停止所有服务（Ctrl+C）
# 重新启动
./scripts/start-dev.sh
```

### 2. 检查后端
```bash
curl http://localhost:8080/api/v1/llm/models
# 应该返回模型列表 JSON
```

### 3. 检查前端
访问 `http://localhost:5173`
- 首页应该正常显示
- 点击"对话"应该看到聊天界面（不是黑框）
- 可以选择模型
- 可以发送消息

### 4. 测试功能
- [ ] AI 对话正常
- [ ] Agent 执行成功
- [ ] RAG 索引和查询正常
- [ ] Graph 执行成功

---

## 📞 还有问题？

1. **查看浏览器控制台**
   - 按 F12 打开开发者工具
   - 查看 Console 和 Network 标签页

2. **查看后端日志**
   - 检查启动后端的终端输出
   - 查找错误信息

3. **检查配置**
   - `.env` 文件是否配置了 API Keys
   - `vite.config.ts` 是否有代理配置
   - 端口是否被占用

4. **重新安装依赖**
   ```bash
   cd web
   rm -rf node_modules package-lock.json
   npm install
   ```

现在应该可以正常工作了！🎉
