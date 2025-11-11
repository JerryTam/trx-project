package main

// @title TRX Project - 后台 API
// @version 1.0
// @description 基于 Gin 框架的现代化 Go Web 服务 - 后台管理接口
// @description 面向管理员的 API 服务，提供用户管理、数据统计等功能
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8081
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description 输入 "Bearer " + 管理员 JWT Token (需要 admin 或 superadmin 角色)

// @tag.name 用户管理
// @tag.description 后台用户管理相关接口

// @tag.name 统计信息
// @tag.description 数据统计相关接口

