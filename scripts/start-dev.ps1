# TRX Project - 同时启动前后端服务
# 用于在 Windows 上同时启动 Frontend 和 Backend 服务

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  🚀 TRX Project 开发环境启动器" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# 检查 Air 是否安装
$airVersion = air -v 2>$null
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ 错误: Air 未安装" -ForegroundColor Red
    Write-Host ""
    Write-Host "请运行以下命令安装:" -ForegroundColor Yellow
    Write-Host "  go install github.com/air-verse/air@latest" -ForegroundColor White
    Write-Host ""
    exit 1
}

Write-Host "✅ Air 已安装" -ForegroundColor Green
Write-Host ""

# 检查端口是否被占用
$port8080 = netstat -ano | Select-String ":8080" | Select-String "LISTENING"
$port8081 = netstat -ano | Select-String ":8081" | Select-String "LISTENING"

if ($port8080) {
    Write-Host "⚠️  警告: 端口 8080 已被占用" -ForegroundColor Yellow
    Write-Host "请先停止占用该端口的进程" -ForegroundColor Yellow
    Write-Host ""
}

if ($port8081) {
    Write-Host "⚠️  警告: 端口 8081 已被占用" -ForegroundColor Yellow
    Write-Host "请先停止占用该端口的进程" -ForegroundColor Yellow
    Write-Host ""
}

if ($port8080 -or $port8081) {
    $continue = Read-Host "是否继续? (y/n)"
    if ($continue -ne "y") {
        exit 0
    }
}

Write-Host "启动方式:" -ForegroundColor Cyan
Write-Host "1. 在两个独立的终端窗口中启动（推荐）" -ForegroundColor White
Write-Host "2. 在后台启动两个服务" -ForegroundColor White
Write-Host ""
$choice = Read-Host "请选择 (1/2)"

if ($choice -eq "1") {
    # 方式 1: 新建两个终端窗口
    Write-Host ""
    Write-Host "正在启动 Frontend 服务..." -ForegroundColor Green
    Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd '$PWD'; .\start-frontend.ps1"
    
    Start-Sleep -Seconds 2
    
    Write-Host "正在启动 Backend 服务..." -ForegroundColor Green
    Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd '$PWD'; .\start-backend.ps1"
    
    Write-Host ""
    Write-Host "✅ 服务已在新窗口中启动" -ForegroundColor Green
    Write-Host ""
    Write-Host "Frontend: http://localhost:8080" -ForegroundColor Cyan
    Write-Host "Backend:  http://localhost:8081" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "提示: 在各自的窗口中按 Ctrl+C 可停止服务" -ForegroundColor Yellow
    
} elseif ($choice -eq "2") {
    # 方式 2: 后台启动
    Write-Host ""
    Write-Host "正在后台启动服务..." -ForegroundColor Green
    
    # 启动 Frontend（后台）
    $frontendJob = Start-Job -ScriptBlock {
        Set-Location $using:PWD
        $env:GO_ENV = "dev"
        $env:AUTO_MIGRATE = "false"
        air -c .air-frontend.toml
    }
    
    # 启动 Backend（后台）
    $backendJob = Start-Job -ScriptBlock {
        Set-Location $using:PWD
        $env:GO_ENV = "dev"
        $env:AUTO_MIGRATE = "false"
        air -c .air-backend.toml
    }
    
    Start-Sleep -Seconds 3
    
    Write-Host ""
    Write-Host "✅ 服务已在后台启动" -ForegroundColor Green
    Write-Host ""
    Write-Host "Frontend Job ID: $($frontendJob.Id)" -ForegroundColor White
    Write-Host "Backend Job ID:  $($backendJob.Id)" -ForegroundColor White
    Write-Host ""
    Write-Host "查看日志:" -ForegroundColor Cyan
    Write-Host "  Receive-Job -Id $($frontendJob.Id) -Keep" -ForegroundColor White
    Write-Host "  Receive-Job -Id $($backendJob.Id) -Keep" -ForegroundColor White
    Write-Host ""
    Write-Host "停止服务:" -ForegroundColor Cyan
    Write-Host "  Stop-Job -Id $($frontendJob.Id)" -ForegroundColor White
    Write-Host "  Stop-Job -Id $($backendJob.Id)" -ForegroundColor White
    Write-Host ""
} else {
    Write-Host "无效的选择" -ForegroundColor Red
    exit 1
}

