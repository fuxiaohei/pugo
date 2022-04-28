package constants

import "time"

const (
	// ConfigFile is the default config file name.
	ConfigFile = "config.toml"
)

const (
	ContentPostsDir = "content/posts"
	ContentPagesDir = "content/pages"
)

var (
	initDirectories = []string{
		ContentPostsDir,
		ContentPagesDir,
		"themes/default",
		"themes/default/static",
		"themes/default/partial",
		"build",
		"assets",
	}
)

func InitDirectories() []string {
	return initDirectories
}

const (
	// WatchPollingDuration is the duration of polling for file changes.
	WatchPollingDuration = time.Second / 2
	// WatchTickerDuaration is the duration of ticker for file changes.
	WatchTickerDuaration = time.Second
)
