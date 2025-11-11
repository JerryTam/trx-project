#!/bin/bash

# 数据库迁移管理脚本

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 默认配置
CONFIG_FILE="${CONFIG_FILE:-config/config.yaml}"

# 打印帮助信息
function print_help() {
    echo -e "${BLUE}数据库迁移管理工具${NC}"
    echo ""
    echo "用法: $0 <command> [options]"
    echo ""
    echo "命令:"
    echo "  up         - 执行所有待执行的迁移"
    echo "  down       - 回滚一个迁移版本"
    echo "  version    - 查看当前迁移版本"
    echo "  force      - 强制设置迁移版本 (需要 VERSION 参数)"
    echo "  goto       - 迁移到指定版本 (需要 VERSION 参数)"
    echo "  drop       - 删除所有表 (危险操作)"
    echo "  create     - 创建新的迁移文件 (需要 NAME 参数)"
    echo ""
    echo "环境变量:"
    echo "  CONFIG_FILE  - 配置文件路径 (默认: config/config.yaml)"
    echo "  VERSION      - 目标版本号 (用于 force 和 goto 命令)"
    echo "  NAME         - 迁移文件名称 (用于 create 命令)"
    echo ""
    echo "示例:"
    echo "  $0 up"
    echo "  $0 version"
    echo "  $0 force VERSION=1"
    echo "  $0 goto VERSION=3"
    echo "  $0 create NAME=add_user_phone"
    echo ""
}

# 执行迁移命令
function run_migrate() {
    local cmd=$1
    local version=$2
    
    if [ -n "$version" ]; then
        go run cmd/migrate/main.go -config "$CONFIG_FILE" -cmd "$cmd" -version "$version"
    else
        go run cmd/migrate/main.go -config "$CONFIG_FILE" -cmd "$cmd"
    fi
}

# 创建新的迁移文件
function create_migration() {
    local name=$1
    
    if [ -z "$name" ]; then
        echo -e "${RED}❌ 错误: 需要提供迁移文件名称${NC}"
        echo "用法: $0 create NAME=<migration_name>"
        exit 1
    fi
    
    # 获取当前最大版本号
    local max_version=$(ls -1 migrations/*.up.sql 2>/dev/null | sed 's/.*\/\([0-9]*\)_.*/\1/' | sort -n | tail -1)
    if [ -z "$max_version" ]; then
        max_version=0
    fi
    
    local new_version=$(printf "%06d" $((max_version + 1)))
    local up_file="migrations/${new_version}_${name}.up.sql"
    local down_file="migrations/${new_version}_${name}.down.sql"
    
    # 创建 up 文件
    cat > "$up_file" <<EOF
-- ${name} - 升级脚本
-- 创建时间: $(date '+%Y-%m-%d %H:%M:%S')

-- TODO: 在这里编写升级 SQL 语句

EOF
    
    # 创建 down 文件
    cat > "$down_file" <<EOF
-- ${name} - 回滚脚本
-- 创建时间: $(date '+%Y-%m-%d %H:%M:%S')

-- TODO: 在这里编写回滚 SQL 语句

EOF
    
    echo -e "${GREEN}✅ 迁移文件创建成功:${NC}"
    echo -e "  📄 $up_file"
    echo -e "  📄 $down_file"
}

# 主逻辑
COMMAND=${1:-help}

case $COMMAND in
    up)
        echo -e "${BLUE}🚀 执行数据库迁移...${NC}"
        run_migrate "up"
        ;;
    down)
        echo -e "${YELLOW}⚠️  回滚数据库迁移...${NC}"
        run_migrate "down"
        ;;
    version)
        echo -e "${BLUE}📌 查询当前迁移版本...${NC}"
        run_migrate "version"
        ;;
    force)
        if [ -z "$VERSION" ]; then
            echo -e "${RED}❌ 错误: 需要提供版本号${NC}"
            echo "用法: $0 force VERSION=<version>"
            exit 1
        fi
        echo -e "${YELLOW}⚠️  强制设置迁移版本...${NC}"
        run_migrate "force" "$VERSION"
        ;;
    goto)
        if [ -z "$VERSION" ]; then
            echo -e "${RED}❌ 错误: 需要提供版本号${NC}"
            echo "用法: $0 goto VERSION=<version>"
            exit 1
        fi
        echo -e "${BLUE}🎯 迁移到指定版本...${NC}"
        run_migrate "goto" "$VERSION"
        ;;
    drop)
        echo -e "${RED}⚠️  警告: 此操作将删除所有表！${NC}"
        run_migrate "drop"
        ;;
    create)
        if [ -z "$NAME" ]; then
            echo -e "${RED}❌ 错误: 需要提供迁移文件名称${NC}"
            echo "用法: $0 create NAME=<migration_name>"
            exit 1
        fi
        create_migration "$NAME"
        ;;
    help|--help|-h)
        print_help
        ;;
    *)
        echo -e "${RED}❌ 未知命令: $COMMAND${NC}"
        echo ""
        print_help
        exit 1
        ;;
esac

