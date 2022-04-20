package builder

import (
	"pugo/internal/model"
	"pugo/internal/zlog"
)

// SourceData is the parsed source data.
type SourceData struct {
	Posts  []*model.Post
	Pages  []*model.Page
	Config *model.Config
}

// NewDefaultSourceData returns a new default source data.
func NewDefaultSourceData() *SourceData {
	return &SourceData{
		Posts:  make([]*model.Post, 0),
		Pages:  make([]*model.Page, 0),
		Config: model.NewDefaultConfig(),
	}
}

func (b *Builder) parseSource() error {
	if err := b.parseConfig(); err != nil {
		zlog.Warn("config: failed to parse", "err", err)
		return err
	}
	if err := b.parseTheme(); err != nil {
		zlog.Warn("theme: failed to parse", "err", err)
		return err
	}
	if err := b.buildPosts(); err != nil {
		zlog.Warn("posts: failed to build", "err", err)
		return err
	}
	if err := b.buildPages(); err != nil {
		zlog.Warn("pages: failed to build", "err", err)
		return err
	}
	return nil
}
