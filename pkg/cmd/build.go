package cmd

import (
	"pugo/pkg/core/generator"

	"github.com/urfave/cli/v2"
)

var (
	genFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  "output",
			Usage: "set output directory, overwrite the config.toml value",
		},
		&cli.BoolFlag{
			Name:  "drafts",
			Usage: "build drafts",
		},
		&cli.BoolFlag{
			Name:  "watch",
			Usage: "watch source files and rebuild when changed",
		},
		&cli.BoolFlag{
			Name:  "archive",
			Usage: "compress built files to one archive",
		},
	}
)

// NewBuild returns a new cli.Command for the build subcommand.
func NewBuild() *cli.Command {
	cmd := &cli.Command{
		Name:        "build",
		Usage:       "build the site",
		Description: "build the site and generate the static files",
		Flags:       append(globalFlags, genFlags...),
		Aliases:     []string{"gen"},
		Action: func(c *cli.Context) error {

			initGlobalFlags(c)

			opt := parseCliOption(c)

			generator.Generate(opt)

			return nil
		},
	}
	return cmd
}

func parseCliOption(c *cli.Context) *generator.Option {
	configFileItem := loadLocalConfigFile()
	var option = generator.Option{
		ConfigFileItem: &configFileItem,
		EnableWatch:    c.Bool("watch"),
		OutputDir:      c.String("output"),
		BuildArchive:   c.Bool("archive"),
	}
	return &option
}
