package cmd

import (
	"fmt"

	"github.com/fromsko/krio/app"
	"github.com/spf13/cobra"
)

// versionCmd 版本命令
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "显示版本信息",
	Long:  `显示 Krio 的版本信息、Git 提交和构建日期。`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("\n🦊 Krio - 智能网页笔记 Agent")
		fmt.Println("=" + "===========================")
		fmt.Printf("版本:     %s\n", app.Version)
		fmt.Printf("提交:     %s\n", app.Commit)
		fmt.Printf("构建日期: %s\n", app.BuildDate)
		fmt.Println("=" + "===========================")
		fmt.Println("\n项目地址: https://github.com/fromsko/krio")
		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
