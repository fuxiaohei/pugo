package generator

import (
	"bytes"
	"path/filepath"
	"pugo/pkg/utils/zlog"
)

func renderIndex(params *renderPostListsParams) error {
	indexTpl := params.Render.GetIndexTemplate()

	// first page
	tplData, _ := buildPostListTemplateData(params, 1)

	buf := bytes.NewBuffer(nil)
	if err := params.Render.Execute(buf, indexTpl, tplData); err != nil {
		zlog.Warn("failed to render index", "err", err)
		return err
	}
	link := "/index.html"
	dstFile := filepath.Join(params.OutputDir, link)
	params.Ctx.SetOutput(dstFile, link, buf)
	zlog.Infof("index generated: %s", dstFile)
	return nil
}
