# Change: Add Web Scraping and Note Generation

## Why
用户需要通过 Claude 调用一个工具,输入网址后自动将网页内容转化为 Obsidian 笔记。目前已有 Obsidian MCP 服务器负责存储,本项目需要实现网页抓取和 AI 总结功能。

## What Changes
- 实现网页内容抓取模块(使用 colly 或类似库)
- 集成 Claude API 进行内容分析和总结
- 创建 MCP 工具 `save_web_note`
- 生成结构化的 Markdown 笔记格式
- 调用现有的 Obsidian MCP 服务器保存笔记

## Impact
- Affected specs: `web-note-agent`
- Affected code:
  - 新增: `internal/scraper` - 网页抓取模块
  - 新增: `internal/summarizer` - AI 总结模块
  - 新增: `internal/mcp/server` - MCP 服务器
  - 新增: `internal/note` - 笔记生成模块
  - 配置: Claude API key 和 Obsidian vault 路径
