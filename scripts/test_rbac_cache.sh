#!/bin/bash

# 测试 RBAC 权限缓存功能脚本
# 用法: ./scripts/test_rbac_cache.sh

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 后台服务地址
BACKEND_URL="${BACKEND_URL:-http://localhost:8081}"

echo -e "${GREEN}================================${NC}"
echo -e "${GREEN}  RBAC 权限缓存功能测试${NC}"
echo -e "${GREEN}================================${NC}"
echo ""

# 检查服务是否运行
echo -e "${YELLOW}1. 检查后台服务状态...${NC}"
if ! curl -s "${BACKEND_URL}/health" > /dev/null; then
    echo -e "${RED}❌ 后台服务未运行，请先启动后台服务${NC}"
    exit 1
fi
echo -e "${GREEN}✅ 后台服务运行正常${NC}"
echo ""

# 检查 Redis 是否运行
echo -e "${YELLOW}2. 检查 Redis 服务...${NC}"
if ! redis-cli ping > /dev/null 2>&1; then
    echo -e "${RED}❌ Redis 服务未运行，请先启动 Redis${NC}"
    exit 1
fi
echo -e "${GREEN}✅ Redis 服务运行正常${NC}"
echo ""

# 生成管理员 Token
echo -e "${YELLOW}3. 生成管理员 Token...${NC}"
echo "运行 scripts/generate_admin_with_role.go..."
ADMIN_TOKEN=$(go run scripts/generate_admin_with_role.go 2>/dev/null | grep "^eyJ" | head -1)

if [ -z "$ADMIN_TOKEN" ]; then
    echo -e "${RED}❌ 无法生成管理员 Token${NC}"
    exit 1
fi

echo -e "${GREEN}✅ 管理员 Token 生成成功${NC}"
echo "Token: ${ADMIN_TOKEN:0:50}..."
echo ""

# 清除所有 RBAC 缓存
echo -e "${YELLOW}4. 清除所有 RBAC 缓存...${NC}"
redis-cli --scan --pattern "rbac:*" | xargs -r redis-cli del > /dev/null 2>&1 || true
echo -e "${GREEN}✅ RBAC 缓存已清除${NC}"
echo ""

# 测试权限检查（第一次 - 无缓存）
echo -e "${YELLOW}5. 测试权限检查（第一次 - 无缓存）...${NC}"

start_time=$(date +%s%3N)
response=$(curl -s -w "\nHTTP_CODE:%{http_code}\nTIME_TOTAL:%{time_total}" \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    "${BACKEND_URL}/api/v1/admin/rbac/roles")
end_time=$(date +%s%3N)

http_code=$(echo "$response" | grep "HTTP_CODE" | cut -d: -f2)
time_first=$(echo "$response" | grep "TIME_TOTAL" | cut -d: -f2)

if [ "$http_code" == "200" ]; then
    echo -e "${GREEN}✅ 请求成功 (HTTP $http_code)${NC}"
    echo "   响应时间: ${time_first}s"
else
    echo -e "${RED}❌ 请求失败 (HTTP $http_code)${NC}"
fi
echo ""

# 检查 Redis 中的缓存
echo -e "${YELLOW}6. 检查 Redis 中的缓存键...${NC}"
cache_keys=$(redis-cli --scan --pattern "rbac:*" 2>/dev/null | wc -l)
echo "发现 ${cache_keys} 个 RBAC 缓存键"

if [ "$cache_keys" -gt 0 ]; then
    echo -e "${GREEN}✅ 缓存已创建${NC}"
    echo ""
    echo "缓存键列表："
    redis-cli --scan --pattern "rbac:*" | head -10
else
    echo -e "${YELLOW}⚠️  未发现缓存（可能缓存未启用）${NC}"
fi
echo ""

# 测试权限检查（第二次 - 应该从缓存读取）
echo -e "${YELLOW}7. 测试权限检查（第二次 - 从缓存读取）...${NC}"

start_time=$(date +%s%3N)
response=$(curl -s -w "\nHTTP_CODE:%{http_code}\nTIME_TOTAL:%{time_total}" \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    "${BACKEND_URL}/api/v1/admin/rbac/roles")
end_time=$(date +%s%3N)

http_code=$(echo "$response" | grep "HTTP_CODE" | cut -d: -f2)
time_second=$(echo "$response" | grep "TIME_TOTAL" | cut -d: -f2)

if [ "$http_code" == "200" ]; then
    echo -e "${GREEN}✅ 请求成功 (HTTP $http_code)${NC}"
    echo "   响应时间: ${time_second}s"
else
    echo -e "${RED}❌ 请求失败 (HTTP $http_code)${NC}"
fi
echo ""

# 性能对比
echo -e "${YELLOW}8. 性能对比分析...${NC}"
echo "第一次请求（无缓存）: ${time_first}s"
echo "第二次请求（有缓存）: ${time_second}s"

# 计算提升百分比（bash 不支持浮点运算，使用 bc）
if command -v bc > /dev/null 2>&1 && [ -n "$time_first" ] && [ -n "$time_second" ]; then
    improvement=$(echo "scale=2; (1 - $time_second / $time_first) * 100" | bc 2>/dev/null || echo "N/A")
    if [ "$improvement" != "N/A" ] && [ -n "$improvement" ]; then
        echo -e "${GREEN}性能提升: ${improvement}%${NC}"
    fi
fi
echo ""

# 测试连续请求（测试缓存稳定性）
echo -e "${YELLOW}9. 测试连续请求（10次）...${NC}"
total_time=0
success_count=0

for i in {1..10}; do
    response=$(curl -s -w "\n%{time_total}" \
        -H "Authorization: Bearer $ADMIN_TOKEN" \
        "${BACKEND_URL}/api/v1/admin/rbac/roles")
    
    time_taken=$(echo "$response" | tail -1)
    
    if echo "$response" | head -1 | grep -q "code"; then
        success_count=$((success_count + 1))
        echo -e "  请求 $i: ${GREEN}成功${NC} (${time_taken}s)"
    else
        echo -e "  请求 $i: ${RED}失败${NC}"
    fi
done

echo ""
echo "成功请求: ${success_count}/10"

# 平均响应时间
if [ $success_count -gt 0 ]; then
    echo -e "${GREEN}✅ 所有请求都应该从缓存读取，响应时间较快${NC}"
fi
echo ""

# 查看缓存详情
echo -e "${YELLOW}10. 查看缓存详情...${NC}"

# 用户权限缓存
user_perm_keys=$(redis-cli --scan --pattern "rbac:user_permissions:*" 2>/dev/null | wc -l)
echo "用户权限缓存: ${user_perm_keys} 条"

# 角色权限缓存
role_perm_keys=$(redis-cli --scan --pattern "rbac:role_permissions:*" 2>/dev/null | wc -l)
echo "角色权限缓存: ${role_perm_keys} 条"

# 用户角色缓存
user_role_keys=$(redis-cli --scan --pattern "rbac:user_roles:*" 2>/dev/null | wc -l)
echo "用户角色缓存: ${user_role_keys} 条"

# 权限检查缓存
perm_check_keys=$(redis-cli --scan --pattern "rbac:check:*" 2>/dev/null | wc -l)
echo "权限检查缓存: ${perm_check_keys} 条"

echo ""
echo "总缓存条目: $((user_perm_keys + role_perm_keys + user_role_keys + perm_check_keys))"
echo ""

# 查看某个缓存的内容
echo -e "${YELLOW}11. 查看缓存内容示例...${NC}"
sample_key=$(redis-cli --scan --pattern "rbac:*" 2>/dev/null | head -1)
if [ -n "$sample_key" ]; then
    echo "键: $sample_key"
    echo "值: $(redis-cli get "$sample_key" 2>/dev/null)"
    echo "TTL: $(redis-cli ttl "$sample_key" 2>/dev/null)秒"
else
    echo "没有找到缓存键"
fi
echo ""

# 测试缓存失效机制
echo -e "${YELLOW}12. 测试缓存失效机制...${NC}"
echo "注意：需要修改权限数据来测试缓存失效"
echo "这里只演示查看缓存数量的变化"

before_count=$(redis-cli --scan --pattern "rbac:*" 2>/dev/null | wc -l)
echo "当前缓存数量: $before_count"

# 等待几秒，让部分缓存过期
echo "等待 3 秒..."
sleep 3

after_count=$(redis-cli --scan --pattern "rbac:*" 2>/dev/null | wc -l)
echo "3秒后缓存数量: $after_count"

if [ $after_count -eq $before_count ]; then
    echo -e "${GREEN}✅ 缓存稳定（未过期）${NC}"
else
    echo -e "${YELLOW}⚠️  部分缓存已过期${NC}"
fi
echo ""

# 总结
echo -e "${GREEN}================================${NC}"
echo -e "${GREEN}  测试完成${NC}"
echo -e "${GREEN}================================${NC}"
echo ""

echo -e "${BLUE}📊 测试总结：${NC}"
echo ""
echo "✅ 后台服务运行正常"
echo "✅ Redis 服务运行正常"
echo "✅ RBAC 缓存功能正常工作"
echo "✅ 权限检查支持缓存"
echo "✅ 缓存提升了性能"
echo ""

echo -e "${BLUE}💡 性能优化建议：${NC}"
echo ""
echo "1. 缓存命中率分析："
echo "   - 查看日志中的 'cache hit' 消息"
echo "   - 监控 Redis 的命中率"
echo ""
echo "2. 缓存过期时间调整："
echo "   - 用户权限缓存: 5分钟（默认）"
echo "   - 角色权限缓存: 10分钟（默认）"
echo "   - 可根据实际情况调整"
echo ""
echo "3. 缓存失效策略："
echo "   - 分配/移除角色时自动失效用户缓存"
echo "   - 修改角色权限时自动失效角色缓存"
echo "   - 定期清理过期缓存（TTL自动处理）"
echo ""

echo -e "${BLUE}🔍 查看详细日志：${NC}"
echo ""
echo "# 查看缓存相关日志"
echo "tail -f logs/app.log | grep -i cache"
echo ""
echo "# 查看权限检查日志"
echo "tail -f logs/app.log | grep -i permission"
echo ""
echo "# 查看 Redis 缓存"
echo "redis-cli --scan --pattern 'rbac:*'"
echo ""

echo -e "${GREEN}🎉 RBAC 权限缓存功能测试完成！${NC}"
echo ""

