package generator

import (
	"bytes"
	"path/filepath"
	"pugo/pkg/utils/zlog"
)

type renderErrorPageParams struct {
	renderBaseParams
	SiteTitle string
}

func renderErrorPage(params *renderErrorPageParams) error {
	notFoundTpl := params.Render.GetTemplate("404")
	tplData := params.Ctx.createTemplateData(map[string]interface{}{
		"current": map[string]interface{}{
			"Title": params.SiteTitle,
		},
	})

	buf := bytes.NewBuffer(nil)
	if err := params.Render.Execute(buf, notFoundTpl, tplData); err != nil {
		zlog.Warn("failed to render 404", "err", err)
		return err
	}
	link := "/404.html"
	dstFile := filepath.Join(params.OutputDir, link)
	params.Ctx.SetOutput(dstFile, link, buf)
	zlog.Infof("404 generated: %s", dstFile)
	return nil
}
