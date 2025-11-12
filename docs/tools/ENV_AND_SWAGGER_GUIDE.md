# 环境配置和 Swagger 文档使用指南

## 🔧 环境配置

### 概述

项目支持根据运行环境加载不同的配置文件，实现开发、测试、生产环境的灵活切换。

### 配置文件结构

```
config/
├── config.yaml          # 默认配置
├── config.dev.yaml      # 开发环境配置 ⭐
├── config.test.yaml     # 测试环境配置 ⭐
└── config.prod.yaml     # 生产环境配置 ⭐
```

### 环境切换方式

#### 方式1: 环境变量（推荐）

```bash
# 开发环境（默认）
export GO_ENV=dev
./bin/frontend

# 测试环境
export GO_ENV=test
./bin/frontend

# 生产环境
export GO_ENV=prod
./bin/frontend
```

**Windows PowerShell:**
```powershell
$env:GO_ENV="prod"
.\bin\frontend.exe
```

**Windows CMD:**
```cmd
set GO_ENV=prod
bin\frontend.exe
```

#### 方式2: 内联环境变量

```bash
# 一次性指定环境
GO_ENV=prod ./bin/frontend

# 或
GO_ENV=test ./bin/backend
```

### 环境配置对比

| 配置项 | 开发环境 (dev) | 测试环境 (test) | 生产环境 (prod) |
|--------|---------------|----------------|-----------------|
| **服务模式** | debug | release | release |
| **数据库** | localhost | test-mysql-server | prod-mysql-server |
| **日志级别** | debug | info | info |
| **日志格式** | console | json | json |
| **Token有效期** | 用户7天/管理员1天 | 用户1天/管理员8小时 | 用户7天/管理员12小时 |
| **连接池大小** | 较小 | 中等 | 较大 |

### 配置详情

#### 开发环境 (config.dev.yaml)

```yaml
server:
  env: dev
  mode: debug

database:
  mysql:
    host: localhost
    database: trx_dev
    max_open_conns: 100

logger:
  level: debug
  encoding: console  # 便于阅读

jwt:
  secret: dev-secret-key
```

#### 测试环境 (config.test.yaml)

```yaml
server:
  env: test
  mode: release

database:
  mysql:
    host: test-mysql-server
    database: trx_test
    password: test_password

logger:
  level: info
  encoding: json

jwt:
  expire_hours: 24  # 缩短有效期
```

#### 生产环境 (config.prod.yaml)

```yaml
server:
  env: prod
  mode: release

database:
  mysql:
    host: prod-mysql-server
    database: trx_prod
    password: CHANGE_ME_IN_PRODUCTION
    max_open_conns: 200  # 更大的连接池

redis:
  password: CHANGE_ME_IN_PRODUCTION
  pool_size: 50

kafka:
  brokers:
    - prod-kafka-server1:9092
    - prod-kafka-server2:9092
    - prod-kafka-server3:9092

logger:
  output_paths:
    - /var/log/trx-project/app.log

jwt:
  secret: CHANGE_ME_IN_PRODUCTION_USE_STRONG_SECRET_KEY
```

### 启动示例

```bash
# 开发环境（默认）
./bin/frontend
✅ 加载环境配置: config/config.dev.yaml (环境: dev)

# 测试环境
GO_ENV=test ./bin/frontend
✅ 加载环境配置: config/config.test.yaml (环境: test)

# 生产环境
GO_ENV=prod ./bin/frontend
✅ 加载环境配置: config/config.prod.yaml (环境: prod)
```

### 配置优先级

1. **环境变量 `GO_ENV`** (最高优先级)
2. 默认值 `dev`

### 注意事项

⚠️ **生产环境配置安全**

1. **不要提交敏感信息到 Git**
   ```bash
   # 生产环境配置建议加入 .gitignore
   echo "config/config.prod.yaml" >> .gitignore
   ```

2. **使用环境变量覆盖敏感配置**
   ```bash
   export DB_PASSWORD="your_secure_password"
   export JWT_SECRET="your_strong_jwt_secret"
   export REDIS_PASSWORD="your_redis_password"
   ```

3. **生产环境部署清单**
   - [ ] 修改数据库密码
   - [ ] 修改 Redis 密码
   - [ ] 修改 JWT Secret
   - [ ] 检查 Kafka 集群配置
   - [ ] 确认日志路径权限
   - [ ] 设置合适的连接池大小

---

## 📚 Swagger 文档

### 访问地址

**前台 API 文档:**
```
http://localhost:8080/swagger/index.html
```

**后台 API 文档:**
```
http://localhost:8081/swagger/index.html
```

### 文档特点

✅ **完整的中文注释** - 所有 API 都有详细的中文说明  
✅ **可交互测试** - 直接在浏览器中测试 API  
✅ **JWT 认证支持** - 内置 Token 认证功能  
✅ **自动生成** - 基于代码注释自动生成文档  

### Swagger UI 功能

#### 1. 查看 API 列表

- 按 Tag 分组展示
- 显示请求方法、路径、描述
- 支持搜索和过滤

#### 2. 查看 API 详情

点击任一 API，可以看到：
- **请求参数**: 参数类型、是否必填、说明
- **请求示例**: JSON 格式的请求体示例
- **响应示例**: 成功和失败的响应格式
- **状态码说明**: 各种 HTTP 状态码的含义

#### 3. 在线测试 API

##### 无需认证的接口（如注册、登录）

1. 点击 **"Try it out"**
2. 填写请求参数
3. 点击 **"Execute"**
4. 查看响应结果

**示例：测试用户注册**
```json
{
  "username": "testuser",
  "email": "test@example.com",
  "password": "password123"
}
```

##### 需要认证的接口

1. **获取 Token**
   - 先调用登录接口获取 JWT Token
   - 或使用工具生成管理员 Token：
     ```bash
     go run scripts/generate_admin_token.go
     ```

2. **配置 Token**
   - 点击右上角 **"Authorize"** 按钮
   - 在弹出框中输入：`Bearer <your_token>`
   - 点击 **"Authorize"**

3. **测试接口**
   - 现在可以测试需要认证的接口了
   - Token 会自动添加到请求头

**示例：前台用户认证**
```
Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**示例：后台管理员认证**
```bash
# 1. 生成管理员 Token
go run scripts/generate_admin_token.go

# 2. 复制输出的 Token
# 3. 在 Swagger UI 中点击 "Authorize"
# 4. 输入: Bearer <复制的Token>
```

### API 分组说明

#### 前台 API (http://localhost:8080/swagger/index.html)

**公开接口** (不需要认证)
- 用户注册 `POST /api/v1/public/register`
- 用户登录 `POST /api/v1/public/login`

**用户接口** (需要用户认证)
- 获取个人信息 `GET /api/v1/user/profile`
- 更新个人信息 `PUT /api/v1/user/profile`
- 获取用户列表 `GET /api/v1/users`
- 删除用户 `DELETE /api/v1/users/{id}`

#### 后台 API (http://localhost:8081/swagger/index.html)

**用户管理** (需要管理员认证)
- 获取用户列表 `GET /api/v1/admin/users`
- 获取用户详情 `GET /api/v1/admin/users/{id}`
- 更新用户状态 `PUT /api/v1/admin/users/{id}/status`
- 删除用户 `DELETE /api/v1/admin/users/{id}`
- 重置用户密码 `POST /api/v1/admin/users/{id}/reset-password`

**统计信息** (需要管理员认证)
- 用户统计 `GET /api/v1/admin/statistics/users`

### 更新 Swagger 文档

#### 添加新的 API

1. **在 Handler 中添加 Swagger 注释**

```go
// CreateProduct 创建商品
// @Summary 创建商品
// @Description 创建新的商品信息
// @Tags 商品管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ProductRequest true "商品信息"
// @Success 201 {object} response.Response{data=model.Product} "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Router /admin/products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
    // 实现代码...
}
```

2. **重新生成 Swagger 文档**

```bash
# 前台
cd cmd/frontend && swag init --parseDependency --parseInternal

# 后台
cd cmd/backend && swag init --parseDependency --parseInternal
```

3. **重新编译和运行**

```bash
go build -o bin/frontend ./cmd/frontend
./bin/frontend
```

4. **刷新浏览器** - 新的 API 文档已生效

### Swagger 注释说明

#### 常用注解

| 注解 | 说明 | 示例 |
|------|------|------|
| `@Summary` | API 简短描述 | `@Summary 用户注册` |
| `@Description` | API 详细描述 | `@Description 创建新用户账号` |
| `@Tags` | API 分组 | `@Tags 公开接口` |
| `@Accept` | 接受的内容类型 | `@Accept json` |
| `@Produce` | 返回的内容类型 | `@Produce json` |
| `@Security` | 认证方式 | `@Security BearerAuth` |
| `@Param` | 参数说明 | `@Param id path int true "用户ID"` |
| `@Success` | 成功响应 | `@Success 200 {object} response.Response` |
| `@Failure` | 失败响应 | `@Failure 400 {object} response.Response` |
| `@Router` | 路由路径 | `@Router /api/v1/users [get]` |

#### 参数类型

- `path` - 路径参数 (如 `/users/{id}`)
- `query` - 查询参数 (如 `?page=1`)
- `body` - 请求体
- `header` - 请求头

#### 数据类型

- `int`, `string`, `bool`
- `object` - 对象，需指定类型 `{object} model.User`
- `array` - 数组，如 `{array} model.User`
- `map[string]interface{}` - Map 类型

### Swagger 配置

主配置在 `cmd/frontend/docs.go` 和 `cmd/backend/docs.go`:

```go
// @title TRX Project - 前台 API
// @version 1.0
// @description 基于 Gin 框架的现代化 Go Web 服务 - 前台接口

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description 输入 "Bearer " + JWT Token
```

### 最佳实践

1. **使用中文注释** - 便于团队阅读理解
2. **详细的参数说明** - 包括类型、必填、示例
3. **完整的响应示例** - 成功和失败情况都要说明
4. **合理的 Tag 分组** - 按功能模块分组
5. **及时更新文档** - 代码变更后重新生成文档

### 故障排查

#### 问题1：Swagger 页面空白

```bash
# 检查是否正确生成文档
ls -l cmd/frontend/docs/
ls -l cmd/backend/docs/

# 重新生成
cd cmd/frontend && swag init --parseDependency --parseInternal
```

#### 问题2：新 API 不显示

```bash
# 确保添加了正确的注释
# 重新生成文档
cd cmd/frontend && swag init --parseDependency --parseInternal

# 重新编译
go build -o bin/frontend ./cmd/frontend

# 重启服务并刷新浏览器（硬刷新 Ctrl+F5）
```

#### 问题3：Token 认证失败

- 确认 Token 格式：必须是 `Bearer <token>`
- 检查 Token 是否过期
- 查看浏览器开发者工具的 Network 标签

---

## 🎯 快速开始

### 1. 启动服务（开发环境）

```bash
# 默认使用开发环境配置
./bin/frontend  # 前台 - http://localhost:8080
./bin/backend   # 后台 - http://localhost:8081
```

### 2. 访问 Swagger 文档

- 前台文档: http://localhost:8080/swagger/index.html
- 后台文档: http://localhost:8081/swagger/index.html

### 3. 测试 API

#### 前台接口测试

1. **注册用户**
   - 打开前台 Swagger 文档
   - 找到 `POST /public/register`
   - 点击 "Try it out"
   - 填写注册信息
   - 执行并获取 Token

2. **使用 Token 访问需认证接口**
   - 点击右上角 "Authorize"
   - 输入 `Bearer <刚才获取的Token>`
   - 测试 `/user/profile` 接口

#### 后台接口测试

1. **生成管理员 Token**
   ```bash
   go run scripts/generate_admin_token.go
   ```

2. **配置认证**
   - 打开后台 Swagger 文档
   - 点击 "Authorize"
   - 输入 `Bearer <管理员Token>`

3. **测试管理接口**
   - 测试 `GET /admin/users` 获取用户列表
   - 测试 `GET /admin/statistics/users` 获取统计信息

### 4. 切换到生产环境

```bash
# 确保已配置好生产环境配置文件
GO_ENV=prod ./bin/frontend
GO_ENV=prod ./bin/backend
```

---

## 📝 总结

### 环境配置功能

✅ 支持多环境配置（dev/test/prod）  
✅ 通过环境变量灵活切换  
✅ 不同环境独立的数据库、Redis、Kafka 配置  
✅ 针对环境优化的日志和性能配置  

### Swagger 文档功能

✅ 完整的中文 API 文档  
✅ 在线交互式测试  
✅ JWT 认证支持  
✅ 前后台独立文档  
✅ 自动生成和更新  

### 下一步

- 根据实际需求调整各环境配置
- 添加更多 API 并更新 Swagger 文档
- 在生产环境部署时修改敏感配置
- 使用 Swagger 文档进行前后端协作

有问题？查看完整文档或联系团队！🚀

