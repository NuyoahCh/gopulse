#!/bin/bash

# 测试高级功能：上下文窗口管理和多模态支持

echo "🧪 测试高级功能"
echo "================================"

BASE_URL="http://localhost:8080/api/v1"

# 颜色定义
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "\n${BLUE}=== 测试 1: 上下文窗口管理 ===${NC}"
echo "发送超长对话，测试自动截断功能"

# 创建一个包含多条消息的长对话
curl -s -X POST $BASE_URL/llm/chat \
  -H "Content-Type: application/json" \
  -d '{
    "model": "doubao-seed-1-6-lite-251015",
    "messages": [
      {"role": "system", "content": "你是一个智能助手"},
      {"role": "user", "content": "第一个问题：什么是人工智能？"},
      {"role": "assistant", "content": "人工智能是计算机科学的一个分支..."},
      {"role": "user", "content": "第二个问题：什么是机器学习？"},
      {"role": "assistant", "content": "机器学习是人工智能的一个子集..."},
      {"role": "user", "content": "第三个问题：什么是深度学习？"},
      {"role": "assistant", "content": "深度学习是机器学习的一种方法..."},
      {"role": "user", "content": "第四个问题：什么是自然语言处理？"},
      {"role": "assistant", "content": "自然语言处理是人工智能的一个领域..."},
      {"role": "user", "content": "第五个问题：总结一下前面所有的概念"}
    ]
  }'
echo ""

echo -e "\n${YELLOW}💡 查看服务器日志，应该看到：${NC}"
echo "   'Context truncated: X -> Y messages, tokens: Z, available: W'"
echo ""

echo -e "\n${BLUE}=== 测试 2: 多模态支持（图像 URL）===${NC}"
echo "使用图像 URL 进行对话"

curl -s -X POST $BASE_URL/multimodal/chat \
  -H "Content-Type: application/json" \
  -d '{
    "text": "这张图片里有什么？",
    "image_url": "https://example.com/image.jpg"
  }'
echo ""

echo -e "\n${BLUE}=== 测试 3: 多模态支持（Base64 图像）===${NC}"
echo "使用 Base64 编码的图像"

curl -s -X POST $BASE_URL/multimodal/chat \
  -H "Content-Type: application/json" \
  -d '{
    "text": "描述这张图片",
    "image_b64": "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==",
    "mime_type": "image/png"
  }'
echo ""

echo -e "\n${GREEN}✅ 测试完成！${NC}"
echo -e "\n${YELLOW}📝 功能说明：${NC}"
echo ""
echo "1. 上下文窗口管理："
echo "   - 自动计算消息的 token 数"
echo "   - 超过限制时自动截断旧消息"
echo "   - 保留系统消息和最新对话"
echo ""
echo "2. 多模态支持："
echo "   - 支持图像 URL"
echo "   - 支持 Base64 编码的图像"
echo "   - 当前使用简化实现（将图像信息附加到文本）"
echo ""
echo "3. 查看服务器日志以确认功能正常工作"
echo ""
