package generator

import (
	"pugo/pkg/core/configs"
	"pugo/pkg/core/constants"
	"pugo/pkg/core/models"
	"pugo/pkg/core/theme"
	"pugo/pkg/utils/zlog"
)

type SiteData struct {
	Posts      []*models.Post
	PostsPager *models.Pager
	Tags       []*models.TagPosts

	Pages []*models.Page

	Config      *configs.Config
	ConfigType  constants.ConfigType
	BuildConfig *configs.Build
	SiteConfig  *configs.Site

	Render *theme.Render
}

// NewSiteData returns a new default sote data.
func NewSiteData() *SiteData {
	return &SiteData{
		Posts: make([]*models.Post, 0),
		Pages: make([]*models.Page, 0),
	}
}

// CreateSiteData creates a new site data from the given config.
func CreateSiteData(item constants.ConfigFileItem) (*SiteData, error) {
	siteData := NewSiteData()

	// load config
	cfg, err := configs.LoadFromFile(item)
	if err != nil {
		zlog.Warnf("load config file failed: %v", err)
		return nil, err
	}
	zlog.Debugf("load config ok: %s", item.File)
	siteData.ConfigType = item.Type
	siteData.Config = cfg
	siteData.BuildConfig = cfg.Build
	siteData.SiteConfig = cfg.Site

	// load theme
	render, err := theme.NewRender(cfg.Theme)
	if err != nil {
		zlog.Warnf("load theme failed: %v", err)
		return nil, err
	}
	siteData.Render = render

	// load contents
	if siteData.Posts, err = models.LoadPosts(); err != nil {
		zlog.Warnf("load posts failed: %v", err)
		return nil, err
	}
	if siteData.Pages, err = models.LoadPages(); err != nil {
		zlog.Warnf("load pages failed: %v", err)
		return nil, err
	}

	siteData.fullfill()

	return siteData, nil
}

// FulFill makes relative data available in source data
func (s *SiteData) fullfill() {

	// set post author data
	for _, post := range s.Posts {
		post.Author = s.assignAuthor(post.AuthorName)
		for _, t := range post.Tags {
			post.TagLinks = append(post.TagLinks, &models.TagLink{Name: t})
		}
	}

	// build tag posts
	s.Tags = models.BuildTagPosts(s.Posts)
	zlog.Infof("load tags ok: %d", len(s.Tags))

	// set page author
	for _, page := range s.Pages {
		page.Author = s.assignAuthor(page.AuthorName)
	}

	// set post pager data
	s.PostsPager = models.NewPager(s.BuildConfig.PostPerPage, len(s.Posts))
	zlog.Infof("load pagination ok: %d", s.PostsPager.PageSize())
}

func (s *SiteData) assignAuthor(name string) *models.Author {
	if name == "" {
		return s.Config.Author[0]
	}
	author := s.Config.GetAuthor(name)
	if author == nil {
		author = models.NewAuthor(name)
	}
	return author
}
