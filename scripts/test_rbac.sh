#!/bin/bash

# RBAC 权限系统测试脚本

YELLOW='\033[1;33m'
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

BASE_URL="http://localhost:8081"

echo -e "${YELLOW}========================================${NC}"
echo -e "${YELLOW}  RBAC 权限系统测试${NC}"
echo -e "${YELLOW}========================================${NC}"
echo ""

# 生成超级管理员 Token
echo -e "${YELLOW}生成超级管理员 Token...${NC}"
SUPER_ADMIN_TOKEN=$(go run scripts/generate_admin_with_role.go 2>/dev/null | grep -A 1 "Token:" | tail -1 | tr -d '\n ')

if [ -z "$SUPER_ADMIN_TOKEN" ]; then
    echo -e "${RED}❌ 无法生成超级管理员 Token${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Token 已生成${NC}"
echo ""

# 1. 测试查看角色列表
echo -e "${YELLOW}1. 测试查看角色列表（需要 rbac:manage 权限）${NC}"
curl -s "$BASE_URL/api/v1/admin/rbac/roles" \
    -H "Authorization: Bearer $SUPER_ADMIN_TOKEN" | jq '.'
echo ""

# 2. 测试查看权限列表
echo -e "${YELLOW}2. 测试查看权限列表（需要 rbac:manage 权限）${NC}"
curl -s "$BASE_URL/api/v1/admin/rbac/permissions" \
    -H "Authorization: Bearer $SUPER_ADMIN_TOKEN" | jq '.'
echo ""

# 3. 测试查看用户角色
echo -e "${YELLOW}3. 测试查看用户角色（用户 ID=1）${NC}"
curl -s "$BASE_URL/api/v1/admin/users/1/roles" \
    -H "Authorization: Bearer $SUPER_ADMIN_TOKEN" | jq '.'
echo ""

# 4. 测试查看用户权限
echo -e "${YELLOW}4. 测试查看用户权限（用户 ID=1）${NC}"
curl -s "$BASE_URL/api/v1/admin/users/1/permissions" \
    -H "Authorization: Bearer $SUPER_ADMIN_TOKEN" | jq '.'
echo ""

# 5. 测试查看用户列表（需要 user:read 权限）
echo -e "${YELLOW}5. 测试查看用户列表（需要 user:read 权限）${NC}"
curl -s "$BASE_URL/api/v1/admin/users" \
    -H "Authorization: Bearer $SUPER_ADMIN_TOKEN" | jq '.'
echo ""

# 6. 测试创建角色
echo -e "${YELLOW}6. 测试创建角色（需要 rbac:manage 权限）${NC}"
curl -s -X POST "$BASE_URL/api/v1/admin/rbac/roles" \
    -H "Authorization: Bearer $SUPER_ADMIN_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{
        "name": "test_role",
        "display_name": "测试角色",
        "description": "用于测试的角色"
    }' | jq '.'
echo ""

# 7. 测试权限拒绝（使用无权限的 Token）
echo -e "${YELLOW}7. 测试权限拒绝（模拟无 rbac:manage 权限）${NC}"
echo -e "${YELLOW}注意：这应该返回 403 Forbidden${NC}"
curl -s "$BASE_URL/api/v1/admin/rbac/roles" | jq '.'
echo ""

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}  RBAC 测试完成！${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo "超级管理员 Token（保存以便后续使用）:"
echo "$SUPER_ADMIN_TOKEN"
echo ""
echo "使用方法："
echo "  export ADMIN_TOKEN=\"$SUPER_ADMIN_TOKEN\""
echo "  curl -H \"Authorization: Bearer \$ADMIN_TOKEN\" $BASE_URL/api/v1/admin/rbac/roles"

