package constants

import (
	"fmt"
	"strings"
)

var (
	// ErrInvalidPostStartLine means the post start line is invalid.
	ErrInvalidPostStartLine = fmt.Errorf("invalid content start")
	// ErrInvalidPostDate means the post date is invalid, it must be format as postDefaultDateLayout
	ErrInvalidPostDate = fmt.Errorf("invalid content date, it must be format as " + strings.Join(postDateLayouts, " or "))
)
