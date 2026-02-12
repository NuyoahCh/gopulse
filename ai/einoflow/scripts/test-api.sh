#!/bin/bash
# API 测试脚本

echo "🧪 测试 EinoFlow API..."
echo ""

BASE_URL="http://localhost:8080"

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 测试函数
test_api() {
    local name=$1
    local method=$2
    local path=$3
    local data=$4
    
    echo -n "测试 $name... "
    
    if [ "$method" = "GET" ]; then
        response=$(curl -s -w "\n%{http_code}" "$BASE_URL$path")
    else
        response=$(curl -s -w "\n%{http_code}" -X "$method" "$BASE_URL$path" \
            -H "Content-Type: application/json" \
            -d "$data")
    fi
    
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')
    
    if [ "$http_code" = "200" ]; then
        echo -e "${GREEN}✅ 成功${NC} (HTTP $http_code)"
        if [ ! -z "$body" ]; then
            echo "   响应: $(echo $body | head -c 100)..."
        fi
    else
        echo -e "${RED}❌ 失败${NC} (HTTP $http_code)"
        if [ ! -z "$body" ]; then
            echo "   错误: $body"
        fi
    fi
    echo ""
}

# 检查后端是否运行
echo "🔍 检查后端服务..."
if ! curl -s "$BASE_URL/api/v1/llm/models" > /dev/null 2>&1; then
    echo -e "${RED}❌ 后端服务未运行！${NC}"
    echo "   请先启动后端: make run"
    exit 1
fi
echo -e "${GREEN}✅ 后端服务正在运行${NC}"
echo ""

# 测试 LLM API
echo "📝 测试 LLM API"
echo "─────────────────"
test_api "获取模型列表" "GET" "/api/v1/llm/models"

test_api "对话测试" "POST" "/api/v1/llm/chat" '{
  "provider": "ark",
  "model": "doubao-seed-1-6-lite-251015",
  "messages": [{"role": "user", "content": "你好"}]
}'

# 测试 Agent API
echo "🤖 测试 Agent API"
echo "─────────────────"
test_api "Agent 执行" "POST" "/api/v1/agent/run" '{
  "task": "解释什么是 Go 语言"
}'

# 测试 RAG API
echo "📚 测试 RAG API"
echo "─────────────────"
test_api "RAG 统计" "GET" "/api/v1/rag/stats"

test_api "RAG 索引" "POST" "/api/v1/rag/index" '{
  "documents": ["Go 是一门编程语言", "Go 由 Google 开发"]
}'

test_api "RAG 查询" "POST" "/api/v1/rag/query" '{
  "query": "Go 是什么？"
}'

# 测试 Graph API
echo "🔀 测试 Graph API"
echo "─────────────────"
test_api "Graph 执行" "POST" "/api/v1/graph/run" '{
  "input": "如何学习 Go 语言？"
}'

echo ""
echo "✅ 测试完成！"
