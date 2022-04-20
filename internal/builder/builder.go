package builder

import (
	"fmt"
	"haisite/internal/model"
	"haisite/internal/zlog"
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

// Builder is the instance for building a site.
type Builder struct {
	configFile string
	config     *model.Config

	render     *Render
	parsedData *ParsedData

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
		config:     model.NewDefaultConfig(),
		outputDir:  opt.OutputDir,
		parsedData: new(ParsedData),
	}
}

// Build builds the site.
func (b *Builder) Build() {
	if err := b.parseConfig(); err != nil {
		zlog.Warn("config: failed to parse", "err", err)
		return
	}
	if err := b.parseTheme(); err != nil {
		zlog.Warn("theme: failed to parse", "err", err)
		return
	}
	if err := b.buildPosts(); err != nil {
		zlog.Warn("posts: failed to build", "err", err)
		return
	}
	if err := b.buildPages(); err != nil {
		zlog.Warn("pages: failed to build", "err", err)
		return
	}
}

func (b *Builder) parseConfig() error {
	fileBytes, err := ioutil.ReadFile(b.configFile)
	if err != nil {
		return fmt.Errorf("failed to read config file: %s", err)
	}
	if err = toml.Unmarshal(fileBytes, b.config); err != nil {
		return fmt.Errorf("failed to parse config file: %s", err)
	}

	// override output directory if empty
	if b.outputDir == "" {
		b.outputDir = b.config.BuildConfig.OutputDir
	}
	if b.outputDir == "" {
		return fmt.Errorf("output directory is empty")
	}

	// zlog.Debug("parsed config", "config", b.config)

	zlog.Info("config: parsed ok", "output", b.outputDir)
	return nil
}

func (b *Builder) parseTheme() error {
	r, err := NewRender(b.config.Theme.Directory, b.config.Theme.ConfigFile)
	if err != nil {
		return err
	}
	b.render = r
	return nil
}
