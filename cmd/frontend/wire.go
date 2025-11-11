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

// initFrontendApp initializes the frontend application with all dependencies
func initFrontendApp(cfg *config.Config) (*gin.Engine, func(), error) {
	wire.Build(
		// Logger
		provideLogger,

		// Database
		provideDB,

		// Redis
		provideRedis,

		// JWT Config
		provideJWTConfig,

		// Repository
		repository.NewUserRepository,

		// Service
		service.NewUserService,

		// Handler
		handler.NewUserHandler,

		// Frontend Router
		provideFrontendRouter,
	)
	return nil, nil, nil
}

// initLogger initializes logger for output
func initLogger(cfg *config.Config) (*zap.Logger, error) {
	return provideLogger(cfg)
}

// initDatabaseAndLogger initializes database and logger for migration
func initDatabaseAndLogger(cfg *config.Config) (*gorm.DB, *zap.Logger, error) {
	zapLogger, err := provideLogger(cfg)
	if err != nil {
		return nil, nil, err
	}
	db, err := provideDB(cfg, zapLogger)
	if err != nil {
		return nil, nil, err
	}
	return db, zapLogger, nil
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

func provideJWTConfig(cfg *config.Config) jwt.Config {
	return jwt.Config{
		Secret:     cfg.JWT.Secret,
		Issuer:     cfg.JWT.Issuer,
		ExpireTime: time.Duration(cfg.JWT.ExpireHours) * time.Hour,
	}
}

func provideFrontendRouter(
	userHandler *handler.UserHandler,
	redisClient *redis.Client,
	logger *zap.Logger,
	cfg *config.Config,
) *gin.Engine {
	return router.SetupFrontend(
		userHandler,
		cfg.JWT.Secret,
		redisClient,
		cfg,
		logger,
		cfg.Server.Mode,
	)
}

