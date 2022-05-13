package generator

import (
	"bytes"
	"path/filepath"
	"pugo/pkg/core/constants"
	"pugo/pkg/core/models"
	"pugo/pkg/ext/sitemap"
	"pugo/pkg/utils"
	"pugo/pkg/utils/zlog"
)

type renderArchivesParams struct {
	renderBaseParams
	Posts        []*models.Post
	ArchivesLink string
}

func renderArchives(params *renderArchivesParams) error {
	archives := models.NewArchives(params.Posts)
	buf := bytes.NewBuffer(nil)
	extData := map[string]interface{}{
		"archives": archives,
		"current": map[string]interface{}{
			"Title":       params.SiteTitle,
			"SubTitle":    params.SiteSubTitle,
			"Description": params.SiteDescription,
			"Link":        params.ArchivesLink,
			"Slug":        params.ArchivesLink,
		},
		"i18n": params.I18n.Get(""),
	}
	tplData := params.Ctx.createTemplateData(extData)
	if err := params.Render.Execute(buf, constants.ArchivesTemplate, tplData); err != nil {
		zlog.Warnf("failed to render archives: %s", err)
		return err
	}
	dstFile := utils.FormatIndexHTML(params.ArchivesLink)
	dstFile = filepath.Join(params.OutputDir, dstFile)
	params.Ctx.SetOutput(dstFile, params.ArchivesLink, buf)
	zlog.Infof("archives generated: %s", dstFile)

	t := params.Posts[0].Date()

	sitemap.Add(&sitemap.URL{Loc: params.ArchivesLink, LastMod: &t})

	return nil
}
