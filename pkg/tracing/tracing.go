package tracing

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.uber.org/zap"
)

// Config OpenTelemetry 配置
type Config struct {
	ServiceName    string // 服务名称
	ServiceVersion string // 服务版本
	Environment    string // 环境 (dev/test/prod)
	JaegerEndpoint string // Jaeger OTLP HTTP 端点
	Enabled        bool   // 是否启用追踪
}

// InitTracer 初始化 OpenTelemetry 追踪器
func InitTracer(cfg *Config, logger *zap.Logger) (func(context.Context) error, error) {
	if !cfg.Enabled {
		logger.Info("OpenTelemetry tracing is disabled")
		// 返回一个空的清理函数
		return func(ctx context.Context) error { return nil }, nil
	}

	logger.Info("Initializing OpenTelemetry tracing",
		zap.String("service", cfg.ServiceName),
		zap.String("version", cfg.ServiceVersion),
		zap.String("environment", cfg.Environment),
		zap.String("endpoint", cfg.JaegerEndpoint))

	// 创建资源
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
			semconv.ServiceVersion(cfg.ServiceVersion),
			semconv.DeploymentEnvironment(cfg.Environment),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// 创建 OTLP HTTP 导出器
	exporter, err := otlptracehttp.New(context.Background(),
		otlptracehttp.WithEndpoint(cfg.JaegerEndpoint),
		otlptracehttp.WithInsecure(), // 开发环境使用 HTTP，生产环境应使用 HTTPS
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP trace exporter: %w", err)
	}

	// 创建追踪提供者
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()), // 开发环境：全部采样；生产环境：可以使用 ParentBased 或概率采样
	)

	// 设置全局追踪提供者
	otel.SetTracerProvider(tp)

	// 设置全局传播器（用于跨服务追踪）
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	logger.Info("OpenTelemetry tracing initialized successfully")

	// 返回清理函数
	cleanup := func(ctx context.Context) error {
		logger.Info("Shutting down OpenTelemetry tracer...")
		if err := tp.Shutdown(ctx); err != nil {
			logger.Error("Failed to shutdown tracer provider", zap.Error(err))
			return err
		}
		logger.Info("OpenTelemetry tracer shut down successfully")
		return nil
	}

	return cleanup, nil
}
