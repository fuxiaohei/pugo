package sitemap

type Config struct {
	Enabled bool   `toml:"enabled"`
	Link    string `toml:"link"`
}

func DefaultConfig() *Config {
	return &Config{
		Enabled: true,
		Link:    "/sitemap.xml",
	}
}
