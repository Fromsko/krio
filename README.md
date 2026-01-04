# Krio - 智能网页笔记 Agent

> 一个基于 trpc-agent-go 框架的智能笔记助手,自动抓取网页内容并生成结构化笔记保存到 Obsidian。

## ✨ 特性

- 🤖 **AI 驱动**: 使用智云 GLM-4 模型进行内容理解和总结
- 🌐 **网页抓取**: 自动提取网页核心内容,去除广告和无关元素
- 📝 **Markdown 笔记**: 生成格式良好的 Markdown 笔记,包含 frontmatter
- 🏷️ **智能标签**: AI 自动生成相关标签,便于分类和检索
- 🔒 **安全防护**: URL 验证和 SSRF 防护
- ⚡ **高性能**:
  - 支持并发处理多个 URL (5 倍速度提升)
  - 智能缓存机制 (100 倍缓存命中速度)
  - 可配置的并发数和 TTL
- 🔧 **易于配置**: YAML 配置文件,支持环境变量覆盖
- 📦 **开箱即用**: 提供预编译二进制文件,无需配置 Go 环境

## 📦 下载预编译版本

访问 [Releases 页面](https://github.com/fromsko/krio/releases) 下载适合你系统的版本：

- **Linux AMD64**: `krio-*-linux-amd64.tar.gz`
- **Linux ARM64**: `krio-*-linux-arm64.tar.gz`
- **macOS Intel**: `krio-*-darwin-amd64.tar.gz`
- **macOS Apple Silicon**: `krio-*-darwin-arm64.tar.gz`
- **Windows**: `krio-*-windows-amd64.zip`

或者使用 `go install` 安装：

```bash
go install github.com/fromsko/krio@latest
```

## 🚀 快速开始

### 前置要求

- Go 1.21+
- Bun (用于 Obsidian MCP 服务器)
- 智云 API Key

### 安装

```bash
# 克隆仓库
git clone https://github.com/fromsko/krio.git
cd krio

# 安装依赖
go mod download

# 构建 CLI 工具
make cli

# 或使用 go build
go build -o krio.exe .

# 初始化配置
./krio.exe init
```

### 配置

编辑生成的配置文件 (`~/.krio.yaml` 或 `config.yaml`):

```yaml
model:
  api_key: "your-api-key"  # 填入你的智云 API Key
```

### 开发模式

```bash
# 安装开发工具
make install-tools

# 启动热重载开发模式
make dev
```

## 📖 使用方法

### CLI 命令行工具

```bash
# 查看帮助
./krio.exe --help

# 初始化配置
./krio.exe init

# 处理单个 URL
./krio.exe run -u https://example.com

# 批量处理文件
./krio.exe run -r urls.txt

# 自定义标签和文件夹
./krio.exe run -u https://example.com -t "tech,ai" -f "Articles"

# 查看缓存统计
./krio.exe cache stats

# 清空缓存
./krio.exe cache clear

# 查看版本信息
./krio.exe version
```

### 单个 URL 处理 (旧方式)

```bash
# 运行演示程序
./bin/server.exe
```

### 批量 URL 处理 (高性能模式)

```bash
# 编译批量演示程序
go build -o batch_demo.exe ./cmd/batch_demo

# 运行批量处理 (自动并发 + 缓存)
./batch_demo.exe
```

### 作为 Function Tool 使用

```go
package main

import (
    "context"
    "github.com/fromsko/krio/internal/tool"
)

// 创建工具
webNoteTool, _ := tool.NewSaveWebNoteTool(cfg)

// 调用工具
req := tool.SaveWebNoteRequest{
    URL:    "https://example.com/article",
    Tags:   []string{"技术", "AI"},
    Folder: "Articles",
}

resp, _ := webNoteTool.SaveWebNote(ctx, req)
fmt.Println(resp.Content)
```

### 批量处理 (并发 + 缓存)

```go
// 批量处理多个 URL (自动并发)
urls := []string{
    "https://example.com/article1",
    "https://example.com/article2",
    "https://example.com/article3",
}

responses := webNoteTool.SaveWebNoteBatch(ctx, urls, tags, "Articles")

// 处理结果
for i, resp := range responses {
    if resp.Success {
        fmt.Printf("✅ %s: %s\n", urls[i], resp.FilePath)
    }
}
```

## ⚙️ 配置说明

### 配置文件

主配置文件: `config.yaml`

```yaml
# 模型配置
model:
  api_key: "your-api-key"
  base_url: "https://open.bigmodel.cn/api/coding/paas/v4"
  model_name: "glm-4.7"
  temperature: 0.7
  max_tokens: 4096

# Obsidian MCP 配置
obsidian_mcp:
  transport: "stdio"
  command: "bun"
  args:
    - "x"
    - "--no-cache"
    - "@andysama/obsidian-mcp-server"
    - "--vault"
    - "/path/to/your/vault"  # 修改为你的 Obsidian vault 路径
  timeout: 30s

# 应用配置
app:
  name: "Web Note Agent"
  version: "1.0.0"
  debug: true

# 网页抓取配置
scraper:
  user_agent: "Mozilla/5.0..."
  timeout: 15s
  max_retries: 3
  retry_delay: 1000ms
  # 性能优化配置
  enable_cache: true        # 启用缓存
  cache_ttl: 1h            # 缓存过期时间
  max_concurrency: 5       # 最大并发数

# 笔记生成配置
note:
  default_folder: "Inbox"
  filename_template: "{{title}}-{{timestamp}}"
  add_timestamp: true
```

### 环境变量

```bash
export MODEL_API_KEY="your-api-key"
export MODEL_BASE_URL="https://open.bigmodel.cn/api/coding/paas/v4"
export MODEL_NAME="glm-4.7"
```

## 📁 项目结构

```
krio/
├── cmd/
│   └── server/           # 主程序入口
├── internal/
│   ├── config/          # 配置管理
│   ├── scraper/         # 网页抓取
│   ├── summarizer/      # AI 总结
│   ├── note/            # 笔记生成
│   └── tool/            # Function Tool
├── pkg/
│   └── logger/          # 日志模块
├── config.yaml          # 实际配置 (不提交)
├── config.example.yaml  # 示例配置
├── Makefile             # 构建脚本
├── .air.toml            # 热重载配置
└── README.md            # 本文档
```

## 🔧 开发

### Makefile 命令

```bash
make dev         # 开发模式 (热重载)
make build       # 构建
make run         # 运行
make test        # 测试
make clean       # 清理
make fmt         # 格式化
make lint        # 代码检查
make deps        # 安装依赖
```

### 代码规范

- 遵循 Go 标准编码规范
- 使用 `gofmt` 格式化代码
- 添加单元测试
- 编写清晰的注释

## 📝 生成的笔记格式

```markdown
---
title: 文章标题
source: https://example.com/article
date: 2026-01-05T12:00:00
tags: ["tag1", "tag2", "tag3"]
filename: article-title-2026-01-05-120000
created_at: 2026-01-05T12:00:00
updated_at: 2026-01-05T12:00:00
id: unique-id
---

# 文章标题

> 一句话概括文章核心内容

## 核心要点

- 要点 1
- 要点 2
- 要点 3
```

## 🔐 安全性

- ✅ URL 验证,防止 SSRF 攻击
- ✅ 私有地址检测
- ✅ 配置文件包含敏感信息,已加入 `.gitignore`
- ✅ 支持环境变量覆盖配置

## 🛠️ 技术栈

- **语言**: Go 1.21+
- **框架**: trpc-agent-go
- **LLM**: 智云 GLM-4
- **网页抓取**: Colly
- **日志**: Zap
- **配置**: YAML

## 📚 相关文档

- [性能优化文档](docs/PERFORMANCE.md) - 并发处理和缓存机制详解
- [trpc-agent-go 文档](references/technical/Trpc-agent-go.md)
- [API 配置参考](references/configuration/api-keys.md)
- [项目规范](openspec/project.md)

## 📄 License

Apache License 2.0

## 🤝 贡献

欢迎提交 Issue 和 Pull Request!

## 📮 联系方式

- GitHub: [@fromsko](https://github.com/fromsko)

---

**Generated with ❤️ by [Claude Code](https://claude.com/claude-code)**
