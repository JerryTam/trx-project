package model

import (
	"time"

	"gorm.io/gorm"
)

// OrderStatus 订单状态
type OrderStatus int

const (
	OrderStatusPending   OrderStatus = 0 // 待支付
	OrderStatusPaid      OrderStatus = 1 // 已支付
	OrderStatusCancelled OrderStatus = 2 // 已取消
	OrderStatusCompleted OrderStatus = 3 // 已完成
)

// OrderStatusText 订单状态文本映射
var OrderStatusText = map[OrderStatus]string{
	OrderStatusPending:   "待支付",
	OrderStatusPaid:      "已支付",
	OrderStatusCancelled: "已取消",
	OrderStatusCompleted: "已完成",
}

// Order 订单模型
type Order struct {
	ID           uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderNo      string         `gorm:"type:varchar(32);uniqueIndex:uk_order_no;not null" json:"order_no"`
	UserID       uint           `gorm:"index:idx_user_id;not null" json:"user_id"`
	ProductName  string         `gorm:"type:varchar(255);not null" json:"product_name"`
	ProductPrice float64        `gorm:"type:decimal(10,2);not null" json:"product_price"`
	Quantity     int            `gorm:"not null;default:1" json:"quantity"`
	TotalAmount  float64        `gorm:"type:decimal(10,2);not null" json:"total_amount"`
	Status       OrderStatus    `gorm:"type:tinyint;index:idx_status;not null;default:0" json:"status"`
	StatusText   string         `gorm:"-" json:"status_text"` // 状态文本（不存储到数据库）
	Remark       string         `gorm:"type:text" json:"remark"`
	PaidAt       *time.Time     `json:"paid_at"`
	CreatedAt    time.Time      `gorm:"index:idx_created_at" json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	User         *User          `gorm:"foreignKey:UserID" json:"user,omitempty"` // 关联用户（可选加载）
}

// TableName 指定表名
func (Order) TableName() string {
	return "orders"
}

// AfterFind GORM 钩子：查询后设置状态文本
func (o *Order) AfterFind(tx *gorm.DB) error {
	o.StatusText = OrderStatusText[o.Status]
	return nil
}

// BeforeCreate GORM 钩子：创建前自动生成订单号
func (o *Order) BeforeCreate(tx *gorm.DB) error {
	if o.OrderNo == "" {
		o.OrderNo = GenerateOrderNo()
	}
	return nil
}

// GenerateOrderNo 生成订单号
// 格式: ORD + 时间戳(14位) + 随机数(4位)
func GenerateOrderNo() string {
	return "ORD" + time.Now().Format("20060102150405") + GenerateRandomString(4)
}

// GenerateRandomString 生成指定长度的随机字符串（数字）
func GenerateRandomString(length int) string {
	const charset = "0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[time.Now().UnixNano()%int64(len(charset))]
		time.Sleep(time.Nanosecond)
	}
	return string(result)
}
