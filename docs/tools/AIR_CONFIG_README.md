# Air 热重载配置说明

本项目使用 Air 实现 Go 应用的热重载开发。由于项目采用前后端分离架构，提供了两个独立的 Air 配置文件。

## 📁 配置文件

### 1. `.air.toml` - Frontend 前端服务配置
用于开发前端服务（端口 8080）

```bash
# 启动前端热重载
air

# 或使用 Makefile
make dev-frontend
```

### 2. `.air-backend.toml` - Backend 后端服务配置
用于开发后端服务（端口 8081）

```bash
# 启动后端热重载
air -c .air-backend.toml

# 或使用 Makefile
make dev-backend
```

## 🔧 配置详解

### 关键配置项

#### `[build]` 构建配置

```toml
bin = "./tmp/frontend"                    # 编译后的二进制文件路径
cmd = "go build -o ./tmp/frontend ./cmd/frontend"  # 构建命令
delay = 1000                              # 文件变化后延迟构建时间（毫秒）
```

#### `exclude_dir` 排除目录
- **Frontend**: 排除 `cmd/backend`, `cmd/migrate` - 避免后端代码变化触发前端重载
- **Backend**: 排除 `cmd/frontend`, `cmd/migrate` - 避免前端代码变化触发后端重载
- **通用**: 排除 `tmp`, `vendor`, `testdata` - 避免临时文件触发重载

#### `exclude_regex` 排除文件正则
```toml
exclude_regex = ["_test.go"]  # 排除测试文件，测试文件变化不触发重载
```

#### `include_ext` 监听的文件扩展名
```toml
include_ext = ["go", "tpl", "tmpl", "html", "yaml"]
```
监听 Go 源码、模板文件和配置文件的变化

#### `full_bin` 完整运行命令
```toml
full_bin = "GO_ENV=dev AUTO_MIGRATE=false ./tmp/frontend"
```
- `GO_ENV=dev`: 设置开发环境，加载 `config.dev.yaml`
- `AUTO_MIGRATE=false`: 开发时不自动运行数据库迁移（避免频繁迁移）

#### `stop_on_error` 错误时停止
```toml
stop_on_error = true  # 编译错误时停止运行，方便调试
```

## 🚀 使用方法

### 方式一：直接使用 Air 命令

```bash
# 1. 安装 Air
go install github.com/cosmtrek/air@latest

# 2. 启动前端服务
air

# 3. 启动后端服务（新终端）
air -c .air-backend.toml
```

### 方式二：使用 Makefile（推荐）

```bash
# 启动前端服务
make dev-frontend

# 启动后端服务（新终端）
make dev-backend

# 同时启动两个服务（需要并行工具）
make dev  # 如果 Makefile 有配置
```

## 📝 开发工作流

### 1. 启动开发环境

```bash
# 终端 1: 启动前端服务
make dev-frontend

# 终端 2: 启动后端服务
make dev-backend
```

### 2. 修改代码
- 保存文件后，Air 会自动检测变化
- 自动重新编译并重启服务
- 在终端查看构建日志和运行日志

### 3. 查看效果
- 前端服务: http://localhost:8080
- 后端服务: http://localhost:8081
- Swagger 文档: 
  - http://localhost:8080/swagger/index.html
  - http://localhost:8081/swagger/index.html

## 🎨 日志颜色说明

```toml
[color]
  app = ""           # 应用输出
  build = "yellow"   # 构建信息（黄色）
  main = "magenta"   # 主要信息（品红色）
  runner = "green"   # 运行信息（绿色）
  watcher = "cyan"   # 监视信息（青色）
```

## 🐛 常见问题

### 1. Air 命令未找到
```bash
# 安装 Air
go install github.com/cosmtrek/air@latest

# 确保 GOPATH/bin 在 PATH 中
export PATH=$PATH:$(go env GOPATH)/bin
```

### 2. 端口已被占用
```bash
# 检查端口占用
# Windows
netstat -ano | findstr :8080
netstat -ano | findstr :8081

# Linux/Mac
lsof -i :8080
lsof -i :8081

# 杀死进程或修改配置文件中的端口
```

### 3. 文件变化不触发重载
- 检查文件扩展名是否在 `include_ext` 中
- 检查文件路径是否在 `exclude_dir` 中
- 尝试重启 Air

### 4. 编译错误后无法恢复
- 修复代码错误后，Air 会自动重试
- 如果卡住，按 `Ctrl+C` 停止后重新启动

## ⚙️ 高级配置

### 自定义环境变量
修改 `full_bin` 配置：

```toml
# 使用测试环境
full_bin = "GO_ENV=test AUTO_MIGRATE=false ./tmp/frontend"

# 启用自动迁移（首次运行）
full_bin = "GO_ENV=dev AUTO_MIGRATE=true ./tmp/frontend"

# 添加其他环境变量
full_bin = "GO_ENV=dev AUTO_MIGRATE=false DEBUG=true ./tmp/frontend"
```

### 调整重载延迟
如果文件变化过于频繁：

```toml
[build]
  delay = 2000  # 增加到 2 秒
  rerun_delay = 1000  # 重新运行延迟
```

### 启用轮询模式
在某些文件系统上（如网络文件系统）：

```toml
[build]
  poll = true
  poll_interval = 1000  # 毫秒
```

## 📚 参考资源

- [Air 官方文档](https://github.com/cosmtrek/air)
- [Air 配置示例](https://github.com/cosmtrek/air/blob/master/air_example.toml)

## 💡 最佳实践

1. **分离前后端开发**：使用两个终端分别运行前后端服务
2. **关闭自动迁移**：开发时设置 `AUTO_MIGRATE=false`，避免频繁迁移
3. **监听配置文件**：`yaml` 文件变化也会触发重载
4. **使用 Makefile**：封装复杂的启动命令
5. **错误时停止**：设置 `stop_on_error = true`，便于调试

---

**更新时间**: 2024
**维护**: trx-project 团队

