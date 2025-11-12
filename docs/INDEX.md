# 文档索引

> 📚 TRX Project 文档导航 - 按功能分类，快速查找

## 📁 目录结构

```
docs/
├── INDEX.md                    # 📖 本文档（文档索引）
├── PROJECT_ORGANIZATION.md     # 📋 项目整理说明
│
├── getting-started/            # 🚀 快速开始
│   ├── QUICK_START.md
│   ├── QUICK_START_REFACTOR.md
│   ├── QUICK_TEST_GUIDE.md
│   ├── STARTUP_GUIDE.md
│   ├── QUICKSTART.md
│   ├── START.md
│   └── QUICK_START_ENV_SWAGGER.md
│
├── architecture/               # 📐 架构设计
│   ├── ARCHITECTURE.md
│   ├── PROJECT_STRUCTURE.md
│   ├── PROJECT_SUMMARY.md
│   ├── SEPARATED_SERVICES.md
│   ├── FRONTEND_BACKEND_SEPARATION.md
│   └── API_FRONTEND_BACKEND.md
│
├── features/                   # ⚡ 功能指南
│   ├── RBAC_GUIDE.md
│   ├── RBAC_CACHE_GUIDE.md
│   ├── RATE_LIMIT_GUIDE.md
│   ├── REQUEST_ID_GUIDE.md
│   ├── MIGRATION_GUIDE.md
│   ├── PROMETHEUS_MONITORING_GUIDE.md
│   ├── OPENTELEMETRY_TRACING_GUIDE.md
│   └── KAFKA_USAGE.md
│
├── implementation/             # 🔧 实现细节
│   ├── RBAC_IMPLEMENTATION_SUMMARY.md
│   ├── RBAC_CACHE_IMPLEMENTATION.md
│   ├── RATE_LIMIT_IMPLEMENTATION.md
│   ├── REQUEST_ID_IMPLEMENTATION.md
│   ├── MIGRATION_IMPLEMENTATION.md
│   ├── PROMETHEUS_IMPLEMENTATION.md
│   ├── OPENTELEMETRY_IMPLEMENTATION.md
│   ├── HANDLER_REFACTOR_SUMMARY.md
│   ├── REFACTOR_COMPLETION_REPORT.md
│   └── PACKAGE_NAME_UPDATE.md
│
├── tools/                      # 🛠️ 开发工具
│   ├── AIR_USAGE_GUIDE.md
│   ├── AIR_WINDOWS_GUIDE.md
│   ├── AIR_CONFIG_README.md
│   ├── SWAGGER_GUIDE.md
│   ├── SWAGGER_SCRIPTS_GUIDE.md
│   ├── SWAGGER_SCRIPTS_IMPLEMENTATION.md
│   ├── SWAGGER_SEPARATION_SUMMARY.md
│   ├── SWAGGER_FIX_GUIDE.md
│   ├── SWAGGER_GO_FILE_FIX.md
│   ├── SWAGGER_GO_FILE_FIX_SUMMARY.md
│   ├── SWAGGER_CROSS_PLATFORM_SUMMARY.md
│   └── ENV_AND_SWAGGER_GUIDE.md
│
├── deployment/                 # 🚀 部署运维
│   ├── DEPLOYMENT.md
│   ├── API.md
│   └── RESPONSE_FORMAT.md
│
└── legacy/                     # 📦 遗留文档
    ├── CHECKLIST.md
    ├── FEATURE_SUMMARY.md
    ├── README_RESPONSE.md
    ├── SUMMARY.md
    └── OPTIMIZATION_RECOMMENDATIONS.md
```

## 🎯 快速导航

### 🚀 我要开始使用项目

**推荐阅读顺序**：

1. [`getting-started/QUICK_START.md`](getting-started/QUICK_START.md) - 快速启动指南
2. [`getting-started/QUICK_TEST_GUIDE.md`](getting-started/QUICK_TEST_GUIDE.md) - API 测试示例
3. [`deployment/API.md`](deployment/API.md) - API 接口文档

**其他快速开始文档**：
- [`getting-started/STARTUP_GUIDE.md`](getting-started/STARTUP_GUIDE.md) - 详细启动步骤
- [`getting-started/QUICK_START_REFACTOR.md`](getting-started/QUICK_START_REFACTOR.md) - 重构后版本
- [`getting-started/QUICK_START_ENV_SWAGGER.md`](getting-started/QUICK_START_ENV_SWAGGER.md) - 环境与 Swagger 配置

---

### 📐 我要了解系统架构

**核心架构文档**：

- [`architecture/ARCHITECTURE.md`](architecture/ARCHITECTURE.md) - 系统架构设计 ⭐
- [`architecture/PROJECT_STRUCTURE.md`](architecture/PROJECT_STRUCTURE.md) - 项目结构说明
- [`architecture/SEPARATED_SERVICES.md`](architecture/SEPARATED_SERVICES.md) - 前后台分离架构 ⭐

**详细设计**：
- [`architecture/FRONTEND_BACKEND_SEPARATION.md`](architecture/FRONTEND_BACKEND_SEPARATION.md) - 服务分离实现
- [`architecture/API_FRONTEND_BACKEND.md`](architecture/API_FRONTEND_BACKEND.md) - API 对比说明
- [`architecture/PROJECT_SUMMARY.md`](architecture/PROJECT_SUMMARY.md) - 项目功能总结

---

### ⚡ 我要使用某个功能

**权限系统 (RBAC)**：
- [`features/RBAC_GUIDE.md`](features/RBAC_GUIDE.md) - RBAC 使用指南 ⭐
- [`features/RBAC_CACHE_GUIDE.md`](features/RBAC_CACHE_GUIDE.md) - 权限缓存优化

**限流与追踪**：
- [`features/RATE_LIMIT_GUIDE.md`](features/RATE_LIMIT_GUIDE.md) - 限流配置指南 ⭐
- [`features/REQUEST_ID_GUIDE.md`](features/REQUEST_ID_GUIDE.md) - 请求追踪指南

**监控与追踪**：
- [`features/PROMETHEUS_MONITORING_GUIDE.md`](features/PROMETHEUS_MONITORING_GUIDE.md) - Prometheus 监控 ⭐
- [`features/OPENTELEMETRY_TRACING_GUIDE.md`](features/OPENTELEMETRY_TRACING_GUIDE.md) - 链路追踪 ⭐

**数据库**：
- [`features/MIGRATION_GUIDE.md`](features/MIGRATION_GUIDE.md) - 数据库迁移管理 ⭐

**消息队列**：
- [`features/KAFKA_USAGE.md`](features/KAFKA_USAGE.md) - Kafka 使用说明

---

### 🔧 我要了解实现细节

**权限系统实现**：
- [`implementation/RBAC_IMPLEMENTATION_SUMMARY.md`](implementation/RBAC_IMPLEMENTATION_SUMMARY.md)
- [`implementation/RBAC_CACHE_IMPLEMENTATION.md`](implementation/RBAC_CACHE_IMPLEMENTATION.md)

**限流与追踪实现**：
- [`implementation/RATE_LIMIT_IMPLEMENTATION.md`](implementation/RATE_LIMIT_IMPLEMENTATION.md)
- [`implementation/REQUEST_ID_IMPLEMENTATION.md`](implementation/REQUEST_ID_IMPLEMENTATION.md)

**监控实现**：
- [`implementation/PROMETHEUS_IMPLEMENTATION.md`](implementation/PROMETHEUS_IMPLEMENTATION.md)
- [`implementation/OPENTELEMETRY_IMPLEMENTATION.md`](implementation/OPENTELEMETRY_IMPLEMENTATION.md)

**数据库迁移实现**：
- [`implementation/MIGRATION_IMPLEMENTATION.md`](implementation/MIGRATION_IMPLEMENTATION.md)

**重构记录**：
- [`implementation/HANDLER_REFACTOR_SUMMARY.md`](implementation/HANDLER_REFACTOR_SUMMARY.md) - Handler 重构
- [`implementation/REFACTOR_COMPLETION_REPORT.md`](implementation/REFACTOR_COMPLETION_REPORT.md) - 重构报告
- [`implementation/PACKAGE_NAME_UPDATE.md`](implementation/PACKAGE_NAME_UPDATE.md) - 包名更新

---

### 🛠️ 我要使用开发工具

**Air 热重载**：
- [`tools/AIR_USAGE_GUIDE.md`](tools/AIR_USAGE_GUIDE.md) - Air 使用指南 ⭐
- [`tools/AIR_WINDOWS_GUIDE.md`](tools/AIR_WINDOWS_GUIDE.md) - Windows 配置
- [`tools/AIR_CONFIG_README.md`](tools/AIR_CONFIG_README.md) - 配置说明

**Swagger 文档**：
- [`tools/SWAGGER_GUIDE.md`](tools/SWAGGER_GUIDE.md) - Swagger 使用指南 ⭐
- [`tools/SWAGGER_SCRIPTS_GUIDE.md`](tools/SWAGGER_SCRIPTS_GUIDE.md) - 脚本使用说明 ⭐
- [`tools/SWAGGER_FIX_GUIDE.md`](tools/SWAGGER_FIX_GUIDE.md) - 问题修复指南
- [`tools/SWAGGER_SEPARATION_SUMMARY.md`](tools/SWAGGER_SEPARATION_SUMMARY.md) - 分离总结
- [`tools/SWAGGER_CROSS_PLATFORM_SUMMARY.md`](tools/SWAGGER_CROSS_PLATFORM_SUMMARY.md) - 跨平台支持
- [`tools/SWAGGER_GO_FILE_FIX.md`](tools/SWAGGER_GO_FILE_FIX.md) - Go 文件修复
- [`tools/SWAGGER_SCRIPTS_IMPLEMENTATION.md`](tools/SWAGGER_SCRIPTS_IMPLEMENTATION.md) - 脚本实现
- [`tools/ENV_AND_SWAGGER_GUIDE.md`](tools/ENV_AND_SWAGGER_GUIDE.md) - 环境与 Swagger

---

### 🚀 我要部署到生产环境

- [`deployment/DEPLOYMENT.md`](deployment/DEPLOYMENT.md) - 生产部署指南 ⭐
- [`deployment/API.md`](deployment/API.md) - API 接口文档
- [`deployment/RESPONSE_FORMAT.md`](deployment/RESPONSE_FORMAT.md) - 响应格式规范

---

## 🔍 按问题查找

### 启动相关问题

| 问题 | 查看文档 |
|------|---------|
| 如何快速启动项目？ | [`getting-started/QUICK_START.md`](getting-started/QUICK_START.md) |
| 启动失败怎么办？ | [`getting-started/STARTUP_GUIDE.md`](getting-started/STARTUP_GUIDE.md) |
| 如何测试 API？ | [`getting-started/QUICK_TEST_GUIDE.md`](getting-started/QUICK_TEST_GUIDE.md) |
| Air 热重载配置？ | [`tools/AIR_USAGE_GUIDE.md`](tools/AIR_USAGE_GUIDE.md) |

### 权限相关问题

| 问题 | 查看文档 |
|------|---------|
| 如何配置 RBAC 权限？ | [`features/RBAC_GUIDE.md`](features/RBAC_GUIDE.md) |
| 权限检查太慢？ | [`features/RBAC_CACHE_GUIDE.md`](features/RBAC_CACHE_GUIDE.md) |
| RBAC 如何实现的？ | [`implementation/RBAC_IMPLEMENTATION_SUMMARY.md`](implementation/RBAC_IMPLEMENTATION_SUMMARY.md) |

### 性能相关问题

| 问题 | 查看文档 |
|------|---------|
| 如何配置限流？ | [`features/RATE_LIMIT_GUIDE.md`](features/RATE_LIMIT_GUIDE.md) |
| 如何优化性能？ | [`features/RBAC_CACHE_GUIDE.md`](features/RBAC_CACHE_GUIDE.md) |
| 如何监控性能？ | [`features/PROMETHEUS_MONITORING_GUIDE.md`](features/PROMETHEUS_MONITORING_GUIDE.md) |

### Swagger 相关问题

| 问题 | 查看文档 |
|------|---------|
| 如何生成 Swagger 文档？ | [`tools/SWAGGER_GUIDE.md`](tools/SWAGGER_GUIDE.md) |
| Swagger 脚本怎么用？ | [`tools/SWAGGER_SCRIPTS_GUIDE.md`](tools/SWAGGER_SCRIPTS_GUIDE.md) |
| Swagger 生成失败？ | [`tools/SWAGGER_FIX_GUIDE.md`](tools/SWAGGER_FIX_GUIDE.md) |

### 数据库相关问题

| 问题 | 查看文档 |
|------|---------|
| 如何执行数据库迁移？ | [`features/MIGRATION_GUIDE.md`](features/MIGRATION_GUIDE.md) |
| 迁移如何实现的？ | [`implementation/MIGRATION_IMPLEMENTATION.md`](implementation/MIGRATION_IMPLEMENTATION.md) |

### 监控相关问题

| 问题 | 查看文档 |
|------|---------|
| 如何配置 Prometheus？ | [`features/PROMETHEUS_MONITORING_GUIDE.md`](features/PROMETHEUS_MONITORING_GUIDE.md) |
| 如何追踪请求？ | [`features/OPENTELEMETRY_TRACING_GUIDE.md`](features/OPENTELEMETRY_TRACING_GUIDE.md) |
| 如何追踪请求 ID？ | [`features/REQUEST_ID_GUIDE.md`](features/REQUEST_ID_GUIDE.md) |

---

## 📚 推荐学习路径

### 🎓 新手入门（必读）

1. [`getting-started/QUICK_START.md`](getting-started/QUICK_START.md) - 快速启动
2. [`architecture/ARCHITECTURE.md`](architecture/ARCHITECTURE.md) - 了解架构
3. [`getting-started/QUICK_TEST_GUIDE.md`](getting-started/QUICK_TEST_GUIDE.md) - 测试 API
4. [`tools/SWAGGER_GUIDE.md`](tools/SWAGGER_GUIDE.md) - API 文档

### 🚀 功能使用（重要）

1. [`features/RBAC_GUIDE.md`](features/RBAC_GUIDE.md) - 权限系统
2. [`features/MIGRATION_GUIDE.md`](features/MIGRATION_GUIDE.md) - 数据库管理
3. [`features/RATE_LIMIT_GUIDE.md`](features/RATE_LIMIT_GUIDE.md) - 限流配置
4. [`features/PROMETHEUS_MONITORING_GUIDE.md`](features/PROMETHEUS_MONITORING_GUIDE.md) - 监控系统

### 💻 深入开发（进阶）

1. [`architecture/SEPARATED_SERVICES.md`](architecture/SEPARATED_SERVICES.md) - 服务分离
2. [`implementation/HANDLER_REFACTOR_SUMMARY.md`](implementation/HANDLER_REFACTOR_SUMMARY.md) - 代码重构
3. [`tools/AIR_USAGE_GUIDE.md`](tools/AIR_USAGE_GUIDE.md) - 开发工具
4. [`implementation/` 目录](implementation/) - 各功能实现细节

### 🎯 生产部署（上线）

1. [`deployment/DEPLOYMENT.md`](deployment/DEPLOYMENT.md) - 部署指南
2. [`features/PROMETHEUS_MONITORING_GUIDE.md`](features/PROMETHEUS_MONITORING_GUIDE.md) - 监控配置
3. [`features/OPENTELEMETRY_TRACING_GUIDE.md`](features/OPENTELEMETRY_TRACING_GUIDE.md) - 链路追踪
4. [`deployment/RESPONSE_FORMAT.md`](deployment/RESPONSE_FORMAT.md) - API 规范

---

## 💡 使用提示

### 如何快速找到需要的文档？

1. **按目录浏览** - 根据功能分类查看对应目录
2. **按问题查找** - 参考上面的"按问题查找"表格
3. **按学习路径** - 跟随推荐的学习路径阅读
4. **全文搜索** - 在 IDE 中搜索关键词

### 文档命名规范

- `*_GUIDE.md` - 使用指南（面向用户）
- `*_IMPLEMENTATION.md` - 实现细节（面向开发者）
- `*_SUMMARY.md` - 功能总结
- `QUICK_*.md` - 快速开始类文档

### 标记说明

- ⭐ - 重要文档，推荐优先阅读
- 📖 - 指南类文档
- 🔧 - 技术实现文档
- 📋 - 总结类文档

---

## 📞 获取帮助

1. **查看主 README** - [`../README.md`](../README.md)
2. **查看脚本说明** - [`../scripts/README.md`](../scripts/README.md)
3. **项目整理说明** - [`PROJECT_ORGANIZATION.md`](PROJECT_ORGANIZATION.md)
4. **提交 Issue** - 在项目仓库提交问题

---

**文档总数**: 53 个  
**最后更新**: 2024-11-12  
**目录结构**: 按功能分类，便于查找和维护
