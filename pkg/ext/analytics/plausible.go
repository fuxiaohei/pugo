package analytics

type PlausibleConfig struct {
	Enabled    bool   `toml:"enabled"`
	DataDomain string `toml:"data_domain"`
}

func defaultPlausibleConfig() *PlausibleConfig {
	return &PlausibleConfig{
		Enabled:    false,
		DataDomain: "",
	}
}
