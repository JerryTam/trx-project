# 🎯 Handler 重构 - 快速开始

## 📋 重构概览

**之前**: 复杂的路径过滤 → **现在**: 简单的目录分离 ✅

## 🏗️ 新的目录结构

```
internal/api/handler/
├── frontend/                    # ✅ 前台 handler
│   └── user_handler.go         # 用户相关接口
└── backend/                     # ✅ 后台 handler
    ├── admin_user_handler.go   # 管理员用户接口
    └── rbac_handler.go         # RBAC 权限接口
```

## 🚀 生成 Swagger 文档

### 方式 1: 使用脚本（推荐）

```bash
# 生成所有文档
./scripts/swagger.sh

# 只生成前台
./scripts/swagger.sh frontend

# 只生成后台
./scripts/swagger.sh backend
```

### 方式 2: 使用 Makefile

```bash
make swag              # 生成所有
make swag-frontend     # 只生成前台
make swag-backend      # 只生成后台
```

## ✨ 核心优势

### 之前（复杂）
```bash
# 1. 生成文档
swag init ...

# 2. 运行过滤脚本
go run scripts/filter_swagger_docs.go ... admin users public

# 3. 需要维护过滤规则列表
```

### 现在（简单）
```bash
# 一条命令搞定！
swag init -g cmd/frontend/main.go ... \
    --exclude internal/api/handler/backend
```

## 📊 验证结果

```bash
# 前台只包含前台接口
$ grep -o '"\/[^"]*":' cmd/frontend/docs/frontend_swagger.json
"/public/login"
"/public/register"
"/user/profile"
"/users"
"/users/{id}"

# 后台只包含后台接口
$ grep -o '"\/[^"]*":' cmd/backend/docs/backend_swagger.json
"/admin/rbac/permissions"
"/admin/rbac/roles"
... (共 12 个)
```

## 💡 添加新 Handler

### 前台 Handler

```bash
# 1. 在前台目录创建文件
vim internal/api/handler/frontend/my_handler.go

# 2. 定义 handler
package frontend

type MyHandler struct {
    // ...
}

func NewMyHandler() *MyHandler {
    return &MyHandler{}
}

// 3. 添加 Swagger 注释
// @Summary 我的接口
// @Router /my/path [get]
func (h *MyHandler) MyMethod(c *gin.Context) {
    // ...
}

# 4. 重新生成文档
./scripts/swagger.sh frontend

# ✅ 完成！无需修改任何过滤规则
```

### 后台 Handler

```bash
# 1. 在后台目录创建文件
vim internal/api/handler/backend/my_admin_handler.go

# 2. 按相同方式实现
package backend

# 3. 重新生成
./scripts/swagger.sh backend

# ✅ 完成！
```

## 🔧 编译和运行

```bash
# 编译
go build ./cmd/frontend
go build ./cmd/backend

# 运行
go run cmd/frontend/main.go  # 8080
go run cmd/backend/main.go   # 8081

# 访问 Swagger UI
# 前台: http://localhost:8080/swagger/index.html
# 后台: http://localhost:8081/swagger/index.html
```

## 📚 相关文档

- [详细重构说明](HANDLER_REFACTOR_SUMMARY.md)
- [Swagger 使用指南](docs/SWAGGER_GUIDE.md)
- [项目 README](README.md)

## ❓ 常见问题

### Q: 旧的过滤脚本还能用吗？

A: 已废弃，请使用新架构。参见 [scripts/DEPRECATED.md](scripts/DEPRECATED.md)

### Q: 如何确定 handler 放在哪个目录？

A:
- **frontend/** - 面向最终用户的接口（注册、登录、个人信息等）
- **backend/** - 面向管理员的接口（用户管理、权限管理、统计等）

### Q: 可以添加更多目录吗？

A: 当然！比如：
- `internal/api/handler/mobile/` - 移动端专用接口
- `internal/api/handler/api/v2/` - API v2 版本
- `internal/api/handler/webhook/` - Webhook 接口

然后在 Swagger 生成时用 `--exclude` 排除其他目录即可。

## 🎉 总结

- ✅ **目录结构清晰** - 一眼看出前后台分离
- ✅ **配置简单** - 一行 `--exclude` 搞定
- ✅ **易于维护** - 新增 handler 无需改配置
- ✅ **完全分离** - Swagger 文档100%独立

开始使用新架构吧！🚀

