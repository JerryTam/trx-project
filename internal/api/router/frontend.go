package router

import (
	"trx-project/internal/api/handler"
	"trx-project/internal/api/middleware"
	"trx-project/pkg/config"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

// SetupFrontend sets up the frontend router
func SetupFrontend(
	userHandler *handler.UserHandler,
	jwtSecret string,
	redisClient *redis.Client,
	cfg *config.Config,
	logger *zap.Logger,
	mode string,
) *gin.Engine {
	// Set gin mode
	gin.SetMode(mode)

	r := gin.New()

	// Apply global middleware
	// 顺序很重要：Recovery → RequestID → Logger → CORS
	r.Use(middleware.Recovery(logger))
	r.Use(middleware.RequestID(logger))  // 请求 ID 追踪
	r.Use(middleware.Logger(logger))     // 日志记录（会包含请求 ID）
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

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "frontend",
		})
	})

	// Swagger 文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 routes
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
