# TRX Project 启动指南

本项目提供了跨平台的启动脚本，支持 **Windows**、**Linux** 和 **macOS**。

## 📋 目录

- [Windows 启动](#windows-启动)
- [Linux/macOS 启动](#linuxmacos-启动)
- [配置文件说明](#配置文件说明)
- [常见问题](#常见问题)

---

## 🪟 Windows 启动

### 前置要求

1. 安装 [Go](https://golang.org/dl/) (推荐 1.21+)
2. 安装 [Air](https://github.com/air-verse/air)
   ```powershell
   go install github.com/air-verse/air@latest
   ```
3. 安装 Git Bash 或 PowerShell

### 启动方式

#### 方式 1: PowerShell 脚本（推荐）

```powershell
# 启动前端服务（端口 8080）
.\start-frontend.ps1

# 启动后端服务（端口 8081）
.\start-backend.ps1

# 同时启动前后端
.\start-dev.ps1
```

#### 方式 2: Bash 脚本（Git Bash）

```bash
# 启动前端服务
./start-frontend.sh

# 启动后端服务
./start-backend.sh

# 同时启动前后端
./start-dev.sh
```

#### 方式 3: 直接使用 Air

```powershell
# 设置环境变量
$env:GO_ENV = "dev"
$env:AUTO_MIGRATE = "false"

# 启动前端
air -c .air-frontend.toml

# 或启动后端
air -c .air-backend.toml
```

---

## 🐧 Linux/macOS 启动

### 前置要求

1. 安装 [Go](https://golang.org/dl/) (推荐 1.21+)
   ```bash
   # Ubuntu/Debian
   sudo apt install golang-go
   
   # macOS
   brew install go
   ```

2. 安装 [Air](https://github.com/air-verse/air)
   ```bash
   go install github.com/air-verse/air@latest
   ```

3. （可选）安装 [tmux](https://github.com/tmux/tmux) - 用于分屏显示
   ```bash
   # Ubuntu/Debian
   sudo apt install tmux
   
   # CentOS/RHEL
   sudo yum install tmux
   
   # macOS
   brew install tmux
   ```

### 启动方式

#### 方式 1: Bash 脚本（推荐）

```bash
# 确保脚本有执行权限
chmod +x start-frontend.sh start-backend.sh start-dev.sh

# 启动前端服务（端口 8080）
./start-frontend.sh

# 启动后端服务（端口 8081）
./start-backend.sh

# 同时启动前后端（提供多种方式选择）
./start-dev.sh
```

#### 方式 2: 使用 tmux（推荐用于开发）

运行 `./start-dev.sh` 并选择选项 `1`：

```
启动方式:
1. 使用 tmux 在分屏中启动（推荐，Linux/macOS）
2. 在后台启动两个服务
3. 依次启动（前台模式，仅启动 Frontend）

请选择 (1/2/3): 1
```

**tmux 常用快捷键：**
- `Ctrl+B` 然后按 `←/→` : 切换面板
- `Ctrl+B` 然后按 `D` : 退出 tmux（服务继续运行）
- `Ctrl+B` 然后按 `&` : 关闭当前窗口

重新连接 tmux 会话：
```bash
tmux attach -t trx-dev
```

#### 方式 3: 后台启动

运行 `./start-dev.sh` 并选择选项 `2`。

查看日志：
```bash
tail -f tmp/frontend.log
tail -f tmp/backend.log
```

停止服务：
```bash
./stop-dev.sh
```

#### 方式 4: 直接使用 Air

```bash
# 设置环境变量
export GO_ENV=dev
export AUTO_MIGRATE=false

# 启动前端
air -c .air-frontend-linux.toml

# 或启动后端
air -c .air-backend-linux.toml
```

---

## 📁 配置文件说明

项目包含多个 Air 配置文件以支持不同平台：

### Windows 配置文件

| 文件 | 用途 | 可执行文件 |
|------|------|-----------|
| `.air-frontend.toml` | Windows 前端配置 | `tmp/frontend.exe` |
| `.air-backend.toml` | Windows 后端配置 | `tmp/backend.exe` |

### Linux/macOS 配置文件

| 文件 | 用途 | 可执行文件 |
|------|------|-----------|
| `.air-frontend-linux.toml` | Linux/macOS 前端配置 | `tmp/frontend` |
| `.air-backend-linux.toml` | Linux/macOS 后端配置 | `tmp/backend` |

### 主要区别

```toml
# Windows 配置
[build]
  bin = "./tmp/frontend.exe"           # 带 .exe 扩展名
  cmd = "go build -o ./tmp/frontend.exe ./cmd/frontend"

# Linux/macOS 配置
[build]
  bin = "./tmp/frontend"               # 无扩展名
  cmd = "go build -o ./tmp/frontend ./cmd/frontend"
```

### 自动选择

启动脚本会自动检测操作系统并选择正确的配置文件：

```bash
# start-frontend.sh 中的自动检测逻辑
detect_os() {
    case "$(uname -s)" in
        Linux*)     echo "linux";;
        Darwin*)    echo "linux";;  # macOS 使用 Linux 配置
        CYGWIN*|MINGW*|MSYS*)    echo "windows";;
        *)          echo "linux";;  # 默认使用 Linux 配置
    esac
}

OS=$(detect_os)
if [ "$OS" = "windows" ]; then
    AIR_CONFIG=".air-frontend.toml"
else
    AIR_CONFIG=".air-frontend-linux.toml"
fi
```

---

## 🔧 环境变量

所有启动脚本都会设置以下环境变量：

| 变量 | 值 | 说明 |
|------|----|----|
| `GO_ENV` | `dev` | 运行环境（dev/prod） |
| `AUTO_MIGRATE` | `false` | 是否自动执行数据库迁移 |

如需修改，可以编辑启动脚本或在命令行中手动设置：

```bash
# Linux/macOS
export GO_ENV=prod
export AUTO_MIGRATE=true

# Windows PowerShell
$env:GO_ENV = "prod"
$env:AUTO_MIGRATE = "true"
```

---

## 🌐 服务地址

启动成功后，可以通过以下地址访问：

### Frontend (端口 8080)

| 地址 | 说明 |
|------|------|
| http://localhost:8080 | 主页 |
| http://localhost:8080/health | 健康检查 |
| http://localhost:8080/swagger/index.html | Swagger API 文档 |
| http://localhost:8080/metrics | Prometheus 指标 |

### Backend (端口 8081)

| 地址 | 说明 |
|------|------|
| http://localhost:8081 | 主页 |
| http://localhost:8081/health | 健康检查 |
| http://localhost:8081/swagger/index.html | Swagger API 文档 |
| http://localhost:8081/metrics | Prometheus 指标 |

---

## ❓ 常见问题

### 1. Air 命令找不到

**错误信息：**
```
bash: air: command not found
```

**解决方法：**
```bash
# 安装 Air
go install github.com/air-verse/air@latest

# 确保 $GOPATH/bin 在 PATH 中
export PATH=$PATH:$(go env GOPATH)/bin

# 验证安装
air -v
```

### 2. 端口已被占用

**错误信息：**
```
bind: address already in use
```

**解决方法：**

**Linux/macOS:**
```bash
# 查找占用端口的进程
lsof -i :8080
lsof -i :8081

# 停止进程
kill <PID>
```

**Windows:**
```powershell
# 查找占用端口的进程
netstat -ano | findstr :8080
netstat -ano | findstr :8081

# 停止进程
taskkill /PID <PID> /F
```

### 3. Windows PowerShell 执行策略限制

**错误信息：**
```
无法加载文件 start-frontend.ps1，因为在此系统上禁止运行脚本。
```

**解决方法：**
```powershell
# 以管理员身份运行 PowerShell
Set-ExecutionPolicy RemoteSigned -Scope CurrentUser

# 或使用 Git Bash 运行 .sh 脚本
./start-frontend.sh
```

### 4. Linux 下没有执行权限

**错误信息：**
```
bash: ./start-frontend.sh: Permission denied
```

**解决方法：**
```bash
# 添加执行权限
chmod +x start-frontend.sh start-backend.sh start-dev.sh stop-dev.sh

# 然后运行
./start-frontend.sh
```

### 5. 数据库连接失败

**错误信息：**
```
Error 1045: Access denied for user 'root'@'localhost'
```

**解决方法：**
1. 检查 `config/config.yaml` 中的数据库配置
2. 确保 MySQL 服务正在运行
3. 验证数据库用户名和密码

### 6. Redis 连接失败

**错误信息：**
```
dial tcp [::1]:6379: connect: connection refused
```

**解决方法：**
1. 确保 Redis 服务正在运行
   ```bash
   # Linux
   sudo systemctl start redis
   
   # macOS
   brew services start redis
   
   # Windows
   redis-server
   ```
2. 检查 `config/config.yaml` 中的 Redis 配置

### 7. 热重载不工作

**可能原因：**
- 文件监控失败
- 排除目录配置不当

**解决方法：**
1. 检查 `.air-*.toml` 配置文件
2. 确保 `tmp` 目录存在
3. 清理并重新编译：
   ```bash
   rm -rf tmp
   mkdir tmp
   ./start-frontend.sh
   ```

---

## 🎯 最佳实践

### 开发环境

1. **使用 tmux（Linux/macOS）**
   - 优点：分屏显示，方便查看前后端日志
   - 启动：`./start-dev.sh` → 选择选项 `1`

2. **使用 PowerShell 脚本（Windows）**
   - 优点：彩色输出，自动清理，友好提示
   - 启动：`.\start-frontend.ps1` 和 `.\start-backend.ps1` 在不同终端

### 生产环境

1. 设置环境变量为生产模式：
   ```bash
   export GO_ENV=prod
   export AUTO_MIGRATE=false
   ```

2. 使用编译后的二进制文件，而不是 Air：
   ```bash
   go build -o frontend ./cmd/frontend
   go build -o backend ./cmd/backend
   
   ./frontend &
   ./backend &
   ```

3. 使用进程管理器（如 systemd、supervisor）管理服务

---

## 📚 相关文档

- [Air 配置说明](./AIR_CONFIG_README.md)
- [Air 使用指南](./AIR_USAGE_GUIDE.md)
- [Air Windows 指南](./AIR_WINDOWS_GUIDE.md)

---

## 🤝 贡献

如果您发现任何问题或有改进建议，欢迎提交 Issue 或 Pull Request。

---

## 📝 许可证

本项目采用 MIT 许可证。

