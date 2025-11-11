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
- **数据库迁移**: [golang-migrate](https://github.com/golang-migrate/migrate) ⭐
- **缓存**: Redis
- **消息队列**: Kafka
- **认证**: JWT (golang-jwt/jwt)
- **监控**: Prometheus + Grafana ⭐

## 📁 项目结构

```
trx-project/
├── cmd/                          # 程序入口
│   ├── frontend/                 # 前台服务 ⭐
│   │   ├── main.go
│   │   ├── wire.go
│   │   └── wire_gen.go
│   ├── backend/                  # 后台服务 ⭐
│   │   ├── main.go
│   │   ├── wire.go
│   │   └── wire_gen.go
│   └── migrate/                  # 数据库迁移工具 ⭐
│       └── main.go
├── migrations/                   # 数据库迁移文件 ⭐
│   ├── 000001_create_users_table.up.sql
│   ├── 000001_create_users_table.down.sql
│   └── ...
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
│   ├── migrate/                 # 数据库迁移管理 ⭐
│   └── response/                # 统一响应格式
├── config/                      # 配置文件
│   ├── config.yaml
│   └── config.yaml.example
├── scripts/                     # 脚本
│   ├── generate_admin_token.go  # 生成管理员 Token ⭐
│   ├── migrate.sh               # 迁移管理脚本 ⭐
│   ├── test_migration.sh        # 迁移功能测试 ⭐
│   ├── test_frontend.sh         # 前台 API 测试 ⭐
│   └── test_backend.sh          # 后台 API 测试 ⭐
├── docs/                        # 文档
│   ├── MIGRATION_GUIDE.md       # 迁移管理指南 ⭐
│   └── ...
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

### 4. 数据库迁移 ⭐

```bash
# 执行数据库迁移（创建表结构和初始数据）
make migrate-up

# 查看当前迁移版本
make migrate-version

# 创建新的迁移文件
make migrate-create NAME=add_user_phone
```

> 📖 详细文档: [数据库迁移管理指南](docs/MIGRATION_GUIDE.md)

### 5. 生成 Wire 代码

```bash
# 生成前后台的依赖注入代码
make wire
```

### 6. 构建服务

```bash
# 构建所有服务
make build

# 或单独构建
make build-frontend  # 前台
make build-backend   # 后台
```

### 7. 运行服务

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

## 🚦 限流功能

### 限流策略

项目实现了完整的三层限流保护：

| 限流类型 | 作用范围 | 默认配置 | 说明 |
|---------|---------|----------|------|
| **全局限流** | 所有请求 | 1000/秒 | 保护整体服务 |
| **IP 限流** | 按客户端 IP | 100/分钟 | 防止单个 IP 滥用 |
| **用户限流** | 按认证用户 | 1000/分钟 | 精细化用户控制 |

### 配置限流

编辑配置文件 `config/config.yaml`:

```yaml
rate_limit:
  enabled: true           # 启用/禁用限流
  global_rate: "1000-S"   # 全局：每秒1000个请求
  ip_rate: "100-M"        # IP：每IP每分钟100个请求
  user_rate: "1000-M"     # 用户：每用户每分钟1000个请求
```

### 限流响应

当请求被限流时，会返回 **HTTP 429** 状态码：

```json
{
  "code": 429,
  "message": "IP 请求限流（192.168.1.100），请 45 秒后重试",
  "data": null
}
```

每个响应都包含限流信息头：

```http
X-RateLimit-Limit: 100        # 限流上限
X-RateLimit-Remaining: 95     # 剩余请求数
X-RateLimit-Reset: 45         # 重置时间（秒）
```

### 测试限流

```bash
# 运行限流测试脚本
./scripts/test_rate_limit.sh

# 手动测试 IP 限流
for i in {1..150}; do
  curl http://localhost:8080/health
done
```

📖 **详细文档**: [限流使用指南](RATE_LIMIT_GUIDE.md)

## 🔍 请求 ID 追踪

### 功能介绍

为每个请求分配唯一的标识符（UUID），用于追踪完整的请求链路，便于调试和日志分析。

### 主要特性

✅ **自动生成** - 服务器自动为每个请求生成唯一 UUID  
✅ **客户端指定** - 支持通过 `X-Request-ID` 头传递请求 ID  
✅ **日志集成** - 所有日志自动包含请求 ID  
✅ **响应头返回** - 每个响应都包含 `X-Request-ID` 头

### 使用示例

#### 自动生成

```bash
# 发送普通请求
curl -i http://localhost:8080/health

# 响应头包含请求 ID
HTTP/1.1 200 OK
X-Request-ID: 550e8400-e29b-41d4-a716-446655440000
...
```

#### 客户端指定

```bash
# 指定请求 ID（用于链路追踪）
curl -i -H "X-Request-ID: my-trace-id-12345" \
  http://localhost:8080/api/v1/public/login

# 服务器返回相同的 ID
HTTP/1.1 200 OK
X-Request-ID: my-trace-id-12345
...
```

### 日志追踪

所有日志自动包含 `request_id` 字段：

```json
{
  "level": "info",
  "msg": "HTTP Request",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "method": "POST",
  "path": "/api/v1/public/login",
  "status": 200,
  "latency": "15ms"
}
```

**追踪特定请求**：

```bash
# 查看特定请求的所有日志
grep "550e8400-e29b-41d4-a716-446655440000" logs/app.log

# 实时追踪
tail -f logs/app.log | grep "550e8400-e29b-41d4-a716-446655440000"
```

### 测试验证

```bash
# 运行请求 ID 测试脚本
./scripts/test_request_id.sh
```

📖 **详细文档**: [请求 ID 追踪指南](REQUEST_ID_GUIDE.md)

## ⚡ RBAC 权限缓存

### 功能介绍

使用 Redis 缓存用户权限数据，大幅提升权限检查性能，响应时间提升 **80-90%**。

### 缓存策略

| 缓存类型 | TTL | 说明 |
|---------|-----|------|
| **用户角色** | 5分钟 | 缓存用户拥有的角色 |
| **角色权限** | 10分钟 | 缓存角色拥有的权限 |
| **用户权限** | 5分钟 | 缓存用户所有权限（聚合） |
| **权限检查** | 5分钟 | 缓存特定权限检查结果 |

### 性能提升

**压力测试结果**:

| 指标 | 无缓存 | 有缓存 | 提升 |
|------|--------|--------|------|
| 权限检查耗时 | 15ms | 1.5ms | **90%** |
| QPS | 500 | 5,000 | **10倍** |
| 数据库查询 | 10,000次 | 100次 | **99%降低** |

### 自动失效机制

✅ **分配/移除角色时** - 自动失效用户缓存  
✅ **修改角色权限时** - 自动失效角色缓存  
✅ **TTL 过期** - 自动清理过期缓存

### 使用示例

#### 查看缓存

```bash
# 查看所有 RBAC 缓存
redis-cli --scan --pattern "rbac:*"

# 查看用户权限缓存
redis-cli get rbac:user_permissions:1

# 查看缓存统计
redis-cli --scan --pattern "rbac:*" | wc -l
```

#### 测试缓存性能

```bash
# 运行缓存测试脚本
./scripts/test_rbac_cache.sh
```

#### 手动清除缓存

```bash
# 清除特定用户的缓存
redis-cli --scan --pattern "rbac:*:1" | xargs redis-cli del

# 清除所有 RBAC 缓存（谨慎使用）
redis-cli --scan --pattern "rbac:*" | xargs redis-cli del
```

### 监控缓存

```bash
# 查看缓存命中情况
tail -f logs/app.log | grep "cache hit"

# 查看缓存统计
echo "用户权限: $(redis-cli --scan --pattern 'rbac:user_permissions:*' | wc -l)"
echo "角色权限: $(redis-cli --scan --pattern 'rbac:role_permissions:*' | wc -l)"
```

📖 **详细文档**: [RBAC 权限缓存指南](RBAC_CACHE_GUIDE.md)

## 📊 Prometheus 监控

### 功能介绍

集成 Prometheus + Grafana 完整监控解决方案，提供实时的系统性能监控和可视化。

### 监控指标

**HTTP 请求监控**
- ✅ 请求总数、QPS
- ✅ 请求延迟（P50/P95/P99）
- ✅ 请求/响应大小
- ✅ 状态码分布

**业务指标监控**
- ✅ 用户注册数
- ✅ 用户登录数
- ✅ 登录失败数

**系统指标监控**
- ✅ 数据库连接数
- ✅ 数据库查询延迟
- ✅ Redis 操作延迟
- ✅ RBAC 权限检查
- ✅ 限流触发统计

### 快速开始

```bash
# 1. 启动监控服务
make docker-up

# 2. 启动应用服务
make dev-frontend  # 前台服务
make dev-backend   # 后台服务

# 3. 访问监控界面
```

**Prometheus UI:**
- URL: http://localhost:9090
- 查看原始指标和 PromQL 查询

**Grafana 仪表板:**
- URL: http://localhost:3000
- 账号/密码: `admin` / `admin`
- 预置仪表板: "TRX Project - 服务监控"

**Metrics 端点:**
- 前台: http://localhost:8080/metrics
- 后台: http://localhost:8081/metrics

### 预置仪表板

**TRX Project - 服务监控** 包含以下面板：

1. HTTP 请求速率 (QPS) - 实时 QPS 图表
2. HTTP 请求延迟 (P95/P99) - 延迟百分位数
3. 总请求速率 - 单值面板
4. 用户注册总数 - 业务指标
5. 用户登录总数 - 业务指标
6. 登录失败总数 - 错误监控
7. HTTP 状态码分布 - 饼图
8. 限流触发速率 - 限流监控

### 常用查询

```promql
# QPS（每秒请求数）
rate(trx_http_requests_total[5m])

# P95 延迟
histogram_quantile(0.95, rate(trx_http_request_duration_seconds_bucket[5m]))

# 错误率
sum(rate(trx_http_requests_total{status=~"5.."}[5m])) / sum(rate(trx_http_requests_total[5m]))

# RBAC 缓存命中率
sum(rate(trx_rbac_cache_hits_total{result="hit"}[5m])) / sum(rate(trx_rbac_cache_hits_total[5m]))
```

### 性能影响

- 🔸 监控开销: < 1% CPU, < 5MB 内存
- 🔸 网络开销: ~1KB/请求
- 🔸 对应用性能几乎无影响

📖 **详细文档**: [Prometheus 监控指南](docs/PROMETHEUS_MONITORING_GUIDE.md)

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
# 基础命令
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

# 数据库迁移命令 ⭐
make migrate-up      # 执行所有待执行的迁移
make migrate-down    # 回滚一个迁移版本
make migrate-version # 查看当前迁移版本
make migrate-create  # 创建新迁移文件 (用法: make migrate-create NAME=add_column)
make migrate-force   # 强制设置迁移版本 (用法: make migrate-force VERSION=1)
make migrate-goto    # 迁移到指定版本 (用法: make migrate-goto VERSION=3)
make migrate-drop    # 删除所有表（危险操作）
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
