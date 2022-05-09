package comments

type Config struct {
	Enabled bool    `toml:"enabled"`
	Disqus  *Disqus `toml:"disqus"`
	Valine  *Valine `toml:"valine"`
	Twikoo  *Twikoo `toml:"twikoo"`
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
	if c.Twikoo.Enabled {
		return "twikoo"
	}
	return "unknown"
}

func DefaultConfig() *Config {
	return &Config{
		Enabled: true,
		Disqus:  DefaultDisqusConfig(),
		Valine:  defaultValine(),
		Twikoo:  defaultTwikoo(),
	}
}
