# Implementation Tasks

## 1. Project Setup
- [x] 1.1 Initialize Go module (`go mod init github.com/yourusername/krio`)
- [x] 1.2 Add dependencies:
  - colly (web scraping)
  - langchaingo (LLM integration)
  - zap (logging)
  - yaml (configuration)
- [x] 1.3 Set up project structure following Go best practices
- [x] 1.4 Create configuration file structure (config.yaml)

## 2. Web Scraping Module
- [x] 2.1 Create `internal/scraper/fetcher.go`
  - 实现网页内容抓取
  - 移除广告、导航等非必要元素
  - 提取标题和正文
- [x] 2.2 Add URL validation
  - 防止 SSRF 攻击
  - 验证 URL 格式
- [x] 2.3 Implement retry logic (最多 3 次)
- [x] 2.4 Add error handling和日志
- [ ] 2.5 Write unit tests

## 3. AI Summarization Module
- [x] 3.1 Create `internal/summarizer/summarizer.go`
  - 集成 GLM-4 API (通过 langchaingo)
  - 实现内容总结 prompt
  - 处理长内容分块
- [x] 3.2 Define summary structure
  - Title
  - One-sentence summary
  - Key points
  - Tags
- [x] 3.3 Add API key management
- [x] 3.4 Implement retry logic (最多 2 次)
- [ ] 3.5 Write unit tests with mock

## 4. Note Generation Module
- [x] 4.1 Create `internal/note/generator.go`
  - 生成 Markdown 格式
  - 包含 frontmatter (title, source, date, tags)
  - 格式化 key points
- [x] 4.2 Implement filename sanitization
- [x] 4.3 Add timestamp for duplicate filenames
- [ ] 4.4 Write unit tests

## 5. Function Tool (替代 MCP Server)
- [x] 5.1 Create `internal/tool/web_note_tool.go`
  - 创建 Function Tool
  - 定义工具 schema
- [x] 5.2 Implement tool handler
  - 接收 URL 参数
  - 调用 scraper → summarizer → generator
  - 返回结果
- [x] 5.3 Add error handling和日志
- [ ] 5.4 Write integration tests

## 6. Configuration
- [x] 6.1 Create `internal/config/config.go`
  - GLM API key
  - Obsidian vault path (可配置,或使用已有的 MCP)
  - 默认文件夹
  - 超时设置
- [x] 6.2 Support environment variables
- [x] 6.3 Create example config file

## 7. Main Entry Point
- [x] 7.1 Create `cmd/server/main.go`
  - 加载配置
  - 初始化模块
  - 启动应用
  - 优雅关闭处理
- [x] 7.2 Add signal handling (SIGTERM, SIGINT)
- [ ] 7.3 Add health check endpoint

## 8. Documentation
- [x] 8.1 Write README.md
  - 项目介绍
  - 安装步骤
  - 配置说明
  - 使用示例
- [x] 8.2 Document API/工具接口
- [x] 8.3 Add development guide

## 9. Testing
- [x] 9.1 End-to-end test with real URLs
- [ ] 9.2 Test error scenarios
- [ ] 9.3 Performance testing
- [ ] 9.4 Security testing (SSRF prevention)

## 10. Deployment
- [x] 10.1 Create build script (Makefile)
- [ ] 10.2 Add Docker support (optional)
- [x] 10.3 Write deployment guide

## 进度总结

### 已完成 ✅
- 项目初始化和依赖管理
- 开发环境配置 (air + Makefile)
- 配置管理系统 (YAML + 环境变量)
- 网页抓取模块 (Colly)
- AI 总结模块 (langchaingo + GLM-4)
- 笔记生成模块 (Markdown + frontmatter)
- Function Tool 实现
- 主程序入口
- README 文档

### 待完成 📝
- 单元测试
- 集成测试
- 性能测试
- Docker 支持
- 健康检查端点
- 与 Obsidian MCP 服务器集成
