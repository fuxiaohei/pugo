package builder

import (
	"haisite/internal/model"
	"haisite/internal/zlog"
	"os"
	"path/filepath"
)

var (
	// ContentPagesDirectory is the directory of pages.
	ContentPagesDirectory = "./content/pages"
)

func (b *Builder) buildPages() error {
	err := filepath.Walk(ContentPagesDirectory, func(path string, info os.FileInfo, err error) error {
		// skip directory
		if info.IsDir() {
			return nil
		}

		// only process markdown files
		if filepath.Ext(path) != ".md" {
			return nil
		}

		page, err := model.NewPageFromFile(path, ContentPagesDirectory)
		if err != nil {
			zlog.Warn("failed to build page", "path", path, "err", err)
			return nil
		}

		// save post into parsed data
		b.parsedData.Pages = append(b.parsedData.Pages, page)
		zlog.Info("pages: parsed ok", "path", path, "title", page.Title)

		return nil
	})

	if err != nil {
		return err
	}
	return nil
}
