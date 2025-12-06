#!/bin/bash
# 安全检查脚本 - 在提交前运行

set -e

echo "🔍 开始安全检查..."
echo ""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

ERRORS=0

# 1. 检查 .env 是否被忽略
echo "📋 检查 1: .env 文件保护"
if git check-ignore .env > /dev/null 2>&1; then
    echo -e "${GREEN}✅ .env 已被 .gitignore 排除${NC}"
else
    echo -e "${RED}❌ 警告：.env 未被忽略！${NC}"
    echo "   请确保 .gitignore 中包含 .env"
    ERRORS=$((ERRORS + 1))
fi
echo ""

# 2. 检查 node_modules 是否被忽略
echo "📋 检查 2: node_modules 保护"
if git check-ignore web/node_modules/ > /dev/null 2>&1; then
    echo -e "${GREEN}✅ node_modules 已被 .gitignore 排除${NC}"
else
    echo -e "${YELLOW}⚠️  警告：node_modules 未被忽略${NC}"
fi
echo ""

# 3. 检查是否有硬编码的 API Key
echo "📋 检查 3: 硬编码密钥检测"
FOUND_KEYS=0

# 检查常见的 API Key 模式
if grep -r "sk-[a-zA-Z0-9]\{20,\}" --include="*.go" --include="*.js" --include="*.ts" . 2>/dev/null | grep -v ".env" | grep -v "node_modules" | grep -v ".git"; then
    echo -e "${RED}❌ 发现可能的 OpenAI API Key！${NC}"
    FOUND_KEYS=1
fi

if grep -r "AIza[a-zA-Z0-9_-]\{35\}" --include="*.go" --include="*.js" --include="*.ts" . 2>/dev/null | grep -v ".env" | grep -v "node_modules" | grep -v ".git"; then
    echo -e "${RED}❌ 发现可能的 Google API Key！${NC}"
    FOUND_KEYS=1
fi

# 检查是否有硬编码的密钥赋值
if grep -rE "(api_key|apiKey|API_KEY)\s*=\s*['\"][a-zA-Z0-9_-]{20,}['\"]" --include="*.go" --include="*.js" --include="*.ts" . 2>/dev/null | grep -v ".env" | grep -v "node_modules" | grep -v ".git" | grep -v "example"; then
    echo -e "${RED}❌ 发现可能的硬编码 API Key！${NC}"
    FOUND_KEYS=1
fi

if [ $FOUND_KEYS -eq 0 ]; then
    echo -e "${GREEN}✅ 未发现硬编码密钥${NC}"
else
    echo -e "${RED}❌ 请移除硬编码的密钥，使用环境变量！${NC}"
    ERRORS=$((ERRORS + 1))
fi
echo ""

# 4. 检查 .env 是否在暂存区（添加或修改，不包括删除）
echo "📋 检查 4: Git 暂存区检查"
if git diff --cached --name-only --diff-filter=AM | grep -q "^.env$"; then
    echo -e "${RED}❌ 警告：.env 文件被添加或修改在暂存区中！${NC}"
    echo "   运行: git reset .env"
    ERRORS=$((ERRORS + 1))
elif git diff --cached --name-only --diff-filter=D | grep -q "^.env$"; then
    echo -e "${GREEN}✅ .env 正在从 git 中删除（正确操作）${NC}"
else
    echo -e "${GREEN}✅ .env 不在暂存区中${NC}"
fi
echo ""

# 5. 检查是否有大文件
echo "📋 检查 5: 大文件检测"
LARGE_FILES=$(git diff --cached --name-only | xargs -I {} du -h {} 2>/dev/null | awk '$1 ~ /M$/ && $1+0 > 10 {print}' || true)
if [ -n "$LARGE_FILES" ]; then
    echo -e "${YELLOW}⚠️  发现大文件（>10MB）：${NC}"
    echo "$LARGE_FILES"
    echo "   考虑使用 Git LFS 或将其添加到 .gitignore"
else
    echo -e "${GREEN}✅ 未发现过大文件${NC}"
fi
echo ""

# 6. 检查 bin/ 目录
echo "📋 检查 6: 编译产物检查"
if git diff --cached --name-only | grep -q "^bin/"; then
    echo -e "${YELLOW}⚠️  警告：bin/ 目录中的文件在暂存区${NC}"
    echo "   编译产物通常不应提交"
else
    echo -e "${GREEN}✅ 无编译产物在暂存区${NC}"
fi
echo ""

# 总结
echo "================================"
if [ $ERRORS -eq 0 ]; then
    echo -e "${GREEN}✅ 安全检查通过！可以安全提交。${NC}"
    exit 0
else
    echo -e "${RED}❌ 发现 $ERRORS 个安全问题，请修复后再提交！${NC}"
    exit 1
fi
