package generator

import (
	"path/filepath"
	"pugo/pkg/core/constants"
	"pugo/pkg/core/watcher"
	"pugo/pkg/ext"
	"pugo/pkg/utils"
	"pugo/pkg/utils/zlog"
	"time"

	"go.uber.org/atomic"
)

// Generate generates current site.
func Generate(opt *Option) error {
	st := time.Now()

	siteData, err := CreateSiteData(*opt.ConfigFileItem, &SiteDataParams{
		WithDrafts: opt.EnableDrafts,
	})
	if err != nil {
		zlog.Warnf("load site data failed: %v", err)
		return err
	}
	if opt.OutputDir == "" {
		opt.OutputDir = siteData.BuildConfig.OutputDir
	}
	zlog.Infof("output dir: %s", opt.OutputDir)

	// TODO: use a method to contains all extensions initialization
	ext.Reload(siteData.Config)

	context := NewContext(siteData, opt)

	if err = Render(siteData, context, opt); err != nil {
		zlog.Warnf("render failed: %v", err)
		return err
	}

	if err = Output(siteData, context, opt); err != nil {
		zlog.Warnf("output failed: %v", err)
		return err
	}
	zlog.Infof("generate %d files finished in %dms", context.getOutputCounter(), time.Since(st).Milliseconds())

	if opt.EnableWatch && !watchFlag.Load() {
		go Watch(opt)
	}
	return nil
}

var (
	nextGenerateTime = time.Now().Add(time.Second * -1)
	watchFlag        = atomic.NewBool(false)
)

func Watch(opt *Option) {

	w, err := watcher.New(constants.WatchPollingDuration)
	if err != nil {
		zlog.Warnf("watch failed: %v", err)
		return
	}

	watchFlag.Store(true)

	// use time loop to handle several events at once
	utils.Ticker(constants.WatchTickerDuaration, func() {
		delta := nextGenerateTime.Sub(time.Now())
		// zlog.Debugf("next generate time: %v", delta)
		if delta > 0 {
			// if nextGenerateTime is in the future, handle it
			nextGenerateTime = time.Now().Add(time.Second * -1)
			Generate(opt)
		}
	})

	var allDirs []string
	baseDir := filepath.Base(opt.OutputDir)
	for _, dir := range constants.InitDirectories() {
		if dir == baseDir {
			continue
		}
		subDirs, err := utils.GetSubDirectories(dir)
		if err != nil {
			zlog.Warnf("get sub directories failed: %v", err)
			allDirs = append(allDirs, dir)
		} else {
			allDirs = append(allDirs, subDirs...)
		}
	}
	allDirs = utils.UniqueStringsSlice(allDirs)
	for _, dir := range allDirs {
		w.Add(dir)
		zlog.Debugf("watching dir: %s", dir)
	}
	go func() {
		for {
			event := <-w.Events()
			zlog.Infof("wathcing event: %s, %v", event.Name, event.Op)
			// set a future time to avoid generating too frequently
			nextGenerateTime = time.Now().Add(time.Hour)
		}
	}()
	zlog.Infof("watching...")
}
