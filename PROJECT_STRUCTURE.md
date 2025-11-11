# 项目结构说明

```
trx-project/
├── cmd/                          # 应用程序入口
│   └── server/                   # 主服务器程序
│       ├── main.go              # 程序入口点
│       ├── wire.go              # Wire 依赖注入配置
│       └── wire_gen.go          # Wire 自动生成的代码
│
├── internal/                     # 私有应用代码
│   ├── api/                     # API 层
│   │   ├── handler/             # HTTP 请求处理器
│   │   │   └── user_handler.go # 用户相关 API 处理
│   │   ├── middleware/          # 中间件
│   │   │   ├── cors.go         # CORS 跨域中间件
│   │   │   ├── logger.go       # 日志中间件
│   │   │   └── recovery.go     # Panic 恢复中间件
│   │   └── router/              # 路由配置
│   │       └── router.go       # 路由定义和初始化
│   │
│   ├── service/                 # 业务逻辑层
│   │   ├── user_service.go     # 用户业务逻辑
│   │   └── user_service_test.go # 用户服务测试
│   │
│   ├── repository/              # 数据访问层
│   │   └── user_repository.go  # 用户数据访问
│   │
│   └── model/                   # 数据模型
│       └── user.go             # 用户数据模型
│
├── pkg/                         # 可被外部使用的公共库
│   ├── config/                  # 配置管理
│   │   └── config.go           # 配置加载和定义
│   ├── logger/                  # 日志封装
│   │   └── logger.go           # Zap 日志初始化
│   ├── database/                # 数据库连接
│   │   └── mysql.go            # MySQL 连接管理
│   ├── cache/                   # 缓存封装
│   │   └── redis.go            # Redis 连接管理
│   └── kafka/                   # 消息队列
│       └── kafka.go            # Kafka Producer/Consumer
│
├── config/                      # 配置文件目录
│   ├── config.yaml             # 主配置文件
│   └── config.yaml.example     # 配置文件示例
│
├── docs/                        # 文档目录
│   ├── API.md                  # API 接口文档
│   ├── ARCHITECTURE.md         # 架构设计文档
│   └── KAFKA_USAGE.md          # Kafka 使用指南
│
├── scripts/                     # 脚本文件
│   ├── init_db.sql             # 数据库初始化脚本
│   └── setup.sh                # 项目初始化脚本
│
├── bin/                         # 编译后的二进制文件（.gitignore）
│   └── server                  # 编译后的服务器程序
│
├── logs/                        # 日志文件目录（.gitignore）
│   ├── app.log                 # 应用日志
│   └── error.log               # 错误日志
│
├── .air.toml                    # Air 热重载配置
├── .gitignore                   # Git 忽略文件配置
├── docker-compose.yml           # Docker Compose 配置
├── go.mod                       # Go 模块定义
├── go.sum                       # Go 模块校验和
├── Makefile                     # Make 命令定义
├── README.md                    # 项目说明文档
├── QUICKSTART.md                # 快速开始指南
└── PROJECT_STRUCTURE.md         # 本文件

```

## 目录说明

### cmd/ - 应用程序入口

存放应用程序的主入口代码。每个子目录代表一个可执行程序。

- `server/`: 主 Web 服务器程序
- 可扩展添加其他程序，如 `consumer/` 用于 Kafka 消费者

### internal/ - 私有应用代码

包含应用程序的核心业务逻辑，外部包无法导入。

#### internal/api/
API 层，处理 HTTP 请求和响应。

- `handler/`: 实现具体的 HTTP 处理逻辑
- `middleware/`: 中间件，如日志、CORS、认证等
- `router/`: 路由配置和初始化

#### internal/service/
业务逻辑层，实现核心业务规则。

- 调用 repository 层进行数据操作
- 处理复杂的业务逻辑
- 与外部服务交互（如 Kafka、Redis）

#### internal/repository/
数据访问层，封装数据库操作。

- 提供 CRUD 接口
- 使用 GORM 进行数据库操作
- 返回原始错误给 service 层处理

#### internal/model/
数据模型层，定义数据结构。

- 数据库表映射
- JSON 序列化/反序列化标签
- 数据验证规则

### pkg/ - 公共库

可被其他项目引用的公共代码。

- `config/`: 配置文件加载和管理
- `logger/`: 日志初始化和封装
- `database/`: 数据库连接管理
- `cache/`: Redis 缓存管理
- `kafka/`: Kafka 消息队列封装

### config/ - 配置文件

存放应用程序的配置文件。

- `config.yaml`: 实际使用的配置（不提交到 Git）
- `config.yaml.example`: 配置模板（提交到 Git）

### docs/ - 文档

项目相关文档。

- `API.md`: API 接口文档
- `ARCHITECTURE.md`: 架构设计文档
- `KAFKA_USAGE.md`: Kafka 使用指南

### scripts/ - 脚本

项目相关的脚本文件。

- `init_db.sql`: 数据库初始化 SQL
- `setup.sh`: 项目一键初始化脚本

## 代码层次关系

```
HTTP Request
    ↓
Router (internal/api/router)
    ↓
Middleware (internal/api/middleware)
    ↓
Handler (internal/api/handler)
    ↓
Service (internal/service)
    ↓
Repository (internal/repository)
    ↓
Database/Cache/MQ
```

## 依赖关系

```
main.go
  ↓
wire (依赖注入)
  ↓
┌─────────────┬─────────────┬─────────────┐
│   Config    │   Logger    │  Database   │
└─────────────┴─────────────┴─────────────┘
                    ↓
            ┌───────────────┐
            │  Repository   │
            └───────────────┘
                    ↓
            ┌───────────────┐
            │    Service    │
            └───────────────┘
                    ↓
            ┌───────────────┐
            │    Handler    │
            └───────────────┘
                    ↓
            ┌───────────────┐
            │     Router    │
            └───────────────┘
```

## 文件命名规范

### Go 文件
- 使用小写字母和下划线
- 例如: `user_service.go`, `user_handler.go`

### 测试文件
- 在文件名后添加 `_test.go`
- 例如: `user_service_test.go`

### 接口定义
- 在实现文件中同时定义接口
- 接口名使用大驼峰
- 实现使用小驼峰

```go
// 接口
type UserService interface {
    Register(...) error
}

// 实现
type userService struct {
    // ...
}
```

## 包导入规范

导入顺序：
1. 标准库
2. 第三方库
3. 本项目内部包

```go
import (
    // 标准库
    "context"
    "fmt"
    
    // 第三方库
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
    
    // 本项目包
    "trx-project/internal/model"
    "trx-project/internal/repository"
)
```

## 添加新功能

### 1. 添加新的数据模型

1. 在 `internal/model/` 创建模型文件
2. 在 `cmd/server/main.go` 的 `autoMigrate()` 中添加

### 2. 添加新的 API

1. 在 `internal/repository/` 创建数据访问接口
2. 在 `internal/service/` 创建业务逻辑
3. 在 `internal/api/handler/` 创建 HTTP 处理器
4. 在 `internal/api/router/router.go` 添加路由
5. 更新 Wire 配置（如果需要）

### 3. 添加新的中间件

1. 在 `internal/api/middleware/` 创建中间件文件
2. 在 `router.go` 中注册中间件

### 4. 添加新的外部服务

1. 在 `pkg/` 创建相应的包
2. 在 `pkg/config/config.go` 添加配置
3. 在 `config/config.yaml` 添加配置项
4. 在 Wire 中配置依赖注入

## 最佳实践

### 1. 错误处理
- Repository 层返回原始错误
- Service 层处理错误并记录日志
- Handler 层转换为 HTTP 响应

### 2. 日志记录
- 使用结构化日志（Zap）
- 记录关键操作
- 不记录敏感信息

### 3. 配置管理
- 使用配置文件管理配置
- 敏感信息使用环境变量
- 提供配置示例文件

### 4. 测试
- 为 Service 层编写单元测试
- 使用 Mock 隔离依赖
- 保持测试覆盖率

### 5. 代码组织
- 遵循单一职责原则
- 保持函数简短
- 使用有意义的命名

## 扩展建议

### 当前已实现
- ✅ 基础 Web 框架（Gin）
- ✅ 日志系统（Zap）
- ✅ 依赖注入（Wire）
- ✅ 数据库 ORM（GORM + MySQL）
- ✅ 缓存（Redis）
- ✅ 消息队列（Kafka）
- ✅ 用户管理 CRUD
- ✅ 配置管理
- ✅ Docker 部署

### 可以添加的功能
- ⬜ JWT 认证和授权
- ⬜ Swagger API 文档
- ⬜ 单元测试和集成测试
- ⬜ API 限流
- ⬜ 熔断器
- ⬜ 分布式追踪（Jaeger）
- ⬜ 监控和指标（Prometheus）
- ⬜ 配置中心（Consul/Etcd）
- ⬜ 服务发现
- ⬜ gRPC 接口
- ⬜ GraphQL 接口
- ⬜ 文件上传/下载
- ⬜ 邮件发送
- ⬜ 短信发送
- ⬜ 定时任务
- ⬜ WebSocket 支持

## 参考资源

- [Go 项目布局标准](https://github.com/golang-standards/project-layout)
- [Gin 框架文档](https://gin-gonic.com/)
- [GORM 文档](https://gorm.io/)
- [Uber Zap 文档](https://pkg.go.dev/go.uber.org/zap)
- [Google Wire 文档](https://github.com/google/wire)

