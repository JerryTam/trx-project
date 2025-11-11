# 前后台接口区分说明

本文档说明项目中前台接口和后台接口的区分和使用方法。

## 接口分类

### 1. 前台接口（Frontend/User API）

**路径前缀**: `/api/v1/public` 或 `/api/v1/user`

**用途**: 面向最终用户的接口

**特点**:
- 部分接口无需认证（如注册、登录）
- 认证接口使用普通用户 Token
- 用户只能操作自己的数据
- 权限受限

### 2. 后台接口（Backend/Admin API）

**路径前缀**: `/api/v1/admin`

**用途**: 面向管理员的接口

**特点**:
- 所有接口都需要管理员认证
- 使用管理员 Token
- 可以管理所有用户数据
- 具有更高的权限

## API 路由结构

```
/api/v1
├── /public                    # 公开接口（无需认证）
│   ├── POST /register        # 用户注册
│   └── POST /login           # 用户登录
│
├── /user                      # 用户接口（需要用户认证）
│   ├── GET  /profile         # 获取当前用户信息
│   └── PUT  /profile         # 更新当前用户信息
│
└── /admin                     # 管理员接口（需要管理员认证）
    ├── /users                # 用户管理
    │   ├── GET    ""         # 获取用户列表
    │   ├── GET    /:id       # 获取用户详情
    │   ├── PUT    /:id/status # 更新用户状态
    │   ├── DELETE /:id       # 删除用户
    │   └── POST   /:id/reset-password # 重置密码
    │
    └── /statistics           # 统计信息
        └── GET /users        # 用户统计
```

## 认证机制

### 前台认证（用户）

**中间件**: `middleware.Auth()`

**Token 格式**: 普通 JWT Token

**请求示例**:
```bash
curl -X GET http://localhost:8080/api/v1/user/profile \
  -H "Authorization: Bearer <user_token>"
```

**特点**:
- 验证用户身份
- 检查 Token 有效性
- 将用户信息存入上下文

### 后台认证（管理员）

**中间件**: `middleware.AdminAuth()`

**Token 格式**: 管理员 JWT Token（临时实现需要 `admin_` 前缀）

**请求示例**:
```bash
curl -X GET http://localhost:8080/api/v1/admin/users \
  -H "Authorization: Bearer admin_<admin_token>"
```

**特点**:
- 验证管理员身份
- 检查管理员权限
- 将管理员信息和角色存入上下文

## 接口详细说明

### 前台接口

#### 1. 用户注册

```
POST /api/v1/public/register
```

**认证**: 无需认证

**请求参数**:
```json
{
  "username": "testuser",
  "email": "test@example.com",
  "password": "password123"
}
```

**响应**:
```json
{
  "code": 200,
  "message": "User registered successfully",
  "data": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com"
  }
}
```

#### 2. 用户登录

```
POST /api/v1/public/login
```

**认证**: 无需认证

**请求参数**:
```json
{
  "username": "testuser",
  "password": "password123"
}
```

**响应**:
```json
{
  "code": 200,
  "message": "Login successful",
  "data": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "token": "<jwt_token>"
  }
}
```

#### 3. 获取个人信息

```
GET /api/v1/user/profile
```

**认证**: 需要用户 Token

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com"
  }
}
```

#### 4. 更新个人信息

```
PUT /api/v1/user/profile
```

**认证**: 需要用户 Token

**请求参数**:
```json
{
  "email": "newemail@example.com"
}
```

### 后台接口

#### 1. 获取用户列表

```
GET /api/v1/admin/users?page=1&page_size=10&status=1&keyword=test
```

**认证**: 需要管理员 Token

**查询参数**:
- `page`: 页码（默认 1）
- `page_size`: 每页数量（默认 10）
- `status`: 用户状态（可选，0=禁用, 1=启用）
- `keyword`: 搜索关键词（可选）

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "username": "user1",
        "email": "user1@example.com",
        "status": 1
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 10
  }
}
```

#### 2. 获取用户详情

```
GET /api/v1/admin/users/:id
```

**认证**: 需要管理员 Token

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "username": "user1",
    "email": "user1@example.com",
    "status": 1,
    "created_at": "2023-01-01T00:00:00Z"
  }
}
```

#### 3. 更新用户状态

```
PUT /api/v1/admin/users/:id/status
```

**认证**: 需要管理员 Token

**请求参数**:
```json
{
  "status": 0
}
```

**说明**: 0=禁用, 1=启用

**响应**:
```json
{
  "code": 200,
  "message": "User status updated successfully",
  "data": {
    "id": 1,
    "status": 0
  }
}
```

#### 4. 删除用户

```
DELETE /api/v1/admin/users/:id
```

**认证**: 需要管理员 Token

**响应**:
```json
{
  "code": 200,
  "message": "User deleted successfully"
}
```

#### 5. 重置用户密码

```
POST /api/v1/admin/users/:id/reset-password
```

**认证**: 需要管理员 Token

**请求参数**:
```json
{
  "new_password": "newpassword123"
}
```

**响应**:
```json
{
  "code": 200,
  "message": "Password reset successfully"
}
```

#### 6. 获取用户统计

```
GET /api/v1/admin/statistics/users
```

**认证**: 需要管理员 Token

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "total_users": 1000,
    "active_users": 850,
    "inactive_users": 150,
    "new_users_today": 12,
    "new_users_week": 89,
    "new_users_month": 356
  }
}
```

## 认证错误响应

### 缺少 Token

```json
{
  "code": 401,
  "message": "Missing authorization token"
}
```

HTTP 状态码: `401 Unauthorized`

### Token 无效

```json
{
  "code": 401,
  "message": "Invalid token"
}
```

HTTP 状态码: `401 Unauthorized`

### 权限不足

```json
{
  "code": 403,
  "message": "Admin access required"
}
```

HTTP 状态码: `403 Forbidden`

## 测试示例

### 前台接口测试

```bash
# 1. 注册
curl -X POST http://localhost:8080/api/v1/public/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'

# 2. 登录
curl -X POST http://localhost:8080/api/v1/public/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'

# 3. 获取个人信息（需要 Token）
curl -X GET http://localhost:8080/api/v1/user/profile \
  -H "Authorization: Bearer <your_token>"
```

### 后台接口测试

```bash
# 1. 获取用户列表（需要管理员 Token）
curl -X GET "http://localhost:8080/api/v1/admin/users?page=1&page_size=10" \
  -H "Authorization: Bearer admin_<admin_token>"

# 2. 获取用户详情
curl -X GET http://localhost:8080/api/v1/admin/users/1 \
  -H "Authorization: Bearer admin_<admin_token>"

# 3. 更新用户状态
curl -X PUT http://localhost:8080/api/v1/admin/users/1/status \
  -H "Authorization: Bearer admin_<admin_token>" \
  -H "Content-Type: application/json" \
  -d '{"status": 0}'

# 4. 重置用户密码
curl -X POST http://localhost:8080/api/v1/admin/users/1/reset-password \
  -H "Authorization: Bearer admin_<admin_token>" \
  -H "Content-Type: application/json" \
  -d '{"new_password": "newpass123"}'

# 5. 获取用户统计
curl -X GET http://localhost:8080/api/v1/admin/statistics/users \
  -H "Authorization: Bearer admin_<admin_token>"
```

## 中间件使用

### 在代码中使用

```go
// 前台路由（需要用户认证）
user := v1.Group("/user")
user.Use(middleware.Auth(logger))
{
    user.GET("/profile", handler.GetProfile)
    user.PUT("/profile", handler.UpdateProfile)
}

// 后台路由（需要管理员认证）
admin := v1.Group("/admin")
admin.Use(middleware.AdminAuth(logger))
{
    admin.GET("/users", adminHandler.ListUsers)
}

// 可选认证（有 Token 则验证，没有也不阻止）
public := v1.Group("/public")
public.Use(middleware.OptionalAuth(logger))
{
    public.GET("/posts", handler.ListPosts) // 登录用户可能看到不同内容
}
```

### 从上下文获取用户信息

```go
import "trx-project/internal/api/middleware"

// 获取用户 ID
func SomeHandler(c *gin.Context) {
    userID, exists := middleware.GetUserID(c)
    if !exists {
        response.Unauthorized(c, "User not authenticated")
        return
    }
    
    // 使用 userID
}

// 获取管理员 ID
func AdminHandler(c *gin.Context) {
    adminID, exists := middleware.GetAdminID(c)
    if !exists {
        response.Unauthorized(c, "Admin not authenticated")
        return
    }
    
    // 使用 adminID
}

// 获取管理员角色
func CheckRole(c *gin.Context) {
    role, exists := middleware.GetAdminRole(c)
    if !exists || role != "superadmin" {
        response.Forbidden(c, "Permission denied")
        return
    }
}
```

## 权限控制

### 权限检查中间件

```go
// 使用权限检查
admin.Group("/sensitive")
admin.Use(middleware.CheckPermission("user:delete", logger))
{
    admin.DELETE("/users/:id", handler.DeleteUser)
}
```

### 自定义权限

在 `middleware/auth.go` 中实现 `CheckPermission` 的真实逻辑：

```go
func CheckPermission(permission string, logger *zap.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        role, _ := GetAdminRole(c)
        
        // 从数据库或缓存获取角色权限
        permissions := getRolePermissions(role)
        
        if !contains(permissions, permission) {
            response.Forbidden(c, "Permission denied")
            c.Abort()
            return
        }
        
        c.Next()
    }
}
```

## 安全建议

### 1. Token 管理
- 使用 JWT 标准实现
- 设置合理的过期时间（如用户 Token 7天，管理员 Token 1天）
- 实现 Token 刷新机制
- 在 Redis 中维护 Token 黑名单

### 2. 密码安全
- 使用 bcrypt 加密（已实现）
- 强制密码复杂度
- 实现密码重试限制
- 记录密码修改日志

### 3. 权限控制
- 实现基于角色的访问控制（RBAC）
- 细粒度权限管理
- 记录敏感操作日志
- 定期审计管理员操作

### 4. API 限流
- 前台接口：更宽松的限制（如 100次/分钟）
- 后台接口：更严格的限制（如 1000次/分钟）
- 按 IP 或 Token 限流

### 5. 数据验证
- 严格的参数验证
- SQL 注入防护（GORM已提供）
- XSS 防护
- CSRF 防护

## TODO

### 短期
- [ ] 实现真实的 JWT Token 生成和验证
- [ ] 完善用户认证逻辑
- [ ] 实现 Token 刷新机制
- [ ] 添加更多前台用户接口

### 中期
- [ ] 实现基于角色的权限控制（RBAC）
- [ ] 添加操作日志记录
- [ ] 实现 API 限流
- [ ] 添加更多后台管理接口

### 长期
- [ ] 实现细粒度权限控制
- [ ] 添加审计日志
- [ ] 实现多租户支持
- [ ] 添加 API 监控和告警

## 相关文档

- [API 文档](API.md) - 完整的 API 接口文档
- [响应格式](RESPONSE_FORMAT.md) - 统一响应格式说明
- [架构设计](ARCHITECTURE.md) - 系统架构设计

