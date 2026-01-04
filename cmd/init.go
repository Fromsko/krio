package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/fromsko/krio/internal/config"
)

var forceInit bool

// initCmd 初始化配置命令
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "初始化配置文件",
	Long:  `在当前目录或用户目录创建默认配置文件。`,
	Run: func(cmd *cobra.Command, args []string) {
		// 确定配置文件路径
		configPath := getConfigPath()

		// 检查配置文件是否存在
		if _, err := os.Stat(configPath); err == nil && !forceInit {
			fmt.Printf("⚠️  配置文件已存在: %s\n", configPath)
			fmt.Println("使用 --force 强制覆盖")
			return
		}

		// 确保目录存在
		configDir := filepath.Dir(configPath)
		if err := os.MkdirAll(configDir, 0755); err != nil {
			fmt.Printf("❌ 创建目录失败: %v\n", err)
			return
		}

		// 创建默认配置
		if err := config.CreateDefault(configPath); err != nil {
			fmt.Printf("❌ 创建配置失败: %v\n", err)
			return
		}

		fmt.Printf("✅ 配置文件已创建: %s\n", configPath)
		fmt.Println("📝 请编辑配置文件,填入你的 API Key:")
		fmt.Println("   model.api_key: \"your-api-key-here\"")
	},
}

func getConfigPath() string {
	if cfgFile != "" {
		return cfgFile
	}

	// 使用标准位置: ~/.config/agent-sko/config.yaml
	return config.GetDefaultConfigPath()
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolVar(&forceInit, "force", false,
		"强制覆盖已存在的配置文件")
}
