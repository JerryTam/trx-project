# 数据库迁移管理指南

本项目使用 [golang-migrate](https://github.com/golang-migrate/migrate) 进行数据库版本化管理。

## 📚 目录

- [基本概念](#基本概念)
- [快速开始](#快速开始)
- [常用命令](#常用命令)
- [迁移文件管理](#迁移文件管理)
- [最佳实践](#最佳实践)
- [故障排查](#故障排查)

---

## 基本概念

### 什么是数据库迁移？

数据库迁移是一种版本化管理数据库结构变更的方法，类似于代码的版本控制。每个迁移包含两个文件：
- **up.sql**: 升级脚本（应用变更）
- **down.sql**: 回滚脚本（撤销变更）

### 迁移版本号

迁移文件使用 6 位数字作为版本号，例如：
```
migrations/
├── 000001_create_users_table.up.sql
├── 000001_create_users_table.down.sql
├── 000002_create_roles_table.up.sql
└── 000002_create_roles_table.down.sql
```

---

## 快速开始

### 1. 初始化数据库

首次部署时，运行所有迁移：

```bash
make migrate-up
# 或
./scripts/migrate.sh up
```

### 2. 查看当前版本

```bash
make migrate-version
# 或
./scripts/migrate.sh version
```

### 3. 创建新迁移

```bash
make migrate-create NAME=add_user_phone
# 或
NAME=add_user_phone ./scripts/migrate.sh create
```

这将创建两个新文件：
- `migrations/000007_add_user_phone.up.sql`
- `migrations/000007_add_user_phone.down.sql`

---

## 常用命令

### Makefile 命令

| 命令 | 说明 | 示例 |
|------|------|------|
| `make migrate-up` | 执行所有待执行的迁移 | `make migrate-up` |
| `make migrate-down` | 回滚一个迁移版本 | `make migrate-down` |
| `make migrate-version` | 查看当前迁移版本 | `make migrate-version` |
| `make migrate-create` | 创建新的迁移文件 | `make migrate-create NAME=add_column` |
| `make migrate-force` | 强制设置迁移版本 | `make migrate-force VERSION=1` |
| `make migrate-goto` | 迁移到指定版本 | `make migrate-goto VERSION=3` |
| `make migrate-drop` | 删除所有表（危险） | `make migrate-drop` |

### Shell 脚本命令

```bash
# 执行所有迁移
./scripts/migrate.sh up

# 回滚一个版本
./scripts/migrate.sh down

# 查看当前版本
./scripts/migrate.sh version

# 创建新迁移文件
NAME=add_user_phone ./scripts/migrate.sh create

# 强制设置版本（修复脏状态）
VERSION=1 ./scripts/migrate.sh force

# 迁移到指定版本
VERSION=3 ./scripts/migrate.sh goto

# 删除所有表（危险操作）
./scripts/migrate.sh drop
```

### Go 命令行工具

```bash
# 执行所有迁移
go run cmd/migrate/main.go -cmd up

# 查看当前版本
go run cmd/migrate/main.go -cmd version

# 强制设置版本
go run cmd/migrate/main.go -cmd force -version 1

# 使用不同的配置文件
go run cmd/migrate/main.go -config config/config.prod.yaml -cmd up
```

---

## 迁移文件管理

### 创建迁移文件

使用工具自动创建：
```bash
make migrate-create NAME=add_user_phone
```

或手动创建文件（遵循命名规范）：
```
migrations/000007_add_user_phone.up.sql
migrations/000007_add_user_phone.down.sql
```

### 编写 UP 脚本

在 `up.sql` 中编写升级 SQL：

```sql
-- 000007_add_user_phone.up.sql
-- 为用户表添加手机号字段

ALTER TABLE `users` 
ADD COLUMN `phone` VARCHAR(20) NULL COMMENT '手机号' AFTER `email`,
ADD INDEX `idx_users_phone` (`phone`);
```

### 编写 DOWN 脚本

在 `down.sql` 中编写回滚 SQL：

```sql
-- 000007_add_user_phone.down.sql
-- 移除用户表的手机号字段

ALTER TABLE `users`
DROP INDEX `idx_users_phone`,
DROP COLUMN `phone`;
```

### 迁移文件示例

#### 创建表

**up.sql**:
```sql
CREATE TABLE IF NOT EXISTS `products` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(100) NOT NULL,
    `price` DECIMAL(10,2) NOT NULL,
    `created_at` DATETIME(3) NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='产品表';
```

**down.sql**:
```sql
DROP TABLE IF EXISTS `products`;
```

#### 修改表结构

**up.sql**:
```sql
ALTER TABLE `users` 
ADD COLUMN `avatar` VARCHAR(255) NULL COMMENT '头像URL',
ADD COLUMN `bio` TEXT NULL COMMENT '个人简介';
```

**down.sql**:
```sql
ALTER TABLE `users`
DROP COLUMN `bio`,
DROP COLUMN `avatar`;
```

#### 数据迁移

**up.sql**:
```sql
-- 将旧的 status 字段值转换为新格式
UPDATE `users` 
SET `status` = CASE 
    WHEN `status` = 1 THEN 'active'
    WHEN `status` = 0 THEN 'inactive'
    ELSE 'unknown'
END;
```

**down.sql**:
```sql
-- 将新的 status 字段值转换回旧格式
UPDATE `users` 
SET `status` = CASE 
    WHEN `status` = 'active' THEN 1
    WHEN `status` = 'inactive' THEN 0
    ELSE -1
END;
```

---

## 最佳实践

### 1. 迁移文件原则

- ✅ **保持小而专注**：每个迁移只做一件事
- ✅ **可逆性**：确保 down 脚本能正确回滚
- ✅ **幂等性**：迁移可以安全地重复执行
- ✅ **测试性**：在开发环境充分测试后再部署
- ❌ **不要修改已应用的迁移**：会导致版本冲突

### 2. 编写安全的迁移

```sql
-- ✅ 好的做法：使用 IF NOT EXISTS
CREATE TABLE IF NOT EXISTS `users` (...);

-- ❌ 不好的做法：直接创建
CREATE TABLE `users` (...);

-- ✅ 好的做法：先检查再删除
ALTER TABLE `users` DROP COLUMN IF EXISTS `temp_column`;

-- ✅ 好的做法：添加默认值避免 NOT NULL 问题
ALTER TABLE `users` 
ADD COLUMN `new_field` VARCHAR(50) NOT NULL DEFAULT '';
```

### 3. 数据迁移策略

对于包含数据的表：
```sql
-- Step 1: 添加新列（允许 NULL）
ALTER TABLE `users` ADD COLUMN `new_email` VARCHAR(100) NULL;

-- Step 2: 迁移数据
UPDATE `users` SET `new_email` = `old_email`;

-- Step 3: 设置为 NOT NULL（如需要）
ALTER TABLE `users` MODIFY COLUMN `new_email` VARCHAR(100) NOT NULL;

-- Step 4: 删除旧列（在后续迁移中）
-- ALTER TABLE `users` DROP COLUMN `old_email`;
```

### 4. 团队协作

- 📝 在合并前先 pull 最新代码，避免版本号冲突
- 📝 给迁移文件起有意义的名称
- 📝 在 PR 中说明迁移的目的和影响
- 📝 在生产环境应用前，在测试环境验证

### 5. 生产环境部署

```bash
# 1. 备份数据库
mysqldump -u root -p mydb > backup_$(date +%Y%m%d_%H%M%S).sql

# 2. 查看当前版本
make migrate-version

# 3. 测试迁移（在测试环境）
make migrate-up

# 4. 应用到生产环境
CONFIG_FILE=config/config.prod.yaml make migrate-up

# 5. 验证版本
CONFIG_FILE=config/config.prod.yaml make migrate-version
```

---

## 自动迁移

### 在应用启动时自动迁移

设置环境变量 `AUTO_MIGRATE=true` 可以在应用启动时自动运行迁移：

```bash
# 启动前台服务并自动迁移
AUTO_MIGRATE=true ./bin/frontend

# 启动后台服务并自动迁移
AUTO_MIGRATE=true ./bin/backend
```

⚠️ **注意**: 生产环境建议手动运行迁移，不要启用自动迁移。

---

## 故障排查

### 1. 迁移处于脏状态 (dirty)

**症状**:
```
error: Dirty database version 3. Fix and force version.
```

**原因**: 迁移执行失败，数据库处于不一致状态。

**解决方案**:
```bash
# 1. 检查数据库状态，手动修复问题
# 2. 确认数据库应该在哪个版本
# 3. 强制设置版本
make migrate-force VERSION=3

# 4. 重新运行迁移
make migrate-up
```

### 2. 迁移文件不存在

**症状**:
```
error: file does not exist
```

**原因**: 数据库记录的版本号对应的迁移文件不存在。

**解决方案**:
```bash
# 1. 检查 migrations 目录
ls migrations/

# 2. 确认缺失的迁移文件版本
make migrate-version

# 3. 恢复缺失的迁移文件，或强制设置到存在的版本
make migrate-force VERSION=<existing_version>
```

### 3. 数据库连接失败

**症状**:
```
error: dial tcp: connect: connection refused
```

**解决方案**:
```bash
# 1. 检查数据库是否运行
docker ps | grep mysql

# 2. 检查配置文件中的数据库连接信息
cat config/config.yaml

# 3. 启动数据库
make docker-up
```

### 4. 权限不足

**症状**:
```
error: Access denied for user
```

**解决方案**:
- 检查数据库用户权限
- 确认配置文件中的用户名和密码正确
- 确保用户有 CREATE、ALTER、DROP 等权限

### 5. 版本号冲突

多人协作时可能出现版本号冲突。

**解决方案**:
```bash
# 1. Pull 最新代码
git pull origin main

# 2. 重命名你的迁移文件，使用新的版本号
mv migrations/000005_my_feature.up.sql migrations/000007_my_feature.up.sql
mv migrations/000005_my_feature.down.sql migrations/000007_my_feature.down.sql
```

---

## 目录结构

```
trx-project/
├── migrations/                 # 迁移文件目录
│   ├── 000001_create_users_table.up.sql
│   ├── 000001_create_users_table.down.sql
│   ├── 000002_create_roles_table.up.sql
│   ├── 000002_create_roles_table.down.sql
│   └── ...
├── pkg/migrate/                # 迁移包
│   └── migrate.go              # 迁移逻辑
├── cmd/migrate/                # 迁移命令行工具
│   └── main.go
├── scripts/                    # 辅助脚本
│   └── migrate.sh              # 迁移脚本
└── Makefile                    # Make 命令
```

---

## 相关资源

- [golang-migrate 官方文档](https://github.com/golang-migrate/migrate)
- [数据库迁移最佳实践](https://github.com/golang-migrate/migrate/blob/master/MIGRATIONS.md)
- [MySQL 数据类型参考](https://dev.mysql.com/doc/refman/8.0/en/data-types.html)

---

## 总结

数据库迁移管理是保证数据库结构变更可追溯、可回滚的重要手段。通过：

1. ✅ 遵循最佳实践
2. ✅ 编写清晰的迁移脚本
3. ✅ 充分测试
4. ✅ 谨慎部署

可以确保数据库变更的安全性和可靠性。

