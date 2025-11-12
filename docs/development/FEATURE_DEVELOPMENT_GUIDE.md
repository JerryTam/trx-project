# 功能模块开发规范 - 完整指南

## 📋 文档说明

本文档提供了在 TRX Project 中开发新功能模块的完整流程和规范，包括详细的步骤说明、代码模板和最佳实践。

**适用场景**: 任何需要添加新的业务功能模块

**参考案例**: 用户订单管理功能

## 🎯 开发流程总览

```
1. 需求分析与设计
   ├── 功能规划
   ├── 数据库设计
   └── API 接口设计

2. 数据库层
   ├── 创建迁移文件
   └── 执行数据库迁移

3. 代码实现（分层架构）
   ├── Model 层（数据模型）
   ├── Repository 层（数据访问）
   ├── Service 层（业务逻辑）
   └── Handler 层（HTTP 处理）

4. 路由与依赖注入
   ├── 注册路由
   └── 配置 Wire

5. 文档与测试
   ├── 生成 Swagger 文档
   ├── 编写测试脚本
   └── 功能测试

6. 代码审查与上线
   ├── 代码审查
   ├── 性能测试
   └── 部署上线
```

---

## 步骤 1: 需求分析与设计

### 1.1 功能规划

**需要明确的内容:**

1. **功能目标** - 这个功能要解决什么问题？
2. **用户角色** - 前台用户 / 后台管理员 / 两者都有？
3. **核心功能** - 列出所有需要实现的功能点
4. **业务规则** - 状态流转、权限控制、数据验证等

**示例（订单管理）:**

```
功能目标: 用户可以创建、查询、支付、取消和完成订单
用户角色: 前台用户
核心功能:
  - 创建订单
  - 查询订单列表（分页）
  - 查询订单详情
  - 支付订单
  - 取消订单
  - 完成订单

业务规则:
  - 订单状态: 待支付 → 已支付 → 已完成
  - 待支付的订单可以取消
  - 只有已支付的订单才能完成
  - 用户只能操作自己的订单
```

### 1.2 数据库设计

**设计要点:**

1. **表结构设计** - 字段定义、数据类型、默认值
2. **索引设计** - 主键、唯一索引、普通索引
3. **关联关系** - 外键、一对多、多对多
4. **软删除** - 是否需要 `deleted_at` 字段

**模板（以订单表为例）:**

```sql
CREATE TABLE IF NOT EXISTS `orders` (
  -- 主键
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '订单ID',
  
  -- 业务字段
  `order_no` varchar(32) NOT NULL COMMENT '订单号',
  `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
  `product_name` varchar(255) NOT NULL COMMENT '商品名称',
  `product_price` decimal(10,2) NOT NULL COMMENT '商品单价',
  `quantity` int NOT NULL DEFAULT 1 COMMENT '购买数量',
  `total_amount` decimal(10,2) NOT NULL COMMENT '订单总金额',
  `status` tinyint NOT NULL DEFAULT 0 COMMENT '状态: 0-待处理, 1-处理中, 2-已完成',
  `remark` text COMMENT '备注',
  
  -- 时间字段（必备）
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  
  -- 索引
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_order_no` (`order_no`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='订单表';
```

**字段命名规范:**

- 使用小写字母和下划线（snake_case）
- ID 字段: `id`, `user_id`, `order_id`
- 时间字段: `created_at`, `updated_at`, `deleted_at`, `paid_at`
- 状态字段: `status`, `is_active`, `is_deleted`
- 金额字段: 使用 `decimal(10,2)` 类型

**索引设计原则:**

- 主键索引: `PRIMARY KEY`
- 唯一索引: 业务唯一字段（如订单号、手机号）
- 普通索引: 经常查询的字段（如用户ID、状态、时间）
- 联合索引: 经常一起查询的字段组合

### 1.3 API 接口设计

**接口规划表格:**

| 功能 | 方法 | 路径 | 请求参数 | 响应数据 | 权限 |
|------|------|------|----------|----------|------|
| 创建XX | POST | `/api/v1/user/xxx` | JSON Body | 创建的对象 | 用户认证 |
| 列表查询 | GET | `/api/v1/user/xxx` | page, page_size | 分页数据 | 用户认证 |
| 详情查询 | GET | `/api/v1/user/xxx/:id` | - | 对象详情 | 用户认证 |
| 更新XX | PUT | `/api/v1/user/xxx/:id` | JSON Body | 更新后的对象 | 用户认证 |
| 删除XX | DELETE | `/api/v1/user/xxx/:id` | - | 成功消息 | 用户认证 |

**RESTful API 设计原则:**

- 使用标准 HTTP 方法（GET, POST, PUT, DELETE, PATCH）
- URL 使用复数名词（`/orders` 而不是 `/order`）
- 使用路径参数表示资源 ID（`/orders/:id`）
- 使用查询参数表示过滤和分页（`?page=1&status=paid`）
- 返回统一的响应格式

**响应格式示例:**

```json
{
  "code": 200,
  "message": "success",
  "data": { ... },
  "timestamp": 1699999999
}
```

---

## 步骤 2: 数据库层实现

### 2.1 创建迁移文件

**推荐方式：使用脚本自动创建（⭐ 推荐）**

```bash
# 方式 1: 使用 migrate.sh 脚本
NAME=create_orders_table ./scripts/migrate.sh create

# 方式 2: 使用 Makefile
make migrate-create NAME=create_orders_table
```

**脚本会自动：**
- ✅ 计算下一个版本号（自动递增）
- ✅ 创建 `up.sql` 和 `down.sql` 两个文件
- ✅ 添加基础模板和注释
- ✅ 使用正确的命名格式

**输出示例:**
```
✅ 迁移文件创建成功:
  📄 migrations/000007_create_orders_table.up.sql
  📄 migrations/000007_create_orders_table.down.sql
```

**手动创建方式（不推荐）:**

如果需要手动创建，命名规范如下：

```
{序号}_{操作}_{表名}.{up|down}.sql

示例:
000007_create_orders_table.up.sql
000007_create_orders_table.down.sql
```

**查看当前最新序号:**

```bash
cd /d/workspace/go/trx-project
ls migrations/*.up.sql | wc -l
```

### 2.2 编写迁移文件

**如果使用脚本创建，文件已经自动生成，只需要编辑 SQL 内容即可。**

**UP 迁移文件模板 (`{number}_create_{table}_table.up.sql`):**

```sql
-- create_{table}_table - 升级脚本
-- 创建时间: 2025-01-12 17:00:00

-- 创建{表名}表
CREATE TABLE IF NOT EXISTS `{table_name}` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  
  -- 业务字段
  `name` varchar(255) NOT NULL COMMENT '名称',
  `status` tinyint NOT NULL DEFAULT 0 COMMENT '状态',
  
  -- 时间字段
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  
  PRIMARY KEY (`id`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='{表名}表';
```

**DOWN 迁移文件模板 (`{number}_create_{table}_table.down.sql`):**

```sql
-- create_{table}_table - 回滚脚本
-- 创建时间: 2025-01-12 17:00:00

-- 删除{表名}表
DROP TABLE IF EXISTS `{table_name}`;
```

### 2.3 执行迁移

```bash
# 查看当前迁移版本
./scripts/migrate.sh version
# 或
make migrate-version

# 执行迁移（向上）
./scripts/migrate.sh up
# 或
make migrate-up

# 回滚迁移（向下）
./scripts/migrate.sh down
# 或
make migrate-down

# 强制设置版本（谨慎使用）
VERSION=7 ./scripts/migrate.sh force
# 或
make migrate-force VERSION=7
```

**验证迁移:**

```bash
# 方式 1: 查看迁移版本
./scripts/migrate.sh version
# 或
make migrate-version

# 方式 2: 直接查看数据库
mysql -u root -p trx_db -e "SHOW TABLES;"
mysql -u root -p trx_db -e "DESC orders;"
```

---

## 步骤 3: 代码实现（分层架构）

### 3.1 Model 层 - 数据模型

**文件位置:** `internal/model/{feature}.go`

**代码模板:**

```go
package model

import (
	"time"
	"gorm.io/gorm"
)

// {Feature}Status {功能}状态枚举
type {Feature}Status int

const (
	{Feature}StatusPending   {Feature}Status = 0 // 待处理
	{Feature}StatusProcessing {Feature}Status = 1 // 处理中
	{Feature}StatusCompleted {Feature}Status = 2 // 已完成
)

// {Feature}StatusText 状态文本映射
var {Feature}StatusText = map[{Feature}Status]string{
	{Feature}StatusPending:   "待处理",
	{Feature}StatusProcessing: "处理中",
	{Feature}StatusCompleted: "已完成",
}

// {Feature} {功能}模型
type {Feature} struct {
	ID          uint              `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string            `gorm:"type:varchar(255);not null" json:"name"`
	Status      {Feature}Status   `gorm:"type:tinyint;index:idx_status;not null;default:0" json:"status"`
	StatusText  string            `gorm:"-" json:"status_text"` // 不存储到数据库
	Remark      string            `gorm:"type:text" json:"remark"`
	CreatedAt   time.Time         `gorm:"index:idx_created_at" json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	DeletedAt   gorm.DeletedAt    `gorm:"index" json:"-"`
}

// TableName 指定表名
func ({Feature}) TableName() string {
	return "{table_name}"
}

// AfterFind GORM 钩子：查询后设置状态文本
func (f *{Feature}) AfterFind(tx *gorm.DB) error {
	f.StatusText = {Feature}StatusText[f.Status]
	return nil
}

// BeforeCreate GORM 钩子：创建前的处理
func (f *{Feature}) BeforeCreate(tx *gorm.DB) error {
	// 如果需要自动生成编号等逻辑，在这里实现
	return nil
}
```

**关键点:**

1. **字段标签**:
   - `gorm` 标签: 定义数据库字段属性
   - `json` 标签: 定义 JSON 序列化字段名

2. **GORM 钩子**:
   - `BeforeCreate`: 创建前执行
   - `AfterFind`: 查询后执行
   - `BeforeUpdate`: 更新前执行
   - `BeforeDelete`: 删除前执行

3. **状态枚举**: 使用 const 定义状态常量，便于维护

### 3.2 Repository 层 - 数据访问

**文件位置:** `internal/repository/{feature}_repository.go`

**代码模板:**

```go
package repository

import (
	"trx-project/internal/model"
	"gorm.io/gorm"
)

// {Feature}Repository {功能}仓储接口
type {Feature}Repository interface {
	Create(item *model.{Feature}) error
	GetByID(id uint) (*model.{Feature}, error)
	GetList(page, pageSize int) ([]model.{Feature}, int64, error)
	Update(item *model.{Feature}) error
	UpdateStatus(id uint, status model.{Feature}Status) error
	Delete(id uint) error
}

type {feature}Repository struct {
	db *gorm.DB
}

// New{Feature}Repository 创建{功能}仓储实例
func New{Feature}Repository(db *gorm.DB) {Feature}Repository {
	return &{feature}Repository{
		db: db,
	}
}

// Create 创建{功能}
func (r *{feature}Repository) Create(item *model.{Feature}) error {
	return r.db.Create(item).Error
}

// GetByID 根据ID获取{功能}
func (r *{feature}Repository) GetByID(id uint) (*model.{Feature}, error) {
	var item model.{Feature}
	err := r.db.First(&item, id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// GetList 获取{功能}列表（分页）
func (r *{feature}Repository) GetList(page, pageSize int) ([]model.{Feature}, int64, error) {
	var items []model.{Feature}
	var total int64

	// 计算总数
	if err := r.db.Model(&model.{Feature}{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := r.db.Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&items).Error

	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// Update 更新{功能}
func (r *{feature}Repository) Update(item *model.{Feature}) error {
	return r.db.Save(item).Error
}

// UpdateStatus 更新{功能}状态
func (r *{feature}Repository) UpdateStatus(id uint, status model.{Feature}Status) error {
	return r.db.Model(&model.{Feature}{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// Delete 删除{功能}（软删除）
func (r *{feature}Repository) Delete(id uint) error {
	return r.db.Delete(&model.{Feature}{}, id).Error
}
```

**Repository 层职责:**

- ✅ 数据库 CRUD 操作
- ✅ 查询条件构建
- ✅ 分页查询
- ✅ 事务处理（如需要）
- ❌ 业务逻辑（应该在 Service 层）
- ❌ 数据验证（应该在 Service 层）

**常用查询方法:**

```go
// 条件查询
func (r *{feature}Repository) GetByCondition(condition map[string]interface{}) ([]model.{Feature}, error) {
	var items []model.{Feature}
	err := r.db.Where(condition).Find(&items).Error
	return items, err
}

// 根据用户ID查询
func (r *{feature}Repository) GetByUserID(userID uint, page, pageSize int) ([]model.{Feature}, int64, error) {
	var items []model.{Feature}
	var total int64

	query := r.db.Model(&model.{Feature}{}).Where("user_id = ?", userID)
	
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&items).Error

	return items, total, err
}

// 批量创建
func (r *{feature}Repository) BatchCreate(items []model.{Feature}) error {
	return r.db.Create(&items).Error
}
```

### 3.3 Service 层 - 业务逻辑

**文件位置:** `internal/service/{feature}_service.go`

**代码模板:**

```go
package service

import (
	"errors"
	"trx-project/internal/model"
	"trx-project/internal/repository"
	
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// {Feature}Service {功能}服务接口
type {Feature}Service interface {
	Create{Feature}(name string, remark string) (*model.{Feature}, error)
	Get{Feature}ByID(id uint) (*model.{Feature}, error)
	Get{Feature}List(page, pageSize int) ([]model.{Feature}, int64, error)
	Update{Feature}(id uint, name string, remark string) (*model.{Feature}, error)
	Delete{Feature}(id uint) error
	Change{Feature}Status(id uint, status model.{Feature}Status) error
}

type {feature}Service struct {
	{feature}Repo repository.{Feature}Repository
	logger        *zap.Logger
}

// New{Feature}Service 创建{功能}服务实例
func New{Feature}Service(
	{feature}Repo repository.{Feature}Repository,
	logger *zap.Logger,
) {Feature}Service {
	return &{feature}Service{
		{feature}Repo: {feature}Repo,
		logger:        logger,
	}
}

// Create{Feature} 创建{功能}
func (s *{feature}Service) Create{Feature}(name string, remark string) (*model.{Feature}, error) {
	// 1. 参数验证
	if name == "" {
		return nil, errors.New("名称不能为空")
	}

	// 2. 创建对象
	item := &model.{Feature}{
		Name:   name,
		Status: model.{Feature}StatusPending,
		Remark: remark,
	}

	// 3. 保存到数据库
	if err := s.{feature}Repo.Create(item); err != nil {
		s.logger.Error("创建{功能}失败",
			zap.String("name", name),
			zap.Error(err))
		return nil, err
	}

	// 4. 记录日志
	s.logger.Info("{功能}创建成功",
		zap.Uint("id", item.ID),
		zap.String("name", name))

	return item, nil
}

// Get{Feature}ByID 根据ID获取{功能}
func (s *{feature}Service) Get{Feature}ByID(id uint) (*model.{Feature}, error) {
	item, err := s.{feature}Repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("{功能}不存在")
		}
		s.logger.Error("获取{功能}失败", zap.Uint("id", id), zap.Error(err))
		return nil, err
	}
	return item, nil
}

// Get{Feature}List 获取{功能}列表
func (s *{feature}Service) Get{Feature}List(page, pageSize int) ([]model.{Feature}, int64, error) {
	// 参数验证
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	items, total, err := s.{feature}Repo.GetList(page, pageSize)
	if err != nil {
		s.logger.Error("获取{功能}列表失败", zap.Error(err))
		return nil, 0, err
	}

	return items, total, nil
}

// Update{Feature} 更新{功能}
func (s *{feature}Service) Update{Feature}(id uint, name string, remark string) (*model.{Feature}, error) {
	// 1. 获取原数据
	item, err := s.Get{Feature}ByID(id)
	if err != nil {
		return nil, err
	}

	// 2. 参数验证
	if name == "" {
		return nil, errors.New("名称不能为空")
	}

	// 3. 更新字段
	item.Name = name
	item.Remark = remark

	// 4. 保存到数据库
	if err := s.{feature}Repo.Update(item); err != nil {
		s.logger.Error("更新{功能}失败",
			zap.Uint("id", id),
			zap.Error(err))
		return nil, err
	}

	s.logger.Info("{功能}更新成功", zap.Uint("id", id))
	return item, nil
}

// Delete{Feature} 删除{功能}
func (s *{feature}Service) Delete{Feature}(id uint) error {
	// 1. 检查是否存在
	if _, err := s.Get{Feature}ByID(id); err != nil {
		return err
	}

	// 2. 执行删除
	if err := s.{feature}Repo.Delete(id); err != nil {
		s.logger.Error("删除{功能}失败",
			zap.Uint("id", id),
			zap.Error(err))
		return err
	}

	s.logger.Info("{功能}删除成功", zap.Uint("id", id))
	return nil
}

// Change{Feature}Status 更改{功能}状态
func (s *{feature}Service) Change{Feature}Status(id uint, status model.{Feature}Status) error {
	// 1. 获取原数据
	item, err := s.Get{Feature}ByID(id)
	if err != nil {
		return err
	}

	// 2. 状态验证（根据业务规则）
	// 例如：只有待处理状态才能更改为处理中
	if item.Status == model.{Feature}StatusCompleted {
		return errors.New("已完成的{功能}不能更改状态")
	}

	// 3. 更新状态
	if err := s.{feature}Repo.UpdateStatus(id, status); err != nil {
		s.logger.Error("更改{功能}状态失败",
			zap.Uint("id", id),
			zap.Int("status", int(status)),
			zap.Error(err))
		return err
	}

	s.logger.Info("{功能}状态更改成功",
		zap.Uint("id", id),
		zap.Int("old_status", int(item.Status)),
		zap.Int("new_status", int(status)))

	return nil
}
```

**Service 层职责:**

- ✅ 业务逻辑处理
- ✅ 参数验证
- ✅ 权限检查
- ✅ 事务控制
- ✅ 日志记录
- ✅ 调用 Repository 层
- ❌ HTTP 请求处理（应该在 Handler 层）
- ❌ 直接操作数据库（应该通过 Repository）

### 3.4 Handler 层 - HTTP 处理

**文件位置:** `internal/api/handler/frontendHandler/{feature}_handler.go` (前台) 或 `internal/api/handler/backendHandler/{feature}_handler.go` (后台)

**代码模板:**

```go
package frontendHandler // 或 backendHandler

import (
	"strconv"
	_ "trx-project/internal/model" // 用于 Swagger 文档生成
	"trx-project/internal/service"
	"trx-project/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// {Feature}Handler {功能}处理器
type {Feature}Handler struct {
	{feature}Service service.{Feature}Service
	logger           *zap.Logger
}

// New{Feature}Handler 创建{功能}处理器实例
func New{Feature}Handler(
	{feature}Service service.{Feature}Service,
	logger *zap.Logger,
) *{Feature}Handler {
	return &{Feature}Handler{
		{feature}Service: {feature}Service,
		logger:           logger,
	}
}

// Create{Feature}Request 创建{功能}请求
type Create{Feature}Request struct {
	Name   string `json:"name" binding:"required" example:"示例名称"`
	Remark string `json:"remark" example:"备注信息"`
}

// Create{Feature} 创建{功能}
// @Summary 创建{功能}
// @Description 创建新的{功能}
// @Tags {功能}管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body Create{Feature}Request true "{功能}信息"
// @Success 200 {object} response.Response{data=model.{Feature}} "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未登录"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /user/{features} [post]
func (h *{Feature}Handler) Create{Feature}(c *gin.Context) {
	// 1. 获取当前登录用户ID（如果需要）
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "请先登录")
		return
	}

	// 2. 绑定请求参数
	var req Create{Feature}Request
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	// 3. 调用服务层
	item, err := h.{feature}Service.Create{Feature}(req.Name, req.Remark)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	// 4. 返回成功响应
	response.Success(c, "{功能}创建成功", item)
}

// Get{Feature} 获取{功能}详情
// @Summary 获取{功能}详情
// @Description 根据ID获取{功能}详情
// @Tags {功能}管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "{功能}ID"
// @Success 200 {object} response.Response{data=model.{Feature}} "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未登录"
// @Failure 404 {object} response.Response "{功能}不存在"
// @Router /user/{features}/{id} [get]
func (h *{Feature}Handler) Get{Feature}(c *gin.Context) {
	// 1. 获取ID参数
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "ID格式错误")
		return
	}

	// 2. 调用服务层
	item, err := h.{feature}Service.Get{Feature}ByID(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	// 3. 返回成功响应
	response.Success(c, "获取成功", item)
}

// Get{Feature}List 获取{功能}列表
// @Summary 获取{功能}列表
// @Description 获取{功能}列表（分页）
// @Tags {功能}管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} response.Response{data=response.PaginationData} "获取成功"
// @Failure 401 {object} response.Response "未登录"
// @Router /user/{features} [get]
func (h *{Feature}Handler) Get{Feature}List(c *gin.Context) {
	// 1. 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// 2. 调用服务层
	items, total, err := h.{feature}Service.Get{Feature}List(page, pageSize)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	// 3. 返回分页数据
	response.Pagination(c, "获取成功", items, total, page, pageSize)
}

// Update{Feature}Request 更新{功能}请求
type Update{Feature}Request struct {
	Name   string `json:"name" binding:"required" example:"新名称"`
	Remark string `json:"remark" example:"新备注"`
}

// Update{Feature} 更新{功能}
// @Summary 更新{功能}
// @Description 更新{功能}信息
// @Tags {功能}管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "{功能}ID"
// @Param body body Update{Feature}Request true "{功能}信息"
// @Success 200 {object} response.Response{data=model.{Feature}} "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未登录"
// @Failure 404 {object} response.Response "{功能}不存在"
// @Router /user/{features}/{id} [put]
func (h *{Feature}Handler) Update{Feature}(c *gin.Context) {
	// 1. 获取ID参数
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "ID格式错误")
		return
	}

	// 2. 绑定请求参数
	var req Update{Feature}Request
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	// 3. 调用服务层
	item, err := h.{feature}Service.Update{Feature}(uint(id), req.Name, req.Remark)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	// 4. 返回成功响应
	response.Success(c, "{功能}更新成功", item)
}

// Delete{Feature} 删除{功能}
// @Summary 删除{功能}
// @Description 删除指定的{功能}
// @Tags {功能}管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "{功能}ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未登录"
// @Failure 404 {object} response.Response "{功能}不存在"
// @Router /user/{features}/{id} [delete]
func (h *{Feature}Handler) Delete{Feature}(c *gin.Context) {
	// 1. 获取ID参数
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "ID格式错误")
		return
	}

	// 2. 调用服务层
	if err := h.{feature}Service.Delete{Feature}(uint(id)); err != nil {
		response.Error(c, err.Error())
		return
	}

	// 3. 返回成功响应
	response.Success(c, "{功能}删除成功", nil)
}
```

**Handler 层职责:**

- ✅ HTTP 请求处理
- ✅ 参数绑定和验证
- ✅ 调用 Service 层
- ✅ 响应格式化
- ✅ Swagger 注释
- ❌ 业务逻辑（应该在 Service 层）
- ❌ 数据库操作（应该通过 Service → Repository）

**Swagger 注释规范:**

```go
// @Summary 简短描述（一句话）
// @Description 详细描述
// @Tags 标签分类（用于 Swagger UI 分组）
// @Accept 接受的内容类型（json/xml）
// @Produce 响应的内容类型（json/xml）
// @Security BearerAuth（如果需要认证）
// @Param 参数名 参数位置 参数类型 是否必须 "描述" 默认值
// @Success 状态码 {响应类型} 响应数据类型 "描述"
// @Failure 状态码 {响应类型} 响应数据类型 "描述"
// @Router 路径 [方法]
```

---

## 步骤 4: 路由与依赖注入

### 4.1 注册路由

**前台路由 (`internal/api/router/frontend.go`):**

```go
// 在 SetupFrontend 函数参数中添加 handler
func SetupFrontend(
	userHandler *frontendHandler.UserHandler,
	{feature}Handler *frontendHandler.{Feature}Handler, // 添加这行
	jwtSecret string,
	redisClient *redis.Client,
	cfg *config.Config,
	logger *zap.Logger,
	mode string,
) *gin.Engine {
	// ... 现有代码 ...

	// 用户接口（需要用户认证）
	user := v1.Group("/user")
	user.Use(middleware.Auth(jwtSecret, logger))
	{
		// ... 现有路由 ...

		// {功能}管理 - 添加这部分
		user.POST("/{features}", {feature}Handler.Create{Feature})
		user.GET("/{features}", {feature}Handler.Get{Feature}List)
		user.GET("/{features}/:id", {feature}Handler.Get{Feature})
		user.PUT("/{features}/:id", {feature}Handler.Update{Feature})
		user.DELETE("/{features}/:id", {feature}Handler.Delete{Feature})
	}

	return r
}
```

**后台路由 (`internal/api/router/backend.go`):**

```go
// 在 SetupBackend 函数参数中添加 handler
func SetupBackend(
	adminUserHandler *backendHandler.AdminUserHandler,
	rbacHandler *backendHandler.RBACHandler,
	{feature}Handler *backendHandler.{Feature}Handler, // 添加这行
	rbacService service.RBACService,
	jwtSecret string,
	redisClient *redis.Client,
	cfg *config.Config,
	logger *zap.Logger,
	mode string,
) *gin.Engine {
	// ... 现有代码 ...

	// 管理员路由
	admin := r.Group("/api/v1/admin")
	admin.Use(middleware.AdminAuth(jwtSecret, logger))
	{
		// ... 现有路由 ...

		// {功能}管理 - 添加这部分
		admin.POST("/{features}", {feature}Handler.Create{Feature})
		admin.GET("/{features}", {feature}Handler.Get{Feature}List)
		admin.GET("/{features}/:id", {feature}Handler.Get{Feature})
		admin.PUT("/{features}/:id", {feature}Handler.Update{Feature})
		admin.DELETE("/{features}/:id", {feature}Handler.Delete{Feature})
	}

	return r
}
```

### 4.2 配置 Wire 依赖注入

**步骤 1: 更新 Router Provider (`cmd/frontend/providers.go` 或 `cmd/backend/providers.go`)**

```go
// 前台示例
func provideFrontendRouter(
	userHandler *frontendHandler.UserHandler,
	{feature}Handler *frontendHandler.{Feature}Handler, // 添加这行
	redisClient *redis.Client,
	logger *zap.Logger,
	cfg *config.Config,
) *gin.Engine {
	return router.SetupFrontend(
		userHandler,
		{feature}Handler, // 添加这行
		cfg.JWT.Secret,
		redisClient,
		cfg,
		logger,
		cfg.Server.Mode,
	)
}
```

**步骤 2: 更新 Wire 配置 (`cmd/frontend/wire.go` 或 `cmd/backend/wire.go`)**

```go
//go:build wireinject
// +build wireinject

package main

import (
	frontendHandler "trx-project/internal/api/handler/frontendHandler"
	"trx-project/internal/repository"
	"trx-project/internal/service"
	"trx-project/pkg/config"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func initFrontendApp(cfg *config.Config) (*gin.Engine, func(), error) {
	wire.Build(
		// Logger
		provideLogger,
		// Database
		provideDB,
		// Redis
		provideRedis,
		// JWT Config
		provideJWTConfig,

		// Repository
		repository.NewUserRepository,
		repository.New{Feature}Repository,  // 添加这行

		// Service
		service.NewUserService,
		service.New{Feature}Service,        // 添加这行

		// Handler
		frontendHandler.NewUserHandler,
		frontendHandler.New{Feature}Handler, // 添加这行

		// Frontend Router
		provideFrontendRouter,
	)
	return nil, nil, nil
}
```

**步骤 3: 生成 Wire 代码**

```bash
# 进入对应的 cmd 目录
cd /d/workspace/go/trx-project/cmd/frontend  # 或 cd cmd/backend

# 生成 Wire 代码
wire

# 应该看到输出：
# wire: trx-project/cmd/frontend: wrote D:\workspace\go\trx-project\cmd\frontend\wire_gen.go
```

**验证 Wire 生成:**

```bash
# 查看生成的文件
cat wire_gen.go

# 检查是否有编译错误
cd /d/workspace/go/trx-project
go build ./cmd/frontend
go build ./cmd/backend
```

---

## 步骤 5: 文档与测试

### 5.1 生成 Swagger 文档

**前台文档:**

```bash
cd /d/workspace/go/trx-project

# 生成前台 Swagger 文档
swag init \
  -g cmd/frontend/main.go \
  -o cmd/frontend/docs \
  --parseDependency \
  --parseInternal \
  --instanceName frontend \
  --exclude internal/api/handler/backendHandler

# 清理旧文件
rm -f cmd/frontend/docs/docs.go
rm -f cmd/frontend/docs/swagger.json
rm -f cmd/frontend/docs/swagger.yaml
```

**后台文档:**

```bash
# 生成后台 Swagger 文档
swag init \
  -g cmd/backend/main.go \
  -o cmd/backend/docs \
  --parseDependency \
  --parseInternal \
  --instanceName backend \
  --exclude internal/api/handler/frontendHandler

# 清理旧文件
rm -f cmd/backend/docs/docs.go
rm -f cmd/backend/docs/swagger.json
rm -f cmd/backend/docs/swagger.yaml
```

**或使用 Makefile:**

```bash
# 生成所有文档
make swag

# 只生成前台文档
make swag-frontend

# 只生成后台文档
make swag-backend
```

**访问 Swagger UI:**

- 前台: http://localhost:8080/swagger/index.html
- 后台: http://localhost:8081/swagger/index.html

### 5.2 编写测试脚本

**创建测试脚本 (`scripts/test_{feature}_api.sh`):**

```bash
#!/bin/bash

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# 服务地址
FRONTEND_URL="http://localhost:8080"

echo "================================================"
echo "      {功能}管理 API 测试"
echo "================================================"
echo ""

# 1. 用户登录获取 Token
echo -e "${YELLOW}1. 用户登录获取 Token${NC}"
LOGIN_RESPONSE=$(curl -s -X POST "$FRONTEND_URL/api/v1/public/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }')

TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
  echo -e "${RED}❌ 登录失败，无法获取 Token${NC}"
  exit 1
fi

echo -e "${GREEN}✅ Token 获取成功${NC}"
echo ""

# 2. 创建{功能}
echo -e "${YELLOW}2. 创建{功能}${NC}"
CREATE_RESPONSE=$(curl -s -X POST "$FRONTEND_URL/api/v1/user/{features}" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "测试{功能}",
    "remark": "这是一个测试"
  }')

echo "$CREATE_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$CREATE_RESPONSE"

ID=$(echo "$CREATE_RESPONSE" | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)

if [ -z "$ID" ]; then
  echo -e "${RED}❌ 创建{功能}失败${NC}"
  exit 1
fi

echo -e "${GREEN}✅ {功能}创建成功，ID: $ID${NC}"
echo ""

# 3. 获取{功能}列表
echo -e "${YELLOW}3. 获取{功能}列表${NC}"
LIST_RESPONSE=$(curl -s -X GET "$FRONTEND_URL/api/v1/user/{features}?page=1&page_size=10" \
  -H "Authorization: Bearer $TOKEN")

echo "$LIST_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$LIST_RESPONSE"
echo ""

# 4. 获取{功能}详情
echo -e "${YELLOW}4. 获取{功能}详情（ID: $ID）${NC}"
GET_RESPONSE=$(curl -s -X GET "$FRONTEND_URL/api/v1/user/{features}/$ID" \
  -H "Authorization: Bearer $TOKEN")

echo "$GET_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$GET_RESPONSE"
echo ""

# 5. 更新{功能}
echo -e "${YELLOW}5. 更新{功能}（ID: $ID）${NC}"
UPDATE_RESPONSE=$(curl -s -X PUT "$FRONTEND_URL/api/v1/user/{features}/$ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "更新后的{功能}",
    "remark": "已更新"
  }')

echo "$UPDATE_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$UPDATE_RESPONSE"
echo -e "${GREEN}✅ {功能}更新成功${NC}"
echo ""

# 6. 删除{功能}
echo -e "${YELLOW}6. 删除{功能}（ID: $ID）${NC}"
DELETE_RESPONSE=$(curl -s -X DELETE "$FRONTEND_URL/api/v1/user/{features}/$ID" \
  -H "Authorization: Bearer $TOKEN")

echo "$DELETE_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$DELETE_RESPONSE"
echo -e "${GREEN}✅ {功能}删除成功${NC}"
echo ""

echo "================================================"
echo -e "${GREEN}✅ 所有测试完成！${NC}"
echo "================================================"
```

**设置执行权限:**

```bash
chmod +x scripts/test_{feature}_api.sh
```

### 5.3 运行测试

```bash
# 1. 启动服务
./bin/frontend  # 或 ./bin/backend

# 2. 在另一个终端运行测试
./scripts/test_{feature}_api.sh
```

### 5.4 编译验证

```bash
cd /d/workspace/go/trx-project

# 编译前台服务
go build -o bin/frontend cmd/frontend/*.go

# 编译后台服务
go build -o bin/backend cmd/backend/*.go

# 检查编译是否成功
echo $?  # 应该输出 0
```

---

## 步骤 6: 代码审查与上线

### 6.1 代码审查清单

**功能完整性:**

- [ ] 所有 CRUD 操作都已实现
- [ ] 业务规则都已正确实现
- [ ] 权限控制都已正确配置
- [ ] 错误处理都已完善

**代码质量:**

- [ ] 遵循项目命名规范
- [ ] 代码有适当的注释
- [ ] 没有硬编码的配置
- [ ] 没有重复代码
- [ ] 使用了统一的错误处理

**数据库:**

- [ ] 迁移文件已创建并测试
- [ ] 索引设计合理
- [ ] 字段类型正确
- [ ] 有软删除字段（如需要）

**API 设计:**

- [ ] URL 命名规范（RESTful）
- [ ] HTTP 方法使用正确
- [ ] 响应格式统一
- [ ] Swagger 注释完整

**测试:**

- [ ] 测试脚本已编写
- [ ] 所有接口都已测试通过
- [ ] 边界情况已测试
- [ ] 错误场景已测试

**文档:**

- [ ] Swagger 文档已生成
- [ ] 功能文档已编写
- [ ] README 已更新

### 6.2 性能测试

**使用 Apache Bench 进行简单的压力测试:**

```bash
# 获取列表接口压力测试
ab -n 1000 -c 10 -H "Authorization: Bearer <token>" \
  http://localhost:8080/api/v1/user/{features}

# 创建接口压力测试
ab -n 100 -c 5 -T "application/json" \
  -H "Authorization: Bearer <token>" \
  -p test_data.json \
  http://localhost:8080/api/v1/user/{features}
```

### 6.3 提交代码

```bash
cd /d/workspace/go/trx-project

# 查看修改的文件
git status

# 添加所有修改
git add .

# 提交（使用规范的提交信息）
git commit -m "feat: 实现{功能}管理模块

- 添加{功能}表迁移文件
- 实现 Model、Repository、Service、Handler 层
- 添加前台/后台路由
- 配置 Wire 依赖注入
- 生成 Swagger 文档
- 添加测试脚本
"

# 推送到远程仓库
git push origin feature/{feature}-management
```

**提交信息规范:**

- `feat`: 新功能
- `fix`: 修复 bug
- `docs`: 文档更新
- `style`: 代码格式调整
- `refactor`: 代码重构
- `test`: 测试相关
- `chore`: 构建/工具相关

### 6.4 部署上线

```bash
# 1. 在服务器上拉取代码
git pull origin main

# 2. 执行数据库迁移
./scripts/migrate.sh up
# 或
make migrate-up

# 3. 编译服务
go build -o bin/frontend cmd/frontend/*.go
go build -o bin/backend cmd/backend/*.go

# 4. 重启服务
systemctl restart trx-frontend
systemctl restart trx-backend

# 5. 检查服务状态
systemctl status trx-frontend
systemctl status trx-backend

# 6. 查看日志
tail -f logs/app.log
```

---

## 📚 附录

### A. 文件清单模板

完成一个功能模块后，应该有以下文件：

```
migrations/
├── {number}_create_{table}_table.up.sql      ✅ 迁移文件（向上）
└── {number}_create_{table}_table.down.sql    ✅ 迁移文件（向下）

internal/
├── model/
│   └── {feature}.go                           ✅ 数据模型
├── repository/
│   └── {feature}_repository.go                ✅ 数据访问层
├── service/
│   └── {feature}_service.go                   ✅ 业务逻辑层
└── api/
    ├── handler/
    │   ├── frontendHandler/
    │   │   └── {feature}_handler.go           ✅ 前台处理器（如需要）
    │   └── backendHandler/
    │       └── {feature}_handler.go           ✅ 后台处理器（如需要）
    └── router/
        ├── frontend.go                        ✅ 前台路由（已更新）
        └── backend.go                         ✅ 后台路由（已更新）

cmd/
├── frontend/
│   ├── wire.go                                ✅ Wire 配置（已更新）
│   ├── wire_gen.go                            ✅ Wire 生成（自动）
│   └── providers.go                           ✅ 依赖提供者（已更新）
└── backend/
    ├── wire.go                                ✅ Wire 配置（已更新）
    ├── wire_gen.go                            ✅ Wire 生成（自动）
    └── providers.go                           ✅ 依赖提供者（已更新）

scripts/
└── test_{feature}_api.sh                      ✅ 测试脚本

docs/
└── features/
    └── {FEATURE}_GUIDE.md                     ✅ 功能文档
```

### B. 常用命令速查

**数据库迁移:**

```bash
# 查看版本
go run cmd/migrate/main.go -cmd version

# 向上迁移
./scripts/migrate.sh up
# 或
make migrate-up

# 向下迁移
./scripts/migrate.sh down
# 或
make migrate-down

# 强制版本
VERSION=7 ./scripts/migrate.sh force
# 或
make migrate-force VERSION=7
```

**Wire 代码生成:**

```bash
# 前台
cd cmd/frontend && wire

# 后台
cd cmd/backend && wire
```

**Swagger 文档生成:**

```bash
# 使用 Makefile
make swag            # 生成所有
make swag-frontend   # 只生成前台
make swag-backend    # 只生成后台

# 或手动执行
swag init -g cmd/frontend/main.go -o cmd/frontend/docs \
  --parseDependency --parseInternal --instanceName frontend \
  --exclude internal/api/handler/backendHandler
```

**编译与运行:**

```bash
# 编译
go build -o bin/frontend cmd/frontend/*.go
go build -o bin/backend cmd/backend/*.go

# 运行
./bin/frontend  # 前台服务
./bin/backend   # 后台服务
```

**测试:**

```bash
# 运行测试脚本
./scripts/test_{feature}_api.sh

# 查看日志
tail -f logs/app.log
```

### C. 常见问题

**Q1: Wire 生成失败怎么办？**

```bash
# 检查 wire.go 语法
go build ./cmd/frontend/wire.go

# 查看详细错误
cd cmd/frontend && wire 2>&1 | more

# 常见原因：
# - import 路径错误
# - 类型不匹配
# - 缺少依赖
```

**Q2: Swagger 文档不显示新接口？**

```bash
# 1. 检查 Swagger 注释格式
# 2. 重新生成文档
make swag

# 3. 清理浏览器缓存
# 4. 检查 Handler 中是否导入了 model
import _ "trx-project/internal/model"
```

**Q3: 数据库迁移失败？**

```bash
# 查看当前版本
go run cmd/migrate/main.go -cmd version

# 查看 SQL 语法
cat migrations/{number}_{name}.up.sql

# 手动执行 SQL（调试用）
mysql -u root -p trx_db < migrations/{number}_{name}.up.sql
```

**Q4: 编译错误？**

```bash
# 查看详细错误
go build -v ./cmd/frontend

# 检查 import 路径
go mod tidy

# 更新依赖
go get -u ./...
```

### D. 最佳实践总结

**1. 分层清晰**
- Model: 只定义数据结构
- Repository: 只做数据库操作
- Service: 只处理业务逻辑
- Handler: 只处理 HTTP 请求

**2. 错误处理**
- 使用 `errors.New()` 创建业务错误
- 使用 `errors.Is()` 判断错误类型
- 记录详细的错误日志
- 返回友好的错误信息

**3. 日志记录**
- 关键操作必须记录日志
- 使用结构化日志（zap.String, zap.Uint）
- 包含足够的上下文信息
- 区分 Info, Warn, Error 级别

**4. 性能优化**
- 使用索引优化查询
- 避免 N+1 查询问题
- 合理使用缓存
- 分页查询大数据集

**5. 安全性**
- 所有接口都要有认证
- 验证用户权限
- 参数验证
- SQL 注入防护（使用 ORM）
- XSS 防护

**6. 可维护性**
- 遵循命名规范
- 编写清晰的注释
- 保持函数简短
- 避免重复代码
- 编写测试

---

## 📝 总结

按照本文档的流程，你可以规范、高效地实现任何新功能模块。

**核心步骤回顾:**

1. ✅ **需求分析** - 明确功能、设计数据库、规划 API
2. ✅ **数据库迁移** - 创建迁移文件、执行迁移
3. ✅ **分层实现** - Model → Repository → Service → Handler
4. ✅ **路由配置** - 注册路由、配置 Wire
5. ✅ **文档测试** - 生成文档、编写测试、执行验证
6. ✅ **代码审查** - 检查质量、性能测试、提交上线

**建议:**

- 严格遵循分层架构原则
- 使用本文档提供的代码模板
- 每一步都进行测试验证
- 保持代码风格统一
- 及时更新文档

---

**参考文档:**

- [订单管理功能实现案例](ORDER_MANAGEMENT_GUIDE.md)
- [项目架构说明](../architecture/ARCHITECTURE.md)
- [Wire 使用指南](https://github.com/google/wire)
- [Swagger 注释规范](https://github.com/swaggo/swag)

**文档版本:** v1.0  
**更新时间:** 2025-01-12  
**维护人员:** 开发团队

