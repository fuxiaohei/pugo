package generator

import (
	"bytes"
	"path/filepath"
	"pugo/pkg/core/models"
	"pugo/pkg/ext/markdown"
	"pugo/pkg/utils"
	"pugo/pkg/utils/zlog"
	"strings"
)

type renderPagesParams struct {
	renderBaseParams
	Pages     []*models.Page
	SiteDesc  string
	SiteTitle string
}

func renderPages(params *renderPagesParams) error {
	var (
		err        error
		dstFile    string
		buf        *bytes.Buffer
		tplData    map[string]interface{}
		descGetter = func(page *models.Page) string {
			if page.Descripition != "" {
				return page.Descripition
			}
			return params.SiteDesc
		}
	)
	// build each page
	for _, pg := range params.Pages {
		pg.Link = "/" + strings.TrimPrefix(pg.Slug, "/")
		dstFile = utils.FormatIndexHTML(pg.Link)

		// convert markdown to html
		if err = pg.Convert(markdown.Get()); err != nil {
			zlog.Warnf("failed to convert markdown page: %s, %s", pg.LocalFile(), err)
			continue
		}

		buf = bytes.NewBuffer(nil)
		extData := map[string]interface{}{
			"page": pg,
			"current": map[string]interface{}{
				"Title":       pg.Title + " - " + params.SiteTitle,
				"Description": descGetter(pg),
			},
		}
		tplData = params.Ctx.createTemplateData(extData)
		if err = params.Render.Execute(buf, pg.Template, tplData); err != nil {
			zlog.Warnf("failed to render page: %s, %s", pg.LocalFile(), err)
			continue
		}
		dstFile = filepath.Join(params.OutputDir, dstFile)
		params.Ctx.SetOutput(dstFile, buf)
		zlog.Infof("page generated: %s", dstFile)

		t := pg.Date()
		params.Ctx.addSitemap(&models.SitemapURL{Loc: pg.Link, LastMod: &t})
	}

	return nil
}
