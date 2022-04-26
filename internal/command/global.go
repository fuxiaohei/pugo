package command

import (
	"pugo/internal/zlog"

	"github.com/urfave/cli/v2"
)

var (
	// GlobalFlags is the list of global flags.
	GlobalFlags = []cli.Flag{
		&cli.BoolFlag{
			Name:  "debug",
			Value: false,
			Usage: "enable debug mode",
		},
	}

	contentPostsDir = "content/posts"
	contentPagesDir = "content/pages"
)

func initGlobalFlags(c *cli.Context) {
	// set debug mode
	if c.Bool("debug") {
		zlog.Init(true)
		zlog.Debug("debug mode enabled")
	} else {
		zlog.Init(false)
	}
}
