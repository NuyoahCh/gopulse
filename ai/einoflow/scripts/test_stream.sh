#!/bin/bash

echo "=== 测试流式对话 ==="
echo ""
echo "发送请求到服务器..."
echo ""

curl -N -X POST http://localhost:8080/api/v1/llm/chat/stream \
  -H "Content-Type: application/json" \
  -d '{
    "model": "doubao-seed-1-6-lite-251015",
    "messages": [
      {
        "role": "user",
        "content": "你好，请用一句话介绍一下你自己"
      }
    ]
  }'

echo ""
echo ""
echo "=== 测试完成 ==="
