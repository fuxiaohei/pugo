package generator

import (
	"bytes"
	"path/filepath"
	"pugo/pkg/core/constants"
	"pugo/pkg/core/models"
	"pugo/pkg/utils"
	"pugo/pkg/zlog"
)

type renderArchivesParams struct {
	renderBaseParams
	Posts        []*models.Post
	SiteTitle    string
	ArchivesLink string
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
