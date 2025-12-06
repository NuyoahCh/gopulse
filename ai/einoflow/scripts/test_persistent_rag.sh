#!/bin/bash

# RAG 持久化存储测试脚本

echo "🧪 测试 RAG 持久化存储功能"
echo "================================"

BASE_URL="http://localhost:8080/api/v1/rag"

# 颜色定义
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "\n${BLUE}步骤 1: 清空现有数据${NC}"
curl -s -X DELETE $BASE_URL/clear | jq '.'

echo -e "\n${BLUE}步骤 2: 查看当前数据（应该为空）${NC}"
curl -s $BASE_URL/stats | jq '.'

echo -e "\n${BLUE}步骤 3: 索引测试文档${NC}"
curl -s -X POST $BASE_URL/index \
  -H "Content-Type: application/json" \
  -d '{
    "documents": [
      "Eino 是字节跳动开源的 LLM 应用框架",
      "Eino 支持 Chain、Agent、RAG、Graph 等功能",
      "Eino 使用 Go 语言编写，性能优秀",
      "EinoFlow 是基于 Eino 构建的完整应用",
      "RAG 系统现在支持 SQLite 持久化存储"
    ]
  }' | jq '.'

echo -e "\n${BLUE}步骤 4: 查看存储的文档${NC}"
curl -s $BASE_URL/stats | jq '.'

echo -e "\n${BLUE}步骤 5: 测试查询功能${NC}"
curl -s -X POST $BASE_URL/query \
  -H "Content-Type: application/json" \
  -d '{
    "query": "Eino 有哪些功能？"
  }' | jq '.'

echo -e "\n${GREEN}✅ 测试完成！${NC}"
echo -e "${YELLOW}💡 提示: 现在重启服务，数据仍然会保留！${NC}"
echo ""
echo "验证持久化："
echo "  1. 停止服务 (Ctrl+C)"
echo "  2. 重新启动 (make run)"
echo "  3. 运行: curl http://localhost:8080/api/v1/rag/stats"
echo "  4. 数据应该仍然存在！"
echo ""
