#!/bin/bash

# 项目验证脚本

echo "================================"
echo "TRX Project 验证"
echo "================================"
echo ""

# 检查 Go 版本
echo "1. Go 版本检查"
go version
echo ""

# 检查项目结构
echo "2. 项目结构检查"
echo "✓ cmd/server - 入口程序"
ls -la cmd/server/*.go | head -3
echo "✓ internal - 私有代码"
ls -d internal/*/ | head -5
echo "✓ pkg - 公共库"
ls -d pkg/*/ | head -10
echo ""

# 检查依赖
echo "3. Go 依赖检查"
go list -m all | head -10
echo ""

# 检查配置文件
echo "4. 配置文件检查"
if [ -f "config/config.yaml.example" ]; then
    echo "✓ config/config.yaml.example 存在"
else
    echo "✗ config/config.yaml.example 不存在"
fi
echo ""

# 检查文档
echo "5. 文档检查"
ls -1 *.md
echo "---"
ls -1 docs/*.md
echo ""

# 编译测试
echo "6. 编译测试"
if go build -o /tmp/server_test ./cmd/server; then
    echo "✅ 项目编译成功"
    rm /tmp/server_test
else
    echo "❌ 项目编译失败"
    exit 1
fi
echo ""

# 统计信息
echo "7. 项目统计"
echo "Go 文件数: $(find . -name '*.go' -not -path './vendor/*' | wc -l)"
echo "测试文件数: $(find . -name '*_test.go' | wc -l)"
echo "文档文件数: $(find . -name '*.md' | wc -l)"
echo "代码行数: $(find . -name '*.go' -not -path './vendor/*' -exec cat {} \; | wc -l)"
echo ""

# 检查响应包
echo "8. 统一响应格式检查"
if [ -f "pkg/response/response.go" ]; then
    echo "✅ pkg/response/response.go 存在"
    grep -c "func Success" pkg/response/response.go | xargs echo "- Success 函数数量:"
    grep -c "func.*Error" pkg/response/response.go | xargs echo "- Error 函数数量:"
else
    echo "❌ pkg/response/response.go 不存在"
fi

if [ -f "pkg/response/code.go" ]; then
    echo "✅ pkg/response/code.go 存在"
    grep -c "^const" pkg/response/code.go | xargs echo "- 状态码常量数量:"
else
    echo "❌ pkg/response/code.go 不存在"
fi
echo ""

echo "================================"
echo "✅ 项目验证完成"
echo "================================"
echo ""
echo "下一步："
echo "1. 启动依赖服务: make docker-up"
echo "2. 运行项目: make run"
echo "3. 测试 API: ./scripts/test_api.sh"
echo ""
