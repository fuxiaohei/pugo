package cmd

import (
	"pugo/pkg/core/constants"
	"pugo/pkg/utils"
	"pugo/pkg/utils/zlog"

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
	zlog.Infof("%s %s", c.App.Name, c.App.Version)
	// set debug mode
	if c.Bool("debug") {
		zlog.Init(true)
		zlog.Debug("debug logging enabled")
	} else {
		zlog.Init(false)
	}
}

func loadLocalConfigFile() constants.ConfigFileItem {
	items := constants.ConfigFiles()
	for _, item := range items {
		if utils.IsFileExist(item.File) {
			return item
		}
	}
	return items[0]
}

func selectConfigFile(ctype constants.ConfigType) *constants.ConfigFileItem {
	items := constants.ConfigFiles()
	for _, item := range items {
		if item.Type == ctype {
			return &item
		}
	}
	return selectConfigFile(constants.ConfigTypeTOML)
}
