# Scripts 目录说明

本目录包含项目的各种实用脚本，用于开发、测试和部署。

## 📁 脚本分类

### Swagger 文档生成

| 脚本 | 平台 | 说明 |
|------|------|------|
| `swagger.sh` | Linux / macOS / Git Bash | Bash 脚本，生成 Swagger 文档 |
| `swagger.bat` | Windows CMD | 批处理脚本，生成 Swagger 文档 |
| `swagger.ps1` | PowerShell | PowerShell 脚本，生成 Swagger 文档 |
| `filter_swagger.go` | 跨平台 | 过滤 Swagger 文档中的接口 |

**使用示例**:
```bash
# Linux / macOS / Git Bash
./scripts/swagger.sh frontend

# Windows CMD
scripts\swagger.bat frontend

# PowerShell
.\scripts\swagger.ps1 frontend
```

**详细文档**: [Swagger 脚本使用指南](../docs/SWAGGER_SCRIPTS_GUIDE.md)

---

### 数据库迁移

| 脚本 | 说明 |
|------|------|
| `migrate.sh` | 数据库迁移管理（up/down/create/drop） |
| `test_migration.sh` | 测试数据库迁移功能 |
| `init_db.sql` | 数据库初始化 SQL |
| `init_rbac.sql` | RBAC 权限初始化 SQL |

**使用示例**:
```bash
# 运行迁移
./scripts/migrate.sh up

# 回滚迁移
./scripts/migrate.sh down

# 创建新迁移
./scripts/migrate.sh create add_user_avatar

# 测试迁移功能
./scripts/test_migration.sh
```

**详细文档**: [数据库迁移管理指南](../docs/MIGRATION_GUIDE.md)

---

### API 测试脚本

| 脚本 | 说明 |
|------|------|
| `test_frontend.sh` | 测试前台 API 接口 |
| `test_backend.sh` | 测试后台 API 接口 |
| `test_api.sh` | 测试基础 API 接口 |
| `test_admin_api.sh` | 测试管理员 API 接口 |

**使用示例**:
```bash
# 测试前台接口
./scripts/test_frontend.sh

# 测试后台接口
./scripts/test_backend.sh
```

---

### 功能测试脚本

| 脚本 | 说明 |
|------|------|
| `test_rbac.sh` | 测试 RBAC 权限功能 |
| `test_rbac_cache.sh` | 测试 RBAC 权限缓存 |
| `test_rate_limit.sh` | 测试限流功能 |
| `test_request_id.sh` | 测试请求 ID 追踪 |
| `test_env_switch.sh` | 测试环境配置切换 |

**使用示例**:
```bash
# 测试 RBAC 权限
./scripts/test_rbac.sh

# 测试限流功能
./scripts/test_rate_limit.sh

# 测试请求 ID
./scripts/test_request_id.sh
```

---

### 启动脚本

| 脚本 | 平台 | 说明 |
|------|------|------|
| `start-frontend.sh` | Linux / macOS / Git Bash | 启动前台服务 |
| `start-frontend.ps1` | PowerShell | 启动前台服务 |
| `start-backend.sh` | Linux / macOS / Git Bash | 启动后台服务 |
| `start-backend.ps1` | PowerShell | 启动后台服务 |
| `start-dev.sh` | Linux / macOS / Git Bash | 同时启动前后台服务（开发模式）|
| `start-dev.ps1` | PowerShell | 同时启动前后台服务（开发模式）|
| `stop-dev.sh` | Linux / macOS / Git Bash | 停止所有开发服务 |

**使用示例**:
```bash
# 启动前台服务（Linux/macOS）
./scripts/start-frontend.sh

# 启动后台服务（Linux/macOS）
./scripts/start-backend.sh

# 同时启动前后台（Linux/macOS）
./scripts/start-dev.sh

# PowerShell
.\scripts\start-frontend.ps1
.\scripts\start-backend.ps1
.\scripts\start-dev.ps1

# 停止所有服务
./scripts/stop-dev.sh
```

---

### Swagger 文档脚本

| 脚本 | 平台 | 说明 |
|------|------|------|
| `regenerate-swagger.sh` | Linux / macOS / Git Bash | 重新生成 Swagger 文档 |
| `regenerate-swagger.ps1` | PowerShell | 重新生成 Swagger 文档 |

**使用示例**:
```bash
# Linux/macOS
./scripts/regenerate-swagger.sh

# PowerShell
.\scripts\regenerate-swagger.ps1
```

---

### 工具脚本

| 脚本 | 说明 |
|------|------|
| `generate_admin_token.go` | 生成管理员 JWT Token |
| `generate_admin_with_role.go` | 生成带角色的管理员账号 |
| `verify.sh` | 验证项目依赖和配置 |

**使用示例**:
```bash
# 生成管理员 Token
go run scripts/generate_admin_token.go

# 生成带角色的管理员
go run scripts/generate_admin_with_role.go

# 验证项目
./scripts/verify.sh
```

---

## 🚀 快速开始

### 1. 第一次使用

```bash
# 添加执行权限（Linux / macOS）
chmod +x scripts/*.sh

# 或单独添加
chmod +x scripts/swagger.sh
chmod +x scripts/migrate.sh
```

### 2. 生成 Swagger 文档

```bash
# 使用最适合你平台的脚本
./scripts/swagger.sh              # Linux / macOS / Git Bash
scripts\swagger.bat               # Windows CMD
.\scripts\swagger.ps1             # PowerShell
```

### 3. 数据库迁移

```bash
# 运行所有迁移
./scripts/migrate.sh up

# 创建新迁移文件
./scripts/migrate.sh create add_new_feature
```

### 4. 测试接口

```bash
# 测试前台接口
./scripts/test_frontend.sh

# 测试后台接口
./scripts/test_backend.sh
```

## 📝 脚本命名规范

- **`swagger.*`** - Swagger 文档相关
- **`migrate.*`** - 数据库迁移相关
- **`test_*.sh`** - 测试脚本
- **`generate_*.go`** - 代码生成工具
- **`init_*.sql`** - 初始化 SQL 脚本

## 💡 使用提示

### Git Bash (Windows) 用户

在 Windows 上使用 Git Bash 时，推荐使用 `.sh` 脚本：

```bash
bash scripts/swagger.sh
bash scripts/migrate.sh
```

### Windows CMD 用户

使用批处理脚本：

```cmd
scripts\swagger.bat
```

### PowerShell 用户

首次使用需要设置执行策略：

```powershell
Set-ExecutionPolicy RemoteSigned -Scope CurrentUser
.\scripts\swagger.ps1
```

### CI/CD 集成

所有脚本都支持非交互式运行，适合 CI/CD 环境：

```yaml
# GitHub Actions 示例
- name: Generate Swagger Docs
  run: bash scripts/swagger.sh

- name: Run Migrations
  run: bash scripts/migrate.sh up
```

## 🔧 故障排查

### 脚本无法执行？

**Linux / macOS**:
```bash
# 确保有执行权限
ls -la scripts/*.sh

# 添加执行权限
chmod +x scripts/swagger.sh
```

**Windows**:
- 在 Git Bash 中使用 `bash scripts/swagger.sh`
- 在 CMD 中直接运行 `scripts\swagger.bat`
- PowerShell 需要设置执行策略

### 找不到命令？

确保已安装必要的工具：

```bash
# 安装 swag
go install github.com/swaggo/swag/cmd/swag@latest

# 安装 migrate
go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

### 脚本运行失败？

1. 检查是否在项目根目录
2. 确保依赖已安装（swag、Go、MySQL 等）
3. 查看详细的错误信息
4. 参考对应的详细文档

## 📚 相关文档

- [Swagger 脚本使用指南](../docs/SWAGGER_SCRIPTS_GUIDE.md)
- [Swagger 文档使用指南](../docs/SWAGGER_GUIDE.md)
- [数据库迁移管理指南](../docs/MIGRATION_GUIDE.md)
- [项目 README](../README.md)

## 🤝 贡献指南

添加新脚本时：

1. 遵循现有的命名规范
2. 添加详细的注释
3. 更新本 README
4. 提供使用示例
5. 考虑跨平台兼容性

## 📞 获取帮助

大多数脚本支持 `help` 参数：

```bash
./scripts/swagger.sh help
./scripts/migrate.sh help
```

查看脚本源码了解更多细节。

