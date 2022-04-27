package models

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
