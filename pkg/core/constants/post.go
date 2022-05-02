package constants

const (
	PostFeedLimit    = 10
	PostTemplate     = "post.html"
	PostListTemplate = "post-list.html"
	ArchivesTemplate = "archives.html"
)

type postMetaSeperator struct {
	StartChars []byte
	EndChars   []byte
	MetaType   ConfigType
}

var (
	postMetaSeperatorList = []postMetaSeperator{
		{
			StartChars: []byte("---\n"),
			EndChars:   []byte("---\n"),
			MetaType:   ConfigTypeYAML,
		},
		{
			StartChars: []byte("```toml\n"),
			EndChars:   []byte("```\n"),
			MetaType:   ConfigTypeTOML,
		},
	}
	postDateLayouts = []string{
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006-01-02",
	}
	postBriefSeperator = []byte("<!--more-->")
)

// PostMetaSeperators returns seperators of post meta.
func PostMetaSeperators() []postMetaSeperator {
	return postMetaSeperatorList
}

// PostDateLayouts returns date layouts of post.
func PostDateLayouts() []string {
	return postDateLayouts
}

// PostBriefSeperator returns seperator of post brief.
func PostBriefSeperator() []byte {
	return postBriefSeperator
}
