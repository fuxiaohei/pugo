package configs

import "strings"

type Site struct {
	Title       string   `toml:"title"`
	SubTitle    string   `toml:"sub_title"`
	Base        string   `toml:"base"`
	Description string   `toml:"description"`
	Keywords    []string `toml:"keywords"`
}

// FullURL returns the full url after the base.
func (sc *Site) FullURL(url string) string {
	return strings.TrimSuffix(sc.Base, "/") + "/" + strings.TrimPrefix(url, "/")
}
