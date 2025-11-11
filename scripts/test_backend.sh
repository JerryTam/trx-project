#!/bin/bash

# 后台服务测试脚本

BASE_URL="http://localhost:8081"
YELLOW='\033[1;33m'
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${YELLOW}========================================${NC}"
echo -e "${YELLOW}  后台服务 (Backend) API 测试${NC}"
echo -e "${YELLOW}========================================${NC}"
echo ""

# 生成管理员 Token
echo -e "${YELLOW}生成管理员 Token...${NC}"
cd "$(dirname "$0")"
ADMIN_TOKEN=$(go run generate_admin_token.go 2>/dev/null | grep "eyJ" | tr -d '\n')
cd - > /dev/null

if [ -z "$ADMIN_TOKEN" ]; then
  echo -e "${RED}❌ 无法生成管理员 Token，请手动运行：${NC}"
  echo "   go run scripts/generate_admin_token.go"
  exit 1
fi

echo -e "${GREEN}✅ 管理员 Token 已生成${NC}"
echo ""

# 1. 健康检查
echo -e "${YELLOW}1. 健康检查${NC}"
curl -s "$BASE_URL/health" | jq '.'
echo ""

# 2. 获取用户列表（需要管理员权限）
echo -e "${YELLOW}2. 获取用户列表（需要管理员权限）${NC}"
curl -s "$BASE_URL/api/v1/admin/users" \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq '.'
echo ""

# 3. 获取用户详情
echo -e "${YELLOW}3. 获取用户详情（ID=1）${NC}"
curl -s "$BASE_URL/api/v1/admin/users/1" \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq '.'
echo ""

# 4. 获取统计信息
echo -e "${YELLOW}4. 获取统计信息${NC}"
curl -s "$BASE_URL/api/v1/admin/statistics/users" \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq '.'
echo ""

# 5. 测试普通用户 Token（应该被拒绝）
echo -e "${YELLOW}5. 测试普通用户 Token（应该被拒绝）${NC}"
USER_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6InVzZXIiLCJyb2xlIjoidXNlciJ9.fake"
curl -s "$BASE_URL/api/v1/admin/users" \
  -H "Authorization: Bearer $USER_TOKEN" | jq '.'
echo ""

# 6. 测试无效 Token
echo -e "${YELLOW}6. 测试无效 Token${NC}"
curl -s "$BASE_URL/api/v1/admin/users" \
  -H "Authorization: Bearer invalid_token" | jq '.'
echo ""

# 7. 测试缺少 Token
echo -e "${YELLOW}7. 测试缺少 Token${NC}"
curl -s "$BASE_URL/api/v1/admin/users" | jq '.'
echo ""

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}  后台服务测试完成！${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo -e "${GREEN}管理员 Token (保存以便后续使用):${NC}"
echo "$ADMIN_TOKEN"

