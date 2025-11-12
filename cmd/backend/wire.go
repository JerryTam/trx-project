//go:build wireinject
// +build wireinject

package main

import (
	"trx-project/internal/api/handler/backendHandler"
	"trx-project/internal/repository"
	"trx-project/internal/service"
	"trx-project/pkg/cache"
	"trx-project/pkg/config"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
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

		// RBAC Cache
		cache.NewRBACCache,

		// JWT Config
		provideAdminJWTConfig,

		// Repository
		repository.NewUserRepository,
		repository.NewRBACRepository,

		// Service
		service.NewUserService,
		service.NewRBACService,

		// Handler
		backendHandler.NewAdminUserHandler,
		backendHandler.NewRBACHandler,

		// Backend Router
		provideBackendRouter,
	)
	return nil, nil, nil
}
