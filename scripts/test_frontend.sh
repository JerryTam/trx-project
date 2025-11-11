#!/bin/bash

# 前台服务测试脚本

BASE_URL="http://localhost:8080"
YELLOW='\033[1;33m'
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${YELLOW}========================================${NC}"
echo -e "${YELLOW}  前台服务 (Frontend) API 测试${NC}"
echo -e "${YELLOW}========================================${NC}"
echo ""

# 1. 健康检查
echo -e "${YELLOW}1. 健康检查${NC}"
curl -s "$BASE_URL/health" | jq '.'
echo ""

# 2. 用户注册
echo -e "${YELLOW}2. 用户注册${NC}"
TIMESTAMP=$(date +%s)
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/public/register" \
  -H "Content-Type: application/json" \
  -d "{
    \"username\": \"user_$TIMESTAMP\",
    \"email\": \"user_$TIMESTAMP@example.com\",
    \"password\": \"password123\"
  }")
echo "$REGISTER_RESPONSE" | jq '.'

# 提取 Token
USER_TOKEN=$(echo "$REGISTER_RESPONSE" | jq -r '.data.token')
echo ""
echo -e "${GREEN}✅ 用户 Token: $USER_TOKEN${NC}"
echo ""

# 3. 用户登录
echo -e "${YELLOW}3. 用户登录${NC}"
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/public/login" \
  -H "Content-Type: application/json" \
  -d "{
    \"username\": \"user_$TIMESTAMP\",
    \"password\": \"password123\"
  }")
echo "$LOGIN_RESPONSE" | jq '.'
echo ""

# 4. 获取个人信息（需要认证）
echo -e "${YELLOW}4. 获取个人信息（需要认证）${NC}"
if [ ! -z "$USER_TOKEN" ]; then
  curl -s "$BASE_URL/api/v1/user/profile" \
    -H "Authorization: Bearer $USER_TOKEN" | jq '.'
else
  echo -e "${RED}❌ 没有 Token，跳过${NC}"
fi
echo ""

# 5. 测试无效 Token
echo -e "${YELLOW}5. 测试无效 Token${NC}"
curl -s "$BASE_URL/api/v1/user/profile" \
  -H "Authorization: Bearer invalid_token" | jq '.'
echo ""

# 6. 测试缺少 Token
echo -e "${YELLOW}6. 测试缺少 Token${NC}"
curl -s "$BASE_URL/api/v1/user/profile" | jq '.'
echo ""

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}  前台服务测试完成！${NC}"
echo -e "${GREEN}========================================${NC}"

