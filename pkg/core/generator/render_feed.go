package generator

import (
	"bytes"
	"encoding/xml"
	"path/filepath"
	"pugo/pkg/core/constants"
	"pugo/pkg/core/models"
	"pugo/pkg/utils/zlog"
)

type renderFeedAtomParams struct {
	renderBaseParams
	Posts         []*models.Post
	SiteConfig    *models.SiteConfig
	PostFeedLimit int
}

func renderFeedAtom(params *renderFeedAtomParams) error {
	var posts []*models.Post
	var limit = params.PostFeedLimit
	if limit <= 0 {
		limit = constants.PostFeedLimit
	}
	if limit > len(params.Posts) {
		posts = params.Posts
	} else {
		posts = params.Posts[:limit]
	}
	var dstFile = "/atom.xml"
	feed := models.BuildAtom(dstFile, posts, params.SiteConfig)
	data, err := xml.Marshal(feed)
	if err != nil {
		zlog.Warnf("failed to marshal atom feed: %s", err)
		return err
	}
	dstFile = filepath.Join(params.OutputDir, dstFile)
	params.Ctx.SetOutput(dstFile, bytes.NewBuffer(data))
	zlog.Infof("atom feed generated: %s", dstFile)
	return nil
}

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
