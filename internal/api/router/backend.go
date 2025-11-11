package router

import (
	"trx-project/internal/api/handler"
	"trx-project/internal/api/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SetupBackend sets up the backend router
func SetupBackend(
	adminUserHandler *handler.AdminUserHandler,
	jwtSecret string,
	logger *zap.Logger,
	mode string,
) *gin.Engine {
	// Set gin mode
	gin.SetMode(mode)

	r := gin.New()

	// Apply global middleware
	r.Use(middleware.Recovery(logger))
	r.Use(middleware.Logger(logger))
	r.Use(middleware.CORS())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "backend",
		})
	})

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// 所有后台接口都需要管理员认证
		admin := v1.Group("/admin")
		admin.Use(middleware.AdminAuth(jwtSecret, logger))
		{
			// 用户管理
			adminUsers := admin.Group("/users")
			{
				adminUsers.GET("", adminUserHandler.ListUsers)                         // 获取用户列表
				adminUsers.GET("/:id", adminUserHandler.GetUser)                       // 获取用户详情
				adminUsers.PUT("/:id/status", adminUserHandler.UpdateUserStatus)       // 更新用户状态
				adminUsers.DELETE("/:id", adminUserHandler.DeleteUser)                 // 删除用户
				adminUsers.POST("/:id/reset-password", adminUserHandler.ResetPassword) // 重置密码
			}

			// 统计信息
			adminStats := admin.Group("/statistics")
			{
				adminStats.GET("/users", adminUserHandler.GetStatistics) // 用户统计
			}

			// 后续可以添加更多管理功能
			// - 内容管理
			// - 订单管理
			// - 系统配置
			// - 操作日志
			// etc.
		}
	}

	return r
}
