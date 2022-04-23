package builder

import (
	"fmt"
	"pugo/internal/zlog"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

// Builder is the instance for building a site.
type Builder struct {
	configFile string

	render   *ThemeRender
	source   *SourceData
	markdown goldmark.Markdown

	outputDir string
}

// OutputDir returns the output directory.
func (b *Builder) OutputDir() string {
	return b.outputDir
}

// Options is the options for building a site.
type Option struct {
	ConfigFile string
	OutputDir  string
}

// NewBuilder returns a new Builder instance.
func NewBuilder(opt *Option) *Builder {
	return &Builder{
		configFile: opt.ConfigFile,
		outputDir:  opt.OutputDir,
		source:     NewDefaultSourceData(),
		markdown: goldmark.New(
			goldmark.WithExtensions(extension.GFM),
			goldmark.WithParserOptions(
				parser.WithAutoHeadingID(),
				parser.WithASTTransformers(
					util.Prioritized(newAstTransformer(), 10000),
				),
			),
			goldmark.WithRendererOptions(
				html.WithHardWraps(),
				html.WithXHTML(),
			),
		),
	}
}

// Build builds the site.
func (b *Builder) Build() {
	st := time.Now()
	if err := b.parseSource(); err != nil {
		zlog.Error("failed to parse source", "err", err)
		return
	}
	ctx, err := b.buildContents()
	if err != nil {
		zlog.Error("failed to build contents", "err", err)
		return
	}
	if err = b.Output(ctx); err != nil {
		zlog.Error("failed to output", "err", err)
		return
	}
	zlog.Info("build: finished", "files", ctx.outputLength(), "duration", time.Since(st).Milliseconds())
}

func (b *Builder) buildContents() (*buildContext, error) {
	ctx := newBuildContext(b.source)
	if ctx == nil {
		return nil, fmt.Errorf("failed to build contents context")
	}

	// set theme copy dir
	if err := b.render.updateCopyDirs(ctx); err != nil {
		zlog.Warn("theme: failed to update copy dirs", "err", err)
		return nil, err
	}

	// build files
	if err := b.buildPosts(ctx); err != nil {
		zlog.Warn("posts: failed to build", "err", err)
		return nil, err
	}
	if err := b.buildPostLists(ctx); err != nil {
		zlog.Warn("posts: failed to build lists", "err", err)
		return nil, err
	}
	if err := b.buildTags(ctx); err != nil {
		zlog.Warn("posts: failed to build tags", "err", err)
		return nil, err
	}
	if err := b.buildIndex(ctx); err != nil {
		zlog.Warn("posts: failed to build index", "err", err)
		return nil, err
	}
	if err := b.buildFeedAtom(ctx); err != nil {
		zlog.Warn("posts: failed to build feed atom", "err", err)
		return nil, err
	}
	if err := b.buildPages(ctx); err != nil {
		zlog.Warn("pages: failed to build", "err", err)
		return nil, err
	}

	return ctx, nil
}
