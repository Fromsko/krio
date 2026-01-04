package cmd

import (
	"os"
	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd 根命令
var rootCmd = &cobra.Command{
	Use:   "krio",
	Short: "智能网页笔记 Agent",
	Long: `Krio 是一个基于 AI 的智能网页笔记工具，
自动抓取网页内容并生成结构化笔记保存到 Obsidian。

特性:
  - AI 驱动: 使用智云 GLM-4 模型进行内容理解和总结
  - 网页抓取: 自动提取网页核心内容,去除广告和无关元素
  - Markdown 笔记: 生成格式良好的 Markdown 笔记
  - 智能标签: AI 自动生成相关标签,便于分类和检索
  - 高性能: 支持并发处理和智能缓存
  - 易于配置: YAML 配置文件,支持环境变量覆盖`,
}

// Execute 执行根命令
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// 在这里可以添加持久化标志和初始化逻辑
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "",
		"配置文件路径 (默认: ./config.yaml)")
}
