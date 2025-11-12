# Air 热重载开发指南

## 🌪️ 什么是 Air？

**Air** 是一个轻量级的 Go 应用热重载工具，类似于 Node.js 的 nodemon。

### 作用对比

| 传统开发流程 | 使用 Air |
|------------|---------|
| 1. 修改代码 | 1. 修改代码 |
| 2. Ctrl+C 停止程序 | 2. **保存文件（自动完成以下步骤）** ✨ |
| 3. 手动运行 `go build` | |
| 4. 手动运行 `./app` | |
| 5. 查看效果 | 3. 查看效果 |

**节省时间：每次修改节省 5-10 秒！** ⚡

---

## 📦 安装 Air

### Windows / macOS / Linux

```bash
# 安装最新版本
go install github.com/air-verse/air@latest

# 验证安装
air -v
```

### 确认 GOPATH/bin 在 PATH 中

**Windows (PowerShell)**
```powershell
$env:Path += ";$(go env GOPATH)\bin"
# 或永久添加到系统环境变量
```

**Linux / macOS (Bash/Zsh)**
```bash
export PATH=$PATH:$(go env GOPATH)/bin
# 添加到 ~/.bashrc 或 ~/.zshrc 使其永久生效
```

---

## 🚀 快速开始

### 1️⃣ 启动前端服务（推荐方式）

```bash
# 方式 1: 使用 Makefile（推荐）
make dev-frontend

# 方式 2: 直接使用 Air
air -c .air-frontend.toml
```

访问：http://localhost:8080

### 2️⃣ 启动后端服务（新终端）

```bash
# 方式 1: 使用 Makefile（推荐）
make dev-backend

# 方式 2: 直接使用 Air
air -c .air-backend.toml
```

访问：http://localhost:8081

### 3️⃣ 修改代码并保存

- 保存任意 `.go` 文件
- Air 自动检测变化
- 自动重新编译
- 自动重启服务
- 在浏览器刷新查看效果

---

## 🎯 实际工作流程示例

### 场景：修改用户注册接口

```bash
# 1. 启动后端服务（终端 1）
make dev-backend

# 输出：
# > Backend server starting on 0.0.0.0:8081
```

```go
// 2. 修改代码 internal/api/handler/user_handler.go
func (h *UserHandler) Register(c *gin.Context) {
    // 添加新的验证逻辑
    if len(req.Password) < 8 {  // 修改：密码至少8位
        response.ValidateError(c, "密码至少需要8个字符")
        return
    }
    // ... 其他代码
}

// 3. 保存文件 (Ctrl+S)
```

```bash
# 4. Air 自动执行（无需手动操作）
# 终端 1 输出：
# > building...
# > running...
# > Backend server starting on 0.0.0.0:8081

# 5. 立即测试新接口
curl -X POST http://localhost:8081/api/v1/public/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","email":"test@example.com","password":"123"}'

# 6. 看到新的错误提示！
# {"code":400,"message":"密码至少需要8个字符"}
```

**整个过程无需停止服务，无需手动编译！** 🎉

---

## 📋 配置文件说明

### `.air-frontend.toml` - 前端服务配置

```toml
[build]
  # 编译输出路径
  bin = "./tmp/frontend"
  
  # 编译命令
  cmd = "go build -o ./tmp/frontend ./cmd/frontend"
  
  # 排除目录（后端代码不会触发前端重载）
  exclude_dir = ["cmd/backend", "cmd/migrate", "tmp", "vendor"]
  
  # 排除测试文件
  exclude_regex = ["_test.go"]
  
  # 监听的文件类型
  include_ext = ["go", "yaml"]
  
  # 运行命令（设置环境变量）
  full_bin = "GO_ENV=dev AUTO_MIGRATE=false ./tmp/frontend"
```

### `.air-backend.toml` - 后端服务配置

```toml
[build]
  bin = "./tmp/backend"
  cmd = "go build -o ./tmp/backend ./cmd/backend"
  exclude_dir = ["cmd/frontend", "cmd/migrate", "tmp", "vendor"]
  exclude_regex = ["_test.go"]
  include_ext = ["go", "yaml"]
  full_bin = "GO_ENV=dev AUTO_MIGRATE=false ./tmp/backend"
```

### 关键配置解释

| 配置项 | 说明 | 作用 |
|-------|------|-----|
| `exclude_dir` | 排除目录 | 前端配置排除后端目录，避免互相触发 |
| `exclude_regex` | 排除正则 | 测试文件改动不触发重载 |
| `include_ext` | 监听扩展名 | 监听 `.go` 和 `.yaml` 文件 |
| `full_bin` | 运行命令 | 设置环境变量（开发环境，不自动迁移） |

---

## 🎨 终端输出说明

### 正常运行输出

```bash
$ make dev-frontend

  __    _   ___  
 / /\  | | | |_) 
/_/--\ |_| |_| \_ v1.63.0, built with Go 1.25.4

watching .
watching cmd
watching cmd\frontend
!exclude tmp
watching internal
watching internal\api
...

building...
running...

[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] GET    /health                   --> main.SetupFrontend.func1 (9 handlers)
[GIN-debug] GET    /metrics                  --> github.com/gin-gonic/gin.WrapH.func1 (9 handlers)
[GIN-debug] POST   /api/v1/public/register   --> trx-project/internal/api/handler.(*UserHandler).Register-fm (9 handlers)

Frontend server starting on 0.0.0.0:8080
```

### 代码修改后输出

```bash
main.go has changed
building...
running...

Frontend server starting on 0.0.0.0:8080
```

### 编译错误输出

```bash
main.go has changed
building...

# command-line-arguments
cmd/frontend/main.go:45:2: undefined: someVariable

Build failed, watching...
```

Air 会等待您修复错误，修复后自动重试！

---

## 💡 使用技巧

### 1. 同时开发前后端

```bash
# 终端 1：前端
make dev-frontend

# 终端 2：后端（新终端窗口）
make dev-backend

# 终端 3：查看日志、测试接口等
curl http://localhost:8080/health
```

### 2. 临时修改环境变量

```bash
# 使用测试环境
GO_ENV=test air -c .air-frontend.toml

# 启用自动迁移（首次启动时）
AUTO_MIGRATE=true air -c .air-backend.toml
```

### 3. 调试编译问题

如果 Air 不工作，尝试手动编译排查：

```bash
# 手动编译前端
go build -o ./tmp/frontend ./cmd/frontend

# 手动编译后端
go build -o ./tmp/backend ./cmd/backend
```

### 4. 清理临时文件

```bash
# 清理编译产物
make clean

# 或手动删除
rm -rf tmp/
```

---

## 🐛 常见问题

### Q1: `air: command not found`

**原因**：Air 未安装或不在 PATH 中

**解决**：
```bash
# 1. 安装 Air
go install github.com/air-verse/air@latest

# 2. 添加到 PATH（Windows PowerShell）
$env:Path += ";$(go env GOPATH)\bin"

# 2. 添加到 PATH（Linux/Mac）
export PATH=$PATH:$(go env GOPATH)/bin
```

### Q2: 端口被占用

**错误**：`bind: address already in use`

**解决**：
```bash
# Windows: 查找占用端口的进程
netstat -ano | findstr :8080
taskkill /F /PID <进程ID>

# Linux/Mac: 查找并杀死进程
lsof -ti:8080 | xargs kill -9
```

### Q3: 修改代码不触发重载

**可能原因**：
1. 文件在排除目录中（`exclude_dir`）
2. 文件扩展名不在监听列表（`include_ext`）
3. 文件是测试文件（`_test.go`）

**解决**：
- 检查配置文件
- 确认文件路径
- 重启 Air

### Q4: 编译错误但 Air 卡住

**解决**：
```bash
# 1. Ctrl+C 停止 Air
# 2. 修复代码错误
# 3. 重新启动 Air
make dev-frontend
```

### Q5: 频繁重启（保存文件时触发多次）

**解决**：增加延迟时间
```toml
[build]
  delay = 2000  # 毫秒，增加到 2 秒
```

---

## 📊 Air vs 传统开发对比

### 开发效率提升

假设每天修改代码 50 次：

| 项目 | 传统开发 | 使用 Air | 节省 |
|------|---------|---------|------|
| 单次操作时间 | 10 秒 | 1 秒 | 9 秒 |
| 每天 50 次 | 500 秒 | 50 秒 | **450 秒** |
| 每周 (5天) | 2,500 秒 | 250 秒 | **37.5 分钟** |
| 每月 (20天) | 10,000 秒 | 1,000 秒 | **2.5 小时** |

**每月节省 2.5 小时！** ⏰

### 开发体验提升

✅ **专注代码编写** - 无需关心编译和重启  
✅ **即时反馈** - 保存即可看到效果  
✅ **减少错误** - 自动化减少人为失误  
✅ **提高士气** - 流畅的开发体验

---

## 🔗 相关资源

- **Air GitHub**: https://github.com/air-verse/air
- **Air 文档**: https://github.com/air-verse/air#readme
- **配置示例**: https://github.com/air-verse/air/blob/master/air_example.toml

---

## 📝 总结

### Air 的三大优势

1. **🚀 极速开发** - 自动化编译和重启流程
2. **🎯 精准控制** - 灵活的配置排除和包含规则
3. **💪 稳定可靠** - 成熟的工具，被广泛使用

### 推荐使用场景

✅ **本地开发** - 日常开发的最佳伴侣  
✅ **快速原型** - 快速验证想法  
✅ **API 调试** - 修改立即生效  
✅ **学习 Go** - 提升学习效率

### 不推荐场景

❌ **生产环境** - 使用编译后的二进制  
❌ **CI/CD** - 使用标准的构建流程  
❌ **性能测试** - 使用优化编译的版本

---

**享受流畅的 Go 开发体验！** 🎉

**更新时间**: 2024  
**维护**: trx-project 团队

