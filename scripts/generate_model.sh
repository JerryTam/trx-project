#!/bin/bash

# GORM Model 自动生成脚本
# 从数据库表结构自动生成对应的 GORM Model 代码

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 默认配置
CONFIG_FILE="${CONFIG_FILE:-config/config.yaml}"
OUTPUT_DIR="${OUTPUT_DIR:-internal/model/generated}"
TABLES="${TABLES:-}"  # 空表示生成所有表

# 打印帮助信息
function print_help() {
    echo -e "${BLUE}GORM Model 自动生成工具${NC}"
    echo ""
    echo "用法: $0 [options]"
    echo ""
    echo "选项:"
    echo "  -c, --config FILE     配置文件路径 (默认: config/config.yaml)"
    echo "  -o, --output DIR      输出目录 (默认: internal/model/generated)"
    echo "  -t, --tables LIST     指定要生成的表，逗号分隔 (默认: 所有表)"
    echo "  -h, --help           显示帮助信息"
    echo ""
    echo "环境变量:"
    echo "  CONFIG_FILE          配置文件路径"
    echo "  OUTPUT_DIR           输出目录"
    echo "  TABLES               要生成的表列表（逗号分隔）"
    echo ""
    echo "示例:"
    echo "  $0                                    # 生成所有表的 Model"
    echo "  $0 -t orders,users                   # 只生成 orders 和 users 表"
    echo "  $0 -o internal/model/generated        # 指定输出目录"
    echo "  CONFIG_FILE=config/config.dev.yaml $0 # 使用不同的配置文件"
    echo ""
}

# 解析命令行参数
while [[ $# -gt 0 ]]; do
    case $1 in
        -c|--config)
            CONFIG_FILE="$2"
            shift 2
            ;;
        -o|--output)
            OUTPUT_DIR="$2"
            shift 2
            ;;
        -t|--tables)
            TABLES="$2"
            shift 2
            ;;
        -h|--help)
            print_help
            exit 0
            ;;
        *)
            echo -e "${RED}❌ 未知参数: $1${NC}"
            print_help
            exit 1
            ;;
    esac
done

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

# 读取配置文件获取数据库连接信息
if [ ! -f "$CONFIG_FILE" ]; then
    echo -e "${RED}❌ 配置文件不存在: $CONFIG_FILE${NC}"
    exit 1
fi

# 从配置文件读取数据库配置（需要 yq 或使用 Go 脚本）
# 这里简化处理，直接从环境变量或配置文件读取
# 实际项目中可以使用 go run 来读取 YAML 配置

echo -e "${BLUE}📋 读取数据库配置...${NC}"

# 使用 Go 脚本读取配置并生成 DSN
DSN=$(go run -tags=dev scripts/read_config_dsn.go "$CONFIG_FILE" 2>/dev/null)

if [ -z "$DSN" ]; then
    echo -e "${YELLOW}⚠️  无法从配置文件读取 DSN，请手动设置环境变量${NC}"
    echo "示例: export DSN=\"user:password@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local\""
    
    if [ -z "$DSN" ]; then
        echo -e "${RED}❌ 未提供数据库连接信息${NC}"
        exit 1
    fi
fi

echo -e "${GREEN}✅ 数据库配置读取成功${NC}"

# 创建输出目录
mkdir -p "$OUTPUT_DIR"
echo -e "${BLUE}📁 输出目录: $OUTPUT_DIR${NC}"

# 构建 gentool 命令
GENTOOL_CMD="gentool -db mysql -dsn \"$DSN\" -outPath \"$OUTPUT_DIR\""

# 如果指定了表，添加 -tables 参数
if [ -n "$TABLES" ]; then
    GENTOOL_CMD="$GENTOOL_CMD -tables \"$TABLES\""
    echo -e "${BLUE}🎯 生成指定表的 Model: $TABLES${NC}"
else
    echo -e "${BLUE}🎯 生成所有表的 Model${NC}"
fi

# 执行生成命令
echo -e "${YELLOW}🚀 开始生成 Model...${NC}"
eval $GENTOOL_CMD

if [ $? -eq 0 ]; then
    echo ""
    echo -e "${GREEN}✅ Model 生成成功！${NC}"
    echo -e "${GREEN}📄 生成的文件位于: $OUTPUT_DIR${NC}"
    echo ""
    echo -e "${YELLOW}💡 提示:${NC}"
    echo "  1. 生成的 Model 在 $OUTPUT_DIR 目录"
    echo "  2. 可以手动复制需要的 Model 到 internal/model/ 目录"
    echo "  3. 根据业务需求添加自定义方法（如 GORM 钩子、状态枚举等）"
else
    echo -e "${RED}❌ Model 生成失败${NC}"
    exit 1
fi

