package analytics

type Config struct {
	GoogleAnalytics *GtagConfig      `toml:"google_analytics"`
	Plausible       *PlausibleConfig `toml:"plausible"`
	Baidu           *BaiduConfig     `toml:"baidu"`
}

func DefaultConfig() *Config {
	return &Config{
		GoogleAnalytics: DefaultGtagConfig(),
		Plausible:       DefaultPlausibleConfig(),
		Baidu:           DefaultBaiduConfig(),
	}
}
