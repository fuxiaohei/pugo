package generator

import (
	"bytes"
	"path/filepath"
	"pugo/pkg/core/constants"
	"pugo/pkg/core/models"
	"pugo/pkg/ext/sitemap"
	"pugo/pkg/utils/zlog"
	"strings"
)

type renderTagsParams struct {
	renderBaseParams
	Tags              []*models.TagPosts
	PostPerPage       int
	TagPageLinkFormat string
}

func renderTags(params *renderTagsParams) error {

	tplName := constants.PostListTemplate

	// build tag pages
	for _, tagData := range params.Tags {
		pager := models.NewPager(params.PostPerPage, len(tagData.Posts))

		total := pager.PageSize()
		for i := 1; i <= total; i++ {
			linkFormat := strings.ReplaceAll(params.TagPageLinkFormat, "{{.Tag}}", tagData.Tag.Name)
			pageItem := pager.Page(i, linkFormat)
			dstFile := pageItem.LocalFile

			buf := bytes.NewBuffer(nil)
			posts := models.PostsPageList(tagData.Posts, pageItem)
			tplData := params.Ctx.createTemplateData(map[string]interface{}{
				"posts": posts,
				"pager": pageItem,
				"tag":   tagData.Tag,
				"current": map[string]interface{}{
					"Title":       tagData.Tag.Name + "-" + params.SiteTitle,
					"Description": tagData.Tag.Name + " - " + params.SiteDescription,
				},
			})

			if err := params.Render.Execute(buf, tplName, tplData); err != nil {
				zlog.Warnf("failed to render tag page: %s, %d, %s", tagData.Tag.Name, i, err)
				return err
			}
			dstFile = filepath.Join(params.OutputDir, dstFile)
			params.Ctx.SetOutput(dstFile, pageItem.Link, buf)
			zlog.Infof("tag page generated: %s", dstFile)

			t := posts[0].Date()

			sitemap.Add(&sitemap.URL{Loc: pageItem.Link, LastMod: &t})

			// tag list index.html
			if i == 1 {
				dstFile = filepath.Join(params.OutputDir, tagData.Tag.LocalFile)
				params.Ctx.SetOutput(dstFile, tagData.Tag.Link, buf)
				sitemap.Add(&sitemap.URL{Loc: tagData.Tag.Link, LastMod: &t})
				zlog.Infof("tag page generated: %s", dstFile)
			}
		}
	}

	return nil
}
