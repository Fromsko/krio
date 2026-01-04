package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/fromsko/krio/internal/config"
	"github.com/fromsko/krio/internal/parser"
	"github.com/fromsko/krio/internal/tool"
	"github.com/fromsko/krio/pkg/logger"
	"go.uber.org/zap"
)

var (
	urlFile   string
	singleURL string
	tags      []string
	folder    string
)

// runCmd 运行命令
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "运行网页笔记生成器",
	Long:  `从 URL 或文件批量生成网页笔记并保存到 Obsidian。`,
	Run: func(cmd *cobra.Command, args []string) {
		// 加载配置
		cfg, err := loadConfig()
		if err != nil {
			fmt.Printf("❌ 加载配置失败: %v\n", err)
			os.Exit(1)
		}

		// 验证配置
		if err := cfg.Validate(); err != nil {
			fmt.Printf("❌ 配置验证失败: %v\n", err)
			os.Exit(1)
		}

		// 初始化日志
		if err := logger.Init(cfg); err != nil {
			fmt.Printf("❌ 初始化日志失败: %v\n", err)
			os.Exit(1)
		}
		defer logger.Sync()

		log := logger.Get()
		log.Info("启动 Krio",
			zap.String("version", cfg.App.Version),
			zap.Bool("debug", cfg.App.Debug),
		)

		// 创建工具
		ctx := context.Background()
		webNoteTool, err := tool.NewSaveWebNoteTool(ctx, cfg)
		if err != nil {
			log.Fatal("创建工具失败", zap.Error(err))
		}

		// 根据参数执行
		switch {
		case singleURL != "":
			runSingleURL(ctx, webNoteTool, singleURL, tags, folder)
		case urlFile != "":
			runFile(ctx, webNoteTool, urlFile, tags, folder)
		default:
			fmt.Println("❌ 请指定 -u <url> 或 -r <file>")
			cmd.Help()
			os.Exit(1)
		}
	},
}

func loadConfig() (*config.Config, error) {
	if cfgFile != "" {
		return config.Load(cfgFile)
	}
	return config.LoadDefault()
}

func runSingleURL(ctx context.Context, webNoteTool *tool.SaveWebNoteTool, url string, tags []string, folder string) {
	log := logger.Get()
	log.Info("处理单个 URL", zap.String("url", url))

	req := tool.SaveWebNoteRequest{
		URL:    url,
		Tags:   tags,
		Folder: folder,
	}

	resp, err := webNoteTool.SaveWebNote(ctx, req)
	if err != nil {
		log.Error("处理失败", zap.String("url", url), zap.Error(err))
		fmt.Printf("❌ 处理失败: %v\n", err)
		return
	}

	if resp.Success {
		fmt.Println("\n" + strings.Repeat("=", 80))
		fmt.Println("✅ 笔记生成成功")
		fmt.Println(strings.Repeat("=", 80))
		fmt.Printf("标题: %s\n", resp.Title)
		fmt.Printf("路径: %s\n", resp.FilePath)
		if len(tags) > 0 {
			fmt.Printf("标签: %v\n", tags)
		}
		fmt.Println(strings.Repeat("=", 80))
	}
}

func runFile(ctx context.Context, webNoteTool *tool.SaveWebNoteTool, filePath string, tags []string, folder string) {
	log := logger.Get()
	log.Info("批量处理文件", zap.String("file", filePath))

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		log.Error("打开文件失败", zap.String("file", filePath), zap.Error(err))
		fmt.Printf("❌ 打开文件失败: %v\n", err)
		return
	}
	defer file.Close()

	// 根据文件扩展名选择解析器
	p := parser.DetectFormat(filePath)

	// 解析 URL
	urls, err := p.Parse(file)
	if err != nil {
		log.Error("解析文件失败", zap.String("file", filePath), zap.Error(err))
		fmt.Printf("❌ 解析文件失败: %v\n", err)
		return
	}

	if len(urls) == 0 {
		fmt.Println("❌ 未找到任何 URL")
		return
	}

	log.Info("找到 URL", zap.Int("count", len(urls)))
	fmt.Printf("\n📝 开始处理 %d 个 URL...\n\n", len(urls))

	// 批量处理
	responses := webNoteTool.SaveWebNoteBatch(ctx, urls, tags, folder)

	// 显示结果
	printResults(urls, responses)
}

func printResults(urls []string, responses []tool.SaveWebNoteResponse) {
	successCount := 0
	failCount := 0

	fmt.Println(strings.Repeat("=", 100))
	fmt.Printf("%-5s %-50s %-20s %s\n", "#", "URL", "标题", "状态")
	fmt.Println(strings.Repeat("=", 100))

	for i, resp := range responses {
		status := "✅ 成功"
		if !resp.Success {
			status = "❌ 失败"
			failCount++
		} else {
			successCount++
		}

		// 截断 URL 显示
		urlDisplay := urls[i]
		if len(urlDisplay) > 47 {
			urlDisplay = urlDisplay[:47] + "..."
		}

		// 截断标题显示
		titleDisplay := resp.Title
		if len(titleDisplay) > 18 {
			titleDisplay = titleDisplay[:18] + "..."
		}

		fmt.Printf("%-5d %-50s %-20s %s\n", i+1, urlDisplay, titleDisplay, status)
	}

	fmt.Println(strings.Repeat("=", 100))
	fmt.Printf("总计: %d 成功, %d 失败\n\n", successCount, failCount)
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVarP(&singleURL, "url", "u", "",
		"单个 URL")
	runCmd.Flags().StringVarP(&urlFile, "require", "r", "",
		"需求文件 (.txt/.md)")
	runCmd.Flags().StringSliceVarP(&tags, "tags", "t", []string{},
		"自定义标签 (逗号分隔)")
	runCmd.Flags().StringVarP(&folder, "folder", "f", "",
		"目标文件夹")
}
