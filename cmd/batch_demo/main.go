package main

import (
	"context"
	"fmt"
	"os"

	"github.com/fromsko/krio/internal/config"
	"github.com/fromsko/krio/internal/tool"
	"github.com/fromsko/krio/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	// 1. 加载配置
	cfg, err := config.LoadDefault()
	if err != nil {
		fmt.Fprintf(os.Stderr, "加载配置失败: %v\n", err)
		os.Exit(1)
	}

	// 2. 验证配置
	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "配置验证失败: %v\n", err)
		os.Exit(1)
	}

	// 3. 初始化日志
	if err := logger.Init(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "初始化日志失败: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	log := logger.Get()
	log.Info("启动批量处理演示",
		zap.String("name", cfg.App.Name),
		zap.String("version", cfg.App.Version),
	)

	// 4. 创建工具
	ctx := context.Background()
	webNoteTool, err := tool.NewSaveWebNoteTool(ctx, cfg)
	if err != nil {
		log.Fatal("创建工具失败", zap.Error(err))
	}

	// 5. 显示缓存状态
	stats := webNoteTool.GetCacheStats()
	log.Info("缓存状态", zap.Any("stats", stats))

	// 6. 批量处理演示
	urls := []string{
		"https://example.com",
		"https://www.iana.org/domains/reserved",
	}

	tags := []string{"demo", "batch-test"}
	folder := "Inbox"

	log.Info("开始批量处理", zap.Int("url_count", len(urls)))
	responses := webNoteTool.SaveWebNoteBatch(ctx, urls, tags, folder)

	// 7. 显示结果
	fmt.Println("\n" + "=" + string(make([]byte, 80)) + "=")
	fmt.Println("批量处理结果:")
	fmt.Println("=" + string(make([]byte, 80)) + "=")

	successCount := 0
	for i, resp := range responses {
		fmt.Printf("\n[%d] URL: %s\n", i+1, urls[i])
		fmt.Printf("    状态: %s\n", map[bool]string{true: "✅ 成功", false: "❌ 失败"}[resp.Success])
		fmt.Printf("    消息: %s\n", resp.Message)
		if resp.Title != "" {
			fmt.Printf("    标题: %s\n", resp.Title)
		}
		if resp.FilePath != "" {
			fmt.Printf("    路径: %s\n", resp.FilePath)
		}
		if resp.Success {
			successCount++
		}
	}

	fmt.Println("\n" + "=" + string(make([]byte, 80)) + "=")
	fmt.Printf("总计: %d 成功 / %d 失败\n", successCount, len(responses)-successCount)
	fmt.Println("=" + string(make([]byte, 80)) + "=")

	// 8. 再次显示缓存状态
	stats = webNoteTool.GetCacheStats()
	log.Info("处理后的缓存状态", zap.Any("stats", stats))
}
