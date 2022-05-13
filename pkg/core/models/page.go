package models

import (
	"os"
	"path/filepath"
	"pugo/pkg/core/constants"
	"pugo/pkg/utils/zlog"
	"strings"
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
		ext := filepath.Ext(path)
		p.Slug, _ = filepath.Rel(contentDir, path)
		p.Slug = strings.TrimSuffix(p.Slug, ext)

		// trim index.html
		if strings.HasSuffix(p.Slug, "/index") {
			p.Slug = strings.TrimSuffix(p.Slug, "/index")
		}
	}

	// fix empty template
	if p.Template == "" {
		p.Template = constants.PageTemplate
	}

	return &Page{
		Post: *p,
	}, nil
}

func LoadPages(withDrafts bool) ([]*Page, error) {
	var pages []*Page
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
		if page.Draft && !withDrafts {
			zlog.Warnf("skip draft page: %s", path)
			return nil
		}

		// save post into parsed data
		pages = append(pages, page)
		zlog.Infof("load page ok: %s", path)

		return nil
	})

	if err != nil {
		return nil, err
	}
	return pages, nil
}
