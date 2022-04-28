package generator

import (
	"path/filepath"
	"pugo/pkg/constants"
	"pugo/pkg/models"
	"pugo/pkg/utils"
	"pugo/pkg/watcher"
	"pugo/pkg/zlog"
	"time"

	"go.uber.org/atomic"
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
		if delta > 0 {
			nextGenerateTime = time.Now().Add(time.Second * -1)
			Generate(opt)
		}
	})

	baseDir := filepath.Base(opt.OutputDir)
	for _, dir := range constants.InitDirectories() {
		if dir == baseDir {
			continue
		}
		w.Add(dir)
		zlog.Debugf("watching dir: %s", dir)
	}
	go func() {
		for {
			event := <-w.Events()
			zlog.Infof("wathcing event: %s, %v", event.Name, event.Op)
			nextGenerateTime = time.Now().Add(constants.WatchPollingDuration)
		}
	}()
	zlog.Infof("watching...")
}
