package backendHandler

import (
	"strconv"
	"trx-project/internal/model"
	"trx-project/internal/service"
	"trx-project/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RBACHandler RBAC 处理器
type RBACHandler struct {
	rbacService service.RBACService
	logger      *zap.Logger
}

// NewRBACHandler 创建 RBAC 处理器
func NewRBACHandler(rbacService service.RBACService, logger *zap.Logger) *RBACHandler {
	return &RBACHandler{
		rbacService: rbacService,
		logger:      logger,
	}
}

// ListRoles 获取角色列表
//
//	@Summary		获取角色列表
//	@Description	获取所有角色列表
//	@Tags			RBAC管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	response.Response{data=[]model.Role}	"成功获取角色列表"
//	@Failure		401	{object}	response.Response						"未授权"
//	@Failure		403	{object}	response.Response						"无权限"
//	@Failure		500	{object}	response.Response						"服务器内部错误"
//	@Router			/admin/rbac/roles [get]
func (h *RBACHandler) ListRoles(c *gin.Context) {
	roles, err := h.rbacService.ListRoles(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to list roles", zap.Error(err))
		response.InternalError(c, "Failed to list roles")
		return
	}

	response.Success(c, roles)
}

// GetRole 获取角色详情
//
//	@Summary		获取角色详情
//	@Description	获取指定角色的详细信息，包括权限列表
//	@Tags			RBAC管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int									true	"角色ID"
//	@Success		200	{object}	response.Response{data=model.Role}	"成功获取角色详情"
//	@Failure		400	{object}	response.Response					"无效的角色ID"
//	@Failure		401	{object}	response.Response					"未授权"
//	@Failure		403	{object}	response.Response					"无权限"
//	@Failure		404	{object}	response.Response					"角色不存在"
//	@Failure		500	{object}	response.Response					"服务器内部错误"
//	@Router			/admin/rbac/roles/{id} [get]
func (h *RBACHandler) GetRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid role ID")
		return
	}

	role, err := h.rbacService.GetRoleWithPermissions(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to get role", zap.Error(err))
		response.NotFound(c, "Role not found")
		return
	}

	response.Success(c, role)
}

// CreateRole 创建角色
//
//	@Summary		创建角色
//	@Description	创建新的角色
//	@Tags			RBAC管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		object{name=string,display_name=string,description=string}	true	"角色信息"
//	@Success		201		{object}	response.Response{data=model.Role}							"创建成功"
//	@Failure		400		{object}	response.Response											"请求参数错误"
//	@Failure		401		{object}	response.Response											"未授权"
//	@Failure		403		{object}	response.Response											"无权限"
//	@Failure		409		{object}	response.Response											"角色名已存在"
//	@Failure		500		{object}	response.Response											"服务器内部错误"
//	@Router			/admin/rbac/roles [post]
func (h *RBACHandler) CreateRole(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		DisplayName string `json:"display_name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidateError(c, err.Error())
		return
	}

	role := &model.Role{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		Status:      1,
	}

	if err := h.rbacService.CreateRole(c.Request.Context(), role); err != nil {
		h.logger.Error("Failed to create role", zap.Error(err))
		if err.Error() == "role name already exists" {
			response.BusinessError(c, response.CodeUserAlreadyExists, err.Error())
			return
		}
		response.InternalError(c, "Failed to create role")
		return
	}

	response.CreatedWithMsg(c, "Role created successfully", role)
}

// AssignPermissionsToRole 为角色分配权限
//
//	@Summary		为角色分配权限
//	@Description	为指定角色分配权限列表
//	@Tags			RBAC管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int								true	"角色ID"
//	@Param			request	body		object{permission_ids=[]int}	true	"权限ID列表"
//	@Success		200		{object}	response.Response				"分配成功"
//	@Failure		400		{object}	response.Response				"请求参数错误"
//	@Failure		401		{object}	response.Response				"未授权"
//	@Failure		403		{object}	response.Response				"无权限"
//	@Failure		500		{object}	response.Response				"服务器内部错误"
//	@Router			/admin/rbac/roles/{id}/permissions [post]
func (h *RBACHandler) AssignPermissionsToRole(c *gin.Context) {
	idStr := c.Param("id")
	roleID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid role ID")
		return
	}

	var req struct {
		PermissionIDs []uint `json:"permission_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidateError(c, err.Error())
		return
	}

	if err := h.rbacService.AssignPermissionsToRole(c.Request.Context(), uint(roleID), req.PermissionIDs); err != nil {
		h.logger.Error("Failed to assign permissions to role", zap.Error(err))
		response.InternalError(c, "Failed to assign permissions")
		return
	}

	response.SuccessWithMsg(c, "Permissions assigned successfully", nil)
}

// ListPermissions 获取权限列表
//
//	@Summary		获取权限列表
//	@Description	获取所有权限列表
//	@Tags			RBAC管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	response.Response{data=[]model.Permission}	"成功获取权限列表"
//	@Failure		401	{object}	response.Response							"未授权"
//	@Failure		403	{object}	response.Response							"无权限"
//	@Failure		500	{object}	response.Response							"服务器内部错误"
//	@Router			/admin/rbac/permissions [get]
func (h *RBACHandler) ListPermissions(c *gin.Context) {
	permissions, err := h.rbacService.ListPermissions(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to list permissions", zap.Error(err))
		response.InternalError(c, "Failed to list permissions")
		return
	}

	response.Success(c, permissions)
}

// AssignRoleToUser 为用户分配角色
//
//	@Summary		为用户分配角色
//	@Description	为指定用户分配角色
//	@Tags			RBAC管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int					true	"用户ID"
//	@Param			request	body		object{role_id=int}	true	"角色ID"
//	@Success		200		{object}	response.Response	"分配成功"
//	@Failure		400		{object}	response.Response	"请求参数错误"
//	@Failure		401		{object}	response.Response	"未授权"
//	@Failure		403		{object}	response.Response	"无权限"
//	@Failure		500		{object}	response.Response	"服务器内部错误"
//	@Router			/admin/users/{id}/role [post]
func (h *RBACHandler) AssignRoleToUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	var req struct {
		RoleID uint `json:"role_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidateError(c, err.Error())
		return
	}

	if err := h.rbacService.AssignRoleToUser(c.Request.Context(), uint(userID), req.RoleID); err != nil {
		h.logger.Error("Failed to assign role to user", zap.Error(err))
		response.InternalError(c, "Failed to assign role")
		return
	}

	response.SuccessWithMsg(c, "Role assigned successfully", nil)
}

// GetUserRoles 获取用户的角色列表
//
//	@Summary		获取用户角色
//	@Description	获取指定用户的角色列表
//	@Tags			RBAC管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int										true	"用户ID"
//	@Success		200	{object}	response.Response{data=[]model.Role}	"成功获取用户角色"
//	@Failure		400	{object}	response.Response						"无效的用户ID"
//	@Failure		401	{object}	response.Response						"未授权"
//	@Failure		403	{object}	response.Response						"无权限"
//	@Failure		500	{object}	response.Response						"服务器内部错误"
//	@Router			/admin/users/{id}/roles [get]
func (h *RBACHandler) GetUserRoles(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	roles, err := h.rbacService.GetUserRoles(c.Request.Context(), uint(userID))
	if err != nil {
		h.logger.Error("Failed to get user roles", zap.Error(err))
		response.InternalError(c, "Failed to get user roles")
		return
	}

	response.Success(c, roles)
}

// GetUserPermissions 获取用户的权限列表
//
//	@Summary		获取用户权限
//	@Description	获取指定用户的所有权限（通过角色继承）
//	@Tags			RBAC管理
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int											true	"用户ID"
//	@Success		200	{object}	response.Response{data=[]model.Permission}	"成功获取用户权限"
//	@Failure		400	{object}	response.Response							"无效的用户ID"
//	@Failure		401	{object}	response.Response							"未授权"
//	@Failure		403	{object}	response.Response							"无权限"
//	@Failure		500	{object}	response.Response							"服务器内部错误"
//	@Router			/admin/users/{id}/permissions [get]
func (h *RBACHandler) GetUserPermissions(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	permissions, err := h.rbacService.GetUserPermissions(c.Request.Context(), uint(userID))
	if err != nil {
		h.logger.Error("Failed to get user permissions", zap.Error(err))
		response.InternalError(c, "Failed to get user permissions")
		return
	}

	response.Success(c, permissions)
}
