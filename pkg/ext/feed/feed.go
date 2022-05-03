package feed

import (
	"bytes"
	"encoding/xml"
	"path/filepath"
	"pugo/pkg/core/models"
	"pugo/pkg/utils/zlog"
)

// RenderParams represents the parameters for rendering the feed.
type RenderParams struct {
	Config      *Config
	Posts       []*models.Post
	SiteBaseURL string
	SiteTitle   string
	OutputDir   string
}

// Render renders the feed.
func Render(params *RenderParams) (*models.OutputFile, error) {
	if params == nil || !params.Config.Enabled {
		zlog.Debugf("feed is disabled")
		return nil, nil
	}
	var posts []*models.Post
	var limit = params.Config.LimitNums
	if limit <= 0 {
		limit = DefaultLimitNums
	}
	if limit > len(params.Posts) {
		posts = params.Posts
	} else {
		posts = params.Posts[:limit]
	}
	feed := BuildAtom(params.Config.Link, posts, params.SiteBaseURL, params.SiteTitle)
	data, err := xml.Marshal(feed)
	if err != nil {
		zlog.Warnf("failed to marshal atom feed: %s", err)
		return nil, err
	}
	return &models.OutputFile{
		Path: filepath.Join(params.OutputDir, params.Config.Link),
		Buf:  bytes.NewBuffer(data),
	}, nil
}
