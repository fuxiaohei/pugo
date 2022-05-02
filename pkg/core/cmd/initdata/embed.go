package initdata

import (
	_ "embed"
)

var (
	//go:embed post.md
	PostBytes []byte
	//go:embed page.md
	PageBytes []byte
)
