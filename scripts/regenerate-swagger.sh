#!/bin/bash
# TRX Project - Swagger 文档重新生成脚本

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

echo ""
echo -e "${CYAN}========================================${NC}"
echo -e "${CYAN}  🔄 重新生成 Swagger 文档${NC}"
echo -e "${CYAN}========================================${NC}"
echo ""

# 检查 swag 是否安装
if ! command -v swag &> /dev/null; then
    echo -e "${RED}❌ 错误: swag 未安装${NC}"
    echo ""
    echo -e "${YELLOW}请运行以下命令安装:${NC}"
    echo "  go install github.com/swaggo/swag/cmd/swag@latest"
    echo ""
    exit 1
fi

echo -e "${GREEN}✅ swag 已安装${NC}"
swag --version
echo ""

# 生成前端文档
echo -e "${YELLOW}📝 生成前端 Swagger 文档...${NC}"
if swag init -g cmd/frontend/main.go -o cmd/frontend/docs --parseDependency --parseInternal; then
    echo -e "${GREEN}✅ 前端文档生成成功${NC}"
else
    echo -e "${RED}❌ 前端文档生成失败${NC}"
    exit 1
fi
echo ""

# 生成后端文档
echo -e "${YELLOW}📝 生成后端 Swagger 文档...${NC}"
if swag init -g cmd/backend/main.go -o cmd/backend/docs --parseDependency --parseInternal; then
    echo -e "${GREEN}✅ 后端文档生成成功${NC}"
else
    echo -e "${RED}❌ 后端文档生成失败${NC}"
    exit 1
fi
echo ""

# 验证生成的文件
echo -e "${YELLOW}🔍 验证生成的文件...${NC}"

FRONTEND_JSON="cmd/frontend/docs/swagger.json"
BACKEND_JSON="cmd/backend/docs/swagger.json"

if [ -f "$FRONTEND_JSON" ]; then
    FRONTEND_SIZE=$(wc -c < "$FRONTEND_JSON")
    if [ $FRONTEND_SIZE -gt 1000 ]; then
        echo -e "${GREEN}✅ 前端 swagger.json: $FRONTEND_SIZE bytes${NC}"
    else
        echo -e "${YELLOW}⚠️  前端 swagger.json 文件过小，可能生成不完整${NC}"
    fi
else
    echo -e "${RED}❌ 前端 swagger.json 文件不存在${NC}"
fi

if [ -f "$BACKEND_JSON" ]; then
    BACKEND_SIZE=$(wc -c < "$BACKEND_JSON")
    if [ $BACKEND_SIZE -gt 1000 ]; then
        echo -e "${GREEN}✅ 后端 swagger.json: $BACKEND_SIZE bytes${NC}"
    else
        echo -e "${YELLOW}⚠️  后端 swagger.json 文件过小，可能生成不完整${NC}"
    fi
else
    echo -e "${RED}❌ 后端 swagger.json 文件不存在${NC}"
fi

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}  ✅ Swagger 文档生成完成！${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo -e "${CYAN}📚 访问地址:${NC}"
echo "  前端: http://localhost:8080/swagger/index.html"
echo "  后端: http://localhost:8081/swagger/index.html"
echo ""
echo -e "${YELLOW}💡 提示:${NC}"
echo "  1. 如果服务正在运行，请重启以加载新文档"
echo "  2. 使用 Air 热重载会自动重新编译"
echo "  3. 清理缓存: rm -rf tmp && mkdir tmp"
echo ""

