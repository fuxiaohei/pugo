package model

// BuildConfig is configuration for building site
type BuildConfig struct {
	OutputDir       string `toml:"output_dir"`
	StaticAssetsDir string `toml:"static_assets_dir"`
	PostLinkFormat  string `toml:"post_link_format"`

	TagLinkFormat     string `toml:"tag_link_format"`
	TagPageLinkFormat string `toml:"tag_page_link_format"`

	PostPerPage        int    `toml:"post_per_page"`
	PostPageLinkFormat string `toml:"post_page_link_format"`
}

// NewDefaultBuildConfig returns a new default build config
func NewDefaultBuildConfig() *BuildConfig {
	return &BuildConfig{
		OutputDir:       "./build",
		StaticAssetsDir: "./assets",
		PostLinkFormat:  "/{{.Date.Year}}/{{.Date.Month}}/{{.Slug}}.html",

		TagLinkFormat:     "/tag/{{.Tag}}.html",
		TagPageLinkFormat: "/tag/{{.Tag}}/{{.Page}}.html",

		PostPerPage:        5,
		PostPageLinkFormat: "/page/{{.Page}}.html",
	}
}
