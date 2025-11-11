package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/ulule/limiter/v3"
	sredis "github.com/ulule/limiter/v3/drivers/store/redis"
	"go.uber.org/zap"
	"trx-project/internal/api/response"
)

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	// 全局限流：每秒请求数
	GlobalRate string
	// IP 限流：每分钟请求数
	IPRate string
	// 用户限流：每分钟请求数
	UserRate string
}

// RateLimiter 限流器
type RateLimiter struct {
	redis       *redis.Client
	logger      *zap.Logger
	globalStore limiter.Store
	ipStore     limiter.Store
	userStore   limiter.Store
}

// NewRateLimiter 创建限流器
func NewRateLimiter(redisClient *redis.Client, logger *zap.Logger) *RateLimiter {
	return &RateLimiter{
		redis:  redisClient,
		logger: logger,
	}
}

// initStore 初始化限流存储
func (rl *RateLimiter) initStore(prefix string) (limiter.Store, error) {
	store, err := sredis.NewStoreWithOptions(rl.redis, limiter.StoreOptions{
		Prefix:   prefix,
		MaxRetry: 3,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create rate limit store: %w", err)
	}
	return store, nil
}

// GlobalRateLimit 全局限流中间件
// 示例: "100-S" = 每秒 100 个请求, "1000-M" = 每分钟 1000 个请求
func (rl *RateLimiter) GlobalRateLimit(rateString string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if rl.globalStore == nil {
			store, err := rl.initStore("rate_limit:global")
			if err != nil {
				rl.logger.Error("Failed to init global rate limit store", zap.Error(err))
				c.Next()
				return
			}
			rl.globalStore = store
		}

		rate, err := limiter.NewRateFromFormatted(rateString)
		if err != nil {
			rl.logger.Error("Invalid rate format", zap.String("rate", rateString), zap.Error(err))
			c.Next()
			return
		}

		limiterInstance := limiter.New(rl.globalStore, rate)
		ctx := context.Background()
		limiterContext, err := limiterInstance.Get(ctx, "global")

		if err != nil {
			rl.logger.Error("Failed to get rate limit context", zap.Error(err))
			c.Next()
			return
		}

		// 设置响应头
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", limiterContext.Limit))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", limiterContext.Remaining))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", limiterContext.Reset))

		if limiterContext.Reached {
			rl.logger.Warn("Global rate limit exceeded",
				zap.String("path", c.Request.URL.Path))

			response.BusinessError(c, response.CodeTooManyRequests, 
				fmt.Sprintf("全局请求限流，请 %d 秒后重试", limiterContext.Reset))
			c.Abort()
			return
		}

		c.Next()
	}
}

// IPRateLimit 基于 IP 的限流中间件
// 示例: "100-M" = 每分钟 100 个请求
func (rl *RateLimiter) IPRateLimit(rateString string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if rl.ipStore == nil {
			store, err := rl.initStore("rate_limit:ip")
			if err != nil {
				rl.logger.Error("Failed to init IP rate limit store", zap.Error(err))
				c.Next()
				return
			}
			rl.ipStore = store
		}

		rate, err := limiter.NewRateFromFormatted(rateString)
		if err != nil {
			rl.logger.Error("Invalid rate format", zap.String("rate", rateString), zap.Error(err))
			c.Next()
			return
		}

		// 获取客户端 IP
		clientIP := c.ClientIP()
		key := fmt.Sprintf("ip:%s", clientIP)

		limiterInstance := limiter.New(rl.ipStore, rate)
		ctx := context.Background()
		limiterContext, err := limiterInstance.Get(ctx, key)

		if err != nil {
			rl.logger.Error("Failed to get rate limit context", zap.Error(err))
			c.Next()
			return
		}

		// 设置响应头
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", limiterContext.Limit))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", limiterContext.Remaining))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", limiterContext.Reset))

		if limiterContext.Reached {
			rl.logger.Warn("IP rate limit exceeded",
				zap.String("ip", clientIP),
				zap.String("path", c.Request.URL.Path))

			response.BusinessError(c, response.CodeTooManyRequests,
				fmt.Sprintf("IP 请求限流（%s），请 %d 秒后重试", clientIP, limiterContext.Reset))
			c.Abort()
			return
		}

		c.Next()
	}
}

// UserRateLimit 基于用户的限流中间件（需要在认证中间件之后使用）
// 示例: "1000-M" = 每分钟 1000 个请求
func (rl *RateLimiter) UserRateLimit(rateString string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if rl.userStore == nil {
			store, err := rl.initStore("rate_limit:user")
			if err != nil {
				rl.logger.Error("Failed to init user rate limit store", zap.Error(err))
				c.Next()
				return
			}
			rl.userStore = store
		}

		rate, err := limiter.NewRateFromFormatted(rateString)
		if err != nil {
			rl.logger.Error("Invalid rate format", zap.String("rate", rateString), zap.Error(err))
			c.Next()
			return
		}

		// 从上下文获取用户 ID（由认证中间件设置）
		userID, exists := c.Get("user_id")
		if !exists {
			// 如果没有用户 ID，跳过用户级别限流
			c.Next()
			return
		}

		key := fmt.Sprintf("user:%v", userID)

		limiterInstance := limiter.New(rl.userStore, rate)
		ctx := context.Background()
		limiterContext, err := limiterInstance.Get(ctx, key)

		if err != nil {
			rl.logger.Error("Failed to get rate limit context", zap.Error(err))
			c.Next()
			return
		}

		// 设置响应头
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", limiterContext.Limit))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", limiterContext.Remaining))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", limiterContext.Reset))

		if limiterContext.Reached {
			rl.logger.Warn("User rate limit exceeded",
				zap.Any("user_id", userID),
				zap.String("path", c.Request.URL.Path))

			response.BusinessError(c, response.CodeTooManyRequests,
				fmt.Sprintf("用户请求限流，请 %d 秒后重试", limiterContext.Reset))
			c.Abort()
			return
		}

		c.Next()
	}
}

// CustomRateLimit 自定义限流中间件，支持自定义 key 生成函数
func (rl *RateLimiter) CustomRateLimit(rateString string, keyFunc func(*gin.Context) string) gin.HandlerFunc {
	store, err := rl.initStore("rate_limit:custom")
	if err != nil {
		rl.logger.Error("Failed to init custom rate limit store", zap.Error(err))
		return func(c *gin.Context) {
			c.Next()
		}
	}

	rate, err := limiter.NewRateFromFormatted(rateString)
	if err != nil {
		rl.logger.Error("Invalid rate format", zap.String("rate", rateString), zap.Error(err))
		return func(c *gin.Context) {
			c.Next()
		}
	}

	return func(c *gin.Context) {
		key := keyFunc(c)
		if key == "" {
			c.Next()
			return
		}

		limiterInstance := limiter.New(store, rate)
		ctx := context.Background()
		limiterContext, err := limiterInstance.Get(ctx, key)

		if err != nil {
			rl.logger.Error("Failed to get rate limit context", zap.Error(err))
			c.Next()
			return
		}

		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", limiterContext.Limit))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", limiterContext.Remaining))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", limiterContext.Reset))

		if limiterContext.Reached {
			rl.logger.Warn("Custom rate limit exceeded",
				zap.String("key", key),
				zap.String("path", c.Request.URL.Path))

			response.BusinessError(c, response.CodeTooManyRequests,
				fmt.Sprintf("请求限流，请 %d 秒后重试", limiterContext.Reset))
			c.Abort()
			return
		}

		c.Next()
	}
}

// CombinedRateLimit 组合限流中间件：同时检查全局、IP 和用户限流
func (rl *RateLimiter) CombinedRateLimit(globalRate, ipRate, userRate string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 全局限流
		if globalRate != "" {
			if !rl.checkRateLimit(c, "global", "global", globalRate, &rl.globalStore, "rate_limit:global") {
				return
			}
		}

		// 2. IP 限流
		if ipRate != "" {
			clientIP := c.ClientIP()
			key := fmt.Sprintf("ip:%s", clientIP)
			if !rl.checkRateLimit(c, "IP", key, ipRate, &rl.ipStore, "rate_limit:ip") {
				return
			}
		}

		// 3. 用户限流（如果已认证）
		if userRate != "" {
			if userID, exists := c.Get("user_id"); exists {
				key := fmt.Sprintf("user:%v", userID)
				if !rl.checkRateLimit(c, "用户", key, userRate, &rl.userStore, "rate_limit:user") {
					return
				}
			}
		}

		c.Next()
	}
}

// checkRateLimit 检查限流的辅助函数
func (rl *RateLimiter) checkRateLimit(c *gin.Context, limitType, key, rateString string, store *limiter.Store, prefix string) bool {
	if *store == nil {
		s, err := rl.initStore(prefix)
		if err != nil {
			rl.logger.Error("Failed to init rate limit store", 
				zap.String("type", limitType), 
				zap.Error(err))
			return true
		}
		*store = s
	}

	rate, err := limiter.NewRateFromFormatted(rateString)
	if err != nil {
		rl.logger.Error("Invalid rate format", 
			zap.String("type", limitType),
			zap.String("rate", rateString), 
			zap.Error(err))
		return true
	}

	limiterInstance := limiter.New(*store, rate)
	ctx := context.Background()
	limiterContext, err := limiterInstance.Get(ctx, key)

	if err != nil {
		rl.logger.Error("Failed to get rate limit context", 
			zap.String("type", limitType),
			zap.Error(err))
		return true
	}

	c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", limiterContext.Limit))
	c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", limiterContext.Remaining))
	c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", limiterContext.Reset))

	if limiterContext.Reached {
		rl.logger.Warn("Rate limit exceeded",
			zap.String("type", limitType),
			zap.String("key", key),
			zap.String("path", c.Request.URL.Path))

		response.BusinessError(c, response.CodeTooManyRequests,
			fmt.Sprintf("%s请求限流，请 %d 秒后重试", limitType, limiterContext.Reset))
		c.Abort()
		return false
	}

	return true
}

