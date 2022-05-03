package ext

import (
	"pugo/pkg/core/configs"
	"pugo/pkg/ext/sitemap"
	"pugo/pkg/utils/zlog"
)

func Reload(cfg *configs.Config) {
	if cfg.Extension.Feed != nil {
		zlog.Debugf("feed reloaded, enabled:%v", cfg.Extension.Feed.Enabled)
	} else {
		zlog.Debugf("feed reloaded, nil, disabled")
	}

	if cfg.Extension.Sitemap != nil {
		sitemap.Init(cfg.Extension.Sitemap, cfg.Site.Base)
		zlog.Debugf("sitemap reloaded, enabled:%v", cfg.Extension.Sitemap.Enabled)
	} else {
		zlog.Debugf("sitemap reloaded, nil, disabled")
	}
}
