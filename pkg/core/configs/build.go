package configs

// Build is configuration for building site
type Build struct {
	OutputDir       string   `toml:"output_dir"`
	StaticAssetsDir []string `toml:"static_assets_dir"`
	PostLinkFormat  string   `toml:"post_link_format"`

	TagLinkFormat     string `toml:"tag_link_format"`
	TagPageLinkFormat string `toml:"tag_page_link_format"`

	PostPerPage        int    `toml:"post_per_page"`
	PostPageLinkFormat string `toml:"post_page_link_format"`

	ArchivesLink string `toml:"archive_link"`

	EnableMinifyHTML bool `toml:"enable_minify_html"`
}

// DefaultBuild returns a new default build config
func DefaultBuild() *Build {
	return &Build{
		OutputDir:       "./build",
		StaticAssetsDir: []string{"./assets"},
		PostLinkFormat:  "/{{.Date.Year}}/{{.Date.Month}}/{{.Slug}}/",

		TagLinkFormat:     "/tag/{{.Tag}}/",
		TagPageLinkFormat: "/tag/{{.Tag}}/{{.Page}}/",

		PostPerPage:        5,
		PostPageLinkFormat: "/page/{{.Page}}/",

		ArchivesLink: "/archives/",

		EnableMinifyHTML: true,
	}
}
