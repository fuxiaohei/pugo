package cmd

import (
	"bytes"
	"fmt"
	"path/filepath"
	"pugo/pkg/core/cmd/initdata"
	"pugo/pkg/core/constants"
	"pugo/pkg/core/models"
	"pugo/pkg/utils"
	"pugo/pkg/utils/zlog"
	"pugo/themes"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/urfave/cli/v2"
)

var (
	initFlags = []cli.Flag{
		&cli.BoolFlag{
			Name:    "yaml",
			Value:   false,
			Aliases: []string{"yml"},
			Usage:   "initialize a new site with yaml config file",
		},
	}
)

// NewInit returns a new cli.Command for the init subcommand.
func NewInit() *cli.Command {
	cmd := &cli.Command{
		Name:  "init",
		Usage: "initialize a new sample site in the current directory",
		Flags: append(globalFlags, initFlags...),
		Action: func(c *cli.Context) error {
			initGlobalFlags(c)

			ctype := constants.ConfigTypeTOML
			if c.Bool("yaml") {
				ctype = constants.ConfigTypeYAML
			}
			configItem, err := initConfigFile(ctype)

			if err != nil {
				zlog.Warnf("failed to initialize config file: %s", err)
				return err
			}
			zlog.Debugf("initialized config file: %s ", configItem.File)
			if err := initDirectories(); err != nil {
				zlog.Warnf("failed to initialize directories: %s", err)
				return err
			}
			if err := initTheme(); err != nil {
				zlog.Warnf("failed to initialize theme: %s", err)
				return err
			}
			if err := initFirstPost(); err != nil {
				zlog.Warnf("failed to initialize first post: %s", err)
				return err
			}
			if err := initFirstPage(); err != nil {
				zlog.Warnf("failed to initialize first page: %s", err)
				return err
			}
			zlog.Infof("initialized sample site in the current directory successfully")
			return nil
		},
	}
	return cmd
}

func initConfigFile(ctype constants.ConfigType) (*constants.ConfigFileItem, error) {
	configFileItem := selectConfigFile(ctype)
	if utils.IsFileExist(configFileItem.File) {
		// FIXME: should we overwrite the config file?
	}
	if configFileItem.Type == constants.ConfigTypeTOML {
		return configFileItem, utils.WriteTOMLFile(configFileItem.File, models.DefaultConfig())
	}
	if configFileItem.Type == constants.ConfigTypeYAML {
		return configFileItem, utils.WriteYAMLFile(configFileItem.File, models.DefaultConfig())
	}
	return nil, fmt.Errorf("unsupported config type: %s", configFileItem.Type)
}

func initDirectories() error {
	for _, dir := range constants.InitDirectories() {
		if err := utils.MkdirAll(dir); err != nil {
			zlog.Warnf("failed to create directory: '%s' ,%s", dir, err)
			return err
		}
		zlog.Debugf("created directory: '%s'", dir)
	}
	return nil
}

func initTheme() error {
	// extract default theme
	if err := extractThemeDir("themes/default", "default"); err != nil {
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
				zlog.Warnf("failed to create theme directory: '%s', %s", dirName, err)
				continue
			}
			zlog.Debugf("created theme directory: '%s'", dirName)
			extractThemeDir(dirName, filepath.Join(dir, file.Name()))
			continue
		}
		filePath := filepath.Join(dir, file.Name())
		data, err := themes.DefaultAssets.ReadFile(filePath)
		if err != nil {
			zlog.Warnf("failed to extract theme file: '%s', %s", filePath, err)
			continue
		}
		dstFile := filepath.Join(topDir, file.Name())
		if err = utils.WriteFile(dstFile, data); err != nil {
			zlog.Warnf("failed to write theme file: '%s', %s", dstFile, err)
			continue
		}
		zlog.Debugf("extracted theme file: '%s'", dstFile)
	}
	return nil
}

func initFirstPost() error {
	post := &models.Post{
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
	zlog.Debugf("create first post: %s", "content/posts/hello-world.md")
	buf.WriteString("```\n")
	buf.Write(initdata.PostBytes)
	buf.WriteString("\n")
	return utils.WriteFile("content/posts/hello-world.md", buf.Bytes())
}

func initFirstPage() error {
	page := &models.Page{
		Post: models.Post{
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
	zlog.Debugf("create fist page: %s", "content/pages/about.md")
	buf.WriteString("```\n")
	buf.Write(initdata.PageBytes)
	buf.WriteString("\n")
	return utils.WriteFile("content/pages/about.md", buf.Bytes())
}
