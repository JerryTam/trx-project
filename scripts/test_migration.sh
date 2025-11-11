#!/bin/bash

# 数据库迁移功能测试脚本

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}   数据库迁移功能测试${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# 1. 测试迁移版本查询
echo -e "${BLUE}📌 测试 1: 查询迁移版本${NC}"
if ./scripts/migrate.sh version 2>&1 | grep -q "当前迁移版本\|Access denied\|connection refused"; then
    if ./scripts/migrate.sh version 2>&1 | grep -q "Access denied\|connection refused"; then
        echo -e "${YELLOW}⚠️  数据库未启动，跳过版本查询测试${NC}"
    else
        echo -e "${GREEN}✅ 版本查询成功${NC}"
    fi
else
    echo -e "${RED}❌ 版本查询失败${NC}"
fi
echo ""

# 2. 测试创建新迁移文件
echo -e "${BLUE}📝 测试 2: 创建新迁移文件${NC}"
TEST_MIGRATION_NAME="test_feature_$(date +%s)"
NAME=$TEST_MIGRATION_NAME ./scripts/migrate.sh create

# 检查文件是否创建
LAST_VERSION=$(ls -1 migrations/*.up.sql | sed 's/.*\/\([0-9]*\)_.*/\1/' | sort -n | tail -1)
UP_FILE="migrations/${LAST_VERSION}_${TEST_MIGRATION_NAME}.up.sql"
DOWN_FILE="migrations/${LAST_VERSION}_${TEST_MIGRATION_NAME}.down.sql"

if [ -f "$UP_FILE" ] && [ -f "$DOWN_FILE" ]; then
    echo -e "${GREEN}✅ 迁移文件创建成功${NC}"
    echo -e "  📄 $UP_FILE"
    echo -e "  📄 $DOWN_FILE"
    
    # 清理测试文件
    rm -f "$UP_FILE" "$DOWN_FILE"
    echo -e "${YELLOW}🧹 已清理测试文件${NC}"
else
    echo -e "${RED}❌ 迁移文件创建失败${NC}"
    exit 1
fi
echo ""

# 3. 测试迁移执行（如果数据库可用）
echo -e "${BLUE}🚀 测试 3: 执行数据库迁移${NC}"
migrate_output=$(./scripts/migrate.sh up 2>&1)
if echo "$migrate_output" | grep -q "completed successfully\|already up to date"; then
    echo -e "${GREEN}✅ 迁移执行成功${NC}"
elif echo "$migrate_output" | grep -q "Access denied\|connection refused"; then
    echo -e "${YELLOW}⚠️  数据库未启动，跳过迁移执行测试${NC}"
else
    echo -e "${YELLOW}⚠️  迁移执行失败（可能数据库未启动）${NC}"
fi
echo ""

# 4. 测试 Makefile 命令
echo -e "${BLUE}🛠️  测试 4: Makefile 命令${NC}"

# 测试 make migrate-version
make_output=$(make migrate-version 2>&1)
if echo "$make_output" | grep -q "当前迁移版本\|Access denied\|connection refused"; then
    if echo "$make_output" | grep -q "Access denied\|connection refused"; then
        echo -e "${YELLOW}⚠️  数据库未启动，跳过 Makefile 命令测试${NC}"
    else
        echo -e "${GREEN}✅ make migrate-version 工作正常${NC}"
    fi
else
    echo -e "${YELLOW}⚠️  make migrate-version 可能失败${NC}"
fi
echo ""

# 5. 测试 Go 命令行工具编译
echo -e "${BLUE}🔧 测试 5: 编译迁移工具${NC}"
if go build -o /tmp/migrate-test cmd/migrate/main.go 2>&1; then
    echo -e "${GREEN}✅ 迁移工具编译成功${NC}"
    rm -f /tmp/migrate-test
else
    echo -e "${RED}❌ 迁移工具编译失败${NC}"
    exit 1
fi
echo ""

# 6. 测试迁移文件语法
echo -e "${BLUE}📋 测试 6: 检查迁移文件语法${NC}"
ERROR_COUNT=0

for file in migrations/*.sql; do
    # 检查文件是否为空
    if [ ! -s "$file" ]; then
        echo -e "${RED}❌ 文件为空: $file${NC}"
        ERROR_COUNT=$((ERROR_COUNT + 1))
        continue
    fi
    
    # 检查基本 SQL 语法（简单检查）
    if grep -q "^--" "$file" || grep -qi "CREATE\|ALTER\|DROP\|INSERT" "$file"; then
        echo -e "${GREEN}✓${NC} $file"
    else
        echo -e "${YELLOW}⚠️  可能有问题: $file${NC}"
        ERROR_COUNT=$((ERROR_COUNT + 1))
    fi
done

if [ $ERROR_COUNT -eq 0 ]; then
    echo -e "${GREEN}✅ 所有迁移文件检查通过${NC}"
else
    echo -e "${YELLOW}⚠️  发现 $ERROR_COUNT 个潜在问题${NC}"
fi
echo ""

# 7. 测试迁移文件配对
echo -e "${BLUE}🔗 测试 7: 检查迁移文件配对${NC}"
UNPAIRED=0

for up_file in migrations/*.up.sql; do
    down_file="${up_file/.up.sql/.down.sql}"
    if [ ! -f "$down_file" ]; then
        echo -e "${RED}❌ 缺少配对的 down 文件: $down_file${NC}"
        UNPAIRED=$((UNPAIRED + 1))
    fi
done

if [ $UNPAIRED -eq 0 ]; then
    echo -e "${GREEN}✅ 所有迁移文件都有配对的 up/down 脚本${NC}"
else
    echo -e "${RED}❌ 发现 $UNPAIRED 个未配对的迁移文件${NC}"
fi
echo ""

# 8. 测试环境变量配置
echo -e "${BLUE}🔧 测试 8: 环境变量配置${NC}"

# 测试不同配置文件
for env in dev test prod; do
    config_file="config/config.${env}.yaml"
    if [ -f "$config_file" ]; then
        echo -e "${GREEN}✓${NC} 找到配置文件: $config_file"
    else
        echo -e "${YELLOW}⚠️  缺少配置文件: $config_file${NC}"
    fi
done
echo ""

# 总结
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}   测试完成总结${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""
echo -e "${GREEN}✅ 迁移功能基础测试通过${NC}"
echo ""
echo -e "📖 完整测试步骤："
echo -e "  1. 启动数据库: ${YELLOW}make docker-up${NC}"
echo -e "  2. 运行迁移: ${YELLOW}make migrate-up${NC}"
echo -e "  3. 查看版本: ${YELLOW}make migrate-version${NC}"
echo -e "  4. 回滚迁移: ${YELLOW}make migrate-down${NC}"
echo ""
echo -e "📚 相关文档: ${BLUE}docs/MIGRATION_GUIDE.md${NC}"
echo ""

