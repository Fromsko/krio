package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config 应用配置
// 配置文件文档: references/configuration/api-keys.md
type Config struct {
	Model       ModelConfig       `yaml:"model"`
	ObsidianMCP ObsidianMCPConfig `yaml:"obsidian_mcp"`
	App         AppConfig         `yaml:"app"`
	Scraper     ScraperConfig     `yaml:"scraper"`
	Note        NoteConfig        `yaml:"note"`
	Logging     LoggingConfig     `yaml:"logging"`
}

// ModelConfig 模型配置
type ModelConfig struct {
	APIKey      string  `yaml:"api_key"`
	BaseURL     string  `yaml:"base_url"`
	ModelName   string  `yaml:"model_name"`
	Temperature float64 `yaml:"temperature"`
	MaxTokens   int     `yaml:"max_tokens"`
}

// ObsidianMCPConfig Obsidian MCP 服务器配置
type ObsidianMCPConfig struct {
	Transport string        `yaml:"transport"`
	Command   string        `yaml:"command"`
	Args      []string      `yaml:"args"`
	Timeout   time.Duration `yaml:"timeout"`
}

// AppConfig 应用配置
type AppConfig struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	Debug   bool   `yaml:"debug"`
}

// ScraperConfig 网页抓取配置
type ScraperConfig struct {
	UserAgent      string        `yaml:"user_agent"`
	Timeout        time.Duration `yaml:"timeout"`
	MaxRetries     int           `yaml:"max_retries"`
	RetryDelay     time.Duration `yaml:"retry_delay"`
	EnableCache    bool          `yaml:"enable_cache"`
	CacheTTL       time.Duration `yaml:"cache_ttl"`
	MaxConcurrency int           `yaml:"max_concurrency"`
}

// NoteConfig 笔记生成配置
type NoteConfig struct {
	DefaultFolder    string `yaml:"default_folder"`
	FilenameTemplate string `yaml:"filename_template"`
	AddTimestamp     bool   `yaml:"add_timestamp"`
}

// LoggingConfig 日志配置
type LoggingConfig struct {
	Level      string `yaml:"level"`
	Output     string `yaml:"output"`
	FilePath   string `yaml:"file_path"`
	JSONFormat bool   `yaml:"json_format"`
}

var globalConfig *Config

// Load 加载配置文件
// 配置文件路径参考: config.yaml
func Load(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 环境变量覆盖
	if apiKey := os.Getenv("MODEL_API_KEY"); apiKey != "" {
		cfg.Model.APIKey = apiKey
	}
	if baseURL := os.Getenv("MODEL_BASE_URL"); baseURL != "" {
		cfg.Model.BaseURL = baseURL
	}
	if modelName := os.Getenv("MODEL_NAME"); modelName != "" {
		cfg.Model.ModelName = modelName
	}

	globalConfig = &cfg
	return &cfg, nil
}

// LoadDefault 加载默认配置文件
func LoadDefault() (*Config, error) {
	// 尝试加载 config.yaml, 如果不存在则尝试 config.example.yaml
	if _, err := os.Stat("config.yaml"); err == nil {
		return Load("config.yaml")
	}
	return Load("config.example.yaml")
}

// Get 获取全局配置
func Get() *Config {
	return globalConfig
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.Model.APIKey == "" || c.Model.APIKey == "your-api-key-here" {
		return fmt.Errorf("model.api_key 未设置,请在 config.yaml 中配置或设置 MODEL_API_KEY 环境变量")
	}
	if c.Model.BaseURL == "" {
		return fmt.Errorf("model.base_url 未设置")
	}
	if c.Model.ModelName == "" {
		return fmt.Errorf("model.model_name 未设置")
	}
	return nil
}
