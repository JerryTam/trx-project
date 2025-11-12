# API 文档

## 基础信息

- 基础 URL: `http://localhost:8080`
- Content-Type: `application/json`

## 健康检查

### 获取服务健康状态

```
GET /health
```

#### 响应示例

```json
{
  "status": "ok"
}
```

## 用户管理 API

### 1. 用户注册

**接口地址**: `/api/v1/users/register`

**请求方法**: `POST`

**请求参数**:

| 参数名   | 类型   | 必填 | 说明               |
| -------- | ------ | ---- | ------------------ |
| username | string | 是   | 用户名 (3-50字符)  |
| email    | string | 是   | 邮箱地址           |
| password | string | 是   | 密码 (最少6字符)   |

**请求示例**:

```bash
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'
```

**成功响应** (HTTP 201):

```json
{
  "code": 200,
  "message": "User registered successfully",
  "data": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "status": 1,
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z"
  }
}
```

### 2. 用户登录

**接口地址**: `/api/v1/users/login`

**请求方法**: `POST`

**请求参数**:

| 参数名   | 类型   | 必填 | 说明   |
| -------- | ------ | ---- | ------ |
| username | string | 是   | 用户名 |
| password | string | 是   | 密码   |

**请求示例**:

```bash
curl -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

**响应示例**:

```json
{
  "code": 200,
  "message": "Login successful",
  "data": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "status": 1,
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z"
  }
}
```

### 3. 获取用户信息

**接口地址**: `/api/v1/users/:id`

**请求方法**: `GET`

**路径参数**:

| 参数名 | 类型 | 必填 | 说明   |
| ------ | ---- | ---- | ------ |
| id     | int  | 是   | 用户ID |

**请求示例**:

```bash
curl -X GET http://localhost:8080/api/v1/users/1
```

**响应示例**:

```json
{
  "code": 200,
  "message": "Success",
  "data": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "status": 1,
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z"
  }
}
```

### 4. 获取用户列表

**接口地址**: `/api/v1/users`

**请求方法**: `GET`

**查询参数**:

| 参数名    | 类型 | 必填 | 默认值 | 说明       |
| --------- | ---- | ---- | ------ | ---------- |
| page      | int  | 否   | 1      | 页码       |
| page_size | int  | 否   | 10     | 每页数量   |

**请求示例**:

```bash
curl -X GET "http://localhost:8080/api/v1/users?page=1&page_size=10"
```

**响应示例**:

```json
{
  "code": 200,
  "message": "Success",
  "data": {
    "users": [
      {
        "id": 1,
        "username": "testuser",
        "email": "test@example.com",
        "status": 1,
        "created_at": "2023-01-01T00:00:00Z",
        "updated_at": "2023-01-01T00:00:00Z"
      }
    ],
    "total": 1,
    "page": 1
  }
}
```

### 5. 删除用户

**接口地址**: `/api/v1/users/:id`

**请求方法**: `DELETE`

**路径参数**:

| 参数名 | 类型 | 必填 | 说明   |
| ------ | ---- | ---- | ------ |
| id     | int  | 是   | 用户ID |

**请求示例**:

```bash
curl -X DELETE http://localhost:8080/api/v1/users/1
```

**响应示例**:

```json
{
  "code": 200,
  "message": "User deleted successfully"
}
```

## 错误码说明

| 错误码 | 说明           |
| ------ | -------------- |
| 200    | 成功           |
| 201    | 创建成功       |
| 400    | 请求参数错误   |
| 401    | 未授权         |
| 404    | 资源不存在     |
| 500    | 服务器内部错误 |

## 错误响应格式

```json
{
  "code": 400,
  "message": "错误信息描述"
}
```

