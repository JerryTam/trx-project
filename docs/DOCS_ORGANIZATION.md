# 📚 文档目录整理说明

> 本文档说明 `docs/` 目录的分类整理结果

## ✅ 整理完成

已将 53 个文档按功能分类整理到 7 个子目录中，便于查找和维护。

## 📁 目录结构

```
docs/
├── INDEX.md                        # 📖 文档总索引（重要！）
├── PROJECT_ORGANIZATION.md         # 项目整理说明
├── DOCS_ORGANIZATION.md            # 本文档（文档分类说明）
│
├── getting-started/                # 🚀 快速开始 (8 个文档)
│   ├── README.md
│   ├── QUICK_START.md              ⭐ 推荐
│   ├── QUICK_TEST_GUIDE.md         ⭐ 推荐
│   ├── STARTUP_GUIDE.md
│   ├── QUICK_START_REFACTOR.md
│   ├── QUICK_START_ENV_SWAGGER.md
│   ├── QUICKSTART.md
│   └── START.md
│
├── architecture/                   # 📐 架构设计 (6 个文档)
│   ├── ARCHITECTURE.md             ⭐ 推荐
│   ├── PROJECT_STRUCTURE.md
│   ├── PROJECT_SUMMARY.md
│   ├── SEPARATED_SERVICES.md       ⭐ 推荐
│   ├── FRONTEND_BACKEND_SEPARATION.md
│   └── API_FRONTEND_BACKEND.md
│
├── features/                       # ⚡ 功能指南 (9 个文档)
│   ├── README.md
│   ├── RBAC_GUIDE.md               ⭐ 推荐
│   ├── RBAC_CACHE_GUIDE.md
│   ├── RATE_LIMIT_GUIDE.md         ⭐ 推荐
│   ├── REQUEST_ID_GUIDE.md
│   ├── MIGRATION_GUIDE.md          ⭐ 推荐
│   ├── PROMETHEUS_MONITORING_GUIDE.md  ⭐ 推荐
│   ├── OPENTELEMETRY_TRACING_GUIDE.md  ⭐ 推荐
│   └── KAFKA_USAGE.md
│
├── implementation/                 # 🔧 实现细节 (10 个文档)
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
├── tools/                          # 🛠️ 开发工具 (13 个文档)
│   ├── README.md
│   ├── AIR_USAGE_GUIDE.md          ⭐ 推荐
│   ├── AIR_WINDOWS_GUIDE.md
│   ├── AIR_CONFIG_README.md
│   ├── SWAGGER_GUIDE.md            ⭐ 推荐
│   ├── SWAGGER_SCRIPTS_GUIDE.md    ⭐ 推荐
│   ├── SWAGGER_FIX_GUIDE.md
│   ├── SWAGGER_SEPARATION_SUMMARY.md
│   ├── SWAGGER_CROSS_PLATFORM_SUMMARY.md
│   ├── SWAGGER_GO_FILE_FIX.md
│   ├── SWAGGER_GO_FILE_FIX_SUMMARY.md
│   ├── SWAGGER_SCRIPTS_IMPLEMENTATION.md
│   └── ENV_AND_SWAGGER_GUIDE.md
│
├── deployment/                     # 🚀 部署运维 (3 个文档)
│   ├── DEPLOYMENT.md               ⭐ 推荐
│   ├── API.md
│   └── RESPONSE_FORMAT.md
│
└── legacy/                         # 📦 遗留文档 (5 个文档)
    ├── CHECKLIST.md
    ├── FEATURE_SUMMARY.md
    ├── README_RESPONSE.md
    ├── SUMMARY.md
    └── OPTIMIZATION_RECOMMENDATIONS.md
```

## 📊 分类统计

| 目录 | 文档数 | 说明 |
|------|--------|------|
| **getting-started/** | 8 | 快速开始和入门指南 |
| **architecture/** | 6 | 系统架构和设计文档 |
| **features/** | 9 | 核心功能使用指南 |
| **implementation/** | 10 | 功能实现细节和重构记录 |
| **tools/** | 13 | 开发工具使用指南 |
| **deployment/** | 3 | 部署运维相关文档 |
| **legacy/** | 5 | 历史总结和遗留文档 |
| **根目录** | 2 | INDEX.md, PROJECT_ORGANIZATION.md |
| **总计** | **56** | 包含子目录 README |

## 🎯 分类说明

### 📁 getting-started/ - 快速开始

**目标用户**: 新手、第一次使用项目的开发者

**包含内容**:
- 快速启动指南
- API 测试示例
- 详细启动步骤
- 环境配置说明

**推荐文档**:
- `QUICK_START.md` - 最推荐，从这里开始
- `QUICK_TEST_GUIDE.md` - 测试 API

---

### 📁 architecture/ - 架构设计

**目标用户**: 需要了解系统设计的开发者、架构师

**包含内容**:
- 系统架构说明
- 项目结构详解
- 前后台分离设计
- API 设计对比

**推荐文档**:
- `ARCHITECTURE.md` - 系统架构
- `SEPARATED_SERVICES.md` - 服务分离

---

### 📁 features/ - 功能指南

**目标用户**: 需要使用特定功能的开发者

**包含内容**:
- RBAC 权限系统
- 限流功能
- 请求追踪
- 数据库迁移
- 监控系统
- 链路追踪
- 消息队列

**推荐文档**:
- `RBAC_GUIDE.md` - 权限系统
- `RATE_LIMIT_GUIDE.md` - 限流
- `MIGRATION_GUIDE.md` - 数据库
- `PROMETHEUS_MONITORING_GUIDE.md` - 监控
- `OPENTELEMETRY_TRACING_GUIDE.md` - 追踪

---

### 📁 implementation/ - 实现细节

**目标用户**: 需要深入了解实现的开发者、贡献者

**包含内容**:
- 各功能的实现细节
- 技术选型说明
- 重构记录
- 包名更新说明

**适合场景**:
- 想要贡献代码
- 需要修改现有功能
- 想要了解技术实现

---

### 📁 tools/ - 开发工具

**目标用户**: 开发者、需要配置开发环境的人员

**包含内容**:
- Air 热重载配置
- Swagger 文档生成
- 跨平台脚本支持
- 工具问题修复

**推荐文档**:
- `AIR_USAGE_GUIDE.md` - 热重载
- `SWAGGER_GUIDE.md` - API 文档
- `SWAGGER_SCRIPTS_GUIDE.md` - 脚本使用

---

### 📁 deployment/ - 部署运维

**目标用户**: 运维人员、需要部署项目的开发者

**包含内容**:
- 生产环境部署指南
- API 接口文档
- 响应格式规范

**推荐文档**:
- `DEPLOYMENT.md` - 部署指南

---

### 📁 legacy/ - 遗留文档

**目标用户**: 需要查看历史文档的人员

**包含内容**:
- 开发检查清单
- 功能总结
- 优化建议
- 历史总结文档

**说明**: 这些文档大部分已被新文档替代，保留作为参考。

---

## 📖 如何查找文档

### 方法 1：使用文档索引 ⭐ 推荐

查看 [`INDEX.md`](INDEX.md) 获取完整的文档索引、分类和快速导航。

### 方法 2：按目录浏览

根据你的需求，直接进入对应的目录：

- 想快速开始 → [`getting-started/`](getting-started/)
- 想了解架构 → [`architecture/`](architecture/)
- 想使用功能 → [`features/`](features/)
- 想了解实现 → [`implementation/`](implementation/)
- 想配置工具 → [`tools/`](tools/)
- 想部署上线 → [`deployment/`](deployment/)

### 方法 3：查看子目录 README

每个子目录都有自己的 `README.md`，提供该目录的文档列表和使用指南。

## 🔄 路径更新

如果你的书签或链接指向旧路径，需要更新：

| 旧路径 | 新路径 |
|--------|--------|
| `docs/QUICK_START.md` | `docs/getting-started/QUICK_START.md` |
| `docs/ARCHITECTURE.md` | `docs/architecture/ARCHITECTURE.md` |
| `docs/RBAC_GUIDE.md` | `docs/features/RBAC_GUIDE.md` |
| `docs/SWAGGER_GUIDE.md` | `docs/tools/SWAGGER_GUIDE.md` |
| `docs/DEPLOYMENT.md` | `docs/deployment/DEPLOYMENT.md` |

## 💡 最佳实践

### 添加新文档时

1. 确定文档类型（指南/实现/工具等）
2. 放入对应的目录
3. 更新该目录的 `README.md`
4. 更新 `docs/INDEX.md`
5. 如果是重要文档，在主 `README.md` 中添加链接

### 查找文档时

1. 先查看 `INDEX.md` 快速定位
2. 进入对应目录查看 `README.md`
3. 选择需要的文档阅读

### 维护文档时

- 保持目录结构清晰
- 及时更新索引
- 标记推荐文档（⭐）
- 添加简短的文档描述

## 🎉 整理效果

### 整理前
- ❌ 53 个文档全在 `docs/` 根目录
- ❌ 难以查找需要的文档
- ❌ 不知道从哪里开始

### 整理后
- ✅ 按功能分类到 7 个子目录
- ✅ 每个目录有清晰的定位
- ✅ 提供完整的索引和导航
- ✅ 新手可以快速找到入门文档
- ✅ 开发者可以快速找到功能文档
- ✅ 维护者可以快速找到实现文档

## 📞 反馈

如果你觉得文档分类不合理，或者有更好的建议，欢迎提 Issue 或 PR！

---

**整理完成时间**: 2024-11-12  
**文档总数**: 56 个（包含子目录 README）  
**主要目录**: 7 个  
**目标**: 提升文档的可查找性和可维护性

