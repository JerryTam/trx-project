package service

import (
	"errors"
	"time"
	"trx-project/internal/model"
	"trx-project/internal/repository"

	// "trx-project/pkg/metrics"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// OrderService 订单服务接口
type OrderService interface {
	CreateOrder(userID uint, productName string, productPrice float64, quantity int, remark string) (*model.Order, error)
	GetOrderByID(id uint, userID uint) (*model.Order, error)
	GetOrderByOrderNo(orderNo string, userID uint) (*model.Order, error)
	GetUserOrders(userID uint, page, pageSize int) ([]model.Order, int64, error)
	PayOrder(id uint, userID uint) error
	CancelOrder(id uint, userID uint) error
	CompleteOrder(id uint, userID uint) error
}

type orderService struct {
	orderRepo repository.OrderRepository
	logger    *zap.Logger
}

// NewOrderService 创建订单服务实例
func NewOrderService(orderRepo repository.OrderRepository, logger *zap.Logger) OrderService {
	return &orderService{
		orderRepo: orderRepo,
		logger:    logger,
	}
}

// CreateOrder 创建订单
func (s *orderService) CreateOrder(userID uint, productName string, productPrice float64, quantity int, remark string) (*model.Order, error) {
	// 参数验证
	if productName == "" {
		return nil, errors.New("商品名称不能为空")
	}
	if productPrice <= 0 {
		return nil, errors.New("商品价格必须大于0")
	}
	if quantity <= 0 {
		return nil, errors.New("购买数量必须大于0")
	}

	// 计算总金额
	totalAmount := productPrice * float64(quantity)

	// 创建订单对象
	order := &model.Order{
		UserID:       userID,
		ProductName:  productName,
		ProductPrice: productPrice,
		Quantity:     quantity,
		TotalAmount:  totalAmount,
		Status:       model.OrderStatusPending,
		Remark:       remark,
	}

	// 保存到数据库
	if err := s.orderRepo.Create(order); err != nil {
		s.logger.Error("创建订单失败",
			zap.Uint("user_id", userID),
			zap.String("product_name", productName),
			zap.Error(err))
		return nil, err
	}

	// 业务指标：订单创建（暂时注释，等待 metrics 更新）
	// metrics.OrderCreated.Inc()

	s.logger.Info("订单创建成功",
		zap.Uint("order_id", order.ID),
		zap.String("order_no", order.OrderNo),
		zap.Uint("user_id", userID))

	return order, nil
}

// GetOrderByID 根据ID获取订单
func (s *orderService) GetOrderByID(id uint, userID uint) (*model.Order, error) {
	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("订单不存在")
		}
		s.logger.Error("获取订单失败", zap.Uint("id", id), zap.Error(err))
		return nil, err
	}

	// 验证订单所有权
	if order.UserID != userID {
		return nil, errors.New("无权访问该订单")
	}

	return order, nil
}

// GetOrderByOrderNo 根据订单号获取订单
func (s *orderService) GetOrderByOrderNo(orderNo string, userID uint) (*model.Order, error) {
	order, err := s.orderRepo.GetByOrderNo(orderNo)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("订单不存在")
		}
		s.logger.Error("获取订单失败", zap.String("order_no", orderNo), zap.Error(err))
		return nil, err
	}

	// 验证订单所有权
	if order.UserID != userID {
		return nil, errors.New("无权访问该订单")
	}

	return order, nil
}

// GetUserOrders 获取用户订单列表
func (s *orderService) GetUserOrders(userID uint, page, pageSize int) ([]model.Order, int64, error) {
	// 参数验证
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	orders, total, err := s.orderRepo.GetByUserID(userID, page, pageSize)
	if err != nil {
		s.logger.Error("获取用户订单列表失败",
			zap.Uint("user_id", userID),
			zap.Error(err))
		return nil, 0, err
	}

	return orders, total, nil
}

// PayOrder 支付订单
func (s *orderService) PayOrder(id uint, userID uint) error {
	// 获取订单
	order, err := s.GetOrderByID(id, userID)
	if err != nil {
		return err
	}

	// 检查订单状态
	if order.Status != model.OrderStatusPending {
		return errors.New("订单状态不正确，无法支付")
	}

	// 更新状态为已支付
	order.Status = model.OrderStatusPaid
	now := time.Now()
	order.PaidAt = &now

	if err := s.orderRepo.Update(order); err != nil {
		s.logger.Error("更新订单状态失败",
			zap.Uint("order_id", id),
			zap.Error(err))
		return err
	}

	// 业务指标：订单支付（暂时注释，等待 metrics 更新）
	// metrics.OrderPaid.Inc()

	s.logger.Info("订单支付成功",
		zap.Uint("order_id", id),
		zap.String("order_no", order.OrderNo))

	return nil
}

// CancelOrder 取消订单
func (s *orderService) CancelOrder(id uint, userID uint) error {
	// 获取订单
	order, err := s.GetOrderByID(id, userID)
	if err != nil {
		return err
	}

	// 检查订单状态（只有待支付的订单可以取消）
	if order.Status != model.OrderStatusPending {
		return errors.New("只有待支付的订单可以取消")
	}

	// 更新状态为已取消
	if err := s.orderRepo.UpdateStatus(id, model.OrderStatusCancelled); err != nil {
		s.logger.Error("取消订单失败",
			zap.Uint("order_id", id),
			zap.Error(err))
		return err
	}

	// 业务指标：订单取消（暂时注释，等待 metrics 更新）
	// metrics.OrderCancelled.Inc()

	s.logger.Info("订单取消成功",
		zap.Uint("order_id", id),
		zap.String("order_no", order.OrderNo))

	return nil
}

// CompleteOrder 完成订单
func (s *orderService) CompleteOrder(id uint, userID uint) error {
	// 获取订单
	order, err := s.GetOrderByID(id, userID)
	if err != nil {
		return err
	}

	// 检查订单状态（只有已支付的订单可以完成）
	if order.Status != model.OrderStatusPaid {
		return errors.New("只有已支付的订单可以完成")
	}

	// 更新状态为已完成
	if err := s.orderRepo.UpdateStatus(id, model.OrderStatusCompleted); err != nil {
		s.logger.Error("完成订单失败",
			zap.Uint("order_id", id),
			zap.Error(err))
		return err
	}

	// 业务指标：订单完成（暂时注释，等待 metrics 更新）
	// metrics.OrderCompleted.Inc()

	s.logger.Info("订单完成成功",
		zap.Uint("order_id", id),
		zap.String("order_no", order.OrderNo))

	return nil
}
