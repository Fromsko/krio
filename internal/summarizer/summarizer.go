package summarizer

import (
	"context"
	"fmt"

	"github.com/fromsko/krio/internal/config"
	"github.com/tidwall/gjson"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

// Summary 总结结果
type Summary struct {
	Title           string   `json:"title"`
	OneSentence     string   `json:"one_sentence"`
	KeyPoints       []string `json:"key_points"`
	Tags            []string `json:"tags"`
	OriginalContent string   `json:"original_content"`
}

// Summarizer 总结器
type Summarizer struct {
	llm  *openai.LLM
	cfg  *config.ModelConfig
}

// NewSummarizer 创建总结器
func NewSummarizer(cfg *config.ModelConfig) (*Summarizer, error) {
	llm, err := openai.New(
		openai.WithToken(cfg.APIKey),
		openai.WithBaseURL(cfg.BaseURL),
		openai.WithModel(cfg.ModelName),
	)
	if err != nil {
		return nil, fmt.Errorf("创建 LLM 失败: %w", err)
	}

	return &Summarizer{
		llm: llm,
		cfg: cfg,
	}, nil
}

// Summarize 总结内容
func (s *Summarizer) Summarize(ctx context.Context, title, content string) (*Summary, error) {
	prompt := s.buildPrompt(title, content)

	// 调用 LLM
	response, err := llms.GenerateFromSinglePrompt(ctx, s.llm, prompt)
	if err != nil {
		return nil, fmt.Errorf("LLM 调用失败: %w", err)
	}

	// 解析结果
	summary, err := s.parseSummary(response)
	if err != nil {
		return nil, fmt.Errorf("解析总结失败: %w", err)
	}

	summary.OriginalContent = content
	return summary, nil
}

// buildPrompt 构建提示词
func (s *Summarizer) buildPrompt(title, content string) string {
	return fmt.Sprintf(`你是一个专业的笔记助手。请深入分析以下网页内容，提取核心知识，重新组织成一份详细且实用的笔记。

网页标题: %s

网页内容:
%s

请按以下 JSON 格式返回笔记:
{
  "title": "文章标题(简洁明了)",
  "one_sentence": "一句话概括文章核心内容",
  "key_points": [
    "章节标题: 详细说明该章节的核心概念和要点",
    "重要概念: 解释概念的定义、原理和重要性",
    "关键步骤: 步骤1 -> 步骤2 -> 步骤3，每一步详细说明",
    "实用技巧: 提取实际应用中的技巧和注意事项",
    "示例说明: 对示例代码或案例进行详细解释说明",
    "常见问题: 列出常见问题和解决方案",
    "总结要点: 归纳总结关键知识点"
  ],
  "tags": ["标签1", "标签2", "标签3"]
}

要求:
1. title: 提取最合适的标题
2. one_sentence: 用一句话概括文章的核心价值,不超过50字
3. key_points: 生成7-15个详细要点，要求:
   - 不仅仅是简单概括，要提取具体知识点
   - 保留重要的技术细节、参数说明、代码示例等
   - 按照逻辑顺序组织（从概念到实践）
   - 每个要点应该是一到两句话的详细说明
   - 包含: 核心概念、关键步骤、注意事项、技巧说明等
   - 对于技术文档，要保留命令、配置项、API说明等重要信息
   - 对于教程类内容，要保留操作步骤和关键代码
4. tags: 生成3-5个相关标签

目标：生成一份内容丰富、信息完整、可以直接作为学习资料使用的详细笔记。

只返回 JSON,不要其他说明文字。`, title, content)
}

// parseSummary 解析总结结果
func (s *Summarizer) parseSummary(response string) (*Summary, error) {
	// gjson 可以直接解析,不需要提取 JSON
	// 它会自动找到第一个 JSON 对象

	// 解析各个字段
	titleResult := gjson.Get(response, "title")
	oneSentenceResult := gjson.Get(response, "one_sentence")
	keyPointsResult := gjson.Get(response, "key_points")
	tagsResult := gjson.Get(response, "tags")

	// 验证必填字段
	if !titleResult.Exists() || titleResult.String() == "" {
		return nil, fmt.Errorf("title 字段为空")
	}
	if !oneSentenceResult.Exists() || oneSentenceResult.String() == "" {
		return nil, fmt.Errorf("one_sentence 字段为空")
	}
	if !keyPointsResult.Exists() || !keyPointsResult.IsArray() {
		return nil, fmt.Errorf("key_points 字段为空或格式错误")
	}

	// 提取 key_points 数组
	var keyPoints []string
	keyPointsResult.ForEach(func(_, result gjson.Result) bool {
		keyPoints = append(keyPoints, result.String())
		return true
	})

	// 提取 tags 数组 (可选)
	var tags []string
	if tagsResult.Exists() && tagsResult.IsArray() {
		tagsResult.ForEach(func(_, result gjson.Result) bool {
			tags = append(tags, result.String())
			return true
		})
	}

	summary := &Summary{
		Title:       titleResult.String(),
		OneSentence: oneSentenceResult.String(),
		KeyPoints:   keyPoints,
		Tags:        tags,
	}

	return summary, nil
}
