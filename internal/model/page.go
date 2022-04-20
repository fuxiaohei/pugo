package model

import "path/filepath"

// Page is the page model.
type Page struct {
	Post
}

// NewPageFromFile creates a new page from file
func NewPageFromFile(path string, contentDir string) (*Page, error) {
	// parse basic info as post
	p, err := NewPostFromFile(path)
	if err != nil {
		return nil, err
	}

	// fix slug empty
	if p.Slug == "" {
		p.Slug, _ = filepath.Rel(contentDir, path)
	}

	return &Page{
		Post: *p,
	}, nil
}
