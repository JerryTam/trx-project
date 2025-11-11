package middleware

import (
	"strings"
	"trx-project/pkg/jwt"
	"trx-project/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Auth 用户认证中间件（前台）
func Auth(jwtSecret string, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Header 获取 token
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			logger.Warn("Missing authorization token")
			response.Unauthorized(c, "Missing authorization token")
			c.Abort()
			return
		}

		// 移除 Bearer 前缀
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// 解析 JWT Token
		claims, err := jwt.ParseToken(tokenString, jwtSecret)
		if err != nil {
			logger.Warn("Invalid token",
				zap.Error(err),
				zap.String("path", c.Request.URL.Path))

			if err == jwt.ErrTokenExpired {
				response.BusinessError(c, response.CodeUserTokenExpired, "Token has expired")
			} else {
				response.Unauthorized(c, "Invalid token")
			}
			c.Abort()
			return
		}

		// 验证角色（用户 token 应该是 "user" 角色）
		if claims.Role != "user" {
			logger.Warn("Invalid user role",
				zap.String("role", claims.Role),
				zap.Uint("user_id", claims.UserID))
			response.Unauthorized(c, "Invalid token")
			c.Abort()
			return
		}

		// TODO: 可选 - 从数据库验证用户状态是否正常

		// 将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Set("token", tokenString)

		logger.Debug("User authenticated",
			zap.Uint("user_id", claims.UserID),
			zap.String("username", claims.Username),
			zap.String("path", c.Request.URL.Path))

		c.Next()
	}
}

// AdminAuth 管理员认证中间件（后台）
func AdminAuth(jwtSecret string, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Header 获取 token
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			logger.Warn("Admin: Missing authorization token")
			response.Unauthorized(c, "Missing authorization token")
			c.Abort()
			return
		}

		// 移除 Bearer 前缀
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// 解析 JWT Token
		claims, err := jwt.ParseToken(tokenString, jwtSecret)
		if err != nil {
			logger.Warn("Invalid admin token",
				zap.Error(err),
				zap.String("path", c.Request.URL.Path))

			if err == jwt.ErrTokenExpired {
				response.BusinessError(c, response.CodeUserTokenExpired, "Token has expired")
			} else {
				response.Unauthorized(c, "Invalid token")
			}
			c.Abort()
			return
		}

		// 验证角色（管理员 token 应该是 "admin" 或 "superadmin" 角色）
		if claims.Role != "admin" && claims.Role != "superadmin" {
			logger.Warn("Not an admin role",
				zap.String("role", claims.Role),
				zap.Uint("user_id", claims.UserID))
			response.Forbidden(c, "Admin access required")
			c.Abort()
			return
		}

		// TODO: 可选 - 从数据库验证管理员状态

		// 将管理员信息存入上下文
		c.Set("admin_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("admin_role", claims.Role)
		c.Set("token", tokenString)

		logger.Debug("Admin authenticated",
			zap.Uint("admin_id", claims.UserID),
			zap.String("username", claims.Username),
			zap.String("role", claims.Role),
			zap.String("path", c.Request.URL.Path))

		c.Next()
	}
}

// OptionalAuth 可选认证中间件
// 如果有 token 则验证，没有也不阻止
func OptionalAuth(jwtSecret string, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			// 没有 token，继续执行，但不设置用户信息
			c.Next()
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// 尝试解析 token
		claims, err := jwt.ParseToken(tokenString, jwtSecret)
		if err == nil {
			// Token 有效，设置用户信息
			c.Set("user_id", claims.UserID)
			c.Set("username", claims.Username)
			c.Set("role", claims.Role)
			c.Set("token", tokenString)
			logger.Debug("Optional auth: User identified",
				zap.Uint("user_id", claims.UserID),
				zap.String("username", claims.Username))
		} else {
			logger.Debug("Optional auth: Invalid token, continue as guest",
				zap.Error(err))
		}

		c.Next()
	}
}

// CheckPermission 检查权限中间件
func CheckPermission(permission string, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文获取管理员角色
		role, exists := c.Get("admin_role")
		if !exists {
			logger.Warn("Permission check: No admin role found")
			response.Forbidden(c, "Permission denied")
			c.Abort()
			return
		}

		// TODO: 实现真实的权限检查逻辑
		// 从数据库或缓存获取角色权限
		// 检查是否拥有指定权限

		// 临时实现：superadmin 拥有所有权限
		if role == "superadmin" {
			c.Next()
			return
		}

		// 其他角色需要检查具体权限
		logger.Warn("Permission check failed",
			zap.String("role", role.(string)),
			zap.String("required_permission", permission))
		response.Forbidden(c, "Permission denied")
		c.Abort()
	}
}

// GetUserID 从上下文获取用户ID
func GetUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	return userID.(uint), true
}

// GetAdminID 从上下文获取管理员ID
func GetAdminID(c *gin.Context) (uint, bool) {
	adminID, exists := c.Get("admin_id")
	if !exists {
		return 0, false
	}
	return adminID.(uint), true
}

// GetAdminRole 从上下文获取管理员角色
func GetAdminRole(c *gin.Context) (string, bool) {
	role, exists := c.Get("admin_role")
	if !exists {
		return "", false
	}
	return role.(string), true
}
