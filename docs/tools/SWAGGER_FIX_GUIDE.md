# Swagger 配置修复指南

## 🔍 问题诊断

如果 Swagger UI 无法加载 JSON 文件（显示空白或 "Failed to load API definition"），通常是以下原因：

### 1. 缺少 Swagger Docs 导入

**症状：**
- 浏览器访问 `/swagger/index.html` 显示空白
- 控制台报错：Failed to load spec

**原因：**
主程序没有导入生成的 Swagger 文档包

**解决方案：**

在 `cmd/frontend/main.go` 和 `cmd/backend/main.go` 中添加：

```go
import (
    // ... 其他导入 ...
    
    _ "trx-project/cmd/frontend/docs" // Swagger 文档 (前端)
    // 或
    _ "trx-project/cmd/backend/docs"  // Swagger 文档 (后端)
)
```

⚠️ **注意**：使用 `_` 空白标识符导入，因为我们只需要包的 `init()` 函数执行。

---

### 2. Handler 未导入 Model 包

**症状：**
- `swag init` 时报错：`cannot find type definition: model.User`
- 生成的 `swagger.json` 中 `paths` 为空 `{}`

**原因：**
Swagger 注释中引用了 `model.User` 等类型，但 Handler 文件没有导入 `model` 包

**解决方案：**

在所有使用了 `model.*` 类型的 Handler 文件中添加导入：

```go
package handler

import (
    "strconv"
    "trx-project/internal/model"    // 添加这行
    "trx-project/internal/service"
    "trx-project/pkg/response"
    // ...
)
```

---

### 3. Swagger 文档未重新生成

**症状：**
- 修改了 API 注释，但 Swagger UI 没有更新
- 新增的 API 接口没有出现在文档中

**原因：**
Swagger 文档是在编译时生成的静态文件，需要重新生成

**解决方案：**

```bash
# 重新生成前端文档
swag init -g cmd/frontend/main.go -o cmd/frontend/docs --parseDependency --parseInternal

# 重新生成后端文档
swag init -g cmd/backend/main.go -o cmd/backend/docs --parseDependency --parseInternal

# 或使用 Makefile
make swag
```

**重要参数说明：**
- `-g`: 指定主入口文件（包含 Swagger 主注释的文件）
- `-o`: 指定输出目录
- `--parseDependency`: 解析外部依赖
- `--parseInternal`: 解析内部包

---

### 4. 路由配置错误

**症状：**
- 访问 `/swagger/index.html` 返回 404

**检查路由配置：**

```go
// internal/api/router/frontend.go 或 backend.go

import (
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
)

// Swagger 文档路由
r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
```

---

## ✅ 完整修复步骤

### 步骤 1：确保导入了 Docs 包

**前端 (`cmd/frontend/main.go`)：**
```go
import (
    // ... 
    _ "trx-project/cmd/frontend/docs"
    // ...
)
```

**后端 (`cmd/backend/main.go`)：**
```go
import (
    // ...
    _ "trx-project/cmd/backend/docs"
    // ...
)
```

### 步骤 2：确保 Handler 导入了 Model

**所有引用了 `model.*` 的 Handler 文件：**
```go
import (
    "trx-project/internal/model"
    // ...
)
```

例如：
- `internal/api/handler/user_handler.go`
- `internal/api/handler/admin_user_handler.go`
- `internal/api/handler/rbac_handler.go`

### 步骤 3：重新生成 Swagger 文档

```bash
# 从项目根目录执行
cd /path/to/trx-project

# 生成前端文档
swag init -g cmd/frontend/main.go -o cmd/frontend/docs --parseDependency --parseInternal

# 生成后端文档
swag init -g cmd/backend/main.go -o cmd/backend/docs --parseDependency --parseInternal
```

**验证生成是否成功：**
```bash
# 检查是否包含 API 路径
head -50 cmd/frontend/docs/swagger.json
head -50 cmd/backend/docs/swagger.json

# 应该看到类似以下内容（而不是空的 paths: {}）:
# "paths": {
#     "/api/v1/public/register": {
#         ...
#     }
# }
```

### 步骤 4：清理并重新编译

```bash
# 清理旧的编译文件
rm -rf tmp
mkdir tmp

# 重新启动服务
./start-frontend.sh  # 或 .\start-frontend.ps1 (Windows)
```

### 步骤 5：验证修复

打开浏览器访问：

**前端：**
- http://localhost:8080/swagger/index.html

**后端：**
- http://localhost:8081/swagger/index.html

**预期结果：**
- ✅ 页面正常加载
- ✅ 显示 API 标题和版本
- ✅ 显示所有 API 端点
- ✅ 可以展开查看详细信息
- ✅ 可以点击 "Try it out" 测试 API

---

## 📋 常见错误及解决方案

### 错误 1: "Failed to load API definition"

**原因**：主程序未导入 docs 包

**解决**：
```go
import _ "trx-project/cmd/frontend/docs"
```

### 错误 2: swagger.json 中 paths 为空

**原因**：
1. Handler 未导入 model 包
2. swag init 命令执行目录不对

**解决**：
```bash
# 确保在项目根目录执行
cd /path/to/trx-project
swag init -g cmd/frontend/main.go -o cmd/frontend/docs --parseDependency --parseInternal
```

### 错误 3: cannot find type definition

**swag init 报错示例：**
```
ParseComment error: cannot find type definition: model.User
```

**解决**：在对应的 Handler 文件中添加：
```go
import "trx-project/internal/model"
```

### 错误 4: 修改注释后文档没有更新

**原因**：没有重新生成 Swagger 文档

**解决**：
```bash
make swag  # 或手动执行 swag init
rm -rf tmp  # 清理编译缓存
./start-frontend.sh  # 重启服务
```

---

## 🔧 开发工作流

### 添加新的 API

1. **编写 Handler 代码并添加 Swagger 注释**
   ```go
   // GetUser 获取用户信息
   // @Summary 获取用户信息
   // @Description 根据用户ID获取用户详细信息
   // @Tags 用户接口
   // @Security BearerAuth
   // @Param id path int true "用户ID"
   // @Success 200 {object} response.Response{data=model.User}
   // @Router /users/{id} [get]
   func (h *UserHandler) GetUser(c *gin.Context) {
       // ...
   }
   ```

2. **重新生成文档**
   ```bash
   make swag
   ```

3. **重启服务**
   ```bash
   # 按 Ctrl+C 停止，然后重新启动
   ./start-frontend.sh
   ```

4. **验证**
   - 访问 http://localhost:8080/swagger/index.html
   - 确认新 API 出现在文档中

---

## 📚 Swagger 注释参考

### 主配置注释 (在 cmd/*/docs.go)

```go
// @title 项目名称
// @version 1.0
// @description 项目描述
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description 输入 "Bearer " + JWT Token
```

### API 端点注释

```go
// FunctionName 函数描述
// @Summary 简短摘要
// @Description 详细描述
// @Tags 标签名称
// @Accept json
// @Produce json
// @Security BearerAuth          // 需要认证的接口添加
// @Param name type dataType required "描述"
// @Success 200 {object} Type "成功描述"
// @Failure 400 {object} Type "失败描述"
// @Router /path [method]
```

### 参数类型

| Swagger 参数 | 含义 | 示例 |
|-------------|------|------|
| `path` | 路径参数 | `/users/{id}` 中的 `id` |
| `query` | 查询参数 | `?page=1` 中的 `page` |
| `body` | 请求体 | JSON 请求体 |
| `header` | 请求头 | `Authorization` |

---

## 🎯 最佳实践

1. **每次修改 API 注释后都重新生成文档**
   ```bash
   make swag
   ```

2. **使用有意义的标签分组**
   ```go
   // @Tags 公开接口
   // @Tags 用户管理
   // @Tags RBAC管理
   ```

3. **提供详细的参数说明和示例**
   ```go
   // @Param page query int false "页码，默认1" default(1)
   // @Param page_size query int false "每页数量，默认10" default(10)
   ```

4. **使用中文注释提高可读性**
   ```go
   // @Summary 用户注册
   // @Description 创建新用户账号，注册成功后返回用户信息和 JWT Token
   ```

5. **定期验证文档完整性**
   - 检查所有 API 是否都有文档
   - 确保参数描述准确
   - 测试示例请求是否有效

---

## 🚀 自动化脚本

### 一键重新生成所有文档

创建 `regenerate-swagger.sh`：

```bash
#!/bin/bash

echo "🔄 重新生成 Swagger 文档..."

# 生成前端文档
echo "📝 生成前端文档..."
swag init -g cmd/frontend/main.go -o cmd/frontend/docs --parseDependency --parseInternal

# 生成后端文档
echo "📝 生成后端文档..."
swag init -g cmd/backend/main.go -o cmd/backend/docs --parseDependency --parseInternal

echo "✅ Swagger 文档生成完成！"
echo ""
echo "📚 访问地址："
echo "  前端: http://localhost:8080/swagger/index.html"
echo "  后端: http://localhost:8081/swagger/index.html"
```

使用：
```bash
chmod +x regenerate-swagger.sh
./regenerate-swagger.sh
```

---

## 📖 相关资源

- [Swagger 官方文档](https://swagger.io/docs/)
- [swag 项目地址](https://github.com/swaggo/swag)
- [gin-swagger 文档](https://github.com/swaggo/gin-swagger)
- [OpenAPI 规范](https://spec.openapis.org/oas/v2.0)

---

## 💡 提示

- Swagger 文档是**编译时生成**的，不是运行时动态生成
- 修改注释后必须重新运行 `swag init`
- Air 热重载会自动重新编译，但你仍需要先运行 `swag init`
- 可以在 `.air-*.toml` 的 `cmd` 中添加 `swag init` 来实现自动化

---

**最后更新：** 2025-11-12

