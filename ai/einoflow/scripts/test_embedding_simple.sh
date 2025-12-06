#!/bin/bash

# 测试真实 Embedding 功能（不依赖 jq）

echo "🧪 测试 ARK Embedding 集成"
echo "================================"

BASE_URL="http://localhost:8080/api/v1/rag"

# 颜色定义
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "\n${BLUE}步骤 1: 清空现有数据${NC}"
curl -s -X DELETE $BASE_URL/clear
echo ""

echo -e "\n${BLUE}步骤 2: 使用真实 Embedding 索引文档${NC}"
echo "（这将使用 doubao-embedding-large-text-250515 模型）"
curl -s -X POST $BASE_URL/index \
  -H "Content-Type: application/json" \
  -d '{
    "documents": [
      "人工智能是计算机科学的一个分支，致力于创建能够执行通常需要人类智能的任务的系统",
      "机器学习是人工智能的一个子集，它使计算机能够从数据中学习而无需明确编程",
      "深度学习是机器学习的一种方法，使用多层神经网络来学习数据的表示",
      "自然语言处理是人工智能的一个领域，专注于计算机与人类语言之间的交互",
      "计算机视觉是人工智能的一个领域，使计算机能够理解和解释视觉信息"
    ]
  }'
echo ""

echo -e "\n${BLUE}步骤 3: 查看存储的文档数量${NC}"
curl -s $BASE_URL/stats
echo ""

echo -e "\n${BLUE}步骤 4: 测试语义搜索（真实 Embedding）${NC}"
echo "查询：什么是深度学习？"
curl -s -X POST $BASE_URL/query \
  -H "Content-Type: application/json" \
  -d '{
    "query": "什么是深度学习？"
  }'
echo ""

echo -e "\n${BLUE}步骤 5: 测试另一个查询${NC}"
echo "查询：NLP 是什么？"
curl -s -X POST $BASE_URL/query \
  -H "Content-Type: application/json" \
  -d '{
    "query": "NLP 是什么？"
  }'
echo ""

echo -e "\n${GREEN}✅ 测试完成！${NC}"
echo -e "\n${YELLOW}💡 检查服务器日志：${NC}"
echo "   查找这条日志："
echo -e "   ${GREEN}'Using ARK Embedding model: doubao-embedding-large-text-250515'${NC}"
echo ""
echo "   如果看到这条日志 → 真实 Embedding 已启用 ✅"
echo "   如果看到 'Using simple character-based embedding' → 使用简单 embedding ⚠️"
echo ""
