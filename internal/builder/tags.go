package builder

import (
	"bytes"
	"pugo/internal/model"
	"pugo/internal/zlog"
	"strings"
)

func (b *Builder) buildTags(ctx *buildContext) error {

	tplName := model.DefaultPostListTemplate

	// build tag pages
	for _, tagData := range b.source.Tags {
		pager := model.NewPager(b.source.BuildConfig.PostPerPage, len(tagData.Posts))

		total := pager.PageSize()
		for i := 1; i <= total; i++ {
			linkFormat := strings.ReplaceAll(b.source.BuildConfig.TagPageLinkFormat, "{{.Tag}}", tagData.Tag.Name)
			pageItem := pager.Page(i, linkFormat)
			dstFile := pageItem.LocalFile

			buf := bytes.NewBuffer(nil)
			posts := model.PostsPageList(tagData.Posts, pageItem)
			tplData := ctx.buildTemplateData(map[string]interface{}{
				"posts": posts,
				"pager": pageItem,
				"tag":   tagData.Tag,
				"current": map[string]interface{}{
					"Title": tagData.Tag.Name + "-" + b.source.Config.Site.Title,
				},
			})

			if err := b.render.Execute(buf, tplName, tplData); err != nil {
				zlog.Warn("failed to render tag page", "page", i, "tag", tagData.Tag.Name, "err", err)
				return err
			}
			ctx.setBuffer(dstFile, buf)
			zlog.Info("tags: rendered ok", "page", i, "tag", tagData.Tag.Name, "dst", dstFile, "size", buf.Len())

			t := posts[0].Date()
			ctx.addSitemap(&model.SitemapURL{Loc: pageItem.Link, LastMod: &t})

			// tag list index.html
			if i == 1 {
				ctx.setBuffer(tagData.Tag.LocalFile, buf)
				ctx.addSitemap(&model.SitemapURL{Loc: tagData.Tag.Link, LastMod: &t})
				zlog.Info("tags: rendered index ok", "page", i, "tag", tagData.Tag.Name, "dst", tagData.Tag.LocalFile, "size", buf.Len())
			}
		}
	}

	zlog.Info("posts: build tags", "tags", len(b.source.Tags))
	return nil
}
