package handler

import (
	"strconv"
	"trx-project/internal/api/middleware"
	"trx-project/internal/service"
	"trx-project/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AdminUserHandler 管理员用户管理处理器
type AdminUserHandler struct {
	service service.UserService
	logger  *zap.Logger
}

// NewAdminUserHandler 创建管理员用户管理处理器
func NewAdminUserHandler(service service.UserService, logger *zap.Logger) *AdminUserHandler {
	return &AdminUserHandler{
		service: service,
		logger:  logger,
	}
}

// ListUsers 获取用户列表（后台）
func (h *AdminUserHandler) ListUsers(c *gin.Context) {
	// 获取管理员信息
	adminID, _ := middleware.GetAdminID(c)
	h.logger.Info("Admin listing users", zap.Uint("admin_id", adminID))

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")   // 可选的状态筛选
	keyword := c.Query("keyword") // 可选的关键词搜索

	h.logger.Debug("Admin list users params",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("status", status),
		zap.String("keyword", keyword))

	users, total, err := h.service.ListUsers(c.Request.Context(), page, pageSize)
	if err != nil {
		h.logger.Error("Admin failed to list users", zap.Error(err))
		response.InternalError(c, "Failed to list users")
		return
	}

	response.PageSuccess(c, users, total, page, pageSize)
}

// GetUser 获取用户详情（后台）
func (h *AdminUserHandler) GetUser(c *gin.Context) {
	adminID, _ := middleware.GetAdminID(c)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	h.logger.Info("Admin getting user detail",
		zap.Uint("admin_id", adminID),
		zap.Uint64("user_id", id))

	user, err := h.service.GetUserByID(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Admin failed to get user", zap.Error(err))
		if err.Error() == "user not found" {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, "Failed to get user")
		return
	}

	response.Success(c, user)
}

// UpdateUserStatus 更新用户状态（后台）
func (h *AdminUserHandler) UpdateUserStatus(c *gin.Context) {
	adminID, _ := middleware.GetAdminID(c)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	var req struct {
		Status int `json:"status" binding:"required,oneof=0 1"` // 0: 禁用, 1: 启用
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidateError(c, err.Error())
		return
	}

	h.logger.Info("Admin updating user status",
		zap.Uint("admin_id", adminID),
		zap.Uint64("user_id", id),
		zap.Int("status", req.Status))

	user, err := h.service.GetUserByID(c.Request.Context(), uint(id))
	if err != nil {
		if err.Error() == "user not found" {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, "Failed to get user")
		return
	}

	user.Status = req.Status
	if err := h.service.UpdateUser(c.Request.Context(), user); err != nil {
		h.logger.Error("Admin failed to update user status", zap.Error(err))
		response.InternalError(c, "Failed to update user status")
		return
	}

	response.SuccessWithMsg(c, "User status updated successfully", user)
}

// DeleteUser 删除用户（后台）
func (h *AdminUserHandler) DeleteUser(c *gin.Context) {
	adminID, _ := middleware.GetAdminID(c)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	h.logger.Info("Admin deleting user",
		zap.Uint("admin_id", adminID),
		zap.Uint64("user_id", id))

	if err := h.service.DeleteUser(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("Admin failed to delete user", zap.Error(err))
		response.InternalError(c, "Failed to delete user")
		return
	}

	response.SuccessWithMsg(c, "User deleted successfully", nil)
}

// GetStatistics 获取用户统计信息（后台）
func (h *AdminUserHandler) GetStatistics(c *gin.Context) {
	adminID, _ := middleware.GetAdminID(c)
	h.logger.Info("Admin getting user statistics", zap.Uint("admin_id", adminID))

	// TODO: 实现真实的统计逻辑
	stats := gin.H{
		"total_users":     1000,
		"active_users":    850,
		"inactive_users":  150,
		"new_users_today": 12,
		"new_users_week":  89,
		"new_users_month": 356,
	}

	response.Success(c, stats)
}

// ResetPassword 重置用户密码（后台）
func (h *AdminUserHandler) ResetPassword(c *gin.Context) {
	adminID, _ := middleware.GetAdminID(c)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	var req struct {
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidateError(c, err.Error())
		return
	}

	h.logger.Info("Admin resetting user password",
		zap.Uint("admin_id", adminID),
		zap.Uint64("user_id", id))

	// TODO: 实现密码重置逻辑
	// 1. 验证用户存在
	// 2. 加密新密码
	// 3. 更新用户密码
	// 4. 记录操作日志
	// 5. 可选：发送通知给用户

	response.SuccessWithMsg(c, "Password reset successfully", nil)
}
