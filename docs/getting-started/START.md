# 🚀 快速启动指南

## 第一次运行（完整步骤）

### 1. 准备环境

```bash
# 确保已安装：
# - Go 1.21+
# - Docker & Docker Compose
# - Make

# 安装 Wire
go install github.com/google/wire/cmd/wire@latest

# 安装 Air（可选，用于热重载）
go install github.com/cosmtrek/air@latest
```

### 2. 安装依赖

```bash
make deps
```

### 3. 启动基础服务

```bash
# 启动 MySQL, Redis, Kafka
make docker-up

# 等待服务启动（约10秒）
sleep 10

# 验证服务状态
docker-compose ps
```

### 4. 配置文件

```bash
# 复制配置示例
cp config/config.yaml.example config/config.yaml

# 如果需要，修改配置
# vim config/config.yaml
```

### 5. 生成 Wire 代码

```bash
make wire
```

### 6. 构建服务

```bash
make build
```

### 7. 运行服务

#### 方式1：直接运行（推荐用于测试）

**终端1 - 前台服务**:
```bash
./bin/frontend
```

**终端2 - 后台服务**:
```bash
./bin/backend
```

#### 方式2：开发模式（推荐用于开发）

**终端1 - 前台服务**:
```bash
make dev-frontend
```

**终端2 - 后台服务**:
```bash
make dev-backend
```

## 验证服务

### 检查健康状态

```bash
# 前台服务
curl http://localhost:8080/health

# 后台服务
curl http://localhost:8081/health
```

应该返回：
```json
{
  "status": "ok",
  "service": "frontend"  // 或 "backend"
}
```

### 测试前台服务

```bash
# 注册用户
curl -X POST http://localhost:8080/api/v1/public/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'

# 会返回用户信息和 Token
```

### 测试后台服务

```bash
# 1. 生成管理员 Token
go run scripts/generate_admin_token.go

# 2. 使用 Token 测试（替换 <token>）
curl http://localhost:8081/api/v1/admin/users \
  -H "Authorization: Bearer <token>"
```

## 自动化测试

```bash
# 启动 Docker 服务
make docker-up

# 等待服务就绪
sleep 10

# 启动前台服务（后台运行）
./bin/frontend &
FRONTEND_PID=$!

# 启动后台服务（后台运行）
./bin/backend &
BACKEND_PID=$!

# 等待服务启动
sleep 3

# 运行前台测试
./scripts/test_frontend.sh

# 运行后台测试
./scripts/test_backend.sh

# 停止服务
kill $FRONTEND_PID $BACKEND_PID
```

## 常见问题

### Q: Wire 生成失败？

```bash
# 确保已安装 Wire
go install github.com/google/wire/cmd/wire@latest

# 清理后重新生成
rm cmd/frontend/wire_gen.go cmd/backend/wire_gen.go
make wire
```

### Q: 构建失败？

```bash
# 清理并重新构建
make clean
make deps
make wire
make build
```

### Q: Docker 服务启动失败？

```bash
# 停止并清理
make docker-down

# 重新启动
make docker-up

# 查看日志
docker-compose logs
```

### Q: 端口被占用？

```bash
# 检查端口占用
# Windows:
netstat -ano | findstr "8080"
netstat -ano | findstr "8081"

# Linux/Mac:
lsof -i :8080
lsof -i :8081

# 修改 config/config.yaml 中的端口，或关闭占用端口的进程
```

### Q: 数据库连接失败？

```bash
# 1. 确保 Docker 服务已启动
docker-compose ps

# 2. 检查 MySQL 状态
docker-compose logs mysql

# 3. 等待 MySQL 完全启动（可能需要30秒）
sleep 30

# 4. 测试连接
docker-compose exec mysql mysql -uroot -proot123456 -e "SELECT 1"
```

## 开发流程

### 修改代码后

```bash
# 如果使用 Air (热重载)，无需手动操作，代码会自动重新编译和重启

# 如果不使用 Air
make build
./bin/frontend   # 或 ./bin/backend
```

### 修改依赖注入

```bash
# 修改 wire.go 后
make wire
make build
```

### 添加新的依赖包

```bash
go get <package>
go mod tidy
```

## 停止服务

### 停止应用

```bash
# Ctrl+C 停止正在运行的服务

# 如果后台运行
pkill frontend
pkill backend
```

### 停止 Docker 服务

```bash
make docker-down
```

## 清理

```bash
# 清理构建文件
make clean

# 清理 Docker 数据（谨慎！会删除所有数据）
docker-compose down -v
```

## 生产环境部署

参见 [docs/DEPLOYMENT.md](docs/DEPLOYMENT.md)

## 下一步

- 阅读 [README.md](README.md) 了解项目详情
- 查看 [SEPARATED_SERVICES.md](SEPARATED_SERVICES.md) 了解架构设计
- 参考 [docs/API.md](docs/API.md) 查看完整 API 文档
- 运行测试：`make test`

---

🎉 **恭喜！你已经成功启动了服务！**

前台服务: http://localhost:8080  
后台服务: http://localhost:8081

