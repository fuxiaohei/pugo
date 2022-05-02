package models

import (
	"os"
	"path/filepath"
	"pugo/pkg/core/constants"
	"pugo/pkg/zlog"
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
		p.Template = constants.PageTemplate
	}

	return &Page{
		Post: *p,
	}, nil
}

func LoadPages(s *SiteData) error {
	err := filepath.Walk(constants.ContentPagesDir, func(path string, info os.FileInfo, err error) error {
		// skip directory
		if info.IsDir() {
			return nil
		}

		// only process markdown files
		if filepath.Ext(path) != ".md" {
			return nil
		}

		page, err := NewPageFromFile(path, constants.ContentPagesDir)
		if err != nil {
			zlog.Warnf("failed to load page: %s, %s", path, err)
			return nil
		}

		// save post into parsed data
		s.Pages = append(s.Pages, page)
		zlog.Infof("load page ok: %s", path)

		return nil
	})

	if err != nil {
		return err
	}
	return nil
}
