package model

type (
	Config struct {
		Site        *SiteConfig  `toml:"site"`
		Author      []*Author    `toml:"author"`
		Theme       *Theme       `toml:"theme"`
		BuildConfig *BuildConfig `toml:"build"`
	}
	SiteConfig struct {
		Title        string   `toml:"title"`
		SubTitle     string   `toml:"sub_title"`
		Base         string   `toml:"base"`
		Descripition string   `toml:"description"`
		Keywords     []string `toml:"keywords"`
	}
	Author struct {
		Name        string `toml:"name"`
		Email       string `toml:"email"`
		Website     string `toml:"website"`
		Avatar      string `toml:"avatar"`
		UseGravatar bool   `toml:"use_gravatar"`
		Slug        string `toml:"slug"`
	}
	BuildConfig struct {
		BuildDir        string `toml:"build_dir"`
		StaticAssetsDir string `toml:"static_assets_dir"`
	}
)

func NewDefaultConfig() *Config {
	return &Config{
		Site: &SiteConfig{
			Title:        "Haisite",
			SubTitle:     "a simple static site generator",
			Base:         "/",
			Descripition: "a simple static site generator with markdown support",
			Keywords:     []string{"site", "generator", "markdown"},
		},
		Author: []*Author{
			{
				Name:        "Hai",
				Email:       "hai@example.com",
				Website:     "haisite.com",
				Avatar:      "",
				UseGravatar: true,
				Slug:        "/author/{.Name}/",
			},
		},
		Theme: &Theme{
			Name:         "default",
			Descripition: "default theme",
			Directory:    "/themes/default",
		},
		BuildConfig: &BuildConfig{
			BuildDir:        "./build",
			StaticAssetsDir: "./assets",
		},
	}
}
