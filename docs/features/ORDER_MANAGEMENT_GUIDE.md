# 用户订单管理功能实现指南

## 📋 功能概述

本文档记录了前台用户订单管理功能的完整实现过程，包括创建订单、查询订单、支付订单、取消订单和完成订单等功能。

## 🎯 功能列表

| 功能 | 接口路径 | 方法 | 说明 |
|------|---------|------|------|
| 创建订单 | `/api/v1/user/orders` | POST | 用户创建新订单 |
| 获取订单列表 | `/api/v1/user/orders` | GET | 获取当前用户的订单列表（分页） |
| 获取订单详情 | `/api/v1/user/orders/:id` | GET | 根据订单 ID 获取详情 |
| 支付订单 | `/api/v1/user/orders/:id/pay` | POST | 用户支付订单 |
| 取消订单 | `/api/v1/user/orders/:id/cancel` | POST | 用户取消订单（仅待支付状态） |
| 完成订单 | `/api/v1/user/orders/:id/complete` | POST | 用户确认收货，完成订单 |

## 🗄️ 数据库设计

### 订单表 (orders)

```sql
CREATE TABLE `orders` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '订单ID',
  `order_no` varchar(32) NOT NULL COMMENT '订单号',
  `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
  `product_name` varchar(255) NOT NULL COMMENT '商品名称',
  `product_price` decimal(10,2) NOT NULL COMMENT '商品单价',
  `quantity` int NOT NULL DEFAULT 1 COMMENT '购买数量',
  `total_amount` decimal(10,2) NOT NULL COMMENT '订单总金额',
  `status` tinyint NOT NULL DEFAULT 0 COMMENT '订单状态: 0-待支付, 1-已支付, 2-已取消, 3-已完成',
  `remark` text COMMENT '订单备注',
  `paid_at` datetime DEFAULT NULL COMMENT '支付时间',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_order_no` (`order_no`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='订单表';
```

### 订单状态说明

| 状态值 | 状态名称 | 说明 |
|-------|---------|------|
| 0 | 待支付 | 订单已创建，等待支付 |
| 1 | 已支付 | 订单已支付，等待发货/收货 |
| 2 | 已取消 | 用户取消订单 |
| 3 | 已完成 | 订单已完成 |

## 📁 代码结构

```
trx-project/
├── migrations/
│   ├── 000007_create_orders_table.up.sql     # 订单表创建迁移
│   └── 000007_create_orders_table.down.sql   # 订单表删除迁移
├── internal/
│   ├── model/
│   │   └── order.go                           # 订单模型
│   ├── repository/
│   │   └── order_repository.go                # 订单数据访问层
│   ├── service/
│   │   └── order_service.go                   # 订单业务逻辑层
│   └── api/
│       ├── handler/
│       │   └── frontendHandler/
│       │       └── order_handler.go           # 订单 HTTP 处理器
│       └── router/
│           └── frontend.go                    # 前台路由（已更新）
├── cmd/
│   └── frontend/
│       ├── wire.go                            # Wire 依赖注入配置
│       ├── wire_gen.go                        # Wire 生成的代码
│       └── providers.go                       # 依赖提供者
└── scripts/
    └── test_order_api.sh                      # 订单 API 测试脚本
```

## 🚀 实现步骤

### 步骤 1: 创建数据库迁移

```bash
# 迁移文件已创建: migrations/000007_create_orders_table.up.sql
# 执行迁移
go run cmd/migrate/main.go -cmd up
```

**输出示例:**
```
✅ 加载环境配置: config\config.dev.yaml (环境: dev)
info    migrate/main.go:59      Running database migrations...
info    migrate/main.go:59      Database migrations completed successfully      {"version": 7, "dirty": false}
✅ 数据库迁移完成
```

### 步骤 2: 定义 Model 层

文件: `internal/model/order.go`

**核心代码:**
```go
type Order struct {
    ID           uint            `json:"id"`
    OrderNo      string          `json:"order_no"`
    UserID       uint            `json:"user_id"`
    ProductName  string          `json:"product_name"`
    ProductPrice float64         `json:"product_price"`
    Quantity     int             `json:"quantity"`
    TotalAmount  float64         `json:"total_amount"`
    Status       OrderStatus     `json:"status"`
    StatusText   string          `json:"status_text"`
    Remark       string          `json:"remark"`
    PaidAt       *time.Time      `json:"paid_at"`
    CreatedAt    time.Time       `json:"created_at"`
    UpdatedAt    time.Time       `json:"updated_at"`
}
```

### 步骤 3: 实现 Repository 层

文件: `internal/repository/order_repository.go`

**主要方法:**
- `Create(order *model.Order)` - 创建订单
- `GetByID(id uint)` - 根据 ID 获取订单
- `GetByUserID(userID, page, pageSize)` - 获取用户订单列表（分页）
- `Update(order *model.Order)` - 更新订单
- `UpdateStatus(id, status)` - 更新订单状态
- `Delete(id uint)` - 删除订单（软删除）

### 步骤 4: 实现 Service 层

文件: `internal/service/order_service.go`

**主要方法:**
- `CreateOrder()` - 创建订单（自动计算总金额、生成订单号）
- `GetOrderByID()` - 获取订单详情（验证所有权）
- `GetUserOrders()` - 获取用户订单列表（分页）
- `PayOrder()` - 支付订单（状态验证）
- `CancelOrder()` - 取消订单（状态验证）
- `CompleteOrder()` - 完成订单（状态验证）

### 步骤 5: 实现 Handler 层

文件: `internal/api/handler/frontendHandler/order_handler.go`

**Swagger 注释示例:**
```go
// @Summary 创建订单
// @Description 用户创建新订单
// @Tags 订单管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body CreateOrderRequest true "订单信息"
// @Success 200 {object} response.Response{data=model.Order} "创建成功"
// @Router /user/orders [post]
func (h *OrderHandler) CreateOrder(c *gin.Context) {
    // ... 实现代码
}
```

### 步骤 6: 注册路由

文件: `internal/api/router/frontend.go`

**路由配置:**
```go
user := v1.Group("/user")
user.Use(middleware.Auth(jwtSecret, logger))
{
    // 订单管理
    user.POST("/orders", orderHandler.CreateOrder)
    user.GET("/orders", orderHandler.GetOrders)
    user.GET("/orders/:id", orderHandler.GetOrder)
    user.POST("/orders/:id/pay", orderHandler.PayOrder)
    user.POST("/orders/:id/cancel", orderHandler.CancelOrder)
    user.POST("/orders/:id/complete", orderHandler.CompleteOrder)
}
```

### 步骤 7: 配置 Wire 依赖注入

文件: `cmd/frontend/wire.go`

**添加依赖:**
```go
wire.Build(
    // ... 现有依赖
    repository.NewOrderRepository,   // 新增
    service.NewOrderService,         // 新增
    frontendHandler.NewOrderHandler, // 新增
    // ...
)
```

**生成 Wire 代码:**
```bash
cd cmd/frontend
wire
```

**输出:**
```
wire: trx-project/cmd/frontend: wrote D:\workspace\go\trx-project\cmd\frontend\wire_gen.go
```

### 步骤 8: 编译服务

```bash
cd /d/workspace/go/trx-project
go build -o bin/frontend cmd/frontend/*.go
```

**输出:**
```
✅ 前台服务编译成功
```

### 步骤 9: 生成 Swagger 文档

```bash
swag init \
  -g cmd/frontend/main.go \
  -o cmd/frontend/docs \
  --parseDependency \
  --parseInternal \
  --instanceName frontend \
  --exclude internal/api/handler/backendHandler
```

**清理旧文件:**
```bash
rm -f cmd/frontend/docs/docs.go cmd/frontend/docs/swagger.json cmd/frontend/docs/swagger.yaml
```

## 🧪 测试

### 方式 1: 使用测试脚本（推荐）

```bash
# 1. 启动前台服务
./bin/frontend

# 2. 在另一个终端运行测试脚本
./scripts/test_order_api.sh
```

**测试脚本会执行以下操作:**
1. 用户注册/登录
2. 创建订单
3. 获取订单列表
4. 获取订单详情
5. 支付订单
6. 完成订单
7. 创建第二个订单
8. 取消订单
9. 获取最终订单列表

### 方式 2: 使用 Swagger UI

1. 访问 Swagger 文档: http://localhost:8080/swagger/index.html
2. 点击右上角 "Authorize" 按钮
3. 输入 Token: `Bearer <your-token>`
4. 测试订单管理接口

### 方式 3: 使用 curl 命令

**1. 注册/登录获取 Token:**
```bash
# 注册
curl -X POST "http://localhost:8080/api/v1/public/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'

# 或登录
curl -X POST "http://localhost:8080/api/v1/public/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

**2. 创建订单:**
```bash
TOKEN="<your-token>"

curl -X POST "http://localhost:8080/api/v1/user/orders" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "product_name": "iPhone 15 Pro",
    "product_price": 7999.00,
    "quantity": 1,
    "remark": "请尽快发货"
  }'
```

**响应示例:**
```json
{
  "code": 200,
  "message": "订单创建成功",
  "data": {
    "id": 1,
    "order_no": "ORD202501121710001234",
    "user_id": 1,
    "product_name": "iPhone 15 Pro",
    "product_price": 7999.00,
    "quantity": 1,
    "total_amount": 7999.00,
    "status": 0,
    "status_text": "待支付",
    "remark": "请尽快发货",
    "paid_at": null,
    "created_at": "2025-01-12T17:10:00Z",
    "updated_at": "2025-01-12T17:10:00Z"
  }
}
```

**3. 获取订单列表:**
```bash
curl -X GET "http://localhost:8080/api/v1/user/orders?page=1&page_size=10" \
  -H "Authorization: Bearer $TOKEN"
```

**4. 支付订单:**
```bash
ORDER_ID=1

curl -X POST "http://localhost:8080/api/v1/user/orders/$ORDER_ID/pay" \
  -H "Authorization: Bearer $TOKEN"
```

**5. 取消订单:**
```bash
curl -X POST "http://localhost:8080/api/v1/user/orders/$ORDER_ID/cancel" \
  -H "Authorization: Bearer $TOKEN"
```

**6. 完成订单:**
```bash
curl -X POST "http://localhost:8080/api/v1/user/orders/$ORDER_ID/complete" \
  -H "Authorization: Bearer $TOKEN"
```

## 📊 订单状态流转

```
┌──────────┐
│ 待支付(0) │ ──┐
└──────────┘   │
     │         │
     │ 支付    │ 取消
     ↓         │
┌──────────┐   │      ┌──────────┐
│ 已支付(1) │   └─────→│ 已取消(2) │
└──────────┘          └──────────┘
     │
     │ 确认收货
     ↓
┌──────────┐
│ 已完成(3) │
└──────────┘
```

**状态转换规则:**
- 待支付 → 已支付: 用户支付订单
- 待支付 → 已取消: 用户取消订单
- 已支付 → 已完成: 用户确认收货

**不允许的操作:**
- ❌ 已支付/已取消/已完成的订单不能再次支付
- ❌ 只有待支付的订单才能取消
- ❌ 只有已支付的订单才能完成

## 🔐 权限控制

所有订单接口都需要用户认证：

1. **JWT 认证**: 所有接口需要在请求头中携带有效的 Token
2. **所有权验证**: 用户只能操作自己的订单，不能访问其他用户的订单

**请求头格式:**
```
Authorization: Bearer <your-jwt-token>
```

## 💡 最佳实践

### 1. 订单号生成

订单号格式: `ORD + 年月日时分秒(14位) + 随机数(4位)`

示例: `ORD202501121710001234`

### 2. 金额计算

总金额 = 商品单价 × 购买数量

在 Service 层自动计算，确保数据一致性。

### 3. 状态验证

在执行订单操作前，必须验证订单状态是否允许该操作。

### 4. 日志记录

关键操作都会记录日志：
- 订单创建
- 订单支付
- 订单取消
- 订单完成

### 5. 错误处理

统一的错误响应格式：
```json
{
  "code": 400,
  "message": "订单状态不正确，无法支付",
  "data": null
}
```

## 🐛 常见问题

### 1. 订单创建失败

**可能原因:**
- 用户未登录（缺少 Token）
- 商品价格或数量参数错误
- 数据库连接失败

**解决方法:**
```bash
# 检查 Token 是否有效
# 检查请求参数是否正确
# 查看日志: tail -f logs/app.log
```

### 2. 无权访问订单

**错误信息:** "无权访问该订单"

**原因:** 尝试访问其他用户的订单

**解决方法:** 确保使用正确的 Token 访问自己的订单

### 3. 订单状态错误

**错误信息:** "订单状态不正确，无法支付"

**原因:** 订单当前状态不允许执行该操作

**解决方法:** 检查订单状态，确保符合状态转换规则

## 📈 性能监控

### Prometheus 指标

订单功能集成了 Prometheus 监控指标：

```
# 订单创建总数
trx_order_created_total

# 订单支付总数
trx_order_paid_total

# 订单取消总数
trx_order_cancelled_total

# 订单完成总数
trx_order_completed_total
```

**查看指标:**
- Prometheus UI: http://localhost:9090
- Grafana: http://localhost:3000

## 🔄 扩展建议

### 1. 支付集成

可以集成第三方支付平台（支付宝、微信支付）：
- 生成支付订单
- 接收支付回调
- 更新订单状态

### 2. 订单商品表

当前设计适合单商品订单，如需支持多商品订单，建议：
- 创建 `order_items` 表
- 一对多关联订单和商品
- 更新订单金额计算逻辑

### 3. 订单评价

支付完成后，用户可以评价订单：
- 创建 `order_reviews` 表
- 关联订单和用户
- 记录评分和评论

### 4. 订单物流

跟踪订单物流信息：
- 创建 `order_logistics` 表
- 记录物流公司和单号
- 提供物流查询接口

## 📚 相关文档

- [项目架构说明](../architecture/ARCHITECTURE.md)
- [API 文档](../deployment/API.md)
- [数据库迁移指南](MIGRATION_GUIDE.md)
- [Swagger 使用指南](../tools/SWAGGER_GUIDE.md)

---

**完成时间**: 2025-01-12  
**维护人员**: 开发团队  
**文档版本**: v1.0

