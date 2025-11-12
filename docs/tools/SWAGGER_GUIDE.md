# Swagger 文档使用指南

## 概述

本项目为前台和后台服务分别生成独立的 Swagger API 文档，确保每个服务只展示自己的接口。

## 文档特性

### 前台文档
- **路径**: `cmd/frontend/docs/`
- **访问地址**: http://localhost:8080/swagger/index.html
- **包含接口**: 
  - `/public/*` - 公开接口（注册、登录等）
  - `/user/*` - 用户接口（个人信息等）
- **排除接口**: 
  - `/admin/*` - 后台管理接口

### 后台文档
- **路径**: `cmd/backend/docs/`
- **访问地址**: http://localhost:8081/swagger/index.html
- **包含接口**: 
  - `/admin/*` - 后台管理接口
  - `/admin/rbac/*` - RBAC 权限管理接口
- **排除接口**: 
  - `/public/*` - 前台公开接口

## 生成文档

### 方式 1: 使用 Makefile（推荐）

```bash
# 生成前台文档
make swag-frontend

# 生成后台文档
make swag-backend

# 生成所有文档
make swag
```

### 方式 2: 使用跨平台脚本

适合没有 `make` 命令的环境（如 Windows）。

**Linux / macOS / Git Bash**:
```bash
./scripts/swagger.sh              # 生成所有文档
./scripts/swagger.sh frontend     # 只生成前台文档
./scripts/swagger.sh backend      # 只生成后台文档
```

**Windows CMD**:
```cmd
scripts\swagger.bat
scripts\swagger.bat frontend
scripts\swagger.bat backend
```

**PowerShell**:
```powershell
.\scripts\swagger.ps1
.\scripts\swagger.ps1 frontend
.\scripts\swagger.ps1 backend
```

💡 **脚本特性**:
- ✅ 自动检查 swag 是否安装
- ✅ 详细的进度输出
- ✅ 彩色状态提示
- ✅ 错误处理和提示

📖 详见：[Swagger 脚本使用指南](SWAGGER_SCRIPTS_GUIDE.md)

### 手动生成

#### 前台文档

```bash
# 1. 生成原始文档
swag init -g cmd/frontend/main.go -o cmd/frontend/docs \
  --parseDependency --parseInternal \
  --instanceName frontend

# 2. 过滤后台接口
go run scripts/filter_swagger.go \
  cmd/frontend/docs/frontend_swagger.json \
  cmd/frontend/docs/frontend_swagger.json \
  admin

# 3. 清理旧文件
rm -f cmd/frontend/docs/docs.go cmd/frontend/docs/swagger.json cmd/frontend/docs/swagger.yaml
```

#### 后台文档

```bash
# 1. 生成原始文档
swag init -g cmd/backend/main.go -o cmd/backend/docs \
  --parseDependency --parseInternal \
  --instanceName backend

# 2. 过滤前台接口
go run scripts/filter_swagger.go \
  cmd/backend/docs/backend_swagger.json \
  cmd/backend/docs/backend_swagger.json \
  public

# 3. 清理旧文件
rm -f cmd/backend/docs/docs.go cmd/backend/docs/swagger.json cmd/backend/docs/swagger.yaml
```

## 工作原理

### 1. Swagger 生成

使用 `swag` 工具扫描代码中的注释，生成 Swagger 文档：

```go
// @Summary 用户注册
// @Description 创建新用户账号
// @Tags 公开接口
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "注册信息"
// @Success 200 {object} response.Response "注册成功"
// @Router /public/register [post]
func (h *UserHandler) Register(c *gin.Context) {
    // ...
}
```

### 2. 实例名称

使用 `--instanceName` 参数为每个服务创建独立的 Swagger 实例：

- 前台: `--instanceName frontend`
- 后台: `--instanceName backend`

生成的文件：
- `frontend_docs.go`、`frontend_swagger.json`、`frontend_swagger.yaml`
- `backend_docs.go`、`backend_swagger.json`、`backend_swagger.yaml`

### 3. 过滤脚本

`scripts/filter_swagger.go` 用于从生成的文档中移除不属于该服务的接口：

```go
// 从前台文档移除 /admin 开头的路径
// 从后台文档移除 /public 开头的路径
```

### 4. 路由注册

在路由器中注册对应的 Swagger 实例：

```go
// frontend router
_ "trx-project/cmd/frontend/docs" // 导入前台文档

r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
    ginSwagger.InstanceName("frontend"))) // 使用 frontend 实例

// backend router
_ "trx-project/cmd/backend/docs" // 导入后台文档

r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
    ginSwagger.InstanceName("backend"))) // 使用 backend 实例
```

## 访问文档

### 启动服务

```bash
# 启动前台服务（端口 8080）
go run cmd/frontend/main.go

# 启动后台服务（端口 8081）
go run cmd/backend/main.go
```

### 访问 Swagger UI

- **前台文档**: http://localhost:8080/swagger/index.html
- **后台文档**: http://localhost:8081/swagger/index.html

### 使用 Token 认证

1. 在 Swagger UI 中点击右上角的 "Authorize" 按钮
2. 输入格式：`Bearer <your-jwt-token>`
3. 示例：`Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...`

## 添加新接口

### 1. 编写 Handler 并添加注释

```go
// @Summary 获取用户列表
// @Description 获取所有用户列表（需要管理员权限）
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param limit query int false "每页数量" default(10)
// @Success 200 {object} response.Response{data=[]model.User} "成功获取用户列表"
// @Failure 401 {object} response.Response "未授权"
// @Router /admin/users [get]
func (h *AdminUserHandler) GetUsers(c *gin.Context) {
    // ...
}
```

### 2. 重新生成文档

```bash
# 如果是后台接口
make swag-backend

# 如果是前台接口
make swag-frontend
```

### 3. 重启服务

重启对应的服务以加载新的文档。

## 常见问题

### Q: 为什么前后台文档都包含所有接口？

A: swag 工具会扫描整个项目的所有 handler。我们使用 `filter_swagger.go` 脚本在生成后自动过滤不属于该服务的接口。

### Q: 如何只显示特定标签的接口？

A: 在 Swagger UI 中，可以使用标签过滤功能。或者修改 `filter_swagger.go` 脚本添加基于标签的过滤。

### Q: 文档没有更新怎么办？

A: 确保：
1. 重新生成了文档（`make swag`）
2. 重启了服务
3. 清除了浏览器缓存

### Q: Git Bash 中路径转换问题

A: `filter_swagger.go` 已处理这个问题。使用不带 `/` 的路径前缀（如 `admin` 而非 `/admin`），脚本会自动添加前缀 `/`。

## 文档注释规范

### 必需标签

```go
// @Summary 接口简短描述（必需）
// @Description 接口详细描述（可选）
// @Tags 标签名称（必需，用于分组）
// @Accept json（必需）
// @Produce json（必需）
// @Router /path [method]（必需）
```

### 参数标签

```go
// @Param name location type required "description" [options]

// 路径参数
// @Param id path int true "用户ID"

// 查询参数
// @Param page query int false "页码" default(1)

// 请求体
// @Param request body RegisterRequest true "注册信息"

// Header
// @Param Authorization header string true "Bearer Token"
```

### 响应标签

```go
// @Success code {type} model "description"
// @Success 200 {object} response.Response "成功"
// @Success 200 {object} response.Response{data=[]model.User} "成功获取用户列表"

// @Failure code {type} model "description"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
```

### 认证标签

```go
// @Security BearerAuth  // 需要 JWT 认证
```

## 最佳实践

1. **保持注释更新**: 修改接口时同步更新注释
2. **使用清晰的标签**: 合理分组接口，方便查找
3. **详细的描述**: 提供足够的信息帮助使用者理解接口
4. **完整的参数说明**: 说明每个参数的类型、是否必需、默认值等
5. **示例数据**: 在可能的情况下提供示例请求和响应
6. **错误码说明**: 列出所有可能的错误响应

## 相关文件

- `scripts/filter_swagger.go` - Swagger 过滤脚本
- `internal/api/router/frontend.go` - 前台路由配置
- `internal/api/router/backend.go` - 后台路由配置
- `cmd/frontend/docs/` - 前台 Swagger 文档目录
- `cmd/backend/docs/` - 后台 Swagger 文档目录

## 参考资料

- [Swag 官方文档](https://github.com/swaggo/swag)
- [Swagger 规范](https://swagger.io/specification/)
- [Gin Swagger 中间件](https://github.com/swaggo/gin-swagger)

