package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/fromsko/krio/internal/config"
	"github.com/fromsko/krio/internal/tool"
	"github.com/fromsko/krio/pkg/logger"
)

// cacheCmd 缓存管理命令
var cacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "缓存管理",
	Long:  `管理网页抓取缓存,包括清空缓存和查看缓存统计。`,
}

// cacheClearCmd 清空缓存命令
var cacheClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "清空缓存",
	Long:  `清空所有网页抓取缓存。`,
	Run: func(cmd *cobra.Command, args []string) {
		// 加载配置
		cfg, err := loadConfigForCache()
		if err != nil {
			fmt.Printf("❌ 加载配置失败: %v\n", err)
			os.Exit(1)
		}

		// 初始化日志
		if err := logger.Init(cfg); err != nil {
			fmt.Printf("❌ 初始化日志失败: %v\n", err)
			os.Exit(1)
		}
		defer logger.Sync()

		// 创建工具
		ctx := context.Background()
		webNoteTool, err := tool.NewSaveWebNoteTool(ctx, cfg)
		if err != nil {
			fmt.Printf("❌ 创建工具失败: %v\n", err)
			os.Exit(1)
		}

		// 清空缓存
		webNoteTool.ClearCache()
		fmt.Println("✅ 缓存已清空")
	},
}

// cacheStatsCmd 缓存统计命令
var cacheStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "缓存统计",
	Long:  `显示缓存统计信息。`,
	Run: func(cmd *cobra.Command, args []string) {
		// 加载配置
		cfg, err := loadConfigForCache()
		if err != nil {
			fmt.Printf("❌ 加载配置失败: %v\n", err)
			os.Exit(1)
		}

		// 初始化日志
		if err := logger.Init(cfg); err != nil {
			fmt.Printf("❌ 初始化日志失败: %v\n", err)
			os.Exit(1)
		}
		defer logger.Sync()

		// 创建工具
		ctx := context.Background()
		webNoteTool, err := tool.NewSaveWebNoteTool(ctx, cfg)
		if err != nil {
			fmt.Printf("❌ 创建工具失败: %v\n", err)
			os.Exit(1)
		}

		// 获取缓存统计
		stats := webNoteTool.GetCacheStats()

		// 显示统计信息
		fmt.Println("\n📊 缓存统计")
		fmt.Println("=" + "===========")
		printStats(stats)
	},

	// 支持输出格式
}

var outputFormat string

func printStats(stats map[string]interface{}) {
	switch outputFormat {
	case "json":
		jsonData, err := json.MarshalIndent(stats, "", "  ")
		if err != nil {
			fmt.Printf("❌ JSON 序列化失败: %v\n", err)
			return
		}
		fmt.Println(string(jsonData))
	default:
		enabled, _ := stats["enabled"].(bool)
		if !enabled {
			fmt.Println("状态: 未启用")
			return
		}

		cacheSize, _ := stats["cache_size"].(int)
		cacheTTL, _ := stats["cache_ttl"].(string)
		maxConcurrency, _ := stats["max_concurrency"].(int)

		fmt.Printf("状态: 已启用\n")
		fmt.Printf("缓存条目: %d\n", cacheSize)
		fmt.Printf("缓存 TTL: %s\n", cacheTTL)
		fmt.Printf("最大并发: %d\n", maxConcurrency)
	}
	fmt.Println("=" + "===========")
}

func loadConfigForCache() (*config.Config, error) {
	if cfgFile != "" {
		return config.Load(cfgFile)
	}
	return config.LoadDefault()
}

func init() {
	rootCmd.AddCommand(cacheCmd)
	cacheCmd.AddCommand(cacheClearCmd)
	cacheCmd.AddCommand(cacheStatsCmd)

	cacheStatsCmd.Flags().StringVarP(&outputFormat, "output", "o", "table",
		"输出格式 (table/json)")
}
