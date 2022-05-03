package feed

const (
	DefaultLimitNums = 10
)

type Config struct {
	Enabled   bool   `toml:"enabled"`
	LimitNums int    `toml:"limit_nums"`
	Link      string `toml:"link"`
}

func DefaultConfig() *Config {
	return &Config{
		Enabled:   true,
		LimitNums: DefaultLimitNums,
		Link:      "/atom.xml",
	}
}
