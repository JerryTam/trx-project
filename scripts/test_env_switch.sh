#!/bin/bash

# 环境切换测试脚本

YELLOW='\033[1;33m'
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${YELLOW}========================================${NC}"
echo -e "${YELLOW}  环境配置切换测试${NC}"
echo -e "${YELLOW}========================================${NC}"
echo ""

# 测试开发环境
echo -e "${YELLOW}1. 测试开发环境 (dev)${NC}"
GO_ENV=dev go run cmd/frontend/main.go > /tmp/frontend_dev.log 2>&1 &
PID_DEV=$!
sleep 2
if grep -q "config.dev.yaml" /tmp/frontend_dev.log; then
    echo -e "${GREEN}✅ 开发环境配置加载成功${NC}"
    cat /tmp/frontend_dev.log | grep "加载环境配置"
else
    echo -e "${RED}❌ 开发环境配置加载失败${NC}"
fi
kill $PID_DEV 2>/dev/null
echo ""

# 测试测试环境
echo -e "${YELLOW}2. 测试测试环境 (test)${NC}"
GO_ENV=test go run cmd/frontend/main.go > /tmp/frontend_test.log 2>&1 &
PID_TEST=$!
sleep 2
if grep -q "config.test.yaml" /tmp/frontend_test.log; then
    echo -e "${GREEN}✅ 测试环境配置加载成功${NC}"
    cat /tmp/frontend_test.log | grep "加载环境配置"
else
    echo -e "${RED}❌ 测试环境配置加载失败${NC}"
fi
kill $PID_TEST 2>/dev/null
echo ""

# 测试生产环境
echo -e "${YELLOW}3. 测试生产环境 (prod)${NC}"
GO_ENV=prod go run cmd/frontend/main.go > /tmp/frontend_prod.log 2>&1 &
PID_PROD=$!
sleep 2
if grep -q "config.prod.yaml" /tmp/frontend_prod.log; then
    echo -e "${GREEN}✅ 生产环境配置加载成功${NC}"
    cat /tmp/frontend_prod.log | grep "加载环境配置"
else
    echo -e "${RED}❌ 生产环境配置加载失败${NC}"
fi
kill $PID_PROD 2>/dev/null
echo ""

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}  环境配置切换测试完成！${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo "配置文件位置："
echo "  - config/config.dev.yaml  (开发环境)"
echo "  - config/config.test.yaml (测试环境)"
echo "  - config/config.prod.yaml (生产环境)"
echo ""
echo "使用方法："
echo "  GO_ENV=dev ./bin/frontend"
echo "  GO_ENV=test ./bin/frontend"
echo "  GO_ENV=prod ./bin/frontend"

