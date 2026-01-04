# API 配置参考

## 概述

本文档记录了项目中使用的 API 密钥和配置信息。

**⚠️ 安全提示:**
- `config.yaml` 包含真实的 API 密钥,已添加到 `.gitignore`
- 分享代码时只提交 `config.example.yaml`
- 生产环境建议使用环境变量或密钥管理服务

## 智云 API (GLM-4)

### API 配置

| 配置项 | 说明 |
|--------|-----|
| API Key | 从智云平台获取 (https://open.bigmodel.cn/) |
| Base URL | `https://open.bigmodel.cn/api/coding/paas/v4` |
| 推荐模型 | `glm-4.7` |

### 使用方式

#### 方式 1: 配置文件 (开发环境)

在 `config.yaml` 中配置:

```yaml
model:
  api_key: "YOUR_API_KEY_HERE"  # 从智云平台获取
  base_url: "https://open.bigmodel.cn/api/coding/paas/v4"
  model_name: "glm-4.7"
```

#### 方式 2: 环境变量 (生产环境)

```bash
export MODEL_API_KEY="YOUR_API_KEY_HERE"
export MODEL_BASE_URL="https://open.bigmodel.cn/api/coding/paas/v4"
export MODEL_NAME="glm-4.7"
```

代码中引用: [internal/config/config.go](../../internal/config/config.go)

## Obsidian MCP 服务器

### 配置

| 配置项 | 说明 |
|--------|-----|
| 命令 | `bun` |
| 参数 | `x --no-cache @andysama/obsidian-mcp-server --vault <your-vault-path>` |
| Vault 路径 | 你的 Obsidian vault 路径 |

### 使用方式

在 `config.yaml` 中配置:

```yaml
obsidian_mcp:
  transport: "stdio"
  command: "bun"
  args:
    - "x"
    - "--no-cache"
    - "@andysama/obsidian-mcp-server"
    - "--vault"
    - "/path/to/your/vault"  # 修改为你的 vault 路径
  timeout: 30
```

代码中引用: [internal/config/config.go](../../internal/config/config.go)

## 配置文件结构

```
krio/
├── config.yaml              # 实际配置 (包含密钥,不提交)
├── config.example.yaml      # 示例配置 (可提交)
├── internal/config/
│   └── config.go           # 配置管理代码
└── references/configuration/
    └── api-keys.md         # 本文档
```

## 代码中使用配置

```go
import "github.com/fromsko/krio/internal/config"

// 加载配置
cfg, err := config.LoadDefault()
if err != nil {
    log.Fatal(err)
}

// 验证配置
if err := cfg.Validate(); err != nil {
    log.Fatal(err)
}

// 使用配置
apiKey := cfg.Model.APIKey
model := cfg.Model.ModelName
```

## 密钥轮换

如果需要更换 API 密钥:

1. 更新 `config.yaml` 中的 `api_key`
2. 或者设置新的环境变量
3. 重启应用

## 相关文档

- [智云 API 官方文档](https://open.bigmodel.cn/)
- [Obsidian MCP 服务器文档](https://github.com/andysama/obsidian-mcp-server)
- [配置管理代码](../../internal/config/config.go)
