package command

import (
	"haisite/internal/builder"

	"github.com/urfave/cli/v2"
)

var (
	buildFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  "output",
			Usage: "set output directory, overwrite the config.toml value",
		},
		&cli.BoolFlag{
			Name:  "drafts",
			Usage: "build drafts",
		},
	}
)

// NewBuildCommand returns a new cli.Command for the build subcommand.
func NewBuildCommand() *cli.Command {
	cmd := &cli.Command{
		Name:        "build",
		Usage:       "build the site",
		Description: "build the site and generate the static files",
		Flags:       append(GlobalFlags, buildFlags...),
		Action: func(c *cli.Context) error {

			initGlobalFlags(c)

			execBuilder(c)

			return nil
		},
	}
	return cmd
}

func execBuilder(c *cli.Context) *builder.Builder {
	var option builder.Option
	option.ConfigFile = DefaultConfigFile

	// parse output directory
	if c.String("output") != "" {
		option.OutputDir = c.String("output")
	}

	b := builder.NewBuilder(&option)
	b.Build()
	return b
}
