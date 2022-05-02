package generator

import "pugo/pkg/core/constants"

// Options is the options for building a site.
type Option struct {
	ConfigFileItem *constants.ConfigFileItem
	OutputDir      string
	EnableWatch    bool
	IsLocalServer  bool // if true, some template should be ignored, such as googleAnalytics
}
