# tRPC-Agent-Go 技术特性

## 模型支持

### 多平台兼容性

tRPC-Agent-Go 提供统一的模型接口，支持主流 LLM 平台：

| 平台 | 模型示例 | 特性 |
|------|----------|------|
| OpenAI | gpt-4o, gpt-4o-mini | 完整功能支持 |
| DeepSeek | deepseek-chat, deepseek-reasoner | 推理模式支持 |
| 腾讯混元 | hunyuan-2.0-thinking | 思考模式支持 |
| Anthropic | claude-3-5-sonnet, claude-3-5-haiku | 高级推理支持 |
| 其他 | OpenAI 兼容 API | 通过配置适配 |

### 统一接口抽象

```go
type Model interface {
    GenerateContent(ctx context.Context, request *Request) (<-chan *Response, error)
    Info() Info
}
```

### 高级模型功能

#### 1. 流式响应
- 实时增量输出
- 低延迟交互体验
- 自动错误处理和恢复

#### 2. 模型动态切换
- **Agent 级别切换**: 影响所有后续请求
- **请求级别切换**: 仅影响单次请求

```go
// Agent 级别切换
agent.SetModel(openai.New("gpt-4o"))

// 请求级别切换
runner.Run(ctx, userID, sessionID, message,
    agent.WithModelName("smart"))
```

#### 3. Token 裁剪 (Token Tailoring)
- 自动处理超长上下文
- 智能消息裁剪策略
- 可配置裁剪参数

```go
model := openai.New("deepseek-chat",
    openai.WithEnableTokenTailoring(true),
    openai.WithMaxInputTokens(10000),
)
```

#### 4. 批量处理 (Batch API)
- 异步批量处理
- 成本优化
- 完整的生命周期管理

#### 5. 重试机制
- 自动错误恢复
- 智能退避策略
- 可配置重试参数

## 工具系统

### 工具类型

#### Function Tools
- 直接调用 Go 函数
- 自动参数序列化
- 支持流式和非流式响应

```go
calculatorTool := function.NewFunctionTool(
    calculator,
    function.WithName("calculator"),
    function.WithDescription("执行数学运算"),
)
```

#### Agent Tool
- 将 Agent 包装为工具
- 支持流式内部转发
- 支持历史作用域控制

#### MCP ToolSet
- 基于 MCP 协议的外部工具
- 支持多种传输方式:
  - STDIO: 标准输入输出
  - SSE: Server-Sent Events
  - Streamable HTTP: 流式 HTTP
- 支持会话重连和动态发现

#### DuckDuckGo Tool
- 基于 DuckDuckGo API 的搜索
- 事实性、百科类信息检索

### 高级工具特性

#### 1. 流式工具支持
```go
type StreamableTool interface {
    StreamableCall(ctx context.Context, jsonArgs []byte) (*StreamReader, error)
    Tool
}
```

#### 2. 并行工具执行
- 多工具并发调用
- 性能优化
- 结果聚合

```go
agent := llmagent.New("assistant",
    llmagent.WithEnableParallelTools(true),
)
```

#### 3. 运行时工具过滤
- 动态控制工具可见性
- 成本优化
- 基于角色的访问控制

```go
// 只允许特定工具
filter := tool.NewIncludeToolNamesFilter("calculator", "time_tool")
runner.Run(ctx, userID, sessionID, message,
    agent.WithToolFilter(filter),
)
```

#### 4. 动态工具集管理
- 运行时添加/删除 ToolSet
- 无需重建 Agent
- 并发安全

```go
agent.AddToolSet(mcpToolSet)
agent.RemoveToolSet("old-toolset")
```

## 会话与记忆管理

### 会话服务

#### 存储方式
- **In-Memory**: 内存存储，适合单机部署
- **Redis**: 分布式存储，支持集群部署

#### 核心功能
- 会话创建与获取
- 事件存储与查询
- 状态管理
- TTL 支持

### 记忆系统

#### 记忆类型
- **长期记忆**: 用户偏好、历史信息
- **会话记忆**: 当前对话上下文
- **临时记忆**: 单次调用临时数据

#### 记忆操作
- 记忆存储与检索
- 记忆搜索与过滤
- 记忆更新与删除

### 状态注入

#### 占位符变量
- `{key}`: 会话状态值
- `{user:subkey}`: 用户态状态
- `{app:subkey}`: 应用态状态
- `{temp:subkey}`: 临时态状态

```go
llm := llmagent.New("agent",
    llmagent.WithInstruction("用户兴趣: {user:topics}"),
)
```

## 可观测性

### OpenTelemetry 集成
- 全链路追踪
- 性能监控
- 自定义指标

```go
// 启动 Langfuse 集成
clean, _ := langfuse.Start(ctx)
defer clean(ctx)
```

### 事件系统
- 实时事件流
- 结构化事件数据
- 完整的执行上下文

### 回调机制
- Agent 级回调
- Model 级回调
- Tool 级回调

```go
callbacks := agent.NewCallbacks()
callbacks.RegisterBeforeAgent(func(ctx context.Context, args *agent.BeforeAgentArgs) (*agent.BeforeAgentResult, error) {
    // 自定义逻辑
    return nil, nil
})
```

## 高级配置

### Provider 统一接口
- 简化多平台切换
- 统一配置选项
- 自动平台适配

```go
model, err := provider.Model(
    "openai",
    "deepseek-chat",
    provider.WithAPIKey(apiKey),
    provider.WithBaseURL(baseURL),
)
```

### 自定义 HTTP Header
- 网关支持
- 专有平台适配
- 代理环境支持

```go
model := openai.New("deepseek-chat",
    openai.WithHeaders(map[string]string{
        "X-Custom-Header": "custom-value",
        "X-Request-ID":    "req-123",
    }),
)
```

### Variant 优化
- 平台特有行为适配
- 自动 API 差异处理
- 环境变量自动配置

```go
// 混元平台适配
model := openai.New("hunyuan-model",
    openai.WithVariant(openai.VariantHunyuan),
)
```

## 性能优化

### 并发处理
- Go 协程并发
- 工具并行执行
- 异步事件处理

### 内存管理
- 自动资源清理
- 流式数据处理
- 缓冲区优化

### 网络优化
- 连接池复用
- 请求合并
- 智能重试

## 安全特性

### 权限控制
- 工具级权限
- 用户级隔离
- 会话级安全

### 输入验证
- 参数验证
- 类型检查
- 长度限制

### 错误处理
- 分层错误处理
- 敏感信息过滤
- 安全错误响应

## 相关文档

- [框架概述](./framework-overview.md)
- [核心架构](./core-architecture.md)
- [开发模式](./development-patterns.md)
- [高级功能](./advanced-features.md)
## 模型支持

### 多平台兼容性

tRPC-Agent-Go 提供统一的模型接口，支持主流 LLM 平台：

| 平台 | 模型示例 | 特性 |
|------|----------|------|
| OpenAI | gpt-4o, gpt-4o-mini | 完整功能支持 |
| DeepSeek | deepseek-chat, deepseek-reasoner | 推理模式支持 |
| 腾讯混元 | hunyuan-2.0-thinking | 思考模式支持 |
| Anthropic | claude-3-5-sonnet, claude-3-5-haiku | 高级推理支持 |
| 其他 | OpenAI 兼容 API | 通过配置适配 |

### 统一接口抽象

```go
type Model interface {
    GenerateContent(ctx context.Context, request *Request) (<-chan *Response, error)
    Info() Info
}
```

### 高级模型功能

#### 1. 流式响应
- 实时增量输出
- 低延迟交互体验
- 自动错误处理和恢复

#### 2. 模型动态切换
- **Agent 级别切换**: 影响所有后续请求
- **请求级别切换**: 仅影响单次请求

```go
// Agent 级别切换
agent.SetModel(openai.New("gpt-4o"))

// 请求级别切换
runner.Run(ctx, userID, sessionID, message,
    agent.WithModelName("smart"))
```

#### 3. Token 裁剪 (Token Tailoring)
- 自动处理超长上下文
- 智能消息裁剪策略
- 可配置裁剪参数

```go
model := openai.New("deepseek-chat",
    openai.WithEnableTokenTailoring(true),
    openai.WithMaxInputTokens(10000),
)
```

#### 4. 批量处理 (Batch API)
- 异步批量处理
- 成本优化
- 完整的生命周期管理

#### 5. 重试机制
- 自动错误恢复
- 智能退避策略
- 可配置重试参数

## 工具系统

### 工具类型

#### Function Tools
- 直接调用 Go 函数
- 自动参数序列化
- 支持流式和非流式响应

```go
calculatorTool := function.NewFunctionTool(
    calculator,
    function.WithName("calculator"),
    function.WithDescription("执行数学运算"),
)
```

#### Agent Tool
- 将 Agent 包装为工具
- 支持流式内部转发
- 支持历史作用域控制

#### MCP ToolSet
- 基于 MCP 协议的外部工具
- 支持多种传输方式:
  - STDIO: 标准输入输出
  - SSE: Server-Sent Events
  - Streamable HTTP: 流式 HTTP
- 支持会话重连和动态发现

#### DuckDuckGo Tool
- 基于 DuckDuckGo API 的搜索
- 事实性、百科类信息检索

### 高级工具特性

#### 1. 流式工具支持
```go
type StreamableTool interface {
    StreamableCall(ctx context.Context, jsonArgs []byte) (*StreamReader, error)
    Tool
}
```

#### 2. 并行工具执行
- 多工具并发调用
- 性能优化
- 结果聚合

```go
agent := llmagent.New("assistant",
    llmagent.WithEnableParallelTools(true),
)
```

#### 3. 运行时工具过滤
- 动态控制工具可见性
- 成本优化
- 基于角色的访问控制

```go
// 只允许特定工具
filter := tool.NewIncludeToolNamesFilter("calculator", "time_tool")
runner.Run(ctx, userID, sessionID, message,
    agent.WithToolFilter(filter),
)
```

#### 4. 动态工具集管理
- 运行时添加/删除 ToolSet
- 无需重建 Agent
- 并发安全

```go
agent.AddToolSet(mcpToolSet)
agent.RemoveToolSet("old-toolset")
```

## 会话与记忆管理

### 会话服务

#### 存储方式
- **In-Memory**: 内存存储，适合单机部署
- **Redis**: 分布式存储，支持集群部署

#### 核心功能
- 会话创建与获取
- 事件存储与查询
- 状态管理
- TTL 支持

### 记忆系统

#### 记忆类型
- **长期记忆**: 用户偏好、历史信息
- **会话记忆**: 当前对话上下文
- **临时记忆**: 单次调用临时数据

#### 记忆操作
- 记忆存储与检索
- 记忆搜索与过滤
- 记忆更新与删除

### 状态注入

#### 占位符变量
- `{key}`: 会话状态值
- `{user:subkey}`: 用户态状态
- `{app:subkey}`: 应用态状态
- `{temp:subkey}`: 临时态状态

```go
llm := llmagent.New("agent",
    llmagent.WithInstruction("用户兴趣: {user:topics}"),
)
```

## 可观测性

### OpenTelemetry 集成
- 全链路追踪
- 性能监控
- 自定义指标

```go
// 启动 Langfuse 集成
clean, _ := langfuse.Start(ctx)
defer clean(ctx)
```

### 事件系统
- 实时事件流
- 结构化事件数据
- 完整的执行上下文

### 回调机制
- Agent 级回调
- Model 级回调
- Tool 级回调

```go
callbacks := agent.NewCallbacks()
callbacks.RegisterBeforeAgent(func(ctx context.Context, args *agent.BeforeAgentArgs) (*agent.BeforeAgentResult, error) {
    // 自定义逻辑
    return nil, nil
})
```

## 高级配置

### Provider 统一接口
- 简化多平台切换
- 统一配置选项
- 自动平台适配

```go
model, err := provider.Model(
    "openai",
    "deepseek-chat",
    provider.WithAPIKey(apiKey),
    provider.WithBaseURL(baseURL),
)
```

### 自定义 HTTP Header
- 网关支持
- 专有平台适配
- 代理环境支持

```go
model := openai.New("deepseek-chat",
    openai.WithHeaders(map[string]string{
        "X-Custom-Header": "custom-value",
        "X-Request-ID":    "req-123",
    }),
)
```

### Variant 优化
- 平台特有行为适配
- 自动 API 差异处理
- 环境变量自动配置

```go
// 混元平台适配
model := openai.New("hunyuan-model",
    openai.WithVariant(openai.VariantHunyuan),
)
```

## 性能优化

### 并发处理
- Go 协程并发
- 工具并行执行
- 异步事件处理

### 内存管理
- 自动资源清理
- 流式数据处理
- 缓冲区优化

### 网络优化
- 连接池复用
- 请求合并
- 智能重试

## 安全特性

### 权限控制
- 工具级权限
- 用户级隔离
- 会话级安全

### 输入验证
- 参数验证
- 类型检查
- 长度限制

### 错误处理
- 分层错误处理
- 敏感信息过滤
- 安全错误响应

## 相关文档

- [框架概述](./framework-overview.md)
- [核心架构](./core-architecture.md)
- [开发模式](./development-patterns.md)
- [高级功能](./advanced-features.md)
