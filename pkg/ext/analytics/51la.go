package analytics

type V51LaConfig struct {
	Enabled bool   `toml:"enabled"`
	ID      string `toml:"id"`
	CK      string `toml:"ck"`
}

func defaultV51LaConfig() *V51LaConfig {
	return &V51LaConfig{
		Enabled: false,
		ID:      "",
		CK:      "",
	}
}
