# 架构设计文档

## 项目概述

本项目是一个基于 Gin 框架的现代化 Go Web 服务，采用清晰的分层架构设计，集成了常用的中间件和基础设施组件。

## 技术栈

- **Web 框架**: Gin v1.x
- **日志库**: Uber Zap
- **依赖注入**: Google Wire
- **ORM**: GORM v2
- **数据库**: MySQL 8.0
- **缓存**: Redis 7
- **消息队列**: Kafka
- **配置管理**: YAML

## 架构设计

### 分层架构

项目采用经典的四层架构设计：

```
┌─────────────────────────────────────┐
│         API Layer (Handler)         │  HTTP 请求处理
├─────────────────────────────────────┤
│       Business Layer (Service)      │  业务逻辑
├─────────────────────────────────────┤
│    Data Access Layer (Repository)   │  数据访问
├─────────────────────────────────────┤
│         Model Layer (Model)         │  数据模型
└─────────────────────────────────────┘
```

#### 1. API 层 (Handler)

- **职责**: 处理 HTTP 请求和响应
- **位置**: `internal/api/handler/`
- **功能**:
  - 参数验证
  - 调用 Service 层
  - 构造 HTTP 响应
  - 错误处理

#### 2. 业务逻辑层 (Service)

- **职责**: 实现核心业务逻辑
- **位置**: `internal/service/`
- **功能**:
  - 业务规则处理
  - 事务管理
  - 调用 Repository 层
  - 缓存管理
  - 消息发送

#### 3. 数据访问层 (Repository)

- **职责**: 数据库操作
- **位置**: `internal/repository/`
- **功能**:
  - CRUD 操作
  - 数据查询
  - 数据持久化

#### 4. 模型层 (Model)

- **职责**: 定义数据结构
- **位置**: `internal/model/`
- **功能**:
  - 数据库表映射
  - 数据验证规则

### 中间件

位置: `internal/api/middleware/`

- **Logger**: 请求日志记录
- **Recovery**: Panic 恢复
- **CORS**: 跨域处理

### 基础设施层

位置: `pkg/`

- **Config**: 配置管理
- **Logger**: 日志封装
- **Database**: 数据库连接
- **Cache**: Redis 缓存
- **Kafka**: 消息队列

## 依赖注入

使用 Google Wire 进行编译时依赖注入：

```
Config → Logger → Database/Redis → Repository → Service → Handler → Router
```

### Wire 工作流程

1. 定义 Provider 函数 (`wire.go`)
2. 运行 `wire` 命令生成代码
3. 生成的代码 (`wire_gen.go`) 包含完整的依赖注入逻辑

## 数据流

### 请求流程

```
Client Request
    ↓
Gin Router
    ↓
Middleware (Logger, Recovery, CORS)
    ↓
Handler (参数验证)
    ↓
Service (业务逻辑)
    ↓
Repository (数据访问)
    ↓
Database/Cache
    ↓
Response
```

### 错误处理

1. **Repository 层**: 返回原始错误
2. **Service 层**: 记录错误日志，返回业务错误
3. **Handler 层**: 转换为 HTTP 响应

## 配置管理

使用 YAML 文件进行配置管理：

- **开发环境**: `config/config.yaml`
- **配置示例**: `config/config.yaml.example`

配置项包括：
- 服务器配置
- 数据库配置
- Redis 配置
- Kafka 配置
- 日志配置

## 日志管理

使用 Zap 结构化日志：

- **日志级别**: Debug, Info, Warn, Error, Fatal
- **日志格式**: JSON/Console
- **日志输出**: 
  - 标准输出
  - 文件 (`logs/app.log`)
  - 错误日志 (`logs/error.log`)

## 数据库设计

### 用户表 (users)

| 字段       | 类型         | 说明         |
| ---------- | ------------ | ------------ |
| id         | BIGINT       | 主键         |
| created_at | DATETIME     | 创建时间     |
| updated_at | DATETIME     | 更新时间     |
| deleted_at | DATETIME     | 删除时间     |
| username   | VARCHAR(50)  | 用户名(唯一) |
| email      | VARCHAR(100) | 邮箱(唯一)   |
| password   | VARCHAR(255) | 密码(加密)   |
| status     | INT          | 状态         |

## 性能优化

### 1. 数据库优化

- 连接池管理
- 索引优化
- 查询优化
- 软删除

### 2. 缓存策略

- Redis 缓存常用数据
- 缓存失效策略
- 缓存预热

### 3. 并发控制

- 数据库连接池
- Redis 连接池
- Goroutine 池

## 安全性

### 1. 密码安全

- 使用 bcrypt 加密
- 密码强度验证

### 2. 输入验证

- 参数类型验证
- 参数范围验证
- SQL 注入防护 (GORM)

### 3. CORS 配置

- 跨域请求控制
- 允许的请求方法

## 扩展性

### 1. 水平扩展

- 无状态设计
- 负载均衡
- 分布式缓存

### 2. 垂直扩展

- 数据库读写分离
- 缓存分层
- 消息队列解耦

## 开发规范

### 1. 代码规范

- 遵循 Go 代码规范
- 使用 golint 和 gofmt
- 编写单元测试

### 2. Git 规范

- 使用语义化提交信息
- 功能分支开发
- Code Review

### 3. 文档规范

- API 文档
- 代码注释
- 架构文档

## 部署架构

```
┌─────────────┐
│   Nginx     │  反向代理/负载均衡
└──────┬──────┘
       │
   ┌───┴────┐
   │        │
┌──▼──┐  ┌──▼──┐
│ App │  │ App │  Go 应用实例
└──┬──┘  └──┬──┘
   │        │
   └───┬────┘
       │
   ┌───▼────┐
   │ MySQL  │  主数据库
   └────────┘
       │
   ┌───▼────┐
   │ Redis  │  缓存
   └────────┘
       │
   ┌───▼────┐
   │ Kafka  │  消息队列
   └────────┘
```

## 监控和日志

### 1. 日志收集

- 应用日志
- 访问日志
- 错误日志

### 2. 性能监控

- 请求延迟
- 数据库性能
- 缓存命中率

### 3. 告警

- 错误率告警
- 性能告警
- 资源告警

## 未来优化

1. 添加认证和授权 (JWT)
2. 集成 Swagger 文档
3. 添加单元测试和集成测试
4. 性能监控和追踪 (Prometheus + Grafana)
5. 配置中心 (Consul/Etcd)
6. 服务发现
7. 分布式追踪 (Jaeger)
8. 限流和熔断

