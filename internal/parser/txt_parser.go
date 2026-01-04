package parser

import (
	"bufio"
	"io"
	"net/url"
	"strings"
)

// TxtParser TXT 文件解析器
type TxtParser struct{}

// NewTxtParser 创建 TXT 解析器
func NewTxtParser() *TxtParser {
	return &TxtParser{}
}

// Parse 解析 TXT 文件，提取 URL
// 支持注释 (# 开头) 和空行
func (p *TxtParser) Parse(r io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(r)
	var urls []string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// 跳过空行
		if line == "" {
			continue
		}

		// 跳过注释行
		if strings.HasPrefix(line, "#") {
			continue
		}

		// 验证是否为有效 URL
		if isValidURL(line) {
			urls = append(urls, line)
		}
	}

	return urls, scanner.Err()
}

// isValidURL 验证 URL 格式
func isValidURL(s string) bool {
	u, err := url.Parse(s)
	if err != nil {
		return false
	}

	// 检查协议和主机
	return u.Scheme == "http" || u.Scheme == "https"
}
