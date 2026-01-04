package note

import (
	"fmt"
	"strings"
	"time"

	"github.com/fromsko/krio/internal/config"
	"github.com/fromsko/krio/internal/summarizer"
	"github.com/rs/xid"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Generator 笔记生成器
type Generator struct {
	cfg *config.NoteConfig
}

// NewGenerator 创建笔记生成器
func NewGenerator(cfg *config.NoteConfig) *Generator {
	return &Generator{cfg: cfg}
}

// Generate 生成 Markdown 笔记
func (g *Generator) Generate(summary *summarizer.Summary, sourceURL string) string {
	now := time.Now()
	timestamp := now.Format("2006-01-02T15:04:05")

	// 生成 frontmatter
	frontmatter := g.generateFrontmatter(summary, sourceURL, timestamp)

	// 生成内容
	content := g.generateContent(summary)

	// 组合
	note := fmt.Sprintf("%s\n\n%s", frontmatter, content)

	return note
}

// generateFrontmatter 生成 frontmatter
func (g *Generator) generateFrontmatter(summary *summarizer.Summary, sourceURL, timestamp string) string {
	tags := g.formatTags(summary.Tags)
	filename := g.generateFilename(summary.Title, summary.Tags)

	frontmatter := fmt.Sprintf(`---
title: %s
source: %s
date: %s
tags: [%s]
filename: %s
created_at: %s
updated_at: %s
id: %s
---
`,
		escapeYAML(summary.Title),
		escapeYAML(sourceURL),
		timestamp,
		tags,
		filename,
		timestamp,
		timestamp,
		xid.New().String(),
	)

	return frontmatter
}

// generateContent 生成内容
func (g *Generator) generateContent(summary *summarizer.Summary) string {
	var sb strings.Builder

	// 标题
	sb.WriteString(fmt.Sprintf("# %s\n\n", summary.Title))

	// 一句话总结
	sb.WriteString(fmt.Sprintf("> %s\n\n", summary.OneSentence))

	// 关键要点
	if len(summary.KeyPoints) > 0 {
		sb.WriteString("## 核心要点\n\n")
		for _, point := range summary.KeyPoints {
			sb.WriteString(fmt.Sprintf("- %s\n", point))
		}
		sb.WriteString("\n")
	}

	// 详细内容 (可选)
	// sb.WriteString("## 详细内容\n\n")
	// sb.WriteString(summary.OriginalContent)
	// sb.WriteString("\n")

	return sb.String()
}

// formatTags 格式化标签
func (g *Generator) formatTags(tags []string) string {
	if len(tags) == 0 {
		return ""
	}

	formatted := make([]string, len(tags))
	for i, tag := range tags {
		formatted[i] = fmt.Sprintf(`"%s"`, tag)
	}
	return strings.Join(formatted, ", ")
}

// generateFilename 生成文件名
func (g *Generator) generateFilename(title string, tags []string) string {
	// 清理标题
	filename := sanitizeFilename(title)

	// 限制长度
	if len(filename) > 50 {
		filename = filename[:50]
	}

	// 添加时间戳
	if g.cfg.AddTimestamp {
		timestamp := time.Now().Format("2006-01-02-150405")
		filename = fmt.Sprintf("%s-%s", filename, timestamp)
	}

	return filename
}

// sanitizeFilename 清理文件名
func sanitizeFilename(name string) string {
	// 转换为小写
	name = strings.ToLower(name)

	// 替换空格为短横线
	name = strings.ReplaceAll(name, " ", "-")

	// 移除非法字符
	invalidChars := []string{
		"<", ">", ":", "\"", "|", "?", "*",
		"\\", "/", "\n", "\r", "\t",
	}
	for _, char := range invalidChars {
		name = strings.ReplaceAll(name, char, "")
	}

	// 移除连续的短横线
	for strings.Contains(name, "--") {
		name = strings.ReplaceAll(name, "--", "-")
	}

	// 移除首尾的短横线和点
	name = strings.Trim(name, "-.")
	name = strings.Trim(name, " ")

	// 如果为空,使用默认名称
	if name == "" {
		name = "untitled"
	}

	return name
}

// escapeYAML 转义 YAML 特殊字符
func escapeYAML(s string) string {
	s = strings.ReplaceAll(s, "'", "''")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "\\r")
	s = strings.ReplaceAll(s, "\t", "\\t")
	return s
}

// TitleCase 标题格式化
func TitleCase(s string) string {
	return cases.Title(language.English).String(s)
}

// GenerateFilename 导出的文件名生成函数
func GenerateFilename(title string) string {
	gen := &Generator{
		cfg: &config.NoteConfig{
			AddTimestamp: true,
		},
	}
	return gen.generateFilename(title, []string{})
}
