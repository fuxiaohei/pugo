package cmd

import (
	"pugo/pkg/zlog"

	"github.com/urfave/cli/v2"
)

var (
	globalFlags = []cli.Flag{
		&cli.BoolFlag{
			Name:  "debug",
			Value: false,
			Usage: "enable debug mode",
		},
	}
)

// GetGlobalFlags returns the global flags.
func GetGlobalFlags() []cli.Flag {
	return globalFlags
}

func initGlobalFlags(c *cli.Context) {
	// set debug mode
	if c.Bool("debug") {
		zlog.Init(true)
		zlog.Debug("debug logging enabled")
	} else {
		zlog.Init(false)
	}
}
