# 部署指南

本文档介绍如何在不同环境中部署 TRX 项目。

## 开发环境部署

### 本地开发

```bash
# 1. 启动依赖服务
docker-compose up -d

# 2. 运行应用
make run

# 或使用热重载
make dev
```

## 测试环境部署

### 使用 Docker

#### 1. 创建 Dockerfile

```dockerfile
# 构建阶段
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 复制依赖文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server

# 运行阶段
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/server .
COPY --from=builder /app/config ./config

# 创建日志目录
RUN mkdir -p logs

# 设置时区
ENV TZ=Asia/Shanghai

# 暴露端口
EXPOSE 8080

# 运行
CMD ["./server"]
```

#### 2. 构建 Docker 镜像

```bash
docker build -t trx-project:latest .
```

#### 3. 运行容器

```bash
docker run -d \
  --name trx-server \
  -p 8080:8080 \
  -v $(pwd)/config:/root/config \
  -v $(pwd)/logs:/root/logs \
  --network=trx-network \
  trx-project:latest
```

### 使用 Docker Compose（完整部署）

创建 `docker-compose.prod.yml`:

```yaml
version: '3.8'

services:
  app:
    build: .
    container_name: trx-app
    ports:
      - "8080:8080"
    volumes:
      - ./config:/root/config
      - ./logs:/root/logs
    depends_on:
      - mysql
      - redis
      - kafka
    environment:
      - TZ=Asia/Shanghai
    networks:
      - trx-network
    restart: always

  mysql:
    image: mysql:8.0
    container_name: trx-mysql
    environment:
      MYSQL_ROOT_PASSWORD: ${MYSQL_ROOT_PASSWORD}
      MYSQL_DATABASE: trx_db
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql
      - ./scripts/init_db.sql:/docker-entrypoint-initdb.d/init.sql
    networks:
      - trx-network
    restart: always

  redis:
    image: redis:7-alpine
    container_name: trx-redis
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    networks:
      - trx-network
    restart: always

  zookeeper:
    image: confluentinc/cp-zookeeper:latest
    container_name: trx-zookeeper
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181
      ZOOKEEPER_TICK_TIME: 2000
    networks:
      - trx-network
    restart: always

  kafka:
    image: confluentinc/cp-kafka:latest
    container_name: trx-kafka
    depends_on:
      - zookeeper
    ports:
      - "9092:9092"
    environment:
      KAFKA_BROKER_ID: 1
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka:9092
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
    networks:
      - trx-network
    restart: always

  nginx:
    image: nginx:alpine
    container_name: trx-nginx
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf
      - ./nginx/conf.d:/etc/nginx/conf.d
      - ./nginx/ssl:/etc/nginx/ssl
    depends_on:
      - app
    networks:
      - trx-network
    restart: always

volumes:
  mysql_data:
  redis_data:

networks:
  trx-network:
    driver: bridge
```

启动：

```bash
docker-compose -f docker-compose.prod.yml up -d
```

## 生产环境部署

### 准备工作

1. **服务器要求**
   - CPU: 2核或以上
   - 内存: 4GB 或以上
   - 磁盘: 20GB 或以上
   - 操作系统: Linux (推荐 Ubuntu 20.04+)

2. **安装必要软件**
   ```bash
   # 更新系统
   sudo apt update && sudo apt upgrade -y
   
   # 安装 Docker
   curl -fsSL https://get.docker.com -o get-docker.sh
   sudo sh get-docker.sh
   
   # 安装 Docker Compose
   sudo curl -L "https://github.com/docker/compose/releases/download/v2.20.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
   sudo chmod +x /usr/local/bin/docker-compose
   ```

### 部署步骤

#### 1. 准备环境变量

创建 `.env` 文件：

```bash
# 数据库配置
MYSQL_ROOT_PASSWORD=your_strong_password
MYSQL_DATABASE=trx_db

# Redis 配置
REDIS_PASSWORD=your_redis_password

# 应用配置
APP_ENV=production
APP_PORT=8080
```

#### 2. 配置 Nginx

创建 `nginx/conf.d/trx-project.conf`:

```nginx
upstream trx_backend {
    server app:8080;
}

server {
    listen 80;
    server_name your-domain.com;

    # 重定向到 HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name your-domain.com;

    # SSL 证书
    ssl_certificate /etc/nginx/ssl/cert.pem;
    ssl_certificate_key /etc/nginx/ssl/key.pem;

    # SSL 配置
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;

    # 日志
    access_log /var/log/nginx/trx-access.log;
    error_log /var/log/nginx/trx-error.log;

    # 客户端最大上传大小
    client_max_body_size 10M;

    location / {
        proxy_pass http://trx_backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # 超时设置
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    # 健康检查
    location /health {
        proxy_pass http://trx_backend/health;
        access_log off;
    }
}
```

#### 3. 部署应用

```bash
# 克隆代码
git clone <your-repo-url> /opt/trx-project
cd /opt/trx-project

# 配置文件
cp config/config.yaml.example config/config.yaml
vim config/config.yaml  # 编辑配置

# 启动服务
docker-compose -f docker-compose.prod.yml up -d

# 查看日志
docker-compose -f docker-compose.prod.yml logs -f
```

#### 4. 验证部署

```bash
# 检查服务状态
docker-compose -f docker-compose.prod.yml ps

# 健康检查
curl http://localhost/health

# 测试 API
curl http://localhost/api/v1/users
```

### 配置 HTTPS（Let's Encrypt）

```bash
# 安装 certbot
sudo apt install certbot python3-certbot-nginx -y

# 获取证书
sudo certbot --nginx -d your-domain.com

# 自动续期（添加到 crontab）
0 12 * * * /usr/bin/certbot renew --quiet
```

## 监控和维护

### 日志管理

#### 配置日志轮转

创建 `/etc/logrotate.d/trx-project`:

```
/opt/trx-project/logs/*.log {
    daily
    rotate 7
    compress
    delaycompress
    missingok
    notifempty
    create 0644 root root
}
```

### 备份策略

#### 数据库备份

```bash
#!/bin/bash
# backup.sh

BACKUP_DIR="/backup/mysql"
DATE=$(date +%Y%m%d_%H%M%S)

# 创建备份目录
mkdir -p $BACKUP_DIR

# 备份数据库
docker exec trx-mysql mysqldump -u root -p${MYSQL_ROOT_PASSWORD} trx_db > $BACKUP_DIR/trx_db_$DATE.sql

# 压缩
gzip $BACKUP_DIR/trx_db_$DATE.sql

# 删除 7 天前的备份
find $BACKUP_DIR -name "*.sql.gz" -mtime +7 -delete

echo "Backup completed: $BACKUP_DIR/trx_db_$DATE.sql.gz"
```

添加到 crontab（每天凌晨 2 点备份）：

```bash
0 2 * * * /opt/trx-project/scripts/backup.sh
```

### 性能监控

#### 使用 Prometheus + Grafana

创建 `docker-compose.monitoring.yml`:

```yaml
version: '3.8'

services:
  prometheus:
    image: prom/prometheus
    container_name: prometheus
    volumes:
      - ./monitoring/prometheus.yml:/etc/prometheus/prometheus.yml
      - prometheus_data:/prometheus
    ports:
      - "9090:9090"
    networks:
      - trx-network

  grafana:
    image: grafana/grafana
    container_name: grafana
    volumes:
      - grafana_data:/var/lib/grafana
    ports:
      - "3000:3000"
    networks:
      - trx-network
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin

volumes:
  prometheus_data:
  grafana_data:

networks:
  trx-network:
    external: true
```

## 更新部署

### 滚动更新

```bash
# 拉取最新代码
cd /opt/trx-project
git pull

# 重新构建镜像
docker-compose -f docker-compose.prod.yml build app

# 重启应用（零停机）
docker-compose -f docker-compose.prod.yml up -d --no-deps --build app
```

### 回滚

```bash
# 查看镜像历史
docker images trx-project

# 回滚到指定版本
docker tag trx-project:previous trx-project:latest
docker-compose -f docker-compose.prod.yml up -d --no-deps app
```

## 故障排查

### 查看日志

```bash
# 应用日志
docker-compose logs -f app

# 数据库日志
docker-compose logs -f mysql

# Nginx 日志
docker-compose logs -f nginx

# 系统日志
tail -f /opt/trx-project/logs/app.log
tail -f /opt/trx-project/logs/error.log
```

### 常见问题

#### 1. 端口被占用

```bash
# 查看端口占用
sudo lsof -i :8080
sudo netstat -tulpn | grep :8080

# 修改端口
vim docker-compose.prod.yml  # 修改端口映射
```

#### 2. 数据库连接失败

```bash
# 检查 MySQL 容器
docker exec -it trx-mysql mysql -uroot -p

# 检查网络
docker network ls
docker network inspect trx-network
```

#### 3. 内存不足

```bash
# 查看内存使用
free -h
docker stats

# 限制容器内存
docker-compose.yml 中添加:
services:
  app:
    mem_limit: 1g
```

## 安全加固

### 1. 防火墙配置

```bash
# 安装 UFW
sudo apt install ufw

# 配置规则
sudo ufw allow 22/tcp   # SSH
sudo ufw allow 80/tcp   # HTTP
sudo ufw allow 443/tcp  # HTTPS
sudo ufw enable
```

### 2. 定期更新

```bash
# 系统更新
sudo apt update && sudo apt upgrade -y

# Docker 镜像更新
docker-compose pull
docker-compose up -d
```

### 3. 密码策略

- 使用强密码
- 定期更换密码
- 使用密钥认证 SSH

## 高可用部署

### 负载均衡

使用多个应用实例 + Nginx 负载均衡：

```yaml
services:
  app1:
    build: .
    ...
  app2:
    build: .
    ...
  app3:
    build: .
    ...

  nginx:
    ...
    # nginx.conf 配置多个 upstream
```

### 数据库主从复制

配置 MySQL 主从复制实现读写分离。

### Redis 集群

使用 Redis Sentinel 或 Redis Cluster 实现高可用。

## 性能优化

1. **启用 Gzip 压缩**（Nginx）
2. **静态资源 CDN**
3. **数据库连接池优化**
4. **Redis 缓存策略**
5. **API 限流**

## 参考资源

- [Docker 官方文档](https://docs.docker.com/)
- [Nginx 官方文档](https://nginx.org/en/docs/)
- [MySQL 官方文档](https://dev.mysql.com/doc/)
- [Redis 官方文档](https://redis.io/documentation)

