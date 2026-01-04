.PHONY: help dev build run test clean install-tools lint fmt deps cli build-release

# 读取版本号
VERSION ?= $(shell cat .version 2>/dev/null || echo "dev")
GIT_COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || date)
LDFLAGS = -ldflags "-X 'github.com/fromsko/krio/app.Version=$(VERSION)' -X 'github.com/fromsko/krio/app.Commit=$(GIT_COMMIT)' -X 'github.com/fromsko/krio/app.BuildDate=$(BUILD_DATE)'"

# 默认目标
help:
	@echo "可用命令:"
	@echo "  make cli         - 构建 CLI 工具"
	@echo "  make build-release - 构建发布版本 (带版本信息)"
	@echo "  make dev         - 启动开发模式 (热重载)"
	@echo "  make build       - 构建演示程序"
	@echo "  make run         - 运行演示程序"
	@echo "  make test        - 运行测试"
	@echo "  make clean       - 清理构建文件"
	@echo "  make fmt         - 格式化代码"
	@echo "  make lint        - 代码检查"
	@echo "  make deps        - 安装依赖"
	@echo "  make install-tools- 安装开发工具"

# 构建 CLI 工具
cli:
	@echo "🔨 构建 CLI 工具..."
	@go build $(LDFLAGS) -o krio.exe .
	@echo "✅ 构建完成: krio.exe"
	@echo ""
	@echo "📝 使用方法:"
	@echo "  ./krio.exe init              # 初始化配置"
	@echo "  ./krio.exe run -u <url>      # 处理单个 URL"
	@echo "  ./krio.exe run -r <file>     # 批量处理文件"
	@echo "  ./krio.exe cache stats       # 查看缓存统计"
	@echo "  ./krio.exe version           # 查看版本信息"
	@echo ""

# 构建发布版本
build-release:
	@echo "🔨 构建发布版本..."
	@echo "版本: $(VERSION)"
	@echo "提交: $(GIT_COMMIT)"
	@echo "日期: $(BUILD_DATE)"
	@go build $(LDFLAGS) -o krio .
	@echo "✅ 构建完成: krio"
	@./krio version

# 开发模式 (热重载)
dev:
	@echo "🚀 启动开发模式..."
	@air

# 构建演示程序
build:
	@echo "🔨 构建演示程序..."
	@go build -o bin/server.exe ./cmd/server
	@echo "✅ 构建完成: bin/server.exe"

# 运行演示程序
run: build
	@echo "▶️  运行演示程序..."
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
