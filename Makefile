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

migrate: ## 运行数据库迁移
	go run cmd/frontend/*.go migrate

.DEFAULT_GOAL := help
