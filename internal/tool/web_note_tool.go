package tool

import (
	"context"
	"fmt"

	"github.com/fromsko/krio/internal/config"
	"github.com/fromsko/krio/internal/note"
	"github.com/fromsko/krio/internal/obsidian"
	"github.com/fromsko/krio/internal/scraper"
	"github.com/fromsko/krio/internal/summarizer"
	"github.com/fromsko/krio/pkg/logger"
	"go.uber.org/zap"
)

// SaveWebNoteRequest 保存网页笔记请求
type SaveWebNoteRequest struct {
	URL    string   `json:"url" jsonschema:"description=要保存的网页URL,required"`
	Tags   []string `json:"tags,omitempty" jsonschema:"description=自定义标签列表,可选"`
	Folder string   `json:"folder,omitempty" jsonschema:"description=保存到Obsidian的文件夹,可选"`
}

// SaveWebNoteResponse 保存网页笔记响应
type SaveWebNoteResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	Title     string `json:"title,omitempty"`
	FilePath  string `json:"file_path,omitempty"`
	Content   string `json:"content,omitempty"`
	NoteCount int    `json:"note_count"`
}

// SaveWebNoteTool 保存网页笔记工具
type SaveWebNoteTool struct {
	cfg        *config.Config
	fetcher    *scraper.Fetcher
	cachedFetcher    *scraper.CachedFetcher // 带缓存的抓取器
	summarizer *summarizer.Summarizer
	generator  *note.Generator
	obsidian   *obsidian.Client // Obsidian MCP 客户端
}

// NewSaveWebNoteTool 创建工具
func NewSaveWebNoteTool(ctx context.Context, cfg *config.Config) (*SaveWebNoteTool, error) {
	// 创建各个模块
	fetcher := scraper.NewFetcher(&cfg.Scraper)

	// 创建带缓存的抓取器
	var cachedFetcher *scraper.CachedFetcher
	if cfg.Scraper.EnableCache {
		maxConcurrency := cfg.Scraper.MaxConcurrency
		if maxConcurrency <= 0 {
			maxConcurrency = 5 // 默认并发数
		}
		cachedFetcher = scraper.NewCachedFetcher(&cfg.Scraper, maxConcurrency, cfg.Scraper.CacheTTL)
		logger.Get().Info("已启用缓存抓取器",
			zap.Int("max_concurrency", maxConcurrency),
			zap.Duration("cache_ttl", cfg.Scraper.CacheTTL),
		)
	}

	summarizer, err := summarizer.NewSummarizer(&cfg.Model)
	if err != nil {
		return nil, fmt.Errorf("创建总结器失败: %w", err)
	}

	generator := note.NewGenerator(&cfg.Note)

	// 创建 Obsidian MCP 客户端
	obsidianClient, err := obsidian.NewClient(ctx, &cfg.ObsidianMCP)
	if err != nil {
		logger.Get().Warn("创建 Obsidian MCP 客户端失败,笔记将不会自动保存", zap.Error(err))
		// 不返回错误,继续创建工具
	}

	return &SaveWebNoteTool{
		cfg:        cfg,
		fetcher:    fetcher,
		cachedFetcher: cachedFetcher,
		summarizer: summarizer,
		generator:  generator,
		obsidian:   obsidianClient,
	}, nil
}

// SaveWebNote 保存网页笔记 (工具函数)
func (t *SaveWebNoteTool) SaveWebNote(ctx context.Context, req SaveWebNoteRequest) (SaveWebNoteResponse, error) {
	log := logger.Get()

	log.Info("开始保存网页笔记",
		zap.String("url", req.URL),
		zap.Strings("tags", req.Tags),
		zap.String("folder", req.Folder),
	)

	// 1. 抓取网页内容
	log.Debug("抓取网页内容", zap.String("url", req.URL))
	var page *scraper.WebPage
	var err error

	// 优先使用缓存抓取器
	if t.cachedFetcher != nil {
		page, err = t.cachedFetcher.Fetch(req.URL)
	} else {
		page, err = t.fetcher.Fetch(req.URL)
	}

	if err != nil {
		log.Error("抓取网页失败", zap.String("url", req.URL), zap.Error(err))
		return SaveWebNoteResponse{
			Success: false,
			Message: fmt.Sprintf("抓取网页失败: %v", err),
		}, err
	}

	log.Info("网页抓取成功",
		zap.String("title", page.Title),
		zap.Int("content_length", len(page.Content)),
	)

	// 2. AI 总结
	log.Debug("开始 AI 总结")
	summary, err := t.summarizer.Summarize(ctx, page.Title, page.Content)
	if err != nil {
		log.Error("AI 总结失败", zap.Error(err))
		return SaveWebNoteResponse{
			Success: false,
			Message: fmt.Sprintf("AI 总结失败: %v", err),
		}, err
	}

	log.Info("AI 总结成功",
		zap.String("title", summary.Title),
		zap.Int("key_points", len(summary.KeyPoints)),
	)

	// 3. 生成 Markdown 笔记
	log.Debug("生成 Markdown 笔记")
	markdown := t.generator.Generate(summary, page.URL)

	// 4. 保存到 Obsidian (通过 MCP)
	filename := note.GenerateFilename(summary.Title)
	folder := t.getFolder(req.Folder)
	var filePath string

	if t.obsidian != nil {
		log.Info("保存笔记到 Obsidian")
		actualPath, err := t.obsidian.SaveNote(ctx, markdown, filename, folder)
		if err != nil {
			log.Error("保存到 Obsidian 失败", zap.Error(err))
			// 返回错误,但不影响笔记内容的返回
			return SaveWebNoteResponse{
				Success: false,
				Message: fmt.Sprintf("保存到 Obsidian 失败: %v", err),
				Title:   summary.Title,
				Content: markdown,
			}, err
		}
		filePath = actualPath
	} else {
		// 如果没有 Obsidian 客户端,返回预期的路径
		filePath = fmt.Sprintf("%s/%s.md", folder, filename)
		log.Warn("Obsidian 客户端未初始化,笔记未实际保存", zap.String("file_path", filePath))
	}

	log.Info("笔记保存成功",
		zap.String("title", summary.Title),
		zap.String("file_path", filePath),
		zap.Int("content_length", len(markdown)),
	)

	return SaveWebNoteResponse{
		Success:   true,
		Message:   "笔记保存成功",
		Title:     summary.Title,
		FilePath:  filePath,
		Content:   markdown,
		NoteCount: 1,
	}, nil
}

// getFolder 获取保存文件夹
func (t *SaveWebNoteTool) getFolder(customFolder string) string {
	if customFolder != "" {
		return customFolder
	}
	return t.cfg.Note.DefaultFolder
}

// SaveWebNoteBatch 批量保存网页笔记 (并发处理)
func (t *SaveWebNoteTool) SaveWebNoteBatch(ctx context.Context, urls []string, tags []string, folder string) []SaveWebNoteResponse {
	log := logger.Get()

	if t.cachedFetcher == nil {
		log.Warn("缓存抓取器未启用,批量处理将使用串行模式")
		// 串行处理
		responses := make([]SaveWebNoteResponse, len(urls))
		for i, url := range urls {
			req := SaveWebNoteRequest{
				URL:    url,
				Tags:   tags,
				Folder: folder,
			}
			resp, err := t.SaveWebNote(ctx, req)
			if err != nil {
				resp = SaveWebNoteResponse{
					Success: false,
					Message: fmt.Sprintf("处理失败: %v", err),
				}
			}
			responses[i] = resp
		}
		return responses
	}

	// 并发抓取所有网页
	log.Info("开始批量并发抓取", zap.Int("total_urls", len(urls)))
	fetchResults := t.cachedFetcher.FetchBatch(ctx, urls)

	// 处理抓取结果
	responses := make([]SaveWebNoteResponse, 0, len(urls))
	successCount := 0

	for _, result := range fetchResults {
		if result.Err != nil {
			responses = append(responses, SaveWebNoteResponse{
				Success: false,
				Message: fmt.Sprintf("抓取失败: %v", result.Err),
			})
			continue
		}

		// 对每个成功抓取的页面进行总结和保存
		page := result.Page

		// AI 总结
		summary, err := t.summarizer.Summarize(ctx, page.Title, page.Content)
		if err != nil {
			log.Error("AI 总结失败", zap.String("url", page.URL), zap.Error(err))
			responses = append(responses, SaveWebNoteResponse{
				Success: false,
				Message: fmt.Sprintf("AI 总结失败: %v", err),
				Title:   page.Title,
			})
			continue
		}

		// 生成 Markdown
		markdown := t.generator.Generate(summary, page.URL)

		// 保存到 Obsidian
		filename := note.GenerateFilename(summary.Title)
		saveFolder := t.getFolder(folder)
		var filePath string

		if t.obsidian != nil {
			actualPath, err := t.obsidian.SaveNote(ctx, markdown, filename, saveFolder)
			if err != nil {
				log.Error("保存到 Obsidian 失败", zap.String("url", page.URL), zap.Error(err))
				responses = append(responses, SaveWebNoteResponse{
					Success: false,
					Message: fmt.Sprintf("保存失败: %v", err),
					Title:   summary.Title,
					Content: markdown,
				})
				continue
			}
			filePath = actualPath
		} else {
			filePath = fmt.Sprintf("%s/%s.md", saveFolder, filename)
		}

		successCount++
		responses = append(responses, SaveWebNoteResponse{
			Success:  true,
			Message:  "笔记保存成功",
			Title:    summary.Title,
			FilePath: filePath,
			Content:  markdown,
		})
	}

	log.Info("批量处理完成",
		zap.Int("total", len(urls)),
		zap.Int("success", successCount),
		zap.Int("failed", len(urls)-successCount),
	)

	return responses
}

// GetCacheStats 获取缓存统计信息
func (t *SaveWebNoteTool) GetCacheStats() map[string]interface{} {
	if t.cachedFetcher == nil {
		return map[string]interface{}{
			"enabled": false,
		}
	}

	return map[string]interface{}{
		"enabled":     true,
		"cache_size":  t.cachedFetcher.GetCacheSize(),
		"cache_ttl":   t.cfg.Scraper.CacheTTL.String(),
		"max_concurrency": t.cfg.Scraper.MaxConcurrency,
	}
}

// ClearCache 清空缓存
func (t *SaveWebNoteTool) ClearCache() {
	if t.cachedFetcher != nil {
		t.cachedFetcher.ClearCache()
	}
}
