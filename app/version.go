package app

import "fmt"

var (
	// Name 应用名称
	Name = "Krio"
	// Version 应用版本 (构建时通过 ldflags 注入)
	Version = "dev"
	// Commit 提交 hash
	Commit = "unknown"
	// BuildDate 构建日期
	BuildDate = "unknown"
)

// String 返回版本信息字符串
func String() string {
	return fmt.Sprintf("%s version %s (commit: %s, built at: %s)",
		Name, Version, Commit, BuildDate)
}
