# 前后台服务完全分离 - 实现文档

## 🎯 架构设计

### 独立服务架构

前后台现在是**完全独立**的两个服务：

```
trx-project/
├── cmd/
│   ├── frontend/         # 前台服务入口 ✅
│   │   ├── main.go
│   │   ├── wire.go
│   │   └── wire_gen.go
│   └── backend/          # 后台服务入口 ✅
│       ├── main.go
│       ├── wire.go
│       └── wire_gen.go
├── internal/
│   └── api/
│       └── router/
│           ├── frontend.go  # 前台路由 ✅
│           └── backend.go   # 后台路由 ✅
└── pkg/
    └── jwt/              # JWT 认证 ✅
        └── jwt.go
```

### 服务隔离

| 特性 | 前台服务 | 后台服务 |
|------|---------|---------|
| 入口程序 | `cmd/frontend/main.go` | `cmd/backend/main.go` |
| 路由配置 | `router.SetupFrontend()` | `router.SetupBackend()` |
| 端口 | **8080** | **8081** |
| 认证方式 | JWT (user 角色) | JWT (admin/superadmin 角色) |
| Wire配置 | `cmd/frontend/wire.go` | `cmd/backend/wire.go` |
| 可执行文件 | `bin/frontend` | `bin/backend` |

## ✅ 已实现功能

### 1. JWT 认证系统

**文件**: `pkg/jwt/jwt.go`

- ✅ JWT Token 生成
- ✅ JWT Token 解析
- ✅ Token 验证
- ✅ Token 刷新
- ✅ 基于角色的 Token（user/admin/superadmin）

### 2. 前台服务（Frontend）

**端口**: `8080`

**路由** (`internal/api/router/frontend.go`):
```
/health                          # 健康检查
/api/v1/public/register          # 注册
/api/v1/public/login             # 登录
/api/v1/user/profile             # 获取个人信息（需认证）
/api/v1/user/profile             # 更新个人信息（需认证）
```

**特点**:
- 面向最终用户
- 返回 user 角色的 JWT Token
- 用户只能操作自己的数据

### 3. 后台服务（Backend）

**端口**: `8081`

**路由** (`internal/api/router/backend.go`):
```
/health                                      # 健康检查
/api/v1/admin/users                          # 用户列表
/api/v1/admin/users/:id                      # 用户详情
/api/v1/admin/users/:id/status               # 更新状态
/api/v1/admin/users/:id                      # 删除用户
/api/v1/admin/users/:id/reset-password       # 重置密码
/api/v1/admin/statistics/users               # 用户统计
```

**特点**:
- 面向管理员
- 需要 admin/superadmin 角色的 JWT Token
- 完整的数据管理权限

### 4. 认证中间件

**文件**: `internal/api/middleware/auth.go`

更新为真实的 JWT 验证：

```go
// 用户认证（前台）
middleware.Auth(jwtSecret, logger)

// 管理员认证（后台）
middleware.AdminAuth(jwtSecret, logger)

// 可选认证
middleware.OptionalAuth(jwtSecret, logger)
```

### 5. Service 层更新

**文件**: `internal/service/user_service.go`

- ✅ 注册返回 JWT Token
- ✅ 登录返回 JWT Token
- ✅ Token 有效期配置（用户7天，管理员1天）

## 🚀 使用方法

### 1. 配置

**文件**: `config/config.yaml`

```yaml
server:
  host: 0.0.0.0
  port: 8080  # 前台端口，后台自动使用 8081

jwt:
  secret: your-secret-key-change-in-production
  issuer: trx-project
  expire_hours: 168      # 用户 Token: 7天
  admin_expire_hours: 24 # 管理员 Token: 1天
```

### 2. 构建服务

```bash
# 构建前台服务
make build-frontend

# 构建后台服务
make build-backend

# 构建全部
make build
```

### 3. 运行服务

```bash
# 运行前台服务（端口 8080）
./bin/frontend

# 运行后台服务（端口 8081）
./bin/backend

# 或使用 make 命令
make run-frontend  # 前台
make run-backend   # 后台
```

### 4. 测试前台服务

```bash
# 1. 注册用户
curl -X POST http://localhost:8080/api/v1/public/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'

# 响应包含 token
{
  "code": 200,
  "message": "User registered successfully",
  "data": {
    "user": {
      "id": 1,
      "username": "testuser",
      "email": "test@example.com"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}

# 2. 登录获取 Token
curl -X POST http://localhost:8080/api/v1/public/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'

# 3. 使用 Token 访问
TOKEN="<your_jwt_token>"
curl http://localhost:8080/api/v1/user/profile \
  -H "Authorization: Bearer $TOKEN"
```

### 5. 测试后台服务

```bash
# 需要管理员 Token（包含 admin 或 superadmin 角色）
ADMIN_TOKEN="<admin_jwt_token>"

# 1. 获取用户列表
curl http://localhost:8081/api/v1/admin/users \
  -H "Authorization: Bearer $ADMIN_TOKEN"

# 2. 更新用户状态
curl -X PUT http://localhost:8081/api/v1/admin/users/1/status \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"status": 0}'

# 3. 获取统计信息
curl http://localhost:8081/api/v1/admin/statistics/users \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

## 🔐 JWT Token 说明

### Token 结构

```json
{
  "user_id": 1,
  "username": "testuser",
  "role": "user",          // user / admin / superadmin
  "iss": "trx-project",
  "exp": 1699999999,
  "iat": 1699999999,
  "nbf": 1699999999
}
```

### 生成 Token

用户注册或登录时自动生成：

```go
// 用户 Token (7天有效期)
token, err := jwt.GenerateToken(userID, username, "user", jwtConfig)

// 管理员 Token (1天有效期)
token, err := jwt.GenerateToken(adminID, username, "admin", adminJWTConfig)
```

### Token 验证

```go
// 解析 Token
claims, err := jwt.ParseToken(tokenString, jwtSecret)

// 检查角色
if claims.Role == "user" {
    // 用户操作
}
if claims.Role == "admin" || claims.Role == "superadmin" {
    // 管理员操作
}
```

## 📊 服务对比

### 前台服务 (Frontend - Port 8080)

```
特点：
- 面向最终用户
- 公开注册/登录接口
- 需要认证才能访问个人信息
- 返回 user 角色 Token
- Token 有效期: 7天

主要功能：
- 用户注册
- 用户登录
- 个人信息管理
- （未来）订单、购物车等
```

### 后台服务 (Backend - Port 8081)

```
特点：
- 面向管理员
- 所有接口都需要管理员认证
- 需要 admin/superadmin 角色 Token
- Token 有效期: 1天
- 完整的数据管理权限

主要功能：
- 用户管理
- 数据统计
- 系统配置
- 操作日志
```

## 🔄 服务间通信

前后台服务是**完全隔离**的：

- ✅ 独立的进程
- ✅ 独立的端口
- ✅ 独立的路由
- ✅ 独立的认证
- ✅ 可以部署在不同服务器

如果需要服务间通信：
- 使用 HTTP API 调用
- 使用消息队列（Kafka）
- 使用共享数据库
- 使用 RPC（gRPC）

## 📝 开发指南

### 添加新的前台接口

1. 在 `internal/api/handler/user_handler.go` 添加处理器
2. 在 `internal/api/router/frontend.go` 注册路由
3. 重新构建: `make build-frontend`

### 添加新的后台接口

1. 在 `internal/api/handler/admin_xxx_handler.go` 添加处理器
2. 在 `internal/api/router/backend.go` 注册路由
3. 在 `cmd/backend/wire.go` 添加依赖注入
4. 重新生成 Wire 代码: `cd cmd/backend && wire`
5. 重新构建: `make build-backend`

### 生成管理员 Token

临时方案 - 创建工具生成管理员 Token:

```go
package main

import (
	"fmt"
	"time"
	"trx-project/pkg/jwt"
)

func main() {
	config := jwt.Config{
		Secret:     "your-secret-key-change-in-production",
		Issuer:     "trx-project",
		ExpireTime: 24 * time.Hour,
	}
	
	// 生成管理员 Token
	token, _ := jwt.GenerateToken(1, "admin", "superadmin", config)
	fmt.Println("Admin Token:", token)
}
```

## 🎨 架构优势

### 1. 完全隔离
- 前后台可以独立部署
- 互不影响
- 独立扩展

### 2. 安全性
- 不同的端口
- 不同的认证方式
- 后台接口不对外暴露

### 3. 可维护性
- 代码分离清晰
- 独立开发和测试
- 独立升级

### 4. 可扩展性
- 可以部署多个前台实例
- 可以部署多个后台实例
- 可以添加负载均衡

## 🔮 后续优化

### 短期
- [ ] 创建管理员 Token 生成工具
- [ ] 完善 Token 刷新机制
- [ ] 添加 Token 黑名单（Redis）
- [ ] 实现操作日志记录

### 中期
- [ ] 添加 API 网关
- [ ] 实现服务注册与发现
- [ ] 添加分布式追踪
- [ ] 实现配置中心

### 长期
- [ ] 微服务化
- [ ] 服务网格（Service Mesh）
- [ ] 容器化部署（Kubernetes）
- [ ] 自动化CI/CD

## 📚 相关文件

```
新增：
- pkg/jwt/jwt.go                    # JWT 实现 ✅
- cmd/frontend/main.go              # 前台入口 ✅
- cmd/frontend/wire.go              # 前台依赖注入 ✅
- cmd/backend/main.go               # 后台入口 ✅
- cmd/backend/wire.go               # 后台依赖注入 ✅
- internal/api/router/frontend.go   # 前台路由 ✅
- internal/api/router/backend.go    # 后台路由 ✅

修改：
- pkg/config/config.go              # 添加 JWT 配置 ✅
- config/config.yaml                # 添加 JWT 配置 ✅
- internal/api/middleware/auth.go   # 使用真实 JWT ✅
- internal/service/user_service.go  # 返回 Token ✅
- internal/api/handler/user_handler.go # 返回 Token ✅
- Makefile                          # 支持前后台构建 ✅
```

## ✅ 验证清单

- [x] JWT 认证实现
- [x] 前台服务独立入口
- [x] 后台服务独立入口
- [x] 前台独立路由
- [x] 后台独立路由
- [x] Wire 依赖注入
- [x] 服务构建成功
- [x] Token 生成和验证
- [x] 角色权限控制

## 🎉 总结

前后台服务已经**完全分离**：

✅ **独立的程序入口**  
✅ **独立的路由配置**  
✅ **独立的端口**（8080/8081）  
✅ **真实的 JWT 认证**  
✅ **基于角色的权限控制**  
✅ **完全隔离的服务**  

现在你有两个完全独立的服务，可以分别部署、扩展和维护！🚀

