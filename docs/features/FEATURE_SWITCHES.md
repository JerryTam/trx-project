# 功能开关配置指南

## 📋 概述

项目中的**限流**和**链路追踪**功能都支持通过配置文件的 `enabled` 开关来启用或禁用。

## 🎛️ 配置位置

所有功能开关都在配置文件 `config/config.yaml` 中：

```yaml
# 限流配置（可选功能）
rate_limit:
  enabled: true           # ⭐ 总开关：启用/禁用限流功能
  global_rate: "1000-S"
  ip_rate: "100-M"
  user_rate: "1000-M"

# OpenTelemetry 链路追踪配置（可选功能）
tracing:
  enabled: true           # ⭐ 总开关：启用/禁用链路追踪
  service_name: "trx-project"
  service_version: "1.0.0"
  jaeger_endpoint: "localhost:4318"
```

## 🚦 限流功能 (Rate Limiting)

### 配置开关

```yaml
rate_limit:
  enabled: false  # ⭐ 设置为 false 禁用所有限流功能
```

### 功能说明

| 开关状态 | 行为 |
|---------|------|
| `enabled: true` | 启用全局、IP、用户三层限流保护 |
| `enabled: false` | 完全禁用限流，所有限流中间件不加载 |

### 启动日志

```bash
# 启用时
✅ 限流功能已启用
   - global_rate: 1000-S
   - ip_rate: 100-M
   - user_rate: 1000-M

# 禁用时
⚠️  限流功能已禁用
```

### 使用场景

#### ✅ 建议启用（enabled: true）

- **生产环境** - 防止系统过载和恶意攻击
- **测试环境** - 验证限流功能正常工作
- **预发布环境** - 模拟生产环境负载

#### ⚠️ 可以禁用（enabled: false）

- **开发环境** - 避免频繁触发限流影响开发调试
- **性能测试** - 测试系统的真实处理能力
- **单元测试** - 简化测试流程

### 性能影响

- **启用时**: < 1ms 延迟，< 1% CPU 开销
- **禁用时**: 无任何性能开销

## 🔍 链路追踪功能 (Tracing)

### 配置开关

```yaml
tracing:
  enabled: false  # ⭐ 设置为 false 禁用链路追踪
```

### 功能说明

| 开关状态 | 行为 |
|---------|------|
| `enabled: true` | 启用 OpenTelemetry 追踪，记录所有请求链路 |
| `enabled: false` | 完全禁用追踪，不初始化 Tracer，不应用中间件 |

### 启动日志

```bash
# 启用时
✅ OpenTelemetry 追踪已启用
   - service: trx-project-frontend
   - endpoint: localhost:4318

# 禁用时
⚠️  OpenTelemetry 追踪已禁用
   提示：如需启用，请在配置文件中设置 tracing.enabled = true
```

### 使用场景

#### ✅ 建议启用（enabled: true）

- **生产环境** - 便于问题排查和性能分析
- **测试环境** - 验证追踪功能正常工作
- **压力测试** - 分析性能瓶颈
- **故障排查** - 追踪请求完整链路

#### ⚠️ 可以禁用（enabled: false）

- **开发环境** - 减少资源占用，加快启动速度
- **本地调试** - 简化日志输出
- **轻量级部署** - 不需要 Jaeger 服务

### 性能影响

- **启用时**: < 1ms 延迟，< 2% CPU 开销，< 10MB 内存
- **禁用时**: 无任何性能开销

## 📝 环境配置建议

### 开发环境 (config.dev.yaml)

```yaml
rate_limit:
  enabled: false  # 🔸 建议禁用，方便开发调试

tracing:
  enabled: false  # 🔸 建议禁用，减少资源占用
```

### 测试环境 (config.test.yaml)

```yaml
rate_limit:
  enabled: true   # ✅ 建议启用，验证限流功能

tracing:
  enabled: true   # ✅ 建议启用，验证追踪功能
```

### 生产环境 (config.prod.yaml)

```yaml
rate_limit:
  enabled: true   # ⭐ 强烈建议启用，保护系统安全

tracing:
  enabled: true   # ⭐ 强烈建议启用，便于问题排查
```

## 🔄 切换功能

### 方法 1: 修改配置文件

```bash
# 编辑配置文件
vim config/config.yaml

# 修改对应的 enabled 值
rate_limit:
  enabled: false  # 改为 false

tracing:
  enabled: false  # 改为 false

# 重启服务
make run-frontend
make run-backend
```

### 方法 2: 使用环境特定配置

```bash
# 使用开发环境配置（限流和追踪都禁用）
export GO_ENV=dev
make run-frontend

# 使用测试环境配置（限流和追踪都启用）
export GO_ENV=test
make run-frontend

# 使用生产环境配置（限流和追踪都启用）
export GO_ENV=prod
make run-frontend
```

## 🧪 验证配置

### 检查限流是否生效

```bash
# 限流启用时，快速请求会触发 429 错误
for i in {1..200}; do
  curl http://localhost:8080/health
done

# 如果看到 "429 Too Many Requests"，说明限流已启用
# 如果所有请求都返回 200，说明限流已禁用
```

### 检查追踪是否生效

```bash
# 1. 查看启动日志
./bin/frontend | grep "追踪"

# 输出示例：
# ✅ OpenTelemetry 追踪已启用 (如果启用)
# ⚠️  OpenTelemetry 追踪已禁用 (如果禁用)

# 2. 检查 Jaeger UI
# 访问 http://localhost:16686
# 如果能看到追踪数据，说明追踪已启用
# 如果没有数据，说明追踪已禁用
```

## 💡 最佳实践

### 1. 分环境配置

创建不同的配置文件：

```bash
config/
├── config.yaml       # 默认配置
├── config.dev.yaml   # 开发环境（功能禁用）
├── config.test.yaml  # 测试环境（功能启用）
└── config.prod.yaml  # 生产环境（功能启用）
```

### 2. 渐进式启用

新环境部署时，建议按以下顺序启用功能：

```
1. 基础功能 ✅
   ↓
2. 链路追踪 ✅ (影响最小)
   ↓
3. 限流功能 ✅ (需要配置合理的阈值)
```

### 3. 监控告警

启用功能后，设置对应的监控：

```yaml
# Prometheus 告警规则
groups:
  - name: trx-project
    rules:
      # 限流触发率告警
      - alert: HighRateLimitTriggers
        expr: rate(trx_rate_limit_triggers_total[5m]) > 10
        annotations:
          summary: "限流触发频率过高"

      # 追踪采样率告警
      - alert: TracingDisabled
        expr: up{job="trx-project"} == 1 and tracing_enabled == 0
        annotations:
          summary: "链路追踪未启用"
```

## 🔧 故障排查

### 限流未生效

```bash
# 1. 检查配置
grep "enabled" config/config.yaml

# 2. 检查 Redis 连接
redis-cli ping

# 3. 查看日志
tail -f logs/app.log | grep "限流"
```

### 追踪未生效

```bash
# 1. 检查配置
grep "enabled" config/config.yaml

# 2. 检查 Jaeger 连接
curl http://localhost:4318

# 3. 查看日志
tail -f logs/app.log | grep "追踪"
```

## 📊 性能对比

| 功能 | 启用时延迟 | 禁用时延迟 | 内存开销 |
|------|-----------|-----------|---------|
| **限流** | +0.5ms | 0ms | +5MB |
| **追踪** | +1.0ms | 0ms | +10MB |
| **两者都启用** | +1.5ms | 0ms | +15MB |

## 🎯 总结

### 🟢 推荐配置

| 环境 | 限流 | 追踪 | 原因 |
|------|------|------|------|
| **开发** | ❌ | ❌ | 方便调试，减少干扰 |
| **测试** | ✅ | ✅ | 验证功能，发现问题 |
| **生产** | ✅ | ✅ | 保护系统，排查问题 |

### 🔑 关键点

1. **限流和追踪都不是必需功能** - 可根据实际需求开关
2. **配置非常灵活** - 通过 `enabled` 开关即可控制
3. **性能影响很小** - 即使启用也不会明显影响系统
4. **推荐生产启用** - 便于保护系统和排查问题
5. **开发环境可禁用** - 提高开发效率

---

**配置文件位置**: 
- 默认配置: `config/config.yaml`
- 开发配置: `config/config.dev.yaml`
- 测试配置: `config/config.test.yaml`
- 生产配置: `config/config.prod.yaml`

**相关文档**:
- [限流使用指南](RATE_LIMIT_GUIDE.md)
- [OpenTelemetry 追踪指南](OPENTELEMETRY_TRACING_GUIDE.md)
- [配置管理文档](../architecture/CONFIG_GUIDE.md)

