package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/fromsko/krio/internal/config"
)

// configCmd 配置管理命令
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "配置管理",
	Long:  `管理和查看配置文件信息。`,
}

// configPathsCmd 显示配置文件路径
var configPathsCmd = &cobra.Command{
	Use:   "paths",
	Short: "显示配置文件路径",
	Long:  `显示所有可能的配置文件路径和优先级。`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("\n📁 配置文件路径优先级 (从高到低):")
		fmt.Println("=" + "=========================================")

		paths := config.GetConfigPaths()
		for i, path := range paths {
			fmt.Printf("%d. %s\n", i+1, path)
		}

		fmt.Println("\n💡 提示:")
		fmt.Println("  - 当前目录的 config.yaml 优先级最高")
		fmt.Println("  - 标准位置: ~/.config/agent-sko/config.yaml")
		fmt.Println("  - 使用 --config 指定配置文件可以覆盖所有默认路径")
		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configPathsCmd)
}
