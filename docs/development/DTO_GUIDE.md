# DTO (Data Transfer Object) 使用指南

## 📋 概述

DTO（Data Transfer Object）用于在 API 层和业务层之间传输数据，将数据库 Model 与 API 响应解耦。

## 🎯 为什么需要 DTO？

### 问题

**直接返回 Model 的问题:**
- ❌ 暴露数据库内部字段（如 `deleted_at`）
- ❌ 可能泄露敏感信息（如密码）
- ❌ 不同接口需要不同的字段组合
- ❌ 无法添加计算字段或格式化字段
- ❌ 数据库结构变更影响 API 响应

### 解决方案

**使用 DTO 的优势:**
- ✅ 只返回前端需要的字段
- ✅ 隐藏敏感信息
- ✅ 不同接口可以使用不同的 DTO
- ✅ 可以添加计算字段、格式化字段
- ✅ 数据库变更不影响 API 响应
- ✅ 更好的 API 版本控制

## 📁 目录结构

```
internal/
├── dto/                    # DTO 层
│   ├── order_dto.go        # 订单相关 DTO
│   ├── user_dto.go         # 用户相关 DTO
│   └── ...                 # 其他 DTO
├── model/                  # Model 层（数据库模型）
│   ├── order.go
│   └── user.go
└── api/
    └── handler/            # Handler 层（使用 DTO）
        └── ...
```

## 🏗️ DTO 设计原则

### 1. 按功能分组

每个功能模块创建对应的 DTO 文件：

```go
// internal/dto/order_dto.go
package dto

type OrderDTO struct { ... }
type OrderListDTO struct { ... }
type OrderSummaryDTO struct { ... }
```

### 2. 不同场景使用不同 DTO

| 场景 | DTO 类型 | 说明 |
|------|---------|------|
| 详情接口 | `OrderDTO` | 包含完整信息 |
| 列表接口 | `OrderListDTO` | 只包含必要字段 |
| 统计接口 | `OrderSummaryDTO` | 包含统计信息 |

### 3. 命名规范

- **详情 DTO**: `{Feature}DTO` (如 `OrderDTO`)
- **列表 DTO**: `{Feature}ListDTO` (如 `OrderListDTO`)
- **统计 DTO**: `{Feature}SummaryDTO` (如 `OrderSummaryDTO`)
- **转换函数**: `To{Feature}DTO()` (如 `ToOrderDTO()`)

## 📝 实现示例

### 1. 创建 DTO 结构

```go
// internal/dto/order_dto.go
package dto

import (
	"time"
	"trx-project/internal/model"
)

// OrderDTO 订单详情 DTO
type OrderDTO struct {
	ID           uint       `json:"id"`
	OrderNo      string     `json:"order_no"`
	ProductName  string     `json:"product_name"`
	TotalAmount  float64    `json:"total_amount"`
	Status       int        `json:"status"`
	StatusText   string     `json:"status_text"`
	CreatedAt    time.Time  `json:"created_at"`
	// 注意：不包含 deleted_at 等数据库内部字段
}

// OrderListDTO 订单列表 DTO（简化版）
type OrderListDTO struct {
	ID          uint      `json:"id"`
	OrderNo     string    `json:"order_no"`
	ProductName string    `json:"product_name"`
	TotalAmount float64   `json:"total_amount"`
	Status      int       `json:"status"`
	StatusText  string    `json:"status_text"`
	CreatedAt   time.Time `json:"created_at"`
}
```

### 2. 创建转换函数

```go
// ToOrderDTO 将 Model 转换为 DTO
func ToOrderDTO(order *model.Order) *OrderDTO {
	if order == nil {
		return nil
	}

	return &OrderDTO{
		ID:           order.ID,
		OrderNo:      order.OrderNo,
		ProductName:  order.ProductName,
		TotalAmount:  order.TotalAmount,
		Status:       int(order.Status),
		StatusText:   order.StatusText,
		CreatedAt:    order.CreatedAt,
	}
}

// ToOrderListDTO 将 Model 转换为列表 DTO
func ToOrderListDTO(order *model.Order) *OrderListDTO {
	if order == nil {
		return nil
	}

	return &OrderListDTO{
		ID:          order.ID,
		OrderNo:     order.OrderNo,
		ProductName: order.ProductName,
		TotalAmount: order.TotalAmount,
		Status:      int(order.Status),
		StatusText:  order.StatusText,
		CreatedAt:   order.CreatedAt,
	}
}

// ToOrderDTOList 批量转换
func ToOrderDTOList(orders []model.Order) []*OrderDTO {
	if len(orders) == 0 {
		return []*OrderDTO{}
	}

	result := make([]*OrderDTO, 0, len(orders))
	for i := range orders {
		result = append(result, ToOrderDTO(&orders[i]))
	}
	return result
}
```

### 3. 在 Handler 中使用

```go
// internal/api/handler/frontendHandler/order_handler.go
package frontendHandler

import (
	"trx-project/internal/dto"
	"trx-project/internal/service"
	"trx-project/pkg/response"
)

func (h *OrderHandler) GetOrder(c *gin.Context) {
	// ... 获取订单逻辑 ...

	// 转换为 DTO 返回
	orderDTO := dto.ToOrderDTO(order)
	response.Success(c, "获取成功", orderDTO)
}

func (h *OrderHandler) GetOrders(c *gin.Context) {
	// ... 获取订单列表逻辑 ...

	// 转换为列表 DTO
	orderListDTOs := dto.ToOrderListDTOList(orders)
	response.Pagination(c, "获取成功", orderListDTOs, total, page, pageSize)
}
```

## 🎨 高级用法

### 1. 添加计算字段

```go
// OrderDTO 可以包含计算字段
type OrderDTO struct {
	ID           uint       `json:"id"`
	TotalAmount  float64    `json:"total_amount"`
	Discount     float64    `json:"discount"`      // 折扣金额（计算字段）
	FinalAmount  float64    `json:"final_amount"`  // 最终金额（计算字段）
	CreatedAt    time.Time  `json:"created_at"`
	CreatedDays  int        `json:"created_days"`  // 创建天数（计算字段）
}

func ToOrderDTO(order *model.Order) *OrderDTO {
	dto := &OrderDTO{
		ID:          order.ID,
		TotalAmount: order.TotalAmount,
		CreatedAt:   order.CreatedAt,
	}
	
	// 计算字段
	dto.Discount = calculateDiscount(order)
	dto.FinalAmount = dto.TotalAmount - dto.Discount
	dto.CreatedDays = int(time.Since(order.CreatedAt).Hours() / 24)
	
	return dto
}
```

### 2. 格式化字段

```go
// OrderDTO 可以包含格式化后的字段
type OrderDTO struct {
	ID           uint   `json:"id"`
	TotalAmount  float64 `json:"total_amount"`
	AmountText   string  `json:"amount_text"`  // 格式化后的金额文本
	CreatedAt    time.Time `json:"created_at"`
	CreatedText  string    `json:"created_text"` // 格式化后的时间文本
}

func ToOrderDTO(order *model.Order) *OrderDTO {
	dto := &OrderDTO{
		ID:          order.ID,
		TotalAmount: order.TotalAmount,
		CreatedAt:   order.CreatedAt,
	}
	
	// 格式化字段
	dto.AmountText = fmt.Sprintf("¥%.2f", order.TotalAmount)
	dto.CreatedText = order.CreatedAt.Format("2006-01-02 15:04:05")
	
	return dto
}
```

### 3. 嵌套 DTO

```go
// OrderDTO 可以嵌套其他 DTO
type OrderDTO struct {
	ID           uint      `json:"id"`
	OrderNo      string    `json:"order_no"`
	User         *UserDTO  `json:"user"`        // 嵌套用户 DTO
	Items        []*OrderItemDTO `json:"items"` // 嵌套订单项 DTO
	CreatedAt    time.Time `json:"created_at"`
}

func ToOrderDTO(order *model.Order) *OrderDTO {
	dto := &OrderDTO{
		ID:        order.ID,
		OrderNo:   order.OrderNo,
		CreatedAt: order.CreatedAt,
	}
	
	// 嵌套转换
	if order.User != nil {
		dto.User = ToUserDTO(order.User)
	}
	
	return dto
}
```

### 4. 条件字段

```go
// OrderDTO 可以根据条件包含不同字段
type OrderDTO struct {
	ID           uint       `json:"id"`
	OrderNo      string     `json:"order_no"`
	TotalAmount  float64    `json:"total_amount"`
	// 敏感字段，只有管理员可以看到
	CostPrice    *float64   `json:"cost_price,omitempty"`  // 成本价
	Profit       *float64   `json:"profit,omitempty"`      // 利润
}

func ToOrderDTO(order *model.Order, isAdmin bool) *OrderDTO {
	dto := &OrderDTO{
		ID:          order.ID,
		OrderNo:     order.OrderNo,
		TotalAmount: order.TotalAmount,
	}
	
	// 只有管理员可以看到成本信息
	if isAdmin {
		costPrice := calculateCostPrice(order)
		profit := order.TotalAmount - costPrice
		dto.CostPrice = &costPrice
		dto.Profit = &profit
	}
	
	return dto
}
```

## 📋 完整示例

### 订单 DTO 完整实现

```go
// internal/dto/order_dto.go
package dto

import (
	"fmt"
	"time"
	"trx-project/internal/model"
)

// OrderDTO 订单详情 DTO
type OrderDTO struct {
	ID           uint       `json:"id"`
	OrderNo      string     `json:"order_no"`
	UserID       uint       `json:"user_id"`
	ProductName  string     `json:"product_name"`
	ProductPrice float64    `json:"product_price"`
	Quantity     int        `json:"quantity"`
	TotalAmount  float64    `json:"total_amount"`
	Status       int        `json:"status"`
	StatusText   string     `json:"status_text"`
	Remark       string     `json:"remark"`
	PaidAt       *time.Time `json:"paid_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	
	// 计算字段
	AmountText   string     `json:"amount_text"`   // 格式化金额
	CreatedText  string     `json:"created_text"`  // 格式化时间
	CanCancel    bool       `json:"can_cancel"`    // 是否可以取消
}

// OrderListDTO 订单列表 DTO
type OrderListDTO struct {
	ID          uint      `json:"id"`
	OrderNo     string    `json:"order_no"`
	ProductName string    `json:"product_name"`
	TotalAmount float64   `json:"total_amount"`
	Status      int       `json:"status"`
	StatusText  string    `json:"status_text"`
	CreatedAt   time.Time `json:"created_at"`
	AmountText  string    `json:"amount_text"`
}

// ToOrderDTO 转换函数
func ToOrderDTO(order *model.Order) *OrderDTO {
	if order == nil {
		return nil
	}

	dto := &OrderDTO{
		ID:           order.ID,
		OrderNo:      order.OrderNo,
		UserID:       order.UserID,
		ProductName:  order.ProductName,
		ProductPrice: order.ProductPrice,
		Quantity:     order.Quantity,
		TotalAmount:  order.TotalAmount,
		Status:       int(order.Status),
		StatusText:   order.StatusText,
		Remark:       order.Remark,
		PaidAt:       order.PaidAt,
		CreatedAt:    order.CreatedAt,
		UpdatedAt:    order.UpdatedAt,
	}

	// 计算字段
	dto.AmountText = fmt.Sprintf("¥%.2f", order.TotalAmount)
	dto.CreatedText = order.CreatedAt.Format("2006-01-02 15:04:05")
	dto.CanCancel = order.Status == model.OrderStatusPending

	return dto
}

// ToOrderListDTO 列表转换函数
func ToOrderListDTO(order *model.Order) *OrderListDTO {
	if order == nil {
		return nil
	}

	return &OrderListDTO{
		ID:          order.ID,
		OrderNo:     order.OrderNo,
		ProductName: order.ProductName,
		TotalAmount: order.TotalAmount,
		Status:      int(order.Status),
		StatusText:  order.StatusText,
		CreatedAt:   order.CreatedAt,
		AmountText:  fmt.Sprintf("¥%.2f", order.TotalAmount),
	}
}

// ToOrderDTOList 批量转换
func ToOrderDTOList(orders []model.Order) []*OrderDTO {
	if len(orders) == 0 {
		return []*OrderDTO{}
	}

	result := make([]*OrderDTO, 0, len(orders))
	for i := range orders {
		result = append(result, ToOrderDTO(&orders[i]))
	}
	return result
}

// ToOrderListDTOList 批量转换列表
func ToOrderListDTOList(orders []model.Order) []*OrderListDTO {
	if len(orders) == 0 {
		return []*OrderListDTO{}
	}

	result := make([]*OrderListDTO, 0, len(orders))
	for i := range orders {
		result = append(result, ToOrderListDTO(&orders[i]))
	}
	return result
}
```

## 🔄 迁移现有代码

### 步骤 1: 创建 DTO

```bash
# 为订单创建 DTO
vim internal/dto/order_dto.go
```

### 步骤 2: 更新 Handler

```go
// 修改前
func (h *OrderHandler) GetOrder(c *gin.Context) {
	order, err := h.orderService.GetOrderByID(id, userID)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, "获取成功", order)  // 直接返回 Model
}

// 修改后
func (h *OrderHandler) GetOrder(c *gin.Context) {
	order, err := h.orderService.GetOrderByID(id, userID)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	orderDTO := dto.ToOrderDTO(order)  // 转换为 DTO
	response.Success(c, "获取成功", orderDTO)
}
```

### 步骤 3: 更新 Swagger 注释

```go
// @Success 200 {object} response.Response{data=dto.OrderDTO} "获取成功"
// 而不是
// @Success 200 {object} response.Response{data=model.Order} "获取成功"
```

## 💡 最佳实践

### 1. 分层清晰

```
Model (数据库层) → DTO (传输层) → Handler (API 层)
```

### 2. 不同场景使用不同 DTO

- **详情接口**: 使用完整的 DTO
- **列表接口**: 使用简化的 ListDTO
- **统计接口**: 使用 SummaryDTO

### 3. 隐藏敏感信息

```go
// ❌ 错误：直接返回 Model，可能泄露密码
response.Success(c, user)

// ✅ 正确：使用 DTO，不包含密码字段
userDTO := dto.ToUserDTO(user)
response.Success(c, userDTO)
```

### 4. 添加计算字段

在 DTO 转换函数中添加业务计算：

```go
func ToOrderDTO(order *model.Order) *OrderDTO {
	dto := &OrderDTO{...}
	
	// 添加计算字段
	dto.AmountText = formatAmount(order.TotalAmount)
	dto.CanCancel = order.Status == model.OrderStatusPending
	
	return dto
}
```

### 5. 性能优化

对于列表接口，使用简化的 DTO：

```go
// 详情接口：完整 DTO
orderDTO := dto.ToOrderDTO(order)

// 列表接口：简化 DTO（减少数据传输）
orderListDTOs := dto.ToOrderListDTOList(orders)
```

## 📚 相关文档

- [功能开发规范](FEATURE_DEVELOPMENT_GUIDE.md)
- [API 响应格式](../deployment/RESPONSE_FORMAT.md)
- [项目架构](../architecture/ARCHITECTURE.md)

## 🎯 总结

使用 DTO 的好处：

1. ✅ **安全性** - 隐藏敏感信息
2. ✅ **灵活性** - 不同接口返回不同字段
3. ✅ **可维护性** - 数据库变更不影响 API
4. ✅ **性能** - 列表接口可以只返回必要字段
5. ✅ **扩展性** - 可以添加计算字段和格式化字段

**推荐流程:**
```
Service 返回 Model → Handler 转换为 DTO → 返回给前端
```

---

**文档版本:** v1.0  
**更新时间:** 2025-01-12

