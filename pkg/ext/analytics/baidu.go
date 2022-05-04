package analytics

type BaiduConfig struct {
	Enabled bool   `toml:"enabled"`
	Hash    string `toml:"hash"`
}

func DefaultBaiduConfig() *BaiduConfig {
	return &BaiduConfig{
		Enabled: false,
		Hash:    "",
	}
}
