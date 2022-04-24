package builder

import (
	"bytes"
	"fmt"
	"html/template"
	"pugo/internal/model"
	"pugo/internal/zlog"
	"sync"

	"go.uber.org/atomic"
)

type buildContext struct {
	lock               sync.RWMutex
	outputs            map[string]*bytes.Buffer
	outputCounter      *atomic.Int64
	globalTemplateData map[string]interface{}
	copingDirs         []*model.CopyDir

	postSlugTemplate *template.Template
	tagLinkTemplate  *template.Template
}

func newBuildContext(s *SourceData) *buildContext {
	ctx := &buildContext{
		outputs:            make(map[string]*bytes.Buffer),
		globalTemplateData: map[string]interface{}{},
		copingDirs:         make([]*model.CopyDir, 0, len(s.BuildConfig.StaticAssetsDir)),
		outputCounter:      atomic.NewInt64(0),
	}

	for _, dir := range s.BuildConfig.StaticAssetsDir {
		ctx.copingDirs = append(ctx.copingDirs, &model.CopyDir{
			SrcDir:  dir,
			DestDir: dir,
		})
	}

	// build post slug template
	tpl, err := template.New("post-slug").Parse(s.BuildConfig.PostLinkFormat)
	if err != nil {
		zlog.Warn("posts: failed to parse post slug template", "err", err)
		return nil
	}
	ctx.postSlugTemplate = tpl
	zlog.Debug("posts: built post slug template", "format", s.BuildConfig.PostLinkFormat)

	// build tag link template
	tpl, err = template.New("tag").Parse(s.BuildConfig.TagLinkFormat)
	if err != nil {
		zlog.Warn("posts: failed to parse tag link template", "err", err)
		return nil
	}
	ctx.tagLinkTemplate = tpl
	zlog.Debug("posts: built post tag link template", "format", s.BuildConfig.TagLinkFormat)

	// prepare global template data
	ctx.globalTemplateData["site"] = s.Config.Site
	ctx.globalTemplateData["menu"] = s.Config.Menu

	// update tag data
	var tagTemplateData []*model.TagLink
	for _, tagData := range s.Tags {
		ctx.updatTagLink(tagData.Tag)
		tagTemplateData = append(tagTemplateData, tagData.Tag)
	}
	ctx.globalTemplateData["tags"] = tagTemplateData

	return ctx
}

func (bc *buildContext) setBuffer(fpath string, buf *bytes.Buffer) {
	bc.lock.Lock()
	defer bc.lock.Unlock()
	bc.outputs[fpath] = buf
}

func (bc *buildContext) getOutputs() map[string]*bytes.Buffer {
	bc.lock.RLock()
	defer bc.lock.RUnlock()
	return bc.outputs
}

func (bc *buildContext) appendCopyDir(srcDir, dstDir string) {
	bc.copingDirs = append(bc.copingDirs, &model.CopyDir{SrcDir: srcDir, DestDir: dstDir})
}

func (bc *buildContext) buildPostLink(p *model.Post) (string, string, error) {
	slugData := map[string]interface{}{
		"Slug": p.Slug,
		"Date": map[string]interface{}{
			"Year":  int(p.Date().Year()),
			"Month": fmt.Sprintf("%02d", int(p.Date().Month())),
			"Day":   fmt.Sprintf("%02d", int(p.Date().Day())),
		},
	}
	var buf bytes.Buffer
	err := bc.postSlugTemplate.Execute(&buf, slugData)
	return buf.String(), model.FormatIndexHTML(buf.String()), err
}

func (bc *buildContext) updatTagLink(t *model.TagLink) {
	data := map[string]interface{}{
		"Tag": t.Name,
	}
	var buf bytes.Buffer
	bc.tagLinkTemplate.Execute(&buf, data)

	t.Link = buf.String()
	t.LocalFile = model.FormatIndexHTML(buf.String())
}

func (bc *buildContext) buildTemplateData(data map[string]interface{}) map[string]interface{} {
	if data == nil {
		data = make(map[string]interface{})
	}
	for k, v := range bc.globalTemplateData {
		// do not override existing data
		if _, ok := data[k]; !ok {
			data[k] = v
		}
	}
	return data
}

// outputLength returns the length of the output buffer
func (bc *buildContext) getOutputCounter() int64 {
	return bc.outputCounter.Load()
}

func (bc *buildContext) incrOutputCounter(delta int64) int64 {
	return bc.outputCounter.Add(delta)
}
