//go:build wireinject
// +build wireinject

package main

import (
	"time"
	"trx-project/internal/api/handler"
	"trx-project/internal/api/router"
	"trx-project/internal/repository"
	"trx-project/internal/service"
	"trx-project/pkg/cache"
	"trx-project/pkg/config"
	"trx-project/pkg/database"
	"trx-project/pkg/jwt"
	"trx-project/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// initBackendApp initializes the backend application with all dependencies
func initBackendApp(cfg *config.Config) (*gin.Engine, func(), error) {
	wire.Build(
		// Logger
		provideLogger,

		// Database
		provideDB,

		// Redis
		provideRedis,

		// JWT Config
		provideAdminJWTConfig,

		// Repository
		repository.NewUserRepository,
		repository.NewRBACRepository,

		// Service
		service.NewUserService,
		service.NewRBACService,

		// Handler
		handler.NewAdminUserHandler,
		handler.NewRBACHandler,

		// Backend Router
		provideBackendRouter,
	)
	return nil, nil, nil
}

// initLogger initializes logger for output
func initLogger(cfg *config.Config) (*zap.Logger, error) {
	return provideLogger(cfg)
}

func provideLogger(cfg *config.Config) (*zap.Logger, error) {
	if err := logger.InitLogger(&cfg.Logger); err != nil {
		return nil, err
	}
	return logger.Logger, nil
}

func provideDB(cfg *config.Config, logger *zap.Logger) (*gorm.DB, error) {
	return database.InitMySQL(&cfg.Database.MySQL, logger)
}

func provideRedis(cfg *config.Config, logger *zap.Logger) (*redis.Client, error) {
	return cache.InitRedis(&cfg.Redis, logger)
}

func provideAdminJWTConfig(cfg *config.Config) jwt.Config {
	return jwt.Config{
		Secret:     cfg.JWT.Secret,
		Issuer:     cfg.JWT.Issuer,
		ExpireTime: time.Duration(cfg.JWT.AdminExpireHours) * time.Hour,
	}
}

func provideBackendRouter(
	adminUserHandler *handler.AdminUserHandler,
	rbacHandler *handler.RBACHandler,
	rbacService service.RBACService,
	logger *zap.Logger,
	cfg *config.Config,
) *gin.Engine {
	return router.SetupBackend(adminUserHandler, rbacHandler, rbacService, cfg.JWT.Secret, logger, cfg.Server.Mode)
}

