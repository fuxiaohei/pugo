package builder

import (
	"bytes"
	"encoding/xml"
	"pugo/internal/model"
	"pugo/internal/zlog"
)

func (b *Builder) buildIndex(ctx *buildContext) error {
	indexTpl := b.render.GetIndexTemplate()

	// first page
	tplData, _ := b.buildPostListTemplateData(ctx, 1)

	buf := bytes.NewBuffer(nil)
	if err := b.render.Execute(buf, indexTpl, tplData); err != nil {
		zlog.Warn("failed to render index", "err", err)
		return err
	}
	ctx.setBuffer("/index.html", buf)
	zlog.Info("posts: index rendered ok", "size", buf.Len())
	return nil
}

func (b *Builder) buildFeedAtom(ctx *buildContext) error {
	var posts []*model.Post
	var limit = b.source.BuildConfig.FeedPostLimit
	if limit <= 0 {
		limit = model.DefaultFeedPostLimit
	}
	if limit > len(b.source.Posts) {
		posts = b.source.Posts
	} else {
		posts = b.source.Posts[:limit]
	}
	var dstFile = "/atom.xml"
	feed := model.BuildAtom(dstFile, posts, b.source.Config.Site)
	data, err := xml.Marshal(feed)
	if err != nil {
		zlog.Warn("failed to marshal atom feed", "err", err)
		return err
	}
	ctx.setBuffer(dstFile, bytes.NewBuffer(data))
	zlog.Info("posts: feed atom rendered ok", "posts", len(posts), "limit", limit, "size", len(data))
	return nil
}

func (b *Builder) buildSitemap(ctx *buildContext) error {
	sitemap := ctx.getSitemap()
	buf := bytes.NewBuffer(nil)
	if err := sitemap.Write(buf); err != nil {
		zlog.Warn("failed to write sitemap", "err", err)
		return err
	}
	var dstFile = "/sitemap.xml"
	ctx.setBuffer(dstFile, buf)
	zlog.Info("sitemap: rendered ok", "size", buf.Len())
	return nil
}
