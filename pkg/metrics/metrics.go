package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics 包含所有 Prometheus 指标
type Metrics struct {
	// HTTP 请求指标
	HTTPRequestsTotal   *prometheus.CounterVec
	HTTPRequestDuration *prometheus.HistogramVec
	HTTPRequestSize     *prometheus.SummaryVec
	HTTPResponseSize    *prometheus.SummaryVec

	// 业务指标
	UserRegistrations *prometheus.CounterVec
	UserLogins        *prometheus.CounterVec
	UserLoginFailures *prometheus.CounterVec

	// 数据库指标
	DBConnections      prometheus.Gauge
	DBQueriesTotal     *prometheus.CounterVec
	DBQueryDuration    *prometheus.HistogramVec
	DBConnectionErrors *prometheus.CounterVec

	// Redis 指标
	RedisOperationsTotal   *prometheus.CounterVec
	RedisOperationDuration *prometheus.HistogramVec
	RedisConnectionErrors  *prometheus.CounterVec

	// RBAC 权限检查指标
	RBACPermissionChecks *prometheus.CounterVec
	RBACCacheHits        *prometheus.CounterVec

	// 限流指标
	RateLimitHits *prometheus.CounterVec
}

// NewMetrics 创建并注册所有指标
func NewMetrics(namespace string) *Metrics {
	m := &Metrics{
		// HTTP 请求指标
		HTTPRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "http_requests_total",
				Help:      "Total number of HTTP requests",
			},
			[]string{"service", "method", "path", "status"},
		),

		HTTPRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "http_request_duration_seconds",
				Help:      "HTTP request latencies in seconds",
				Buckets:   prometheus.DefBuckets, // 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10
			},
			[]string{"service", "method", "path"},
		),

		HTTPRequestSize: promauto.NewSummaryVec(
			prometheus.SummaryOpts{
				Namespace: namespace,
				Name:      "http_request_size_bytes",
				Help:      "HTTP request sizes in bytes",
			},
			[]string{"service", "method", "path"},
		),

		HTTPResponseSize: promauto.NewSummaryVec(
			prometheus.SummaryOpts{
				Namespace: namespace,
				Name:      "http_response_size_bytes",
				Help:      "HTTP response sizes in bytes",
			},
			[]string{"service", "method", "path"},
		),

		// 业务指标
		UserRegistrations: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "user_registrations_total",
				Help:      "Total number of user registrations",
			},
			[]string{"service", "status"},
		),

		UserLogins: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "user_logins_total",
				Help:      "Total number of successful user logins",
			},
			[]string{"service"},
		),

		UserLoginFailures: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "user_login_failures_total",
				Help:      "Total number of failed user login attempts",
			},
			[]string{"service", "reason"},
		),

		// 数据库指标
		DBConnections: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "db_connections",
				Help:      "Current number of database connections",
			},
		),

		DBQueriesTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "db_queries_total",
				Help:      "Total number of database queries",
			},
			[]string{"service", "operation", "table"},
		),

		DBQueryDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "db_query_duration_seconds",
				Help:      "Database query latencies in seconds",
				Buckets:   []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1},
			},
			[]string{"service", "operation", "table"},
		),

		DBConnectionErrors: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "db_connection_errors_total",
				Help:      "Total number of database connection errors",
			},
			[]string{"service"},
		),

		// Redis 指标
		RedisOperationsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "redis_operations_total",
				Help:      "Total number of Redis operations",
			},
			[]string{"service", "operation", "status"},
		),

		RedisOperationDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "redis_operation_duration_seconds",
				Help:      "Redis operation latencies in seconds",
				Buckets:   []float64{0.0001, 0.0005, 0.001, 0.005, 0.01, 0.025, 0.05, 0.1},
			},
			[]string{"service", "operation"},
		),

		RedisConnectionErrors: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "redis_connection_errors_total",
				Help:      "Total number of Redis connection errors",
			},
			[]string{"service"},
		),

		// RBAC 权限检查指标
		RBACPermissionChecks: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "rbac_permission_checks_total",
				Help:      "Total number of RBAC permission checks",
			},
			[]string{"service", "permission", "result"},
		),

		RBACCacheHits: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "rbac_cache_hits_total",
				Help:      "Total number of RBAC cache hits/misses",
			},
			[]string{"service", "cache_type", "result"},
		),

		// 限流指标
		RateLimitHits: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "rate_limit_hits_total",
				Help:      "Total number of rate limit hits",
			},
			[]string{"service", "limit_type", "identifier"},
		),
	}

	return m
}
