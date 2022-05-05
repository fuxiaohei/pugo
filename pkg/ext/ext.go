package ext

import (
	"pugo/pkg/core/configs"
	"pugo/pkg/ext/sitemap"
	"pugo/pkg/utils/zlog"
)

func Reload(cfg *configs.Config) {
	ext := cfg.Extension
	if ext.Feed != nil {
		zlog.Debugf("feed reloaded, enabled:%v", ext.Feed.Enabled)
	} else {
		zlog.Debugf("feed reloaded, nil, disabled")
	}

	if ext.Sitemap != nil {
		sitemap.Init(ext.Sitemap, cfg.Site.Base)
		zlog.Debugf("sitemap reloaded, enabled:%v", ext.Sitemap.Enabled)
	} else {
		zlog.Debugf("sitemap reloaded, nil, disabled")
	}

	as := ext.Analytics
	if as.GoogleAnalytics.Enabled {
		zlog.Debugf("analytics: GoogleAnalytics enabled")
	}
	if as.Plausible.Enabled {
		zlog.Debugf("analytics: Plausible enabled")
	}

	ct := ext.Comments
	if ct.Enabled {
		zlog.Debugf("comments: enabled")
		if ct.Disqus.Enabled {
			zlog.Debugf("comments: Disqus enabled")
		}
	}
}
