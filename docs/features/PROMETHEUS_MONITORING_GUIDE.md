# Prometheus 监控集成指南

## 📋 概述

本项目已集成 Prometheus + Grafana 完整监控解决方案，提供实时的系统性能监控和可视化。

---

## 🏗️ 架构

```
应用服务 (Frontend/Backend)
    ↓ 暴露 /metrics 端点
Prometheus
    ↓ 采集指标数据
Grafana
    ↓ 可视化展示
用户 (浏览器访问仪表板)
```

---

## 🚀 快速开始

### 1. 启动监控服务

```bash
# 启动所有服务（包括 Prometheus 和 Grafana）
make docker-up

# 等待服务启动（约 30 秒）
```

### 2. 启动应用服务

```bash
# 启动前台服务（端口 8080）
make dev-frontend

# 启动后台服务（端口 8081）
make dev-backend
```

### 3. 访问监控界面

**Prometheus UI:**
- URL: http://localhost:9090
- 查看原始指标和执行 PromQL 查询

**Grafana 仪表板:**
- URL: http://localhost:3000
- 默认账号: `admin`
- 默认密码: `admin`

首次登录后会提示修改密码，可以选择跳过。

### 4. 查看仪表板

在 Grafana 中：
1. 点击左侧菜单 "Dashboards"
2. 选择 "TRX Project - 服务监控"
3. 查看实时监控数据

---

## 📊 监控指标

### HTTP 请求指标

| 指标名称 | 类型 | 说明 |
|---------|------|------|
| `trx_http_requests_total` | Counter | HTTP 请求总数 |
| `trx_http_request_duration_seconds` | Histogram | HTTP 请求延迟 |
| `trx_http_request_size_bytes` | Summary | HTTP 请求大小 |
| `trx_http_response_size_bytes` | Summary | HTTP 响应大小 |

**标签 (Labels):**
- `service`: 服务名称 (frontend/backend)
- `method`: HTTP 方法 (GET/POST/PUT/DELETE)
- `path`: 请求路径
- `status`: HTTP 状态码

**示例查询:**
```promql
# 前台服务的 QPS
rate(trx_http_requests_total{service="frontend"}[5m])

# P95 延迟
histogram_quantile(0.95, rate(trx_http_request_duration_seconds_bucket[5m]))

# 错误率
rate(trx_http_requests_total{status=~"5.."}[5m]) / rate(trx_http_requests_total[5m])
```

### 业务指标

| 指标名称 | 类型 | 说明 |
|---------|------|------|
| `trx_user_registrations_total` | Counter | 用户注册总数 |
| `trx_user_logins_total` | Counter | 用户登录总数 |
| `trx_user_login_failures_total` | Counter | 登录失败总数 |

**标签:**
- `service`: 服务名称
- `status`: 操作状态 (success/failure)
- `reason`: 失败原因

### 数据库指标

| 指标名称 | 类型 | 说明 |
|---------|------|------|
| `trx_db_connections` | Gauge | 当前数据库连接数 |
| `trx_db_queries_total` | Counter | 数据库查询总数 |
| `trx_db_query_duration_seconds` | Histogram | 数据库查询延迟 |
| `trx_db_connection_errors_total` | Counter | 数据库连接错误数 |

**标签:**
- `service`: 服务名称
- `operation`: 操作类型 (SELECT/INSERT/UPDATE/DELETE)
- `table`: 表名

### Redis 指标

| 指标名称 | 类型 | 说明 |
|---------|------|------|
| `trx_redis_operations_total` | Counter | Redis 操作总数 |
| `trx_redis_operation_duration_seconds` | Histogram | Redis 操作延迟 |
| `trx_redis_connection_errors_total` | Counter | Redis 连接错误数 |

**标签:**
- `service`: 服务名称
- `operation`: 操作类型 (GET/SET/DEL 等)
- `status`: 操作状态 (success/error)

### RBAC 权限指标

| 指标名称 | 类型 | 说明 |
|---------|------|------|
| `trx_rbac_permission_checks_total` | Counter | 权限检查总数 |
| `trx_rbac_cache_hits_total` | Counter | RBAC 缓存命中数 |

**标签:**
- `service`: 服务名称
- `permission`: 权限代码
- `result`: 检查结果 (allowed/denied)
- `cache_type`: 缓存类型
- `result`: 缓存结果 (hit/miss)

### 限流指标

| 指标名称 | 类型 | 说明 |
|---------|------|------|
| `trx_rate_limit_hits_total` | Counter | 限流触发总数 |

**标签:**
- `service`: 服务名称
- `limit_type`: 限流类型 (global/ip/user)
- `identifier`: 标识符 (IP 地址或用户 ID)

---

## 📈 常用 PromQL 查询

### 性能监控

```promql
# QPS（每秒请求数）
rate(trx_http_requests_total[5m])

# P50 延迟
histogram_quantile(0.50, rate(trx_http_request_duration_seconds_bucket[5m]))

# P95 延迟
histogram_quantile(0.95, rate(trx_http_request_duration_seconds_bucket[5m]))

# P99 延迟
histogram_quantile(0.99, rate(trx_http_request_duration_seconds_bucket[5m]))

# 平均响应时间
rate(trx_http_request_duration_seconds_sum[5m]) / rate(trx_http_request_duration_seconds_count[5m])
```

### 错误监控

```promql
# 4xx 错误率
sum(rate(trx_http_requests_total{status=~"4.."}[5m])) by (service) / sum(rate(trx_http_requests_total[5m])) by (service)

# 5xx 错误率
sum(rate(trx_http_requests_total{status=~"5.."}[5m])) by (service) / sum(rate(trx_http_requests_total[5m])) by (service)

# 登录失败率
rate(trx_user_login_failures_total[5m]) / rate(trx_user_logins_total[5m])
```

### 业务监控

```promql
# 注册用户增长速率
rate(trx_user_registrations_total[5m])

# 活跃用户数（最近 5 分钟登录）
sum(increase(trx_user_logins_total[5m]))

# 限流触发频率
rate(trx_rate_limit_hits_total[5m])
```

### 缓存监控

```promql
# RBAC 缓存命中率
sum(rate(trx_rbac_cache_hits_total{result="hit"}[5m])) / sum(rate(trx_rbac_cache_hits_total[5m]))

# Redis 操作成功率
sum(rate(trx_redis_operations_total{status="success"}[5m])) / sum(rate(trx_redis_operations_total[5m]))
```

---

## 🎨 Grafana 仪表板

### 默认仪表板功能

**TRX Project - 服务监控** 包含以下面板：

1. **HTTP 请求速率 (QPS)** - 实时 QPS 图表
2. **HTTP 请求延迟 (P95/P99)** - 延迟百分位数
3. **总请求速率** - 单值面板显示总 QPS
4. **用户注册总数** - 业务指标
5. **用户登录总数** - 业务指标
6. **登录失败总数** - 错误监控
7. **HTTP 状态码分布** - 饼图
8. **限流触发速率** - 限流监控

### 自定义仪表板

#### 创建新面板

1. 进入仪表板
2. 点击右上角 "Add panel"
3. 选择 "Add a new panel"
4. 在查询编辑器中输入 PromQL
5. 选择可视化类型（时间序列、统计、饼图等）
6. 配置面板选项
7. 点击 "Apply" 保存

#### 导出/导入仪表板

**导出:**
1. 打开仪表板
2. 点击顶部设置图标 ⚙️
3. 选择 "JSON Model"
4. 复制 JSON 内容

**导入:**
1. 点击左侧 "+" → "Import"
2. 粘贴 JSON 或上传文件
3. 点击 "Load"

---

## ⚙️ 配置说明

### Prometheus 配置

文件: `config/prometheus.yml`

```yaml
scrape_configs:
  - job_name: 'frontend'
    metrics_path: '/metrics'
    static_configs:
      - targets: ['host.docker.internal:8080']
        labels:
          service: 'frontend'
```

**关键配置:**
- `scrape_interval`: 采集间隔（默认 15s）
- `evaluation_interval`: 规则评估间隔（默认 15s）
- `retention.time`: 数据保留时间（默认 15天）

### Grafana 数据源配置

文件: `config/grafana/provisioning/datasources/prometheus.yml`

```yaml
datasources:
  - name: Prometheus
    type: prometheus
    url: http://prometheus:9090
    isDefault: true
```

---

## 🔧 故障排查

### 1. 指标数据不显示

**症状**: Grafana 仪表板显示"No Data"

**解决方案:**

```bash
# 1. 检查应用服务是否启动
curl http://localhost:8080/health
curl http://localhost:8081/health

# 2. 检查 metrics 端点
curl http://localhost:8080/metrics
curl http://localhost:8081/metrics

# 3. 检查 Prometheus 是否能访问服务
# 访问 http://localhost:9090/targets
# 确认 frontend 和 backend 目标状态为 UP

# 4. 检查 Docker 网络
docker ps | grep prometheus
docker ps | grep grafana
```

### 2. Prometheus 无法采集数据

**症状**: Prometheus targets 页面显示服务为 DOWN

**原因**: Docker 容器无法访问宿主机服务

**解决方案:**

**Windows/Mac:**
```yaml
# config/prometheus.yml 已配置
targets: ['host.docker.internal:8080']
```

**Linux:**
```yaml
# 修改 config/prometheus.yml
targets: ['172.17.0.1:8080']
# 或使用 host 网络模式
```

### 3. Grafana 无法连接 Prometheus

**症状**: Grafana 数据源测试失败

**解决方案:**

```bash
# 1. 检查 Prometheus 容器
docker logs trx-prometheus

# 2. 检查网络连通性
docker exec trx-grafana ping prometheus

# 3. 重启服务
docker-compose restart prometheus grafana
```

### 4. 仪表板显示异常

**症状**: 图表显示不正确或缺失数据

**解决方案:**

```bash
# 1. 检查时间范围（右上角）
# 确保选择了有数据的时间段

# 2. 验证 PromQL 查询
# 在 Prometheus UI 中测试查询

# 3. 检查数据源配置
# Dashboard settings → Data source

# 4. 重新加载仪表板
# 点击刷新按钮或按 Ctrl+R
```

---

## 📊 告警配置（高级）

### 创建告警规则

创建文件 `config/prometheus/rules/alerts.yml`:

```yaml
groups:
  - name: trx_alerts
    interval: 30s
    rules:
      # 高错误率告警
      - alert: HighErrorRate
        expr: |
          sum(rate(trx_http_requests_total{status=~"5.."}[5m])) by (service) 
          / 
          sum(rate(trx_http_requests_total[5m])) by (service) 
          > 0.05
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High 5xx error rate detected"
          description: "Service {{ $labels.service }} has error rate of {{ $value }}"

      # 高延迟告警
      - alert: HighLatency
        expr: |
          histogram_quantile(0.95, 
            rate(trx_http_request_duration_seconds_bucket[5m])
          ) > 1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High latency detected"
          description: "P95 latency is {{ $value }}s"

      # 登录失败率告警
      - alert: HighLoginFailureRate
        expr: |
          rate(trx_user_login_failures_total[5m]) 
          / 
          rate(trx_user_logins_total[5m]) 
          > 0.3
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High login failure rate"
          description: "Login failure rate is {{ $value }}"
```

更新 `config/prometheus.yml`:

```yaml
rule_files:
  - "rules/*.yml"

alerting:
  alertmanagers:
    - static_configs:
        - targets:
          - alertmanager:9093
```

---

## 🎯 最佳实践

### 1. 指标命名

- ✅ 使用清晰的命名空间: `trx_*`
- ✅ 使用标准后缀: `_total`, `_seconds`, `_bytes`
- ✅ 使用有意义的标签

### 2. 查询优化

- ✅ 使用适当的时间范围: `[5m]` 而不是 `[1h]`
- ✅ 避免过多的标签基数
- ✅ 使用 `rate()` 而不是 `increase()` 计算速率

### 3. 仪表板设计

- ✅ 按功能分组面板
- ✅ 使用合适的可视化类型
- ✅ 添加清晰的标题和说明
- ✅ 设置合理的刷新间隔（10s-30s）

### 4. 性能考虑

- ✅ 控制指标数量
- ✅ 定期清理过期数据
- ✅ 监控 Prometheus 自身性能

---

## 📚 参考资源

- [Prometheus 官方文档](https://prometheus.io/docs/)
- [Grafana 官方文档](https://grafana.com/docs/)
- [PromQL 查询语法](https://prometheus.io/docs/prometheus/latest/querying/basics/)
- [Grafana 仪表板设计指南](https://grafana.com/docs/grafana/latest/dashboards/)

---

## 📝 总结

### 已实现功能

✅ **完整的监控指标**
- HTTP 请求监控
- 业务指标追踪
- 数据库性能监控
- Redis 操作监控
- RBAC 权限监控
- 限流监控

✅ **自动化部署**
- Docker Compose 一键启动
- 自动配置数据源
- 预置仪表板

✅ **实时可视化**
- Grafana 仪表板
- 多维度图表
- 实时更新

✅ **生产就绪**
- 性能优化
- 错误监控
- 故障排查

### 下一步

1. 配置告警通知（邮件/钉钉/Slack）
2. 添加更多业务指标
3. 创建更多专用仪表板
4. 设置数据备份策略

---

**维护者**: TRX Project Team
**更新时间**: 2024-11
**版本**: 1.0.0

