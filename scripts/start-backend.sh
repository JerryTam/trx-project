#!/bin/bash
# TRX Project - Backend 启动脚本
# 跨平台支持: Linux/macOS/Windows(Git Bash)

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# 设置环境变量
export GO_ENV=dev
export AUTO_MIGRATE=false

# 清理旧的编译文件
if [ -f "tmp/backend" ] || [ -f "tmp/backend.exe" ]; then
    echo -e "${YELLOW}清理旧的编译文件...${NC}"
    rm -f tmp/backend tmp/backend.exe
fi

# 显示启动信息
echo ""
echo -e "${CYAN}========================================${NC}"
echo -e "${CYAN}  🚀 启动 Backend 服务 (热重载模式)${NC}"
echo -e "${CYAN}========================================${NC}"
echo ""
echo -e "${GREEN}环境变量:${NC}"
echo -e "  GO_ENV       = ${GO_ENV}"
echo -e "  AUTO_MIGRATE = ${AUTO_MIGRATE}"
echo ""
echo -e "${GREEN}服务地址:${NC}"
echo -e "  HTTP:    http://localhost:8081"
echo -e "  Health:  http://localhost:8081/health"
echo -e "  Swagger: http://localhost:8081/swagger/index.html"
echo ""
echo -e "${YELLOW}提示: 按 Ctrl+C 停止服务${NC}"
echo ""

# 检测操作系统并选择配置文件
detect_os() {
    case "$(uname -s)" in
        Linux*)     echo "linux";;
        Darwin*)    echo "linux";;  # macOS 使用 Linux 配置
        CYGWIN*|MINGW*|MSYS*)    echo "windows";;
        *)          echo "linux";;  # 默认使用 Linux 配置
    esac
}

OS=$(detect_os)
if [ "$OS" = "windows" ]; then
    AIR_CONFIG=".air-backend.toml"
else
    AIR_CONFIG=".air-backend-linux.toml"
fi

echo -e "${GREEN}使用配置文件: ${AIR_CONFIG}${NC}"
echo ""

# 启动 Air
air -c $AIR_CONFIG

