package generator

import (
	"pugo/pkg/models"
	"pugo/pkg/zlog"
	"time"
)

// Generate generates current site.
func Generate(opt *Option) error {
	st := time.Now()

	siteData, err := models.LoadSiteData(opt.ConfigFile)
	if err != nil {
		zlog.Warnf("load site data failed: %v", err)
		return err
	}
	if opt.OutputDir == "" {
		opt.OutputDir = siteData.BuildConfig.OutputDir
	}

	context := NewContext(siteData)

	if err = Render(siteData, context, opt); err != nil {
		zlog.Warnf("render failed: %v", err)
		return err
	}

	if err = Output(siteData, context, opt.OutputDir); err != nil {
		zlog.Warnf("output failed: %v", err)
		return err
	}
	zlog.Infof("generate %d files finished in %dms", context.getOutputCounter(), time.Since(st).Milliseconds())
	return nil
}
