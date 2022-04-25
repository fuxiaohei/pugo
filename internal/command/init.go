package command

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"pugo/internal/model"
	"pugo/internal/utils"
	"pugo/internal/zlog"
	"pugo/themes"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/urfave/cli/v2"

	_ "embed"
)

// NewInitCommand returns a new cli.Command for the init subcommand.
func NewInitCommand() *cli.Command {
	cmd := &cli.Command{
		Name:  "init",
		Usage: "initialize a new sample site in the current directory",
		Flags: GlobalFlags,
		Action: func(c *cli.Context) error {

			initGlobalFlags(c)

			// get directory name
			dir := "./"
			if c.Args().Len() > 0 {
				dir = filepath.Join(dir, c.Args().Get(0))
			}
			zlog.Info("init directory", "dir", dir)

			// create site directory
			if err := utils.MkdirAll(dir); err != nil {
				zlog.Warn("failed to create directory", "dir", dir, "err", err)
				return err
			}

			// create default config file
			if err := createDefaultConfigFile(dir); err != nil {
				zlog.Warn("failed to create default config file", "err", err)
				return err
			}

			// create init directory
			if err := createDirectories(dir); err != nil {
				return err
			}

			// create demo contents
			if err := createDemoContents(dir); err != nil {
				return err
			}

			// create demo post
			if err := createDemoPost(dir); err != nil {
				return err
			}

			// create demo page
			if err := createDemoPage(dir); err != nil {
				return err
			}

			return nil
		},
	}
	return cmd
}

const (
	// DefaultConfigFile is the default config file name.
	DefaultConfigFile = "config.toml"
)

var (
	initDirectories = []string{
		"content/posts",
		"content/pages",
		"themes/default",
		"build",
		"assets",
	}
)

func createDefaultConfigFile(dir string) error {
	configFile := filepath.Join(dir, DefaultConfigFile)
	if utils.IsFileExist(configFile) {
		//return fmt.Errorf("config file %s already exists", configFile)
	}
	cfg := model.NewDefaultConfig()

	buffer := bytes.NewBuffer(nil)
	if err := toml.NewEncoder(buffer).Encode(cfg); err != nil {
		return err
	}

	return ioutil.WriteFile(configFile, buffer.Bytes(), 0644)
}

func createDirectories(topDir string) error {
	for _, dir := range initDirectories {
		realDir := filepath.Join(topDir, dir)
		if err := utils.MkdirAll(realDir); err != nil {
			zlog.Warn("failed to create directory", "dir", realDir, "err", err)
			return err
		}
		zlog.Debug("create directory", "dir", realDir)
	}
	return nil
}

func createDemoContents(topDir string) error {
	zlog.Debug("create demo contents")

	// extract default theme
	if err := extractThemeDir(filepath.Join(topDir, "themes/default"), "default"); err != nil {
		zlog.Warn("failed to extract default theme", "err", err)
		return err
	}

	return nil
}

func extractThemeDir(topDir, dir string) error {
	files, err := themes.DefaultAssets.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, file := range files {
		if file.IsDir() {
			dirName := filepath.Join(topDir, file.Name())
			if err := utils.MkdirAll(dirName); err != nil {
				zlog.Warn("failed to create theme directory", "dir", dirName, "err", err)
				continue
			}
			zlog.Debug("create theme directory", "dir", dirName)
			if err := extractThemeDir(dirName, filepath.Join(dir, file.Name())); err != nil {
				zlog.Warn("failed to extract theme directory", "dir", dirName, "err", err)
			}
			continue
		}
		filePath := filepath.Join(dir, file.Name())
		data, err := themes.DefaultAssets.ReadFile(filePath)
		if err != nil {
			zlog.Warn("failed to extract theme file", "name", file.Name(), "err", err)
			continue
		}
		dstFile := filepath.Join(topDir, file.Name())
		if err = utils.WriteFile(dstFile, data); err != nil {
			zlog.Warn("failed to write theme file", "name", file.Name(), "err", err)
			continue
		}
		zlog.Debug("create theme file", "name", dstFile)
	}
	return nil
}

var (
	//go:embed init_post.md
	initPostBytes []byte
	//go:embed init_page.md
	initPageBytes []byte
)

func createDemoPost(topDir string) error {
	post := &model.Post{
		Title:        "Hello World",
		Slug:         "hello-world",
		Descripition: "this is a demo post",
		DateString:   time.Now().Format("2006-01-02 15:04:05"),
		Tags:         []string{"hello"},
		Template:     "post.html",
		AuthorName:   "admin",
	}
	buf := bytes.NewBufferString("```toml\n")
	if err := toml.NewEncoder(buf).Encode(post); err != nil {
		return err
	}
	zlog.Debug("create demo post", "file", "content/posts/hello-world.md")
	buf.WriteString("```\n")
	buf.Write(initPostBytes)
	buf.WriteString("\n")
	return utils.WriteFile(filepath.Join(topDir, "content/posts/hello-world.md"), buf.Bytes())
}

func createDemoPage(topDir string) error {
	page := &model.Page{
		Post: model.Post{
			Title:        "About",
			Slug:         "about/",
			Descripition: "this is a demo page",
			DateString:   time.Now().Format("2006-01-02 15:04:05"),
			Template:     "page.html",
			AuthorName:   "admin",
		},
	}
	buf := bytes.NewBufferString("```toml\n")
	if err := toml.NewEncoder(buf).Encode(page); err != nil {
		return err
	}
	zlog.Debug("create demo page", "file", "content/pages/about.md")
	buf.WriteString("```\n")
	buf.Write(initPageBytes)
	buf.WriteString("\n")
	return utils.WriteFile(filepath.Join(topDir, "content/pages/about.md"), buf.Bytes())
}
