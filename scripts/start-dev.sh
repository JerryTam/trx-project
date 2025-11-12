#!/bin/bash
# TRX Project - 同时启动前后端服务
# 跨平台支持: Linux/macOS/Windows(Git Bash)

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

echo ""
echo -e "${CYAN}========================================${NC}"
echo -e "${CYAN}  🚀 TRX Project 开发环境启动器${NC}"
echo -e "${CYAN}========================================${NC}"
echo ""

# 检查 Air 是否安装
if ! command -v air &> /dev/null; then
    echo -e "${RED}❌ 错误: Air 未安装${NC}"
    echo ""
    echo -e "${YELLOW}请运行以下命令安装:${NC}"
    echo "  go install github.com/air-verse/air@latest"
    echo ""
    exit 1
fi

echo -e "${GREEN}✅ Air 已安装${NC}"
echo ""

# 检查端口是否被占用
check_port() {
    local port=$1
    if command -v lsof &> /dev/null; then
        # macOS/Linux with lsof
        lsof -i :$port &> /dev/null
    elif command -v netstat &> /dev/null; then
        # Linux/Windows with netstat
        netstat -an | grep ":$port " | grep LISTEN &> /dev/null
    else
        # 无法检查
        return 1
    fi
}

PORT_WARNING=false

if check_port 8080; then
    echo -e "${YELLOW}⚠️  警告: 端口 8080 已被占用${NC}"
    echo -e "${YELLOW}请先停止占用该端口的进程${NC}"
    echo ""
    PORT_WARNING=true
fi

if check_port 8081; then
    echo -e "${YELLOW}⚠️  警告: 端口 8081 已被占用${NC}"
    echo -e "${YELLOW}请先停止占用该端口的进程${NC}"
    echo ""
    PORT_WARNING=true
fi

if [ "$PORT_WARNING" = true ]; then
    read -p "是否继续? (y/n): " continue
    if [ "$continue" != "y" ]; then
        exit 0
    fi
fi

# 检测操作系统
detect_os() {
    case "$(uname -s)" in
        Linux*)     echo "linux";;
        Darwin*)    echo "macos";;
        CYGWIN*|MINGW*|MSYS*)    echo "windows";;
        *)          echo "unknown";;
    esac
}

OS=$(detect_os)
echo -e "${GREEN}检测到操作系统: ${OS}${NC}"

# 选择配置文件
if [ "$OS" = "windows" ]; then
    FRONTEND_CONFIG=".air-frontend.toml"
    BACKEND_CONFIG=".air-backend.toml"
else
    FRONTEND_CONFIG=".air-frontend-linux.toml"
    BACKEND_CONFIG=".air-backend-linux.toml"
fi

echo -e "${GREEN}Frontend 配置: ${FRONTEND_CONFIG}${NC}"
echo -e "${GREEN}Backend 配置:  ${BACKEND_CONFIG}${NC}"
echo ""

# 设置环境变量
export GO_ENV=dev
export AUTO_MIGRATE=false

echo -e "${CYAN}启动方式:${NC}"
echo "1. 使用 tmux 在分屏中启动（推荐，Linux/macOS）"
echo "2. 在后台启动两个服务"
echo "3. 依次启动（前台模式，仅启动 Frontend）"
echo ""
read -p "请选择 (1/2/3): " choice

case $choice in
    1)
        # 使用 tmux
        if ! command -v tmux &> /dev/null; then
            echo -e "${RED}❌ tmux 未安装${NC}"
            echo ""
            echo -e "${YELLOW}安装方法:${NC}"
            echo "  Ubuntu/Debian: sudo apt install tmux"
            echo "  macOS:        brew install tmux"
            echo "  CentOS/RHEL:  sudo yum install tmux"
            echo ""
            exit 1
        fi

        echo ""
        echo -e "${GREEN}正在创建 tmux 会话...${NC}"
        
        # 创建新的 tmux 会话
        tmux new-session -d -s trx-dev

        # 分割窗口
        tmux split-window -h

        # 在左侧窗口启动 Frontend
        tmux select-pane -t 0
        tmux send-keys "cd '$PWD' && export GO_ENV=dev AUTO_MIGRATE=false && air -c $FRONTEND_CONFIG" C-m

        # 在右侧窗口启动 Backend
        tmux select-pane -t 1
        tmux send-keys "cd '$PWD' && export GO_ENV=dev AUTO_MIGRATE=false && air -c $BACKEND_CONFIG" C-m

        echo ""
        echo -e "${GREEN}✅ 服务已在 tmux 中启动${NC}"
        echo ""
        echo -e "${CYAN}服务地址:${NC}"
        echo "  Frontend: http://localhost:8080"
        echo "  Backend:  http://localhost:8081"
        echo ""
        echo -e "${YELLOW}查看服务:${NC}"
        echo "  tmux attach -t trx-dev"
        echo ""
        echo -e "${YELLOW}停止服务:${NC}"
        echo "  在 tmux 中按 Ctrl+C 停止各个服务"
        echo "  或运行: tmux kill-session -t trx-dev"
        echo ""
        echo -e "${YELLOW}tmux 快捷键:${NC}"
        echo "  Ctrl+B 然后按 ←/→  : 切换面板"
        echo "  Ctrl+B 然后按 D    : 退出 tmux（服务继续运行）"
        echo "  Ctrl+B 然后按 &    : 关闭当前窗口"
        echo ""
        
        # 询问是否立即连接
        read -p "是否立即连接到 tmux 会话? (y/n): " attach
        if [ "$attach" = "y" ]; then
            tmux attach -t trx-dev
        fi
        ;;
        
    2)
        # 后台启动
        echo ""
        echo -e "${GREEN}正在后台启动服务...${NC}"
        
        # 启动 Frontend
        nohup bash -c "export GO_ENV=dev AUTO_MIGRATE=false; air -c $FRONTEND_CONFIG" > tmp/frontend.log 2>&1 &
        FRONTEND_PID=$!
        
        # 启动 Backend
        nohup bash -c "export GO_ENV=dev AUTO_MIGRATE=false; air -c $BACKEND_CONFIG" > tmp/backend.log 2>&1 &
        BACKEND_PID=$!
        
        # 保存 PID
        echo $FRONTEND_PID > tmp/frontend.pid
        echo $BACKEND_PID > tmp/backend.pid
        
        sleep 2
        
        echo ""
        echo -e "${GREEN}✅ 服务已在后台启动${NC}"
        echo ""
        echo "Frontend PID: $FRONTEND_PID"
        echo "Backend PID:  $BACKEND_PID"
        echo ""
        echo -e "${CYAN}服务地址:${NC}"
        echo "  Frontend: http://localhost:8080"
        echo "  Backend:  http://localhost:8081"
        echo ""
        echo -e "${YELLOW}查看日志:${NC}"
        echo "  tail -f tmp/frontend.log"
        echo "  tail -f tmp/backend.log"
        echo ""
        echo -e "${YELLOW}停止服务:${NC}"
        echo "  kill $FRONTEND_PID $BACKEND_PID"
        echo "  或运行: ./stop-dev.sh"
        echo ""
        ;;
        
    3)
        # 前台启动（仅 Frontend）
        echo ""
        echo -e "${YELLOW}注意: 此模式仅启动 Frontend，Backend 需要在另一个终端手动启动${NC}"
        echo ""
        ./start-frontend.sh
        ;;
        
    *)
        echo -e "${RED}无效的选择${NC}"
        exit 1
        ;;
esac

