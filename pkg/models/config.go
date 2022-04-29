package models

import (
	"fmt"
	"pugo/pkg/constants"
	"pugo/pkg/theme"
	"pugo/pkg/utils"
	"strings"
)

type (
	SiteConfig struct {
		Title       string   `toml:"title"`
		SubTitle    string   `toml:"sub_title"`
		Base        string   `toml:"base"`
		Description string   `toml:"description"`
		Keywords    []string `toml:"keywords"`
	}
	Config struct {
		Site        *SiteConfig  `toml:"site"`
		Menu        []*Menu      `toml:"menu"`
		Theme       *theme.Theme `toml:"theme"`
		BuildConfig *BuildConfig `toml:"build"`
		Author      []*Author    `toml:"author"`
	}
)

// GetAuthor gets the author by the given name
func (c *Config) GetAuthor(name string) *Author {
	for _, author := range c.Author {
		if author.Name == name || author.Slug == name {
			return author
		}
	}
	return nil
}

// DefaultConfig returns a new default config.
func DefaultConfig() *Config {
	return &Config{
		Site: &SiteConfig{
			Title:       "PuGo",
			SubTitle:    "a simple static site generator",
			Base:        "http://localhost:18080",
			Description: "a simple static site generator with markdown support",
			Keywords:    []string{"site", "generator", "markdown"},
		},
		Menu: []*Menu{
			{
				Title: "Home",
				Slug:  "/",
				Blank: false,
			},
			{
				Title: "Archives",
				Slug:  "/archives/",
				Blank: false,
			},
			{
				Title: "About",
				Slug:  "/about/",
				Blank: false,
			},
		},
		Author: []*Author{NewAuthor("admin")},
		Theme: &theme.Theme{
			Directory:  "./themes/default",
			ConfigFile: "theme_config.toml",
		},
		BuildConfig: DefaultBuildConfig(),
	}
}

// DefaultBuildConfig returns a new default build config
func DefaultBuildConfig() *BuildConfig {
	return &BuildConfig{
		OutputDir:       "./build",
		StaticAssetsDir: []string{"./assets"},
		PostLinkFormat:  "/{{.Date.Year}}/{{.Date.Month}}/{{.Slug}}/",

		TagLinkFormat:     "/tag/{{.Tag}}/",
		TagPageLinkFormat: "/tag/{{.Tag}}/{{.Page}}/",

		PostPerPage:        5,
		PostPageLinkFormat: "/page/{{.Page}}/",

		ArchivesLink: "/archives/",

		FeedPostLimit: constants.PostFeedLimit,

		EnableMinifyHTML: true,
	}
}

// LoadConfigFromFile loads config file.
func LoadConfigFromFile(item constants.ConfigFileItem) (*Config, error) {
	config := DefaultConfig()
	if item.Type == constants.ConfigTypeTOML {
		if err := utils.LoadTOMLFile(item.File, config); err != nil {
			return nil, fmt.Errorf("failed to parse config file: %s", err)
		}
		return config, nil
	}
	if item.Type == constants.ConfigTypeYAML {
		if err := utils.LoadYAMLFile(item.File, config); err != nil {
			return nil, fmt.Errorf("failed to parse config file: %s", err)
		}
		return config, nil
	}
	return nil, fmt.Errorf("unsupported config file type: %s", item.Type)
}

// FullURL returns the full url after the base.
func (sc *SiteConfig) FullURL(url string) string {
	return strings.TrimSuffix(sc.Base, "/") + "/" + strings.TrimPrefix(url, "/")
}
