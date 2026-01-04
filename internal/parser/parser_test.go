package parser

import (
	"strings"
	"testing"
)

func TestTxtParser_Parse(t *testing.T) {
	input := `# 注释行
https://example.com

https://test.com
# 另一个注释
https://github.com/test
`

	parser := NewTxtParser()
	urls, err := parser.Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(urls) != 3 {
		t.Errorf("Expected 3 URLs, got %d", len(urls))
	}

	expected := []string{
		"https://example.com",
		"https://test.com",
		"https://github.com/test",
	}

	for i, url := range urls {
		if url != expected[i] {
			t.Errorf("URL %d: expected %s, got %s", i, expected[i], url)
		}
	}
}

func TestMdParser_Parse(t *testing.T) {
	input := `# 学习计划

## 前端框架
- [React](https://react.dev/) - 现代 UI 库
- [Vue](https://vuejs.org/) - 渐进式框架

## 其他资源
详见: https://example.com/resources
测试: https://test.com/page
`

	parser := NewMdParser()
	urls, err := parser.Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(urls) != 4 {
		t.Errorf("Expected 4 URLs, got %d", len(urls))
		t.Logf("Parsed URLs: %v", urls)
	}

	// 验证包含的 URL (注意: cleanURL 会移除尾部斜杠)
	expectedURLs := []string{
		"https://react.dev",      // 尾部斜杠被移除
		"https://vuejs.org",      // 尾部斜杠被移除
		"https://example.com/resources",
		"https://test.com/page",
	}

	for i, expected := range expectedURLs {
		if i >= len(urls) {
			t.Errorf("Not enough URLs, expected at least %d", len(expectedURLs))
			break
		}
		if urls[i] != expected {
			t.Errorf("URL %d: expected %s, got %s", i, expected, urls[i])
		}
	}
}
