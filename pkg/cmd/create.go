package cmd

import (
	"bytes"
	"errors"
	"path/filepath"
	"pugo/pkg/constants"
	"pugo/pkg/models"
	"pugo/pkg/utils"
	"pugo/pkg/zlog"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/urfave/cli/v2"
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
			config, err := models.LoadConfigFromFile(constants.ConfigFile)
			if err != nil {
				zlog.Warnf("load config file failed: %v", err)
				return err
			}

			if typename == "post" {
				return createSamplePost(filename, config)
			}
			if typename == "page" {
				return createSamplePage(filename, config)
			}
			return nil

		},
	}
	return cmd
}

func createSamplePost(filename string, cfg *models.Config) error {
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

	buf := bytes.NewBufferString("```toml\n")
	if err := toml.NewEncoder(buf).Encode(post); err != nil {
		return err
	}
	buf.WriteString("```\n")
	buf.WriteString("this is an empty post\n")
	zlog.Info("create post", "filename", fpath)
	return utils.WriteFile(fpath, buf.Bytes())
}

func createSamplePage(filename string, cfg *models.Config) error {
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
	buf := bytes.NewBufferString("```toml\n")
	if err := toml.NewEncoder(buf).Encode(page); err != nil {
		return err
	}
	buf.WriteString("```\n")
	buf.WriteString("this is an empty page\n")
	zlog.Info("create page", "filename", fpath)
	return utils.WriteFile(fpath, buf.Bytes())
}

var replacer = strings.NewReplacer("-", " ", "_", " ")

func normalizeTitle(name string) string {
	name = replacer.Replace(name)
	return strings.Title(strings.ToLower(name))
}
