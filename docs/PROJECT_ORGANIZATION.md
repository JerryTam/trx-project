# 项目目录整理说明

> 本文档说明项目目录结构的整理和优化

## 📋 整理内容

### ✅ 已完成的整理

#### 1. 文档整理

**移动到 `docs/` 目录的文档**:

所有根目录下的 `.md` 文档（除了 `README.md`）已移动到 `docs/` 目录，包括：

- **快速开始类**: 
  - `QUICK_START.md`
  - `QUICK_START_REFACTOR.md`
  - `QUICK_TEST_GUIDE.md`
  - `QUICKSTART.md`
  - `START.md`

- **架构设计类**:
  - `ARCHITECTURE.md` (已在 docs/)
  - `PROJECT_STRUCTURE.md`
  - `PROJECT_SUMMARY.md`
  - `FRONTEND_BACKEND_SEPARATION.md`
  - `SEPARATED_SERVICES.md`

- **功能实现类**:
  - `RBAC_GUIDE.md`
  - `RBAC_IMPLEMENTATION_SUMMARY.md`
  - `RBAC_CACHE_GUIDE.md`
  - `RBAC_CACHE_IMPLEMENTATION.md`
  - `RATE_LIMIT_GUIDE.md`
  - `RATE_LIMIT_IMPLEMENTATION.md`
  - `REQUEST_ID_GUIDE.md`
  - `REQUEST_ID_IMPLEMENTATION.md`

- **重构相关**:
  - `HANDLER_REFACTOR_SUMMARY.md`
  - `REFACTOR_COMPLETION_REPORT.md`
  - `PACKAGE_NAME_UPDATE.md`

- **Swagger 相关**:
  - `ENV_AND_SWAGGER_GUIDE.md`
  - `QUICK_START_ENV_SWAGGER.md`
  - `SWAGGER_CROSS_PLATFORM_SUMMARY.md`
  - `SWAGGER_GO_FILE_FIX_SUMMARY.md`

- **其他**:
  - `AIR_CONFIG_README.md`
  - `CHECKLIST.md`
  - `FEATURE_SUMMARY.md`
  - `OPTIMIZATION_RECOMMENDATIONS.md`
  - `README_RESPONSE.md`
  - `SUMMARY.md`

#### 2. 脚本整理

**移动到 `scripts/` 目录的脚本**:

- **启动脚本**:
  - `start-frontend.sh`
  - `start-frontend.ps1`
  - `start-backend.sh`
  - `start-backend.ps1`
  - `start-dev.sh`
  - `start-dev.ps1`
  - `stop-dev.sh`

- **Swagger 脚本**:
  - `regenerate-swagger.sh`
  - `regenerate-swagger.ps1`

#### 3. 临时文件清理

**已删除的临时文件**:

- `backend.exe` - 临时编译文件
- `frontend.exe` - 临时编译文件
- `API_ROUTES.txt` - 临时文本文件
- `TREE.txt` - 临时文本文件

#### 4. .gitignore 更新

**新增忽略规则**:

```gitignore
# 临时文件
tmp/
temp/
*.txt
*.exe
!scripts/*.exe
```

这样可以：
- ✅ 忽略根目录的 `.exe` 和 `.txt` 文件
- ✅ 允许 scripts 目录下的可执行文件（如果有）
- ✅ 自动忽略临时生成的文本文件

#### 5. 新增文档索引

**创建的新文档**:

- `docs/INDEX.md` - 完整的文档索引，便于快速查找
- `docs/PROJECT_ORGANIZATION.md` - 本文档，说明整理内容

## 📁 整理后的目录结构

### 根目录（简洁）

```
trx-project/
├── README.md                    # 主说明文档（唯一的根目录文档）
├── Makefile                     # 构建脚本
├── docker-compose.yml           # Docker 配置
├── go.mod / go.sum              # Go 依赖
├── .gitignore                   # Git 忽略规则
├── .air-*.toml                  # Air 配置文件
├── bin/                         # 可执行文件
├── cmd/                         # 程序入口
├── config/                      # 配置文件
├── docs/                        # 📚 所有文档
├── internal/                    # 内部代码
├── logs/                        # 日志文件
├── migrations/                  # 数据库迁移
├── pkg/                         # 公共包
├── scripts/                     # 🔧 所有脚本
└── tmp/                         # 临时文件
```

### docs/ 目录（完整文档）

```
docs/
├── INDEX.md                              # 📖 文档索引（重要！）
├── PROJECT_ORGANIZATION.md               # 本文档
│
├── # 快速开始
├── QUICK_START.md
├── QUICK_START_REFACTOR.md
├── QUICK_TEST_GUIDE.md
├── STARTUP_GUIDE.md
│
├── # 架构设计
├── ARCHITECTURE.md
├── PROJECT_STRUCTURE.md
├── SEPARATED_SERVICES.md
│
├── # 权限系统
├── RBAC_GUIDE.md
├── RBAC_CACHE_GUIDE.md
├── RBAC_IMPLEMENTATION_SUMMARY.md
│
├── # 监控追踪
├── PROMETHEUS_MONITORING_GUIDE.md
├── OPENTELEMETRY_TRACING_GUIDE.md
│
├── # API 文档
├── API.md
├── SWAGGER_GUIDE.md
├── SWAGGER_SCRIPTS_GUIDE.md
├── RESPONSE_FORMAT.md
│
├── # 数据库
├── MIGRATION_GUIDE.md
├── MIGRATION_IMPLEMENTATION.md
│
├── # 开发工具
├── AIR_USAGE_GUIDE.md
├── AIR_WINDOWS_GUIDE.md
│
└── # 部署运维
    └── DEPLOYMENT.md
```

### scripts/ 目录（实用脚本）

```
scripts/
├── README.md                             # 📖 脚本说明文档
│
├── # 启动脚本
├── start-frontend.sh
├── start-frontend.ps1
├── start-backend.sh
├── start-backend.ps1
├── start-dev.sh
├── start-dev.ps1
├── stop-dev.sh
│
├── # Swagger
├── swagger.sh / swagger.bat / swagger.ps1
├── regenerate-swagger.sh
├── regenerate-swagger.ps1
│
├── # 数据库迁移
├── migrate.sh
├── test_migration.sh
├── init_db.sql
├── init_rbac.sql
│
├── # API 测试
├── test_frontend.sh
├── test_backend.sh
├── test_api.sh
├── test_admin_api.sh
│
├── # 功能测试
├── test_rbac.sh
├── test_rbac_cache.sh
├── test_rate_limit.sh
├── test_request_id.sh
│
└── # 工具
    ├── generate_admin_token.go
    ├── generate_admin_with_role.go
    └── verify.sh
```

## 🎯 整理的优势

### 1. 根目录清爽

- ✅ 只保留必要的配置文件和核心文档（`README.md`）
- ✅ 目录结构一目了然
- ✅ 符合 Go 项目最佳实践

### 2. 文档集中管理

- ✅ 所有文档都在 `docs/` 目录
- ✅ 创建了文档索引 (`INDEX.md`)
- ✅ 便于查找和维护

### 3. 脚本统一管理

- ✅ 所有脚本都在 `scripts/` 目录
- ✅ 更新了 scripts 的 README
- ✅ 按功能分类清晰

### 4. 自动化忽略

- ✅ 更新了 `.gitignore`
- ✅ 自动忽略临时文件
- ✅ 避免误提交

## 📖 如何查找文档

### 方法 1：使用文档索引

查看 [`docs/INDEX.md`](INDEX.md) 获取完整的文档索引和导航。

### 方法 2：按需查找

根据需求快速定位：

| 需求 | 查看文档 |
|------|---------|
| 快速启动项目 | [`docs/QUICK_START.md`](QUICK_START.md) |
| 测试 API | [`docs/QUICK_TEST_GUIDE.md`](QUICK_TEST_GUIDE.md) |
| 了解架构 | [`docs/ARCHITECTURE.md`](ARCHITECTURE.md) |
| 配置权限 | [`docs/RBAC_GUIDE.md`](RBAC_GUIDE.md) |
| 数据库迁移 | [`docs/MIGRATION_GUIDE.md`](MIGRATION_GUIDE.md) |
| 生成 Swagger | [`docs/SWAGGER_GUIDE.md`](SWAGGER_GUIDE.md) |
| 监控配置 | [`docs/PROMETHEUS_MONITORING_GUIDE.md`](PROMETHEUS_MONITORING_GUIDE.md) |
| 生产部署 | [`docs/DEPLOYMENT.md`](DEPLOYMENT.md) |

### 方法 3：使用脚本

查看 [`scripts/README.md`](../scripts/README.md) 获取所有脚本的使用说明。

## 🔄 迁移指南

如果你有旧的链接或书签：

### 文档链接更新

| 旧路径 | 新路径 |
|--------|--------|
| `/QUICK_START.md` | `/docs/QUICK_START.md` |
| `/RBAC_GUIDE.md` | `/docs/RBAC_GUIDE.md` |
| `/MIGRATION_GUIDE.md` | `/docs/MIGRATION_GUIDE.md` |
| 其他根目录 .md 文件 | `/docs/xxx.md` |

### 脚本路径更新

| 旧路径 | 新路径 |
|--------|--------|
| `/start-frontend.sh` | `/scripts/start-frontend.sh` |
| `/start-backend.sh` | `/scripts/start-backend.sh` |
| `/regenerate-swagger.sh` | `/scripts/regenerate-swagger.sh` |

### 在代码中的引用

如果代码或其他文档中引用了旧路径，需要更新：

```bash
# 旧
./start-frontend.sh

# 新
./scripts/start-frontend.sh
```

```markdown
<!-- 旧 -->
详见：[RBAC 指南](RBAC_GUIDE.md)

<!-- 新 -->
详见：[RBAC 指南](docs/RBAC_GUIDE.md)
```

## 💡 最佳实践

### 添加新文档时

1. 将文档放在 `docs/` 目录
2. 更新 `docs/INDEX.md`
3. 在主 `README.md` 中添加相关链接（如果重要）

### 添加新脚本时

1. 将脚本放在 `scripts/` 目录
2. 更新 `scripts/README.md`
3. 添加执行权限（Linux/macOS）
4. 考虑跨平台兼容性

### 保持根目录整洁

- ❌ 不要在根目录创建临时文件
- ❌ 不要在根目录添加新的 .md 文档
- ✅ 使用 `tmp/` 目录存放临时文件
- ✅ 新文档放在 `docs/` 目录
- ✅ 新脚本放在 `scripts/` 目录

## 📞 反馈

如果有任何问题或建议，欢迎提 Issue 或 PR！

---

**整理完成时间**: 2024-11-12  
**整理目的**: 提升项目的可维护性和专业性

