package analytics

type Config struct {
	GoogleAnalytics *GtagConfig      `toml:"google_analytics"`
	Plausible       *PlausibleConfig `toml:"plausible"`
	Baidu           *BaiduConfig     `toml:"baidu"`
	V51La           *V51LaConfig     `toml:"v51la"`
}

func DefaultConfig() *Config {
	return &Config{
		GoogleAnalytics: defaultGtagConfig(),
		Plausible:       defaultPlausibleConfig(),
		Baidu:           defaultBaiduConfig(),
		V51La:           defaultV51LaConfig(),
	}
}
