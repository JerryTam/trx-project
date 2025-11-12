# 功能开发快速参考卡片

> 💡 这是 [功能模块开发规范](FEATURE_DEVELOPMENT_GUIDE.md) 的快速参考版本

## ⚡ 10 分钟快速上手

### 1️⃣ 数据库迁移 (2分钟)

```bash
# 方式 1: 使用脚本自动创建迁移文件（⭐ 推荐）
NAME=create_{table}_table ./scripts/migrate.sh create

# 方式 2: 使用 Makefile
make migrate-create NAME=create_{table}_table

# 编辑生成的迁移文件
vim migrations/000XXX_create_{table}_table.up.sql
vim migrations/000XXX_create_{table}_table.down.sql

# 执行迁移
./scripts/migrate.sh up
# 或
make migrate-up
```

### 2️⃣ Model 层 (1分钟)

```bash
# 创建文件
vim internal/model/{feature}.go
```

```go
package model

type {Feature} struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	Name      string     `json:"name"`
	Status    int        `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func ({Feature}) TableName() string {
	return "{table_name}"
}
```

### 3️⃣ Repository 层 (2分钟)

```bash
vim internal/repository/{feature}_repository.go
```

```go
package repository

type {Feature}Repository interface {
	Create(item *model.{Feature}) error
	GetByID(id uint) (*model.{Feature}, error)
	GetList(page, pageSize int) ([]model.{Feature}, int64, error)
	Update(item *model.{Feature}) error
	Delete(id uint) error
}

type {feature}Repository struct {
	db *gorm.DB
}

func New{Feature}Repository(db *gorm.DB) {Feature}Repository {
	return &{feature}Repository{db: db}
}

// 实现所有接口方法...
```

### 4️⃣ Service 层 (2分钟)

```bash
vim internal/service/{feature}_service.go
```

```go
package service

type {Feature}Service interface {
	Create{Feature}(name string) (*model.{Feature}, error)
	Get{Feature}ByID(id uint) (*model.{Feature}, error)
	// ... 其他方法
}

type {feature}Service struct {
	{feature}Repo repository.{Feature}Repository
	logger        *zap.Logger
}

func New{Feature}Service(
	{feature}Repo repository.{Feature}Repository,
	logger *zap.Logger,
) {Feature}Service {
	return &{feature}Service{
		{feature}Repo: {feature}Repo,
		logger:        logger,
	}
}

// 实现所有接口方法...
```

### 5️⃣ Handler 层 (2分钟)

```bash
vim internal/api/handler/frontendHandler/{feature}_handler.go
```

```go
package frontendHandler

import _ "trx-project/internal/model" // 必须导入！

type {Feature}Handler struct {
	{feature}Service service.{Feature}Service
	logger           *zap.Logger
}

func New{Feature}Handler(
	{feature}Service service.{Feature}Service,
	logger *zap.Logger,
) *{Feature}Handler {
	return &{Feature}Handler{
		{feature}Service: {feature}Service,
		logger:           logger,
	}
}

// @Summary 创建{功能}
// @Tags {功能}管理
// @Security BearerAuth
// @Router /user/{features} [post]
func (h *{Feature}Handler) Create{Feature}(c *gin.Context) {
	// ... 实现
}
```

### 6️⃣ 注册路由 (1分钟)

```go
// internal/api/router/frontend.go
func SetupFrontend(
	// ... 现有参数
	{feature}Handler *frontendHandler.{Feature}Handler, // 添加
	// ... 其他参数
) *gin.Engine {
	// ...
	user := v1.Group("/user")
	{
		user.POST("/{features}", {feature}Handler.Create{Feature})
		user.GET("/{features}", {feature}Handler.Get{Feature}List)
		user.GET("/{features}/:id", {feature}Handler.Get{Feature})
		// ...
	}
}
```

### 7️⃣ 配置 Wire (1分钟)

```go
// cmd/frontend/wire.go
wire.Build(
	// ... 现有依赖
	repository.New{Feature}Repository,   // 添加
	service.New{Feature}Service,         // 添加
	frontendHandler.New{Feature}Handler, // 添加
	// ...
)

// cmd/frontend/providers.go
func provideFrontendRouter(
	// ... 现有参数
	{feature}Handler *frontendHandler.{Feature}Handler, // 添加
	// ... 其他参数
) *gin.Engine {
	return router.SetupFrontend(
		// ... 现有参数
		{feature}Handler, // 添加
		// ...
	)
}
```

### 8️⃣ 生成与编译 (1分钟)

```bash
# 生成 Wire 代码
cd cmd/frontend && wire && cd ../..

# 生成 Swagger 文档
make swag-frontend

# 编译
go build -o bin/frontend cmd/frontend/*.go
```

---

## 📝 替换规则

在使用模板时，按以下规则替换：

| 占位符 | 说明 | 示例 |
|--------|------|------|
| `{Feature}` | 大驼峰功能名 | `Order`, `Product`, `Comment` |
| `{feature}` | 小驼峰功能名 | `order`, `product`, `comment` |
| `{features}` | 复数URL路径 | `orders`, `products`, `comments` |
| `{table_name}` | 数据库表名 | `orders`, `products`, `comments` |
| `{功能}` | 中文功能名 | `订单`, `商品`, `评论` |
| `{number}` | 迁移序号 | `000007`, `000008` |

---

## 🎯 检查清单

创建新功能前，快速检查：

### 文件创建
- [ ] `migrations/{number}_create_{table}_table.up.sql`
- [ ] `migrations/{number}_create_{table}_table.down.sql`
- [ ] `internal/model/{feature}.go`
- [ ] `internal/repository/{feature}_repository.go`
- [ ] `internal/service/{feature}_service.go`
- [ ] `internal/api/handler/{frontend|backend}Handler/{feature}_handler.go`
- [ ] `scripts/test_{feature}_api.sh`

### 代码修改
- [ ] `internal/api/router/frontend.go` 或 `backend.go`
- [ ] `cmd/frontend/wire.go` 或 `backend/wire.go`
- [ ] `cmd/frontend/providers.go` 或 `backend/providers.go`

### 执行命令
- [ ] `go run cmd/migrate/main.go -cmd up`
- [ ] `cd cmd/frontend && wire`
- [ ] `make swag-frontend`
- [ ] `go build -o bin/frontend cmd/frontend/*.go`
- [ ] `./scripts/test_{feature}_api.sh`

---

## 🚀 常用命令

```bash
# 数据库迁移
go run cmd/migrate/main.go -cmd up
go run cmd/migrate/main.go -cmd version

# Wire 生成
cd cmd/frontend && wire
cd cmd/backend && wire

# Swagger 生成
make swag                # 全部
make swag-frontend       # 前台
make swag-backend        # 后台

# 编译
go build -o bin/frontend cmd/frontend/*.go
go build -o bin/backend cmd/backend/*.go

# 运行
./bin/frontend
./bin/backend

# 测试
./scripts/test_{feature}_api.sh
```

---

## 💡 常见错误与解决

### Wire 生成失败

```bash
# 检查语法
go build ./cmd/frontend/wire.go

# 查看错误
cd cmd/frontend && wire 2>&1
```

**常见原因:**
- Import 路径错误
- 参数类型不匹配
- 忘记添加到 wire.Build

### Swagger 找不到类型

```go
// Handler 文件必须导入 model
import _ "trx-project/internal/model"
```

### 编译失败

```bash
# 更新依赖
go mod tidy

# 查看详细错误
go build -v ./cmd/frontend
```

---

## 📖 更多信息

详细说明请查看: [功能模块开发规范完整版](FEATURE_DEVELOPMENT_GUIDE.md)

---

**文档版本:** v1.0  
**更新时间:** 2025-01-12

