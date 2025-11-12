#!/bin/bash

# GORM Model 自动生成脚本（简化版）
# 需要手动提供数据库连接信息

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# 默认配置
DSN="${DSN:-}"
OUTPUT_DIR="${OUTPUT_DIR:-internal/model/generated}"
TABLES="${TABLES:-}"

# 打印帮助信息
function print_help() {
    echo -e "${BLUE}GORM Model 自动生成工具（简化版）${NC}"
    echo ""
    echo "用法: DSN=\"<dsn>\" $0 [options]"
    echo ""
    echo "环境变量:"
    echo "  DSN           数据库连接字符串 (必需)"
    echo "                格式: user:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
    echo "  OUTPUT_DIR    输出目录 (默认: internal/model/generated)"
    echo "  TABLES        要生成的表，逗号分隔 (默认: 所有表)"
    echo ""
    echo "示例:"
    echo "  DSN=\"root:password@tcp(localhost:3306)/trx_db?charset=utf8mb4&parseTime=True&loc=Local\" $0"
    echo "  DSN=\"...\" TABLES=\"orders,users\" $0"
    echo "  DSN=\"...\" OUTPUT_DIR=\"internal/model/generated\" $0"
    echo ""
}

# 检查参数
if [ "$1" == "-h" ] || [ "$1" == "--help" ]; then
    print_help
    exit 0
fi

# 检查 DSN
if [ -z "$DSN" ]; then
    echo -e "${RED}❌ 错误: 需要提供 DSN 环境变量${NC}"
    echo ""
    print_help
    exit 1
fi

# 检查 gentool 是否安装
if ! command -v gentool &> /dev/null; then
    echo -e "${YELLOW}⚠️  gentool 未安装，正在安装...${NC}"
    go install gorm.io/gen/tools/gentool@latest
    if [ $? -ne 0 ]; then
        echo -e "${RED}❌ 安装 gentool 失败${NC}"
        exit 1
    fi
    echo -e "${GREEN}✅ gentool 安装成功${NC}"
fi

# 创建输出目录
mkdir -p "$OUTPUT_DIR"
echo -e "${BLUE}📁 输出目录: $OUTPUT_DIR${NC}"

# 构建命令
if [ -n "$TABLES" ]; then
    echo -e "${BLUE}🎯 生成指定表的 Model: $TABLES${NC}"
    gentool -db mysql -dsn "$DSN" -tables "$TABLES" -outPath "$OUTPUT_DIR"
else
    echo -e "${BLUE}🎯 生成所有表的 Model${NC}"
    gentool -db mysql -dsn "$DSN" -outPath "$OUTPUT_DIR"
fi

if [ $? -eq 0 ]; then
    echo ""
    echo -e "${GREEN}✅ Model 生成成功！${NC}"
    echo -e "${GREEN}📄 生成的文件位于: $OUTPUT_DIR${NC}"
    echo ""
    echo -e "${YELLOW}💡 下一步:${NC}"
    echo "  1. 查看生成的文件: ls -la $OUTPUT_DIR"
    echo "  2. 复制需要的 Model 到 internal/model/ 目录"
    echo "  3. 根据业务需求添加自定义方法"
else
    echo -e "${RED}❌ Model 生成失败${NC}"
    exit 1
fi

