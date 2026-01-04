package tool

import (
	"context"
	"fmt"

	"github.com/fromsko/krio/internal/config"
	"github.com/fromsko/krio/internal/note"
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
	summarizer *summarizer.Summarizer
	generator  *note.Generator
}

// NewSaveWebNoteTool 创建工具
func NewSaveWebNoteTool(cfg *config.Config) (*SaveWebNoteTool, error) {
	// 创建各个模块
	fetcher := scraper.NewFetcher(&cfg.Scraper)

	summarizer, err := summarizer.NewSummarizer(&cfg.Model)
	if err != nil {
		return nil, fmt.Errorf("创建总结器失败: %w", err)
	}

	generator := note.NewGenerator(&cfg.Note)

	return &SaveWebNoteTool{
		cfg:        cfg,
		fetcher:    fetcher,
		summarizer: summarizer,
		generator:  generator,
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
	page, err := t.fetcher.Fetch(req.URL)
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
	// TODO: 实现调用 Obsidian MCP 服务器的逻辑
	// 现在先返回内容,手动保存

	log.Info("笔记生成成功",
		zap.String("title", summary.Title),
		zap.Int("content_length", len(markdown)),
	)

	return SaveWebNoteResponse{
		Success:   true,
		Message:   "笔记生成成功",
		Title:     summary.Title,
		FilePath:  fmt.Sprintf("%s/%s.md", t.getFolder(req.Folder), note.GenerateFilename(summary.Title)),
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
