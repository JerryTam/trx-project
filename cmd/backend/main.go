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
	"trx-project/pkg/tracing"

	"go.uber.org/zap"
)

func main() {
	// 加载配置
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 优先初始化日志记录器用于迁移
	logger, err := initLogger(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	// 初始化 OpenTelemetry 链路追踪
	tracingCleanup, err := tracing.InitTracer(&tracing.Config{
		ServiceName:    cfg.Tracing.ServiceName + "-backend",
		ServiceVersion: cfg.Tracing.ServiceVersion,
		Environment:    cfg.Server.Env,
		JaegerEndpoint: cfg.Tracing.JaegerEndpoint,
		Enabled:        cfg.Tracing.Enabled,
	}, logger)
	if err != nil {
		log.Fatalf("Failed to initialize tracing: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := tracingCleanup(ctx); err != nil {
			logger.Error("Failed to cleanup tracing", zap.Error(err))
		}
	}()

	// 如果启用了 AUTO_MIGRATE，则运行数据库迁移
	if os.Getenv("AUTO_MIGRATE") == "true" {
		logger.Info("AUTO_MIGRATE enabled, running database migrations...")
		if err := runMigrations(cfg, logger); err != nil {
			log.Fatalf("Failed to run migrations: %v", err)
		}
		logger.Info("Database migrations completed successfully")
	}

	// 使用 Wire 初始化应用
	router, cleanup, err := initBackendApp(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize backend app: %v", err)
	}
	defer cleanup()

	// 后端使用不同的端口
	backendPort := cfg.Server.Port + 1 // 前端 8080，后端 8081
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, backendPort)

	srv := &http.Server{
		Addr:           addr,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// 在 goroutine 中启动服务器
	go func() {
		logger.Sugar().Infof("Backend server starting on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start backend server: %v", err)
		}
	}()

	// 等待中断信号以优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Sugar().Info("Shutting down backend server...")

	// 带超时的优雅关闭
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
