package analytics

type GtagConfig struct {
	Enabled bool   `toml:"enabled"`
	UID     string `toml:"uid"`
}

func defaultGtagConfig() *GtagConfig {
	return &GtagConfig{
		Enabled: false,
		UID:     "",
	}
}
