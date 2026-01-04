# External References

此目录用于存放项目的外部参考资料，这些资料将作为 AI 智能体的知识库。

## 目录结构

```
references/
├── README.md              # 本文件
├── technical/            # 技术文档
│   ├── index.md          # tRPC-Agent-Go 框架介绍
│   ├── Trpc-agent-go.md  # 官方 README
│   ├── agent.md          # Agent 使用文档
│   ├── model.md          # Model 模块文档
│   ├── tool.md           # Tool 工具文档
│   ├── runner.md         # Runner 执行器文档
│   ├── session.md        # Session 会话管理
│   ├── memory.md         # Memory 记忆管理
│   ├── knowledge.md      # Knowledge 知识库
│   ├── planner.md        # Planner 规划器
│   ├── graph.md          # Graph 工作流
│   ├── multiagent.md     # Multi-Agent 多智能体
│   ├── callbacks.md      # Callbacks 回调机制
│   ├── event.md          # Event 事件系统
│   ├── observability.md  # Observability 可观测性
│   ├── custom-agent.md   # 自定义 Agent
│   ├── skill.md          # Agent Skills
│   ├── artifact.md       # Artifacts 工件管理
│   ├── evaluation.md     # Evaluation 评测
│   ├── plugin.md         # Plugin 插件系统
│   ├── agui.md           # AG-UI 用户交互
│   ├── a2a.md            # A2A Agent 互通
│   ├── dify.md           # Dify 集成
│   ├── ecosystem.md      # 生态系统
│   └── wikipedia.md      # Wikipedia 集成
├── business/             # 业务规则
│   ├── workflows/        # 工作流程
│   └── policies/         # 策略文档
├── domain/               # 领域知识
│   ├── concepts/         # 概念定义
│   └── examples/         # 示例
└── external/             # 外部资源
    ├── links/            # 外部链接
    └── papers/           # 论文/文章
```

## 技术文档索引

### 核心组件
- **[index.md](technical/index.md)** - tRPC-Agent-Go 框架介绍，快速了解框架特性
- **[Trpc-agent-go.md](technical/Trpc-agent-go.md)** - 官方 README，包含完整的使用指南
- **[agent.md](technical/agent.md)** - Agent 使用文档，核心执行单元
- **[model.md](technical/model.md)** - Model 模块，LLM 抽象层
- **[tool.md](technical/tool.md)** - Tool 工具系统，外部服务集成

### 执行与管理
- **[runner.md](technical/runner.md)** - Runner 执行器，管理 Agent 执行流程
- **[session.md](technical/session.md)** - Session 会话管理
- **[memory.md](technical/memory.md)** - Memory 记忆管理
- **[knowledge.md](technical/knowledge.md)** - Knowledge 知识库，RAG 能力

### 高级特性
- **[planner.md](technical/planner.md)** - Planner 规划器，任务分解
- **[graph.md](technical/graph.md)** - Graph 工作流，图式编排
- **[multiagent.md](technical/multiagent.md)** - Multi-Agent 多智能体协作
- **[callbacks.md](technical/callbacks.md)** - Callbacks 回调机制
- **[event.md](technical/event.md)** - Event 事件系统

### 扩展与集成
- **[custom-agent.md](technical/custom-agent.md)** - 自定义 Agent
- **[skill.md](technical/skill.md)** - Agent Skills，可复用工作流
- **[artifact.md](technical/artifact.md)** - Artifacts 工件管理
- **[evaluation.md](technical/evaluation.md)** - Evaluation 评测框架
- **[plugin.md](technical/plugin.md)** - Plugin 插件系统
- **[agui.md](technical/agui.md)** - AG-UI 用户交互
- **[a2a.md](technical/a2a.md)** - A2A Agent 互通
- **[dify.md](technical/dify.md)** - Dify 集成
- **[ecosystem.md](technical/ecosystem.md)** - 生态系统
- **[wikipedia.md](technical/wikipedia.md)** - Wikipedia 集成

### 运维与监控
- **[observability.md](technical/observability.md)** - Observability 可观测性

## 使用指南

### 添加参考资料

1. 根据资料类型选择合适的子目录
2. 使用清晰的文件命名 (如 `user-auth-flow.md`)
3. 在文件顶部添加元数据:

```markdown
---
title: 资料标题
type: technical/business/domain/external
tags: [tag1, tag2]
updated: 2024-01-01
---

内容...
```

### 引用参考资料

在代码或文档中引用参考资料时，使用相对路径:

```
详见 references/technical/agent.md
```

### 索引管理

定期更新本目录的索引文件，方便查找和检索。

## 注意事项

- 不要在此目录存放敏感信息 (密码、密钥等)
- 保持文档格式统一 (推荐使用 Markdown)
- 定期清理过时的参考资料
- 重要资料建议进行版本控制

## 快速开始

如果你是第一次使用 tRPC-Agent-Go 框架，建议按以下顺序阅读：

1. 先读 **[index.md](technical/index.md)** 了解框架概况
2. 再读 **[Trpc-agent-go.md](technical/Trpc-agent-go.md)** 查看官方示例
3. 然后根据需要深入各个组件文档：
   - 想 Agent 怎么工作 → 读 **[agent.md](technical/agent.md)**
   - 想集成 LLM → 读 **[model.md](technical/model.md)**
   - 想添加工具 → 读 **[tool.md](technical/tool.md)**
   - 想管理会话 → 读 **[runner.md](technical/runner.md)**
