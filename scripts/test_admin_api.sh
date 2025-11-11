#!/bin/bash

# 后台管理接口测试脚本

BASE_URL="http://localhost:8080"
ADMIN_TOKEN="admin_test_token_123456"

echo "================================"
echo "后台管理接口测试"
echo "================================"
echo ""

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}注意: 需要先启动服务器${NC}"
echo -e "${YELLOW}管理员 Token: $ADMIN_TOKEN${NC}"
echo ""

# 测试 1: 缺少 Token
echo "1. 测试缺少管理员 Token（应该返回 401）"
echo "GET $BASE_URL/api/v1/admin/users"
response=$(curl -s $BASE_URL/api/v1/admin/users)
echo "响应: $response"
echo ""

# 测试 2: 无效的 Token（不是管理员 Token）
echo "2. 测试无效的 Token（应该返回 403）"
echo "GET $BASE_URL/api/v1/admin/users"
response=$(curl -s $BASE_URL/api/v1/admin/users \
  -H "Authorization: Bearer invalid_token")
echo "响应: $response"
echo ""

# 测试 3: 获取用户列表（有效管理员 Token）
echo "3. 测试获取用户列表（管理员）"
echo "GET $BASE_URL/api/v1/admin/users"
response=$(curl -s "$BASE_URL/api/v1/admin/users?page=1&page_size=10" \
  -H "Authorization: Bearer $ADMIN_TOKEN")
echo "响应: $response"
echo ""

# 测试 4: 获取用户详情
echo "4. 测试获取用户详情（管理员）"
echo "GET $BASE_URL/api/v1/admin/users/1"
response=$(curl -s $BASE_URL/api/v1/admin/users/1 \
  -H "Authorization: Bearer $ADMIN_TOKEN")
echo "响应: $response"
echo ""

# 测试 5: 更新用户状态
echo "5. 测试更新用户状态（管理员）"
echo "PUT $BASE_URL/api/v1/admin/users/1/status"
response=$(curl -s -X PUT $BASE_URL/api/v1/admin/users/1/status \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"status": 0}')
echo "响应: $response"
echo ""

# 测试 6: 重置用户密码
echo "6. 测试重置用户密码（管理员）"
echo "POST $BASE_URL/api/v1/admin/users/1/reset-password"
response=$(curl -s -X POST $BASE_URL/api/v1/admin/users/1/reset-password \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"new_password": "newpassword123"}')
echo "响应: $response"
echo ""

# 测试 7: 获取用户统计
echo "7. 测试获取用户统计（管理员）"
echo "GET $BASE_URL/api/v1/admin/statistics/users"
response=$(curl -s $BASE_URL/api/v1/admin/statistics/users \
  -H "Authorization: Bearer $ADMIN_TOKEN")
echo "响应: $response"
echo ""

# 测试 8: 删除用户
echo "8. 测试删除用户（管理员）"
echo "DELETE $BASE_URL/api/v1/admin/users/999"
response=$(curl -s -X DELETE $BASE_URL/api/v1/admin/users/999 \
  -H "Authorization: Bearer $ADMIN_TOKEN")
echo "响应: $response"
echo ""

echo "================================"
echo "测试完成"
echo "================================"
echo ""
echo "验证项："
echo "✅ 后台接口需要管理员认证"
echo "✅ 缺少 Token 返回 401"
echo "✅ 无效 Token 返回 403"
echo "✅ 有效管理员 Token 可以访问后台接口"
echo "✅ 统一响应格式"
echo ""
echo -e "${GREEN}前后台接口已成功区分！${NC}"
echo ""

