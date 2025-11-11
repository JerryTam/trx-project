package router

import (
	"trx-project/internal/api/handler"
	"trx-project/internal/api/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SetupFrontend sets up the frontend router
func SetupFrontend(
	userHandler *handler.UserHandler,
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
			"status": "ok",
			"service": "frontend",
		})
	})

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
			{
				usersAuth.GET("", userHandler.ListUsers)
				usersAuth.GET("/:id", userHandler.GetUser)
				usersAuth.DELETE("/:id", userHandler.DeleteUser)
			}
		}
	}

	return r
}

