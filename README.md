# TRX Project

> 基于 Gin 框架的企业级 Go Web 服务 - **前后台完全分离架构**

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

## ✨ 项目特点

- 🎯 **前后台完全分离** - 独立服务、独立端口、独立认证
- 🔐 **完整的 RBAC 权限系统** - 角色、权限、用户管理
- ⚡ **高性能缓存** - Redis 缓存，权限检查性能提升 90%
- 🚦 **三层限流保护** - 全局、IP、用户三级限流
- 📊 **监控可观测性** - Prometheus + Grafana + Jaeger
- 🔍 **分布式追踪** - OpenTelemetry + Jaeger 链路追踪
- 📖 **完整的 API 文档** - Swagger 自动生成
- 🗄️ **数据库版本管理** - golang-migrate 迁移工具

## 🛠️ 技术栈

| 类型 | 技术 |
|------|------|
| **Web 框架** | Gin |
| **依赖注入** | Google Wire |
| **数据库** | MySQL + GORM |
| **缓存** | Redis |
| **消息队列** | Kafka |
| **监控** | Prometheus + Grafana |
| **链路追踪** | OpenTelemetry + Jaeger |
| **日志** | Uber Zap |
| **认证** | JWT |

## 📁 项目结构

```
trx-project/
├── cmd/                        # 程序入口
│   ├── frontend/               # 前台服务 (端口 8080)
│   ├── backend/                # 后台服务 (端口 8081)
│   └── migrate/                # 数据库迁移工具
├── internal/                   # 内部代码
│   ├── api/
│   │   ├── handler/            # HTTP 处理器
│   │   ├── middleware/         # 中间件
│   │   └── router/             # 路由
│   ├── model/                  # 数据模型
│   ├── repository/             # 数据访问层
│   └── service/                # 业务逻辑层
├── pkg/                        # 公共包
├── config/                     # 配置文件
├── migrations/                 # 数据库迁移文件
├── scripts/                    # 工具脚本
└── docs/                       # 文档
```

## 🚀 快速开始

### 1. 环境要求

- Go 1.21+
- MySQL 8.0+
- Redis 6.0+
- Docker & Docker Compose (可选)

### 2. 安装依赖

```bash
# 克隆项目
git clone <repository-url>
cd trx-project

# 安装 Go 依赖
make deps

# 安装工具
go install github.com/google/wire/cmd/wire@latest
go install github.com/cosmtrek/air@latest
```

### 3. 启动基础服务

```bash
# 使用 Docker 启动 MySQL、Redis、Kafka、Prometheus、Grafana、Jaeger
make docker-up

# 或者手动启动这些服务
```

### 4. 配置

```bash
# 复制配置文件
cp config/config.yaml.example config/config.yaml

# 编辑配置（根据实际环境调整）
vim config/config.yaml
```

### 5. 数据库迁移

```bash
# 执行数据库迁移（创建表结构和初始数据）
make migrate-up

# 查看当前迁移版本
make migrate-version
```

### 6. 生成代码

```bash
# 生成 Wire 依赖注入代码
make wire

# 生成 Swagger 文档
make swag
```

### 7. 运行服务

**开发模式（推荐）**

```bash
# 终端1：启动前台服务
./scripts/start-frontend.sh

# 终端2：启动后台服务
./scripts/start-backend.sh
```

**生产模式**

```bash
# 构建
make build

# 运行前台服务（端口 8080）
./bin/frontend

# 运行后台服务（端口 8081）
./bin/backend
```

## 📡 服务访问

### 前台服务 (端口 8080)

- **健康检查**: http://localhost:8080/health
- **API 文档**: http://localhost:8080/swagger/index.html
- **Metrics**: http://localhost:8080/metrics

### 后台服务 (端口 8081)

- **健康检查**: http://localhost:8081/health
- **API 文档**: http://localhost:8081/swagger/index.html
- **Metrics**: http://localhost:8081/metrics

### 监控服务

- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3000 (admin/admin)
- **Jaeger**: http://localhost:16686

## 🔐 JWT 认证

### 生成管理员 Token

```bash
go run scripts/generate_admin_token.go
```

输出示例：
```
✅ 管理员 Token 生成成功!

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

## 🧪 测试

```bash
# 测试前台服务
./scripts/test_frontend.sh

# 测试后台服务
./scripts/test_backend.sh

# 测试限流功能
./scripts/test_rate_limit.sh

# 测试 RBAC 缓存
./scripts/test_rbac_cache.sh

# 单元测试
make test
```

## 🛠️ 常用命令

```bash
# 构建
make build              # 构建所有服务
make build-frontend     # 构建前台
make build-backend      # 构建后台

# 运行
make run-frontend       # 运行前台
make run-backend        # 运行后台

# 开发
make dev-frontend       # 开发模式（热重载）
make dev-backend        # 开发模式（热重载）

# 数据库迁移
make migrate-up         # 执行迁移
make migrate-down       # 回滚迁移
make migrate-version    # 查看版本
make migrate-create NAME=add_column  # 创建新迁移

# Swagger
make swag               # 生成所有文档
make swag-frontend      # 生成前台文档
make swag-backend       # 生成后台文档

# Docker
make docker-up          # 启动服务
make docker-down        # 停止服务

# 清理
make clean              # 清理构建文件
```

## 📚 详细文档

> 💡 **文档导航**: 查看 [docs/INDEX.md](docs/INDEX.md) 获取完整的文档索引和分类

### 🚀 快速开始
- [快速开始指南](docs/getting-started/QUICK_START.md) - 详细的启动步骤 ⭐
- [快速测试指南](docs/getting-started/QUICK_TEST_GUIDE.md) - API 测试示例
- [启动指南](docs/getting-started/STARTUP_GUIDE.md) - 详细启动步骤

### 👨‍💻 开发指南
- [功能开发规范](docs/development/FEATURE_DEVELOPMENT_GUIDE.md) - 新功能开发完整流程 ⭐⭐⭐
- [快速参考卡片](docs/development/QUICK_REFERENCE.md) - 10分钟快速上手
- [订单管理案例](docs/features/ORDER_MANAGEMENT_GUIDE.md) - 完整实现案例

### 📐 架构与设计
- [系统架构](docs/architecture/ARCHITECTURE.md) - 系统架构说明 ⭐
- [项目结构](docs/architecture/PROJECT_STRUCTURE.md) - 目录结构详解
- [服务分离](docs/architecture/SEPARATED_SERVICES.md) - 前后台分离设计 ⭐

### ⚡ 核心功能
- [RBAC 权限指南](docs/features/RBAC_GUIDE.md) - 权限系统使用 ⭐
- [RBAC 缓存指南](docs/features/RBAC_CACHE_GUIDE.md) - 缓存优化
- [限流指南](docs/features/RATE_LIMIT_GUIDE.md) - 限流配置使用 ⭐
- [请求追踪](docs/features/REQUEST_ID_GUIDE.md) - 请求 ID 追踪
- [数据库迁移](docs/features/MIGRATION_GUIDE.md) - 迁移工具使用 ⭐

### 📊 监控与追踪
- [Prometheus 监控](docs/features/PROMETHEUS_MONITORING_GUIDE.md) - 监控配置 ⭐
- [OpenTelemetry 追踪](docs/features/OPENTELEMETRY_TRACING_GUIDE.md) - 链路追踪 ⭐

### 🛠️ 开发工具
- [Swagger 使用指南](docs/tools/SWAGGER_GUIDE.md) - API 文档生成 ⭐
- [Swagger 脚本](docs/tools/SWAGGER_SCRIPTS_GUIDE.md) - 脚本使用说明
- [Air 使用指南](docs/tools/AIR_USAGE_GUIDE.md) - 热重载开发
- [脚本使用说明](scripts/README.md) - 工具脚本

### 🚀 部署运维
- [部署指南](docs/deployment/DEPLOYMENT.md) - 生产环境部署 ⭐
- [API 文档](docs/deployment/API.md) - 接口说明
- [响应格式](docs/deployment/RESPONSE_FORMAT.md) - API 响应规范

## 🎯 核心特性详解

### 🔐 RBAC 权限系统

完整的基于角色的访问控制（RBAC）系统：

- ✅ 角色管理（创建、分配、查询）
- ✅ 权限管理（定义、分配、检查）
- ✅ 用户-角色关联
- ✅ 角色-权限关联
- ✅ Redis 缓存优化（性能提升 90%）

### 🚦 限流保护

三层限流策略保护系统安全：

| 限流类型 | 作用范围 | 默认配置 |
|---------|---------|----------|
| 全局限流 | 所有请求 | 1000/秒 |
| IP 限流 | 按 IP | 100/分钟 |
| 用户限流 | 按用户 | 1000/分钟 |

### 📊 监控可观测性

完整的监控解决方案：

- **Prometheus** - 指标采集（HTTP、业务、系统指标）
- **Grafana** - 可视化仪表板
- **Jaeger** - 分布式链路追踪
- **Request ID** - 请求全链路追踪

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

MIT License

---

⭐ 如果这个项目对你有帮助，请给个 Star！
