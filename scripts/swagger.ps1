# Swagger 文档生成脚本 - PowerShell 版本
# 使用方法: .\scripts\swagger.ps1 [frontend|backend|all]

param(
    [Parameter(Position=0)]
    [ValidateSet('frontend', 'backend', 'all', 'help')]
    [string]$Target = 'all'
)

# 切换到项目根目录
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$ProjectRoot = Split-Path -Parent $ScriptDir
Set-Location $ProjectRoot

# 颜色输出函数
function Write-Info {
    param([string]$Message)
    Write-Host "ℹ️  $Message" -ForegroundColor Blue
}

function Write-Success {
    param([string]$Message)
    Write-Host "✅ $Message" -ForegroundColor Green
}

function Write-Error {
    param([string]$Message)
    Write-Host "❌ $Message" -ForegroundColor Red
}

function Write-Warning {
    param([string]$Message)
    Write-Host "⚠️  $Message" -ForegroundColor Yellow
}

# 检查 swag 是否安装
function Test-Swag {
    try {
        $null = Get-Command swag -ErrorAction Stop
        $version = swag --version 2>&1
        Write-Success "swag 已安装: $version"
        return $true
    }
    catch {
        Write-Error "swag 未安装！"
        Write-Host ""
        Write-Host "请运行以下命令安装："
        Write-Host "  go install github.com/swaggo/swag/cmd/swag@latest"
        Write-Host ""
        return $false
    }
}

# 生成前台文档
function New-FrontendDocs {
    Write-Info "生成前台 Swagger 文档..."
    
    # 生成文档（排除后台 handler）
    $output = swag init `
        -g cmd/frontend/main.go `
        -o cmd/frontend/docs `
        --parseDependency `
        --parseInternal `
        --instanceName frontend `
        --exclude internal/api/handler/backendHandler 2>&1
    
    if ($LASTEXITCODE -ne 0) {
        Write-Error "前台文档生成失败！"
        Write-Host $output
        exit 1
    }
    
    # 清理旧文件
    Remove-Item -Path "cmd/frontend/docs/docs.go" -ErrorAction SilentlyContinue
    Remove-Item -Path "cmd/frontend/docs/swagger.json" -ErrorAction SilentlyContinue
    Remove-Item -Path "cmd/frontend/docs/swagger.yaml" -ErrorAction SilentlyContinue
    
    Write-Success "前台文档生成完成！"
    Write-Info "文档位置: cmd/frontend/docs/frontend_swagger.json"
    Write-Info "访问地址: http://localhost:8080/swagger/index.html"
    Write-Host ""
}

# 生成后台文档
function New-BackendDocs {
    Write-Info "生成后台 Swagger 文档..."
    
    # 生成文档（排除前台 handler）
    $output = swag init `
        -g cmd/backend/main.go `
        -o cmd/backend/docs `
        --parseDependency `
        --parseInternal `
        --instanceName backend `
        --exclude internal/api/handler/frontendhandler 2>&1
    
    if ($LASTEXITCODE -ne 0) {
        Write-Error "后台文档生成失败！"
        Write-Host $output
        exit 1
    }
    
    # 清理旧文件
    Remove-Item -Path "cmd/backend/docs/docs.go" -ErrorAction SilentlyContinue
    Remove-Item -Path "cmd/backend/docs/swagger.json" -ErrorAction SilentlyContinue
    Remove-Item -Path "cmd/backend/docs/swagger.yaml" -ErrorAction SilentlyContinue
    
    Write-Success "后台文档生成完成！"
    Write-Info "文档位置: cmd/backend/docs/backend_swagger.json"
    Write-Info "访问地址: http://localhost:8081/swagger/index.html"
    Write-Host ""
}

# 显示帮助信息
function Show-Help {
    Write-Host "Swagger 文档生成脚本 (PowerShell)"
    Write-Host ""
    Write-Host "使用方法:"
    Write-Host "  .\scripts\swagger.ps1 [frontend|backend|all]"
    Write-Host ""
    Write-Host "参数:"
    Write-Host "  frontend    生成前台 Swagger 文档"
    Write-Host "  backend     生成后台 Swagger 文档"
    Write-Host "  all         生成所有 Swagger 文档 (默认)"
    Write-Host ""
    Write-Host "示例:"
    Write-Host "  .\scripts\swagger.ps1              # 生成所有文档"
    Write-Host "  .\scripts\swagger.ps1 frontend     # 只生成前台文档"
    Write-Host "  .\scripts\swagger.ps1 backend      # 只生成后台文档"
    Write-Host ""
}

# 主函数
function Main {
    Write-Host ""
    Write-Host "================================" -ForegroundColor Cyan
    Write-Host "   Swagger 文档生成器 (PowerShell)" -ForegroundColor Cyan
    Write-Host "================================" -ForegroundColor Cyan
    Write-Host ""
    
    # 显示帮助
    if ($Target -eq 'help') {
        Show-Help
        exit 0
    }
    
    # 检查依赖
    if (-not (Test-Swag)) {
        exit 1
    }
    Write-Host ""
    
    # 根据参数执行
    switch ($Target) {
        'frontend' {
            New-FrontendDocs
        }
        'backend' {
            New-BackendDocs
        }
        'all' {
            New-FrontendDocs
            New-BackendDocs
            Write-Success "所有文档生成完成！"
        }
    }
    
    Write-Host "================================" -ForegroundColor Cyan
    Write-Success "完成！"
    Write-Host ""
}

# 运行主函数
Main

