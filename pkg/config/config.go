package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server    ServerConfig    `yaml:"server"`
	Database  DatabaseConfig  `yaml:"database"`
	Redis     RedisConfig     `yaml:"redis"`
	Kafka     KafkaConfig     `yaml:"kafka"`
	Logger    LoggerConfig    `yaml:"logger"`
	JWT       JWTConfig       `yaml:"jwt"`
	RateLimit RateLimitConfig `yaml:"rate_limit"`
	Tracing   TracingConfig   `yaml:"tracing"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"`
	Env  string `yaml:"env"` // 运行环境: dev, test, prod
}

type DatabaseConfig struct {
	MySQL MySQLConfig `yaml:"mysql"`
}

type MySQLConfig struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	Database     string `yaml:"database"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
	MaxOpenConns int    `yaml:"max_open_conns"`
	MaxLifetime  int    `yaml:"max_lifetime"`
	LogLevel     string `yaml:"log_level"`
}

type RedisConfig struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	Password     string `yaml:"password"`
	DB           int    `yaml:"db"`
	PoolSize     int    `yaml:"pool_size"`
	MinIdleConns int    `yaml:"min_idle_conns"`
}

type KafkaConfig struct {
	Brokers []string `yaml:"brokers"`
	GroupID string   `yaml:"group_id"`
	Topics  []string `yaml:"topics"`
}

type LoggerConfig struct {
	Level            string   `yaml:"level"`
	Encoding         string   `yaml:"encoding"`
	OutputPaths      []string `yaml:"output_paths"`
	ErrorOutputPaths []string `yaml:"error_output_paths"`
}

type JWTConfig struct {
	Secret           string `yaml:"secret"`
	Issuer           string `yaml:"issuer"`
	ExpireHours      int    `yaml:"expire_hours"`       // 用户 Token 过期时间（小时）
	AdminExpireHours int    `yaml:"admin_expire_hours"` // 管理员 Token 过期时间（小时）
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	Enabled    bool   `yaml:"enabled"`     // 是否启用限流
	GlobalRate string `yaml:"global_rate"` // 全局限流：例如 "1000-S" (每秒1000个请求)
	IPRate     string `yaml:"ip_rate"`     // IP限流：例如 "100-M" (每分钟100个请求)
	UserRate   string `yaml:"user_rate"`   // 用户限流：例如 "1000-M" (每分钟1000个请求)
}

// TracingConfig OpenTelemetry 追踪配置
type TracingConfig struct {
	Enabled        bool   `yaml:"enabled"`         // 是否启用追踪
	ServiceName    string `yaml:"service_name"`    // 服务名称
	ServiceVersion string `yaml:"service_version"` // 服务版本
	JaegerEndpoint string `yaml:"jaeger_endpoint"` // Jaeger OTLP HTTP 端点
}

// Load 根据环境加载配置文件
// 优先级: 环境变量 GO_ENV > 命令行参数 > 默认值 (dev)
// 配置文件命名: config.{env}.yaml
func Load(path string) (*Config, error) {
	// 获取环境配置
	env := GetEnv()

	// 如果 path 为空，使用默认路径
	if path == "" {
		path = "config/config.yaml"
	}

	// 构建环境特定的配置文件路径
	dir := filepath.Dir(path)
	ext := filepath.Ext(path)
	base := filepath.Base(path)
	nameWithoutExt := base[:len(base)-len(ext)]

	// 尝试加载环境特定配置文件，如 config.dev.yaml
	envConfigPath := filepath.Join(dir, fmt.Sprintf("%s.%s%s", nameWithoutExt, env, ext))

	// 如果环境配置文件存在，使用它；否则使用默认配置文件
	configPath := path
	if _, err := os.Stat(envConfigPath); err == nil {
		configPath = envConfigPath
		fmt.Printf("✅ 加载环境配置: %s (环境: %s)\n", configPath, env)
	} else {
		fmt.Printf("⚠️  环境配置文件不存在 (%s)，使用默认配置: %s\n", envConfigPath, path)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 设置当前环境
	cfg.Server.Env = env

	return &cfg, nil
}

// GetEnv 获取当前运行环境
// 优先级: 1. GO_ENV 环境变量 2. 默认 dev
func GetEnv() string {
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "dev"
	}
	return env
}

// GetDSN 返回 MySQL DSN 字符串
func (m *MySQLConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		m.Username, m.Password, m.Host, m.Port, m.Database)
}

// GetAddress 返回 Redis 地址
func (r *RedisConfig) GetAddress() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}
