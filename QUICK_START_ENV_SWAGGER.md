# 🚀 环境配置和 Swagger 快速开始

## 1️⃣ 测试环境切换

### 启动服务查看环境加载

```bash
# 开发环境（默认）
./bin/frontend
# 输出: ✅ 加载环境配置: config/config.dev.yaml (环境: dev)

# 测试环境
GO_ENV=test ./bin/frontend
# 输出: ✅ 加载环境配置: config/config.test.yaml (环境: test)

# 生产环境
GO_ENV=prod ./bin/frontend
# 输出: ✅ 加载环境配置: config/config.prod.yaml (环境: prod)
```

### Windows 用户

**PowerShell:**
```powershell
$env:GO_ENV="test"
.\bin\frontend.exe
```

**CMD:**
```cmd
set GO_ENV=test
bin\frontend.exe
```

## 2️⃣ 访问 Swagger 文档

### 启动服务

```bash
# 终端1 - 前台服务
./bin/frontend

# 终端2 - 后台服务
./bin/backend
```

### 打开浏览器

**前台 API 文档:**
```
http://localhost:8080/swagger/index.html
```

**后台 API 文档:**
```
http://localhost:8081/swagger/index.html
```

## 3️⃣ 测试前台 API（使用 Swagger UI）

### 步骤1: 注册用户

1. 打开 http://localhost:8080/swagger/index.html
2. 找到 **"公开接口"** 分组
3. 点击 **"POST /public/register"**
4. 点击 **"Try it out"**
5. 填写注册信息：

```json
{
  "username": "testuser",
  "email": "test@example.com",
  "password": "password123"
}
```

6. 点击 **"Execute"**
7. 查看响应，复制返回的 Token

### 步骤2: 配置认证

1. 点击页面右上角的 **"Authorize"** 🔒 按钮
2. 在弹出框中输入：
```
Bearer <刚才复制的Token>
```
3. 点击 **"Authorize"** 按钮
4. 关闭弹出框

### 步骤3: 测试需要认证的接口

1. 找到 **"用户接口"** 分组
2. 点击 **"GET /user/profile"**
3. 点击 **"Try it out"**
4. 点击 **"Execute"**
5. 成功返回用户信息！✅

## 4️⃣ 测试后台 API（使用 Swagger UI）

### 步骤1: 生成管理员 Token

```bash
go run scripts/generate_admin_token.go
```

输出示例：
```
✅ 管理员 Token 生成成功!

Token:
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

使用示例:
  export ADMIN_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### 步骤2: 配置管理员认证

1. 打开 http://localhost:8081/swagger/index.html
2. 点击页面右上角的 **"Authorize"** 🔒 按钮
3. 在弹出框中输入：
```
Bearer <刚才生成的管理员Token>
```
4. 点击 **"Authorize"** 按钮
5. 关闭弹出框

### 步骤3: 测试管理接口

1. 找到 **"用户管理"** 分组
2. 点击 **"GET /admin/users"**
3. 点击 **"Try it out"**
4. 可以设置查询参数：
   - `page`: 1
   - `page_size`: 10
5. 点击 **"Execute"**
6. 成功返回用户列表！✅

### 步骤4: 测试统计接口

1. 找到 **"统计信息"** 分组
2. 点击 **"GET /admin/statistics/users"**
3. 点击 **"Try it out"**
4. 点击 **"Execute"**
5. 查看用户统计信息

## 5️⃣ 使用 curl 测试（命令行方式）

### 前台接口

```bash
# 1. 注册用户
curl -X POST http://localhost:8080/api/v1/public/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'

# 保存返回的 Token
USER_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# 2. 获取个人信息
curl http://localhost:8080/api/v1/user/profile \
  -H "Authorization: Bearer $USER_TOKEN"
```

### 后台接口

```bash
# 1. 生成管理员 Token
ADMIN_TOKEN=$(go run scripts/generate_admin_token.go | grep "eyJ" | tr -d '\n')

# 2. 获取用户列表
curl http://localhost:8081/api/v1/admin/users \
  -H "Authorization: Bearer $ADMIN_TOKEN"

# 3. 获取统计信息
curl http://localhost:8081/api/v1/admin/statistics/users \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

## 6️⃣ 环境配置对比测试

### 测试数据库配置

```bash
# 开发环境 - 使用 localhost
./bin/frontend
# 连接: localhost:3306/trx_dev

# 测试环境 - 使用测试服务器
GO_ENV=test ./bin/frontend
# 连接: test-mysql-server:3306/trx_test

# 生产环境 - 使用生产服务器
GO_ENV=prod ./bin/frontend
# 连接: prod-mysql-server:3306/trx_prod
```

### 测试日志格式

```bash
# 开发环境 - console 格式（便于阅读）
./bin/frontend
# 日志: 2024/11/11 21:50:00 INFO User registered...

# 生产环境 - json 格式（便于解析）
GO_ENV=prod ./bin/frontend
# 日志: {"level":"info","ts":1699999999,"msg":"User registered",...}
```

## 7️⃣ Swagger 文档特点演示

### 中文注释

所有 API 都有完整的中文说明：
- **Summary**: 简短的中文描述
- **Description**: 详细的中文说明
- **Parameters**: 参数的中文描述
- **Responses**: 响应的中文说明

### 参数示例

点击任何 API，都能看到：
- 参数类型
- 是否必填
- 默认值
- 示例值

### 响应示例

展示了完整的响应格式：
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "user": {
      "id": 1,
      "username": "testuser",
      "email": "test@example.com"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  },
  "timestamp": 1699999999
}
```

## 8️⃣ 常见操作

### 更新 Swagger 文档

```bash
# 添加新 API 后，重新生成文档
cd cmd/frontend && swag init --parseDependency --parseInternal

# 或后台
cd cmd/backend && swag init --parseDependency --parseInternal

# 重新编译
go build -o bin/frontend ./cmd/frontend

# 重启服务，刷新浏览器即可看到新文档
```

### 切换环境测试

```bash
# 开发环境测试
GO_ENV=dev ./bin/frontend &
curl http://localhost:8080/health

# 测试环境测试
GO_ENV=test ./bin/frontend &
curl http://localhost:8080/health
```

### 查看当前环境配置

启动服务时会显示：
```
✅ 加载环境配置: config/config.dev.yaml (环境: dev)
```

## 9️⃣ 故障排查

### Swagger 页面空白？

```bash
# 检查是否生成了文档
ls -l cmd/frontend/docs/
ls -l cmd/backend/docs/

# 重新生成
cd cmd/frontend && swag init --parseDependency --parseInternal
```

### Token 认证失败？

1. 检查 Token 格式：必须是 `Bearer <token>`
2. 确认 Token 未过期
3. 查看浏览器控制台的 Network 标签
4. 确认后台使用的是管理员 Token，前台使用的是用户 Token

### 环境配置未生效？

```bash
# 确认环境变量已设置
echo $GO_ENV

# 查看启动日志，确认加载了正确的配置文件
./bin/frontend
# 应该显示: ✅ 加载环境配置: config/config.xxx.yaml
```

## 🎯 总结

现在你已经掌握了：

✅ **环境配置**
- 通过 `GO_ENV` 切换环境
- 不同环境使用不同的配置文件
- 启动时可以看到加载的配置

✅ **Swagger 文档**
- 访问前后台 API 文档
- 在线测试 API
- JWT 认证配置
- 查看完整的中文注释

✅ **实际测试**
- 前台用户注册/登录
- 后台管理员认证
- 使用 Swagger UI 交互测试
- 使用 curl 命令行测试

**下一步**: 查看 [ENV_AND_SWAGGER_GUIDE.md](ENV_AND_SWAGGER_GUIDE.md) 了解更多详细信息！

---

**有问题？**
- 查看日志输出
- 访问 Swagger 文档测试
- 检查环境变量配置

**祝你使用愉快！** 🚀

