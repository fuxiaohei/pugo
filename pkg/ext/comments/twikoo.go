package comments

type Twikoo struct {
	Enabled bool   `toml:"enabled"`
	EID     string `toml:"eid"`
}

func defaultTwikoo() *Twikoo {
	return &Twikoo{
		Enabled: false,
		EID:     "",
	}
}
