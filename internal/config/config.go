package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
// 优先级: 当前目录 > .config/agent-sko/ > ~/.krio.yaml > config.example.yaml
func LoadDefault() (*Config, error) {
	// 1. 当前目录的 config.yaml (优先级最高,方便开发)
	if _, err := os.Stat("config.yaml"); err == nil {
		return Load("config.yaml")
	}

	// 2. 用户目录 .config/agent-sko/config.yaml (标准位置)
	homeDir, err := os.UserHomeDir()
	if err == nil {
		configDir := filepath.Join(homeDir, ".config", "agent-sko")
		configPath := filepath.Join(configDir, "config.yaml")
		if _, err := os.Stat(configPath); err == nil {
			return Load(configPath)
		}
	}

	// 3. 用户目录 .krio.yaml (兼容旧版本)
	if _, err := os.Stat(filepath.Join(homeDir, ".krio.yaml")); err == nil {
		return Load(filepath.Join(homeDir, ".krio.yaml"))
	}

	// 4. 最后尝试 config.example.yaml
	return Load("config.example.yaml")
}

// GetConfigPaths 获取所有可能的配置文件路径
func GetConfigPaths() []string {
	paths := []string{}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return paths
	}

	// 当前目录
	paths = append(paths, "config.yaml (当前目录)")

	// 标准位置
	configDir := filepath.Join(homeDir, ".config", "agent-sko")
	paths = append(paths, filepath.Join(configDir, "config.yaml"))

	// 兼容位置
	paths = append(paths, filepath.Join(homeDir, ".krio.yaml"))

	return paths
}

// GetDefaultConfigPath 获取默认配置文件路径
// 返回 .config/agent-sko/config.yaml 的完整路径
func GetDefaultConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "config.yaml"
	}
	configDir := filepath.Join(homeDir, ".config", "agent-sko")
	return filepath.Join(configDir, "config.yaml")
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

// Exists 检查配置文件是否存在
func Exists() bool {
	if _, err := os.Stat("config.yaml"); err == nil {
		return true
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false
	}
	configPath := filepath.Join(homeDir, ".krio.yaml")
	if _, err := os.Stat(configPath); err == nil {
		return true
	}
	return false
}

// CreateDefault 创建默认配置文件
func CreateDefault(path string) error {
	defaultConfig := `# Krio 配置文件
# 生成时间: {{timestamp}}

# 模型配置
model:
  # 智云 API Key (必填)
  api_key: "your-api-key-here"
  # API Base URL
  base_url: "https://open.bigmodel.cn/api/coding/paas/v4"
  # 使用的模型
  model_name: "glm-4.7"
  # 温度参数 (0.0 - 1.0)
  temperature: 0.7
  # 最大 token 数
  max_tokens: 4096

# Obsidian MCP 服务器配置
obsidian_mcp:
  # MCP 服务器传输方式
  transport: "stdio"
  # 执行命令
  command: "bun"
  # 命令参数
  args:
    - "x"
    - "--no-cache"
    - "@andysama/obsidian-mcp-server"
    - "--vault"
    - "D:/notes/Fromsko"
  # 超时时间
  timeout: 30s

# 应用配置
app:
  # 应用名称
  name: "Web Note Agent"
  # 版本
  version: "1.0.0"
  # 调试模式
  debug: true

# 网页抓取配置
scraper:
  # 用户代理
  user_agent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
  # 请求超时
  timeout: 15s
  # 最大重试次数
  max_retries: 3
  # 重试延迟
  retry_delay: 1000ms
  # 性能优化配置
  enable_cache: true        # 启用缓存
  cache_ttl: 1h            # 缓存过期时间
  max_concurrency: 5       # 最大并发数

# 笔记生成配置
note:
  # 默认保存文件夹
  default_folder: "Inbox"
  # 笔记命名模板
  filename_template: "{{title}}-{{timestamp}}"
  # 是否添加时间戳
  add_timestamp: true

# 日志配置
logging:
  # 日志级别: debug, info, warn, error
  level: "info"
  # 日志输出: console, file
  output: "console"
  # 日志文件路径 (当 output=file 时)
  file_path: "logs/app.log"
  # 是否启用 JSON 格式
  json_format: false
`

	// 替换时间戳
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	defaultConfig = strings.ReplaceAll(defaultConfig, "{{timestamp}}", timestamp)

	// 写入文件
	if err := os.WriteFile(path, []byte(defaultConfig), 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}

	return nil
}
