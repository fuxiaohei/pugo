package configs

import (
	"pugo/pkg/ext/analytics"
	"pugo/pkg/ext/feed"
	"pugo/pkg/ext/sitemap"
)

type Extension struct {
	Feed      *feed.Config      `toml:"feed"`
	Sitemap   *sitemap.Config   `toml:"sitemap"`
	Analytics *analytics.Config `toml:"analytics"`
}

func defaultExtension() *Extension {
	return &Extension{
		Feed:      feed.DefaultConfig(),
		Sitemap:   sitemap.DefaultConfig(),
		Analytics: analytics.DefaultConfig(),
	}
}
