# API 响应格式说明

本文档说明项目的统一 API 响应格式规范。

## 统一响应结构

所有 API 接口都遵循统一的响应格式：

```json
{
  "code": 200,
  "message": "success",
  "data": {}
}
```

### 字段说明

| 字段    | 类型   | 必填 | 说明                           |
| ------- | ------ | ---- | ------------------------------ |
| code    | int    | 是   | 业务状态码（非 HTTP 状态码）    |
| message | string | 是   | 响应消息                       |
| data    | any    | 否   | 响应数据，成功时返回，失败时省略 |

## 状态码定义

### 通用状态码

| 状态码 | 说明           | HTTP 状态码 |
| ------ | -------------- | ----------- |
| 200    | 成功           | 200         |
| 400    | 错误的请求     | 400         |
| 401    | 未授权         | 401         |
| 403    | 禁止访问       | 403         |
| 404    | 未找到         | 404         |
| 500    | 内部错误       | 500         |
| 10001  | 参数验证错误   | 400         |

### 用户相关状态码 (20xxx)

| 状态码 | 说明           |
| ------ | -------------- |
| 20001  | 用户不存在     |
| 20002  | 用户已存在     |
| 20003  | 密码错误       |
| 20004  | 用户已禁用     |
| 20005  | Token 无效     |
| 20006  | Token 过期     |
| 20007  | 权限不足       |

### 数据库相关状态码 (30xxx)

| 状态码 | 说明           |
| ------ | -------------- |
| 30001  | 数据库错误     |
| 30002  | 记录不存在     |
| 30003  | 记录已存在     |

### 缓存相关状态码 (40xxx)

| 状态码 | 说明       |
| ------ | ---------- |
| 40001  | 缓存错误   |

### 消息队列相关状态码 (50xxx)

| 状态码 | 说明        |
| ------ | ----------- |
| 50001  | Kafka 错误  |

### 第三方服务相关状态码 (60xxx)

| 状态码 | 说明               |
| ------ | ------------------ |
| 60001  | 第三方服务错误     |

## 响应示例

### 1. 成功响应

#### 普通成功

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com"
  }
}
```

#### 创建成功

```json
{
  "code": 200,
  "message": "User registered successfully",
  "data": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "status": 1,
    "created_at": "2023-01-01T00:00:00Z"
  }
}
```

HTTP 状态码：`201 Created`

#### 分页数据

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "username": "user1",
        "email": "user1@example.com"
      },
      {
        "id": 2,
        "username": "user2",
        "email": "user2@example.com"
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 10
  }
}
```

### 2. 错误响应

#### 参数验证错误

```json
{
  "code": 10001,
  "message": "Key: 'RegisterRequest.Username' Error:Field validation for 'Username' failed on the 'required' tag"
}
```

HTTP 状态码：`400 Bad Request`

#### 业务错误 - 用户已存在

```json
{
  "code": 20002,
  "message": "username already exists"
}
```

HTTP 状态码：`200 OK`

#### 未授权

```json
{
  "code": 401,
  "message": "invalid username or password"
}
```

HTTP 状态码：`401 Unauthorized`

#### 资源不存在

```json
{
  "code": 404,
  "message": "user not found"
}
```

HTTP 状态码：`404 Not Found`

#### 内部错误

```json
{
  "code": 500,
  "message": "internal server error"
}
```

HTTP 状态码：`500 Internal Server Error`

## 使用方法

### 在 Handler 中使用

```go
package handler

import (
    "trx-project/pkg/response"
    "github.com/gin-gonic/gin"
)

// 成功响应
func (h *Handler) GetUser(c *gin.Context) {
    user := getUserFromDB()
    response.Success(c, user)
}

// 成功响应（自定义消息）
func (h *Handler) CreateUser(c *gin.Context) {
    user := createUser()
    response.SuccessWithMsg(c, "User created successfully", user)
}

// 创建成功响应
func (h *Handler) Register(c *gin.Context) {
    user := register()
    response.Created(c, user)
    // 或
    response.CreatedWithMsg(c, "Registered successfully", user)
}

// 参数验证错误
func (h *Handler) ValidateError(c *gin.Context) {
    response.ValidateError(c, "Invalid parameters")
}

// 业务错误
func (h *Handler) BusinessError(c *gin.Context) {
    response.BusinessError(c, response.CodeUserAlreadyExists, "User already exists")
}

// 错误请求
func (h *Handler) BadRequest(c *gin.Context) {
    response.BadRequest(c, "Invalid request")
}

// 未授权
func (h *Handler) Unauthorized(c *gin.Context) {
    response.Unauthorized(c, "Invalid credentials")
}

// 禁止访问
func (h *Handler) Forbidden(c *gin.Context) {
    response.Forbidden(c, "Access denied")
}

// 未找到
func (h *Handler) NotFound(c *gin.Context) {
    response.NotFound(c, "Resource not found")
}

// 内部错误
func (h *Handler) InternalError(c *gin.Context) {
    response.InternalError(c, "Something went wrong")
}

// 分页响应
func (h *Handler) ListUsers(c *gin.Context) {
    users, total, page, pageSize := getUsers()
    response.PageSuccess(c, users, total, page, pageSize)
}

// 自定义错误响应
func (h *Handler) CustomError(c *gin.Context) {
    response.Error(c, 400, response.CodeValidateError, "Custom error message")
}
```

## 辅助函数

### response.Success(c, data)
返回成功响应，默认消息 "success"

### response.SuccessWithMsg(c, message, data)
返回成功响应，自定义消息

### response.Created(c, data)
返回创建成功响应（HTTP 201），默认消息 "created successfully"

### response.CreatedWithMsg(c, message, data)
返回创建成功响应（HTTP 201），自定义消息

### response.BadRequest(c, message)
返回错误请求响应（HTTP 400）

### response.Unauthorized(c, message)
返回未授权响应（HTTP 401）

### response.Forbidden(c, message)
返回禁止访问响应（HTTP 403）

### response.NotFound(c, message)
返回未找到响应（HTTP 404）

### response.InternalError(c, message)
返回内部错误响应（HTTP 500）

### response.ValidateError(c, message)
返回参数验证错误响应（HTTP 400，code 10001）

### response.BusinessError(c, code, message)
返回业务错误响应（HTTP 200，自定义 code）

### response.PageSuccess(c, list, total, page, pageSize)
返回分页数据响应

### response.Error(c, httpCode, code, message)
自定义错误响应

## 添加新的状态码

1. 在 `pkg/response/code.go` 中定义新的状态码常量：

```go
const (
    CodeOrderNotFound = 30001 // 订单不存在
    CodeOrderCanceled = 30002 // 订单已取消
)
```

2. 在 `CodeMessage` map 中添加对应的默认消息：

```go
var CodeMessage = map[int]string{
    CodeOrderNotFound: "order not found",
    CodeOrderCanceled: "order canceled",
}
```

3. 使用新的状态码：

```go
if order == nil {
    response.BusinessError(c, response.CodeOrderNotFound, "Order not found")
    return
}
```

## 最佳实践

### 1. 业务错误 vs HTTP 错误

- **业务错误**：使用 `BusinessError`，HTTP 状态码 200，自定义 code
  - 例如：用户已存在、余额不足、订单已取消
  
- **HTTP 错误**：使用对应的 HTTP 错误函数
  - 400: 参数错误
  - 401: 未登录/认证失败
  - 403: 权限不足
  - 404: 资源不存在
  - 500: 服务器内部错误

### 2. 错误消息

- 使用清晰、友好的错误消息
- 不要暴露敏感的系统信息
- 根据环境决定错误详细程度（开发环境可以详细，生产环境简洁）

### 3. 日志记录

在返回错误响应前，记录详细的错误日志：

```go
if err != nil {
    h.logger.Error("Failed to create user", 
        zap.Error(err),
        zap.String("username", username))
    response.InternalError(c, "Failed to create user")
    return
}
```

### 4. 数据返回

- 成功时返回必要的数据
- 敏感信息（如密码）不要返回
- 使用 JSON tag 控制字段序列化：`json:"-"` 不序列化

```go
type User struct {
    ID       uint   `json:"id"`
    Username string `json:"username"`
    Password string `json:"-"` // 不序列化到 JSON
}
```

## 前端处理建议

### Axios 拦截器示例

```javascript
// 请求拦截器
axios.interceptors.response.use(
  response => {
    const { code, message, data } = response.data
    
    if (code === 200) {
      return data
    } else if (code === 401) {
      // 跳转到登录页
      router.push('/login')
      return Promise.reject(new Error(message))
    } else {
      // 显示错误消息
      Message.error(message)
      return Promise.reject(new Error(message))
    }
  },
  error => {
    // HTTP 错误
    Message.error('网络错误，请稍后重试')
    return Promise.reject(error)
  }
)
```

### TypeScript 类型定义

```typescript
interface ApiResponse<T = any> {
  code: number
  message: string
  data?: T
}

interface PageData<T> {
  list: T[]
  total: number
  page: number
  page_size: number
}

// 使用示例
const response: ApiResponse<User> = await api.getUser(1)
const listResponse: ApiResponse<PageData<User>> = await api.getUserList()
```

## 测试

### 单元测试示例

```go
func TestSuccessResponse(t *testing.T) {
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    
    response.Success(c, gin.H{"key": "value"})
    
    assert.Equal(t, 200, w.Code)
    
    var resp response.Response
    json.Unmarshal(w.Body.Bytes(), &resp)
    
    assert.Equal(t, 200, resp.Code)
    assert.Equal(t, "success", resp.Message)
}
```

## 版本兼容性

如果需要修改响应格式，请：

1. 考虑向后兼容性
2. 使用 API 版本控制（如 `/api/v2/`）
3. 提前通知前端开发者
4. 保留旧版本 API 一段时间

## 参考资源

- [pkg/response/response.go](../pkg/response/response.go) - 响应函数实现
- [pkg/response/code.go](../pkg/response/code.go) - 状态码定义
- [internal/api/handler/user_handler.go](../internal/api/handler/user_handler.go) - 使用示例

