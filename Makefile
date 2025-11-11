.PHONY: help build run test clean wire deps docker-up docker-down frontend backend

help: ## 显示帮助信息
	@echo "可用的命令："
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

deps: ## 安装依赖
	go mod download
	go mod tidy

wire: ## 生成 Wire 依赖注入代码
	cd cmd/frontend && wire
	cd cmd/backend && wire
	@echo "✅ Wire 代码生成完成"

build: ## 编译前后台服务
	go build -o bin/frontend cmd/frontend/*.go
	go build -o bin/backend cmd/backend/*.go
	@echo "✅ 前后台服务构建完成"

build-frontend: ## 编译前台服务
	go build -o bin/frontend cmd/frontend/*.go
	@echo "✅ 前台服务构建完成"

build-backend: ## 编译后台服务
	go build -o bin/backend cmd/backend/*.go
	@echo "✅ 后台服务构建完成"

run-frontend: ## 运行前台服务（端口 8080）
	./bin/frontend

run-backend: ## 运行后台服务（端口 8081）
	./bin/backend

run: ## 同时运行前后台服务
	@echo "启动前台服务（端口 8080）..."
	./bin/frontend &
	@echo "启动后台服务（端口 8081）..."
	./bin/backend &
	@echo "✅ 前后台服务已启动"
	@echo "前台: http://localhost:8080"
	@echo "后台: http://localhost:8081"

dev-frontend: ## 开发模式运行前台（热重载）
	air -c .air-frontend.toml

dev-backend: ## 开发模式运行后台（热重载）
	air -c .air-backend.toml

test: ## 运行测试
	go test -v -cover ./...

clean: ## 清理构建文件
	rm -rf bin/
	rm -rf logs/
	go clean
	@echo "✅ 清理完成"

docker-up: ## 启动 Docker 服务 (MySQL, Redis, Kafka)
	docker-compose up -d
	@echo "✅ Docker 服务已启动"

docker-down: ## 停止 Docker 服务
	docker-compose down
	@echo "✅ Docker 服务已停止"

# ==================== 数据库迁移 ====================

migrate-up: ## 执行所有待执行的迁移
	@./scripts/migrate.sh up

migrate-down: ## 回滚一个迁移版本
	@./scripts/migrate.sh down

migrate-version: ## 查看当前迁移版本
	@./scripts/migrate.sh version

migrate-force: ## 强制设置迁移版本 (用法: make migrate-force VERSION=1)
	@./scripts/migrate.sh force

migrate-goto: ## 迁移到指定版本 (用法: make migrate-goto VERSION=3)
	@./scripts/migrate.sh goto

migrate-create: ## 创建新的迁移文件 (用法: make migrate-create NAME=add_user_phone)
	@./scripts/migrate.sh create

migrate-drop: ## 删除所有表 (危险操作)
	@./scripts/migrate.sh drop

# Swagger 文档生成
swag-frontend: ## 生成前台 Swagger 文档
	swag init -g cmd/frontend/main.go -o cmd/frontend/docs --parseDependency --parseInternal

swag-backend: ## 生成后台 Swagger 文档
	swag init -g cmd/backend/main.go -o cmd/backend/docs --parseDependency --parseInternal

swag: swag-frontend swag-backend ## 生成所有 Swagger 文档
	@echo "✅ Swagger 文档生成完成"

.DEFAULT_GOAL := help
