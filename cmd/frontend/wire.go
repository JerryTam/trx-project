//go:build wireinject
// +build wireinject

package main

import (
	"trx-project/internal/api/handler/frontendhandler"
	"trx-project/internal/repository"
	"trx-project/internal/service"
	"trx-project/pkg/config"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
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
		frontendhandler.NewUserHandler,

		// Frontend Router
		provideFrontendRouter,
	)
	return nil, nil, nil
}
