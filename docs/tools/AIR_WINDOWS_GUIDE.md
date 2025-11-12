# Air 在 Windows 上的使用指南

## 🪟 Windows 特殊配置说明

由于 Windows 系统的特殊性，Air 配置需要做一些调整。

### ⚠️ 关键问题

之前的错误：
```
CMD will not recognize non .exe file for execution
Process Exit with Code 0
```

**原因**：Windows 需要 `.exe` 扩展名才能执行程序。

### ✅ 已修复的配置

#### `.air-frontend.toml`
```toml
[build]
  bin = "./tmp/frontend.exe"           # ✅ 添加 .exe
  cmd = "go build -o ./tmp/frontend.exe ./cmd/frontend"  # ✅ 添加 .exe
  full_bin = "GO_ENV=dev AUTO_MIGRATE=false ./tmp/frontend.exe"  # ✅ 添加 .exe
```

#### `.air-backend.toml`
```toml
[build]
  bin = "./tmp/backend.exe"            # ✅ 添加 .exe
  cmd = "go build -o ./tmp/backend.exe ./cmd/backend"   # ✅ 添加 .exe
  full_bin = "GO_ENV=dev AUTO_MIGRATE=false ./tmp/backend.exe"  # ✅ 添加 .exe
```

---

## 🚀 现在重新启动服务

### 步骤 1：清理旧文件

```bash
# PowerShell
Remove-Item -Recurse -Force tmp

# 或 CMD
rmdir /s /q tmp

# 或使用 Makefile
make clean
```

### 步骤 2：重新启动前端服务

```bash
air -c .air-frontend.toml
```

**应该看到正常输出**：
```
building...
running...
[GIN-debug] Listening and serving HTTP on 0.0.0.0:8080
Frontend server starting on 0.0.0.0:8080
```

### 步骤 3：测试访问

打开浏览器访问：
- http://localhost:8080/health
- http://localhost:8080/swagger/index.html

---

## 🔍 验证服务是否正常运行

### 方法 1：使用 PowerShell

```powershell
# 检查端口是否在监听
netstat -ano | findstr :8080

# 应该看到类似输出：
# TCP    0.0.0.0:8080    0.0.0.0:0    LISTENING    12345
```

### 方法 2：使用 curl（如果已安装）

```powershell
curl http://localhost:8080/health

# 应该返回：
# {"status":"ok","service":"frontend"}
```

### 方法 3：使用浏览器

直接访问：http://localhost:8080/health

---

## 🐛 Windows 特有问题排查

### 问题 1：防火墙阻止

**症状**：程序启动但无法访问

**解决**：
```powershell
# 1. 临时关闭防火墙测试
# Windows 设置 → 更新和安全 → Windows 安全中心 → 防火墙和网络保护

# 2. 或添加防火墙规则
New-NetFirewallRule -DisplayName "TRX Frontend" -Direction Inbound -LocalPort 8080 -Protocol TCP -Action Allow
New-NetFirewallRule -DisplayName "TRX Backend" -Direction Inbound -LocalPort 8081 -Protocol TCP -Action Allow
```

### 问题 2：端口被占用

**症状**：`bind: address already in use`

**解决**：
```powershell
# 1. 查找占用端口的进程
netstat -ano | findstr :8080

# 2. 记下 PID（最后一列的数字）
# 例如：TCP  0.0.0.0:8080  0.0.0.0:0  LISTENING  12345

# 3. 杀死进程
taskkill /F /PID 12345

# 或使用 GUI 方式：任务管理器 → 详细信息 → 找到 PID → 结束任务
```

### 问题 3：路径问题

**症状**：`cannot find the path specified`

**解决**：
```powershell
# 确保在项目根目录
cd D:\workspace\go\trx-project

# 创建 tmp 目录
mkdir tmp -Force

# 检查配置文件是否存在
Test-Path .air-frontend.toml
Test-Path .air-backend.toml
```

### 问题 4：环境变量设置失败

**症状**：`GO_ENV=dev` 不生效

**解决**：
```powershell
# Windows PowerShell 方式 1：使用 full_bin（已配置）
# Air 会自动处理

# 方式 2：手动设置环境变量
$env:GO_ENV = "dev"
$env:AUTO_MIGRATE = "false"
air -c .air-frontend.toml

# 方式 3：创建启动脚本（推荐）
# 见下方脚本部分
```

### 问题 5：权限问题

**症状**：`Access is denied`

**解决**：
```powershell
# 以管理员身份运行 PowerShell
# 右键点击 PowerShell → 以管理员身份运行
```

---

## 📝 推荐的启动脚本

### `start-frontend.ps1` - PowerShell 脚本

```powershell
# 设置环境变量
$env:GO_ENV = "dev"
$env:AUTO_MIGRATE = "false"

# 清理旧的编译文件
if (Test-Path "tmp/frontend.exe") {
    Remove-Item "tmp/frontend.exe" -Force
}

# 启动 Air
Write-Host "🚀 Starting Frontend Service..." -ForegroundColor Green
air -c .air-frontend.toml
```

### `start-backend.ps1` - PowerShell 脚本

```powershell
# 设置环境变量
$env:GO_ENV = "dev"
$env:AUTO_MIGRATE = "false"

# 清理旧的编译文件
if (Test-Path "tmp/backend.exe") {
    Remove-Item "tmp/backend.exe" -Force
}

# 启动 Air
Write-Host "🚀 Starting Backend Service..." -ForegroundColor Green
air -c .air-backend.toml
```

### 使用脚本

```powershell
# 创建脚本文件
notepad start-frontend.ps1
notepad start-backend.ps1

# 复制上面的内容到文件中保存

# 可能需要修改执行策略（首次运行）
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser

# 启动服务
.\start-frontend.ps1
.\start-backend.ps1
```

---

## 🎯 完整的开发流程（Windows）

### 1. 打开两个 PowerShell 终端

**终端 1（前端）**：
```powershell
cd D:\workspace\go\trx-project
air -c .air-frontend.toml
```

**终端 2（后端）**：
```powershell
cd D:\workspace\go\trx-project
air -c .air-backend.toml
```

### 2. 等待服务启动

看到以下输出表示成功：
```
Frontend server starting on 0.0.0.0:8080
```

### 3. 测试访问

浏览器打开：
- http://localhost:8080/health
- http://localhost:8081/health

### 4. 修改代码

编辑任意 `.go` 文件，保存后 Air 会自动重新编译和重启。

### 5. 停止服务

在终端按 `Ctrl+C`

---

## 💡 Windows 开发技巧

### 技巧 1：使用 Windows Terminal

**推荐使用 Windows Terminal**（而不是 CMD 或单独的 PowerShell）：

1. 安装 Windows Terminal（Microsoft Store）
2. 可以在一个窗口中开多个标签页
3. 支持分屏显示

```powershell
# 在 Windows Terminal 中
# 按 Ctrl+Shift+T 新建标签页
# 按 Alt+Shift++ 垂直分屏
# 按 Alt+Shift+- 水平分屏
```

### 技巧 2：使用 VS Code 集成终端

在 VS Code 中：
1. 按 `` Ctrl+` `` 打开终端
2. 点击 `+` 旁边的下拉菜单 → 分割终端
3. 一个运行前端，一个运行后端

### 技巧 3：创建任务配置

在 `.vscode/tasks.json` 中已经配置好了：

```json
{
    "label": "启动 Frontend (Air 热重载)",
    "type": "shell",
    "command": "air",
    "args": ["-c", ".air-frontend.toml"]
}
```

使用方法：
1. 按 `Ctrl+Shift+P`
2. 输入 "Run Task"
3. 选择 "启动 Frontend (Air 热重载)"

### 技巧 4：检查服务状态

创建 `check-services.ps1`：

```powershell
$frontend = netstat -ano | findstr :8080
$backend = netstat -ano | findstr :8081

Write-Host "`n=== Service Status ===" -ForegroundColor Cyan

if ($frontend) {
    Write-Host "✅ Frontend (8080): Running" -ForegroundColor Green
    Write-Host $frontend
} else {
    Write-Host "❌ Frontend (8080): Not Running" -ForegroundColor Red
}

if ($backend) {
    Write-Host "✅ Backend (8081): Running" -ForegroundColor Green
    Write-Host $backend
} else {
    Write-Host "❌ Backend (8081): Not Running" -ForegroundColor Red
}
```

---

## 📊 Windows vs Linux/Mac 对比

| 项目 | Windows | Linux/Mac |
|------|---------|-----------|
| 可执行文件扩展名 | `.exe` 必需 | 无扩展名 |
| 环境变量设置 | `$env:VAR="value"` | `export VAR=value` |
| 端口检查 | `netstat -ano` | `lsof -i` |
| 进程终止 | `taskkill /F /PID` | `kill -9` |
| 路径分隔符 | `\` 或 `/` | `/` |

---

## 🔗 相关资源

- **Air Windows Issues**: https://github.com/air-verse/air/issues?q=windows
- **PowerShell 文档**: https://docs.microsoft.com/powershell/
- **Windows Terminal**: https://github.com/microsoft/terminal

---

## ✅ 检查清单

在 Windows 上使用 Air 前，确保：

- [ ] Go 已安装并在 PATH 中
- [ ] Air 已安装 (`go install github.com/air-verse/air@latest`)
- [ ] GOPATH/bin 在 PATH 中
- [ ] 配置文件使用 `.exe` 扩展名
- [ ] 端口 8080 和 8081 未被占用
- [ ] 防火墙允许这些端口
- [ ] 在项目根目录运行命令
- [ ] 依赖已安装 (`go mod download`)

---

**现在可以正常使用 Air 了！** 🎉

**更新时间**: 2024  
**维护**: trx-project 团队

