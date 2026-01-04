package parser

import (
	"io"
	"strings"
)

// Parser 解析器接口
type Parser interface {
	Parse(r io.Reader) ([]string, error)
}

// DetectFormat 检测文件格式
func DetectFormat(filename string) Parser {
	switch {
	case strings.HasSuffix(filename, ".txt"):
		return NewTxtParser()
	case strings.HasSuffix(filename, ".md"):
		return NewMdParser()
	default:
		return NewTxtParser() // 默认使用 TXT 解析器
	}
}
