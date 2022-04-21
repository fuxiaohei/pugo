package builder

import (
	"bytes"
	"pugo/internal/zlog"
)

func (b *Builder) buildIndex(ctx *buildContext) error {
	indexTpl := b.render.GetIndexTemplate()

	// first page
	tplData, _ := b.buildPostListTemplateData(ctx, 1)

	buf := bytes.NewBuffer(nil)
	if err := b.render.Execute(buf, indexTpl, tplData); err != nil {
		zlog.Warn("failed to render index", "err", err)
		return err
	}
	ctx.setBuffer("/index.html", buf)
	zlog.Info("posts: index rendered ok")
	return nil
}
