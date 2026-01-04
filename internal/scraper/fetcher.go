package scraper

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/fromsko/krio/internal/config"
	"github.com/gocolly/colly/v2"
)

// WebPage 网页内容
type WebPage struct {
	URL     string
	Title   string
	Content string
}

// Fetcher 网页抓取器
type Fetcher struct {
	cfg *config.ScraperConfig
}

// NewFetcher 创建抓取器
func NewFetcher(cfg *config.ScraperConfig) *Fetcher {
	return &Fetcher{cfg: cfg}
}

// Fetch 抓取网页内容
func (f *Fetcher) Fetch(urlStr string) (*WebPage, error) {
	// 验证 URL
	if err := f.validateURL(urlStr); err != nil {
		return nil, fmt.Errorf("URL 验证失败: %w", err)
	}

	var page *WebPage
	var fetchErr error

	// 重试逻辑
	for i := 0; i <= f.cfg.MaxRetries; i++ {
		page, fetchErr = f.fetchOnce(urlStr)
		if fetchErr == nil {
			return page, nil
		}

		if i < f.cfg.MaxRetries {
			time.Sleep(f.cfg.RetryDelay)
		}
	}

	return nil, fmt.Errorf("抓取失败(已重试 %d 次): %w", f.cfg.MaxRetries, fetchErr)
}

// fetchOnce 单次抓取
func (f *Fetcher) fetchOnce(urlStr string) (*WebPage, error) {
	c := colly.NewCollector(
		colly.UserAgent(f.cfg.UserAgent),
		colly.MaxDepth(1),
		colly.Async(false),
		// debug.Debugger(&debug.LogDebugger{}), // 调试时启用
	)

	// 设置超时
	c.SetRequestTimeout(time.Duration(f.cfg.Timeout) * time.Second)

	page := &WebPage{}
	var errorMsg string

	// 抓取标题
	c.OnHTML("title", func(e *colly.HTMLElement) {
		page.Title = strings.TrimSpace(e.Text)
	})

	// 如果没有 title 标签,尝试 h1
	c.OnHTML("h1", func(e *colly.HTMLElement) {
		if page.Title == "" {
			page.Title = strings.TrimSpace(e.Text)
		}
	})

	// 抓取主要内容
	c.OnHTML("body", func(e *colly.HTMLElement) {
		// 移除不需要的元素
		e.ForEach("script, style, nav, header, footer, iframe, noscript", func(i int, el *colly.HTMLElement) {
			_ = el.DOM.Remove()
		})

		// 提取文本内容
		page.Content = strings.TrimSpace(e.Text)
	})

	// 错误处理
	c.OnError(func(r *colly.Response, err error) {
		errorMsg = fmt.Sprintf("请求失败: %s", err)
	})

	// 设置 URL
	page.URL = urlStr

	// 开始抓取
	if err := c.Visit(urlStr); err != nil {
		return nil, err
	}

	c.Wait()

	if errorMsg != "" {
		return nil, fmt.Errorf("%s", errorMsg)
	}

	// 验证内容
	if page.Content == "" {
		return nil, fmt.Errorf("未获取到内容")
	}

	// 限制内容长度 (避免 token 过多)
	if len(page.Content) > 50000 {
		page.Content = page.Content[:50000] + "\n\n...(内容过长,已截断)"
	}

	return page, nil
}

// validateURL 验证 URL
func (f *Fetcher) validateURL(urlStr string) error {
	u, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("URL 格式错误: %w", err)
	}

	// 检查协议
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("不支持的协议: %s", u.Scheme)
	}

	// 防止 SSRF 攻击
	if f.isPrivateURL(u) {
		return fmt.Errorf("不允许访问私有地址: %s", urlStr)
	}

	return nil
}

// isPrivateURL 检查是否为私有地址
func (f *Fetcher) isPrivateURL(u *url.URL) bool {
	host := u.Hostname()

	// localhost
	if host == "localhost" || host == "127.0.0.1" {
		return true
	}

	// 私有 IP 段
	privatePrefixes := []string{
		"10.", "172.16.", "172.17.", "172.18.", "172.19.",
		"172.20.", "172.21.", "172.22.", "172.23.", "172.24.",
		"172.25.", "172.26.", "172.27.", "172.28.", "172.29.",
		"172.30.", "172.31.", "192.168.",
	}

	for _, prefix := range privatePrefixes {
		if strings.HasPrefix(host, prefix) {
			return true
		}
	}

	return false
}
