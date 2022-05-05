package comments

type Disqus struct {
	Enabled   bool   `toml:"enabled"`
	Shortname string `toml:"shortname"`
}

func DefaultDisqusConfig() *Disqus {
	return &Disqus{
		Enabled:   false,
		Shortname: "",
	}
}
