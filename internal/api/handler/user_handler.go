package handler

import (
	"strconv"
	"trx-project/internal/service"
	"trx-project/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserHandler struct {
	service service.UserService
	logger  *zap.Logger
}

// NewUserHandler creates a new user handler
func NewUserHandler(service service.UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		service: service,
		logger:  logger,
	}
}

// RegisterRequest represents user registration request
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest represents user login request
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Register handles user registration
func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidateError(c, err.Error())
		return
	}

	user, token, err := h.service.Register(c.Request.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		h.logger.Error("Failed to register user", zap.Error(err))
		// 根据错误类型返回不同的响应
		if err.Error() == "username already exists" {
			response.BusinessError(c, response.CodeUserAlreadyExists, err.Error())
			return
		}
		if err.Error() == "email already exists" {
			response.BusinessError(c, response.CodeUserAlreadyExists, err.Error())
			return
		}
		response.InternalError(c, "Failed to register user")
		return
	}

	response.CreatedWithMsg(c, "User registered successfully", gin.H{
		"user":  user,
		"token": token,
	})
}

// Login handles user login
func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidateError(c, err.Error())
		return
	}

	user, token, err := h.service.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		h.logger.Error("Failed to login", zap.Error(err))
		// 根据错误类型返回不同的响应
		if err.Error() == "invalid username or password" {
			response.Unauthorized(c, err.Error())
			return
		}
		if err.Error() == "user account is inactive" {
			response.BusinessError(c, response.CodeUserDisabled, err.Error())
			return
		}
		response.InternalError(c, "Failed to login")
		return
	}

	response.SuccessWithMsg(c, "Login successful", gin.H{
		"user":  user,
		"token": token,
	})
}

// GetUser handles getting user by ID（前台 - 获取当前用户信息）
func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	// 前台用户只能获取自己的信息
	// TODO: 从认证中间件获取当前用户ID进行验证
	// userID, _ := middleware.GetUserID(c)
	// if uint(id) != userID {
	//     response.Forbidden(c, "You can only access your own information")
	//     return
	// }

	user, err := h.service.GetUserByID(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to get user", zap.Error(err))
		if err.Error() == "user not found" {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, "Failed to get user")
		return
	}

	response.Success(c, user)
}

// GetProfile 获取当前用户信息（前台 - 需要认证）
func (h *UserHandler) GetProfile(c *gin.Context) {
	// TODO: 从认证中间件获取当前用户ID
	// userID, exists := middleware.GetUserID(c)
	// if !exists {
	//     response.Unauthorized(c, "User not authenticated")
	//     return
	// }

	// 临时实现：从路径参数获取
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	user, err := h.service.GetUserByID(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to get profile", zap.Error(err))
		if err.Error() == "user not found" {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, "Failed to get profile")
		return
	}

	response.Success(c, user)
}

// UpdateProfile 更新当前用户信息（前台 - 需要认证）
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	// TODO: 从认证中间件获取当前用户ID
	// userID, exists := middleware.GetUserID(c)
	// if !exists {
	//     response.Unauthorized(c, "User not authenticated")
	//     return
	// }

	var req struct {
		Email string `json:"email" binding:"omitempty,email"`
		// 可以添加其他可更新的字段，如昵称、头像等
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidateError(c, err.Error())
		return
	}

	// TODO: 实现更新逻辑
	response.SuccessWithMsg(c, "Profile updated successfully", nil)
}

// ListUsers handles listing users
func (h *UserHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	users, total, err := h.service.ListUsers(c.Request.Context(), page, pageSize)
	if err != nil {
		h.logger.Error("Failed to list users", zap.Error(err))
		response.InternalError(c, "Failed to list users")
		return
	}

	response.PageSuccess(c, users, total, page, pageSize)
}

// DeleteUser handles user deletion
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	if err := h.service.DeleteUser(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("Failed to delete user", zap.Error(err))
		response.InternalError(c, "Failed to delete user")
		return
	}

	response.SuccessWithMsg(c, "User deleted successfully", nil)
}

