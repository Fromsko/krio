package parser

import (
	"bufio"
	"io"
	"regexp"
	"strings"
)

// MdParser Markdown 文件解析器
type MdParser struct{}

// NewMdParser 创建 Markdown 解析器
func NewMdParser() *MdParser {
	return &MdParser{}
}

// Parse 解析 Markdown 文件，提取 URL
// 支持链接语法和纯文本 URL
func (p *MdParser) Parse(r io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(r)
	var urls []string

	// URL 正则表达式
	// 匹配 http:// 或 https:// 开头的 URL
	urlRegex := regexp.MustCompile(`https?://[^\s\]]+`)

	for scanner.Scan() {
		line := scanner.Text()

		// 提取所有 URL
		matches := urlRegex.FindAllString(line, -1)
		for _, match := range matches {
			// 清理 URL (移除可能的尾部标点)
			url := cleanURL(match)
			if url != "" {
				urls = append(urls, url)
			}
		}
	}

	return urls, scanner.Err()
}

// cleanURL 清理 URL
// 移除尾部可能的标点符号和 Markdown 语法字符
func cleanURL(rawURL string) string {
	// 移除尾部的标点符号
	url := strings.TrimRight(rawURL, ".,;:!?()[]{}\"'`")

	// 移除尾部的 / ) 等字符
	url = strings.TrimRight(url, "/)")

	return url
}
