# 快速开始指南

本指南将帮助你快速启动和运行 TRX 项目。

## 前置要求

确保你的系统已安装以下软件：

- **Go 1.21+**: [下载 Go](https://golang.org/dl/)
- **Docker**: [安装 Docker](https://docs.docker.com/get-docker/)
- **Docker Compose**: [安装 Docker Compose](https://docs.docker.com/compose/install/)
- **Git**: [安装 Git](https://git-scm.com/downloads)

## 一键启动（推荐）

### 方式一：使用 setup 脚本

```bash
# 运行初始化脚本（Linux/Mac）
./scripts/setup.sh

# Windows 用户可以手动执行以下步骤
```

### 方式二：手动步骤

#### 1. 安装依赖

```bash
# 安装 Go 依赖
go mod download
go mod tidy

# 安装 Wire
go install github.com/google/wire/cmd/wire@latest
```

#### 2. 生成 Wire 代码

```bash
cd cmd/server
wire
cd ../..
```

#### 3. 配置文件

```bash
# 复制配置示例
cp config/config.yaml.example config/config.yaml

# 根据需要编辑配置文件
# vim config/config.yaml
```

#### 4. 启动依赖服务

```bash
# 启动 MySQL, Redis, Kafka
docker-compose up -d

# 等待服务启动（约10秒）
sleep 10

# 检查服务状态
docker-compose ps
```

#### 5. 创建日志目录

```bash
mkdir -p logs
```

#### 6. 构建项目

```bash
# 构建可执行文件
go build -o bin/server ./cmd/server

# 或使用 Make
make build
```

#### 7. 运行服务

```bash
# 直接运行
./bin/server

# 或使用 Make
make run

# 或使用开发模式（热重载）
make dev  # 需要安装 air: go install github.com/air-verse/air@latest
```

## 验证安装

### 1. 健康检查

```bash
curl http://localhost:8080/health
```

预期响应：
```json
{
  "status": "ok"
}
```

### 2. 注册用户

```bash
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'
```

### 3. 用户登录

```bash
curl -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

### 4. 获取用户列表

```bash
curl http://localhost:8080/api/v1/users?page=1&page_size=10
```

## 常用命令

### Make 命令

```bash
make help         # 显示所有可用命令
make deps         # 安装依赖
make wire         # 生成 Wire 代码
make build        # 构建项目
make run          # 运行项目
make test         # 运行测试
make clean        # 清理构建文件
make docker-up    # 启动 Docker 服务
make docker-down  # 停止 Docker 服务
```

### Docker Compose 命令

```bash
# 启动所有服务
docker-compose up -d

# 停止所有服务
docker-compose down

# 查看服务状态
docker-compose ps

# 查看服务日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f mysql
docker-compose logs -f redis
docker-compose logs -f kafka

# 重启服务
docker-compose restart

# 重新构建并启动
docker-compose up -d --build
```

### Kafka 命令

```bash
# 进入 Kafka 容器
docker exec -it trx-kafka bash

# 创建 Topic
kafka-topics.sh --create --topic user-events --bootstrap-server localhost:9092 --partitions 3 --replication-factor 1

# 列出所有 Topics
kafka-topics.sh --list --bootstrap-server localhost:9092

# 查看 Topic 详情
kafka-topics.sh --describe --topic user-events --bootstrap-server localhost:9092

# 生产消息（测试）
kafka-console-producer.sh --topic user-events --bootstrap-server localhost:9092

# 消费消息（测试）
kafka-console-consumer.sh --topic user-events --from-beginning --bootstrap-server localhost:9092
```

### MySQL 命令

```bash
# 连接到 MySQL
docker exec -it trx-mysql mysql -uroot -ppassword

# 查看数据库
SHOW DATABASES;

# 使用数据库
USE trx_db;

# 查看表
SHOW TABLES;

# 查看用户表
SELECT * FROM users;
```

### Redis 命令

```bash
# 连接到 Redis
docker exec -it trx-redis redis-cli

# 查看所有键
KEYS *

# 获取值
GET key_name

# 清空数据库
FLUSHDB
```

## 开发模式

### 使用 Air 进行热重载

1. 安装 Air：
```bash
go install github.com/air-verse/air@latest
```

2. 运行开发模式：
```bash
make dev
# 或
air
```

现在，当你修改代码时，服务会自动重新编译和重启。

## 故障排查

### 端口已被占用

如果遇到端口冲突，修改 `config/config.yaml` 或 `docker-compose.yml` 中的端口：

```yaml
# config/config.yaml
server:
  port: 8081  # 修改为其他端口

# docker-compose.yml
services:
  mysql:
    ports:
      - "3307:3306"  # 修改为其他端口
```

### Docker 服务启动失败

```bash
# 查看日志
docker-compose logs

# 重新创建容器
docker-compose down -v
docker-compose up -d
```

### Wire 生成失败

```bash
# 确保在正确的目录
cd cmd/server

# 重新安装 Wire
go install github.com/google/wire/cmd/wire@latest

# 再次生成
wire
```

### 数据库连接失败

1. 确保 MySQL 容器正在运行：
```bash
docker-compose ps
```

2. 检查配置文件中的数据库配置

3. 等待 MySQL 完全启动（大约 10-20 秒）

## 下一步

- 阅读 [API 文档](docs/API.md)
- 了解 [项目架构](docs/ARCHITECTURE.md)
- 学习 [Kafka 使用](docs/KAFKA_USAGE.md)
- 编写测试代码
- 添加新的功能模块

## 获取帮助

- 查看 [README.md](README.md) 了解项目详情
- 查看 [docs/](docs/) 目录下的详细文档
- 提交 Issue 报告问题

## 停止服务

### 停止应用

```bash
# 如果是直接运行
Ctrl + C

# 或者找到进程并终止
ps aux | grep server
kill -9 <PID>
```

### 停止 Docker 服务

```bash
# 停止服务但保留数据
docker-compose stop

# 停止服务并删除容器（保留数据卷）
docker-compose down

# 停止服务并删除所有数据
docker-compose down -v
```

## 生产部署

在生产环境部署时，请注意：

1. 修改配置文件中的敏感信息
2. 使用环境变量管理配置
3. 启用 HTTPS
4. 配置日志轮转
5. 设置监控和告警
6. 配置数据备份
7. 使用反向代理（Nginx）
8. 配置防火墙规则

更多部署信息，请参考 [ARCHITECTURE.md](docs/ARCHITECTURE.md)。

