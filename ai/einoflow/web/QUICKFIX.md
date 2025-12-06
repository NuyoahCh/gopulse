# 🔧 快速修复指南

## ❌ 问题：react-router-dom 未找到

### 错误信息
```
Failed to resolve import "react-router-dom" from "src/App.tsx". Does the file exist?
```

### 原因
`react-router-dom` 依赖在 `package.json` 中添加了，但还没有安装到 `node_modules`。

### ✅ 解决方案

#### 方法 1: 手动安装（最快）
```bash
cd web
npm install
```

#### 方法 2: 重新运行启动脚本
更新后的 `start-dev.sh` 会自动检测并安装依赖：
```bash
./scripts/start-dev.sh
```

### 验证安装
```bash
cd web
npm list react-router-dom
```

应该看到：
```
einoflow-frontend@0.1.0
└── react-router-dom@6.30.2
```

---

## 🚀 完整启动步骤

### 首次启动
```bash
# 1. 安装前端依赖
cd web
npm install
cd ..

# 2. 启动开发环境
./scripts/start-dev.sh
```

### 后续启动
```bash
# 直接运行启动脚本即可
./scripts/start-dev.sh
```

---

## 📝 其他常见问题

### 1. 端口被占用
**错误**: `Error: listen EADDRINUSE: address already in use :::8080`

**解决**:
```bash
# 查找占用端口的进程
lsof -i :8080
lsof -i :5173

# 杀死进程
kill -9 <PID>
```

### 2. 依赖版本冲突
**解决**:
```bash
cd web
rm -rf node_modules package-lock.json
npm install
```

### 3. TypeScript 错误
**解决**:
```bash
cd web
npm run build  # 检查类型错误
```

### 4. 后端连接失败
**检查**:
- 后端是否在运行：`curl http://localhost:8080/api/v1/llm/models`
- `.env` 文件是否配置了 API Keys
- 防火墙是否阻止了连接

---

## ✅ 验证一切正常

### 1. 检查后端
```bash
curl http://localhost:8080/api/v1/llm/models
```

应该返回模型列表的 JSON。

### 2. 检查前端
访问 `http://localhost:5173`，应该看到 EinoFlow 首页。

### 3. 测试功能
- 点击"对话"，测试 AI 对话
- 点击"Agent"，测试 Agent 功能
- 点击"RAG"，测试文档索引和查询
- 点击"Graph"，测试多步骤处理

---

## 🎉 现在应该可以正常使用了！

如果还有问题，请检查：
1. Node.js 版本 >= 16
2. Go 版本 >= 1.21
3. `.env` 文件已配置
4. 所有依赖已安装

需要帮助？查看 `web/SETUP.md` 获取详细说明。
