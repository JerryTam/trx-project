# 请求 ID 追踪功能使用指南

## 📋 目录

- [功能概述](#功能概述)
- [实现原理](#实现原理)
- [使用方法](#使用方法)
- [日志追踪](#日志追踪)
- [测试验证](#测试验证)
- [最佳实践](#最佳实践)
- [常见问题](#常见问题)
- [集成示例](#集成示例)

---

## 🎯 功能概述

请求 ID 追踪（Request ID Tracing）是一个用于追踪完整请求链路的功能，为每个请求分配一个唯一的标识符，便于调试、日志分析和问题排查。

### 主要特性

✅ **自动生成唯一 ID**
- 使用 UUID v4 生成全局唯一标识符
- 格式：`xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx`

✅ **客户端指定 ID**
- 支持客户端通过 `X-Request-ID` 头传递请求 ID
- 适用于跨服务调用和链路追踪

✅ **完整的日志集成**
- 所有日志自动包含请求 ID
- 方便通过请求 ID 过滤和追踪

✅ **响应头返回**
- 每个响应都包含 `X-Request-ID` 头
- 客户端可以获取并追踪请求

✅ **零侵入式**
- 基于中间件实现
- 不需要修改业务代码

---

## 🔧 实现原理

### 架构设计

```
客户端请求
    ↓
[1] 检查 X-Request-ID 请求头
    ├─ 存在 → 使用客户端提供的 ID
    └─ 不存在 → 生成新的 UUID
    ↓
[2] 存储到 Gin Context
    └─ key: "request_id"
    ↓
[3] 添加到响应头
    └─ X-Request-ID: <request-id>
    ↓
[4] 记录到日志
    └─ 所有日志字段自动包含 request_id
    ↓
[5] 业务处理
    └─ 可通过 middleware.GetRequestID(c) 获取
    ↓
返回响应（包含 X-Request-ID 头）
```

### 中间件执行顺序

```
Recovery 
    ↓
RequestID  ← 生成或获取请求 ID
    ↓
Logger     ← 日志包含请求 ID
    ↓
CORS
    ↓
限流
    ↓
认证
    ↓
业务逻辑
```

**重要**：RequestID 中间件必须在 Logger 中间件之前，这样日志才能包含请求 ID。

---

## 🚀 使用方法

### 1. 服务器自动生成

**最简单的方式**：不需要做任何特殊处理，服务器会自动为每个请求生成唯一 ID。

```bash
# 发送普通请求
curl http://localhost:8080/api/v1/public/login \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"test"}'
```

**响应头**：
```http
HTTP/1.1 200 OK
X-Request-ID: 550e8400-e29b-41d4-a716-446655440000
Content-Type: application/json
...
```

### 2. 客户端指定 ID

**用于请求链追踪**：客户端可以生成并传递请求 ID，在整个调用链中保持一致。

```bash
# 指定请求 ID
curl http://localhost:8080/api/v1/public/login \
  -H "X-Request-ID: my-custom-request-id-12345" \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"test"}'
```

**响应头**：
```http
HTTP/1.1 200 OK
X-Request-ID: my-custom-request-id-12345  ← 返回相同的 ID
Content-Type: application/json
...
```

### 3. 在代码中获取请求 ID

如果你需要在业务代码中获取当前请求的 ID：

```go
import "trx-project/internal/api/middleware"

func MyHandler(c *gin.Context) {
    // 获取请求 ID
    requestID := middleware.GetRequestID(c)
    
    // 使用请求 ID
    logger.Info("Processing request",
        zap.String("request_id", requestID),
        zap.String("user", "john"))
    
    // 或者创建带请求 ID 的 logger
    reqLogger := middleware.RequestIDToLogger(c, logger)
    reqLogger.Info("This log will include request_id automatically")
}
```

---

## 📝 日志追踪

### 日志格式

启用请求 ID 后，所有日志都会自动包含 `request_id` 字段：

**JSON 格式**：
```json
{
  "level": "info",
  "ts": "2024-11-11T10:30:15.123Z",
  "msg": "HTTP Request",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "method": "POST",
  "path": "/api/v1/public/login",
  "status": 200,
  "latency": "15.234ms",
  "client_ip": "192.168.1.100"
}
```

**Console 格式**：
```
2024-11-11T10:30:15.123Z  INFO  HTTP Request
  request_id=550e8400-e29b-41d4-a716-446655440000
  method=POST
  path=/api/v1/public/login
  status=200
  latency=15.234ms
  client_ip=192.168.1.100
```

### 追踪完整请求

#### 1. 实时追踪

```bash
# 追踪所有请求
tail -f logs/app.log | grep request_id

# 追踪特定请求 ID
tail -f logs/app.log | grep "550e8400-e29b-41d4-a716-446655440000"
```

#### 2. 历史查询

```bash
# 查找特定请求的所有日志
grep "550e8400-e29b-41d4-a716-446655440000" logs/app.log

# 只显示错误日志
grep "550e8400-e29b-41d4-a716-446655440000" logs/error.log

# 使用 jq 解析 JSON 日志
cat logs/app.log | jq 'select(.request_id=="550e8400-e29b-41d4-a716-446655440000")'
```

#### 3. 统计分析

```bash
# 统计每个请求 ID 的日志数量
grep -o '"request_id":"[^"]*"' logs/app.log | sort | uniq -c | sort -rn

# 找出日志最多的请求（可能是问题请求）
grep -o '"request_id":"[^"]*"' logs/app.log | sort | uniq -c | sort -rn | head -10
```

---

## 🧪 测试验证

### 方法 1: 使用测试脚本

```bash
# 运行完整的测试套件
./scripts/test_request_id.sh
```

测试内容：
- ✅ 服务器自动生成请求 ID
- ✅ 客户端指定请求 ID
- ✅ 请求链追踪
- ✅ 前后台服务支持
- ✅ 不同 API 端点
- ✅ 并发请求唯一性

### 方法 2: 手动测试

#### 测试自动生成

```bash
# 发送请求并查看响应头
curl -i http://localhost:8080/health

# 响应示例
HTTP/1.1 200 OK
X-Request-ID: 550e8400-e29b-41d4-a716-446655440000
Content-Type: application/json
...
```

#### 测试客户端指定

```bash
# 指定请求 ID
curl -i -H "X-Request-ID: test-12345" http://localhost:8080/health

# 验证返回的 ID 是否相同
HTTP/1.1 200 OK
X-Request-ID: test-12345
...
```

#### 测试请求链追踪

```bash
# 使用相同的请求 ID 发送多个请求
REQUEST_ID="trace-$(date +%s)"

# 请求 1
curl -H "X-Request-ID: ${REQUEST_ID}" \
  http://localhost:8080/api/v1/public/login \
  -d '{"username":"user1","password":"pass1"}'

# 请求 2
curl -H "X-Request-ID: ${REQUEST_ID}" \
  http://localhost:8080/api/v1/user/profile \
  -H "Authorization: Bearer <token>"

# 查看日志
grep "${REQUEST_ID}" logs/app.log
```

---

## 💡 最佳实践

### 1. 微服务链路追踪

在微服务架构中，使用请求 ID 追踪完整的调用链：

```
用户请求
    ↓ (request_id: abc-123)
前端服务
    ↓ (传递 X-Request-ID: abc-123)
后端服务
    ↓ (传递 X-Request-ID: abc-123)
数据库服务
    ↓ (传递 X-Request-ID: abc-123)
缓存服务
```

**实现示例**：

```go
// 调用其他服务时传递请求 ID
func CallOtherService(c *gin.Context, url string) error {
    requestID := middleware.GetRequestID(c)
    
    req, _ := http.NewRequest("GET", url, nil)
    req.Header.Set("X-Request-ID", requestID)  // 传递请求 ID
    
    resp, err := http.DefaultClient.Do(req)
    // ...
}
```

### 2. 前端集成

前端可以生成并维护请求 ID：

```javascript
// 生成请求 ID
function generateRequestId() {
    return 'web-' + Date.now() + '-' + Math.random().toString(36).substr(2, 9);
}

// 发送请求时携带请求 ID
async function apiCall(url, options = {}) {
    const requestId = generateRequestId();
    
    const response = await fetch(url, {
        ...options,
        headers: {
            ...options.headers,
            'X-Request-ID': requestId,
        },
    });
    
    // 从响应头获取请求 ID（验证）
    const returnedId = response.headers.get('X-Request-ID');
    console.log('Request ID:', returnedId);
    
    return response.json();
}
```

### 3. 错误追踪

当发生错误时，将请求 ID 返回给客户端：

```go
func ErrorHandler(c *gin.Context, err error) {
    requestID := middleware.GetRequestID(c)
    
    logger.Error("Request failed",
        zap.String("request_id", requestID),
        zap.Error(err))
    
    c.JSON(500, gin.H{
        "error": "Internal Server Error",
        "request_id": requestID,  // 返回给客户端
        "message": "Please contact support with this request ID",
    })
}
```

客户端可以将请求 ID 提供给客服或技术支持，方便快速定位问题。

### 4. 性能监控

结合请求 ID 进行性能分析：

```go
func PerformanceTracker(c *gin.Context) {
    requestID := middleware.GetRequestID(c)
    start := time.Now()
    
    // 处理请求
    c.Next()
    
    // 记录性能指标
    duration := time.Since(start)
    logger.Info("Request performance",
        zap.String("request_id", requestID),
        zap.Duration("duration", duration),
        zap.String("path", c.Request.URL.Path))
    
    // 如果请求太慢，记录告警
    if duration > 1*time.Second {
        logger.Warn("Slow request detected",
            zap.String("request_id", requestID),
            zap.Duration("duration", duration))
    }
}
```

### 5. 日志聚合

在使用日志聚合系统（如 ELK、Loki）时，请求 ID 是重要的索引字段：

```yaml
# Logstash 配置示例
filter {
  json {
    source => "message"
  }
  
  # 提取请求 ID 作为独立字段
  if [request_id] {
    mutate {
      add_field => {
        "trace_id" => "%{request_id}"
      }
    }
  }
}

# Kibana 查询
request_id: "550e8400-e29b-41d4-a716-446655440000"
```

---

## ❓ 常见问题

### Q1: 请求 ID 是否会重复？

**A**: 几乎不可能。UUID v4 的碰撞概率极低（约 1/10^36），在实际应用中可以认为是唯一的。

### Q2: 请求 ID 对性能有影响吗？

**A**: 影响极小：
- UUID 生成时间: < 1μs
- Context 存储/读取: < 0.1μs
- 总体性能损失: < 0.1%

### Q3: 如何在业务代码中使用请求 ID？

**A**: 有两种方式：

```go
// 方式 1: 直接获取
requestID := middleware.GetRequestID(c)
logger.Info("Processing", zap.String("request_id", requestID))

// 方式 2: 创建带请求 ID 的 logger
reqLogger := middleware.RequestIDToLogger(c, logger)
reqLogger.Info("Processing")  // 自动包含 request_id
```

### Q4: 可以自定义请求 ID 格式吗？

**A**: 可以。修改 `internal/api/middleware/request_id.go`：

```go
// 自定义格式
requestID := fmt.Sprintf("req-%s-%d", 
    uuid.New().String()[:8], 
    time.Now().Unix())
```

### Q5: 如何追踪跨服务的请求？

**A**: 在调用其他服务时传递 `X-Request-ID` 头：

```go
func CallExternalAPI(c *gin.Context) {
    requestID := middleware.GetRequestID(c)
    
    req, _ := http.NewRequest("GET", "http://other-service/api", nil)
    req.Header.Set("X-Request-ID", requestID)  // 传递
    
    resp, _ := http.DefaultClient.Do(req)
    // ...
}
```

### Q6: 请求 ID 会存储到数据库吗？

**A**: 默认不会。但你可以在业务逻辑中获取并存储：

```go
func CreateOrder(c *gin.Context, order *Order) {
    requestID := middleware.GetRequestID(c)
    order.RequestID = requestID  // 存储到数据库
    
    db.Create(order)
}
```

---

## 🔗 集成示例

### 示例 1: 用户注册流程追踪

```go
// Handler
func (h *UserHandler) Register(c *gin.Context) {
    // 获取带请求 ID 的 logger
    logger := middleware.RequestIDToLogger(c, h.logger)
    
    logger.Info("User registration started")
    
    var req RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        logger.Error("Invalid request", zap.Error(err))
        response.BadRequest(c, err.Error())
        return
    }
    
    logger.Info("Creating user", zap.String("username", req.Username))
    
    user, err := h.service.Register(c.Request.Context(), req.Username, req.Email, req.Password)
    if err != nil {
        logger.Error("Registration failed", zap.Error(err))
        response.InternalError(c, "Failed to register")
        return
    }
    
    logger.Info("User registered successfully", zap.Uint("user_id", user.ID))
    response.Success(c, user)
}
```

**日志输出**：
```
INFO  User registration started  request_id=abc-123
INFO  Creating user  request_id=abc-123 username=john
INFO  User registered successfully  request_id=abc-123 user_id=42
```

### 示例 2: API 调用链追踪

```go
// 服务 A 调用服务 B
func ServiceAHandler(c *gin.Context) {
    requestID := middleware.GetRequestID(c)
    logger := middleware.RequestIDToLogger(c, globalLogger)
    
    logger.Info("Service A: Processing request")
    
    // 调用服务 B，传递请求 ID
    resp, err := callServiceB(requestID)
    if err != nil {
        logger.Error("Service B call failed", zap.Error(err))
        return
    }
    
    logger.Info("Service A: Request completed")
    c.JSON(200, resp)
}

func callServiceB(requestID string) (interface{}, error) {
    req, _ := http.NewRequest("GET", "http://service-b/api", nil)
    req.Header.Set("X-Request-ID", requestID)
    
    resp, err := http.DefaultClient.Do(req)
    // ...
}
```

**完整链路日志**：
```
[Service A] INFO  Service A: Processing request  request_id=trace-456
[Service A] INFO  Calling Service B  request_id=trace-456
[Service B] INFO  Service B: Processing request  request_id=trace-456
[Service B] INFO  Service B: Request completed  request_id=trace-456
[Service A] INFO  Service A: Request completed  request_id=trace-456
```

### 示例 3: 错误上下文追踪

```go
func ComplexOperation(c *gin.Context) {
    requestID := middleware.GetRequestID(c)
    logger := middleware.RequestIDToLogger(c, globalLogger)
    
    // 步骤 1
    logger.Info("Step 1: Validating input")
    if err := validateInput(); err != nil {
        logger.Error("Validation failed", zap.Error(err))
        response.BadRequest(c, "Invalid input")
        return
    }
    
    // 步骤 2
    logger.Info("Step 2: Processing data")
    if err := processData(); err != nil {
        logger.Error("Processing failed", zap.Error(err))
        c.JSON(500, gin.H{
            "error": "Processing failed",
            "request_id": requestID,  // 返回给客户端
        })
        return
    }
    
    // 步骤 3
    logger.Info("Step 3: Saving result")
    if err := saveResult(); err != nil {
        logger.Error("Save failed", zap.Error(err))
        c.JSON(500, gin.H{
            "error": "Save failed",
            "request_id": requestID,
        })
        return
    }
    
    logger.Info("Operation completed successfully")
    response.Success(c, "OK")
}
```

---

## 📚 相关文档

- [README.md](README.md) - 项目主文档
- [OPTIMIZATION_RECOMMENDATIONS.md](OPTIMIZATION_RECOMMENDATIONS.md) - 优化建议
- [RATE_LIMIT_GUIDE.md](RATE_LIMIT_GUIDE.md) - 限流功能指南

---

## 🎓 进阶话题

### 分布式追踪集成

请求 ID 可以作为分布式追踪系统（如 Jaeger、Zipkin）的 Trace ID 的一部分：

```go
import "github.com/opentracing/opentracing-go"

func TracingMiddleware(c *gin.Context) {
    requestID := middleware.GetRequestID(c)
    
    // 创建 span，使用请求 ID 作为标签
    span := opentracing.StartSpan("http.request")
    span.SetTag("request_id", requestID)
    defer span.Finish()
    
    c.Next()
}
```

### 关联 ID 模式

除了请求 ID，还可以添加其他关联 ID：

```go
// 关联 ID 常量
const (
    RequestIDKey     = "request_id"
    CorrelationIDKey = "correlation_id"  // 业务关联 ID
    SessionIDKey     = "session_id"      // 会话 ID
    UserIDKey        = "user_id"         // 用户 ID
)

// 在日志中包含所有关联 ID
logger.Info("Processing order",
    zap.String("request_id", requestID),
    zap.String("correlation_id", orderID),
    zap.String("session_id", sessionID),
    zap.Uint("user_id", userID))
```

---

**文档更新时间**: 2024-11-11  
**版本**: 1.0.0  
**作者**: AI Assistant

