package comments

type Valine struct {
	Enabled   bool   `toml:"enabled"`
	APPID     string `toml:"app_id"`
	APPKey    string `toml:"app_key"`
	ServerURL string `toml:"server_url"`
	CDN       string `toml:"cdn"`
}

func defaultValine() *Valine {
	return &Valine{
		Enabled:   false,
		APPID:     "",
		APPKey:    "",
		ServerURL: "",
		CDN:       "https://unpkg.com/valine@latest/dist/Valine.min.js",
	}
}
