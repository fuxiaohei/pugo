package builder

import (
	"bytes"
	"os"
	"path/filepath"
	"pugo/internal/model"
	"pugo/internal/zlog"
	"sort"
)

var (
	// ContentPostsDirectory is the directory of posts.
	ContentPostsDirectory = "./content/posts"
)

func (b *Builder) parsePosts() error {
	err := filepath.Walk(ContentPostsDirectory, func(path string, info os.FileInfo, err error) error {
		// skip directory
		if info.IsDir() {
			return nil
		}

		// only process markdown files
		if filepath.Ext(path) != ".md" {
			return nil
		}

		post, err := model.NewPostFromFile(path)
		if err != nil {
			zlog.Warn("failed to build post", "path", path, "err", err)
			return nil
		}

		// save post into parsed data
		b.source.Posts = append(b.source.Posts, post)
		zlog.Info("posts: parsed ok", "path", path, "title", post.Title)

		return nil
	})

	if err != nil {
		return err
	}
	sort.Slice(b.source.Posts, func(i, j int) bool {
		// order by date desc
		return b.source.Posts[i].Date().Unix() > b.source.Posts[j].Date().Unix()
	})
	return nil
}

func (b *Builder) buildPosts(ctx *buildContext) error {
	var (
		err        error
		dstFile    string
		link       string
		buf        *bytes.Buffer
		tplData    map[string]interface{}
		descGetter = func(post *model.Post) string {
			if post.Descripition != "" {
				return post.Descripition
			}
			return b.source.Config.Site.Description
		}
	)

	// build each post
	for _, p := range b.source.Posts {

		// build link
		link, dstFile, err = ctx.buildPostLink(p)
		if err != nil {
			zlog.Warn("failed to build post slug dstFile", "title", p.Title, "path", p.LocalFile(), "slug", p.Slug, "err", err)
			continue
		}
		p.Link = link

		// convert markdown
		if err := p.Convert(b.markdown); err != nil {
			zlog.Warn("failed to convert markdown post", "title", p.Title, "path", p.LocalFile(), "err", err)
			continue
		}

		buf = bytes.NewBuffer(nil)
		extData := map[string]interface{}{
			"post": p,
			"current": map[string]interface{}{
				"Title":       p.Title + " - " + b.source.Config.Site.Title,
				"Description": descGetter(p),
			},
		}
		tplData = ctx.buildTemplateData(extData)
		if err = b.render.Execute(buf, p.Template, tplData); err != nil {
			zlog.Warn("failed to render post", "title", p.Title, "path", p.LocalFile(), "err", err)
			continue
		}

		// save buffer to write content file later
		ctx.setBuffer(dstFile, buf)

		zlog.Info("posts: rendered ok", "title", p.Title, "path", p.LocalFile(), "dst", dstFile, "size", buf.Len())
	}

	// create pager
	b.source.PostsPager = model.NewPager(b.source.Config.BuildConfig.PostPerPage, len(b.source.Posts))
	zlog.Info("posts: pager created ok", "total", len(b.source.Posts), "per_page", b.source.Config.BuildConfig.PostPerPage, "pages", b.source.PostsPager.PageSize())
	return nil
}

func (b *Builder) buildPostLists(ctx *buildContext) error {
	total := b.source.PostsPager.PageSize()
	tplName := model.DefaultPostListTemplate
	for i := 1; i <= total; i++ {
		// build each page list
		buf := bytes.NewBuffer(nil)
		tplData, pageItem := b.buildPostListTemplateData(ctx, i)
		if err := b.render.Execute(buf, tplName, tplData); err != nil {
			zlog.Warn("failed to render post list", "page", i, "err", err)
			return err
		}
		dstFile := pageItem.LocalFile
		ctx.setBuffer(dstFile, buf)
		zlog.Info("posts: rendered page list ok", "page", i, "dst", dstFile, "size", buf.Len())
	}
	return nil
}

func (b *Builder) buildPostListTemplateData(ctx *buildContext, page int) (map[string]interface{}, *model.PagerItem) {
	pageItem := b.source.PostsPager.Page(page, b.source.Config.BuildConfig.PostPageLinkFormat)
	tplData := ctx.buildTemplateData(map[string]interface{}{
		"posts": model.PostsPageList(b.source.Posts, pageItem),
		"pager": pageItem,
		"current": map[string]interface{}{
			"Title": b.source.Config.Site.Title,
		},
	})
	return tplData, pageItem
}
