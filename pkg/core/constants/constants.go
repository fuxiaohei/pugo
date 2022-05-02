package constants

import "time"

// ConfigType is config file type
type ConfigType string

const (
	// ConfigTypeTOML is the config type for TOML
	ConfigTypeTOML ConfigType = "toml"
	// ConfigTypeYAML is the config type for YAML
	ConfigTypeYAML ConfigType = "yaml"
)

// ConfigFileItem returns a config file item
type ConfigFileItem struct {
	Type ConfigType
	File string
}

// ConfigFiles returns default config files list
func ConfigFiles() []ConfigFileItem {
	return []ConfigFileItem{
		{Type: ConfigTypeTOML, File: "config.toml"},
		{Type: ConfigTypeYAML, File: "config.yaml"},
	}
}

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
