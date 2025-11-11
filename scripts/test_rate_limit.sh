#!/bin/bash

# 测试限流功能脚本
# 用法: ./scripts/test_rate_limit.sh

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 前台服务地址
FRONTEND_URL="${FRONTEND_URL:-http://localhost:8080}"
# 后台服务地址
BACKEND_URL="${BACKEND_URL:-http://localhost:8081}"

echo -e "${GREEN}================================${NC}"
echo -e "${GREEN}  限流功能测试${NC}"
echo -e "${GREEN}================================${NC}"
echo ""

# 检查服务是否运行
echo -e "${YELLOW}1. 检查服务状态...${NC}"
if ! curl -s "${FRONTEND_URL}/health" > /dev/null; then
    echo -e "${RED}❌ 前台服务未运行，请先启动前台服务${NC}"
    exit 1
fi

if ! curl -s "${BACKEND_URL}/health" > /dev/null; then
    echo -e "${RED}❌ 后台服务未运行，请先启动后台服务${NC}"
    exit 1
fi
echo -e "${GREEN}✅ 服务运行正常${NC}"
echo ""

# 测试 IP 限流
echo -e "${YELLOW}2. 测试 IP 限流（快速连续请求）...${NC}"
echo "发送 10 个连续请求到前台服务..."
success_count=0
rate_limited_count=0

for i in {1..10}; do
    response=$(curl -s -w "\n%{http_code}" "${FRONTEND_URL}/api/v1/public/login" \
        -H "Content-Type: application/json" \
        -d '{"username":"test","password":"test"}')
    
    status_code=$(echo "$response" | tail -n1)
    
    if [ "$status_code" == "200" ] || [ "$status_code" == "400" ] || [ "$status_code" == "401" ]; then
        success_count=$((success_count + 1))
        echo -e "  请求 $i: ${GREEN}成功 (${status_code})${NC}"
    elif [ "$status_code" == "429" ]; then
        rate_limited_count=$((rate_limited_count + 1))
        echo -e "  请求 $i: ${RED}被限流 (429)${NC}"
        # 提取限流响应信息
        body=$(echo "$response" | sed '$d')
        echo "    响应: $body"
    else
        echo -e "  请求 $i: ${YELLOW}其他状态 (${status_code})${NC}"
    fi
    
    # 短暂延迟
    sleep 0.1
done

echo ""
echo "结果统计:"
echo "  成功请求: ${success_count}"
echo "  被限流: ${rate_limited_count}"
echo ""

# 测试响应头
echo -e "${YELLOW}3. 检查限流响应头...${NC}"
response=$(curl -s -i "${FRONTEND_URL}/api/v1/public/login" \
    -H "Content-Type: application/json" \
    -d '{"username":"test","password":"test"}' 2>&1)

if echo "$response" | grep -q "X-RateLimit-Limit"; then
    limit=$(echo "$response" | grep "X-RateLimit-Limit" | cut -d' ' -f2 | tr -d '\r')
    remaining=$(echo "$response" | grep "X-RateLimit-Remaining" | cut -d' ' -f2 | tr -d '\r')
    reset=$(echo "$response" | grep "X-RateLimit-Reset" | cut -d' ' -f2 | tr -d '\r')
    
    echo -e "${GREEN}✅ 找到限流响应头:${NC}"
    echo "  X-RateLimit-Limit: $limit"
    echo "  X-RateLimit-Remaining: $remaining"
    echo "  X-RateLimit-Reset: $reset"
else
    echo -e "${YELLOW}⚠️  未找到限流响应头（可能限流未启用）${NC}"
fi
echo ""

# 测试用户级别限流（需要先注册和登录）
echo -e "${YELLOW}4. 测试用户级别限流...${NC}"

# 生成随机用户名
RANDOM_USER="testuser_$(date +%s)"

# 注册用户
echo "注册测试用户: $RANDOM_USER"
register_response=$(curl -s "${FRONTEND_URL}/api/v1/public/register" \
    -H "Content-Type: application/json" \
    -d "{\"username\":\"${RANDOM_USER}\",\"email\":\"${RANDOM_USER}@test.com\",\"password\":\"Test123456\"}")

# 提取 token
token=$(echo "$register_response" | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$token" ]; then
    echo -e "${YELLOW}⚠️  无法获取 token，跳过用户级别限流测试${NC}"
else
    echo -e "${GREEN}✅ 获取到 token${NC}"
    echo ""
    
    echo "使用 token 发送多个请求..."
    user_success=0
    user_limited=0
    
    for i in {1..15}; do
        status_code=$(curl -s -o /dev/null -w "%{http_code}" \
            "${FRONTEND_URL}/api/v1/user/profile" \
            -H "Authorization: Bearer $token")
        
        if [ "$status_code" == "200" ]; then
            user_success=$((user_success + 1))
            echo -e "  用户请求 $i: ${GREEN}成功 (200)${NC}"
        elif [ "$status_code" == "429" ]; then
            user_limited=$((user_limited + 1))
            echo -e "  用户请求 $i: ${RED}被限流 (429)${NC}"
        else
            echo -e "  用户请求 $i: 状态 ${status_code}"
        fi
        
        sleep 0.1
    done
    
    echo ""
    echo "用户级别限流测试结果:"
    echo "  成功请求: ${user_success}"
    echo "  被限流: ${user_limited}"
fi
echo ""

# 测试全局限流（需要大量并发请求，这里只做基本测试）
echo -e "${YELLOW}5. 测试全局限流（并发请求）...${NC}"
echo "发送 20 个并发请求..."

# 使用后台任务并发请求
pids=()
temp_dir=$(mktemp -d)

for i in {1..20}; do
    (
        status_code=$(curl -s -o /dev/null -w "%{http_code}" "${FRONTEND_URL}/health")
        echo "$status_code" > "${temp_dir}/result_${i}.txt"
    ) &
    pids+=($!)
done

# 等待所有请求完成
for pid in "${pids[@]}"; do
    wait $pid
done

# 统计结果
global_success=0
global_limited=0

for result_file in "${temp_dir}"/result_*.txt; do
    status_code=$(cat "$result_file")
    if [ "$status_code" == "200" ]; then
        global_success=$((global_success + 1))
    elif [ "$status_code" == "429" ]; then
        global_limited=$((global_limited + 1))
    fi
done

rm -rf "${temp_dir}"

echo ""
echo "全局限流测试结果:"
echo "  成功请求: ${global_success}"
echo "  被限流: ${global_limited}"
echo ""

# 总结
echo -e "${GREEN}================================${NC}"
echo -e "${GREEN}  测试完成${NC}"
echo -e "${GREEN}================================${NC}"
echo ""
echo "💡 提示："
echo "1. 如果所有请求都成功，可能是："
echo "   - 限流未启用（检查配置文件 rate_limit.enabled）"
echo "   - 限流阈值设置过高"
echo "   - Redis 未正确连接"
echo ""
echo "2. 如果部分请求被限流，说明限流功能正常工作"
echo ""
echo "3. 可以通过修改配置文件调整限流参数："
echo "   - config/config.yaml 中的 rate_limit 部分"
echo "   - global_rate: 全局限流"
echo "   - ip_rate: IP 限流"
echo "   - user_rate: 用户限流"
echo ""
echo "4. 查看 Redis 中的限流数据："
echo "   redis-cli keys 'rate_limit:*'"
echo ""

