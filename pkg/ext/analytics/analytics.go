package analytics

type Config struct {
	GoogleAnalytics *GtagConfig      `toml:"google_analytics"`
	Plausible       *PlausibleConfig `toml:"plausible"`
}

func DefaultConfig() *Config {
	return &Config{
		GoogleAnalytics: DefaultGtagConfig(),
		Plausible:       DefaultPlausibleConfig(),
	}
}
