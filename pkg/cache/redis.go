package cache

import (
	"context"
	"fmt"
	"trx-project/pkg/config"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// InitRedis initializes Redis client
func InitRedis(cfg *config.RedisConfig, logger *zap.Logger) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.GetAddress(),
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
	})

	// Test connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	logger.Info("Redis connected successfully")
	return client, nil
}

// CloseRedis closes Redis client
func CloseRedis(client *redis.Client) error {
	return client.Close()
}

