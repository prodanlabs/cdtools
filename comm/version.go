package comm

import (
	"fmt"
	"runtime"
)

var (
	// 初始化为 unknown，如果编译时没有传入这些值，则为 unknown
	Version        = "unknown"
	GitCommitLog   = "unknown"
	BuildTime      = "unknown"
	BuildGoVersion = "unknown"
)

// 返回多行格式
func StringifyMultiLine() string {
	return fmt.Sprintf("Version=%s\nGitCommitLog=%s\nBuildTime=%s\nGoVersion=%s\nruntime=%s/%s\n",
		Version, GitCommitLog, BuildTime, BuildGoVersion, runtime.GOOS, runtime.GOARCH)
}