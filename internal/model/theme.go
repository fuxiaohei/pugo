package model

// Theme is the theme of the site.
type Theme struct {
	Name         string `toml:"name"`
	Descripition string `toml:"descripition"`
	Directory    string `toml:"directory"`
}
