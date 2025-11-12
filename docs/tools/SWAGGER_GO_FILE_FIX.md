# Swagger Go 文件过滤问题修复

## 问题描述

用户反馈：虽然前台的 `frontend_swagger.json` 文件中已经过滤掉了 `/admin` 接口，但通过浏览器访问 `http://localhost:8080/swagger/index.html` 仍然能看到这些后台接口。

## 根本原因

`swag init` 命令会生成 **三个文件**：

1. **`{instance}_swagger.json`** - JSON 格式的文档
2. **`{instance}_swagger.yaml`** - YAML 格式的文档  
3. **`{instance}_docs.go`** - **Go 代码格式的文档** ⭐

关键问题：**Gin Swagger 实际上是从 Go 文件（`.go`）中读取文档内容的**，而不是从 JSON 文件读取！

### 文件结构

```go
// cmd/frontend/docs/frontend_docs.go
package docs

import "github.com/swaggo/swag"

const docTemplatefrontend = `{
    "swagger": "2.0",
    "paths": {
        "/admin/rbac/permissions": {  // ← 这里包含了后台接口！
            ...
        },
        "/public/login": {
            ...
        }
    }
}`
```

### 之前的过滤脚本

之前的 `filter_swagger.go` 只过滤了 JSON 文件：

```bash
go run scripts/filter_swagger.go \
    cmd/frontend/docs/frontend_swagger.json \  # 只过滤 JSON
    cmd/frontend/docs/frontend_swagger.json \
    admin users
```

**结果**：
- ✅ JSON 文件被正确过滤
- ❌ Go 文件未被过滤
- ❌ Swagger UI 仍显示所有接口（因为它读取的是 Go 文件）

## 解决方案

### 1. 创建新的过滤工具

创建了 `scripts/filter_swagger_docs.go`，它能够：

1. 过滤 JSON 文件（移除不需要的接口）
2. **过滤 Go 文件**（用过滤后的 JSON 替换 Go 文件中的模板内容）

### 2. 工作原理

```go
func filterGoFile(goFile, jsonFile, instanceName string) error {
    // 1. 读取过滤后的 JSON 内容
    jsonData, err := os.ReadFile(jsonFile)
    
    // 2. 读取原始 Go 文件
    goData, err := os.ReadFile(goFile)
    goContent := string(goData)
    
    // 3. 找到 Go 文件中的模板定义
    //    格式: const docTemplate{instanceName} = `...`
    templateStart := "const docTemplate" + instanceName + " = `"
    startIdx := strings.Index(goContent, templateStart)
    contentStart := startIdx + len(templateStart)
    
    // 4. 找到模板结束位置（结束的 `）
    endIdx := strings.Index(goContent[contentStart:], "`")
    contentEnd := contentStart + endIdx
    
    // 5. 替换模板内容
    newGoContent := goContent[:contentStart] + jsonStr + goContent[contentEnd:]
    
    // 6. 写回 Go 文件
    os.WriteFile(goFile, []byte(newGoContent), 0644)
}
```

### 3. 新的使用方式

```bash
# 使用新工具，同时过滤 JSON 和 Go 文件
go run scripts/filter_swagger_docs.go \
    cmd/frontend/docs \  # 文档目录
    frontend \           # 实例名称
    admin users          # 要排除的路径前缀
```

## 修复效果

### 前台文档（frontend）

**过滤前**:
```json
{
    "paths": {
        "/admin/rbac/permissions": {...},  // ❌ 不应该有
        "/admin/users": {...},             // ❌ 不应该有
        "/public/login": {...},            // ✅
        "/user/profile": {...}             // ✅
    }
}
```

**过滤后**:
```json
{
    "paths": {
        "/public/login": {...},    // ✅ 只有前台接口
        "/user/profile": {...}     // ✅
    }
}
```

- ✅ 移除了 14 个 `/admin` 和 `/users` 开头的接口
- ✅ 保留了 3 个前台接口
- ✅ **JSON 和 Go 文件都被正确过滤**

### 后台文档（backend）

**过滤前**:
```json
{
    "paths": {
        "/admin/rbac/permissions": {...},  // ✅
        "/admin/users": {...},             // ✅
        "/public/login": {...},            // ❌ 不应该有
        "/user/profile": {...}             // ❌ 不应该有
    }
}
```

**过滤后**:
```json
{
    "paths": {
        "/admin/rbac/permissions": {...},  // ✅ 只有后台接口
        "/admin/users": {...},             // ✅
        ...
    }
}
```

- ✅ 移除了 5 个 `/public` 和 `/user` 开头的接口
- ✅ 保留了 12 个后台接口
- ✅ **JSON 和 Go 文件都被正确过滤**

## 验证测试

### 1. 编译测试

```bash
$ go build ./cmd/frontend
✅ 编译成功

$ go build ./cmd/backend
✅ 编译成功
```

### 2. 文档检查

```bash
# 前台 Go 文件中不应有 /admin 接口
$ grep "/admin" cmd/frontend/docs/frontend_docs.go
No matches found ✅

# 后台 Go 文件中不应有 /public 接口
$ grep "/public" cmd/backend/docs/backend_docs.go
No matches found ✅
```

### 3. Swagger UI 验证

启动服务后，访问 Swagger UI：

**前台** (`http://localhost:8080/swagger/index.html`):
- ✅ 只显示 3 个接口：
  - `POST /public/login`
  - `POST /public/register`
  - `GET /user/profile`
- ❌ 不再显示任何 `/admin/*` 接口

**后台** (`http://localhost:8081/swagger/index.html`):
- ✅ 只显示 12 个接口：
  - `/admin/rbac/*` 系列
  - `/admin/users/*` 系列
  - `/admin/statistics/*` 系列
- ❌ 不再显示任何 `/public/*` 或 `/user/*` 接口

## 更新的文件

### 新增文件

1. **`scripts/filter_swagger_docs.go`** - 新的过滤工具，同时处理 JSON 和 Go 文件
2. **`docs/SWAGGER_GO_FILE_FIX.md`** - 本文档

### 更新的文件

1. **`scripts/swagger.sh`** - 使用新的过滤工具
2. **`scripts/swagger.bat`** - 使用新的过滤工具
3. **`scripts/swagger.ps1`** - 使用新的过滤工具
4. **`Makefile`** - 使用新的过滤工具

### 更新内容示例

**Bash 脚本** (`swagger.sh`):
```bash
# 旧版本
go run scripts/filter_swagger.go \
    cmd/frontend/docs/frontend_swagger.json \
    cmd/frontend/docs/frontend_swagger.json \
    admin users

# 新版本 ✅
go run scripts/filter_swagger_docs.go \
    cmd/frontend/docs \
    frontend \
    admin users
```

**Makefile**:
```makefile
# 旧版本
@go run scripts/filter_swagger.go \
    cmd/frontend/docs/frontend_swagger.json \
    cmd/frontend/docs/frontend_swagger.json \
    admin users

# 新版本 ✅
@go run scripts/filter_swagger_docs.go \
    cmd/frontend/docs \
    frontend \
    admin users
```

## 使用说明

### 重新生成文档

如果你已经运行过旧版本的脚本，需要重新生成文档：

```bash
# 使用 Makefile
make swag-frontend
make swag-backend

# 或使用脚本
./scripts/swagger.sh frontend
./scripts/swagger.sh backend

# 或生成所有文档
./scripts/swagger.sh
make swag
```

### 验证修复

1. **重新生成文档**（使用上面的命令）

2. **启动服务**:
   ```bash
   go run cmd/frontend/main.go  # 端口 8080
   go run cmd/backend/main.go   # 端口 8081
   ```

3. **访问 Swagger UI**:
   - 前台：http://localhost:8080/swagger/index.html
   - 后台：http://localhost:8081/swagger/index.html

4. **检查接口列表**:
   - ✅ 前台只显示前台接口
   - ✅ 后台只显示后台接口

## 技术细节

### Swagger 文档加载流程

```
1. 应用启动
   ↓
2. 导入 docs 包
   import _ "trx-project/cmd/frontend/docs"
   ↓
3. 注册 Swagger 文档
   docs.SwaggerInfo{instanceName} 读取 docTemplate{instanceName}
   ↓
4. Swagger UI 请求文档
   Gin Swagger 返回从 Go 文件读取的内容
   ↓
5. 浏览器显示 Swagger UI
```

### 为什么是 Go 文件而不是 JSON？

1. **编译时嵌入**: Go 文件在编译时被嵌入到二进制文件中
2. **无需额外文件**: 部署时不需要携带 JSON 文件
3. **性能更好**: 直接从内存读取，无需文件 I/O

### 旧工具 vs 新工具对比

| 特性 | filter_swagger.go | filter_swagger_docs.go |
|------|-------------------|------------------------|
| 过滤 JSON | ✅ | ✅ |
| 过滤 Go 文件 | ❌ | ✅ |
| Swagger UI 生效 | ❌ | ✅ |
| 使用方式 | 指定 JSON 文件路径 | 指定文档目录和实例名 |

## 常见问题

### Q: 为什么之前没发现这个问题？

A: 因为在开发过程中，我们主要关注 JSON 文件的内容，而没有实际启动服务并访问 Swagger UI 进行验证。

### Q: 旧的 filter_swagger.go 还需要保留吗？

A: 可以保留作为备份，但推荐使用新的 `filter_swagger_docs.go`。

### Q: 如果只想过滤 JSON 怎么办？

A: 仍然可以使用旧工具：
```bash
go run scripts/filter_swagger.go <json-file> <output-file> <prefixes...>
```

### Q: 如果 Go 文件格式发生变化怎么办？

A: 新工具使用字符串查找和替换，只要 `const docTemplate{instanceName} = ` 这个模式不变，就能正常工作。

## 总结

✅ **问题根源**: Swagger UI 读取的是 Go 文件，而不是 JSON 文件  
✅ **解决方案**: 创建新工具，同时过滤 JSON 和 Go 文件  
✅ **修复验证**: 编译成功，Swagger UI 只显示对应服务的接口  
✅ **向后兼容**: 所有现有脚本和 Makefile 已更新  

现在，前后台的 Swagger 文档真正实现了完全分离！🎉

