# Swagger 文档分离实现总结

## 问题描述

在实现前后台服务分离后，Swagger 文档生成时出现了一个问题：前台和后台的 Swagger 文档包含了相同的接口列表，即所有接口（前台 + 后台）都出现在两份文档中。

**预期行为**：
- 前台文档应该只包含前台接口（`/public/*`、`/user/*`）
- 后台文档应该只包含后台接口（`/admin/*`）

## 根本原因

`swag init` 工具在生成文档时，会扫描整个项目的所有 handler 文件，收集所有的 Swagger 注释。这导致：

1. 前台文档生成时扫描了 `internal/api/handler/` 下的所有文件，包括：
   - `user_handler.go`（前台）
   - `admin_user_handler.go`（后台）
   - `rbac_handler.go`（后台）

2. 后台文档生成时也扫描了相同的所有文件

3. `--exclude` 参数在实际测试中没有生效

## 解决方案

### 1. 使用独立的实例名称

为前后台生成不同的 Swagger 实例：

```bash
# 前台
swag init -g cmd/frontend/main.go -o cmd/frontend/docs --instanceName frontend

# 后台
swag init -g cmd/backend/main.go -o cmd/backend/docs --instanceName backend
```

这会生成：
- `cmd/frontend/docs/frontend_docs.go`、`frontend_swagger.json`、`frontend_swagger.yaml`
- `cmd/backend/docs/backend_docs.go`、`backend_swagger.json`、`backend_swagger.yaml`

### 2. 创建过滤脚本

创建 `scripts/filter_swagger.go` 脚本，用于从生成的 Swagger 文档中移除不属于该服务的接口：

**功能**：
- 读取 Swagger JSON 文件
- 根据路径前缀过滤掉不需要的接口
- 支持多个排除前缀
- 将过滤后的结果写回文件

**使用方式**：
```bash
# 前台：移除 /admin 和 /users 开头的接口
go run scripts/filter_swagger.go \
  cmd/frontend/docs/frontend_swagger.json \
  cmd/frontend/docs/frontend_swagger.json \
  admin users

# 后台：移除 /public 和 /user 开头的接口
go run scripts/filter_swagger.go \
  cmd/backend/docs/backend_swagger.json \
  cmd/backend/docs/backend_swagger.json \
  public user
```

**处理 Git Bash 路径转换问题**：
- Git Bash 会自动将 `/admin` 转换为 Windows 路径（如 `C:/Program Files/Git/admin`）
- 解决方案：脚本接受不带 `/` 的前缀（如 `admin`），然后自动添加 `/`

### 3. 更新 Makefile

集成文档生成和过滤流程：

```makefile
swag-frontend:
	@swag init -g cmd/frontend/main.go -o cmd/frontend/docs --instanceName frontend
	@go run scripts/filter_swagger.go \
		cmd/frontend/docs/frontend_swagger.json \
		cmd/frontend/docs/frontend_swagger.json \
		admin users
	@rm -f cmd/frontend/docs/docs.go cmd/frontend/docs/swagger.json

swag-backend:
	@swag init -g cmd/backend/main.go -o cmd/backend/docs --instanceName backend
	@go run scripts/filter_swagger.go \
		cmd/backend/docs/backend_swagger.json \
		cmd/backend/docs/backend_swagger.json \
		public user
	@rm -f cmd/backend/docs/docs.go cmd/backend/docs/swagger.json
```

### 4. 更新路由配置

在路由器中导入对应的 docs 包并指定实例名称：

**frontend.go**:
```go
import (
    _ "trx-project/cmd/frontend/docs" // 导入前台文档
    // ...
)

r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
    ginSwagger.InstanceName("frontend"))) // 使用 frontend 实例
```

**backend.go**:
```go
import (
    _ "trx-project/cmd/backend/docs" // 导入后台文档
    // ...
)

r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
    ginSwagger.InstanceName("backend"))) // 使用 backend 实例
```

## 实现效果

### 前台文档（http://localhost:8080/swagger/index.html）

**包含接口** (3个):
```
/public/login          - 用户登录
/public/register       - 用户注册
/user/profile          - 获取/更新个人信息
```

**排除接口**:
- `/admin/*` - 所有后台管理接口
- `/users` - 内部用户列表接口

### 后台文档（http://localhost:8081/swagger/index.html）

**包含接口** (12个):
```
/admin/rbac/permissions                - 权限管理
/admin/rbac/roles                      - 角色管理
/admin/rbac/roles/{id}                 - 角色详情/更新/删除
/admin/rbac/roles/{id}/permissions     - 角色权限关联
/admin/statistics/users                - 用户统计
/admin/users                           - 用户列表/创建
/admin/users/{id}                      - 用户详情/更新/删除
/admin/users/{id}/permissions          - 用户权限查询
/admin/users/{id}/reset-password       - 重置密码
/admin/users/{id}/role                 - 分配角色
/admin/users/{id}/roles                - 用户角色查询
/admin/users/{id}/status               - 用户状态管理
```

**排除接口**:
- `/public/*` - 所有前台公开接口
- `/user/*` - 前台用户接口

## 工作流程

1. **开发者添加新接口**
   - 在 handler 中添加 Swagger 注释
   - 在路由器中注册路由

2. **生成文档**
   ```bash
   make swag  # 或 make swag-frontend / make swag-backend
   ```

3. **自动处理**
   - swag 扫描代码生成完整文档
   - 过滤脚本移除不属于该服务的接口
   - 清理旧的文档文件

4. **访问文档**
   - 启动服务
   - 访问对应的 Swagger UI

## 技术细节

### 为什么不在 handler 层面分离？

**考虑过的方案**：
1. 将 handler 分成 `frontend/` 和 `backend/` 子目录
2. 使用 build tags

**为什么不采用**：
- 需要大规模重构现有代码
- handler 之间可能有共享的逻辑
- 过滤方案更简单、更灵活

### 为什么文档中有些接口路径不完整？

**现象**：
- Swagger 注释中的路径是 `/users/{id}`
- 但实际 URL 是 `/api/v1/admin/users/{id}`

**原因**：
- Swagger `@Router` 注释中的路径是相对于 `@BasePath` 的
- `@BasePath` 在 main.go 中定义为 `/api/v1`
- 路由组会添加额外的前缀（如 `/admin`）

**影响**：
- 不影响文档的正确性
- Swagger UI 会自动组合完整的 URL

### 过滤脚本的局限性

**基于路径前缀过滤**：
- 简单、高效
- 依赖于良好的路径命名规范

**不支持**：
- 基于标签（tags）过滤
- 基于文件名过滤
- 复杂的过滤规则

**未来改进方向**：
- 支持基于标签的过滤
- 支持正则表达式匹配
- 支持配置文件定义过滤规则

## 相关文件

- `scripts/filter_swagger.go` - 过滤脚本
- `Makefile` - 文档生成配置
- `internal/api/router/frontend.go` - 前台路由配置
- `internal/api/router/backend.go` - 后台路由配置
- `cmd/frontend/docs/` - 前台文档目录
- `cmd/backend/docs/` - 后台文档目录
- `docs/SWAGGER_GUIDE.md` - Swagger 使用指南

## 总结

通过 **Swagger 实例分离 + 后处理过滤** 的方案，我们成功实现了前后台 API 文档的完全分离，确保：

✅ 每个服务只展示自己的接口
✅ 文档生成流程自动化
✅ 开发体验良好（使用 Makefile 一键生成）
✅ 维护成本低（无需修改现有代码结构）

这个方案在不改变现有架构的前提下，优雅地解决了文档分离的问题。

