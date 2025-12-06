# 新功能问题排查指南

## 🐛 问题 1：流式输出出现两个消息框

### 症状
在 Chat 页面使用流式输出时，会同时显示两个消息框：
1. 一个空白的 loading 消息框（带转圈动画）
2. 一个正在接收流式内容的消息框

### 原因
代码在流式模式下同时：
- 设置了 `loading=true`，触发显示 loading 消息框
- 添加了空的 assistant 消息，用于接收流式内容

### 解决方案
✅ **已修复**：在流式模式下立即关闭 loading 状态

```tsx
// web/src/pages/ChatPage.tsx
if (streaming) {
  // 立即关闭 loading 状态，避免显示两个消息框
  setLoading(false);
  
  // 添加空的 assistant 消息用于流式更新
  const assistantMessage: ChatMessage = {
    role: 'assistant',
    content: '',
  };
  setMessages((prev) => [...prev, assistantMessage]);
  // ...
}
```

### 验证修复
1. 重启前端服务
2. 在 Chat 页面勾选"流式输出"
3. 发送消息
4. 应该只看到一个消息框，内容逐字显示

---

## 🐛 问题 2：天气查询无法调用

### 症状
在 Agent 页面输入天气相关问题，但无法获取实时天气信息。

### 可能原因

#### 原因 1：工具执行器未正确初始化
检查 `internal/api/agent_handler.go` 中是否正确创建了工具执行器。

**解决方案：**
```go
func NewAgentHandler(chatModel model.ChatModel) *AgentHandler {
    // 创建工具注册表
    toolRegistry := tools.NewRegistry("./data/einoflow.db", "./data/files")
    toolExecutor := tools.NewExecutor(toolRegistry)

    return &AgentHandler{
        chatModel:    chatModel,
        toolRegistry: toolRegistry,
        toolExecutor: toolExecutor,
    }
}
```

#### 原因 2：关键词识别失败
检查输入的问题是否包含天气关键词。

**支持的关键词：**
- 中文：天气、气温、温度、下雨、晴天、阴天
- 英文：weather

**测试命令：**
```bash
# 应该触发天气查询
curl -X POST http://localhost:8080/api/v1/agent/run \
  -H "Content-Type: application/json" \
  -d '{"task": "北京今天天气怎么样？"}'

# 不会触发天气查询（没有关键词）
curl -X POST http://localhost:8080/api/v1/agent/run \
  -H "Content-Type: application/json" \
  -d '{"task": "北京今天怎么样？"}'
```

#### 原因 3：城市名识别失败
检查输入的城市是否在支持列表中。

**支持的城市：**
北京、上海、广州、深圳、杭州、成都、重庆、武汉、西安、南京、天津、苏州、郑州、长沙、沈阳、青岛、厦门、大连、宁波、无锡、福州、济南、哈尔滨、长春

**解决方案：**
- 使用完整城市名（"北京"而不是"BJ"）
- 如果城市不在列表中，会使用默认城市"北京"

#### 原因 4：网络连接问题
天气 API（wttr.in）需要网络连接。

**检查网络：**
```bash
# 直接测试天气 API
curl "https://wttr.in/Beijing?format=j1&lang=zh"

# 如果返回 JSON 数据，说明网络正常
```

**如果网络不通：**
- 检查防火墙设置
- 检查代理配置
- 尝试使用 VPN

#### 原因 5：API 响应解析失败
检查后端日志，查看是否有解析错误。

**查看日志：**
```bash
# 启动后端时查看日志
go run cmd/server/main.go

# 发送天气查询请求
curl -X POST http://localhost:8080/api/v1/agent/run \
  -H "Content-Type: application/json" \
  -d '{"task": "北京今天天气怎么样？"}'

# 查看日志输出
[Agent] Task: 北京今天天气怎么样？, IsWeatherQuery: true, HasToolExecutor: true
[ToolExecutor] Executing tool: weather, params: map[location:北京]
[ToolExecutor] Calling weather tool for location: 北京
[ToolExecutor] Weather result: {...}, error: <nil>
```

### 调试步骤

#### 步骤 1：检查服务状态
```bash
# 检查后端是否运行
curl http://localhost:8080/health

# 检查前端是否运行
curl http://localhost:5173
```

#### 步骤 2：检查配置
```bash
# 确认 .env 文件存在
ls -la .env

# 确认 API Key 已配置
cat .env | grep API_KEY
```

#### 步骤 3：测试 API 直接调用
```bash
# 使用测试脚本
./test_weather.sh

# 或手动测试
curl -X POST http://localhost:8080/api/v1/agent/run \
  -H "Content-Type: application/json" \
  -d '{"task": "北京今天天气怎么样？"}' | jq '.'
```

#### 步骤 4：查看详细日志
在后端代码中已添加日志输出，重启后端查看：

```bash
go run cmd/server/main.go
```

发送请求后，应该看到类似输出：
```
[Agent] Task: 北京今天天气怎么样？, IsWeatherQuery: true, HasToolExecutor: true
[ToolExecutor] Executing tool: weather, params: map[location:北京]
[ToolExecutor] Calling weather tool for location: 北京
[ToolExecutor] Weather result: {"location":"北京",...}, error: <nil>
[Agent] Weather query result: 北京当前天气：..., error: <nil>
```

#### 步骤 5：测试天气 API
```bash
# 直接测试 wttr.in API
curl "https://wttr.in/北京?format=j1&lang=zh" | jq '.current_condition[0]'
```

### 快速修复检查清单

- [ ] 后端服务正在运行
- [ ] 前端服务正在运行
- [ ] .env 文件已配置 API Key
- [ ] 网络连接正常
- [ ] wttr.in API 可访问
- [ ] 输入包含天气关键词
- [ ] 城市名在支持列表中
- [ ] 工具执行器已初始化
- [ ] 查看后端日志无错误

---

## 🐛 问题 3：Markdown 渲染不正常

### 症状
AI 回复的 Markdown 内容没有正确格式化。

### 可能原因

#### 原因 1：前端依赖未安装
```bash
# 检查依赖
cd web
npm list react-markdown remark-gfm rehype-highlight

# 如果缺失，重新安装
npm install react-markdown remark-gfm rehype-highlight rehype-raw
```

#### 原因 2：CSS 样式未加载
检查浏览器控制台是否有 CSS 加载错误。

**解决方案：**
```bash
# 重新构建前端
cd web
npm run build
npm run dev
```

#### 原因 3：组件未正确导入
检查页面是否导入了 MarkdownRenderer 组件。

```tsx
import { MarkdownRenderer } from '../components/MarkdownRenderer';

// 使用
<MarkdownRenderer content={message.content} />
```

### 验证修复
1. 在 Chat 页面输入：`请用 Markdown 格式介绍 Python`
2. 检查返回内容是否有：
   - 标题样式
   - 代码高亮
   - 列表格式
   - 表格边框

---

## 🐛 问题 4：RAG 文件上传失败

### 症状
上传文件时返回错误或无响应。

### 可能原因

#### 原因 1：文件太大
**限制：** 10MB

**解决方案：**
- 压缩文件
- 分割成多个小文件
- 使用文本索引功能

#### 原因 2：文件格式不支持
**支持格式：** TXT, MD

**解决方案：**
- 转换为 TXT 格式
- 复制内容到文本索引区域

#### 原因 3：文件编码问题
**要求：** UTF-8 编码

**检查编码：**
```bash
file -I your_file.txt
```

**转换编码：**
```bash
iconv -f GBK -t UTF-8 input.txt > output.txt
```

### 验证修复
```bash
# 创建测试文件
echo "测试内容" > test.txt

# 上传测试
curl -X POST http://localhost:8080/api/v1/rag/upload \
  -F "file=@test.txt"

# 应该返回
{
  "message": "File uploaded and indexed successfully",
  "filename": "test.txt",
  "document_count": 1,
  "total_count": 1
}
```

---

## 📊 性能问题

### 问题：天气查询很慢

**原因：** wttr.in API 响应时间较长

**解决方案：**
1. 添加缓存（未来优化）
2. 使用其他天气 API
3. 设置更短的超时时间

### 问题：文件上传很慢

**原因：** 文件太大或网络慢

**解决方案：**
1. 压缩文件
2. 使用更快的网络
3. 增加超时时间

---

## 🔍 日志分析

### 查看后端日志
```bash
# 启动时查看所有日志
go run cmd/server/main.go 2>&1 | tee server.log

# 过滤天气相关日志
tail -f server.log | grep -E "\[Agent\]|\[ToolExecutor\]"

# 过滤 RAG 相关日志
tail -f server.log | grep -i "rag\|upload"
```

### 查看前端日志
打开浏览器开发者工具（F12）：
1. Console 标签：查看 JavaScript 错误
2. Network 标签：查看 API 请求
3. Application 标签：查看本地存储

---

## 💡 常见错误代码

### 400 Bad Request
- 请求参数错误
- 文件格式不支持
- 文件太大

### 500 Internal Server Error
- 后端服务异常
- API Key 无效
- 数据库错误

### 504 Gateway Timeout
- 天气 API 超时
- 网络连接问题

---

## 📞 获取帮助

如果以上方法都无法解决问题：

1. **查看完整文档**
   - `docs/NEW_FEATURES.md`
   - `docs/TESTING_GUIDE.md`
   - `QUICKSTART_NEW_FEATURES.md`

2. **提交 Issue**
   - 包含错误信息
   - 包含后端日志
   - 包含复现步骤

3. **联系支持**
   - GitHub Issues
   - 项目维护者

---

## ✅ 验证所有功能正常

运行完整测试：

```bash
# 1. 测试天气查询
./test_weather.sh

# 2. 测试文件上传
curl -X POST http://localhost:8080/api/v1/rag/upload \
  -F "file=@test_document.txt"

# 3. 测试 Markdown 渲染
# 在前端 Chat 页面输入：
# "请用 Markdown 格式介绍 Python，包括代码示例和表格"

# 4. 查看统计
curl http://localhost:8080/api/v1/rag/stats
```

如果所有测试通过，说明功能正常！🎉
