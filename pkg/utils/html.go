package utils

import (
	"path/filepath"
	"strings"
)

// FormatIndexHTML append index.html to the end of the post.
func FormatIndexHTML(slug string) string {
	if strings.HasSuffix(slug, ".html") {
		return slug
	}
	return filepath.Join(slug, "index.html")
}
