# Swagger 脚本使用指南

## 概述

为了方便在不同平台上生成 Swagger 文档，我们提供了三个跨平台脚本：

- **`swagger.sh`** - Bash 脚本（Linux、macOS、Git Bash）
- **`swagger.bat`** - Windows 批处理脚本（CMD）
- **`swagger.ps1`** - PowerShell 脚本（Windows、Linux、macOS）

这些脚本提供了与 Makefile 相同的功能，但更易于在没有 `make` 命令的环境中使用。

## 平台支持

### Linux / macOS

**推荐使用**: `swagger.sh`

```bash
# 添加执行权限（首次使用）
chmod +x scripts/swagger.sh

# 生成所有文档
./scripts/swagger.sh

# 只生成前台文档
./scripts/swagger.sh frontend

# 只生成后台文档
./scripts/swagger.sh backend

# 查看帮助
./scripts/swagger.sh help
```

**或使用 PowerShell** (如果已安装):

```bash
pwsh scripts/swagger.ps1
pwsh scripts/swagger.ps1 frontend
```

### Windows

#### 方式 1: Git Bash（推荐）

```bash
# 在 Git Bash 中运行
bash scripts/swagger.sh
bash scripts/swagger.sh frontend
bash scripts/swagger.sh backend
```

#### 方式 2: CMD（命令提示符）

```cmd
REM 生成所有文档
scripts\swagger.bat

REM 只生成前台文档
scripts\swagger.bat frontend

REM 只生成后台文档
scripts\swagger.bat backend

REM 查看帮助
scripts\swagger.bat help
```

#### 方式 3: PowerShell

```powershell
# 生成所有文档
.\scripts\swagger.ps1

# 只生成前台文档
.\scripts\swagger.ps1 frontend

# 只生成后台文档
.\scripts\swagger.ps1 backend

# 查看帮助
.\scripts\swagger.ps1 help
```

**注意**: 首次运行 PowerShell 脚本可能需要修改执行策略：

```powershell
# 允许运行本地脚本
Set-ExecutionPolicy RemoteSigned -Scope CurrentUser
```

## 功能特性

### ✅ 自动检查依赖

所有脚本会自动检查 `swag` 是否已安装：

```bash
✅ swag 已安装: swag version v1.16.4
```

如果未安装，会提示安装命令：

```bash
❌ swag 未安装！

请运行以下命令安装：
  go install github.com/swaggo/swag/cmd/swag@latest
```

### 📝 详细的进度输出

生成过程中会显示清晰的状态信息：

```
🚀 Swagger 文档生成器
================================

✅ swag 已安装: swag version v1.16.4

ℹ️  生成前台 Swagger 文档...
ℹ️  过滤后台接口 (/admin, /users)...
📊 Filtered 14 paths with prefixes [/admin /users]
✅ 过滤成功！
✅ 前台文档生成完成！
ℹ️  文档位置: cmd/frontend/docs/frontend_swagger.json
ℹ️  访问地址: http://localhost:8080/swagger/index.html

ℹ️  生成后台 Swagger 文档...
ℹ️  过滤前台接口 (/public, /user)...
📊 Filtered 5 paths with prefixes [/public /user]
✅ 过滤成功！
✅ 后台文档生成完成！
ℹ️  文档位置: cmd/backend/docs/backend_swagger.json
ℹ️  访问地址: http://localhost:8081/swagger/index.html

================================
✅ 完成！
```

### 🎯 选择性生成

支持三种生成模式：

1. **只生成前台文档** - 适合只修改了前台接口
2. **只生成后台文档** - 适合只修改了后台接口
3. **生成所有文档** - 完整生成（默认）

### 🧹 自动清理

脚本会自动清理旧的文档文件：

- `docs.go` - 旧格式的文档
- `swagger.json` - 未使用实例名的文档
- `swagger.yaml` - 未使用实例名的文档

### 🔄 自动过滤

- **前台文档**: 自动移除 `/admin` 和 `/users` 接口
- **后台文档**: 自动移除 `/public` 和 `/user` 接口

## 脚本对比

| 特性 | swagger.sh | swagger.bat | swagger.ps1 | Makefile |
|------|------------|-------------|-------------|----------|
| Linux 支持 | ✅ | ❌ | ✅* | ✅ |
| macOS 支持 | ✅ | ❌ | ✅* | ✅ |
| Windows CMD | ❌ | ✅ | ❌ | ❌** |
| Windows PowerShell | ❌ | ❌ | ✅ | ❌** |
| Git Bash (Windows) | ✅ | ❌ | ❌ | ❌** |
| 彩色输出 | ✅ | ❌ | ✅ | ✅ |
| 依赖检查 | ✅ | ✅ | ✅ | ❌ |
| 错误处理 | ✅ | ✅ | ✅ | ✅ |

\* 需要安装 PowerShell Core  
\** 需要安装 make 工具

## 工作流程

所有脚本执行相同的流程：

```
1. 检查 swag 是否安装
   ↓
2. 运行 swag init 生成原始文档
   ↓
3. 调用 filter_swagger.go 过滤接口
   ↓
4. 清理旧的文档文件
   ↓
5. 显示完成信息
```

## 常见问题

### Q: 哪个脚本最适合我？

**Linux/macOS 用户**:
- 使用 `swagger.sh`（Bash 脚本）

**Windows 用户**:
- **Git Bash**: 使用 `swagger.sh`（最推荐）
- **CMD**: 使用 `swagger.bat`
- **PowerShell**: 使用 `swagger.ps1`

### Q: 脚本和 Makefile 有什么区别？

**功能**: 完全相同，都执行相同的文档生成流程

**优势对比**:
- **Makefile**: 
  - ✅ 开发者熟悉
  - ✅ 可以定义其他构建任务
  - ❌ Windows 上需要额外安装 make
  
- **Shell 脚本**: 
  - ✅ 跨平台兼容性好
  - ✅ 不需要额外工具
  - ✅ 更好的依赖检查和错误提示

### Q: PowerShell 提示无法运行脚本？

这是 Windows 的安全策略。运行以下命令：

```powershell
Set-ExecutionPolicy RemoteSigned -Scope CurrentUser
```

### Q: Git Bash 中批处理脚本显示乱码？

这是正常的。在 Git Bash 中请使用 `swagger.sh`，在 CMD 中使用 `swagger.bat`。

### Q: 如何只重新生成某个服务的文档？

```bash
# 只更新前台
./scripts/swagger.sh frontend

# 只更新后台
./scripts/swagger.sh backend
```

这样可以节省时间，特别是在大型项目中。

### Q: 脚本失败了怎么办？

脚本会显示详细的错误信息：

```bash
❌ swag 未安装！
❌ 前台文档生成失败！
❌ 过滤失败！
```

根据错误提示：
1. 检查 swag 是否正确安装
2. 检查代码中的 Swagger 注释是否正确
3. 查看详细的错误输出

### Q: 可以在 CI/CD 中使用这些脚本吗？

可以！脚本在失败时会返回非零退出码，适合 CI/CD 使用：

```yaml
# GitHub Actions 示例
- name: Generate Swagger Docs
  run: bash scripts/swagger.sh

# GitLab CI 示例
swagger:
  script:
    - bash scripts/swagger.sh
```

## 示例

### 开发工作流

```bash
# 1. 修改前台接口的 Handler
vim internal/api/handler/user_handler.go

# 2. 只重新生成前台文档
./scripts/swagger.sh frontend

# 3. 启动前台服务查看文档
go run cmd/frontend/main.go

# 4. 访问 http://localhost:8080/swagger/index.html
```

### 发布前检查

```bash
# 生成所有文档
./scripts/swagger.sh

# 验证编译
go build ./cmd/frontend
go build ./cmd/backend

# 提交代码
git add .
git commit -m "docs: update swagger documentation"
```

### CI/CD 集成

```bash
# .gitlab-ci.yml
test:swagger:
  stage: test
  script:
    - go install github.com/swaggo/swag/cmd/swag@latest
    - bash scripts/swagger.sh
    - test -f cmd/frontend/docs/frontend_swagger.json
    - test -f cmd/backend/docs/backend_swagger.json
```

## 脚本源码

- **Bash 脚本**: `scripts/swagger.sh`
- **批处理脚本**: `scripts/swagger.bat`
- **PowerShell 脚本**: `scripts/swagger.ps1`
- **过滤脚本**: `scripts/filter_swagger.go`

## 相关文档

- [Swagger 文档使用指南](SWAGGER_GUIDE.md)
- [Swagger 分离实现总结](SWAGGER_SEPARATION_SUMMARY.md)
- [项目 README](../README.md)

## 总结

通过提供多个跨平台脚本，我们确保了在任何环境下都能轻松生成 Swagger 文档：

✅ **跨平台** - Linux、macOS、Windows 全支持  
✅ **易用性** - 简单的命令行参数  
✅ **自动化** - 包含生成、过滤、清理全流程  
✅ **友好提示** - 彩色输出和详细信息  
✅ **错误处理** - 清晰的错误提示和退出码

选择适合你环境的脚本，享受便捷的文档生成体验！

