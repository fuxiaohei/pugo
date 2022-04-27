package constants

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
		"build",
		"assets",
	}
)

func InitDirectories() []string {
	return initDirectories
}
