package theme

// Config is the theme config.
type Config struct {
	Name             string   `toml:"name"`
	IndexTemplate    string   `toml:"index_template"`
	NotFoundTemplate string   `toml:"not_found_template"`
	Extension        []string `toml:"extension"`
	StaticDirs       []string `toml:"static_dirs"`
	EnableDarkMode   bool     `toml:"enable_dark_mode"`
	ShowPuGoVersion  bool     `toml:"show_pugo_version"`
}

// NewDefaultConfig returns default theme config
func NewDefaultConfig() *Config {
	return &Config{
		Name:             "theme",
		Extension:        []string{".html"},
		IndexTemplate:    "post-list.html",
		NotFoundTemplate: "404.html",
		StaticDirs:       []string{"static"},
		EnableDarkMode:   false,
		ShowPuGoVersion:  false,
	}
}
