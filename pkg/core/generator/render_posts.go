package generator

import (
	"bytes"
	"path/filepath"
	"pugo/pkg/core/constants"
	"pugo/pkg/core/models"
	"pugo/pkg/ext/markdown"
	"pugo/pkg/ext/sitemap"
	"pugo/pkg/utils/zlog"
)

type renderPostsParams struct {
	renderBaseParams
	Posts     []*models.Post
	SiteDesc  string
	SiteTitle string
}

func renderPosts(params *renderPostsParams) error {
	var (
		err        error
		dstFile    string
		link       string
		buf        *bytes.Buffer
		tplData    map[string]interface{}
		descGetter = func(post *models.Post) string {
			if post.Descripition != "" {
				return post.Descripition
			}
			return params.SiteDesc
		}
	)

	// build each post
	for _, p := range params.Posts {

		// build link
		link, dstFile, err = params.Ctx.createPostLink(p)
		if err != nil {
			zlog.Warnf("failed to build post link: %s, %s", p.LocalFile(), err)
			continue
		}
		p.Link = link
		dstFile = filepath.Join(params.OutputDir, dstFile)

		// convert markdown
		if err := p.Convert(markdown.Get()); err != nil {
			zlog.Warnf("failed to convert markdown post: %s, %s", p.LocalFile(), err)
			continue
		}

		buf = bytes.NewBuffer(nil)
		extData := map[string]interface{}{
			"post": p,
			"current": map[string]interface{}{
				"Title":       p.Title + " - " + params.SiteTitle,
				"Description": descGetter(p),
			},
		}
		tplData = params.Ctx.createTemplateData(extData)
		if err = params.Render.Execute(buf, p.Template, tplData); err != nil {
			zlog.Debugf("failed to render post: %s, %s", p.LocalFile(), err)
			continue
		}

		// save buffer to write content file later
		params.Ctx.SetOutput(dstFile, buf)
		zlog.Infof("post generated: %s", dstFile)

		t := p.Date()
		sitemap.Add(&sitemap.URL{Loc: link, LastMod: &t})

	}

	return nil
}

type renderPostListsParams struct {
	renderBaseParams
	Pager              *models.Pager
	Posts              []*models.Post
	PostPageLinkFormat string
	SiteTitle          string
}

func buildPostListTemplateData(params *renderPostListsParams, page int) (map[string]interface{}, *models.PagerItem) {
	pageItem := params.Pager.Page(page, params.PostPageLinkFormat)
	tplData := params.Ctx.createTemplateData(map[string]interface{}{
		"posts": models.PostsPageList(params.Posts, pageItem),
		"pager": pageItem,
		"current": map[string]interface{}{
			"Title": params.SiteTitle,
		},
	})
	return tplData, pageItem
}

func renderPostLists(params *renderPostListsParams) error {
	total := params.Pager.PageSize()
	tplName := constants.PostListTemplate
	for i := 1; i <= total; i++ {
		// build each page list
		buf := bytes.NewBuffer(nil)
		tplData, pageItem := buildPostListTemplateData(params, i)
		if err := params.Render.Execute(buf, tplName, tplData); err != nil {
			zlog.Warnf("failed to render post list: %d, %s", i, err)
			return err
		}
		dstFile := filepath.Join(params.OutputDir, pageItem.LocalFile)
		params.Ctx.SetOutput(dstFile, buf)
		sitemap.Add(&sitemap.URL{Loc: pageItem.Link})
		zlog.Infof("post list generated: %s", dstFile)
	}
	return nil
}
