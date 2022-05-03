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

// FullURL returns the full url after the base.
func FullURL(base, url string) string {
	return strings.TrimSuffix(base, "/") + "/" + strings.TrimPrefix(url, "/")
}
