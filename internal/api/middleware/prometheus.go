package middleware

import (
	"strconv"
	"time"
	"trx-project/pkg/metrics"

	"github.com/gin-gonic/gin"
)

// PrometheusMiddleware 返回 Prometheus 监控中间件
func PrometheusMiddleware(m *metrics.Metrics, serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// 记录请求大小
		if c.Request.ContentLength > 0 {
			m.HTTPRequestSize.WithLabelValues(serviceName, method, path).Observe(float64(c.Request.ContentLength))
		}

		// 处理请求
		c.Next()

		// 请求处理完成后记录指标
		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())

		// HTTP 请求总数
		m.HTTPRequestsTotal.WithLabelValues(serviceName, method, path, status).Inc()

		// HTTP 请求延迟
		m.HTTPRequestDuration.WithLabelValues(serviceName, method, path).Observe(duration)

		// 响应大小
		m.HTTPResponseSize.WithLabelValues(serviceName, method, path).Observe(float64(c.Writer.Size()))
	}
}
