package main

import (
	"time"
	"trx-project/internal/api/handler/frontendHandler"
	"trx-project/internal/api/router"
	"trx-project/pkg/cache"
	"trx-project/pkg/config"
	"trx-project/pkg/database"
	"trx-project/pkg/jwt"
	"trx-project/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

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

func provideJWTConfig(cfg *config.Config) jwt.Config {
	return jwt.Config{
		Secret:     cfg.JWT.Secret,
		Issuer:     cfg.JWT.Issuer,
		ExpireTime: time.Duration(cfg.JWT.ExpireHours) * time.Hour,
	}
}

func provideFrontendRouter(
	userHandler *frontendHandler.UserHandler,
	orderHandler *frontendHandler.OrderHandler,
	redisClient *redis.Client,
	logger *zap.Logger,
	cfg *config.Config,
) *gin.Engine {
	return router.SetupFrontend(
		userHandler,
		orderHandler,
		cfg.JWT.Secret,
		redisClient,
		cfg,
		logger,
		cfg.Server.Mode,
	)
}
