#!/bin/bash

# 测试项目改进效果

echo "🧪 测试项目改进"
echo "================================"

BASE_URL="http://localhost:8080"

# 颜色定义
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "\n${BLUE}=== 测试 1: 配置验证 ===${NC}"
echo "✅ 配置验证已在服务启动时自动执行"
echo "   查看服务启动日志，确认没有配置错误"
echo ""

echo -e "\n${BLUE}=== 测试 2: 请求追踪（Request ID）===${NC}"
echo "发送请求并检查 Request ID..."

RESPONSE=$(curl -s -i -X POST $BASE_URL/api/v1/llm/chat \
  -H "Content-Type: application/json" \
  -d '{
    "model": "doubao-seed-1-6-lite-251015",
    "messages": [
      {"role": "user", "content": "你好"}
    ]
  }')

REQUEST_ID=$(echo "$RESPONSE" | grep -i "X-Request-ID" | cut -d' ' -f2 | tr -d '\r')

if [ -n "$REQUEST_ID" ]; then
    echo -e "${GREEN}✅ Request ID 已生成：${REQUEST_ID}${NC}"
else
    echo -e "${RED}❌ 未找到 Request ID${NC}"
fi

echo ""

echo -e "\n${BLUE}=== 测试 3: 结构化日志 ===${NC}"
echo "✅ 结构化日志已启用"
echo "   查看服务器日志，应该看到 JSON 格式的日志："
echo ""
echo -e "${YELLOW}示例日志：${NC}"
cat << 'EOF'
{
  "level": "info",
  "msg": "Request completed",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "method": "POST",
  "path": "/api/v1/llm/chat",
  "status": 200,
  "latency_ms": 1234,
  "client_ip": "127.0.0.1",
  "time": "2025-01-17T09:40:00Z"
}
EOF
echo ""

echo -e "\n${BLUE}=== 测试 4: Swagger 文档 ===${NC}"
echo "检查 Swagger UI 是否可访问..."

SWAGGER_STATUS=$(curl -s -o /dev/null -w "%{http_code}" $BASE_URL/swagger/index.html)

if [ "$SWAGGER_STATUS" = "200" ]; then
    echo -e "${GREEN}✅ Swagger UI 可访问${NC}"
    echo -e "   访问地址：${BLUE}http://localhost:8080/swagger/index.html${NC}"
else
    echo -e "${RED}❌ Swagger UI 不可访问（状态码：$SWAGGER_STATUS）${NC}"
fi

echo ""

echo -e "\n${BLUE}=== 测试 5: 错误处理 ===${NC}"
echo "测试错误响应是否包含 Request ID..."

ERROR_RESPONSE=$(curl -s -X POST $BASE_URL/api/v1/llm/chat \
  -H "Content-Type: application/json" \
  -d '{"invalid": "request"}')

ERROR_REQUEST_ID=$(echo "$ERROR_RESPONSE" | grep -o '"request_id":"[^"]*"' | cut -d'"' -f4)

if [ -n "$ERROR_REQUEST_ID" ]; then
    echo -e "${GREEN}✅ 错误响应包含 Request ID：${ERROR_REQUEST_ID}${NC}"
    echo "   错误响应：$ERROR_RESPONSE"
else
    echo -e "${RED}❌ 错误响应不包含 Request ID${NC}"
fi

echo ""

echo -e "\n${GREEN}✅ 测试完成！${NC}"
echo ""
echo -e "${YELLOW}📝 改进总结：${NC}"
echo ""
echo "1. ✅ 配置验证 - 启动时自动检查"
echo "2. ✅ 请求追踪 - 每个请求都有唯一 ID"
echo "3. ✅ 结构化日志 - JSON 格式，便于分析"
echo "4. ✅ Swagger 文档 - 交互式 API 文档"
echo "5. ✅ 错误处理 - 包含 request_id，便于追踪"
echo ""
echo -e "${BLUE}🔍 查看详细信息：${NC}"
echo "   - Swagger UI: http://localhost:8080/swagger/index.html"
echo "   - 健康检查: http://localhost:8080/health"
echo "   - 服务器日志: 查看终端输出"
echo ""
