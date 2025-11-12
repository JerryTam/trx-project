# TRX Project - Swagger 文档重新生成脚本 (Windows PowerShell)

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  🔄 重新生成 Swagger 文档" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# 检查 swag 是否安装
$swagVersion = swag --version 2>$null
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ 错误: swag 未安装" -ForegroundColor Red
    Write-Host ""
    Write-Host "请运行以下命令安装:" -ForegroundColor Yellow
    Write-Host "  go install github.com/swaggo/swag/cmd/swag@latest" -ForegroundColor White
    Write-Host ""
    exit 1
}

Write-Host "✅ swag 已安装" -ForegroundColor Green
swag --version
Write-Host ""

# 生成前端文档
Write-Host "📝 生成前端 Swagger 文档..." -ForegroundColor Yellow
$result = swag init -g cmd/frontend/main.go -o cmd/frontend/docs --parseDependency --parseInternal
if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ 前端文档生成成功" -ForegroundColor Green
} else {
    Write-Host "❌ 前端文档生成失败" -ForegroundColor Red
    exit 1
}
Write-Host ""

# 生成后端文档
Write-Host "📝 生成后端 Swagger 文档..." -ForegroundColor Yellow
$result = swag init -g cmd/backend/main.go -o cmd/backend/docs --parseDependency --parseInternal
if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ 后端文档生成成功" -ForegroundColor Green
} else {
    Write-Host "❌ 后端文档生成失败" -ForegroundColor Red
    exit 1
}
Write-Host ""

# 验证生成的文件
Write-Host "🔍 验证生成的文件..." -ForegroundColor Yellow

$frontendJson = "cmd\frontend\docs\swagger.json"
$backendJson = "cmd\backend\docs\swagger.json"

if (Test-Path $frontendJson) {
    $frontendSize = (Get-Item $frontendJson).Length
    if ($frontendSize -gt 1000) {
        Write-Host "✅ 前端 swagger.json: $frontendSize bytes" -ForegroundColor Green
    } else {
        Write-Host "⚠️  前端 swagger.json 文件过小，可能生成不完整" -ForegroundColor Yellow
    }
} else {
    Write-Host "❌ 前端 swagger.json 文件不存在" -ForegroundColor Red
}

if (Test-Path $backendJson) {
    $backendSize = (Get-Item $backendJson).Length
    if ($backendSize -gt 1000) {
        Write-Host "✅ 后端 swagger.json: $backendSize bytes" -ForegroundColor Green
    } else {
        Write-Host "⚠️  后端 swagger.json 文件过小，可能生成不完整" -ForegroundColor Yellow
    }
} else {
    Write-Host "❌ 后端 swagger.json 文件不存在" -ForegroundColor Red
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Green
Write-Host "  ✅ Swagger 文档生成完成！" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Green
Write-Host ""
Write-Host "📚 访问地址:" -ForegroundColor Cyan
Write-Host "  前端: http://localhost:8080/swagger/index.html" -ForegroundColor White
Write-Host "  后端: http://localhost:8081/swagger/index.html" -ForegroundColor White
Write-Host ""
Write-Host "💡 提示:" -ForegroundColor Yellow
Write-Host "  1. 如果服务正在运行，请重启以加载新文档" -ForegroundColor White
Write-Host "  2. 使用 Air 热重载会自动重新编译" -ForegroundColor White
Write-Host "  3. 清理缓存: Remove-Item tmp -Recurse -Force; New-Item -ItemType Directory tmp" -ForegroundColor White
Write-Host ""

