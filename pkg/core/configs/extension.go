package configs

import (
	"pugo/pkg/ext/feed"
	"pugo/pkg/ext/sitemap"
)

type Extension struct {
	Feed    *feed.Config    `toml:"feed"`
	Sitemap *sitemap.Config `toml:"sitemap"`
}
