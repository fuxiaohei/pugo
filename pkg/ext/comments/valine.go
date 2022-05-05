package comments

type Valine struct {
	Enabled   bool   `toml:"enabled"`
	APPID     string `toml:"app_id"`
	APPKey    string `toml:"app_key"`
	ServerURL string `toml:"server_url"`
}

func defaultValine() *Valine {
	return &Valine{
		Enabled:   false,
		APPID:     "",
		APPKey:    "",
		ServerURL: "",
	}
}
