# 数据库迁移管理实现总结

## 📋 实现概述

本项目集成了 [golang-migrate](https://github.com/golang-migrate/migrate) 进行数据库版本化管理，提供了完整的数据库迁移解决方案。

---

## 🏗️ 架构设计

### 核心组件

```
数据库迁移系统
├── migrations/              # 迁移文件目录
│   ├── *.up.sql            # 升级脚本
│   └── *.down.sql          # 回滚脚本
├── pkg/migrate/            # 迁移管理包
│   └── migrate.go          # 核心逻辑
├── cmd/migrate/            # 命令行工具
│   └── main.go
└── scripts/                # 辅助脚本
    ├── migrate.sh          # Shell 脚本
    └── test_migration.sh   # 测试脚本
```

### 设计原则

1. **版本化管理**: 每个数据库变更都有唯一的版本号
2. **可回滚性**: 每个迁移都包含 up 和 down 脚本
3. **环境适配**: 支持开发、测试、生产等多环境
4. **自动化**: 可选的自动迁移功能
5. **安全性**: 提供多重检查和确认机制

---

## 📦 实现的功能

### 1. 迁移文件管理

#### 自动编号
- 使用 6 位数字版本号：`000001`, `000002`, ...
- 自动识别当前最大版本号
- 创建新迁移时自动递增

#### 文件命名规范
```
000001_create_users_table.up.sql
000001_create_users_table.down.sql
```

#### 初始迁移文件
- `000001`: 创建 users 表
- `000002`: 创建 roles 表
- `000003`: 创建 permissions 表
- `000004`: 创建 role_permissions 关联表
- `000005`: 创建 user_roles 关联表
- `000006`: 初始化基础数据（角色和权限）

### 2. 迁移管理包 (pkg/migrate/)

#### Migrator 结构
```go
type Migrator struct {
    migrate *migrate.Migrate
    logger  *zap.Logger
}
```

#### 核心方法
- `Up()`: 执行所有待执行的迁移
- `Down()`: 回滚一个版本
- `Migrate(version)`: 迁移到指定版本
- `Version()`: 获取当前版本
- `Force(version)`: 强制设置版本
- `Drop()`: 删除所有表

#### 特性
- 日志集成：使用 Zap 记录所有操作
- 错误处理：详细的错误信息和恢复建议
- 状态检查：支持检查迁移状态和脏标志

### 3. 命令行工具 (cmd/migrate/)

#### 命令支持
```bash
go run cmd/migrate/main.go -cmd <command> [options]
```

| 命令 | 说明 | 示例 |
|------|------|------|
| up | 执行所有迁移 | `-cmd up` |
| down | 回滚一个版本 | `-cmd down` |
| version | 查看当前版本 | `-cmd version` |
| force | 强制设置版本 | `-cmd force -version 1` |
| goto | 迁移到指定版本 | `-cmd goto -version 3` |
| drop | 删除所有表 | `-cmd drop` |

#### 配置支持
```bash
# 使用不同环境的配置
-config config/config.dev.yaml
-config config/config.prod.yaml
```

### 4. Shell 脚本 (scripts/migrate.sh)

#### 功能
- 彩色输出和友好提示
- 环境变量支持
- 错误处理
- 二次确认（危险操作）

#### 使用方式
```bash
# 基本命令
./scripts/migrate.sh up
./scripts/migrate.sh version
./scripts/migrate.sh create

# 使用环境变量
CONFIG_FILE=config/config.prod.yaml ./scripts/migrate.sh up
VERSION=3 ./scripts/migrate.sh goto
NAME=add_user_phone ./scripts/migrate.sh create
```

### 5. Makefile 集成

#### 迁移命令
```bash
make migrate-up         # 执行迁移
make migrate-down       # 回滚
make migrate-version    # 查看版本
make migrate-create     # 创建新迁移
make migrate-force      # 强制版本
make migrate-goto       # 指定版本
make migrate-drop       # 删除所有表
```

#### 参数传递
```bash
make migrate-create NAME=add_user_phone
make migrate-force VERSION=1
make migrate-goto VERSION=3
```

### 6. 应用启动集成

#### 自动迁移
```bash
# 启动时自动执行迁移
AUTO_MIGRATE=true ./bin/frontend
AUTO_MIGRATE=true ./bin/backend
```

#### 实现细节
```go
// 在 main.go 中
if os.Getenv("AUTO_MIGRATE") == "true" {
    logger.Info("Running database migrations...")
    if err := runMigrations(cfg, logger); err != nil {
        log.Fatalf("Failed to run migrations: %v", err)
    }
}
```

#### 使用场景
- ✅ 开发环境：方便快速启动
- ⚠️ 测试环境：可选启用
- ❌ 生产环境：建议手动执行

---

## 📄 文件清单

### 新增文件

#### 迁移文件 (12 个)
```
migrations/000001_create_users_table.up.sql
migrations/000001_create_users_table.down.sql
migrations/000002_create_roles_table.up.sql
migrations/000002_create_roles_table.down.sql
migrations/000003_create_permissions_table.up.sql
migrations/000003_create_permissions_table.down.sql
migrations/000004_create_role_permissions_table.up.sql
migrations/000004_create_role_permissions_table.down.sql
migrations/000005_create_user_roles_table.up.sql
migrations/000005_create_user_roles_table.down.sql
migrations/000006_seed_initial_data.up.sql
migrations/000006_seed_initial_data.down.sql
```

#### 核心代码
```
pkg/migrate/migrate.go              # 迁移管理核心逻辑 (195 行)
cmd/migrate/main.go                 # 命令行工具 (130 行)
```

#### 脚本
```
scripts/migrate.sh                  # Shell 脚本 (150 行)
scripts/test_migration.sh           # 测试脚本 (130 行)
```

#### 文档
```
docs/MIGRATION_GUIDE.md             # 使用指南 (600+ 行)
docs/MIGRATION_IMPLEMENTATION.md    # 实现总结 (本文件)
```

### 修改文件

#### 应用入口
```
cmd/frontend/main.go                # 添加迁移逻辑
cmd/frontend/wire.go                # 添加 initLogger
cmd/backend/main.go                 # 添加迁移逻辑
```

#### 配置和构建
```
Makefile                            # 添加迁移命令
README.md                           # 更新文档
```

---

## 🔄 工作流程

### 开发环境工作流

```bash
# 1. 启动基础服务
make docker-up

# 2. 执行迁移
make migrate-up

# 3. 启动应用（自动迁移）
AUTO_MIGRATE=true make dev-frontend
AUTO_MIGRATE=true make dev-backend

# 4. 开发新功能（需要数据库变更）
make migrate-create NAME=add_user_avatar

# 5. 编辑迁移文件
vim migrations/000007_add_user_avatar.up.sql
vim migrations/000007_add_user_avatar.down.sql

# 6. 应用迁移
make migrate-up

# 7. 测试回滚
make migrate-down

# 8. 确认无误后提交
git add migrations/
git commit -m "feat: add user avatar field"
```

### 生产环境部署流程

```bash
# 1. 备份数据库
mysqldump -u root -p mydb > backup_$(date +%Y%m%d).sql

# 2. 查看当前版本
CONFIG_FILE=config/config.prod.yaml make migrate-version

# 3. 在测试环境验证
CONFIG_FILE=config/config.test.yaml make migrate-up

# 4. 应用到生产
CONFIG_FILE=config/config.prod.yaml make migrate-up

# 5. 验证结果
CONFIG_FILE=config/config.prod.yaml make migrate-version

# 6. 启动应用（不使用自动迁移）
./bin/frontend
./bin/backend
```

---

## 🎯 最佳实践

### 1. 迁移文件编写

✅ **好的做法**:
```sql
-- 使用 IF NOT EXISTS
CREATE TABLE IF NOT EXISTS `users` (...);

-- 添加默认值
ALTER TABLE `users` 
ADD COLUMN `status` INT NOT NULL DEFAULT 1;

-- 检查后删除
ALTER TABLE `users` 
DROP COLUMN IF EXISTS `temp_field`;
```

❌ **不好的做法**:
```sql
-- 直接创建（可能失败）
CREATE TABLE `users` (...);

-- 没有默认值的 NOT NULL 字段
ALTER TABLE `users` 
ADD COLUMN `status` INT NOT NULL;
```

### 2. 数据迁移

对于包含数据的表，分多步进行：

```sql
-- Step 1: 添加新列（允许 NULL）
ALTER TABLE `users` ADD COLUMN `new_email` VARCHAR(100) NULL;

-- Step 2: 迁移数据
UPDATE `users` SET `new_email` = `old_email`;

-- Step 3: 设置为 NOT NULL
ALTER TABLE `users` 
MODIFY COLUMN `new_email` VARCHAR(100) NOT NULL;

-- Step 4: 删除旧列（在后续迁移中）
-- 不要在同一个迁移中删除旧列，给应用一个适应期
```

### 3. 团队协作

- 📝 **沟通**: 迁移前通知团队成员
- 📝 **同步**: 定期同步代码，避免版本冲突
- 📝 **审查**: PR 中重点审查迁移文件
- 📝 **文档**: 复杂迁移要写详细注释

### 4. 版本控制

```bash
# ✅ 正确：迁移文件进入版本控制
git add migrations/
git commit -m "feat: add user avatar migration"

# ❌ 错误：修改已应用的迁移
# 永远不要修改已经应用到生产环境的迁移文件
```

---

## 🛡️ 安全机制

### 1. 脏状态检测
```bash
# 如果迁移失败，系统会标记为 dirty
# 需要手动修复后使用 force 命令
make migrate-force VERSION=3
```

### 2. 二次确认
```bash
# drop 命令需要输入 'YES' 确认
make migrate-drop
⚠️  警告：此操作将删除所有表！请输入 'YES' 确认:
```

### 3. 日志记录
所有迁移操作都会记录详细日志：
```
2024-01-15 10:30:00 INFO  Running database migrations...
2024-01-15 10:30:01 INFO  Database migrations completed successfully version=6
```

### 4. 环境隔离
不同环境使用不同的配置文件：
```bash
# 开发
CONFIG_FILE=config/config.dev.yaml

# 测试
CONFIG_FILE=config/config.test.yaml

# 生产
CONFIG_FILE=config/config.prod.yaml
```

---

## 📊 性能考虑

### 1. 迁移执行时间
- 表结构变更：通常很快（< 1s）
- 数据迁移：取决于数据量
- 索引创建：可能较慢（大表）

### 2. 优化建议
```sql
-- ✅ 为大表添加索引时使用 ALGORITHM=INPLACE
ALTER TABLE `large_table` 
ADD INDEX `idx_field` (`field`) 
ALGORITHM=INPLACE;

-- ✅ 批量更新数据
UPDATE `users` SET `status` = 1 
WHERE `status` IS NULL 
LIMIT 1000;
```

### 3. 生产环境建议
- 🔸 大表变更在低峰期执行
- 🔸 提前评估执行时间
- 🔸 准备回滚方案
- 🔸 监控数据库性能

---

## 🔧 故障排查

### 常见问题

#### 1. 迁移失败 (dirty state)
```bash
# 症状
error: Dirty database version 3

# 解决
# 1. 检查数据库，手动修复问题
# 2. 强制设置版本
make migrate-force VERSION=3
# 3. 重新运行
make migrate-up
```

#### 2. 连接失败
```bash
# 症状
error: dial tcp: connection refused

# 解决
make docker-up  # 启动数据库
```

#### 3. 版本冲突
```bash
# 多人协作时可能出现相同版本号

# 解决
# 1. Pull 最新代码
git pull origin main
# 2. 重命名迁移文件
mv migrations/000005_my.up.sql migrations/000007_my.up.sql
```

---

## 📈 未来改进

### 计划中的增强
1. **Web UI**: 可视化迁移管理界面
2. **并行迁移**: 支持多数据库并行迁移
3. **回滚策略**: 更智能的回滚机制
4. **性能监控**: 记录每次迁移的执行时间
5. **通知集成**: 迁移完成后发送通知

### 可能的扩展
1. **多数据库支持**: PostgreSQL, SQLite
2. **数据校验**: 迁移后自动校验数据一致性
3. **灰度迁移**: 分批次迁移大表
4. **版本回退**: 支持回退到任意历史版本

---

## 📚 参考资源

- [golang-migrate 官方文档](https://github.com/golang-migrate/migrate)
- [数据库迁移最佳实践](https://github.com/golang-migrate/migrate/blob/master/MIGRATIONS.md)
- [项目迁移使用指南](MIGRATION_GUIDE.md)

---

## 📝 总结

### 实现成果

✅ **完整的迁移系统**
- 6 个初始迁移文件
- 完善的工具链
- 详细的文档

✅ **多种使用方式**
- Makefile 命令
- Shell 脚本
- Go 命令行工具
- 自动迁移

✅ **开发友好**
- 简单易用
- 快速创建
- 充分测试

✅ **生产就绪**
- 安全机制
- 错误恢复
- 环境隔离

### 关键优势

1. **可追溯性**: 所有数据库变更都有版本记录
2. **可回滚性**: 任何变更都可以安全回滚
3. **团队协作**: 避免数据库结构冲突
4. **自动化**: 减少手动操作，降低错误
5. **文档完善**: 详细的使用和故障排查指南

---

## 🎓 学习建议

1. **阅读文档**: 从 [MIGRATION_GUIDE.md](MIGRATION_GUIDE.md) 开始
2. **实践操作**: 在开发环境尝试各种命令
3. **阅读代码**: 理解 `pkg/migrate/migrate.go` 的实现
4. **编写迁移**: 尝试创建自己的迁移文件
5. **测试回滚**: 熟悉回滚流程

---

**实现时间**: 2024-11
**维护者**: TRX Project Team
**版本**: 1.0.0

