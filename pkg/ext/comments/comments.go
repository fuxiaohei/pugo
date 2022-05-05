package comments

type Config struct {
	Enabled bool    `toml:"enabled"`
	Disqus  *Disqus `toml:"disqus"`
	Valine  *Valine `toml:"valine"`
}

// Current returns the name of the enabled comment system.
// order by disqus, valine
func (c *Config) Current() string {
	if c.Disqus.Enabled {
		return "disqus"
	}
	if c.Valine.Enabled {
		return "valine"
	}
	return "unknown"
}

func DefaultConfig() *Config {
	return &Config{
		Enabled: false,
		Disqus:  DefaultDisqusConfig(),
		Valine:  defaultValine(),
	}
}
