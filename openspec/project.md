# Project Context

## Purpose
基于 trpc-agent-go 创建一个带有 MCP (Model Context Protocol) 服务的智能笔记 Agent。用户只需提供一个网址,Agent 会自动:
1. 抓取网页内容
2. 使用 AI 理解并总结内容
3. 生成 Markdown 格式的结构化笔记
4. 保存到用户的 Obsidian 笔记库中

这是一个"一键存笔记"的工具,让信息收集和知识管理变得自动化。

## Tech Stack
- **Go** - 主要开发语言
- **trpc-agent-go** - Agent 框架
- **MCP (Model Context Protocol)** - 与 AI 模型的通信协议
- **gRPC/TRPC** - 服务通信框架

## Project Conventions

### Code Style
- 使用 `gofmt` 进行代码格式化
- 遵循 Go 官方编码规范
- 包名使用小写单词,不使用下划线或驼峰
- 导出的函数和类型必须添加注释
- 错误处理使用显式返回,不忽略错误

### Architecture Patterns
- **分层架构**: Handler → Service → Repository
- **依赖注入**: 使用接口和依赖注入提高可测试性
- **中间件模式**: 使用中间件处理横切关注点(日志、认证等)
- **插件化设计**: Agent 能力作为插件模块,易于扩展

### Testing Strategy
- 单元测试覆盖率目标: >80%
- 使用表驱动测试(table-driven tests)
- 关键业务逻辑必须有集成测试
- 使用 mock 接口进行依赖隔离

### Git Workflow
- **主分支**: `main` - 稳定版本
- **开发分支**: `develop` - 开发集成分支
- **特性分支**: `feature/<capability-name>` - 新功能开发
- **修复分支**: `fix/<issue-name>` - Bug 修复
- **Commit 规范**: 使用 Conventional Commits 格式
  - `feat:` 新功能
  - `fix:` Bug 修复
  - `docs:` 文档更新
  - `refactor:` 代码重构
  - `test:` 测试相关

## Domain Context

### 核心概念
- **Agent**: 智能笔记助手,自动将网页内容转化为笔记
- **Note**: Markdown 格式的笔记,存储在 Obsidian vault 中
- **MCP Server**: 提供 `save_web_note` 工具给 Claude 等 AI 模型调用
- **Obsidian Vault**: 本地 Markdown 笔记库路径

### 业务流程
1. 用户通过 Claude(或其他支持 MCP 的 AI)发送网址
2. MCP 工具接收 URL 请求
3. Agent 抓取网页内容
4. 使用 AI 提取关键信息并总结
5. 生成 Markdown 笔记(包含标题、摘要、要点、标签等)
6. 保存到 Obsidian vault 的指定目录
7. 返回笔记路径给用户

## Important Constraints
- **性能**: 单次笔记生成应在 30 秒内完成
- **可靠性**: 失败重试机制,最多 3 次重试
- **安全性**: 验证输入 URL,防止 SSRF 攻击
- **资源限制**: 内存使用 < 500MB,支持并发处理

## External Dependencies
- **AI 模型服务**: Claude API (用于内容理解和总结)
- **内容抓取**: HTTP 客户端库(如 colly)用于网页抓取
- **Markdown 处理**: goldmark 或类似库
- **Obsidian**: 本地文件系统,直接写入 Markdown 文件
- **日志**: 结构化日志库(如 zap)
