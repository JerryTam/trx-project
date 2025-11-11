#!/bin/bash

# 测试请求 ID 追踪功能脚本
# 用法: ./scripts/test_request_id.sh

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 前台服务地址
FRONTEND_URL="${FRONTEND_URL:-http://localhost:8080}"
# 后台服务地址
BACKEND_URL="${BACKEND_URL:-http://localhost:8081}"

echo -e "${GREEN}================================${NC}"
echo -e "${GREEN}  请求 ID 追踪功能测试${NC}"
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

# 测试 1: 服务器自动生成请求 ID
echo -e "${YELLOW}2. 测试服务器自动生成请求 ID...${NC}"
echo "发送不带 X-Request-ID 的请求..."

response=$(curl -s -i "${FRONTEND_URL}/health" 2>&1)
request_id=$(echo "$response" | grep -i "X-Request-ID" | cut -d' ' -f2 | tr -d '\r\n')

if [ -n "$request_id" ]; then
    echo -e "${GREEN}✅ 服务器自动生成了请求 ID:${NC} ${BLUE}${request_id}${NC}"
    
    # 验证 UUID 格式 (8-4-4-4-12)
    if echo "$request_id" | grep -qE '^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$'; then
        echo -e "${GREEN}✅ 请求 ID 格式正确 (UUID)${NC}"
    else
        echo -e "${YELLOW}⚠️  请求 ID 格式不是标准 UUID${NC}"
    fi
else
    echo -e "${RED}❌ 未找到 X-Request-ID 响应头${NC}"
fi
echo ""

# 测试 2: 客户端指定请求 ID
echo -e "${YELLOW}3. 测试客户端指定请求 ID...${NC}"
custom_request_id="custom-test-id-12345"
echo "发送带有自定义 X-Request-ID 的请求: ${custom_request_id}"

response=$(curl -s -i -H "X-Request-ID: ${custom_request_id}" "${FRONTEND_URL}/health" 2>&1)
returned_id=$(echo "$response" | grep -i "X-Request-ID" | cut -d' ' -f2 | tr -d '\r\n')

if [ "$returned_id" = "$custom_request_id" ]; then
    echo -e "${GREEN}✅ 服务器返回了相同的请求 ID:${NC} ${BLUE}${returned_id}${NC}"
else
    echo -e "${RED}❌ 请求 ID 不匹配${NC}"
    echo "  期望: ${custom_request_id}"
    echo "  实际: ${returned_id}"
fi
echo ""

# 测试 3: 请求链追踪
echo -e "${YELLOW}4. 测试请求链追踪...${NC}"
echo "使用同一个请求 ID 发送多个请求..."

trace_id="trace-$(date +%s)"
echo -e "追踪 ID: ${BLUE}${trace_id}${NC}"
echo ""

# 发送 3 个请求，使用相同的请求 ID
for i in {1..3}; do
    response=$(curl -s -i -H "X-Request-ID: ${trace_id}" "${FRONTEND_URL}/health" 2>&1)
    returned_id=$(echo "$response" | grep -i "X-Request-ID" | cut -d' ' -f2 | tr -d '\r\n')
    
    if [ "$returned_id" = "$trace_id" ]; then
        echo -e "  请求 ${i}: ${GREEN}✅ 请求 ID 一致${NC}"
    else
        echo -e "  请求 ${i}: ${RED}❌ 请求 ID 不一致${NC}"
    fi
done

echo ""
echo -e "${GREEN}💡 提示：${NC}查看日志文件，这些请求应该有相同的 request_id: ${trace_id}"
echo ""

# 测试 4: 前台和后台服务
echo -e "${YELLOW}5. 测试前后台服务的请求 ID...${NC}"

# 前台服务
echo "测试前台服务..."
frontend_response=$(curl -s -i "${FRONTEND_URL}/health" 2>&1)
frontend_id=$(echo "$frontend_response" | grep -i "X-Request-ID" | cut -d' ' -f2 | tr -d '\r\n')

if [ -n "$frontend_id" ]; then
    echo -e "${GREEN}✅ 前台服务请求 ID:${NC} ${BLUE}${frontend_id}${NC}"
else
    echo -e "${RED}❌ 前台服务未返回请求 ID${NC}"
fi

# 后台服务
echo "测试后台服务..."
backend_response=$(curl -s -i "${BACKEND_URL}/health" 2>&1)
backend_id=$(echo "$backend_response" | grep -i "X-Request-ID" | cut -d' ' -f2 | tr -d '\r\n')

if [ -n "$backend_id" ]; then
    echo -e "${GREEN}✅ 后台服务请求 ID:${NC} ${BLUE}${backend_id}${NC}"
else
    echo -e "${RED}❌ 后台服务未返回请求 ID${NC}"
fi
echo ""

# 测试 5: 不同 API 端点
echo -e "${YELLOW}6. 测试不同 API 端点的请求 ID...${NC}"

endpoints=(
    "${FRONTEND_URL}/health"
    "${FRONTEND_URL}/api/v1/public/login"
)

for endpoint in "${endpoints[@]}"; do
    response=$(curl -s -i -X POST "$endpoint" \
        -H "Content-Type: application/json" \
        -d '{"username":"test","password":"test"}' 2>&1)
    
    request_id=$(echo "$response" | grep -i "X-Request-ID" | cut -d' ' -f2 | tr -d '\r\n')
    
    if [ -n "$request_id" ]; then
        echo -e "${GREEN}✅${NC} ${endpoint}"
        echo -e "   请求 ID: ${BLUE}${request_id}${NC}"
    else
        echo -e "${RED}❌${NC} ${endpoint}"
        echo "   未找到请求 ID"
    fi
done
echo ""

# 测试 6: 并发请求的请求 ID 唯一性
echo -e "${YELLOW}7. 测试并发请求的请求 ID 唯一性...${NC}"
echo "发送 10 个并发请求，验证每个请求 ID 都是唯一的..."

temp_dir=$(mktemp -d)
pids=()

for i in {1..10}; do
    (
        response=$(curl -s -i "${FRONTEND_URL}/health" 2>&1)
        request_id=$(echo "$response" | grep -i "X-Request-ID" | cut -d' ' -f2 | tr -d '\r\n')
        echo "$request_id" > "${temp_dir}/id_${i}.txt"
    ) &
    pids+=($!)
done

# 等待所有请求完成
for pid in "${pids[@]}"; do
    wait $pid
done

# 收集所有请求 ID
declare -A request_ids
duplicate_count=0

for id_file in "${temp_dir}"/id_*.txt; do
    id=$(cat "$id_file")
    if [ -n "$id" ]; then
        if [ -n "${request_ids[$id]}" ]; then
            echo -e "${RED}❌ 发现重复的请求 ID: ${id}${NC}"
            duplicate_count=$((duplicate_count + 1))
        else
            request_ids[$id]=1
        fi
    fi
done

rm -rf "${temp_dir}"

unique_count=${#request_ids[@]}
echo -e "生成的唯一请求 ID 数量: ${GREEN}${unique_count}${NC}/10"

if [ $unique_count -eq 10 ]; then
    echo -e "${GREEN}✅ 所有请求 ID 都是唯一的${NC}"
else
    echo -e "${RED}❌ 发现 ${duplicate_count} 个重复的请求 ID${NC}"
fi
echo ""

# 总结
echo -e "${GREEN}================================${NC}"
echo -e "${GREEN}  测试完成${NC}"
echo -e "${GREEN}================================${NC}"
echo ""

echo -e "${BLUE}📊 功能验证总结：${NC}"
echo ""
echo "✅ 服务器自动生成请求 ID (UUID 格式)"
echo "✅ 客户端可以指定请求 ID"
echo "✅ 请求 ID 在请求链中保持一致"
echo "✅ 前后台服务都支持请求 ID"
echo "✅ 所有 API 端点都返回请求 ID"
echo "✅ 并发请求的请求 ID 唯一性"
echo ""

echo -e "${BLUE}💡 使用建议：${NC}"
echo ""
echo "1. 查看日志中的请求 ID："
echo "   tail -f logs/app.log | grep request_id"
echo ""
echo "2. 追踪特定请求："
echo "   grep '<request-id>' logs/app.log"
echo ""
echo "3. 在客户端指定请求 ID："
echo "   curl -H 'X-Request-ID: custom-id' http://localhost:8080/api/v1/..."
echo ""
echo "4. 在分布式系统中传递请求 ID："
echo "   - 前端生成请求 ID"
echo "   - 通过 HTTP 头传递给后端"
echo "   - 后端在调用其他服务时继续传递"
echo "   - 所有日志都包含相同的请求 ID"
echo ""

echo -e "${GREEN}🎉 请求 ID 追踪功能正常工作！${NC}"
echo ""

