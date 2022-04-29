package cmd

import (
	"pugo/pkg/generator"

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
	}

	// parse output directory
	if c.String("output") != "" {
		option.OutputDir = c.String("output")
	}

	// enable watching source files
	if c.Bool("watch") {
		option.EnableWatch = true

	}

	return &option
}
