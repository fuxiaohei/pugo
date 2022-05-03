package configs

import "pugo/pkg/ext/feed"

type Extension struct {
	Feed *feed.Config `toml:"feed"`
}
