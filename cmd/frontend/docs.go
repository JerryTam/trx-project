package main

// @title TRX Project - 前台 API
// @version 1.0
// @description 基于 Gin 框架的现代化 Go Web 服务 - 前台接口
// @description 面向最终用户的 API 服务，提供用户注册、登录、个人信息管理等功能
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description 输入 "Bearer " + JWT Token

// @tag.name 公开接口
// @tag.description 无需认证的公开接口

// @tag.name 用户接口
// @tag.description 需要用户认证的接口
