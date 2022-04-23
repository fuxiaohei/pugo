package builder

import (
	"bytes"
	"os"
	"path/filepath"
	"pugo/internal/model"
	"pugo/internal/zlog"
	"strings"
)

var (
	// ContentPagesDirectory is the directory of pages.
	ContentPagesDirectory = "./content/pages"
)

func (b *Builder) parsePages() error {
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
		b.source.Pages = append(b.source.Pages, page)
		zlog.Info("pages: parsed ok", "path", path, "title", page.Title)

		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

func (b *Builder) buildPages(ctx *buildContext) error {
	var (
		err        error
		dstFile    string
		buf        *bytes.Buffer
		tplData    map[string]interface{}
		descGetter = func(page *model.Page) string {
			if page.Descripition != "" {
				return page.Descripition
			}
			return b.source.Config.Site.Description
		}
	)
	// build each page
	for _, pg := range b.source.Pages {
		pg.Link = "/" + strings.TrimPrefix(pg.Slug, "/")
		dstFile = model.FormatIndexHTML(pg.Link)

		// convert markdown to html
		if err = pg.Convert(b.markdown); err != nil {
			zlog.Warn("failed to convert markdown post", "title", pg.Title, "path", pg.LocalFile(), "err", err)
			continue
		}

		buf = bytes.NewBuffer(nil)
		extData := map[string]interface{}{
			"page": pg,
			"current": map[string]interface{}{
				"Title":       pg.Title + " - " + b.source.Config.Site.Title,
				"Description": descGetter(pg),
			},
		}
		tplData = ctx.buildTemplateData(extData)
		if err = b.render.Execute(buf, pg.Template, tplData); err != nil {
			zlog.Warn("failed to render page", "title", pg.Title, "path", pg.LocalFile(), "err", err)
			continue
		}

		ctx.setBuffer(dstFile, buf)
		zlog.Info("pages: rendered ok", "title", pg.Title, "path", pg.LocalFile(), "dst", dstFile, "size", buf.Len())
	}
	return nil
}
