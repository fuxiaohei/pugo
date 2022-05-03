package generator

import (
	"bytes"
	"path/filepath"
	"pugo/pkg/ext/sitemap"
	"pugo/pkg/utils/zlog"
)

func renderSitemap(link, outputDir string, ctx *Context) error {
	buf := bytes.NewBuffer(nil)
	if err := sitemap.Write(buf); err != nil {
		zlog.Warnf("failed to marshal sitemap: %s", err)
		return err
	}
	var dstFile = link
	dstFile = filepath.Join(outputDir, dstFile)
	ctx.SetOutput(dstFile, buf)
	zlog.Infof("sitemap generated: %s", dstFile)
	return nil
}
