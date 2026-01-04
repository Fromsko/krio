package obsidian

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fromsko/krio/internal/config"
	"github.com/fromsko/krio/pkg/logger"
	"trpc.group/trpc-go/trpc-agent-go/tool"
	"trpc.group/trpc-go/trpc-agent-go/tool/mcp"
	"go.uber.org/zap"
)

// Client Obsidian MCP 客户端
type Client struct {
	toolSet *mcp.ToolSet
	cfg     *config.ObsidianMCPConfig
}

// NewClient 创建 Obsidian MCP 客户端
func NewClient(ctx context.Context, cfg *config.ObsidianMCPConfig) (*Client, error) {
	log := logger.Get()

	// 创建 MCP ToolSet
	toolSet := mcp.NewMCPToolSet(
		mcp.ConnectionConfig{
			Transport: cfg.Transport,
			Command:   cfg.Command,
			Args:      cfg.Args,
			Timeout:   cfg.Timeout,
		},
	)

	// 显式初始化
	log.Info("初始化 Obsidian MCP 服务器",
		zap.String("command", cfg.Command),
		zap.Strings("args", cfg.Args),
	)

	if err := toolSet.Init(ctx); err != nil {
		return nil, fmt.Errorf("初始化 Obsidian MCP 工具集失败: %w", err)
	}

	log.Info("Obsidian MCP 服务器初始化成功")

	return &Client{
		toolSet: toolSet,
		cfg:     cfg,
	}, nil
}

// SaveNote 保存笔记到 Obsidian
func (c *Client) SaveNote(ctx context.Context, content, filename, folder string) (string, error) {
	log := logger.Get()

	// 获取工具列表
	tools := c.toolSet.Tools(ctx)

	// 查找 create_note 工具
	var createNoteTool tool.Tool
	for _, t := range tools {
		decl := t.Declaration()
		if decl != nil && decl.Name == "create_note" {
			createNoteTool = t
			log.Info("找到 create_note 工具")
			break
		}
	}

	if createNoteTool == nil {
		return "", fmt.Errorf("未找到 create_note 工具")
	}

	// 准备参数 - 尝试不同的参数组合
	// 根据错误信息 "path" 参数缺失,尝试使用 "path" 作为参数名
	fullPath := folder + "/" + filename + ".md"
	args := map[string]interface{}{
		"path":    fullPath,
		"content": content,
	}

	log.Info("保存笔记到 Obsidian",
		zap.String("path", fullPath),
		zap.Int("content_length", len(content)),
	)

	// 序列化为 JSON
	jsonArgs, err := json.Marshal(args)
	if err != nil {
		return "", fmt.Errorf("序列化参数失败: %w", err)
	}

	// 调用工具
	if callable, ok := createNoteTool.(tool.CallableTool); ok {
		result, err := callable.Call(ctx, jsonArgs)
		if err != nil {
			return "", fmt.Errorf("调用 create_note 工具失败: %w", err)
		}

		log.Debug("MCP 工具返回原始结果", zap.Any("result", result))

		// 解析返回结果 - 期望格式: [{"type":"text","text":"笔记创建成功: path/to/file.md"}]
		resultArray, ok := result.([]interface{})
		if !ok || len(resultArray) == 0 {
			return fullPath, nil // 返回我们传入的路径作为fallback
		}

		firstItem := resultArray[0]
		itemMap, ok := firstItem.(map[string]interface{})
		if !ok {
			return fullPath, nil
		}

		if text, ok := itemMap["text"].(string); ok {
			// 提取文件路径
			// 格式: "笔记创建成功: Inbox/example-domain-2026-01-05-015713.md"
			if idx := strings.Index(text, ": "); idx > 0 {
				filePath := text[idx+2:]
				log.Info("笔记保存成功", zap.String("file_path", filePath))
				return filePath, nil
			}
		}

		// Fallback: 返回我们传入的路径
		log.Info("笔记保存成功", zap.String("file_path", fullPath))
		return fullPath, nil
	}

	return "", fmt.Errorf("工具不支持调用")
}

// Close 关闭客户端
func (c *Client) Close() error {
	log := logger.Get()
	log.Info("关闭 Obsidian MCP 客户端")
	return c.toolSet.Close()
}
