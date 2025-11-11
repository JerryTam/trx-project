# TRX Project

> 基于 Gin 框架的现代化 Go Web 服务 - **前后台完全分离架构**

## ✨ 架构特点

### 🎯 完全分离的前后台服务

```
前台服务 (Frontend)              后台服务 (Backend)
├── 端口: 8080                   ├── 端口: 8081
├── 面向: 最终用户                ├── 面向: 管理员
├── 认证: JWT (user 角色)         ├── 认证: JWT (admin/superadmin 角色)
├── Token: 7天有效期              ├── Token: 1天有效期
└── 入口: cmd/frontend/          └── 入口: cmd/backend/
```

### 🛠️ 技术栈

- **Web 框架**: [Gin](https://gin-gonic.com/)
- **日志**: [Uber Zap](https://github.com/uber-go/zap)
- **依赖注入**: [Google Wire](https://github.com/google/wire)
- **数据库**: MySQL + [GORM](https://gorm.io/)
- **缓存**: Redis
- **消息队列**: Kafka
- **认证**: JWT (golang-jwt/jwt)

## 📁 项目结构

```
trx-project/
├── cmd/                          # 程序入口
│   ├── frontend/                 # 前台服务 ⭐
│   │   ├── main.go
│   │   ├── wire.go
│   │   └── wire_gen.go
│   └── backend/                  # 后台服务 ⭐
│       ├── main.go
│       ├── wire.go
│       └── wire_gen.go
├── internal/                     # 内部代码
│   ├── api/
│   │   ├── handler/             # HTTP 处理器
│   │   ├── middleware/          # 中间件 (JWT 认证)
│   │   └── router/
│   │       ├── frontend.go      # 前台路由 ⭐
│   │       └── backend.go       # 后台路由 ⭐
│   ├── model/                   # 数据模型
│   ├── repository/              # 数据访问层
│   └── service/                 # 业务逻辑层
├── pkg/                         # 公共包
│   ├── cache/                   # Redis 封装
│   ├── config/                  # 配置管理
│   ├── database/                # 数据库初始化
│   ├── jwt/                     # JWT 认证 ⭐
│   ├── kafka/                   # Kafka 封装
│   ├── logger/                  # 日志封装
│   └── response/                # 统一响应格式
├── config/                      # 配置文件
│   ├── config.yaml
│   └── config.yaml.example
├── scripts/                     # 脚本
│   ├── generate_admin_token.go  # 生成管理员 Token ⭐
│   ├── test_frontend.sh         # 前台 API 测试 ⭐
│   ├── test_backend.sh          # 后台 API 测试 ⭐
│   ├── init_db.sql
│   └── setup.sh
├── docs/                        # 文档
├── docker-compose.yml           # Docker 编排
├── Makefile                     # 构建脚本
└── README.md
```

## 🚀 快速开始

### 1. 安装依赖

```bash
# 安装 Go 依赖
make deps

# 安装 Wire (如果还没有)
go install github.com/google/wire/cmd/wire@latest

# 安装 Air (热重载，可选)
go install github.com/cosmtrek/air@latest
```

### 2. 启动基础服务

```bash
# 启动 MySQL, Redis, Kafka
make docker-up
```

### 3. 配置

```bash
# 复制配置文件
cp config/config.yaml.example config/config.yaml

# 编辑配置（特别是 JWT secret）
vim config/config.yaml
```

### 4. 生成 Wire 代码

```bash
# 生成前后台的依赖注入代码
make wire
```

### 5. 构建服务

```bash
# 构建所有服务
make build

# 或单独构建
make build-frontend  # 前台
make build-backend   # 后台
```

### 6. 运行服务

#### 方式1：生产模式

```bash
# 运行前台服务（端口 8080）
./bin/frontend

# 运行后台服务（端口 8081）
./bin/backend
```

#### 方式2：开发模式（热重载）

```bash
# 终端1：运行前台
make dev-frontend

# 终端2：运行后台
make dev-backend
```

## 🔐 JWT 认证

### 生成管理员 Token

```bash
go run scripts/generate_admin_token.go
```

输出示例：
```
✅ 管理员 Token 生成成功!

Token 信息:
  User ID:  1
  Username: admin
  Role:     superadmin
  Expires:  1 day

Token:
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

使用示例:
  export ADMIN_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  curl -H "Authorization: Bearer $ADMIN_TOKEN" http://localhost:8081/api/v1/admin/users
```

### Token 类型

| 类型 | 角色 | 有效期 | 用途 |
|------|------|--------|------|
| 用户 Token | `user` | 7天 | 前台用户认证 |
| 管理员 Token | `admin` / `superadmin` | 1天 | 后台管理认证 |

## 📡 API 接口

### 前台服务 (Port 8080)

```bash
# 健康检查
GET http://localhost:8080/health

# 用户注册
POST http://localhost:8080/api/v1/public/register
{
  "username": "testuser",
  "email": "test@example.com",
  "password": "password123"
}

# 用户登录
POST http://localhost:8080/api/v1/public/login
{
  "username": "testuser",
  "password": "password123"
}

# 获取个人信息（需要认证）
GET http://localhost:8080/api/v1/user/profile
Authorization: Bearer <user_token>

# 更新个人信息（需要认证）
PUT http://localhost:8080/api/v1/user/profile
Authorization: Bearer <user_token>
```

### 后台服务 (Port 8081)

```bash
# 健康检查
GET http://localhost:8081/health

# 获取用户列表（需要管理员权限）
GET http://localhost:8081/api/v1/admin/users
Authorization: Bearer <admin_token>

# 获取用户详情
GET http://localhost:8081/api/v1/admin/users/:id
Authorization: Bearer <admin_token>

# 更新用户状态
PUT http://localhost:8081/api/v1/admin/users/:id/status
Authorization: Bearer <admin_token>

# 删除用户
DELETE http://localhost:8081/api/v1/admin/users/:id
Authorization: Bearer <admin_token>

# 重置密码
POST http://localhost:8081/api/v1/admin/users/:id/reset-password
Authorization: Bearer <admin_token>

# 获取统计信息
GET http://localhost:8081/api/v1/admin/statistics/users
Authorization: Bearer <admin_token>
```

## 🧪 测试

### 自动化测试脚本

```bash
# 测试前台服务
./scripts/test_frontend.sh

# 测试后台服务
./scripts/test_backend.sh

# 单元测试
make test
```

### 手动测试

```bash
# 1. 注册用户并获取 Token
curl -X POST http://localhost:8080/api/v1/public/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'

# 2. 使用 Token 访问
USER_TOKEN="<从注册响应中获取>"
curl http://localhost:8080/api/v1/user/profile \
  -H "Authorization: Bearer $USER_TOKEN"

# 3. 生成管理员 Token
ADMIN_TOKEN=$(go run scripts/generate_admin_token.go | grep "eyJ" | tr -d '\n')

# 4. 使用管理员 Token 访问后台
curl http://localhost:8081/api/v1/admin/users \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

## 🛠️ Makefile 命令

```bash
make help           # 显示所有可用命令
make deps           # 安装依赖
make wire           # 生成 Wire 代码
make build          # 构建所有服务
make build-frontend # 构建前台服务
make build-backend  # 构建后台服务
make run-frontend   # 运行前台服务
make run-backend    # 运行后台服务
make dev-frontend   # 开发模式（热重载）
make dev-backend    # 开发模式（热重载）
make test           # 运行测试
make clean          # 清理构建文件
make docker-up      # 启动 Docker 服务
make docker-down    # 停止 Docker 服务
```

## 📊 统一响应格式

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "user": {
      "id": 1,
      "username": "testuser",
      "email": "test@example.com"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  },
  "timestamp": 1699999999
}
```

详见: [docs/RESPONSE_FORMAT.md](docs/RESPONSE_FORMAT.md)

## 📚 文档

- [API 文档](docs/API.md)
- [架构文档](docs/ARCHITECTURE.md)
- [前后台分离说明](SEPARATED_SERVICES.md)
- [部署指南](docs/DEPLOYMENT.md)
- [Kafka 使用指南](docs/KAFKA_USAGE.md)
- [响应格式说明](docs/RESPONSE_FORMAT.md)

## 🔄 开发流程

### 添加新的前台接口

1. 在 `internal/api/handler/user_handler.go` 添加处理器
2. 在 `internal/api/router/frontend.go` 注册路由
3. 重新构建: `make build-frontend`

### 添加新的后台接口

1. 创建 handler: `internal/api/handler/admin_xxx_handler.go`
2. 在 `internal/api/router/backend.go` 注册路由
3. 在 `cmd/backend/wire.go` 添加依赖注入
4. 重新生成 Wire 代码: `cd cmd/backend && wire`
5. 重新构建: `make build-backend`

### 修改依赖注入

```bash
# 修改 cmd/frontend/wire.go 或 cmd/backend/wire.go 后
cd cmd/frontend && wire  # 或 cd cmd/backend && wire
cd ../..
make build
```

## 🐳 Docker 部署

```bash
# 启动所有服务（MySQL, Redis, Kafka）
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

## 🎯 架构优势

### ✅ 完全分离
- 前后台独立部署
- 互不影响
- 独立扩展

### ✅ 安全性
- 不同端口
- 不同认证方式
- 后台接口不对外暴露

### ✅ 可维护性
- 代码分离清晰
- 独立开发和测试
- 独立升级

### ✅ 可扩展性
- 可部署多个实例
- 可添加负载均衡
- 微服务化准备

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可

MIT License

## 📮 联系方式

- 作者: Your Name
- Email: your.email@example.com
- GitHub: https://github.com/yourusername/trx-project

---

⭐ 如果这个项目对你有帮助，请给个 Star！
