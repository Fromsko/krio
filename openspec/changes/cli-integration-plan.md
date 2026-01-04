# CLI 工具集成实现计划

## 变更信息

- **创建时间**: 2026-01-05
- **状态**: 📋 规划中
- **优先级**: 高
- **复杂度**: 中高
- **预计工期**: 2-3 周

## 概述

将 Krio 从简单演示程序改造为功能完整的 CLI 工具，支持：

1. **命令行接口**: 使用 [spf13/cobra](https://github.com/spf13/cobra) 构建
2. **配置初始化**: 自动生成默认配置文件
3. **URL 处理**: 支持命令行指定 URL
4. **文件读取**: 支持 `.txt` 和 `.md` 文件批量导入 URL
5. **智能提取**: AI 自动提取任务和学习内容

## 背景与动机

### 当前限制

1. **硬编码演示**: URL 在代码中硬编码
2. **无 CLI 接口**: 只能运行演示，无法灵活使用
3. **配置复杂**: 需要手动编辑配置文件
4. **批量处理困难**: 需要编写代码才能批量处理

### 用户需求

1. **快速上手**: 开箱即用，自动配置
2. **灵活使用**: 支持多种输入方式
3. **批量处理**: 方便处理大量 URL
4. **智能提取**: AI 自动解析文件内容

## 技术方案

### 1. 技术栈

**核心库**: [spf13/cobra](https://github.com/spf13/cobra)

**选择理由**:
- ✅ Go 生态最流行的 CLI 框架
- ✅ 功能强大 (子命令、参数、标志位)
- ✅ 自动生成文档和补全
- ✅ 广泛使用 (kubectl, docker-compose, hugo)

**依赖库**:
- `github.com/spf13/cobra` - CLI 框架
- `github.com/spf13/viper` - 配置管理 (可选)
- 现有依赖不变

### 2. 命令结构设计

```bash
krio
├── init         # 初始化配置
├── run          # 运行主程序
│   ├── -u, --url       # 单个 URL
│   ├── -r, --require   # 需求文件 (.txt/.md)
│   ├── -t, --tags      # 自定义标签
│   ├── -f, --folder    # 目标文件夹
│   └── --batch         # 批量模式
├── cache        # 缓存管理
│   ├── clear         # 清空缓存
│   └── stats         # 缓存统计
└── version      # 版本信息
```

### 3. 文件格式支持

#### 3.1 TXT 文件格式

```
# 支持：
# 1. 每行一个 URL
https://example.com/article1
https://example.com/article2
https://example.com/article3

# 2. 支持注释 (# 开头)
# 这是注释
https://example.com/article4

# 3. 支持空行（自动忽略）

https://example.com/article5
```

#### 3.2 Markdown 文件格式

```markdown
# 学习任务清单

## 前端框架
- [ ] 学习 React: https://react.dev/learn
- [ ] 学习 Vue: https://vuejs.org/guide/

## 后端框架
- [ ] 学习 Go: https://go.dev/doc/
- [ ] 学习 Python: https://docs.python.org/3/

## 其他资源
有用的博客: https://blog.example.com
```

**AI 提取规则**:
- 提取所有 HTTP/HTTPS URL
- 忽略代码块内的 URL
- 支持列表项、链接、纯文本中的 URL
- 自动提取上下文作为标签

### 4. AI 智能提取

#### 4.1 提取流程

```go
// 读取文件内容
content := readFile(filePath)

// AI 提取 URL 和元数据
extracted := ai.ExtractURLs(content)

// 结构化结果
type ExtractedResult struct {
    URLs []URLMetadata
}

type URLMetadata struct {
    URL         string
    Title       string   // 从上下文提取
    Tags        []string // 从周围文本提取
    Priority    int      // 优先级
    Description string   // 描述
}
```

#### 4.2 AI 提示词

```go
prompt := `你是网页内容提取助手。请分析以下文本内容，提取所有需要学习或保存的 URL。

文本内容:
%s

请按以下 JSON 格式返回:
{
  "urls": [
    {
      "url": "https://example.com",
      "title": "网页标题或描述",
      "tags": ["标签1", "标签2"],
      "priority": 1,
      "description": "简短描述"
    }
  ]
}

规则:
1. url: 完整的 HTTP/HTTPS URL
2. title: 根据上下文提取标题或生成描述
3. tags: 从周围文本提取 2-5 个相关标签
4. priority: 1-5 (1=最高优先级)
5. description: 为什么需要学习/保存这个网页

只返回 JSON，不要其他内容。`
```

## 实现计划

### 阶段 1: Cobra 集成 (3-5 天)

#### 任务清单

- [ ] **1.1 项目结构调整**
  - [ ] 创建 `cmd/root.go` - 根命令
  - [ ] 重构 `cmd/server/main.go` → `cmd/run.go`
  - [ ] 创建 `cmd/init.go` - 初始化命令
  - [ ] 创建 `cmd/cache.go` - 缓存管理
  - [ ] 创建 `cmd/version.go` - 版本信息

- [ ] **1.2 核心命令实现**
  - [ ] `krio init` - 初始化配置
  - [ ] `krio run -u <url>` - 单 URL 运行
  - [ ] `krio run -r <file>` - 文件批量运行
  - [ ] `krio cache clear` - 清空缓存
  - [ ] `krio cache stats` - 缓存统计
  - [ ] `krio version` - 版本信息

- [ ] **1.3 参数和标志位**
  ```go
  // run 命令
  var (
      urlFile   string // -r, --require
      singleURL string // -u, --url
      tags      []string // -t, --tags
      folder    string // -f, --folder
      batch     bool // --batch
      config    string // -c, --config
  )
  ```

- [ ] **1.4 配置文件管理**
  - [ ] 检测配置文件是否存在
  - [ ] 自动生成默认配置
  - [ ] 支持自定义配置路径
  - [ ] 环境变量支持

#### 文件结构

```
cmd/
├── root.go           # 根命令 (新增)
├── init.go           # init 命令 (新增)
├── run.go            # run 命令 (重构)
├── cache.go          # cache 命令 (新增)
└── version.go        # version 命令 (新增)
```

#### 代码示例

**root.go**
```go
package cmd

import (
    "os"
    "github.com/spf13/cobra"
    "github.com/fromsko/krio/internal/config"
)

var cfgFile string

var rootCmd = &cobra.Command{
    Use:   "krio",
    Short: "智能网页笔记 Agent",
    Long: `Krio 是一个基于 AI 的智能网页笔记工具，
自动抓取网页内容并生成结构化笔记保存到 Obsidian。`,
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        os.Exit(1)
    }
}

func init() {
    cobra.OnInitialize(initConfig)
    rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "",
        "配置文件路径 (默认: $HOME/.krio.yaml)")
}
```

**init.go**
```go
package cmd

import (
    "fmt"
    "os"
    "github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
    Use:   "init",
    Short: "初始化配置文件",
    Long:  `在当前目录或用户目录创建默认配置文件。`,
    Run: func(cmd *cobra.Command, args []string) {
        // 检查配置文件是否存在
        if config.Exists() {
            fmt.Println("⚠️  配置文件已存在")
            return
        }

        // 创建默认配置
        if err := config.CreateDefault(); err != nil {
            fmt.Printf("❌ 创建配置失败: %v\n", err)
            return
        }

        fmt.Println("✅ 配置文件已创建")
        fmt.Println("📝 请编辑配置文件,填入 API Key")
    },
}

func init() {
    rootCmd.AddCommand(initCmd)
}
```

**run.go**
```go
package cmd

import (
    "context"
    "fmt"
    "github.com/spf13/cobra"
    "github.com/fromsko/krio/internal/config"
    "github.com/fromsko/krio/internal/tool"
)

var (
    urlFile   string
    singleURL string
    tags      []string
    folder    string
    batchMode bool
)

var runCmd = &cobra.Command{
    Use:   "run",
    Short: "运行网页笔记生成器",
    Long:  `从 URL 或文件批量生成网页笔记并保存到 Obsidian。`,
    Run: func(cmd *cobra.Command, args []string) {
        // 加载配置
        cfg, err := config.LoadDefault()
        if err != nil {
            fmt.Printf("❌ 加载配置失败: %v\n", err)
            return
        }

        // 创建工具
        ctx := context.Background()
        webNoteTool, err := tool.NewSaveWebNoteTool(ctx, cfg)
        if err != nil {
            fmt.Printf("❌ 创建工具失败: %v\n", err)
            return
        }

        // 根据参数执行
        switch {
        case singleURL != "":
            runSingleURL(ctx, webNoteTool, singleURL)
        case urlFile != "":
            runFile(ctx, webNoteTool, urlFile)
        default:
            fmt.Println("❌ 请指定 -u <url> 或 -r <file>")
            cmd.Help()
        }
    },
}

func init() {
    rootCmd.AddCommand(runCmd)
    runCmd.Flags().StringVarP(&singleURL, "url", "u", "",
        "单个 URL")
    runCmd.Flags().StringVarP(&urlFile, "require", "r", "",
        "需求文件 (.txt/.md)")
    runCmd.Flags().StringSliceVarP(&tags, "tags", "t", []string{},
        "自定义标签")
    runCmd.Flags().StringVarP(&folder, "folder", "f", "",
        "目标文件夹")
}
```

### 阶段 2: 文件解析 (2-3 天)

#### 任务清单

- [ ] **2.1 TXT 文件解析**
  - [ ] 实现 `parser/txt_parser.go`
  - [ ] 支持注释 (# 开头)
  - [ ] 支持空行过滤
  - [ ] URL 格式验证

- [ ] **2.2 Markdown 文件解析**
  - [ ] 实现 `parser/md_parser.go`
  - [ ] 提取链接语法: `[text](url)`
  - [ ] 提取纯文本 URL
  - [ ] 保留上下文信息

- [ ] **2.3 AI 智能提取**
  - [ ] 实现 `summarizer/extractor.go`
  - [ ] AI 提取 URL 和元数据
  - [ ] 自动生成标签和描述
  - [ ] 优先级判断

#### 文件结构

```
internal/parser/
├── parser.go         # 解析器接口
├── txt_parser.go     # TXT 解析器
├── md_parser.go      # Markdown 解析器
└── extractor.go      # AI 提取器

internal/summarizer/
├── extractor.go      # URL 提取 (新增)
└── summarizer.go     # (已有)
```

#### 代码示例

**parser.go**
```go
package parser

import (
    "io"
    "github.com/fromsko/krio/internal/summarizer"
)

// Parser 解析器接口
type Parser interface {
    Parse(r io.Reader) ([]string, error)
    ParseWithAI(r io.Reader, ai *summarizer.Extractor) ([]summarizer.URLMetadata, error)
}

// DetectFormat 检测文件格式
func DetectFormat(filename string) Parser {
    switch {
    case strings.HasSuffix(filename, ".txt"):
        return NewTxtParser()
    case strings.HasSuffix(filename, ".md"):
        return NewMdParser()
    default:
        return NewTxtParser() // 默认
    }
}
```

**txt_parser.go**
```go
package parser

import (
    "bufio"
    "io"
    "strings"
)

type TxtParser struct{}

func NewTxtParser() *TxtParser {
    return &TxtParser{}
}

// Parse 解析 TXT 文件，提取 URL
func (p *TxtParser) Parse(r io.Reader) ([]string, error) {
    scanner := bufio.NewScanner(r)
    var urls []string

    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())

        // 跳过空行和注释
        if line == "" || strings.HasPrefix(line, "#") {
            continue
        }

        // 验证 URL
        if isValidURL(line) {
            urls = append(urls, line)
        }
    }

    return urls, scanner.Err()
}
```

**extractor.go**
```go
package summarizer

import (
    "context"
    "fmt"
)

// URLMetadata URL 元数据
type URLMetadata struct {
    URL         string   `json:"url"`
    Title       string   `json:"title"`
    Tags        []string `json:"tags"`
    Priority    int      `json:"priority"`
    Description string   `json:"description"`
}

// Extractor AI 提取器
type Extractor struct {
    sum *Summarizer
}

func NewExtractor(sum *Summarizer) *Extractor {
    return &Extractor{sum: sum}
}

// ExtractFromText 从文本中提取 URL 和元数据
func (e *Extractor) ExtractFromText(ctx context.Context, content string) ([]URLMetadata, error) {
    prompt := e.buildExtractPrompt(content)

    // 调用 AI
    response, err := llms.GenerateFromSinglePrompt(ctx, e.sum.llm, prompt)
    if err != nil {
        return nil, fmt.Errorf("AI 提取失败: %w", err)
    }

    // 解析结果
    return e.parseExtractResult(response)
}
```

### 阶段 3: 增强功能 (2-3 天)

#### 任务清单

- [ ] **3.1 交互式模式**
  - [ ] 支持命令行交互输入
  - [ ] 实时进度显示
  - [ ] 错误重试机制

- [ ] **3.2 输出格式**
  - [ ] 表格输出 (批量结果)
  - [ ] JSON 输出 (API 集成)
  - [ ] 详细日志模式

- [ ] **3.3 命令补全**
  - [ ] Bash 自动补全
  - [ ] Zsh 自动补全
  - [ ] PowerShell 自动补全

- [ ] **3.4 配置验证**
  - [ ] 配置文件语法检查
  - [ ] API Key 验证
  - [ ] Obsidian 连接测试

#### 代码示例

**输出格式**
```go
// 表格输出
func printTableResults(responses []SaveWebNoteResponse) {
    table := tablewriter.NewWriter(os.Stdout)
    table.SetHeader([]string{"#", "URL", "Title", "Status", "Path"})

    for i, resp := range responses {
        status := "✅"
        if !resp.Success {
            status = "❌"
        }
        table.Append([]string{
            fmt.Sprintf("%d", i+1),
            truncate(url, 40),
            truncate(resp.Title, 30),
            status,
            resp.FilePath,
        })
    }

    table.Render()
}

// JSON 输出
func printJSONResults(responses []SaveWebNoteResponse) {
    json.NewEncoder(os.Stdout).Encode(responses)
}
```

### 阶段 4: 文档和测试 (2-3 天)

#### 任务清单

- [ ] **4.1 用户文档**
  - [ ] CLI 使用指南
  - [ ] 配置文件说明
  - [ ] 文件格式规范
  - [ ] 常见问题 FAQ

- [ ] **4.2 开发文档**
  - [ ] 命令开发指南
  - [ ] 插件系统设计
  - [ ] API 文档

- [ ] **4.3 测试**
  - [ ] 命令行测试
  - [ ] 文件解析测试
  - [ ] AI 提取测试
  - [ ] 集成测试

#### 文档结构

```
docs/
├── cli/
│   ├── user-guide.md       # 用户指南
│   ├── commands.md         # 命令参考
│   ├── file-formats.md     # 文件格式
│   └── faq.md              # 常见问题
└── development/
    ├── command-dev.md      # 命令开发
    └── testing.md          # 测试指南
```

## 使用示例

### 基础使用

```bash
# 1. 初始化配置
krio init

# 2. 编辑配置文件
vim ~/.krio.yaml

# 3. 运行单个 URL
krio run -u https://example.com

# 4. 批量处理 TXT 文件
krio run -r urls.txt

# 5. 批量处理 Markdown 文件
krio run -r learning-list.md

# 6. 自定义标签和文件夹
krio run -u https://example.com -t "tech,ai" -f "Articles"
```

### 高级使用

```bash
# 1. 查看缓存统计
krio cache stats

# 2. 清空缓存
krio cache clear

# 3. JSON 输出
krio run -r urls.txt --output json

# 4. 详细日志
krio run -r urls.txt --log-level debug

# 5. 自定义配置文件
krio -c /path/to/config.yaml run -r urls.txt
```

### 文件格式示例

**urls.txt**
```txt
# 学习资源
https://go.dev/doc/
https://vuejs.org/guide/
https://react.dev/learn
```

**learning.md**
```markdown
# 学习计划

## 前端
- [ ] [React 文档](https://react.dev/learn)
- [ ] [Vue 指南](https://vuejs.org/guide/)

## 后端
- [ ] Go 官方文档: https://go.dev/doc/

## 其他资源
详见: https://example.com/resources
```

## API 设计

### 命令行接口

```bash
# 全局参数
krio [global flags] command [command flags] [arguments]

# 全局标志
--config string   # 配置文件路径
--log-level       # 日志级别
--verbose         # 详细输出
--version         # 版本信息
```

### 命令参考

```bash
# init 命令
krio init [flags]
Flags:
  --force      # 覆盖已存在的配置
  --path       # 指定配置文件路径

# run 命令
krio run [flags]
Flags:
  -u, --url string         # 单个 URL
  -r, --require string     # 需求文件
  -t, --tags strings       # 自定义标签
  -f, --folder string      # 目标文件夹
  --batch                  # 批量模式
  --output string          # 输出格式 (table/json)
  --concurrency int        # 并发数

# cache 命令
krio cache [subcommand]
Subcommands:
  clear                   # 清空缓存
  stats                   # 缓存统计

# version 命令
krio version [flags]
Flags:
  --json                  # JSON 输出
```

## 配置管理

### 配置文件优先级

```
命令行参数 > 环境变量 > 配置文件 > 默认值
```

### 配置文件位置

```
1. --config 指定的路径
2. ./config.yaml (当前目录)
3. ~/.krio.yaml (用户目录)
4. /etc/krio/config.yaml (系统目录)
```

### 环境变量

```bash
export KRio_CONFIG=/path/to/config.yaml
export MODEL_API_KEY=your-api-key
export KRio_LOG_LEVEL=debug
```

## 性能考虑

### 批量处理优化

```go
// 智能批量大小
const (
    MinBatchSize = 1
    MaxBatchSize = 50
    DefaultConcurrency = 5
)

// 根据 URL 数量自动调整
func calculateConcurrency(urlCount int) int {
    switch {
    case urlCount < 5:
        return 1
    case urlCount < 20:
        return 3
    case urlCount < 50:
        return 5
    default:
        return 10
    }
}
```

### 内存管理

```go
// 流式读取大文件
func parseLargeFile(filePath string) <-string {
    ch := make(chan string)

    go func() {
        file, _ := os.Open(filePath)
        defer file.Close()

        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
            ch <- scanner.Text()
        }
        close(ch)
    }()

    return ch
}
```

## 测试策略

### 单元测试

```go
// parser/txt_parser_test.go
func TestTxtParser_Parse(t *testing.T) {
    input := `# 注释
https://example.com
https://test.com`

    parser := NewTxtParser()
    urls, _ := parser.Parse(strings.NewReader(input))

    assert.Equal(t, 2, len(urls))
    assert.Equal(t, "https://example.com", urls[0])
}
```

### 集成测试

```go
// cmd/run_test.go
func TestRunCmd(t *testing.T) {
    // 创建临时配置
    cfg := createTempConfig(t)
    defer os.Remove(cfg)

    // 测试单 URL
    output := executeCommand(t, "run", "-u", "https://example.com")
    assert.Contains(t, output, "✅")
}
```

## 风险与挑战

### 技术风险

1. **AI 提取准确性**
   - 风险: AI 可能误提取或遗漏 URL
   - 缓解: 结合正则表达式和 AI 提取，提高准确性

2. **大文件处理**
   - 风险: 大文件可能导致内存溢出
   - 缓解: 流式读取，批量处理

3. **并发控制**
   - 风险: 高并发可能触发 API 限流
   - 缓解: 自适应并发控制

### 兼容性风险

1. **不同平台**
   - Windows: 路径分隔符、命令补全
   - Linux/macOS: 权限、符号链接

2. **配置迁移**
   - 旧版本配置兼容性
   - 平滑升级路径

## 时间估算

| 阶段 | 任务 | 预计时间 |
|------|------|----------|
| 阶段 1 | Cobra 集成 | 3-5 天 |
| 阶段 2 | 文件解析 | 2-3 天 |
| 阶段 3 | 增强功能 | 2-3 天 |
| 阶段 4 | 文档和测试 | 2-3 天 |
| **总计** | | **9-14 天** |

## 成功标准

### 功能完整性

- [x] 支持所有计划的命令
- [x] 支持 TXT 和 Markdown 文件
- [x] AI 智能提取功能
- [x] 完整的错误处理

### 用户体验

- [x] 开箱即用 (自动配置)
- [x] 清晰的错误提示
- [x] 实时进度反馈
- [x] 完善的文档

### 性能指标

- [x] 单个 URL: <5 秒
- [x] 批量 10 个: <30 秒
- [x] 内存占用: <200MB
- [x] 缓存命中: <100ms

## 后续优化

### 短期 (1-2 周)

- [ ] 添加配置向导 (交互式)
- [ ] 支持更多文件格式 (JSON, CSV)
- [ ] 添加进度条显示

### 中期 (1-2 月)

- [ ] 支持插件系统
- [ ] Web UI 界面
- [ ] 定时任务支持

### 长期 (3-6 月)

- [ ] 分布式处理
- [ ] 云端同步
- [ ] 多语言支持

## 参考资料

### Cobra 文档

- [Cobra 官方文档](https://github.com/spf13/cobra/blob/main/README.md)
- [Cobra 用户指南](https://github.com/spf13/cobra/blob/main/site/content/user_guide.md)
- [命令行最佳实践](https://clig.dev/)

### 项目文档

- [OpenSpec 项目规范](../project.md)
- [性能优化文档](../../docs/PERFORMANCE.md)

---

**状态**: 📋 待评审
**下一步**: 阶段 1 - Cobra 集成
**预计完成**: 2026-01-19
