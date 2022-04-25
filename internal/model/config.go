package model

import (
	"fmt"
	"strings"
)

var (
	// ErrNoAuthor means no author found.
	ErrNoAuthor = fmt.Errorf("no author")
)

type (
	Config struct {
		Site        *SiteConfig  `toml:"site"`
		Menu        []*Menu      `toml:"menu"`
		Theme       *Theme       `toml:"theme"`
		BuildConfig *BuildConfig `toml:"build"`
		Author      []*Author    `toml:"author"`
	}
	SiteConfig struct {
		Title       string   `toml:"title"`
		SubTitle    string   `toml:"sub_title"`
		Base        string   `toml:"base"`
		Description string   `toml:"description"`
		Keywords    []string `toml:"keywords"`
	}
	// Theme is the theme of the site.
	Theme struct {
		Directory  string `toml:"directory"`
		ConfigFile string `toml:"config_file"`
	}
)

// FullURL returns the full url after the base.
func (sc *SiteConfig) FullURL(url string) string {
	return strings.TrimSuffix(sc.Base, "/") + "/" + strings.TrimPrefix(url, "/")
}

// GetAuthor gets the author by the given name
func (c *Config) GetAuthor(name string) *Author {
	for _, author := range c.Author {
		if author.Name == name || author.Slug == name {
			return author
		}
	}
	return nil
}

// Check checks the config is valid for building contents.
func (c *Config) Check() error {
	if len(c.Author) == 0 {
		return ErrNoAuthor
	}
	return nil
}

// NewDefaultConfig returns a new default config.
func NewDefaultConfig() *Config {
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
		Author: []*Author{NewDemoAuthor("admin")},
		Theme: &Theme{
			Directory:  "./themes/default",
			ConfigFile: "theme_config.toml",
		},
		BuildConfig: NewDefaultBuildConfig(),
	}
}
