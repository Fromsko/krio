package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

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
	log.Info("启动应用",
		zap.String("name", cfg.App.Name),
		zap.String("version", cfg.App.Version),
		zap.Bool("debug", cfg.App.Debug),
	)

	// 4. 创建工具
	ctx := context.Background()
	webNoteTool, err := tool.NewSaveWebNoteTool(ctx, cfg)
	if err != nil {
		log.Fatal("创建工具失败", zap.Error(err))
	}

	// 5. 演示使用
	demo(ctx, webNoteTool)

	log.Info("应用退出")
}

// demo 演示工具使用
func demo(ctx context.Context, webNoteTool *tool.SaveWebNoteTool) {
	log := logger.Get()

	// 示例 URL
	url := "https://example.com"
	log.Info("开始演示", zap.String("url", url))

	// 调用工具
	req := tool.SaveWebNoteRequest{
		URL:    url,
		Tags:   []string{"demo", "test"},
		Folder: "Inbox",
	}

	resp, err := webNoteTool.SaveWebNote(ctx, req)
	if err != nil {
		log.Error("保存笔记失败", zap.Error(err))
		return
	}

	log.Info("保存笔记成功",
		zap.Bool("success", resp.Success),
		zap.String("message", resp.Message),
		zap.String("title", resp.Title),
		zap.String("file_path", resp.FilePath),
	)

	// 只在成功时打印笔记内容
	if resp.Success {
		separator := "--------------------------------------------------------------------------------"
		fmt.Println("\n" + separator)
		fmt.Println("生成的笔记内容:")
		fmt.Println(separator)
		fmt.Println(resp.Content)
		fmt.Println(separator)
		fmt.Printf("✅ 笔记已保存到: %s\n", resp.FilePath)
		fmt.Println(separator)
	}
}

// setupSignalHandling 设置信号处理
func setupSignalHandling() chan struct{} {
	done := make(chan struct{})

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		logger.Get().Info("收到信号,开始优雅关闭", zap.String("signal", sig.String()))
		close(done)
	}()

	return done
}
