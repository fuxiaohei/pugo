package comments

type Twikoo struct {
	Enabled bool   `toml:"enabled"`
	EID     string `toml:"eid"`
	CDN     string `toml:"cdn"`
}

func defaultTwikoo() *Twikoo {
	return &Twikoo{
		Enabled: false,
		EID:     "",
		CDN:     "https://cdn.jsdelivr.net/npm/twikoo@1.5.9/dist/twikoo.all.min.js",
	}
}
