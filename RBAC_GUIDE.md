# RBAC 权限系统使用指南

## 📚 概述

后台 API 已实现完整的 RBAC（Role-Based Access Control，基于角色的访问控制）权限系统，提供细粒度的权限管理。

## 🏗️ 架构设计

### 核心概念

```
用户 (User) ──┐
             ├──> 用户角色 (UserRole) ──> 角色 (Role) ──> 角色权限 (RolePermission) ──> 权限 (Permission)
用户 (User) ──┘
```

**关系说明：**
- 一个用户可以有多个角色
- 一个角色可以有多个权限
- 用户通过角色获得权限

### 数据模型

#### 1. 角色 (Role)
```go
type Role struct {
    ID          uint
    Name        string  // 角色名称：superadmin, admin, editor, viewer
    DisplayName string  // 显示名称：超级管理员、管理员、编辑、查看者
    Description string  // 角色描述
    Status      int     // 状态：1-启用 0-禁用
}
```

#### 2. 权限 (Permission)
```go
type Permission struct {
    ID          uint
    Code        string  // 权限编码：user:read, user:write
    Name        string  // 权限名称：查看用户、编辑用户
    Resource    string  // 资源：user, role, permission, statistics
    Action      string  // 操作：read, write, delete
    Description string  // 权限描述
}
```

#### 3. 用户角色关联 (UserRole)
```go
type UserRole struct {
    UserID  uint
    RoleID  uint
}
```

#### 4. 角色权限关联 (RolePermission)
```go
type RolePermission struct {
    RoleID       uint
    PermissionID uint
}
```

## 🎭 默认角色

系统预置了 4 个角色：

| 角色 | 名称 | 权限 | 说明 |
|------|------|------|------|
| **superadmin** | 超级管理员 | 所有权限 | 完全控制，包括 RBAC 管理 |
| **admin** | 管理员 | user:*, statistics:read | 用户管理和统计查看 |
| **editor** | 编辑员 | user:read, user:write, statistics:read | 查看和编辑，不能删除 |
| **viewer** | 查看者 | user:read, statistics:read | 只能查看 |

## 🔑 默认权限

| 权限代码 | 权限名称 | 资源 | 操作 | 说明 |
|---------|---------|------|------|------|
| **user:read** | 查看用户 | user | read | 查看用户列表和详情 |
| **user:write** | 编辑用户 | user | write | 编辑用户信息、状态、重置密码 |
| **user:delete** | 删除用户 | user | delete | 删除用户 |
| **rbac:manage** | RBAC管理 | rbac | manage | 管理角色和权限 |
| **statistics:read** | 查看统计 | statistics | read | 查看统计信息 |

## 🚀 快速开始

### 1. 初始化 RBAC 数据

```bash
# 连接数据库
mysql -u root -p trx_dev

# 执行初始化脚本
source scripts/init_rbac.sql
# 或
mysql -u root -p trx_dev < scripts/init_rbac.sql
```

**脚本会自动：**
- ✅ 创建 RBAC 相关表
- ✅ 插入默认角色
- ✅ 插入默认权限
- ✅ 配置角色权限关系
- ✅ 为用户 ID=1 分配超级管理员角色

### 2. 生成管理员 Token

```bash
go run scripts/generate_admin_with_role.go
```

输出示例：
```
==================== 超级管理员 Token ====================
✅ 超级管理员 Token 生成成功!

Token 信息:
  User ID:  1
  Username: admin
  Role:     superadmin
  权限:     所有权限（包括 RBAC 管理）

Token:
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### 3. 启动后台服务

```bash
./bin/backend
```

### 4. 测试 RBAC 接口

```bash
# 保存 Token
export ADMIN_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# 1. 查看所有角色
curl -H "Authorization: Bearer $ADMIN_TOKEN" \
  http://localhost:8081/api/v1/admin/rbac/roles

# 2. 查看所有权限
curl -H "Authorization: Bearer $ADMIN_TOKEN" \
  http://localhost:8081/api/v1/admin/rbac/permissions

# 3. 查看用户角色
curl -H "Authorization: Bearer $ADMIN_TOKEN" \
  http://localhost:8081/api/v1/admin/users/1/roles

# 4. 查看用户权限
curl -H "Authorization: Bearer $ADMIN_TOKEN" \
  http://localhost:8081/api/v1/admin/users/1/permissions
```

## 📡 API 接口

### RBAC 管理接口

**所有 RBAC 管理接口都需要 `rbac:manage` 权限**

```
GET    /api/v1/admin/rbac/roles                      # 获取角色列表
GET    /api/v1/admin/rbac/roles/:id                  # 获取角色详情
POST   /api/v1/admin/rbac/roles                      # 创建角色
POST   /api/v1/admin/rbac/roles/:id/permissions      # 为角色分配权限
GET    /api/v1/admin/rbac/permissions                # 获取权限列表
```

### 用户角色管理接口

**需要 `rbac:manage` 权限**

```
POST   /api/v1/admin/users/:user_id/role             # 为用户分配角色
GET    /api/v1/admin/users/:user_id/roles            # 获取用户角色
GET    /api/v1/admin/users/:user_id/permissions      # 获取用户权限
```

### 用户管理接口（带权限控制）

```
GET    /api/v1/admin/users                           # 需要 user:read
GET    /api/v1/admin/users/:id                       # 需要 user:read
PUT    /api/v1/admin/users/:id/status                # 需要 user:write
POST   /api/v1/admin/users/:id/reset-password        # 需要 user:write
DELETE /api/v1/admin/users/:id                       # 需要 user:delete
```

### 统计信息接口

```
GET    /api/v1/admin/statistics/users                # 需要 statistics:read
```

## 🔧 使用示例

### 1. 创建新角色

```bash
curl -X POST http://localhost:8081/api/v1/admin/rbac/roles \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "moderator",
    "display_name": "版主",
    "description": "内容审核人员"
  }'
```

### 2. 为角色分配权限

```bash
# 为 moderator 角色分配 user:read 和 user:write 权限
curl -X POST http://localhost:8081/api/v1/admin/rbac/roles/5/permissions \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "permission_ids": [1, 2]
  }'
```

### 3. 为用户分配角色

```bash
# 为用户 ID=2 分配 editor 角色
curl -X POST http://localhost:8081/api/v1/admin/users/2/role \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "role_id": 3
  }'
```

### 4. 查看用户的所有权限

```bash
curl -H "Authorization: Bearer $ADMIN_TOKEN" \
  http://localhost:8081/api/v1/admin/users/2/permissions
```

### 5. 权限验证

```bash
# 使用普通管理员 Token (只有 user:read, user:write, user:delete)
NORMAL_ADMIN_TOKEN="..."

# ✅ 可以访问 - 需要 user:read
curl -H "Authorization: Bearer $NORMAL_ADMIN_TOKEN" \
  http://localhost:8081/api/v1/admin/users

# ❌ 拒绝访问 - 需要 rbac:manage
curl -H "Authorization: Bearer $NORMAL_ADMIN_TOKEN" \
  http://localhost:8081/api/v1/admin/rbac/roles
# 返回: {"code":403,"message":"Permission denied: rbac:manage"}
```

## 🛠️ 中间件使用

### RequirePermission - 单一权限

```go
// 要求必须有 user:write 权限
adminUsers.PUT("/:id/status", 
    middleware.RequirePermission("user:write", rbacService, logger),
    adminUserHandler.UpdateUserStatus)
```

### RequireAnyPermission - 任一权限

```go
// 只要有 user:read 或 user:write 任一权限即可
adminUsers.GET("/:id", 
    middleware.RequireAnyPermission([]string{"user:read", "user:write"}, rbacService, logger),
    adminUserHandler.GetUser)
```

### RequireAllPermissions - 所有权限

```go
// 必须同时拥有 user:write 和 rbac:manage 权限
adminUsers.PUT("/:id/role", 
    middleware.RequireAllPermissions([]string{"user:write", "rbac:manage"}, rbacService, logger),
    adminUserHandler.UpdateUserRole)
```

## 📊 权限检查流程

```
请求进入
    ↓
1. AdminAuth 中间件验证 JWT Token
    ↓
2. 从 Token 中提取 admin_id
    ↓
3. RequirePermission 中间件检查权限
    ↓
4. 查询数据库：user_roles → roles → role_permissions → permissions
    ↓
5. 检查用户是否拥有所需权限
    ↓
6. ✅ 有权限：继续执行 Handler
   ❌ 无权限：返回 403 Forbidden
```

## 🔍 数据库查询

### 查看用户的所有角色

```sql
SELECT 
  u.id as user_id,
  u.username,
  r.name as role_name,
  r.display_name
FROM users u
JOIN user_roles ur ON ur.user_id = u.id
JOIN roles r ON r.id = ur.role_id
WHERE u.id = 1;
```

### 查看用户的所有权限

```sql
SELECT DISTINCT
  u.id as user_id,
  u.username,
  p.code as permission_code,
  p.name as permission_name
FROM users u
JOIN user_roles ur ON ur.user_id = u.id
JOIN role_permissions rp ON rp.role_id = ur.role_id
JOIN permissions p ON p.id = rp.permission_id
WHERE u.id = 1 AND p.status = 1;
```

### 查看角色的所有权限

```sql
SELECT 
  r.name as role_name,
  p.code as permission_code,
  p.name as permission_name
FROM roles r
JOIN role_permissions rp ON rp.role_id = r.id
JOIN permissions p ON p.id = rp.permission_id
WHERE r.name = 'admin';
```

## 🎯 应用场景

### 场景1：多级管理员

- **超级管理员**: 完全控制
- **部门管理员**: 只能管理本部门用户
- **客服人员**: 只能查看和编辑用户信息
- **数据分析师**: 只能查看统计信息

### 场景2：细粒度权限

```
内容管理：
- content:read (查看内容)
- content:write (编辑内容)
- content:publish (发布内容)
- content:delete (删除内容)

订单管理：
- order:read (查看订单)
- order:process (处理订单)
- order:refund (退款)
- order:delete (删除订单)
```

### 场景3：动态权限

根据业务需要动态添加权限：

```bash
# 添加新权限
INSERT INTO permissions (code, name, resource, action, description)
VALUES ('product:manage', '商品管理', 'product', 'manage', '管理商品信息');

# 为角色分配新权限
INSERT INTO role_permissions (role_id, permission_id)
VALUES (2, (SELECT id FROM permissions WHERE code = 'product:manage'));
```

## ⚠️ 注意事项

### 1. 超级管理员保护

- 不要删除 superadmin 角色
- 至少保留一个超级管理员用户
- 建议创建专门的超级管理员账号，不用于日常操作

### 2. 权限编码规范

```
格式：<resource>:<action>
示例：
- user:read
- user:write
- user:delete
- order:manage
- statistics:read
```

### 3. 角色设计原则

- **最小权限原则**: 只给必要的权限
- **职责分离**: 不同角色有明确的职责边界
- **易于理解**: 角色名称和权限描述要清晰

### 4. 性能考虑

- 权限检查会查询数据库，建议：
  - 使用 Redis 缓存用户权限
  - 定期清理无效的角色权限关联
  - 监控权限查询性能

## 🔄 扩展建议

### 1. 添加权限缓存

```go
// 缓存用户权限到 Redis
func (s *rbacService) GetUserPermissionsWithCache(ctx context.Context, userID uint) ([]*model.Permission, error) {
    // 1. 尝试从 Redis 获取
    cacheKey := fmt.Sprintf("user_permissions:%d", userID)
    // ...
    
    // 2. 缓存未命中，从数据库查询
    permissions, err := s.repo.GetUserPermissions(ctx, userID)
    
    // 3. 写入缓存
    // ...
    
    return permissions, nil
}
```

### 2. 添加数据权限

```go
// 用户只能看到自己部门的数据
type DataPermission struct {
    UserID       uint
    DepartmentID uint
    CanViewAll   bool
}
```

### 3. 添加操作日志

```go
// 记录权限相关的操作
type RBACLog struct {
    UserID      uint
    Action      string  // assign_role, revoke_role, etc.
    TargetType  string  // user, role, permission
    TargetID    uint
    Description string
}
```

### 4. 添加权限审计

```go
// 定期审计权限分配情况
- 哪些用户有超级管理员权限？
- 哪些角色长期未使用？
- 哪些权限从未被使用？
```

## 📚 相关文件

```
新增文件：
- internal/model/rbac.go                    # RBAC 数据模型
- internal/repository/rbac_repository.go    # RBAC 数据访问层
- internal/service/rbac_service.go          # RBAC 业务逻辑层
- internal/api/handler/rbac_handler.go      # RBAC API 处理器
- internal/api/middleware/rbac.go           # RBAC 权限中间件
- scripts/init_rbac.sql                     # RBAC 初始化脚本
- scripts/generate_admin_with_role.go       # 生成管理员 Token 工具
- RBAC_GUIDE.md                             # 本文档

修改文件：
- internal/api/router/backend.go            # 应用 RBAC 权限检查
- cmd/backend/wire.go                       # 添加 RBAC 依赖注入
```

## 🎉 总结

RBAC 权限系统已完整实现：

✅ **完整的数据模型** - Role、Permission、关联表  
✅ **灵活的权限管理** - 支持多角色、多权限  
✅ **细粒度控制** - API 级别的权限检查  
✅ **易于扩展** - 可以动态添加角色和权限  
✅ **生产就绪** - 包含初始化脚本和管理工具  

现在后台 API 拥有企业级的权限管理能力！🚀

---

**有问题？**
- 查看 Swagger 文档: http://localhost:8081/swagger/index.html
- 检查数据库中的角色和权限配置
- 使用 `generate_admin_with_role.go` 生成测试 Token

