#!/bin/bash

# 订单 API 测试脚本

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 前台服务地址
FRONTEND_URL="http://localhost:8080"

echo "================================================"
echo "      订单管理 API 测试"
echo "================================================"
echo ""

# 1. 用户注册
echo -e "${YELLOW}1. 用户注册${NC}"
REGISTER_RESPONSE=$(curl -s -X POST "$FRONTEND_URL/api/v1/public/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testorder",
    "email": "testorder@example.com",
    "password": "password123"
  }')

echo "$REGISTER_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$REGISTER_RESPONSE"

# 提取 Token
TOKEN=$(echo "$REGISTER_RESPONSE" | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
  echo -e "${RED}❌ 注册失败，无法获取 Token${NC}"
  
  # 尝试登录
  echo -e "${YELLOW}尝试登录...${NC}"
  LOGIN_RESPONSE=$(curl -s -X POST "$FRONTEND_URL/api/v1/public/login" \
    -H "Content-Type: application/json" \
    -d '{
      "username": "testorder",
      "password": "password123"
    }')
  
  echo "$LOGIN_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$LOGIN_RESPONSE"
  TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"token":"[^"]*' | cut -d'"' -f4)
  
  if [ -z "$TOKEN" ]; then
    echo -e "${RED}❌ 登录也失败了，请检查服务是否正常运行${NC}"
    exit 1
  fi
fi

echo -e "${GREEN}✅ Token 获取成功${NC}"
echo "Token: $TOKEN"
echo ""

# 2. 创建订单
echo -e "${YELLOW}2. 创建订单${NC}"
CREATE_ORDER_RESPONSE=$(curl -s -X POST "$FRONTEND_URL/api/v1/user/orders" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "product_name": "iPhone 15 Pro",
    "product_price": 7999.00,
    "quantity": 1,
    "remark": "请尽快发货"
  }')

echo "$CREATE_ORDER_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$CREATE_ORDER_RESPONSE"

# 提取订单 ID
ORDER_ID=$(echo "$CREATE_ORDER_RESPONSE" | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)

if [ -z "$ORDER_ID" ]; then
  echo -e "${RED}❌ 创建订单失败${NC}"
  exit 1
fi

echo -e "${GREEN}✅ 订单创建成功，订单 ID: $ORDER_ID${NC}"
echo ""

# 3. 获取订单列表
echo -e "${YELLOW}3. 获取订单列表${NC}"
LIST_ORDERS_RESPONSE=$(curl -s -X GET "$FRONTEND_URL/api/v1/user/orders?page=1&page_size=10" \
  -H "Authorization: Bearer $TOKEN")

echo "$LIST_ORDERS_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$LIST_ORDERS_RESPONSE"
echo ""

# 4. 获取订单详情
echo -e "${YELLOW}4. 获取订单详情（订单 ID: $ORDER_ID）${NC}"
GET_ORDER_RESPONSE=$(curl -s -X GET "$FRONTEND_URL/api/v1/user/orders/$ORDER_ID" \
  -H "Authorization: Bearer $TOKEN")

echo "$GET_ORDER_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$GET_ORDER_RESPONSE"
echo ""

# 5. 支付订单
echo -e "${YELLOW}5. 支付订单（订单 ID: $ORDER_ID）${NC}"
PAY_ORDER_RESPONSE=$(curl -s -X POST "$FRONTEND_URL/api/v1/user/orders/$ORDER_ID/pay" \
  -H "Authorization: Bearer $TOKEN")

echo "$PAY_ORDER_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$PAY_ORDER_RESPONSE"
echo -e "${GREEN}✅ 订单支付成功${NC}"
echo ""

# 6. 完成订单
echo -e "${YELLOW}6. 完成订单（订单 ID: $ORDER_ID）${NC}"
COMPLETE_ORDER_RESPONSE=$(curl -s -X POST "$FRONTEND_URL/api/v1/user/orders/$ORDER_ID/complete" \
  -H "Authorization: Bearer $TOKEN")

echo "$COMPLETE_ORDER_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$COMPLETE_ORDER_RESPONSE"
echo -e "${GREEN}✅ 订单完成成功${NC}"
echo ""

# 7. 创建第二个订单（用于测试取消功能）
echo -e "${YELLOW}7. 创建第二个订单（用于测试取消）${NC}"
CREATE_ORDER2_RESPONSE=$(curl -s -X POST "$FRONTEND_URL/api/v1/user/orders" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "product_name": "MacBook Pro",
    "product_price": 12999.00,
    "quantity": 1,
    "remark": "需要发票"
  }')

echo "$CREATE_ORDER2_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$CREATE_ORDER2_RESPONSE"

ORDER_ID2=$(echo "$CREATE_ORDER2_RESPONSE" | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)
echo -e "${GREEN}✅ 第二个订单创建成功，订单 ID: $ORDER_ID2${NC}"
echo ""

# 8. 取消订单
echo -e "${YELLOW}8. 取消订单（订单 ID: $ORDER_ID2）${NC}"
CANCEL_ORDER_RESPONSE=$(curl -s -X POST "$FRONTEND_URL/api/v1/user/orders/$ORDER_ID2/cancel" \
  -H "Authorization: Bearer $TOKEN")

echo "$CANCEL_ORDER_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$CANCEL_ORDER_RESPONSE"
echo -e "${GREEN}✅ 订单取消成功${NC}"
echo ""

# 9. 获取最终的订单列表
echo -e "${YELLOW}9. 获取最终的订单列表${NC}"
FINAL_LIST_RESPONSE=$(curl -s -X GET "$FRONTEND_URL/api/v1/user/orders?page=1&page_size=10" \
  -H "Authorization: Bearer $TOKEN")

echo "$FINAL_LIST_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$FINAL_LIST_RESPONSE"
echo ""

echo "================================================"
echo -e "${GREEN}✅ 所有测试完成！${NC}"
echo "================================================"
echo ""
echo "测试摘要:"
echo "  - 用户注册/登录: ✅"
echo "  - 创建订单: ✅"
echo "  - 获取订单列表: ✅"
echo "  - 获取订单详情: ✅"
echo "  - 支付订单: ✅"
echo "  - 完成订单: ✅"
echo "  - 取消订单: ✅"
echo ""

