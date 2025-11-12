package router

import (
	"trx-project/internal/api/handler/frontendHandler"
	"trx-project/internal/api/middleware"
	"trx-project/pkg/config"
	"trx-project/pkg/metrics"

	_ "trx-project/cmd/frontend/docs" // Swagger 文档

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.uber.org/zap"
)

// SetupFrontend 设置前端路由器
func SetupFrontend(
	userHandler *frontendHandler.UserHandler,
	jwtSecret string,
	redisClient *redis.Client,
	cfg *config.Config,
	logger *zap.Logger,
	mode string,
) *gin.Engine {
	// 设置 gin 模式
	gin.SetMode(mode)

	r := gin.New()

	// 创建 Prometheus 指标
	m := metrics.NewMetrics("trx")

	// 应用全局中间件
	// 顺序很重要：Recovery → OpenTelemetry → RequestID → Prometheus → Logger → CORS
	r.Use(middleware.Recovery(logger))

	// OpenTelemetry 链路追踪
	if cfg.Tracing.Enabled {
		r.Use(otelgin.Middleware(cfg.Tracing.ServiceName))
		logger.Info("OpenTelemetry tracing enabled for frontend")
	}

	r.Use(middleware.RequestID(logger))                   // 请求 ID 追踪
	r.Use(middleware.PrometheusMiddleware(m, "frontend")) // Prometheus 监控
	r.Use(middleware.Logger(logger))                      // 日志记录（会包含请求 ID）
	r.Use(middleware.CORS())

	// 限流中间件
	if cfg.RateLimit.Enabled {
		logger.Info("Rate limiting enabled",
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

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "frontend",
		})
	})

	// Prometheus metrics 端点
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Swagger 文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.InstanceName("frontend")))

	// API v1 路由
	v1 := r.Group("/api/v1")
	{
		// 公开接口（无需认证）
		public := v1.Group("/public")
		{
			public.POST("/register", userHandler.Register)
			public.POST("/login", userHandler.Login)
		}

		// 用户接口（需要用户认证）
		user := v1.Group("/user")
		user.Use(middleware.Auth(jwtSecret, logger))
		// 用户级别限流（需要在认证中间件之后）
		if cfg.RateLimit.Enabled && cfg.RateLimit.UserRate != "" {
			rateLimiter := middleware.NewRateLimiter(redisClient, logger)
			user.Use(rateLimiter.UserRateLimit(cfg.RateLimit.UserRate))
		}
		{
			user.GET("/profile", userHandler.GetProfile)
			user.PUT("/profile", userHandler.UpdateProfile)
		}

		// 兼容旧接口（临时保留）
		users := v1.Group("/users")
		{
			users.POST("/register", userHandler.Register)
			users.POST("/login", userHandler.Login)
			// 以下接口需要认证
			usersAuth := users.Group("")
			usersAuth.Use(middleware.Auth(jwtSecret, logger))
			// 用户级别限流
			if cfg.RateLimit.Enabled && cfg.RateLimit.UserRate != "" {
				rateLimiter := middleware.NewRateLimiter(redisClient, logger)
				usersAuth.Use(rateLimiter.UserRateLimit(cfg.RateLimit.UserRate))
			}
			{
				usersAuth.GET("", userHandler.ListUsers)
				usersAuth.GET("/:id", userHandler.GetUser)
				usersAuth.DELETE("/:id", userHandler.DeleteUser)
			}
		}
	}

	return r
}
