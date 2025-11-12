package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// RequestID 常量定义
const (
	// RequestIDKey 在 context 中存储请求 ID 的 key
	RequestIDKey = "request_id"
	// RequestIDHeader 请求头中请求 ID 的 key
	RequestIDHeader = "X-Request-ID"
)

// RequestID 请求 ID 中间件
// 功能：
// 1. 从请求头读取 X-Request-ID，如果不存在则生成新的 UUID
// 2. 将请求 ID 存储到 context 中
// 3. 将请求 ID 添加到响应头
// 4. 在日志中记录请求 ID
func RequestID(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 尝试从请求头获取请求 ID
		requestID := c.GetHeader(RequestIDHeader)

		// 2. 如果请求头中没有，生成新的 UUID
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// 3. 将请求 ID 存储到 context 中
		c.Set(RequestIDKey, requestID)

		// 4. 将请求 ID 添加到响应头
		c.Header(RequestIDHeader, requestID)

		// 5. 记录请求开始日志（包含请求 ID）
		logger.Debug("Request started",
			zap.String("request_id", requestID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("client_ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
		)

		// 继续处理请求
		c.Next()

		// 6. 记录请求结束日志（包含请求 ID 和状态码）
		logger.Debug("Request completed",
			zap.String("request_id", requestID),
			zap.Int("status_code", c.Writer.Status()),
			zap.Int("response_size", c.Writer.Size()),
		)
	}
}

// GetRequestID 从 context 中获取请求 ID
// 这是一个辅助函数，方便在其他地方获取当前请求的 ID
func GetRequestID(c *gin.Context) string {
	if requestID, exists := c.Get(RequestIDKey); exists {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return ""
}

// RequestIDToLogger 将请求 ID 添加到日志记录器
// 返回一个带有请求 ID 字段的新 logger
func RequestIDToLogger(c *gin.Context, logger *zap.Logger) *zap.Logger {
	requestID := GetRequestID(c)
	if requestID != "" {
		return logger.With(zap.String("request_id", requestID))
	}
	return logger
}
