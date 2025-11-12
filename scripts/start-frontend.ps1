# TRX Project - Frontend 启动脚本
# 用于在 Windows 上正确设置环境变量并启动 Air 热重载

# 设置环境变量
$env:GO_ENV = "dev"
$env:AUTO_MIGRATE = "false"

# 清理旧的编译文件
if (Test-Path "tmp/frontend.exe") {
    Write-Host "清理旧的编译文件..." -ForegroundColor Yellow
    Remove-Item "tmp/frontend.exe" -Force
}

# 显示启动信息
Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  🚀 启动 Frontend 服务 (热重载模式)" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "环境变量:" -ForegroundColor Green
Write-Host "  GO_ENV       = $env:GO_ENV" -ForegroundColor White
Write-Host "  AUTO_MIGRATE = $env:AUTO_MIGRATE" -ForegroundColor White
Write-Host ""
Write-Host "服务地址:" -ForegroundColor Green
Write-Host "  HTTP:    http://localhost:8080" -ForegroundColor White
Write-Host "  Health:  http://localhost:8080/health" -ForegroundColor White
Write-Host "  Swagger: http://localhost:8080/swagger/index.html" -ForegroundColor White
Write-Host ""
Write-Host "提示: 按 Ctrl+C 停止服务" -ForegroundColor Yellow
Write-Host ""

# 启动 Air
air -c .air-frontend.toml

