package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"pugo/internal/utils"
	"pugo/internal/zlog"
)

// Output outputs contents to destination directory.
func (b *Builder) Output(ctx *buildContext) error {
	if err := b.outputCompiledFiles(ctx); err != nil {
		return err
	}
	if err := b.copyAssets(ctx); err != nil {
		return err
	}
	return nil
}

func (b *Builder) outputCompiledFiles(ctx *buildContext) error {
	if b.outputDir == "" {
		return fmt.Errorf("output directory is empty")
	}
	outputs := ctx.getOutputs()
	for fpath, buf := range outputs {
		fpath = filepath.Join(b.outputDir, fpath)
		if err := utils.WriteFile(fpath, buf.Bytes()); err != nil {
			zlog.Warn("output: failed to write", "path", fpath, "err", err)
			continue
		}
		zlog.Debug("output: write ok", "path", fpath)
	}
	return nil
}

func (b *Builder) copyAssets(ctx *buildContext) error {
	for _, dirData := range ctx.copingDirs {
		err := filepath.Walk(dirData.SourceDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			relPath, err := filepath.Rel(dirData.SourceDir, path)
			if err != nil {
				return nil
			}
			dstPath := filepath.Join(dirData.DstDir, relPath)
			dstPath = filepath.Join(b.outputDir, dstPath)
			if err := utils.CopyFile(path, dstPath); err != nil {
				zlog.Warn("copyAssets: failed to copy", "src", path, "dst", dstPath, "err", err)
				return err
			}
			zlog.Info("copyAssets: copied ok", "src", path, "dst", dstPath)
			return nil
		})
		if err != nil {
			zlog.Warn("copyAssets: failed to copy", "src", dirData.SourceDir, "dst", dirData.DstDir, "err", err)
			return err
		}
	}
	return nil
}
