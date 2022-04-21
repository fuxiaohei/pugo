package model

import (
	"fmt"
	"net/url"
)

var (
	// ErrNoAuthor means no author found.
	ErrNoAuthor = fmt.Errorf("no author")
)

type (
	Config struct {
		Site        *SiteConfig  `toml:"site"`
		Menu        []*Menu      `toml:"menu"`
		Author      []*Author    `toml:"author"`
		Theme       *Theme       `toml:"theme"`
		BuildConfig *BuildConfig `toml:"build"`
	}
	SiteConfig struct {
		Title       string   `toml:"title"`
		SubTitle    string   `toml:"sub_title"`
		Base        string   `toml:"base"`
		Description string   `toml:"description"`
		Keywords    []string `toml:"keywords"`
	}
	Author struct {
		Name        string `toml:"name"`
		Email       string `toml:"email"`
		Website     string `toml:"website"`
		Avatar      string `toml:"avatar"`
		UseGravatar bool   `toml:"use_gravatar"`
		Slug        string `toml:"slug"`
	}
	// Theme is the theme of the site.
	Theme struct {
		Directory  string `toml:"directory"`
		ConfigFile string `toml:"config_file"`
	}
)

// NewDemoAuthor return a new author with demo fulfilled information.
func NewDemoAuthor(name string) *Author {
	return &Author{
		Name:        name,
		Email:       name + "@example.com",
		UseGravatar: true,
		Slug:        "/author/" + url.PathEscape(name),
	}
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
			Base:        "/",
			Description: "a simple static site generator with markdown support",
			Keywords:    []string{"site", "generator", "markdown"},
		},
		Menu: []*Menu{
			{
				Title: "Index",
				Slug:  "/index.html",
				Blank: false,
			},
			{
				Title: "Archive",
				Slug:  "/archive.html",
				Blank: false,
			},
			{
				Title: "About",
				Slug:  "/about.html",
				Blank: false,
			},
		},
		Author: []*Author{
			{
				Name:        "admin",
				Email:       "admin@example.com",
				Website:     "http://pugo.io",
				Avatar:      "",
				UseGravatar: true,
				Slug:        "admin",
			},
		},
		Theme: &Theme{
			Directory:  "./theme/default",
			ConfigFile: "theme_config.toml",
		},
		BuildConfig: NewDefaultBuildConfig(),
	}
}
