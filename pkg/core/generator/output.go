package generator

import (
	"bytes"
	"os"
	"path/filepath"
	"pugo/pkg/core/theme"
	"pugo/pkg/ext/markdown"
	"pugo/pkg/utils"
	"pugo/pkg/utils/zlog"
)

// Output outputs contents to destination directory.
func Output(s *SiteData, ctx *Context, outputDir string) error {
	if err := updateThemeCopyDirs(s.Render, ctx); err != nil {
		zlog.Warn("theme: failed to update copy dirs", "err", err)
		return err
	}
	if err := outputFiles(s, ctx); err != nil {
		return err
	}
	if err := copyAssets(outputDir, ctx); err != nil {
		return err
	}
	return nil
}

func updateThemeCopyDirs(r *theme.Render, ctx *Context) error {
	staticDirs := r.GetStaticDirs()
	themeDir := r.GetDir()
	for _, dir := range staticDirs {
		ctx.appendCopyDir(filepath.Join(themeDir, dir), dir)
	}
	return nil
}

func outputFiles(s *SiteData, ctx *Context) error {
	var (
		err   error
		fpath string
		buf   *bytes.Buffer
	)

	outputs := ctx.GetOutputs()
	for _, outputFile := range outputs {
		fpath = outputFile.Path
		buf = outputFile.Buf

		data := buf.Bytes()
		dataLen := len(data)
		if s.BuildConfig.EnableMinifyHTML {
			data, err = markdown.MinifyHTML(data)
			if err != nil {
				zlog.Warnf("output: failed to minify: %s, %s", fpath, err)
				data = buf.Bytes()
			} else {
				zlog.Debugf("minified ok: %s, %d -> %d", fpath, dataLen, len(data))
			}
		}
		if err = utils.WriteFile(fpath, data); err != nil {
			zlog.Warnf("output: failed to write file: %s, %s", fpath, err)
			continue
		}
		ctx.incrOutputCounter(1)
	}
	return nil
}

func copyAssets(outputDir string, ctx *Context) error {
	for _, dirData := range ctx.copingDirs {
		err := filepath.Walk(dirData.SrcDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			if utils.IsTempFile(path) {
				zlog.Debugf("skip temp file: %s", path)
				return nil
			}
			relPath, err := filepath.Rel(dirData.SrcDir, path)
			if err != nil {
				return nil
			}
			dstPath := filepath.Join(dirData.DestDir, relPath)
			dstPath = filepath.Join(outputDir, dstPath)
			if err := utils.CopyFile(path, dstPath); err != nil {
				zlog.Warnf("failed to copy file: %s, %s", dstPath, err)
				return err
			}
			zlog.Infof("assets copied: %s", dstPath)
			ctx.incrOutputCounter(1)
			return nil
		})
		if err != nil {
			zlog.Warnf("failed to copy assets: %s, %s", dirData.SrcDir, err)
			return err
		}
	}
	return nil
}
