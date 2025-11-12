# 快速测试指南

本指南帮助你快速测试前后台接口分离功能。

## 准备工作

### 1. 启动依赖服务

```bash
# 启动 MySQL、Redis、Kafka
make docker-up

# 等待服务启动（约10秒）
sleep 10

# 检查服务状态
docker-compose ps
```

### 2. 启动应用

```bash
# 创建日志目录
mkdir -p logs

# 运行应用
make run

# 或者使用热重载模式
make dev
```

## 测试前台接口

### 方式一：使用测试脚本

```bash
./scripts/test_api.sh
```

### 方式二：手动测试

#### 1. 注册用户

```bash
curl -X POST http://localhost:8080/api/v1/public/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'
```

**预期响应**:
```json
{
  "code": 200,
  "message": "User registered successfully",
  "data": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "status": 1
  }
}
```

#### 2. 用户登录

```bash
curl -X POST http://localhost:8080/api/v1/public/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

**预期响应**:
```json
{
  "code": 200,
  "message": "Login successful",
  "data": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com"
  }
}
```

#### 3. 测试参数验证

```bash
curl -X POST http://localhost:8080/api/v1/public/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "ab",
    "email": "invalid-email",
    "password": "123"
  }'
```

**预期响应**:
```json
{
  "code": 10001,
  "message": "Key: 'RegisterRequest.Username' Error:Field validation for 'Username' failed on the 'min' tag"
}
```

## 测试后台接口

### 方式一：使用测试脚本

```bash
./scripts/test_admin_api.sh
```

### 方式二：手动测试

**管理员 Token**: `admin_test_token_123456`

#### 1. 测试无 Token 访问（应该失败）

```bash
curl http://localhost:8080/api/v1/admin/users
```

**预期响应**:
```json
{
  "code": 401,
  "message": "Missing authorization token"
}
```

#### 2. 测试无效 Token（应该失败）

```bash
curl http://localhost:8080/api/v1/admin/users \
  -H "Authorization: Bearer invalid_token"
```

**预期响应**:
```json
{
  "code": 403,
  "message": "Admin access required"
}
```

#### 3. 测试有效管理员 Token（应该成功）

```bash
# 设置管理员 Token
ADMIN_TOKEN="admin_test_token_123456"

# 获取用户列表
curl "http://localhost:8080/api/v1/admin/users?page=1&page_size=10" \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

**预期响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [...],
    "total": 1,
    "page": 1,
    "page_size": 10
  }
}
```

#### 4. 获取用户详情

```bash
curl http://localhost:8080/api/v1/admin/users/1 \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

#### 5. 更新用户状态

```bash
curl -X PUT http://localhost:8080/api/v1/admin/users/1/status \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"status": 0}'
```

**预期响应**:
```json
{
  "code": 200,
  "message": "User status updated successfully",
  "data": {
    "id": 1,
    "status": 0
  }
}
```

#### 6. 重置用户密码

```bash
curl -X POST http://localhost:8080/api/v1/admin/users/1/reset-password \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"new_password": "newpassword123"}'
```

#### 7. 获取用户统计

```bash
curl http://localhost:8080/api/v1/admin/statistics/users \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

**预期响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "total_users": 1000,
    "active_users": 850,
    "inactive_users": 150,
    "new_users_today": 12,
    "new_users_week": 89,
    "new_users_month": 356
  }
}
```

## 验证清单

### ✅ 前台接口
- [ ] 用户注册成功
- [ ] 用户登录成功
- [ ] 参数验证错误正确返回
- [ ] 用户已存在错误正确返回
- [ ] 密码错误正确返回

### ✅ 后台接口
- [ ] 无 Token 访问返回 401
- [ ] 无效 Token 返回 403
- [ ] 有效管理员 Token 可以访问
- [ ] 获取用户列表成功（分页）
- [ ] 获取用户详情成功
- [ ] 更新用户状态成功
- [ ] 重置密码功能正常
- [ ] 获取统计信息成功

### ✅ 响应格式
- [ ] 所有响应包含 code 和 message
- [ ] 成功响应包含 data
- [ ] 分页响应格式正确
- [ ] 错误响应格式正确

## 常见问题

### 1. 服务启动失败

**问题**: 端口被占用

```bash
# 检查端口占用
lsof -i :8080  # Mac/Linux
netstat -ano | findstr :8080  # Windows

# 修改配置文件中的端口
vim config/config.yaml
```

### 2. 数据库连接失败

**问题**: MySQL 未启动或连接配置错误

```bash
# 检查 MySQL 容器
docker-compose ps

# 重启 MySQL
docker-compose restart mysql

# 查看 MySQL 日志
docker-compose logs mysql
```

### 3. 测试脚本权限问题

**问题**: Permission denied

```bash
# 添加执行权限
chmod +x scripts/*.sh
```

### 4. 测试时返回 404

**问题**: 路由未正确配置

```bash
# 检查日志
tail -f logs/app.log

# 确认路由
curl http://localhost:8080/health
```

## 使用 Postman 测试

### 导入集合

创建以下 Postman 集合：

#### 前台接口

```json
{
  "name": "前台接口",
  "requests": [
    {
      "name": "注册",
      "method": "POST",
      "url": "http://localhost:8080/api/v1/public/register",
      "body": {
        "username": "testuser",
        "email": "test@example.com",
        "password": "password123"
      }
    },
    {
      "name": "登录",
      "method": "POST",
      "url": "http://localhost:8080/api/v1/public/login",
      "body": {
        "username": "testuser",
        "password": "password123"
      }
    }
  ]
}
```

#### 后台接口

```json
{
  "name": "后台接口",
  "auth": {
    "type": "bearer",
    "bearer": "admin_test_token_123456"
  },
  "requests": [
    {
      "name": "用户列表",
      "method": "GET",
      "url": "http://localhost:8080/api/v1/admin/users?page=1&page_size=10"
    },
    {
      "name": "用户详情",
      "method": "GET",
      "url": "http://localhost:8080/api/v1/admin/users/1"
    },
    {
      "name": "更新状态",
      "method": "PUT",
      "url": "http://localhost:8080/api/v1/admin/users/1/status",
      "body": {
        "status": 0
      }
    }
  ]
}
```

## 性能测试

使用 Apache Bench 进行简单的性能测试：

```bash
# 测试健康检查接口
ab -n 1000 -c 10 http://localhost:8080/health

# 测试用户列表接口（后台）
ab -n 100 -c 5 -H "Authorization: Bearer admin_test_token_123456" \
   http://localhost:8080/api/v1/admin/users
```

## 下一步

完成测试后，你可以：

1. **集成 JWT**: 实现真实的 Token 生成和验证
2. **添加权限**: 实现基于角色的权限控制
3. **完善功能**: 添加更多管理接口
4. **编写测试**: 添加单元测试和集成测试
5. **性能优化**: 添加缓存、限流等

## 相关文档

- [前后台接口详细文档](docs/API_FRONTEND_BACKEND.md)
- [API 接口文档](docs/API.md)
- [快速开始指南](QUICKSTART.md)
- [实现总结](FRONTEND_BACKEND_SEPARATION.md)

## 获取帮助

遇到问题？

1. 查看日志: `tail -f logs/app.log`
2. 检查配置: `cat config/config.yaml`
3. 查看文档: `docs/` 目录
4. 运行验证: `./scripts/verify.sh`

祝测试顺利！🚀

