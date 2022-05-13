package configs

// Language is the language configuration
type Language struct {
	Enabled bool   `toml:"enabled"`
	Lang    string `toml:"lang"`
}

func defaultLanguage() *Language {
	return &Language{
		Enabled: true,
		Lang:    "en",
	}
}
