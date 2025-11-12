package middleware

import (
	"trx-project/internal/service"
	"trx-project/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RequirePermission RBAC 权限检查中间件
// 要求用户必须拥有指定的权限才能访问
func RequirePermission(permissionCode string, rbacService service.RBACService, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文获取管理员 ID
		adminID, exists := c.Get("admin_id")
		if !exists {
			logger.Warn("Admin ID not found in context")
			response.Unauthorized(c, "Unauthorized")
			c.Abort()
			return
		}

		userID, ok := adminID.(uint)
		if !ok {
			logger.Error("Invalid admin ID type")
			response.InternalError(c, "Internal error")
			c.Abort()
			return
		}

		// 检查权限
		err := rbacService.CheckPermission(c.Request.Context(), userID, permissionCode)
		if err != nil {
			logger.Warn("Permission denied",
				zap.Uint("admin_id", userID),
				zap.String("permission", permissionCode),
				zap.Error(err))
			response.Forbidden(c, "Permission denied: "+permissionCode)
			c.Abort()
			return
		}

		logger.Debug("Permission granted",
			zap.Uint("admin_id", userID),
			zap.String("permission", permissionCode))

		c.Next()
	}
}

// RequireAnyPermission 要求用户拥有任一权限即可访问
func RequireAnyPermission(permissionCodes []string, rbacService service.RBACService, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		adminID, exists := c.Get("admin_id")
		if !exists {
			response.Unauthorized(c, "Unauthorized")
			c.Abort()
			return
		}

		userID, ok := adminID.(uint)
		if !ok {
			response.InternalError(c, "Internal error")
			c.Abort()
			return
		}

		// 检查是否拥有任一权限
		hasPermission := false
		for _, permCode := range permissionCodes {
			err := rbacService.CheckPermission(c.Request.Context(), userID, permCode)
			if err == nil {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			logger.Warn("Permission denied (require any)",
				zap.Uint("admin_id", userID),
				zap.Strings("permissions", permissionCodes))
			response.Forbidden(c, "Permission denied")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAllPermissions 要求用户拥有所有权限才能访问
func RequireAllPermissions(permissionCodes []string, rbacService service.RBACService, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		adminID, exists := c.Get("admin_id")
		if !exists {
			response.Unauthorized(c, "Unauthorized")
			c.Abort()
			return
		}

		userID, ok := adminID.(uint)
		if !ok {
			response.InternalError(c, "Internal error")
			c.Abort()
			return
		}

		// 检查是否拥有所有权限
		for _, permCode := range permissionCodes {
			err := rbacService.CheckPermission(c.Request.Context(), userID, permCode)
			if err != nil {
				logger.Warn("Permission denied (require all)",
					zap.Uint("admin_id", userID),
					zap.String("missing_permission", permCode))
				response.Forbidden(c, "Permission denied: "+permCode)
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
