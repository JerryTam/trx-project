# TRX Project - 快速开始

## 🚀 一键启动

### Windows

```powershell
# PowerShell
.\start-frontend.ps1   # 启动前端（端口 8080）
.\start-backend.ps1    # 启动后端（端口 8081）
.\start-dev.ps1        # 同时启动前后端

# 或使用 Git Bash
./start-frontend.sh
./start-backend.sh
./start-dev.sh
```

### Linux/macOS

```bash
# 首次运行，添加执行权限
chmod +x *.sh

# 启动服务
./start-frontend.sh    # 启动前端（端口 8080）
./start-backend.sh     # 启动后端（端口 8081）
./start-dev.sh         # 同时启动前后端（支持 tmux 分屏）
```

## 📋 前置要求

1. **Go 1.21+**
   ```bash
   go version
   ```

2. **Air 热重载工具**
   ```bash
   go install github.com/air-verse/air@latest
   ```

3. **MySQL** - 数据库
   
4. **Redis** - 缓存和限流

5. **(可选) tmux** - Linux/macOS 分屏显示
   ```bash
   # Ubuntu/Debian
   sudo apt install tmux
   
   # macOS
   brew install tmux
   ```

## 🌐 访问地址

启动成功后访问：

### Frontend (8080)
- 主页: http://localhost:8080
- 健康检查: http://localhost:8080/health
- API 文档: http://localhost:8080/swagger/index.html

### Backend (8081)
- 主页: http://localhost:8081
- 健康检查: http://localhost:8081/health
- API 文档: http://localhost:8081/swagger/index.html

## 📖 详细文档

查看 [完整启动指南](docs/STARTUP_GUIDE.md) 了解更多选项和配置。

## 🔧 配置

编辑 `config/config.yaml` 修改数据库、Redis 等配置：

```yaml
server:
  name: "trx-project"
  host: "0.0.0.0"
  port: 8080

database:
  host: "localhost"
  port: 3306
  user: "root"
  password: "your_password"
  dbname: "trx_db"

redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0
```

## ⚠️ 常见问题

### Air 命令找不到

```bash
go install github.com/air-verse/air@latest
export PATH=$PATH:$(go env GOPATH)/bin
```

### 端口已被占用

**Linux/macOS:**
```bash
lsof -i :8080    # 查找占用进程
kill <PID>       # 停止进程
```

**Windows:**
```powershell
netstat -ano | findstr :8080
taskkill /PID <PID> /F
```

### 停止服务

**Linux/macOS:**
```bash
./stop-dev.sh    # 停止所有服务
```

**Windows/所有平台:**
```
按 Ctrl+C 停止前台运行的服务
```

## 🎯 开发建议

### Windows 开发者
- 推荐使用 PowerShell 脚本（彩色输出，自动清理）
- 在不同终端窗口分别启动前后端

### Linux/macOS 开发者
- 推荐使用 `./start-dev.sh` → 选择 tmux 模式
- 可以在分屏中同时查看前后端日志
- `tmux attach -t trx-dev` 重新连接会话

## 📚 项目结构

```
trx-project/
├── cmd/                    # 入口文件
│   ├── frontend/          # 前端服务
│   └── backend/           # 后端服务
├── internal/              # 内部包
│   ├── api/              # API 层
│   ├── model/            # 数据模型
│   ├── repository/       # 数据访问层
│   └── service/          # 业务逻辑层
├── pkg/                   # 公共包
│   ├── config/           # 配置管理
│   ├── database/         # 数据库
│   ├── cache/            # 缓存
│   ├── jwt/              # JWT 认证
│   └── logger/           # 日志
├── config/                # 配置文件
├── docs/                  # 文档
├── .air-*.toml           # Air 配置文件
├── start-*.sh            # Linux/macOS 启动脚本
├── start-*.ps1           # Windows 启动脚本
└── stop-dev.sh           # 停止脚本
```

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

MIT License

