# Swagger 跨平台脚本实现总结

## 实现背景

在之前实现 Swagger 文档分离后，用户提出需要跨平台的脚本来生成文档。主要原因：

1. **Windows 用户**: `make` 命令不可用或需要额外安装
2. **跨平台需求**: 团队成员使用不同操作系统
3. **CI/CD 集成**: 需要在不同环境中运行
4. **易用性**: 希望有简单的命令行工具

## 解决方案

创建了三个跨平台脚本，覆盖所有主流平台和 Shell 环境：

### 1. Bash 脚本 (`swagger.sh`)

**适用平台**:
- ✅ Linux
- ✅ macOS
- ✅ Windows (Git Bash)

**特点**:
- ANSI 颜色输出
- 完整的错误处理
- 详细的进度信息
- 自动依赖检查

**实现亮点**:
```bash
#!/bin/bash
set -e  # 任何命令失败立即退出

# 彩色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'

# 自动定位项目根目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
cd "$PROJECT_ROOT"
```

### 2. Windows 批处理脚本 (`swagger.bat`)

**适用平台**:
- ✅ Windows CMD
- ✅ Windows PowerShell (兼容模式)

**特点**:
- 原生 Windows 支持
- 不依赖 Unix 工具
- 清晰的状态输出
- 标准的退出码

**实现亮点**:
```batch
@echo off
setlocal enabledelayedexpansion

REM 切换到项目根目录
cd /d "%~dp0\.."

REM 检查依赖
where swag >nul 2>nul
if %errorlevel% neq 0 (
    echo [错误] swag 未安装！
    exit /b 1
)

REM 使用函数组织代码
call :generate_frontend
```

### 3. PowerShell 脚本 (`swagger.ps1`)

**适用平台**:
- ✅ Windows (PowerShell 5.1+)
- ✅ Linux (PowerShell Core)
- ✅ macOS (PowerShell Core)

**特点**:
- 现代化语法
- 强类型参数验证
- 彩色输出支持
- 对象化错误处理

**实现亮点**:
```powershell
param(
    [Parameter(Position=0)]
    [ValidateSet('frontend', 'backend', 'all', 'help')]
    [string]$Target = 'all'
)

function Write-Success {
    param([string]$Message)
    Write-Host "✅ $Message" -ForegroundColor Green
}

try {
    $null = Get-Command swag -ErrorAction Stop
    return $true
}
catch {
    Write-Error "swag 未安装！"
    return $false
}
```

## 统一的工作流程

所有脚本执行相同的流程：

```
┌─────────────────────────────┐
│ 1. 检查 swag 是否安装       │
└─────────────┬───────────────┘
              ↓
┌─────────────────────────────┐
│ 2. 生成原始 Swagger 文档    │
│    - swag init              │
│    - 指定 instanceName      │
└─────────────┬───────────────┘
              ↓
┌─────────────────────────────┐
│ 3. 过滤不需要的接口         │
│    - go run filter_swagger  │
│    - 移除对方服务的接口     │
└─────────────┬───────────────┘
              ↓
┌─────────────────────────────┐
│ 4. 清理旧文档文件           │
│    - 删除 docs.go           │
│    - 删除 swagger.json      │
└─────────────┬───────────────┘
              ↓
┌─────────────────────────────┐
│ 5. 显示成功信息和访问地址   │
└─────────────────────────────┘
```

## 核心功能实现

### 1. 依赖检查

**Bash**:
```bash
check_swag() {
    if ! command -v swag &> /dev/null; then
        print_error "swag 未安装！"
        echo "请运行: go install github.com/swaggo/swag/cmd/swag@latest"
        exit 1
    fi
}
```

**PowerShell**:
```powershell
function Test-Swag {
    try {
        $null = Get-Command swag -ErrorAction Stop
        return $true
    }
    catch {
        Write-Error "swag 未安装！"
        return $false
    }
}
```

### 2. 参数解析

**Bash**:
```bash
TARGET="${1:-all}"  # 默认值为 all

case "$TARGET" in
    frontend) generate_frontend ;;
    backend) generate_backend ;;
    all) generate_frontend && generate_backend ;;
    *) echo "未知参数"; exit 1 ;;
esac
```

**PowerShell**:
```powershell
param(
    [ValidateSet('frontend', 'backend', 'all')]
    [string]$Target = 'all'
)

switch ($Target) {
    'frontend' { New-FrontendDocs }
    'backend' { New-BackendDocs }
    'all' { 
        New-FrontendDocs
        New-BackendDocs 
    }
}
```

### 3. 错误处理

**Bash**:
```bash
set -e  # 任何命令失败立即退出

if [ $? -ne 0 ]; then
    print_error "生成失败！"
    exit 1
fi
```

**Batch**:
```batch
if %errorlevel% neq 0 (
    echo [错误] 生成失败！
    exit /b 1
)
```

**PowerShell**:
```powershell
if ($LASTEXITCODE -ne 0) {
    Write-Error "生成失败！"
    exit 1
}
```

## 平台兼容性处理

### 路径处理

**Bash/PowerShell**:
```bash
# Unix 风格路径分隔符
cmd/frontend/docs/frontend_swagger.json
```

**Batch**:
```batch
REM Windows 风格路径分隔符
cmd\frontend\docs\frontend_swagger.json
```

### 命令续行

**Bash**:
```bash
swag init \
    -g cmd/frontend/main.go \
    -o cmd/frontend/docs \
    --instanceName frontend
```

**Batch**:
```batch
swag init ^
    -g cmd/frontend/main.go ^
    -o cmd/frontend/docs ^
    --instanceName frontend
```

**PowerShell**:
```powershell
swag init `
    -g cmd/frontend/main.go `
    -o cmd/frontend/docs `
    --instanceName frontend
```

### 文件删除

**Bash**:
```bash
rm -f file.txt
```

**Batch**:
```batch
if exist file.txt del /f /q file.txt
```

**PowerShell**:
```powershell
Remove-Item -Path "file.txt" -ErrorAction SilentlyContinue
```

## 用户体验优化

### 1. 彩色输出

使用不同颜色表示不同状态：
- 🔵 蓝色 - 信息提示
- 🟢 绿色 - 成功消息
- 🔴 红色 - 错误消息
- 🟡 黄色 - 警告消息

### 2. 进度反馈

```
🚀 Swagger 文档生成器
================================

✅ swag 已安装: swag version v1.16.4

ℹ️  生成前台 Swagger 文档...
ℹ️  过滤后台接口 (/admin, /users)...
📊 Filtered 14 paths with prefixes [/admin /users]
✅ 前台文档生成完成！
ℹ️  文档位置: cmd/frontend/docs/frontend_swagger.json
ℹ️  访问地址: http://localhost:8080/swagger/index.html
```

### 3. 帮助信息

所有脚本都支持 `help` 参数：

```bash
./scripts/swagger.sh help
scripts\swagger.bat help
.\scripts\swagger.ps1 help
```

## 测试验证

### 测试场景

1. ✅ **依赖检查** - 未安装 swag 时的提示
2. ✅ **正常生成** - 生成前台、后台、全部文档
3. ✅ **错误处理** - 生成失败时的错误提示
4. ✅ **参数验证** - 无效参数的处理
5. ✅ **文件清理** - 旧文件的正确删除
6. ✅ **过滤功能** - 接口正确过滤

### 测试结果

```bash
# 前台文档测试
$ ./scripts/swagger.sh frontend
✅ 生成成功
✅ 只包含 3 个接口：/public/login, /public/register, /user/profile

# 后台文档测试
$ ./scripts/swagger.sh backend
✅ 生成成功
✅ 只包含 12 个接口：/admin/* 相关接口

# 全部生成测试
$ ./scripts/swagger.sh
✅ 前台和后台文档都生成成功
```

## 文档更新

### 新增文档

1. **`docs/SWAGGER_SCRIPTS_GUIDE.md`**
   - 详细的脚本使用指南
   - 平台支持说明
   - 常见问题解答

2. **`scripts/README.md`**
   - 脚本目录总览
   - 快速开始指南
   - 故障排查

3. **`docs/SWAGGER_SCRIPTS_IMPLEMENTATION.md`** (本文档)
   - 实现细节和技术决策

### 更新文档

1. **`README.md`**
   - 添加跨平台脚本使用说明
   - 更新项目结构

2. **`docs/SWAGGER_GUIDE.md`**
   - 补充脚本生成方式

## 技术决策

### 为什么提供三个脚本而不是一个？

**优点**:
- ✅ 每个平台使用原生工具
- ✅ 更好的兼容性
- ✅ 更少的依赖
- ✅ 更好的错误处理

**缺点**:
- ❌ 需要维护多个脚本
- ❌ 代码重复

**决策**: 优点远大于缺点，特别是考虑到用户体验。

### 为什么不使用 Go 程序？

**考虑过的方案**: 用 Go 写一个跨平台工具

**不采用的原因**:
- 需要编译
- 增加项目复杂度
- Shell 脚本更简单直观
- 已经有 `filter_swagger.go`

### 为什么保留 Makefile？

**原因**:
- 开发者熟悉 make
- Makefile 可以定义其他任务
- 两者可以共存

**结论**: Makefile 和脚本互为补充，用户可以选择。

## 性能对比

| 操作 | Makefile | Bash 脚本 | Batch 脚本 | PS 脚本 |
|------|----------|-----------|------------|---------|
| 启动时间 | 快 | 快 | 中等 | 慢 |
| 执行速度 | 快 | 快 | 快 | 快 |
| 依赖检查 | ❌ | ✅ | ✅ | ✅ |
| 错误提示 | 基础 | 详细 | 详细 | 详细 |

## CI/CD 集成

### GitHub Actions

```yaml
name: Generate Swagger Docs

on: [push, pull_request]

jobs:
  docs:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
      
      - name: Install swag
        run: go install github.com/swaggo/swag/cmd/swag@latest
      
      - name: Generate docs
        run: bash scripts/swagger.sh
      
      - name: Check files exist
        run: |
          test -f cmd/frontend/docs/frontend_swagger.json
          test -f cmd/backend/docs/backend_swagger.json
```

### GitLab CI

```yaml
swagger:
  stage: build
  script:
    - go install github.com/swaggo/swag/cmd/swag@latest
    - bash scripts/swagger.sh
  artifacts:
    paths:
      - cmd/frontend/docs/
      - cmd/backend/docs/
```

## 未来改进

### 短期

1. [ ] 添加 `--watch` 模式，文件变化时自动重新生成
2. [ ] 支持自定义过滤规则配置文件
3. [ ] 添加更多的验证检查

### 长期

1. [ ] 开发 Web UI 用于文档生成
2. [ ] 集成到开发工具插件（VS Code、GoLand）
3. [ ] 支持更多文档格式（OpenAPI 3.0、AsyncAPI）

## 总结

通过实现三个跨平台脚本，我们实现了：

✅ **全平台支持** - Linux、macOS、Windows 全覆盖  
✅ **易用性** - 简单的命令行界面  
✅ **自动化** - 完整的生成、过滤、清理流程  
✅ **友好提示** - 彩色输出和详细信息  
✅ **错误处理** - 清晰的错误提示  
✅ **CI/CD 就绪** - 适合自动化环境  

这些脚本与现有的 Makefile 互为补充，为不同环境的用户提供了最佳的使用体验。

## 相关文件

- `scripts/swagger.sh` - Bash 脚本
- `scripts/swagger.bat` - 批处理脚本
- `scripts/swagger.ps1` - PowerShell 脚本
- `scripts/filter_swagger.go` - 过滤工具
- `scripts/README.md` - 脚本目录说明
- `docs/SWAGGER_SCRIPTS_GUIDE.md` - 使用指南
- `Makefile` - Make 构建文件

