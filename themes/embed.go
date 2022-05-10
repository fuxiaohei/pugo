package themes

import "embed"

//go:embed default/* docs/*
var DefaultAssets embed.FS
