package models

import (
	"pugo/pkg/theme"
	"pugo/pkg/zlog"
)

type SiteData struct {
	Posts      []*Post
	PostsPager *Pager
	Tags       []*TagPosts

	Pages []*Page

	Config      *Config
	BuildConfig *BuildConfig
	SiteConfig  *SiteConfig

	Render *theme.Render
}

// NewSiteData returns a new default sote data.
func NewSiteData() *SiteData {
	return &SiteData{
		Posts: make([]*Post, 0),
		Pages: make([]*Page, 0),
	}
}

func LoadSiteData(configFile string) (*SiteData, error) {
	siteData := NewSiteData()

	// load config
	cfg, err := LoadConfigFromFile(configFile)
	if err != nil {
		zlog.Warnf("load config file failed: %v", err)
		return nil, err
	}
	siteData.Config = cfg
	siteData.BuildConfig = cfg.BuildConfig
	siteData.SiteConfig = cfg.Site

	// load theme
	render, err := theme.NewRender(cfg.Theme)
	if err != nil {
		zlog.Warnf("load theme failed: %v", err)
		return nil, err
	}
	siteData.Render = render

	// load contents
	if err = LoadPosts(siteData); err != nil {
		zlog.Warnf("load posts failed: %v", err)
		return nil, err
	}
	if err = LoadPages(siteData); err != nil {
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
			post.TagLinks = append(post.TagLinks, &TagLink{Name: t})
		}
	}

	// build tag posts
	s.Tags = BuildTagPosts(s.Posts)
	zlog.Infof("load tags ok: %d", len(s.Tags))

	// set page author
	for _, page := range s.Pages {
		page.Author = s.assignAuthor(page.AuthorName)
	}

	// set post pager data
	s.PostsPager = NewPager(s.BuildConfig.PostPerPage, len(s.Posts))
	zlog.Infof("load pagination ok: %d", s.PostsPager.PageSize())
}

func (s *SiteData) assignAuthor(name string) *Author {
	if name == "" {
		return s.Config.Author[0]
	}
	author := s.Config.GetAuthor(name)
	if author == nil {
		author = NewAuthor(name)
	}
	return author
}
