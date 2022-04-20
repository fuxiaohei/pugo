package command

import (
	"pugo/internal/server"

	"github.com/urfave/cli/v2"
)

var (
	serveFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  "port",
			Value: "18080",
			Usage: "port to listen",
		},
	}
)

// NewServeCommand returns a new cli.Command for the serve subcommand.
func NewServeCommand() *cli.Command {
	flags := append(GlobalFlags, buildFlags...)
	flags = append(flags, serveFlags...)
	cmd := &cli.Command{
		Name:        "serve",
		Usage:       "serve the site",
		Description: "start serve the site, auto rebuild when source files changed",
		Flags:       flags,
		Action: func(c *cli.Context) error {
			initGlobalFlags(c)

			builder := execBuilder(c)

			s := server.New(builder.OutputDir(), c.Int("port"))
			s.Run()
			return nil
		},
	}
	return cmd
}
