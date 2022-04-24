package theme

// Config is the theme config.
type Config struct {
	Name          string   `toml:"name"`
	IndexTemplate string   `toml:"index_template"`
	Extension     []string `toml:"extension"`
	StaticDirs    []string `toml:"static_dirs"`
}

// NewDefaultConfig returns default theme config
func NewDefaultConfig() *Config {
	return &Config{
		Name:          "theme",
		Extension:     []string{".html"},
		IndexTemplate: "post-list.html",
		StaticDirs:    []string{"static"},
	}
}
