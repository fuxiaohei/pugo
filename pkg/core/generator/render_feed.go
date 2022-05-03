package generator

import (
	"bytes"
	"path/filepath"
	"pugo/pkg/utils/zlog"
)

func renderSitemap(outputDir string, ctx *Context) error {
	sitemap := ctx.getSitemap()
	buf := bytes.NewBuffer(nil)
	if err := sitemap.Write(buf); err != nil {
		zlog.Warnf("failed to marshal sitemap: %s", err)
		return err
	}
	var dstFile = "/sitemap.xml"
	dstFile = filepath.Join(outputDir, dstFile)
	ctx.SetOutput(dstFile, buf)
	zlog.Infof("sitemap generated: %s", dstFile)
	return nil
}
