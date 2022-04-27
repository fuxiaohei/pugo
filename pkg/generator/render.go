package generator

import (
	"pugo/pkg/models"
	"pugo/pkg/zlog"
)

func Render(siteData *models.SiteData, context *Context, opt *Option) error {
	if err := renderPosts(&renderPostsParams{
		Ctx:       context,
		Posts:     siteData.Posts,
		SiteDesc:  siteData.SiteConfig.Description,
		SiteTitle: siteData.SiteConfig.Title,
		Render:    siteData.Render,
		OutputDir: opt.OutputDir,
	}); err != nil {
		zlog.Warnf("render posts failed: %v", err)
		return err
	}
	postListParams := &renderPostListsParams{
		Ctx:                context,
		Pager:              siteData.PostsPager,
		Render:             siteData.Render,
		Posts:              siteData.Posts,
		PostPageLinkFormat: siteData.BuildConfig.PostPageLinkFormat,
		SiteTitle:          siteData.SiteConfig.Title,
		OutputDir:          opt.OutputDir,
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
		Ctx:               context,
		Tags:              siteData.Tags,
		PostPerPage:       siteData.BuildConfig.PostPerPage,
		TagPageLinkFormat: siteData.BuildConfig.TagPageLinkFormat,
		SiteTitle:         siteData.SiteConfig.Title,
		Render:            siteData.Render,
		OutputDir:         opt.OutputDir,
	}); err != nil {
		zlog.Warnf("render tags failed: %v", err)
		return err
	}
	if err := renderArchives(&renderArchivesParams{
		Ctx:          context,
		Posts:        siteData.Posts,
		SiteTitle:    siteData.SiteConfig.Title,
		Render:       siteData.Render,
		OutputDir:    opt.OutputDir,
		ArchivesLink: siteData.BuildConfig.ArchivesLink,
	}); err != nil {
		zlog.Warnf("render archives failed: %v", err)
		return err
	}

	if err := renderPages(&renderPagesParams{
		Ctx:       context,
		Pages:     siteData.Pages,
		SiteDesc:  siteData.SiteConfig.Description,
		SiteTitle: siteData.SiteConfig.Title,
		Render:    siteData.Render,
		OutputDir: opt.OutputDir,
	}); err != nil {
		zlog.Warnf("render pages failed: %v", err)
		return err
	}

	if err := renderFeedAtom(&renderFeedAtomParams{
		Ctx:           context,
		Posts:         siteData.Posts,
		SiteConfig:    siteData.SiteConfig,
		PostFeedLimit: siteData.BuildConfig.FeedPostLimit,
		OutputDir:     opt.OutputDir,
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
