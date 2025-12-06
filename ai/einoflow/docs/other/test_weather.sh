#!/bin/bash

echo "测试天气查询功能..."
echo ""

# 测试 1: 北京天气
echo "1. 测试北京天气查询："
curl -X POST http://localhost:8080/api/v1/agent/run \
  -H "Content-Type: application/json" \
  -d '{"task": "北京今天天气怎么样？"}' \
  2>/dev/null | jq '.'

echo ""
echo "---"
echo ""

# 测试 2: 上海天气
echo "2. 测试上海天气查询："
curl -X POST http://localhost:8080/api/v1/agent/run \
  -H "Content-Type: application/json" \
  -d '{"task": "上海的温度是多少？"}' \
  2>/dev/null | jq '.'

echo ""
echo "---"
echo ""

# 测试 3: 深圳天气
echo "3. 测试深圳天气查询："
curl -X POST http://localhost:8080/api/v1/agent/run \
  -H "Content-Type: application/json" \
  -d '{"task": "深圳会下雨吗？"}' \
  2>/dev/null | jq '.'

echo ""
echo "测试完成！"
