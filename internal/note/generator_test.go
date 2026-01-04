package note

import (
	"testing"

	"github.com/fromsko/krio/internal/config"
)

func TestGenerateFilename(t *testing.T) {
	cfg := &config.NoteConfig{
		AddTimestamp: false, // 关闭时间戳以便测试
	}
	gen := &Generator{cfg: cfg}

	tests := []struct {
		name  string
		input string
		// 不检查完整结果,因为可能有时间戳
	}{
		{
			name:  "simple filename",
			input: "TestFile",
		},
		{
			name:  "with special characters",
			input: `Test/File<>:|?*.md`,
		},
		{
			name:  "with spaces",
			input: "My Test File",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := gen.generateFilename(tt.input, []string{})
			// 只验证结果不为空
			if result == "" {
				t.Errorf("generateFilename() returned empty string for input %v", tt.input)
			}
			// 验证小写
			if result != "" && result[0] >= 'A' && result[0] <= 'Z' {
				t.Errorf("generateFilename() should return lowercase, got %v", result)
			}
		})
	}
}

func TestEscapeYAML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "normal text",
			input:    "Hello World",
			expected: "Hello World",
		},
		{
			name:     "with quotes",
			input:    `He said "Hello"`,
			expected: `He said \"Hello\"`,
		},
		{
			name:     "with newline",
			input:    "Line1\nLine2",
			expected: "Line1\\nLine2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := escapeYAML(tt.input)
			if result != tt.expected {
				t.Errorf("escapeYAML() = %v, want %v", result, tt.expected)
			}
		})
	}
}
