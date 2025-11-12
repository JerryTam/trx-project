package frontendHandler

import (
	"strconv"
	_ "trx-project/internal/model" // 用于 Swagger 文档生成
	"trx-project/internal/service"
	"trx-project/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserHandler struct {
	service service.UserService
	logger  *zap.Logger
}

// NewUserHandler 创建新的用户处理器
func NewUserHandler(service service.UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		service: service,
		logger:  logger,
	}
}

// RegisterRequest 用户注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50" example:"testuser"` // 用户名，3-50个字符
	Email    string `json:"email" binding:"required,email" example:"test@example.com"`   // 邮箱地址
	Password string `json:"password" binding:"required,min=6" example:"password123"`     // 密码，最少6个字符
}

// LoginRequest 用户登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"testuser"`    // 用户名
	Password string `json:"password" binding:"required" example:"password123"` // 密码
}

// Register 用户注册
//
//	@Summary		用户注册
//	@Description	创建新用户账号，注册成功后返回用户信息和 JWT Token
//	@Tags			公开接口
//	@Accept			json
//	@Produce		json
//	@Param			request	body		RegisterRequest									true	"注册信息"
//	@Success		201		{object}	response.Response{data=map[string]interface{}}	"注册成功，返回用户信息和 Token"
//	@Failure		400		{object}	response.Response								"请求参数错误"
//	@Failure		409		{object}	response.Response								"用户名或邮箱已存在"
//	@Failure		500		{object}	response.Response								"服务器内部错误"
//	@Router			/public/register [post]
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

// Login 用户登录
//
//	@Summary		用户登录
//	@Description	使用用户名和密码登录，成功后返回用户信息和 JWT Token
//	@Tags			公开接口
//	@Accept			json
//	@Produce		json
//	@Param			request	body		LoginRequest									true	"登录信息"
//	@Success		200		{object}	response.Response{data=map[string]interface{}}	"登录成功，返回用户信息和 Token"
//	@Failure		400		{object}	response.Response								"请求参数错误"
//	@Failure		401		{object}	response.Response								"用户名或密码错误"
//	@Failure		403		{object}	response.Response								"账号已被禁用"
//	@Failure		500		{object}	response.Response								"服务器内部错误"
//	@Router			/public/login [post]
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

// GetUser 获取用户信息
//
//	@Summary		获取用户信息（前台）
//	@Description	根据用户ID获取用户信息，用户只能查看自己的信息
//	@Tags			用户接口
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int									true	"用户ID"
//	@Success		200	{object}	response.Response{data=model.User}	"成功获取用户信息"
//	@Failure		400	{object}	response.Response					"无效的用户ID"
//	@Failure		401	{object}	response.Response					"未授权"
//	@Failure		403	{object}	response.Response					"无权访问"
//	@Failure		404	{object}	response.Response					"用户不存在"
//	@Failure		500	{object}	response.Response					"服务器内部错误"
//	@Router			/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	// 前台用户只能获取自己的信息
	// TODO: 从认证中间件获取当前用户 ID 进行验证
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

// GetProfile 获取个人信息
//
//	@Summary		获取当前登录用户的个人信息
//	@Description	获取当前登录用户的详细信息
//	@Tags			用户接口
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	response.Response{data=model.User}	"成功获取个人信息"
//	@Failure		401	{object}	response.Response					"未授权或Token无效"
//	@Failure		404	{object}	response.Response					"用户不存在"
//	@Failure		500	{object}	response.Response					"服务器内部错误"
//	@Router			/user/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	// TODO: 从认证中间件获取当前用户 ID
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

// UpdateProfile 更新个人信息
//
//	@Summary		更新当前登录用户的个人信息
//	@Description	更新当前登录用户的个人信息（如邮箱等）
//	@Tags			用户接口
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		object{email=string}	true	"更新信息"
//	@Success		200		{object}	response.Response		"更新成功"
//	@Failure		400		{object}	response.Response		"请求参数错误"
//	@Failure		401		{object}	response.Response		"未授权或Token无效"
//	@Failure		500		{object}	response.Response		"服务器内部错误"
//	@Router			/user/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	// TODO: 从认证中间件获取当前用户 ID
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

// ListUsers 获取用户列表
//
//	@Summary		获取用户列表
//	@Description	分页获取用户列表（需要认证）
//	@Tags			用户接口
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			page		query		int												false	"页码，默认1"	default(1)
//	@Param			page_size	query		int												false	"每页数量，默认10"	default(10)
//	@Success		200			{object}	response.Response{data=map[string]interface{}}	"成功获取用户列表"
//	@Failure		400			{object}	response.Response								"请求参数错误"
//	@Failure		401			{object}	response.Response								"未授权"
//	@Failure		500			{object}	response.Response								"服务器内部错误"
//	@Router			/users [get]
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

// DeleteUser 处理用户删除
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
