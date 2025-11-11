package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"trx-project/pkg/config"
	"trx-project/pkg/migrate"

	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger first for migration
	logger, err := initLogger(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	// Run database migrations if AUTO_MIGRATE is enabled
	if os.Getenv("AUTO_MIGRATE") == "true" {
		logger.Info("AUTO_MIGRATE enabled, running database migrations...")
		if err := runMigrations(cfg, logger); err != nil {
			log.Fatalf("Failed to run migrations: %v", err)
		}
		logger.Info("Database migrations completed successfully")
	}

	// Initialize app with Wire
	router, cleanup, err := initBackendApp(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize backend app: %v", err)
	}
	defer cleanup()

	// Backend uses different port
	backendPort := cfg.Server.Port + 1 // Frontend 8080, Backend 8081
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, backendPort)

	srv := &http.Server{
		Addr:           addr,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// Start server in a goroutine
	go func() {
		logger.Sugar().Infof("Backend server starting on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start backend server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Sugar().Info("Shutting down backend server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Backend server forced to shutdown: %v", err)
	}

	logger.Sugar().Info("Backend server exited")
}

// runMigrations 运行数据库迁移
func runMigrations(cfg *config.Config, logger *zap.Logger) error {
	// 构建数据库连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=true",
		cfg.Database.MySQL.Username,
		cfg.Database.MySQL.Password,
		cfg.Database.MySQL.Host,
		cfg.Database.MySQL.Port,
		cfg.Database.MySQL.Database,
	)

	// 创建迁移器
	migrator, err := migrate.NewMigrator(&migrate.Config{
		MigrationsPath: "file://migrations",
		DatabaseURL:    dsn,
		Logger:         logger,
	})
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}
	defer migrator.Close()

	// 执行迁移
	return migrator.Up()
}
