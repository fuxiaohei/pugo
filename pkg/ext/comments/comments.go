package comments

type Config struct {
	Enabled bool    `toml:"enabled"`
	Disqus  *Disqus `toml:"disqus"`
}

func DefaultConfig() *Config {
	return &Config{
		Enabled: false,
		Disqus:  DefaultDisqusConfig(),
	}
}
