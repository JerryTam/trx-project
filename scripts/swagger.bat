@echo off
REM Swagger 文档生成脚本 - Windows 批处理版本
REM 使用方法: scripts\swagger.bat [frontend|backend|all]

setlocal enabledelayedexpansion

REM 切换到项目根目录
cd /d "%~dp0\.."

REM 解析参数
set "TARGET=%~1"
if "%TARGET%"=="" set "TARGET=all"

echo.
echo ================================
echo   Swagger 文档生成器 (Windows)
echo ================================
echo.

REM 检查 swag 是否安装
where swag >nul 2>nul
if %errorlevel% neq 0 (
    echo [错误] swag 未安装！
    echo.
    echo 请运行以下命令安装：
    echo   go install github.com/swaggo/swag/cmd/swag@latest
    echo.
    exit /b 1
)

echo [信息] swag 已安装
echo.

REM 根据参数执行对应操作
if /i "%TARGET%"=="frontend" (
    call :generate_frontend
    if !errorlevel! neq 0 exit /b !errorlevel!
) else if /i "%TARGET%"=="backend" (
    call :generate_backend
    if !errorlevel! neq 0 exit /b !errorlevel!
) else if /i "%TARGET%"=="all" (
    call :generate_frontend
    if !errorlevel! neq 0 exit /b !errorlevel!
    call :generate_backend
    if !errorlevel! neq 0 exit /b !errorlevel!
    echo [成功] 所有文档生成完成！
) else if /i "%TARGET%"=="-h" (
    call :show_help
    exit /b 0
) else if /i "%TARGET%"=="--help" (
    call :show_help
    exit /b 0
) else if /i "%TARGET%"=="help" (
    call :show_help
    exit /b 0
) else (
    echo [错误] 未知参数: %TARGET%
    echo.
    call :show_help
    exit /b 1
)

echo.
echo ================================
echo [成功] 完成！
echo ================================
echo.
exit /b 0

REM ========== 函数定义 ==========

:generate_frontend
echo [信息] 生成前台 Swagger 文档...
echo.

REM 生成文档（排除后台 handler）
swag init ^
    -g cmd/frontend/main.go ^
    -o cmd/frontend/docs ^
    --parseDependency ^
    --parseInternal ^
    --instanceName frontend ^
    --exclude internal/api/handler/backendHandler

if %errorlevel% neq 0 (
    echo [错误] 前台文档生成失败！
    exit /b 1
)

REM 清理旧文件
if exist cmd\frontend\docs\docs.go del /f /q cmd\frontend\docs\docs.go
if exist cmd\frontend\docs\swagger.json del /f /q cmd\frontend\docs\swagger.json
if exist cmd\frontend\docs\swagger.yaml del /f /q cmd\frontend\docs\swagger.yaml

echo.
echo [成功] 前台文档生成完成！
echo [信息] 文档位置: cmd\frontend\docs\frontend_swagger.json
echo [信息] 访问地址: http://localhost:8080/swagger/index.html
echo.
exit /b 0

:generate_backend
echo [信息] 生成后台 Swagger 文档...
echo.

REM 生成文档（排除前台 handler）
swag init ^
    -g cmd/backend/main.go ^
    -o cmd/backend/docs ^
    --parseDependency ^
    --parseInternal ^
    --instanceName backend ^
    --exclude internal/api/handler/frontendhandler

if %errorlevel% neq 0 (
    echo [错误] 后台文档生成失败！
    exit /b 1
)

REM 清理旧文件
if exist cmd\backend\docs\docs.go del /f /q cmd\backend\docs\docs.go
if exist cmd\backend\docs\swagger.json del /f /q cmd\backend\docs\swagger.json
if exist cmd\backend\docs\swagger.yaml del /f /q cmd\backend\docs\swagger.yaml

echo.
echo [成功] 后台文档生成完成！
echo [信息] 文档位置: cmd\backend\docs\backend_swagger.json
echo [信息] 访问地址: http://localhost:8081/swagger/index.html
echo.
exit /b 0

:show_help
echo Swagger 文档生成脚本 (Windows)
echo.
echo 使用方法:
echo   scripts\swagger.bat [frontend^|backend^|all]
echo.
echo 参数:
echo   frontend    生成前台 Swagger 文档
echo   backend     生成后台 Swagger 文档
echo   all         生成所有 Swagger 文档 (默认)
echo.
echo 示例:
echo   scripts\swagger.bat              # 生成所有文档
echo   scripts\swagger.bat frontend     # 只生成前台文档
echo   scripts\swagger.bat backend      # 只生成后台文档
echo.
exit /b 0

