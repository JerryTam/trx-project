package frontendHandler

import (
	"strconv"
	_ "trx-project/internal/model" // 用于 Swagger 文档生成
	"trx-project/internal/service"
	"trx-project/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// OrderHandler 订单处理器
type OrderHandler struct {
	orderService service.OrderService
	logger       *zap.Logger
}

// NewOrderHandler 创建订单处理器实例
func NewOrderHandler(orderService service.OrderService, logger *zap.Logger) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
		logger:       logger,
	}
}

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	ProductName  string  `json:"product_name" binding:"required" example:"iPhone 15 Pro"`
	ProductPrice float64 `json:"product_price" binding:"required,gt=0" example:"7999"`
	Quantity     int     `json:"quantity" binding:"required,gt=0" example:"1"`
	Remark       string  `json:"remark" example:"请尽快发货"`
}

// CreateOrder 创建订单
// @Summary 创建订单
// @Description 用户创建新订单
// @Tags 订单管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body CreateOrderRequest true "订单信息"
// @Success 200 {object} response.Response{data=model.Order} "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未登录"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /user/orders [post]
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	// 获取当前登录用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "请先登录")
		return
	}

	// 绑定请求参数
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	// 调用服务层创建订单
	order, err := h.orderService.CreateOrder(
		userID.(uint),
		req.ProductName,
		req.ProductPrice,
		req.Quantity,
		req.Remark,
	)

	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, "订单创建成功", order)
}

// GetOrder 获取订单详情
// @Summary 获取订单详情
// @Description 根据订单ID获取订单详情
// @Tags 订单管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "订单ID"
// @Success 200 {object} response.Response{data=model.Order} "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未登录"
// @Failure 404 {object} response.Response "订单不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /user/orders/{id} [get]
func (h *OrderHandler) GetOrder(c *gin.Context) {
	// 获取当前登录用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "请先登录")
		return
	}

	// 获取订单ID
	orderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "订单ID格式错误")
		return
	}

	// 调用服务层获取订单
	order, err := h.orderService.GetOrderByID(uint(orderID), userID.(uint))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, "获取成功", order)
}

// GetOrders 获取用户订单列表
// @Summary 获取用户订单列表
// @Description 获取当前用户的订单列表（分页）
// @Tags 订单管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} response.Response{data=response.PaginationData} "获取成功"
// @Failure 401 {object} response.Response "未登录"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /user/orders [get]
func (h *OrderHandler) GetOrders(c *gin.Context) {
	// 获取当前登录用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "请先登录")
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// 调用服务层获取订单列表
	orders, total, err := h.orderService.GetUserOrders(userID.(uint), page, pageSize)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	// 返回分页数据
	response.Pagination(c, "获取成功", orders, total, page, pageSize)
}

// PayOrder 支付订单
// @Summary 支付订单
// @Description 用户支付订单
// @Tags 订单管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "订单ID"
// @Success 200 {object} response.Response "支付成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未登录"
// @Failure 404 {object} response.Response "订单不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /user/orders/{id}/pay [post]
func (h *OrderHandler) PayOrder(c *gin.Context) {
	// 获取当前登录用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "请先登录")
		return
	}

	// 获取订单ID
	orderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "订单ID格式错误")
		return
	}

	// 调用服务层支付订单
	if err := h.orderService.PayOrder(uint(orderID), userID.(uint)); err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, "支付成功", nil)
}

// CancelOrder 取消订单
// @Summary 取消订单
// @Description 用户取消订单
// @Tags 订单管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "订单ID"
// @Success 200 {object} response.Response "取消成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未登录"
// @Failure 404 {object} response.Response "订单不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /user/orders/{id}/cancel [post]
func (h *OrderHandler) CancelOrder(c *gin.Context) {
	// 获取当前登录用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "请先登录")
		return
	}

	// 获取订单ID
	orderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "订单ID格式错误")
		return
	}

	// 调用服务层取消订单
	if err := h.orderService.CancelOrder(uint(orderID), userID.(uint)); err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, "取消成功", nil)
}

// CompleteOrder 完成订单
// @Summary 完成订单
// @Description 用户确认收货，完成订单
// @Tags 订单管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "订单ID"
// @Success 200 {object} response.Response "完成成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未登录"
// @Failure 404 {object} response.Response "订单不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /user/orders/{id}/complete [post]
func (h *OrderHandler) CompleteOrder(c *gin.Context) {
	// 获取当前登录用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "请先登录")
		return
	}

	// 获取订单ID
	orderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "订单ID格式错误")
		return
	}

	// 调用服务层完成订单
	if err := h.orderService.CompleteOrder(uint(orderID), userID.(uint)); err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, "订单已完成", nil)
}

