package model

import (
	"path/filepath"
)

const (
	// DefaultPageTemplate is the default page template.
	DefaultPageTemplate = "page.html"
)

// Page is the page model.
type Page struct {
	Post
}

// NewPageFromFile creates a new page from file
func NewPageFromFile(path, contentDir string) (*Page, error) {
	// parse basic info as post
	p, err := parseContentBase(path)
	if err != nil {
		return nil, err
	}

	// fix slug empty
	if p.Slug == "" {
		p.Slug, _ = filepath.Rel(contentDir, path)
	}

	// fix empty template
	if p.Template == "" {
		p.Template = DefaultPageTemplate
	}

	return &Page{
		Post: *p,
	}, nil
}
