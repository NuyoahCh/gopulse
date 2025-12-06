#!/bin/bash

echo "=== EinoFlow 快速启动 ==="
echo ""

# 检查 .env 文件
if [ ! -f .env ]; then
    echo "❌ .env 文件不存在"
    echo "请复制 .env.example 到 .env 并配置 API Keys"
    exit 1
fi

# 检查依赖
echo "📦 检查依赖..."
go mod download
go mod tidy

# 创建必要的目录
echo "📁 创建目录..."
mkdir -p data/documents data/vector_store bin

# 编译项目
echo "🔨 编译项目..."
make build

echo ""
echo "✅ 准备完成！"
echo ""
echo "运行选项:"
echo "  1. 启动 Web 服务: make run"
echo "  2. 运行演示程序: make demo"
echo "  3. 查看帮助: make help"
echo ""