#!/bin/bash

# Swagger 文档生成脚本 - 支持 Linux、macOS、Git Bash (Windows)
# 使用方法: ./scripts/swagger.sh [frontend|backend|all]

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 项目根目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# 切换到项目根目录
cd "$PROJECT_ROOT"

# 打印带颜色的消息
print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

# 检查 swag 是否安装
check_swag() {
    if ! command -v swag &> /dev/null; then
        print_error "swag 未安装！"
        echo ""
        echo "请运行以下命令安装："
        echo "  go install github.com/swaggo/swag/cmd/swag@latest"
        echo ""
        exit 1
    fi
    print_success "swag 已安装: $(swag --version)"
}

# 生成前台文档
generate_frontend() {
    print_info "生成前台 Swagger 文档..."
    
    # 扫描时排除后台 handler 目录
    swag init \
        -g cmd/frontend/main.go \
        -o cmd/frontend/docs \
        --parseDependency \
        --parseInternal \
        --instanceName frontend \
        --exclude internal/api/handler/backendHandler
    
    if [ $? -ne 0 ]; then
        print_error "前台文档生成失败！"
        exit 1
    fi
    
    # 清理旧文件
    rm -f cmd/frontend/docs/docs.go
    rm -f cmd/frontend/docs/swagger.json
    rm -f cmd/frontend/docs/swagger.yaml
    
    print_success "前台文档生成完成！"
    print_info "文档位置: cmd/frontend/docs/frontend_swagger.json"
    print_info "访问地址: http://localhost:8080/swagger/index.html"
    echo ""
}

# 生成后台文档
generate_backend() {
    print_info "生成后台 Swagger 文档..."
    
    # 扫描时排除前台 handler 目录
    swag init \
        -g cmd/backend/main.go \
        -o cmd/backend/docs \
        --parseDependency \
        --parseInternal \
        --instanceName backend \
        --exclude internal/api/handler/frontendHandler
    
    if [ $? -ne 0 ]; then
        print_error "后台文档生成失败！"
        exit 1
    fi
    
    # 清理旧文件
    rm -f cmd/backend/docs/docs.go
    rm -f cmd/backend/docs/swagger.json
    rm -f cmd/backend/docs/swagger.yaml
    
    print_success "后台文档生成完成！"
    print_info "文档位置: cmd/backend/docs/backend_swagger.json"
    print_info "访问地址: http://localhost:8081/swagger/index.html"
    echo ""
}

# 显示使用帮助
show_help() {
    echo "Swagger 文档生成脚本"
    echo ""
    echo "使用方法:"
    echo "  $0 [frontend|backend|all]"
    echo ""
    echo "参数:"
    echo "  frontend    生成前台 Swagger 文档"
    echo "  backend     生成后台 Swagger 文档"
    echo "  all         生成所有 Swagger 文档（默认）"
    echo ""
    echo "示例:"
    echo "  $0              # 生成所有文档"
    echo "  $0 frontend     # 只生成前台文档"
    echo "  $0 backend      # 只生成后台文档"
    echo ""
}

# 主函数
main() {
    echo ""
    echo "🚀 Swagger 文档生成器"
    echo "================================"
    echo ""
    
    # 检查依赖
    check_swag
    echo ""
    
    # 解析参数
    TARGET="${1:-all}"
    
    case "$TARGET" in
        frontend)
            generate_frontend
            ;;
        backend)
            generate_backend
            ;;
        all)
            generate_frontend
            generate_backend
            print_success "所有文档生成完成！"
            ;;
        -h|--help|help)
            show_help
            exit 0
            ;;
        *)
            print_error "未知参数: $TARGET"
            echo ""
            show_help
            exit 1
            ;;
    esac
    
    echo "================================"
    print_success "完成！"
}

# 运行主函数
main "$@"

