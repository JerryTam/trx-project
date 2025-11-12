#!/bin/bash
# TRX Project - 停止开发服务
# 跨平台支持: Linux/macOS/Windows(Git Bash)

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

echo ""
echo -e "${CYAN}========================================${NC}"
echo -e "${CYAN}  🛑 停止 TRX Project 开发服务${NC}"
echo -e "${CYAN}========================================${NC}"
echo ""

STOPPED=false

# 方式 1: 通过 PID 文件停止
if [ -f "tmp/frontend.pid" ] || [ -f "tmp/backend.pid" ]; then
    echo -e "${YELLOW}通过 PID 文件停止服务...${NC}"
    
    if [ -f "tmp/frontend.pid" ]; then
        FRONTEND_PID=$(cat tmp/frontend.pid)
        if kill -0 $FRONTEND_PID 2>/dev/null; then
            kill $FRONTEND_PID
            echo -e "${GREEN}✅ Frontend 服务已停止 (PID: $FRONTEND_PID)${NC}"
            STOPPED=true
        fi
        rm -f tmp/frontend.pid
    fi
    
    if [ -f "tmp/backend.pid" ]; then
        BACKEND_PID=$(cat tmp/backend.pid)
        if kill -0 $BACKEND_PID 2>/dev/null; then
            kill $BACKEND_PID
            echo -e "${GREEN}✅ Backend 服务已停止 (PID: $BACKEND_PID)${NC}"
            STOPPED=true
        fi
        rm -f tmp/backend.pid
    fi
fi

# 方式 2: 通过端口查找并停止
stop_by_port() {
    local port=$1
    local name=$2
    
    if command -v lsof &> /dev/null; then
        # macOS/Linux with lsof
        PID=$(lsof -ti :$port)
        if [ ! -z "$PID" ]; then
            kill $PID
            echo -e "${GREEN}✅ ${name} 服务已停止 (端口 ${port}, PID: $PID)${NC}"
            STOPPED=true
        fi
    elif command -v netstat &> /dev/null; then
        # Linux with netstat
        if [[ "$OSTYPE" == "linux-gnu"* ]]; then
            PID=$(netstat -tlnp 2>/dev/null | grep ":$port " | awk '{print $7}' | cut -d'/' -f1)
            if [ ! -z "$PID" ]; then
                kill $PID
                echo -e "${GREEN}✅ ${name} 服务已停止 (端口 ${port}, PID: $PID)${NC}"
                STOPPED=true
            fi
        fi
    fi
}

echo ""
echo -e "${YELLOW}检查端口占用...${NC}"
stop_by_port 8080 "Frontend"
stop_by_port 8081 "Backend"

# 方式 3: 停止 tmux 会话
if tmux has-session -t trx-dev 2>/dev/null; then
    echo ""
    echo -e "${YELLOW}发现 tmux 会话 'trx-dev'${NC}"
    read -p "是否停止 tmux 会话? (y/n): " kill_tmux
    if [ "$kill_tmux" = "y" ]; then
        tmux kill-session -t trx-dev
        echo -e "${GREEN}✅ tmux 会话已停止${NC}"
        STOPPED=true
    fi
fi

# 方式 4: 查找所有 air 进程
echo ""
echo -e "${YELLOW}检查 air 进程...${NC}"
if command -v pgrep &> /dev/null; then
    AIR_PIDS=$(pgrep -f "air -c .air-")
    if [ ! -z "$AIR_PIDS" ]; then
        echo "发现 air 进程: $AIR_PIDS"
        read -p "是否停止所有 air 进程? (y/n): " kill_air
        if [ "$kill_air" = "y" ]; then
            kill $AIR_PIDS
            echo -e "${GREEN}✅ air 进程已停止${NC}"
            STOPPED=true
        fi
    fi
fi

echo ""
if [ "$STOPPED" = true ]; then
    echo -e "${GREEN}✅ 服务已停止${NC}"
else
    echo -e "${YELLOW}⚠️  未发现运行中的服务${NC}"
fi
echo ""

