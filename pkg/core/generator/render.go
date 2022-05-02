package generator

import (
	"pugo/pkg/core/models"
	"pugo/pkg/core/theme"
	"pugo/pkg/zlog"
)

type renderBaseParams struct {
	Ctx       *Context
	Render    *theme.Render
	OutputDir string
}

func newRenderBaseParams(siteData *models.SiteData, context *Context, opt *Option) renderBaseParams {
	return renderBaseParams{
		Ctx:       context,
		Render:    siteData.Render,
		OutputDir: opt.OutputDir,
	}
}

func Render(siteData *models.SiteData, context *Context, opt *Option) error {
	renderBase := newRenderBaseParams(siteData, context, opt)
	if err := renderPosts(&renderPostsParams{
		renderBaseParams: renderBase,
		Posts:            siteData.Posts,
		SiteDesc:         siteData.SiteConfig.Description,
		SiteTitle:        siteData.SiteConfig.Title,
	}); err != nil {
		zlog.Warnf("render posts failed: %v", err)
		return err
	}
	postListParams := &renderPostListsParams{
		renderBaseParams:   renderBase,
		Pager:              siteData.PostsPager,
		Posts:              siteData.Posts,
		PostPageLinkFormat: siteData.BuildConfig.PostPageLinkFormat,
		SiteTitle:          siteData.SiteConfig.Title,
	}
	if err := renderPostLists(postListParams); err != nil {
		zlog.Warnf("render post lists failed: %v", err)
		return err
	}
	if err := renderIndex(postListParams); err != nil {
		zlog.Warnf("render index failed: %v", err)
		return err
	}
	if err := renderTags(&renderTagsParams{
		renderBaseParams:  renderBase,
		Tags:              siteData.Tags,
		PostPerPage:       siteData.BuildConfig.PostPerPage,
		TagPageLinkFormat: siteData.BuildConfig.TagPageLinkFormat,
		SiteTitle:         siteData.SiteConfig.Title,
	}); err != nil {
		zlog.Warnf("render tags failed: %v", err)
		return err
	}
	if err := renderArchives(&renderArchivesParams{
		renderBaseParams: renderBase,
		Posts:            siteData.Posts,
		SiteTitle:        siteData.SiteConfig.Title,
		ArchivesLink:     siteData.BuildConfig.ArchivesLink,
	}); err != nil {
		zlog.Warnf("render archives failed: %v", err)
		return err
	}

	if err := renderPages(&renderPagesParams{
		renderBaseParams: renderBase,
		Pages:            siteData.Pages,
		SiteDesc:         siteData.SiteConfig.Description,
		SiteTitle:        siteData.SiteConfig.Title,
	}); err != nil {
		zlog.Warnf("render pages failed: %v", err)
		return err
	}

	if err := renderFeedAtom(&renderFeedAtomParams{
		renderBaseParams: renderBase,
		Posts:            siteData.Posts,
		SiteConfig:       siteData.SiteConfig,
		PostFeedLimit:    siteData.BuildConfig.FeedPostLimit,
	}); err != nil {
		zlog.Warnf("render feed atom failed: %v", err)
		return err
	}

	if err := renderSitemap(opt.OutputDir, context); err != nil {
		zlog.Warnf("render sitemap failed: %v", err)
		return err
	}

	return nil

}
