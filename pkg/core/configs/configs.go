package configs

import (
	"fmt"
	"pugo/pkg/core/constants"
	"pugo/pkg/core/models"
	"pugo/pkg/core/theme"
	"pugo/pkg/ext/feed"
	"pugo/pkg/utils"
)

type Config struct {
	Site      *Site            `toml:"site"`
	Menu      []*models.Menu   `toml:"menu"`
	Theme     *theme.Theme     `toml:"theme"`
	Build     *Build           `toml:"build"`
	Author    []*models.Author `toml:"author"`
	Extension *Extension       `toml:"extension"`
}

// GetAuthor gets the author by the given name
func (c *Config) GetAuthor(name string) *models.Author {
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
		Site: &Site{
			Title:       "PuGo",
			SubTitle:    "a simple static site generator",
			Base:        "http://localhost:18080",
			Description: "a simple static site generator with markdown support",
			Keywords:    []string{"site", "generator", "markdown"},
		},
		Menu: []*models.Menu{
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
		Author: []*models.Author{models.NewAuthor("admin")},
		Theme: &theme.Theme{
			Directory:  "./themes/default",
			ConfigFile: "theme_config.toml",
		},
		Build: DefaultBuild(),
		Extension: &Extension{
			Feed: feed.DefaultConfig(),
		},
	}
}

// LoadFromFile loads config file.
func LoadFromFile(item constants.ConfigFileItem) (*Config, error) {
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
