package generator

import "pugo/pkg/constants"

// Options is the options for building a site.
type Option struct {
	ConfigFileItem *constants.ConfigFileItem
	OutputDir      string
	EnableWatch    bool
}
