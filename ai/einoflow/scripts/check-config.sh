#!/bin/bash

# 配置检查脚本
echo "🔍 检查 EinoFlow 配置..."
echo ""

# 检查 .env 文件
echo "1️⃣  检查 .env 文件..."
if [ ! -f ".env" ]; then
    echo "❌ .env 文件不存在"
    echo "   请运行: cp .env.example .env"
    exit 1
else
    echo "✅ .env 文件存在"
fi
echo ""

# 检查 .env 是否被 Git 跟踪
echo "2️⃣  检查 .env 文件安全性..."
if git ls-files --error-unmatch .env 2>/dev/null; then
    echo "❌ 危险！.env 文件被 Git 跟踪"
    echo "   立即运行: git rm --cached .env"
    exit 1
else
    echo "✅ .env 文件未被 Git 跟踪（安全）"
fi
echo ""

echo "✅ 配置检查完成！"
