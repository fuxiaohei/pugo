package generator

import "pugo/pkg/core/constants"

// Options is the options for building a site.
type Option struct {
	ConfigFileItem *constants.ConfigFileItem
	OutputDir      string
	EnableWatch    bool // if true, watch source files and rebuild when changed
	EnableDrafts   bool // if true, render drafts
	IsLocalServer  bool // if true, some template should be ignored, such as googleAnalytics
	BuildArchive   bool // if true, build archive
}
