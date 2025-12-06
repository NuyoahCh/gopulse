#!/bin/bash
# 开发环境启动脚本

set -e

echo "🚀 启动 EinoFlow 开发环境..."
echo ""

# 检查是否在项目根目录
if [ ! -f "go.mod" ]; then
    echo "❌ 错误：请在项目根目录运行此脚本"
    exit 1
fi

# 检查 .env 文件
if [ ! -f ".env" ]; then
    echo "⚠️  警告：.env 文件不存在"
    echo "   请复制 .env.example 并配置 API Keys"
    echo ""
    read -p "是否继续？(y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# 检查并安装前端依赖
echo "📦 检查前端依赖..."
cd web
if [ ! -d "node_modules" ] || [ ! -f "node_modules/.package-lock.json" ]; then
    echo "   安装前端依赖..."
    npm install
    echo "✅ 前端依赖安装完成"
else
    # 检查 package.json 是否有更新
    if [ "package.json" -nt "node_modules/.package-lock.json" ]; then
        echo "   检测到依赖更新，重新安装..."
        npm install
        echo "✅ 前端依赖更新完成"
    else
        echo "✅ 前端依赖已是最新"
    fi
fi
cd ..
echo ""

# 启动后端
echo "🔧 启动后端服务..."
echo "   地址: http://localhost:8080"
echo ""

# 在后台启动后端
go run cmd/server/main.go &
BACKEND_PID=$!

# 等待后端启动
echo "⏳ 等待后端启动..."
sleep 3

# 检查后端是否启动成功
if ! curl -s http://localhost:8080/api/v1/llm/models > /dev/null 2>&1; then
    echo "⚠️  后端可能未完全启动，但继续启动前端..."
else
    echo "✅ 后端启动成功"
fi
echo ""

# 启动前端
echo "🎨 启动前端服务..."
echo "   地址: http://localhost:5173"
echo ""

cd web
npm run dev &
FRONTEND_PID=$!
cd ..

echo ""
echo "✅ 开发环境启动完成！"
echo ""
echo "📱 访问地址："
echo "   前端: http://localhost:5173"
echo "   后端: http://localhost:8080"
echo ""
echo "🛑 停止服务："
echo "   按 Ctrl+C 停止"
echo ""

# 等待用户中断
trap "echo ''; echo '🛑 停止服务...'; kill $BACKEND_PID $FRONTEND_PID 2>/dev/null; echo '✅ 服务已停止'; exit 0" INT TERM

# 保持脚本运行
wait
