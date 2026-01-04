.PHONY: help dev build run test clean install-tools lint fmt deps

# 默认目标
help:
	@echo "可用命令:"
	@echo "  make dev         - 启动开发模式 (热重载)"
	@echo "  make build       - 构建项目"
	@echo "  make run         - 运行程序"
	@echo "  make test        - 运行测试"
	@echo "  make clean       - 清理构建文件"
	@echo "  make fmt         - 格式化代码"
	@echo "  make lint        - 代码检查"
	@echo "  make deps        - 安装依赖"
	@echo "  make install-tools- 安装开发工具"

# 开发模式 (热重载)
dev:
	@echo "🚀 启动开发模式..."
	@air

# 构建
build:
	@echo "🔨 构建项目..."
	@go build -o bin/server.exe ./cmd/server
	@echo "✅ 构建完成: bin/server.exe"

# 运行
run: build
	@echo "▶️  运行程序..."
	@./bin/server.exe

# 测试
test:
	@echo "🧪 运行测试..."
	@go test -v -race -cover ./...

# 清理
clean:
	@echo "🧹 清理构建文件..."
	@if exist tmp rmdir /s /q tmp
	@if exist bin rmdir /s /q bin
	@echo "✅ 清理完成"

# 格式化代码
fmt:
	@echo "✨ 格式化代码..."
	@go fmt ./...
	@goimports -w .

# 代码检查
lint:
	@echo "🔍 代码检查..."
	@golangci-lint run ./...

# 安装依赖
deps:
	@echo "📦 安装依赖..."
	@go mod download
	@go mod tidy

# 安装开发工具
install-tools:
	@echo "🛠️  安装开发工具..."
	@go install github.com/cosmtrek/air@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "✅ 开发工具安装完成"

# 初始化项目
init: install-tools deps
	@echo "🎉 项目初始化完成!"
	@echo "运行 'make dev' 开始开发"
