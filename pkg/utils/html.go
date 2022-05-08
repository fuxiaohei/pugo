package utils

import (
	"bytes"
	"path/filepath"
	"strings"
	"unicode"
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

var stripHTMLReplacer = strings.NewReplacer("\n", " ", "</p>", "\n", "<br>", "\n", "<br />", "\n")

// https://github.com/gohugoio/hugo/blob/release-0.98.0/helpers/content.go#L110, Apache License Version 2.0
// StripHTML accepts a string, strips out all HTML tags and returns it.
func StripHTML(s string) string {
	// Shortcut strings with no tags in them
	if !strings.ContainsAny(s, "<>") {
		return s
	}
	s = stripHTMLReplacer.Replace(s)

	// Walk through the string removing all tags
	b := bytes.NewBuffer(nil)
	var inTag, isSpace, wasSpace bool
	for _, r := range s {
		if !inTag {
			isSpace = false
		}

		switch {
		case r == '<':
			inTag = true
		case r == '>':
			inTag = false
		case unicode.IsSpace(r):
			isSpace = true
			fallthrough
		default:
			if !inTag && (!isSpace || (isSpace && !wasSpace)) {
				b.WriteRune(r)
			}
		}

		wasSpace = isSpace

	}
	return b.String()
}
