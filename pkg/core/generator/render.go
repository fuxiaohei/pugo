package generator

import (
	"pugo/pkg/core/i18n"
	"pugo/pkg/core/theme"
	"pugo/pkg/ext/feed"
	"pugo/pkg/ext/sitemap"
	"pugo/pkg/utils/zlog"
)

type renderBaseParams struct {
	Ctx             *Context
	Render          *theme.Render
	OutputDir       string
	SiteTitle       string
	SiteSubTitle    string
	SiteDescription string
	I18n            *i18n.I18n
}

func newRenderBaseParams(siteData *SiteData, context *Context, opt *Option) renderBaseParams {
	return renderBaseParams{
		Ctx:             context,
		Render:          siteData.Render,
		OutputDir:       opt.OutputDir,
		SiteTitle:       siteData.SiteConfig.Title,
		SiteSubTitle:    siteData.SiteConfig.SubTitle,
		SiteDescription: siteData.SiteConfig.Description,
		I18n:            siteData.I18n,
	}
}

func Render(siteData *SiteData, context *Context, opt *Option) error {
	renderBase := newRenderBaseParams(siteData, context, opt)
	if err := renderPosts(&renderPostsParams{
		renderBaseParams: renderBase,
		Posts:            siteData.Posts,
	}); err != nil {
		zlog.Warnf("render posts failed: %v", err)
		return err
	}
	postListParams := &renderPostListsParams{
		renderBaseParams:   renderBase,
		Pager:              siteData.PostsPager,
		Posts:              siteData.Posts,
		PostPageLinkFormat: siteData.BuildConfig.PostPageLinkFormat,
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
	}); err != nil {
		zlog.Warnf("render tags failed: %v", err)
		return err
	}
	if err := renderArchives(&renderArchivesParams{
		renderBaseParams: renderBase,
		Posts:            siteData.Posts,
		ArchivesLink:     siteData.BuildConfig.ArchivesLink,
	}); err != nil {
		zlog.Warnf("render archives failed: %v", err)
		return err
	}

	if err := renderPages(&renderPagesParams{
		renderBaseParams: renderBase,
		Pages:            siteData.Pages,
	}); err != nil {
		zlog.Warnf("render pages failed: %v", err)
		return err
	}

	if err := renderErrorPage(&renderErrorPageParams{
		renderBaseParams: renderBase,
		SiteTitle:        siteData.SiteConfig.Title,
	}); err != nil {
		zlog.Warnf("render error page failed: %v", err)
		return err
	}

	// render feed
	out, err := feed.Render(&feed.RenderParams{
		Posts:       siteData.Posts,
		SiteBaseURL: siteData.SiteConfig.Base,
		SiteTitle:   siteData.SiteConfig.Title,
		OutputDir:   opt.OutputDir,
		Config:      siteData.Config.Extension.Feed,
	})
	if err != nil {
		zlog.Warnf("render feed failed: %v", err)
		return err
	}
	if out != nil {
		context.SetOutput(out.Path, out.Link, out.Buf)
		zlog.Infof("atom feed generated: %s", out.Path)
	}

	// render sitemap
	out, err = sitemap.Render(siteData.Config.Extension.Sitemap, opt.OutputDir)
	if err != nil {
		zlog.Warnf("render sitemap failed: %v", err)
		return err
	}
	if out != nil {
		context.SetOutput(out.Path, out.Link, out.Buf)
		zlog.Infof("sitemap generated: %s", out.Path)
	}

	return nil

}
