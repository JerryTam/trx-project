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
)

func main() {
	// Load configuration
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
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

	// Get logger for output
	logger, _ := initLogger(cfg)

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

