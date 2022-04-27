package generator

import (
	"bytes"
	"path/filepath"
	"pugo/pkg/constants"
	"pugo/pkg/models"
	"pugo/pkg/theme"
	"pugo/pkg/utils"
	"pugo/pkg/zlog"
)

type renderArchivesParams struct {
	Ctx          *Context
	Posts        []*models.Post
	SiteTitle    string
	Render       *theme.Render
	ArchivesLink string
	OutputDir    string
}

func renderArchives(params *renderArchivesParams) error {
	archives := models.NewArchives(params.Posts)
	buf := bytes.NewBuffer(nil)
	extData := map[string]interface{}{
		"archives": archives,
		"current": map[string]interface{}{
			"Title": params.SiteTitle,
		},
	}
	tplData := params.Ctx.createTemplateData(extData)
	if err := params.Render.Execute(buf, constants.ArchivesTemplate, tplData); err != nil {
		zlog.Warnf("failed to render archives: %s", err)
		return err
	}
	dstFile := utils.FormatIndexHTML(params.ArchivesLink)
	dstFile = filepath.Join(params.OutputDir, dstFile)
	params.Ctx.SetOutput(dstFile, buf)
	zlog.Infof("archives generated: %s", dstFile)

	t := params.Posts[0].Date()
	params.Ctx.addSitemap(&models.SitemapURL{Loc: params.ArchivesLink, LastMod: &t})
	return nil
}
