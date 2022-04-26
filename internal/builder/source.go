package builder

import (
	"fmt"
	"pugo/internal/model"
	"pugo/internal/theme"
	"pugo/internal/utils"
	"pugo/internal/zlog"

	"github.com/tdewolff/minify/v2"
	mhtml "github.com/tdewolff/minify/v2/html"
)

// SourceData is the parsed source data.
type SourceData struct {
	Posts       []*model.Post
	PostsPager  *model.Pager
	Tags        []*model.TagPosts
	Pages       []*model.Page
	Config      *model.Config
	BuildConfig *model.BuildConfig
}

// NewDefaultSourceData returns a new default source data.
func NewDefaultSourceData() *SourceData {
	return &SourceData{
		Posts:  make([]*model.Post, 0),
		Pages:  make([]*model.Page, 0),
		Config: model.NewDefaultConfig(),
	}
}

// LoadConfig loads config file.
func LoadConfig(configFile string) (*model.Config, error) {
	config := model.NewDefaultConfig()
	if err := utils.LoadTOML(configFile, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %s", err)
	}
	return config, nil
}

func (b *Builder) parseConfig() error {
	if err := utils.LoadTOML(b.configFile, b.source.Config); err != nil {
		return fmt.Errorf("failed to parse config file: %s", err)
	}

	// move BuildConfig to top
	b.source.BuildConfig = b.source.Config.BuildConfig

	// override output directory if empty
	if b.outputDir == "" {
		b.outputDir = b.source.BuildConfig.OutputDir
	}
	if b.outputDir == "" {
		return fmt.Errorf("output directory is empty")
	}

	if err := b.source.Config.Check(); err != nil {
		return err
	}

	zlog.Info("config: parsed ok", "output", b.outputDir)
	return nil
}

func (b *Builder) parseTheme() error {
	r, err := theme.NewRender(b.source.Config.Theme)
	if err != nil {
		return err
	}
	b.render = r
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

	// prepare minifer
	if b.source.BuildConfig.EnableMinifyHTML {
		m := minify.New()
		m.Add("text/html", &mhtml.Minifier{
			KeepComments:            false,
			KeepConditionalComments: true,
			KeepDefaultAttrVals:     true,
			KeepDocumentTags:        true,
			KeepEndTags:             false,
			KeepQuotes:              true,
			KeepWhitespace:          false,
		})
		b.minifier = m
	}

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

	// set page author
	for _, page := range s.Pages {
		page.Author = s.assignAuthor(page.AuthorName)
	}
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
