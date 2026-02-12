# EinoFlow 前端设置指南

## 🎯 功能概览

EinoFlow 前端提供了完整的 AI 功能界面：

- **AI 对话** - 与多个 AI 模型进行智能对话，支持流式输出
- **AI Agent** - 智能 Agent 完成复杂任务（写作、分析、代码生成）
- **RAG 检索** - 文档索引和智能检索问答
- **Graph 编排** - 多步骤图编排处理复杂问题

## 📦 安装依赖

```bash
cd web
npm install
```

这会安装所有必要的依赖，包括：
- React 18 + TypeScript
- React Router DOM（路由）
- Axios（API 请求）
- TailwindCSS（样式）
- Lucide React（图标）

## 🚀 启动开发服务器

### 1. 启动后端服务

首先确保后端服务正在运行：

```bash
# 在项目根目录
make run

# 或者
go run cmd/server/main.go
```

后端会在 `http://localhost:8080` 启动。

### 2. 启动前端开发服务器

```bash
cd web
npm run dev
```

前端会在 `http://localhost:5173` 启动（Vite 默认端口）。

## 🔧 配置说明

### API 代理配置

前端通过 Vite 代理连接后端 API。配置在 `vite.config.ts` 中：

```typescript
export default defineConfig({
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
});
```

这样前端的 `/api/*` 请求会自动转发到后端的 `http://localhost:8080/api/*`。

### 环境变量

如果需要自定义配置，可以创建 `.env` 文件：

```bash
# web/.env
VITE_API_BASE_URL=http://localhost:8080
```

## 📱 页面说明

### 1. 首页 (`/`)
- 功能展示和导航
- 技术栈介绍
- 快速开始入口

### 2. AI 对话 (`/chat`)
- 多模型选择（豆包、OpenAI）
- 实时对话
- 流式输出支持
- 对话历史

### 3. AI Agent (`/agent`)
- 任务描述输入
- Agent 执行
- 结果展示
- 执行步骤追踪

### 4. RAG 检索 (`/rag`)
- 文档索引
- 智能查询
- 来源追溯
- 统计信息

### 5. Graph 编排 (`/graph`)
- 复杂问题输入
- 多步骤执行
- 步骤可视化
- 最终结果

## 🎨 UI 组件

项目使用 shadcn/ui 风格的组件：

- `Button` - 按钮组件
- `Card` - 卡片容器
- `Badge` - 标签徽章

所有组件都在 `src/components/ui/` 目录中。

## 🔌 API 客户端

API 客户端在 `src/api/` 目录中：

- `client.ts` - Axios 基础配置
- `llm.ts` - LLM 对话 API
- `agent.ts` - Agent API
- `rag.ts` - RAG API
- `graph.ts` - Graph API

### 使用示例

```typescript
import { chat, chatStream } from './api/llm';

// 普通对话
const response = await chat({
  provider: 'ark',
  model: 'doubao-seed-1-6-lite-251015',
  messages: [{ role: 'user', content: 'Hello' }],
});

// 流式对话
await chatStream(
  { provider: 'ark', model: '...', messages: [...], stream: true },
  (content) => console.log('Chunk:', content),
  () => console.log('Done'),
  (error) => console.error('Error:', error)
);
```

## 🏗️ 项目结构

```
web/
├── src/
│   ├── api/              # API 客户端
│   │   ├── client.ts     # Axios 配置
│   │   ├── llm.ts        # LLM API
│   │   ├── agent.ts      # Agent API
│   │   ├── rag.ts        # RAG API
│   │   └── graph.ts      # Graph API
│   ├── components/       # UI 组件
│   │   └── ui/           # 基础 UI 组件
│   ├── pages/            # 页面组件
│   │   ├── HomePage.tsx  # 首页
│   │   ├── ChatPage.tsx  # 对话页面
│   │   ├── AgentPage.tsx # Agent 页面
│   │   ├── RAGPage.tsx   # RAG 页面
│   │   └── GraphPage.tsx # Graph 页面
│   ├── App.tsx           # 主应用和路由
│   └── main.tsx          # 入口文件
├── package.json          # 依赖配置
├── tsconfig.json         # TypeScript 配置
├── tailwind.config.ts    # TailwindCSS 配置
└── vite.config.ts        # Vite 配置
```

## 🐛 常见问题

### 1. 端口冲突

如果 5173 端口被占用，Vite 会自动使用下一个可用端口（5174、5175...）。

### 2. API 连接失败

确保：
- 后端服务正在运行（`http://localhost:8080`）
- `.env` 文件中配置了正确的 API Keys
- 防火墙没有阻止连接

### 3. TypeScript 错误

如果看到 TypeScript 错误，尝试：
```bash
npm run build  # 检查类型错误
```

### 4. 样式不生效

确保 TailwindCSS 正确配置：
```bash
# 重新安装依赖
rm -rf node_modules package-lock.json
npm install
```

## 📝 开发建议

### 添加新页面

1. 在 `src/pages/` 创建新组件
2. 在 `src/App.tsx` 添加路由
3. 在导航栏添加链接

### 添加新 API

1. 在 `src/api/` 创建新文件
2. 定义 TypeScript 接口
3. 实现 API 调用函数

### 样式定制

使用 TailwindCSS 工具类：
```tsx
<div className="rounded-lg bg-blue-600 p-4 text-white">
  Content
</div>
```

## 🚀 生产构建

```bash
npm run build
```

构建产物在 `dist/` 目录中，可以部署到任何静态文件服务器。

### 部署到 Nginx

```nginx
server {
    listen 80;
    server_name your-domain.com;
    root /path/to/dist;
    
    location / {
        try_files $uri $uri/ /index.html;
    }
    
    location /api {
        proxy_pass http://localhost:8080;
    }
}
```

## 🎉 完成！

现在你可以：

1. 访问 `http://localhost:5173` 查看前端
2. 点击导航栏切换不同功能
3. 测试所有 AI 功能

享受使用 EinoFlow！🚀
