# API 统一响应格式已实现 ✅

## 改动说明

已成功为项目实现统一的 API 响应格式，主要改动包括：

### 1. 新增响应处理包

**文件**: `pkg/response/response.go`

提供了一系列统一的响应辅助函数：

- ✅ `Success(c, data)` - 成功响应
- ✅ `SuccessWithMsg(c, message, data)` - 成功响应（自定义消息）
- ✅ `Created(c, data)` - 创建成功响应（HTTP 201）
- ✅ `CreatedWithMsg(c, message, data)` - 创建成功响应（自定义消息）
- ✅ `BadRequest(c, message)` - 错误请求（HTTP 400）
- ✅ `Unauthorized(c, message)` - 未授权（HTTP 401）
- ✅ `Forbidden(c, message)` - 禁止访问（HTTP 403）
- ✅ `NotFound(c, message)` - 未找到（HTTP 404）
- ✅ `InternalError(c, message)` - 内部错误（HTTP 500）
- ✅ `ValidateError(c, message)` - 参数验证错误
- ✅ `BusinessError(c, code, message)` - 业务错误
- ✅ `PageSuccess(c, list, total, page, pageSize)` - 分页响应

### 2. 新增状态码定义

**文件**: `pkg/response/code.go`

定义了完整的业务状态码体系：

- **通用状态码**: 200, 400, 401, 403, 404, 500, 10001
- **用户相关**: 20001-20007
- **数据库相关**: 30001-30003
- **缓存相关**: 40001
- **Kafka 相关**: 50001
- **第三方服务**: 60001

### 3. 更新 Handler 层

**文件**: `internal/api/handler/user_handler.go`

所有接口已统一使用新的响应格式：

#### 注册接口
```go
// 成功
response.CreatedWithMsg(c, "User registered successfully", user)

// 用户已存在
response.BusinessError(c, response.CodeUserAlreadyExists, err.Error())

// 参数验证失败
response.ValidateError(c, err.Error())
```

#### 登录接口
```go
// 成功
response.SuccessWithMsg(c, "Login successful", user)

// 认证失败
response.Unauthorized(c, err.Error())

// 用户已禁用
response.BusinessError(c, response.CodeUserDisabled, err.Error())
```

#### 获取用户
```go
// 成功
response.Success(c, user)

// 用户不存在
response.NotFound(c, err.Error())

// 参数错误
response.BadRequest(c, "Invalid user ID")
```

#### 用户列表
```go
// 分页响应
response.PageSuccess(c, users, total, page, pageSize)
```

#### 删除用户
```go
// 成功
response.SuccessWithMsg(c, "User deleted successfully", nil)

// 失败
response.InternalError(c, "Failed to delete user")
```

## 统一响应格式

### 成功响应

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "username": "testuser"
  }
}
```

### 分页响应

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [...],
    "total": 100,
    "page": 1,
    "page_size": 10
  }
}
```

### 错误响应

```json
{
  "code": 20002,
  "message": "username already exists"
}
```

## 使用示例

### 基本使用

```go
import "trx-project/pkg/response"

// 成功响应
func GetUser(c *gin.Context) {
    user := fetchUser()
    response.Success(c, user)
}

// 错误响应
func CreateUser(c *gin.Context) {
    if err := validate(); err != nil {
        response.ValidateError(c, err.Error())
        return
    }
    
    if userExists() {
        response.BusinessError(c, response.CodeUserAlreadyExists, "User already exists")
        return
    }
    
    user := createUser()
    response.Created(c, user)
}
```

### 分页使用

```go
func ListUsers(c *gin.Context) {
    page := getPage()
    pageSize := getPageSize()
    users, total := fetchUsers(page, pageSize)
    
    response.PageSuccess(c, users, total, page, pageSize)
}
```

## 状态码使用规范

### HTTP 状态码 vs 业务状态码

- **HTTP 状态码**: 表示 HTTP 请求的处理结果
  - 200: 成功
  - 201: 创建成功
  - 400: 请求错误
  - 401: 未认证
  - 403: 无权限
  - 404: 资源不存在
  - 500: 服务器错误

- **业务状态码** (code 字段): 表示业务逻辑的处理结果
  - 200: 业务成功
  - 20002: 用户已存在（业务错误，但 HTTP 返回 200）
  - 10001: 参数验证失败（HTTP 返回 400）

### 何时使用业务错误

使用 `BusinessError` 当：
- 业务规则不满足（如余额不足、订单已取消）
- 数据状态不允许操作（如已完成的订单不能修改）
- 业务约束冲突（如用户名已存在）

这些情况下 HTTP 状态码返回 200，但 code 字段使用特定的业务错误码。

### 何时使用 HTTP 错误

使用对应的 HTTP 错误函数当：
- 参数格式错误（`BadRequest`）
- 未登录或认证失败（`Unauthorized`）
- 权限不足（`Forbidden`）
- 资源不存在（`NotFound`）
- 服务器异常（`InternalError`）

## 新增状态码

如需添加新的业务状态码：

1. 在 `pkg/response/code.go` 中定义常量：

```go
const (
    CodeOrderNotFound = 30001 // 订单不存在
)
```

2. 添加默认消息：

```go
var CodeMessage = map[int]string{
    CodeOrderNotFound: "order not found",
}
```

3. 使用：

```go
response.BusinessError(c, response.CodeOrderNotFound, "Order not found")
```

## 文档

详细文档请查看：
- [docs/RESPONSE_FORMAT.md](docs/RESPONSE_FORMAT.md) - 完整的响应格式说明和使用指南

## 验证

项目已成功构建并验证：

```bash
✅ 构建成功！统一响应格式已实现。
```

测试命令：

```bash
# 构建项目
go build -o bin/server ./cmd/server

# 运行项目
./bin/server

# 测试 API
curl http://localhost:8080/api/v1/users
```

## 后续建议

1. ✅ **已完成**: 统一响应格式
2. ✅ **已完成**: 状态码定义
3. ✅ **已完成**: Handler 层改造
4. ⬜ **建议**: 添加请求追踪 ID (Trace ID)
5. ⬜ **建议**: 添加响应时间统计
6. ⬜ **建议**: 实现错误国际化
7. ⬜ **建议**: 添加更多业务错误码

## 优势

✅ **统一规范**: 所有接口遵循相同的响应格式  
✅ **易于维护**: 集中管理响应逻辑和状态码  
✅ **类型安全**: 使用常量定义状态码，避免硬编码  
✅ **易于扩展**: 可轻松添加新的状态码和响应类型  
✅ **前端友好**: 统一的响应格式便于前端处理  
✅ **文档完善**: 提供详细的使用说明和示例  

## 示例对比

### 改造前

```go
c.JSON(http.StatusOK, gin.H{
    "code": 200,
    "message": "success",
    "data": user,
})
```

### 改造后

```go
response.Success(c, user)
```

更简洁、更统一、更易维护！🎉

