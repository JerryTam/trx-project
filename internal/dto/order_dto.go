package dto

import (
	"time"
	"trx-project/internal/model"
)

// OrderDTO 订单响应 DTO
// 用于 API 响应，不包含敏感信息和数据库内部字段
type OrderDTO struct {
	ID           uint       `json:"id"`            // 订单ID
	OrderNo      string     `json:"order_no"`      // 订单号
	UserID       uint       `json:"user_id"`       // 用户ID
	ProductName  string     `json:"product_name"`  // 商品名称
	ProductPrice float64    `json:"product_price"` // 商品单价
	Quantity     int        `json:"quantity"`      // 购买数量
	TotalAmount  float64    `json:"total_amount"`  // 订单总金额
	Status       int        `json:"status"`        // 订单状态
	StatusText   string     `json:"status_text"`   // 状态文本
	Remark       string     `json:"remark"`        // 订单备注
	PaidAt       *time.Time `json:"paid_at"`       // 支付时间
	CreatedAt    time.Time  `json:"created_at"`    // 创建时间
	UpdatedAt    time.Time  `json:"updated_at"`    // 更新时间
	// 注意：不包含 deleted_at 字段，因为这是数据库内部字段
}

// OrderListDTO 订单列表响应 DTO（简化版）
// 用于列表接口，只包含必要字段
type OrderListDTO struct {
	ID          uint      `json:"id"`
	OrderNo     string    `json:"order_no"`
	ProductName string    `json:"product_name"`
	TotalAmount float64   `json:"total_amount"`
	Status      int       `json:"status"`
	StatusText  string    `json:"status_text"`
	CreatedAt   time.Time `json:"created_at"`
}

// OrderSummaryDTO 订单统计 DTO
// 用于订单统计接口
type OrderSummaryDTO struct {
	TotalOrders     int     `json:"total_orders"`     // 总订单数
	TotalAmount     float64 `json:"total_amount"`     // 总金额
	PendingOrders   int     `json:"pending_orders"`   // 待支付订单数
	PaidOrders      int     `json:"paid_orders"`      // 已支付订单数
	CompletedOrders int     `json:"completed_orders"` // 已完成订单数
}

// ToOrderDTO 将 Model 转换为 DTO
func ToOrderDTO(order *model.Order) *OrderDTO {
	if order == nil {
		return nil
	}

	return &OrderDTO{
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

// ToOrderDTOList 批量转换 Model 列表为 DTO 列表
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

// ToOrderListDTOList 批量转换 Model 列表为列表 DTO
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
