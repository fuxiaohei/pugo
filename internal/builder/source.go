package builder

import (
	"fmt"
	"io/ioutil"
	"pugo/internal/model"
	"pugo/internal/zlog"

	"github.com/BurntSushi/toml"
)

// SourceData is the parsed source data.
type SourceData struct {
	Posts      []*model.Post
	PostsPager *model.Pager
	Tags       []*model.TagPosts
	Pages      []*model.Page
	Config     *model.Config
}

// NewDefaultSourceData returns a new default source data.
func NewDefaultSourceData() *SourceData {
	return &SourceData{
		Posts:  make([]*model.Post, 0),
		Pages:  make([]*model.Page, 0),
		Config: model.NewDefaultConfig(),
	}
}

func (b *Builder) parseConfig() error {
	fileBytes, err := ioutil.ReadFile(b.configFile)
	if err != nil {
		return fmt.Errorf("failed to read config file: %s", err)
	}
	if err = toml.Unmarshal(fileBytes, b.source.Config); err != nil {
		return fmt.Errorf("failed to parse config file: %s", err)
	}

	// override output directory if empty
	if b.outputDir == "" {
		b.outputDir = b.source.Config.BuildConfig.OutputDir
	}
	if b.outputDir == "" {
		return fmt.Errorf("output directory is empty")
	}

	if err = b.source.Config.Check(); err != nil {
		return err
	}

	// zlog.Debug("parsed config", "config", b.config)

	zlog.Info("config: parsed ok", "output", b.outputDir)
	return nil
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
	if err := b.parsePosts(); err != nil {
		zlog.Warn("posts: failed to build", "err", err)
		return err
	}
	if err := b.parsePages(); err != nil {
		zlog.Warn("pages: failed to build", "err", err)
		return err
	}

	// make relative data available
	b.source.FulFill()

	return nil
}

// FulFill makes relative data available in source data
func (s *SourceData) FulFill() {

	// set post author data
	for _, post := range s.Posts {
		post.Author = s.assignAuthor(post.AuthorName)
		for _, t := range post.Tags {
			post.TagLinks = append(post.TagLinks, &model.TagLink{Name: t})
		}
	}

	// build tag posts
	s.Tags = model.BuildTagPosts(s.Posts)
	zlog.Info("posts: parsed tags ok", "tags", len(s.Tags))
}

func (s *SourceData) assignAuthor(name string) *model.Author {
	if name == "" {
		return s.Config.Author[0]
	}
	author := s.Config.GetAuthor(name)
	if author == nil {
		author = model.NewDemoAuthor(name)
	}
	return author
}
