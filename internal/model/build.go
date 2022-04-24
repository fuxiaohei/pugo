package model

const (
	// DefaultFeedPostLimit is the default limit of feed post
	DefaultFeedPostLimit = 10
)

// BuildConfig is configuration for building site
type BuildConfig struct {
	OutputDir       string   `toml:"output_dir"`
	StaticAssetsDir []string `toml:"static_assets_dir"`
	PostLinkFormat  string   `toml:"post_link_format"`

	TagLinkFormat     string `toml:"tag_link_format"`
	TagPageLinkFormat string `toml:"tag_page_link_format"`

	PostPerPage        int    `toml:"post_per_page"`
	PostPageLinkFormat string `toml:"post_page_link_format"`

	ArchivesLink string `toml:"archive_link"`

	FeedPostLimit int `toml:"feed_post_limit"`

	EnableMinifyHTML bool `toml:"enable_minify_html"`
}

// NewDefaultBuildConfig returns a new default build config
func NewDefaultBuildConfig() *BuildConfig {
	return &BuildConfig{
		OutputDir:       "./build",
		StaticAssetsDir: []string{"./assets"},
		PostLinkFormat:  "/{{.Date.Year}}/{{.Date.Month}}/{{.Slug}}/",

		TagLinkFormat:     "/tag/{{.Tag}}/",
		TagPageLinkFormat: "/tag/{{.Tag}}/{{.Page}}/",

		PostPerPage:        5,
		PostPageLinkFormat: "/page/{{.Page}}/",

		ArchivesLink: "/archives/",

		FeedPostLimit: DefaultFeedPostLimit,

		EnableMinifyHTML: true,
	}
}
