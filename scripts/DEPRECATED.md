# 已废弃的脚本

以下脚本已被新的架构替代，不再使用：

## ❌ filter_swagger.go (已废弃)

**原用途**: 过滤 Swagger JSON 文件中的特定路径

**废弃原因**: 
- 采用目录结构分离后，使用 `swag --exclude` 参数更简单
- 无需后处理过滤
- 维护成本高

**替代方案**: 使用 `swag init --exclude` 参数

## ❌ filter_swagger_docs.go (已废弃)

**原用途**: 过滤 Swagger JSON 和 Go 文件中的特定路径

**废弃原因**:
- Handler 目录已按 frontend/backend 分离
- 使用 `swag --exclude` 在生成时就能正确分离
- 不再需要后处理

**替代方案**: 使用 `swag init --exclude` 参数

## 新架构

### Handler 目录结构

```
internal/api/handler/
├── frontend/       # 前台 handler
│   └── user_handler.go
└── backend/        # 后台 handler
    ├── admin_user_handler.go
    └── rbac_handler.go
```

### Swagger 生成

```bash
# 前台 - 排除后台目录
swag init -g cmd/frontend/main.go ... \
    --exclude internal/api/handler/backend

# 后台 - 排除前台目录
swag init -g cmd/backend/main.go ... \
    --exclude internal/api/handler/frontend
```

## 迁移说明

如果你的代码还在使用这些脚本，请参考以下文档进行迁移：

- [Handler 重构总结](../HANDLER_REFACTOR_SUMMARY.md)
- [Swagger 使用指南](../docs/SWAGGER_GUIDE.md)

## 删除计划

这些脚本将在下一个主要版本中移除。请尽快迁移到新架构。

