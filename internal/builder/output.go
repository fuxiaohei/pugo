package builder

import (
	"fmt"
	"path/filepath"
	"pugo/internal/utils"
	"pugo/internal/zlog"
)

// Output outputs contents to destination directory.
func (b *Builder) Output(ctx *buildContext) error {
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
		zlog.Info("output: write ok", "path", fpath)
	}
	return nil
}
