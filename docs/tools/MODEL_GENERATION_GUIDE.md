# GORM Model 自动生成指南

## 📋 概述

本指南介绍如何使用工具从数据库表结构自动生成对应的 GORM Model 代码，避免手动编写 Model 的繁琐工作。

## 🎯 为什么需要自动生成？

**传统方式的问题:**
- ❌ 手动编写 Model 容易出错
- ❌ 字段类型映射容易遗漏
- ❌ 数据库结构变更后需要手动同步
- ❌ 工作量大，效率低

**自动生成的优势:**
- ✅ 自动从数据库读取表结构
- ✅ 自动生成正确的 GORM 标签
- ✅ 自动处理字段类型映射
- ✅ 快速同步数据库变更
- ✅ 减少手动错误

## 🛠️ 工具选择

### 推荐工具: gorm.io/gen (gentool)

**GORM 官方提供的代码生成工具**，功能强大，维护活跃。

**特点:**
- ✅ GORM 官方维护
- ✅ 支持 MySQL、PostgreSQL、SQLite、SQL Server
- ✅ 自动生成 CRUD 方法
- ✅ 支持自定义查询
- ✅ 类型安全

## 📦 安装工具

```bash
# 安装 gentool
go install gorm.io/gen/tools/gentool@latest

# 验证安装
gentool --version
```

## 🚀 快速开始

### 方式 1: 使用脚本（推荐）

**生成所有表的 Model:**

```bash
# 设置数据库连接信息
export DSN="root:password@tcp(localhost:3306)/trx_db?charset=utf8mb4&parseTime=True&loc=Local"

# 执行生成
./scripts/generate_model_simple.sh
```

**生成指定表的 Model:**

```bash
# 只生成 orders 和 users 表
DSN="root:password@tcp(localhost:3306)/trx_db?charset=utf8mb4&parseTime=True&loc=Local" \
TABLES="orders,users" \
./scripts/generate_model_simple.sh
```

**使用 Makefile:**

```bash
# 生成所有表
make model-gen DSN="root:password@tcp(localhost:3306)/trx_db?charset=utf8mb4&parseTime=True&loc=Local"

# 生成指定表
make model-gen-table DSN="root:password@tcp(localhost:3306)/trx_db?charset=utf8mb4&parseTime=True&loc=Local" TABLES="orders,users"
```

### 方式 2: 直接使用 gentool

```bash
# 生成所有表
gentool -db mysql \
  -dsn "root:password@tcp(localhost:3306)/trx_db?charset=utf8mb4&parseTime=True&loc=Local" \
  -outPath "./internal/model/generated"

# 生成指定表
gentool -db mysql \
  -dsn "root:password@tcp(localhost:3306)/trx_db?charset=utf8mb4&parseTime=True&loc=Local" \
  -tables "orders,users" \
  -outPath "./internal/model/generated"
```

## 📁 生成的文件结构

生成后的目录结构：

```
internal/model/generated/
├── gen.go                    # 生成器配置
├── gen_model.go              # Model 定义
├── query/
│   ├── gen.go                # 查询方法
│   └── {table}.gen.go        # 每个表的查询方法
└── model/
    └── {table}.gen.go        # 每个表的 Model 定义
```

## 🔄 完整工作流程

### 步骤 1: 创建数据库迁移

```bash
# 创建迁移文件
NAME=create_orders_table ./scripts/migrate.sh create

# 编辑迁移文件
vim migrations/000007_create_orders_table.up.sql
```

### 步骤 2: 执行迁移

```bash
# 执行迁移，创建表结构
./scripts/migrate.sh up
```

### 步骤 3: 生成 Model

```bash
# 从数据库生成 Model
DSN="root:password@tcp(localhost:3306)/trx_db?charset=utf8mb4&parseTime=True&loc=Local" \
TABLES="orders" \
./scripts/generate_model_simple.sh
```

### 步骤 4: 复制并自定义 Model

```bash
# 查看生成的文件
ls -la internal/model/generated/model/

# 复制到 model 目录
cp internal/model/generated/model/orders.gen.go internal/model/order.go

# 编辑并添加自定义逻辑
vim internal/model/order.go
```

**添加自定义内容:**

```go
package model

import (
	"time"
	"gorm.io/gorm"
)

// OrderStatus 订单状态枚举（手动添加）
type OrderStatus int

const (
	OrderStatusPending   OrderStatus = 0
	OrderStatusPaid      OrderStatus = 1
	OrderStatusCancelled OrderStatus = 2
	OrderStatusCompleted OrderStatus = 3
)

// Order 订单模型（从生成的文件复制，并添加自定义内容）
type Order struct {
	ID           uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderNo      string         `gorm:"type:varchar(32);uniqueIndex:uk_order_no;not null" json:"order_no"`
	UserID       uint           `gorm:"index:idx_user_id;not null" json:"user_id"`
	// ... 其他字段
	
	// 自定义字段（不存储到数据库）
	StatusText   string         `gorm:"-" json:"status_text"`
}

// AfterFind GORM 钩子（手动添加）
func (o *Order) AfterFind(tx *gorm.DB) error {
	o.StatusText = OrderStatusText[o.Status]
	return nil
}
```

## 📝 生成的文件示例

**生成的 Model (orders.gen.go):**

```go
package model

import (
	"time"
	"gorm.io/gorm"
)

// Order 订单模型
type Order struct {
	ID           uint           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	OrderNo      string         `gorm:"column:order_no;type:varchar(32);not null;uniqueIndex:uk_order_no" json:"order_no"`
	UserID       uint           `gorm:"column:user_id;type:bigint unsigned;not null;index:idx_user_id" json:"user_id"`
	ProductName  string         `gorm:"column:product_name;type:varchar(255);not null" json:"product_name"`
	ProductPrice float64        `gorm:"column:product_price;type:decimal(10,2);not null" json:"product_price"`
	Quantity     int            `gorm:"column:quantity;type:int;not null;default:1" json:"quantity"`
	TotalAmount  float64        `gorm:"column:total_amount;type:decimal(10,2);not null" json:"total_amount"`
	Status       int            `gorm:"column:status;type:tinyint;not null;default:0;index:idx_status" json:"status"`
	Remark       string         `gorm:"column:remark;type:text" json:"remark"`
	PaidAt       *time.Time     `gorm:"column:paid_at;type:datetime" json:"paid_at"`
	CreatedAt    time.Time      `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;index:idx_created_at" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;type:datetime;index" json:"-"`
}

// TableName 指定表名
func (Order) TableName() string {
	return "orders"
}
```

## ⚙️ 配置选项

### gentool 参数说明

| 参数 | 说明 | 示例 |
|------|------|------|
| `-db` | 数据库类型 | `mysql`, `postgres`, `sqlite`, `sqlserver` |
| `-dsn` | 数据库连接字符串 | `user:pass@tcp(localhost:3306)/db?charset=utf8mb4` |
| `-tables` | 指定要生成的表（逗号分隔） | `orders,users` |
| `-outPath` | 输出目录 | `./internal/model/generated` |
| `-onlyModel` | 只生成 Model，不生成查询方法 | `-onlyModel` |
| `-fieldWithIndex` | 为索引字段添加标签 | `-fieldWithIndex` |
| `-fieldWithType` | 为字段添加类型标签 | `-fieldWithType` |

### 完整命令示例

```bash
gentool -db mysql \
  -dsn "root:password@tcp(localhost:3306)/trx_db?charset=utf8mb4&parseTime=True&loc=Local" \
  -tables "orders,users" \
  -outPath "./internal/model/generated" \
  -onlyModel \
  -fieldWithIndex \
  -fieldWithType
```

## 💡 最佳实践

### 1. 工作流程

```
创建迁移文件 → 执行迁移 → 生成 Model → 复制并自定义
```

### 2. 生成后的处理

**必须做的:**
- ✅ 复制生成的文件到 `internal/model/` 目录
- ✅ 添加业务相关的枚举类型
- ✅ 添加 GORM 钩子方法（如需要）
- ✅ 添加自定义字段（如 StatusText）
- ✅ 检查并调整字段标签

**可选做的:**
- 🔸 添加验证方法
- 🔸 添加业务方法
- 🔸 添加辅助函数

### 3. 不要直接使用生成的文件

**原因:**
- 生成的文件会在下次生成时被覆盖
- 需要添加业务逻辑
- 需要自定义字段和方法

**正确做法:**
```bash
# 1. 生成到临时目录
gentool ... -outPath "./internal/model/generated"

# 2. 复制到正式目录
cp internal/model/generated/model/orders.gen.go internal/model/order.go

# 3. 编辑并添加自定义内容
vim internal/model/order.go
```

### 4. 版本控制

**建议:**
- ✅ 将 `internal/model/` 目录加入版本控制
- ❌ 不要将 `internal/model/generated/` 加入版本控制（添加到 .gitignore）

**.gitignore 配置:**

```
# 生成的 Model 文件
internal/model/generated/
```

## 🔧 常见问题

### Q1: 生成失败，提示连接数据库失败？

**解决方法:**
```bash
# 1. 检查数据库是否运行
mysql -u root -p -e "SELECT 1"

# 2. 检查 DSN 格式是否正确
# 格式: user:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local

# 3. 测试连接
mysql -u root -p -h localhost -P 3306 trx_db
```

### Q2: 生成的字段类型不对？

**原因:** 数据库字段类型映射问题

**解决方法:**
- 检查数据库字段类型
- 手动调整生成的 Model 中的类型
- 或使用 `-fieldWithType` 参数

### Q3: 如何只生成 Model，不生成查询方法？

```bash
gentool ... -onlyModel
```

### Q4: 生成的代码格式不符合项目规范？

**解决方法:**
```bash
# 使用 gofmt 格式化
gofmt -w internal/model/generated/

# 或使用 goimports
goimports -w internal/model/generated/
```

### Q5: 如何从配置文件读取 DSN？

**使用脚本:**

```bash
# 方式 1: 使用简化脚本（需要手动设置 DSN）
DSN="..." ./scripts/generate_model_simple.sh

# 方式 2: 从配置文件读取（需要实现读取配置的脚本）
./scripts/generate_model.sh -c config/config.yaml
```

## 📚 相关文档

- [GORM Gen 官方文档](https://gorm.io/zh_CN/gen/gen_tool.html)
- [功能开发规范](../development/FEATURE_DEVELOPMENT_GUIDE.md)
- [数据库迁移指南](../features/MIGRATION_GUIDE.md)

## 🎯 总结

使用自动生成工具可以：

1. ✅ **提高效率** - 快速生成 Model 代码
2. ✅ **减少错误** - 自动处理类型映射
3. ✅ **保持同步** - 数据库变更后快速更新 Model
4. ✅ **标准化** - 统一的代码风格

**推荐流程:**
```
迁移文件 → 执行迁移 → 生成 Model → 复制并自定义 → 使用
```

---

**文档版本:** v1.0  
**更新时间:** 2025-01-12

