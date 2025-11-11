# Kafka 使用指南

## 简介

本项目集成了 Kafka 消息队列，用于异步消息处理和服务解耦。

## Kafka 组件

项目提供了 Kafka Producer 和 Consumer 的封装：

- **Producer**: 发送消息到 Kafka Topic
- **Consumer**: 从 Kafka Topic 消费消息

## 使用示例

### 1. 发送消息 (Producer)

```go
package main

import (
    "context"
    "encoding/json"
    "trx-project/pkg/kafka"
    "trx-project/pkg/config"
    "trx-project/pkg/logger"
)

func main() {
    // 初始化配置和日志
    cfg, _ := config.Load("config/config.yaml")
    logger.InitLogger(&cfg.Logger)
    
    // 创建 Producer
    producer := kafka.NewProducer(&cfg.Kafka, logger.Logger)
    defer producer.Close()
    
    // 准备消息数据
    type UserEvent struct {
        UserID   uint   `json:"user_id"`
        Username string `json:"username"`
        Action   string `json:"action"`
    }
    
    event := UserEvent{
        UserID:   1,
        Username: "testuser",
        Action:   "register",
    }
    
    // 序列化消息
    message, _ := json.Marshal(event)
    
    // 发送消息
    ctx := context.Background()
    err := producer.SendMessage(ctx, "user-events", "user-1", message)
    if err != nil {
        logger.Error("Failed to send message", zap.Error(err))
    }
}
```

### 2. 消费消息 (Consumer)

```go
package main

import (
    "context"
    "encoding/json"
    "trx-project/pkg/kafka"
    "trx-project/pkg/config"
    "trx-project/pkg/logger"
    "go.uber.org/zap"
)

func main() {
    // 初始化配置和日志
    cfg, _ := config.Load("config/config.yaml")
    logger.InitLogger(&cfg.Logger)
    
    // 创建 Consumer
    consumer := kafka.NewConsumer(&cfg.Kafka, "user-events", logger.Logger)
    defer consumer.Close()
    
    // 消费消息
    ctx := context.Background()
    for {
        msg, err := consumer.ReadMessage(ctx)
        if err != nil {
            logger.Error("Failed to read message", zap.Error(err))
            continue
        }
        
        // 处理消息
        logger.Info("Received message",
            zap.String("topic", msg.Topic),
            zap.String("key", string(msg.Key)),
            zap.String("value", string(msg.Value)))
        
        // 解析消息
        type UserEvent struct {
            UserID   uint   `json:"user_id"`
            Username string `json:"username"`
            Action   string `json:"action"`
        }
        
        var event UserEvent
        if err := json.Unmarshal(msg.Value, &event); err != nil {
            logger.Error("Failed to unmarshal message", zap.Error(err))
            continue
        }
        
        // 业务处理
        processUserEvent(event)
    }
}

func processUserEvent(event UserEvent) {
    // 处理用户事件
    logger.Info("Processing user event",
        zap.Uint("user_id", event.UserID),
        zap.String("action", event.Action))
}
```

### 3. 在 Service 层集成 Kafka

```go
package service

import (
    "context"
    "encoding/json"
    "trx-project/internal/model"
    "trx-project/internal/repository"
    "trx-project/pkg/kafka"
    
    "github.com/redis/go-redis/v9"
    "go.uber.org/zap"
)

type UserService interface {
    Register(ctx context.Context, username, email, password string) (*model.User, error)
}

type userService struct {
    repo     repository.UserRepository
    redis    *redis.Client
    producer *kafka.Producer
    logger   *zap.Logger
}

func NewUserService(
    repo repository.UserRepository,
    redis *redis.Client,
    producer *kafka.Producer,
    logger *zap.Logger,
) UserService {
    return &userService{
        repo:     repo,
        redis:    redis,
        producer: producer,
        logger:   logger,
    }
}

func (s *userService) Register(ctx context.Context, username, email, password string) (*model.User, error) {
    // 创建用户
    user := &model.User{
        Username: username,
        Email:    email,
        Password: hashPassword(password),
        Status:   1,
    }
    
    if err := s.repo.Create(ctx, user); err != nil {
        return nil, err
    }
    
    // 发送注册事件到 Kafka
    event := map[string]interface{}{
        "user_id":  user.ID,
        "username": user.Username,
        "email":    user.Email,
        "action":   "register",
    }
    
    eventData, _ := json.Marshal(event)
    if err := s.producer.SendMessage(ctx, "user-events", username, eventData); err != nil {
        s.logger.Error("Failed to send user event", zap.Error(err))
        // 不影响主流程，仅记录错误
    }
    
    return user, nil
}
```

## 配置说明

在 `config/config.yaml` 中配置 Kafka：

```yaml
kafka:
  brokers:
    - localhost:9092
    # - broker2:9092  # 可配置多个 broker
  group_id: trx-group
  topics:
    - user-events
    - order-events
    - notification-events
```

## 常见使用场景

### 1. 异步通知

用户注册后，发送欢迎邮件：

```go
// 在注册成功后发送事件
event := map[string]interface{}{
    "type": "welcome_email",
    "user_id": user.ID,
    "email": user.Email,
}
producer.SendMessage(ctx, "notification-events", userID, eventData)

// Consumer 处理发送邮件
```

### 2. 数据同步

将用户数据同步到其他服务：

```go
// 用户信息更新后发送事件
event := map[string]interface{}{
    "type": "user_updated",
    "user_id": user.ID,
    "changes": changes,
}
producer.SendMessage(ctx, "user-events", userID, eventData)
```

### 3. 日志收集

收集业务日志：

```go
// 记录用户行为
event := map[string]interface{}{
    "user_id": userID,
    "action": "login",
    "timestamp": time.Now(),
    "ip": clientIP,
}
producer.SendMessage(ctx, "user-behavior", userID, eventData)
```

### 4. 事件溯源

记录所有业务事件，用于审计和数据恢复：

```go
// 记录订单创建事件
event := map[string]interface{}{
    "event_type": "order_created",
    "order_id": orderID,
    "user_id": userID,
    "amount": amount,
    "timestamp": time.Now(),
}
producer.SendMessage(ctx, "order-events", orderID, eventData)
```

## 最佳实践

### 1. 消息格式

使用统一的消息格式：

```go
type Event struct {
    ID        string                 `json:"id"`
    Type      string                 `json:"type"`
    Timestamp time.Time              `json:"timestamp"`
    Data      map[string]interface{} `json:"data"`
}
```

### 2. 错误处理

- Producer: 发送失败不应影响主流程
- Consumer: 处理失败应有重试机制

### 3. 幂等性

确保消息消费的幂等性，防止重复处理：

```go
func (h *EventHandler) HandleMessage(msg kafka.Message) error {
    // 检查是否已处理
    processed, _ := h.redis.Get(ctx, "processed:"+messageID).Result()
    if processed == "1" {
        return nil // 已处理，跳过
    }
    
    // 处理消息
    if err := h.processMessage(msg); err != nil {
        return err
    }
    
    // 标记为已处理
    h.redis.Set(ctx, "processed:"+messageID, "1", 24*time.Hour)
    return nil
}
```

### 4. 监控

监控 Kafka 消息的：
- 发送成功率
- 消费延迟
- 消息堆积

## 创建独立的 Consumer 服务

创建 `cmd/consumer/main.go`:

```go
package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"
    "trx-project/pkg/config"
    "trx-project/pkg/kafka"
    "trx-project/pkg/logger"
    
    "go.uber.org/zap"
)

func main() {
    // 加载配置
    cfg, err := config.Load("config/config.yaml")
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }
    
    // 初始化日志
    if err := logger.InitLogger(&cfg.Logger); err != nil {
        log.Fatalf("Failed to init logger: %v", err)
    }
    defer logger.Sync()
    
    // 创建 Consumer
    consumer := kafka.NewConsumer(&cfg.Kafka, "user-events", logger.Logger)
    defer consumer.Close()
    
    // 启动消费
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    
    go func() {
        for {
            select {
            case <-ctx.Done():
                return
            default:
                msg, err := consumer.ReadMessage(ctx)
                if err != nil {
                    logger.Error("Failed to read message", zap.Error(err))
                    continue
                }
                
                // 处理消息
                handleMessage(msg)
            }
        }
    }()
    
    // 等待退出信号
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    logger.Info("Consumer shutting down...")
}

func handleMessage(msg kafka.Message) {
    logger.Info("Processing message",
        zap.String("topic", msg.Topic),
        zap.String("key", string(msg.Key)),
        zap.ByteString("value", msg.Value))
    
    // 实现具体的消息处理逻辑
}
```

## 故障排查

### 1. 连接失败

检查 Kafka broker 地址和端口是否正确：

```bash
# 测试连接
telnet localhost 9092
```

### 2. 消息发送失败

- 检查 Topic 是否存在
- 检查 Kafka 服务状态
- 查看日志输出

### 3. 消息堆积

- 增加 Consumer 数量
- 优化消息处理逻辑
- 检查消费性能

## 相关命令

```bash
# 创建 Topic
kafka-topics.sh --create --topic user-events --bootstrap-server localhost:9092

# 查看 Topic 列表
kafka-topics.sh --list --bootstrap-server localhost:9092

# 查看 Topic 详情
kafka-topics.sh --describe --topic user-events --bootstrap-server localhost:9092

# 查看消费者组
kafka-consumer-groups.sh --list --bootstrap-server localhost:9092

# 查看消费者组详情
kafka-consumer-groups.sh --describe --group trx-group --bootstrap-server localhost:9092
```

