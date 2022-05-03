package cmd

import (
	"bytes"
	"errors"
	"path/filepath"
	"pugo/pkg/core/configs"
	"pugo/pkg/core/constants"
	"pugo/pkg/core/models"
	"pugo/pkg/utils"
	"pugo/pkg/utils/zlog"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

// NewCreate returns a new cli.Command for the create subcommand.
func NewCreate() *cli.Command {
	cmd := &cli.Command{
		Name:        "create",
		Usage:       "create new post or page",
		Description: "create new post or page with default template",
		Flags:       append(globalFlags, genFlags...),
		Action: func(c *cli.Context) error {

			initGlobalFlags(c)

			if c.Args().Len() == 0 {
				return cli.Exit("create command requires 'post' or 'page' as first param", 1)
			}
			filename := c.Args().Get(1)
			if filename == "" {
				return cli.Exit("create command requires 'filename' as second param", 1)
			}

			typename := strings.ToLower(c.Args().Get(0))

			if typename != "post" && typename != "page" {
				return cli.Exit("create command requires 'post' or 'page' param", 1)
			}

			// load config
			configFileItem := loadLocalConfigFile()
			config, err := configs.LoadFromFile(configFileItem)
			if err != nil {
				zlog.Warnf("load config file failed: %v", err)
				return err
			}
			zlog.Debugf("load config file: %s", configFileItem.File)

			if typename == "post" {
				return createSamplePost(filename, config, configFileItem)
			}
			if typename == "page" {
				return createSamplePage(filename, config, configFileItem)
			}
			return nil

		},
	}
	return cmd
}

func createContentMeta(content interface{}, configType constants.ConfigType) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer(nil)
	for _, separator := range constants.PostMetaSeperators() {
		if separator.MetaType == configType {
			buf.Write(separator.StartChars)

			if configType == constants.ConfigTypeTOML {
				if err := toml.NewEncoder(buf).Encode(content); err != nil {
					return nil, err
				}
			} else if configType == constants.ConfigTypeYAML {
				data, err := yaml.Marshal(content)
				if err != nil {
					return nil, err
				}
				buf.Write(data)
			}

			buf.Write(separator.EndChars)
			break
		}
	}
	return buf, nil
}

func createSamplePost(filename string, cfg *configs.Config, item constants.ConfigFileItem) error {
	fpath := filepath.Join(constants.ContentPostsDir, filename)
	if filepath.Ext(fpath) != ".md" {
		return errors.New("file extension must be .md")
	}
	basename := filepath.Base(fpath)
	slug := strings.TrimSuffix(basename, filepath.Ext(basename))

	post := &models.Post{
		Title:        normalizeTitle(slug),
		Slug:         slug,
		Descripition: "",
		DateString:   time.Now().Format("2006-01-02 15:04:05"),
		Tags:         []string{"post"},
		Template:     "post.html",
		AuthorName:   cfg.Author[0].Name,
	}
	buf, err := createContentMeta(post, item.Type)
	if err != nil {
		return err
	}
	buf.WriteString("this is an empty post\n")
	zlog.Info("create post", "filename", fpath)
	return utils.WriteFile(fpath, buf.Bytes())
}

func createSamplePage(filename string, cfg *configs.Config, item constants.ConfigFileItem) error {
	fpath := filepath.Join(constants.ContentPagesDir, filename)
	if filepath.Ext(fpath) != ".md" {
		return errors.New("file extension must be .md")
	}
	basename := filepath.Base(fpath)
	slug := strings.TrimSuffix(basename, filepath.Ext(basename))

	page := &models.Page{
		Post: models.Post{
			Title:        normalizeTitle(slug),
			Slug:         slug,
			Descripition: "this is an empty page",
			DateString:   time.Now().Format("2006-01-02 15:04:05"),
			Template:     "page.html",
			AuthorName:   cfg.Author[0].Name,
		},
	}
	buf, err := createContentMeta(page, item.Type)
	if err != nil {
		return err
	}
	buf.WriteString("this is an empty page\n")
	zlog.Info("create page", "filename", fpath)
	return utils.WriteFile(fpath, buf.Bytes())
}

var replacer = strings.NewReplacer("-", " ", "_", " ")

func normalizeTitle(name string) string {
	name = replacer.Replace(name)
	return strings.Title(strings.ToLower(name))
}
