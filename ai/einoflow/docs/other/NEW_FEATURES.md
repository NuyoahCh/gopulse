# 新功能说明

## 📅 更新日期：2025-11-18

本次更新新增了三个重要功能，显著提升了 EinoFlow 的实用性和用户体验。

---

## 🌤️ 功能一：天气查询工具（MCP 集成）

### 功能描述
Agent 现在可以查询实时天气信息，支持中国主要城市的天气查询。

### 技术实现
- 使用 `wttr.in` 免费天气 API
- 支持中文城市名查询
- 自动识别天气相关问题
- 提供温度、湿度、风速、体感温度等详细信息

### 使用方法

#### API 调用
```bash
curl -X POST http://localhost:8080/api/v1/agent/run \
  -H "Content-Type: application/json" \
  -d '{
    "task": "北京今天天气怎么样？"
  }'
```

#### 前端使用
1. 访问 Agent 页面
2. 输入天气查询问题，例如：
   - "北京今天天气怎么样？"
   - "上海的气温是多少？"
   - "深圳会下雨吗？"
3. Agent 会自动调用天气工具并返回实时天气信息

### 支持的城市
北京、上海、广州、深圳、杭州、成都、重庆、武汉、西安、南京、天津、苏州、郑州、长沙、沈阳、青岛、厦门、大连、宁波、无锡、福州、济南、哈尔滨、长春等主要城市。

### 响应示例
```json
{
  "answer": "北京当前天气：晴，温度5°C，体感温度2°C，湿度45%，风速15 km/h"
}
```

---

## 📁 功能二：RAG 文件上传

### 功能描述
RAG 模块现在支持直接上传文件进行索引，不再局限于手动输入文本。

### 支持的文件格式
- `.txt` - 纯文本文件
- `.md` - Markdown 文件
- `.pdf` - PDF 文档（计划支持）
- `.doc/.docx` - Word 文档（计划支持）

### 文件限制
- 最大文件大小：10MB
- 自动分块：每 500 字符一块
- 支持中文分词

### 使用方法

#### API 调用
```bash
curl -X POST http://localhost:8080/api/v1/rag/upload \
  -F "file=@/path/to/your/document.txt"
```

#### 前端使用
1. 访问 RAG 页面
2. 点击"索引文档"标签
3. 在"上传文档文件"区域选择文件
4. 点击"上传并索引"按钮
5. 系统会自动分块并索引文件内容

### 响应示例
```json
{
  "message": "File uploaded and indexed successfully",
  "filename": "document.txt",
  "document_count": 25,
  "total_count": 125
}
```

### 工作流程
1. **文件上传** → 用户选择本地文件
2. **内容提取** → 读取文件内容
3. **智能分块** → 按 500 字符分块（保持语义完整性）
4. **向量化** → 使用 Embedding 模型生成向量
5. **存储索引** → 保存到向量数据库
6. **查询使用** → 在查询时自动检索相关内容

---

## 🎨 功能三：Markdown 渲染优化

### 功能描述
前端现在支持完整的 Markdown 渲染，AI 回复内容更加美观易读。

### 支持的 Markdown 特性
- ✅ 标题（H1-H6）
- ✅ 粗体、斜体、删除线
- ✅ 代码块（带语法高亮）
- ✅ 行内代码
- ✅ 列表（有序、无序）
- ✅ 表格
- ✅ 引用块
- ✅ 链接（自动在新标签页打开）
- ✅ 图片
- ✅ 任务列表
- ✅ 分隔线

### 代码高亮
支持多种编程语言的语法高亮：
- JavaScript/TypeScript
- Python
- Go
- Java
- C/C++
- Shell/Bash
- SQL
- JSON/YAML
- 等等...

### 样式特点
- **代码块**：深色主题，带语言标签
- **表格**：响应式设计，自动滚动
- **链接**：蓝色高亮，悬停下划线
- **引用**：左侧边框，斜体文字
- **列表**：清晰的项目符号和缩进

### 使用示例

当 AI 返回以下 Markdown 内容时：

```markdown
## 代码示例

这是一个 Python 函数：

\`\`\`python
def hello_world():
    print("Hello, World!")
\`\`\`

### 特性列表
- 支持语法高亮
- 自动格式化
- 响应式设计
```

前端会自动渲染为格式化的内容，包括：
- 标题样式
- 代码块高亮
- 列表格式化

---

## 🔧 技术细节

### 天气工具实现
```go
// internal/tools/weather.go
type WeatherTool struct {
    client *http.Client
}

func (t *WeatherTool) GetWeather(ctx context.Context, location string) (string, error) {
    url := fmt.Sprintf("https://wttr.in/%s?format=j1&lang=zh", location)
    // ... 调用 API 并解析响应
}
```

### 文件上传实现
```go
// internal/api/rag_handler.go
func (h *RAGHandler) UploadFile(c *gin.Context) {
    file, _ := c.FormFile("file")
    content, _ := io.ReadAll(file)
    chunks := h.splitText(string(content), 500)
    // ... 向量化并存储
}
```

### Markdown 渲染实现
```tsx
// web/src/components/MarkdownRenderer.tsx
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import rehypeHighlight from 'rehype-highlight';

export const MarkdownRenderer: React.FC<MarkdownRendererProps> = ({ content }) => {
  return (
    <ReactMarkdown
      remarkPlugins={[remarkGfm]}
      rehypePlugins={[rehypeHighlight]}
      components={{ /* 自定义组件样式 */ }}
    >
      {content}
    </ReactMarkdown>
  );
};
```

---

## 📊 性能优化

### 天气查询
- 超时时间：10 秒
- 降级策略：API 失败时返回友好提示
- 缓存：可考虑添加 5 分钟缓存（未来优化）

### 文件上传
- 分块大小：500 字符（可配置）
- 并发处理：支持多文件队列（未来优化）
- 进度显示：实时显示上传进度

### Markdown 渲染
- 懒加载：大文档分段渲染
- 代码高亮：使用 highlight.js
- 性能：React.memo 优化重渲染

---

## 🐛 已知问题

### 天气工具
- ❌ 部分小城市可能无法查询
- ❌ 国外城市支持有限
- ✅ 解决方案：提示用户使用其他渠道查询

### 文件上传
- ❌ PDF 解析需要额外库（计划中）
- ❌ Word 文档需要格式转换（计划中）
- ✅ 当前支持：TXT 和 MD 文件完美支持

### Markdown 渲染
- ❌ 数学公式（LaTeX）暂不支持
- ❌ 图表（Mermaid）暂不支持
- ✅ 基础 Markdown 完全支持

---

## 🚀 未来计划

### 短期（1-2 周）
- [ ] 添加 PDF 文件解析支持
- [ ] 优化天气查询缓存
- [ ] 支持更多城市天气查询

### 中期（1 个月）
- [ ] 添加数学公式渲染（KaTeX）
- [ ] 支持图表渲染（Mermaid）
- [ ] 文件上传进度条

### 长期（2-3 个月）
- [ ] 多语言天气查询
- [ ] 文件批量上传
- [ ] 自定义 Markdown 主题

---

## 📝 更新日志

### v1.1.0 (2025-11-18)
- ✨ 新增天气查询工具
- ✨ 新增 RAG 文件上传功能
- ✨ 优化 Markdown 渲染
- 🐛 修复配置文件硬编码问题
- 📚 完善安全配置文档

---

## 💡 使用建议

### 天气查询
- 使用完整的城市名（如"北京"而不是"BJ"）
- 询问具体问题（如"温度"、"会下雨吗"）
- 如果查询失败，尝试使用其他城市名

### 文件上传
- 优先使用 TXT 和 MD 格式
- 确保文件编码为 UTF-8
- 大文件建议先分割再上传

### Markdown 使用
- 在提问时可以要求 AI 使用 Markdown 格式
- 代码块记得指定语言以获得语法高亮
- 表格数据会自动格式化

---

## 🤝 贡献指南

如果你想为这些功能贡献代码：

1. Fork 项目
2. 创建功能分支：`git checkout -b feature/your-feature`
3. 提交更改：`git commit -am 'Add some feature'`
4. 推送分支：`git push origin feature/your-feature`
5. 提交 Pull Request

---

## 📞 反馈与支持

如有问题或建议，请：
- 提交 Issue：[GitHub Issues](https://github.com/your-org/einoflow/issues)
- 查看文档：`docs/` 目录
- 联系维护者：support@example.com
