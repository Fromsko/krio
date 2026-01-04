# Krio - 智能网页笔记 Agent

> 一个基于 trpc-agent-go 框架的智能笔记助手,自动抓取网页内容并生成结构化笔记保存到 Obsidian。

## ✨ 特性

- 🤖 **AI 驱动**: 使用智云 GLM-4 模型进行内容理解和总结
- 🌐 **网页抓取**: 自动提取网页核心内容,去除广告和无关元素
- 📝 **Markdown 笔记**: 生成格式良好的 Markdown 笔记,包含 frontmatter
- 🏷️ **智能标签**: AI 自动生成相关标签,便于分类和检索
- 🔒 **安全防护**: URL 验证和 SSRF 防护
- ⚡ **高性能**: Go 语言实现,支持并发处理
- 🔧 **易于配置**: YAML 配置文件,支持环境变量覆盖

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

# 配置 API Key
cp config.example.yaml config.yaml
# 编辑 config.yaml,填入你的 API Key

# 构建
make build

# 运行
./bin/server.exe
```

### 开发模式

```bash
# 安装开发工具
make install-tools

# 启动热重载开发模式
make dev
```

## 📖 使用方法

### 命令行使用

```bash
# 直接运行程序
./bin/server.exe

# 程序会演示抓取 example.com 并生成笔记
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
    - "D:/notes/Fromsko"
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
