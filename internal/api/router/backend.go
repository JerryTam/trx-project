package router

import (
	"trx-project/internal/api/handler"
	"trx-project/internal/api/middleware"
	"trx-project/internal/service"
	"trx-project/pkg/config"
	"trx-project/pkg/metrics"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.uber.org/zap"
)

// SetupBackend sets up the backend router
func SetupBackend(
	adminUserHandler *handler.AdminUserHandler,
	rbacHandler *handler.RBACHandler,
	rbacService service.RBACService,
	jwtSecret string,
	redisClient *redis.Client,
	cfg *config.Config,
	logger *zap.Logger,
	mode string,
) *gin.Engine {
	// Set gin mode
	gin.SetMode(mode)

	r := gin.New()

	// 创建 Prometheus 指标
	m := metrics.NewMetrics("trx")

	// Apply global middleware
	// 顺序很重要：Recovery → OpenTelemetry → RequestID → Prometheus → Logger → CORS
	r.Use(middleware.Recovery(logger))

	// OpenTelemetry 链路追踪
	if cfg.Tracing.Enabled {
		r.Use(otelgin.Middleware(cfg.Tracing.ServiceName))
		logger.Info("OpenTelemetry tracing enabled for backend")
	}

	r.Use(middleware.RequestID(logger))                  // 请求 ID 追踪
	r.Use(middleware.PrometheusMiddleware(m, "backend")) // Prometheus 监控
	r.Use(middleware.Logger(logger))                     // 日志记录（会包含请求 ID）
	r.Use(middleware.CORS())

	// 限流中间件
	if cfg.RateLimit.Enabled {
		logger.Info("Rate limiting enabled for backend",
			zap.String("global_rate", cfg.RateLimit.GlobalRate),
			zap.String("ip_rate", cfg.RateLimit.IPRate),
			zap.String("user_rate", cfg.RateLimit.UserRate))

		rateLimiter := middleware.NewRateLimiter(redisClient, logger)

		// 应用全局限流和 IP 限流
		if cfg.RateLimit.GlobalRate != "" {
			r.Use(rateLimiter.GlobalRateLimit(cfg.RateLimit.GlobalRate))
		}
		if cfg.RateLimit.IPRate != "" {
			r.Use(rateLimiter.IPRateLimit(cfg.RateLimit.IPRate))
		}
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "backend",
		})
	})

	// Prometheus metrics 端点
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Swagger 文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// 所有后台接口都需要管理员认证
		admin := v1.Group("/admin")
		admin.Use(middleware.AdminAuth(jwtSecret, logger))
		// 管理员用户级别限流（需要在认证中间件之后）
		if cfg.RateLimit.Enabled && cfg.RateLimit.UserRate != "" {
			rateLimiter := middleware.NewRateLimiter(redisClient, logger)
			admin.Use(rateLimiter.UserRateLimit(cfg.RateLimit.UserRate))
		}
		{
			// ==================== RBAC 管理 ====================
			rbac := admin.Group("/rbac")
			rbac.Use(middleware.RequirePermission("rbac:manage", rbacService, logger)) // 需要 RBAC 管理权限
			{
				// 角色管理
				rbac.GET("/roles", rbacHandler.ListRoles)                                // 获取角色列表
				rbac.GET("/roles/:id", rbacHandler.GetRole)                              // 获取角色详情
				rbac.POST("/roles", rbacHandler.CreateRole)                              // 创建角色
				rbac.POST("/roles/:id/permissions", rbacHandler.AssignPermissionsToRole) // 为角色分配权限

				// 权限管理
				rbac.GET("/permissions", rbacHandler.ListPermissions) // 获取权限列表
			}

			// ==================== 用户管理 ====================
			adminUsers := admin.Group("/users")
			{
				// 查看用户（需要 user:read 权限）
				adminUsers.GET("",
					middleware.RequirePermission("user:read", rbacService, logger),
					adminUserHandler.ListUsers)
				adminUsers.GET("/:id",
					middleware.RequirePermission("user:read", rbacService, logger),
					adminUserHandler.GetUser)

				// 修改用户（需要 user:write 权限）
				adminUsers.PUT("/:id/status",
					middleware.RequirePermission("user:write", rbacService, logger),
					adminUserHandler.UpdateUserStatus)
				adminUsers.POST("/:id/reset-password",
					middleware.RequirePermission("user:write", rbacService, logger),
					adminUserHandler.ResetPassword)

				// 删除用户（需要 user:delete 权限）
				adminUsers.DELETE("/:id",
					middleware.RequirePermission("user:delete", rbacService, logger),
					adminUserHandler.DeleteUser)

				// 用户角色管理（需要 rbac:manage 权限）
				adminUsers.POST("/:user_id/role",
					middleware.RequirePermission("rbac:manage", rbacService, logger),
					rbacHandler.AssignRoleToUser)
				adminUsers.GET("/:user_id/roles",
					middleware.RequirePermission("rbac:manage", rbacService, logger),
					rbacHandler.GetUserRoles)
				adminUsers.GET("/:user_id/permissions",
					middleware.RequirePermission("rbac:manage", rbacService, logger),
					rbacHandler.GetUserPermissions)
			}

			// ==================== 统计信息 ====================
			adminStats := admin.Group("/statistics")
			adminStats.Use(middleware.RequirePermission("statistics:read", rbacService, logger))
			{
				adminStats.GET("/users", adminUserHandler.GetStatistics) // 用户统计
			}
		}
	}

	return r
}
