package generator

import (
	"bytes"
	"fmt"
	"html/template"
	"pugo/pkg/core/constants"
	"pugo/pkg/core/models"
	"pugo/pkg/utils"
	"pugo/pkg/utils/zlog"
	"sort"
	"sync"

	"go.uber.org/atomic"
)

type Context struct {
	outputs       sync.Map
	outputCounter *atomic.Int64

	templateData map[string]interface{}

	copingDirs []*models.CopyDir

	postSlugTemplate *template.Template
	tagLinkTemplate  *template.Template

	allLinkFiles sync.Map
}

func NewContext(s *SiteData, opt *Option) *Context {
	ctx := &Context{
		templateData:  map[string]interface{}{},
		copingDirs:    make([]*models.CopyDir, 0, len(s.BuildConfig.StaticAssetsDir)),
		outputCounter: atomic.NewInt64(0),
	}

	for _, dir := range s.BuildConfig.StaticAssetsDir {
		ctx.copingDirs = append(ctx.copingDirs, &models.CopyDir{
			SrcDir:  dir,
			DestDir: dir,
		})
	}

	ctx.copingDirs = append(ctx.copingDirs, &models.CopyDir{
		SrcDir:  "content/public",
		DestDir: ".",
	})

	// build post slug template
	tpl, err := template.New("post-slug").Parse(s.BuildConfig.PostLinkFormat)
	if err != nil {
		zlog.Warn("posts: failed to parse post slug template", "err", err)
		return nil
	}
	ctx.postSlugTemplate = tpl
	zlog.Debugf("load post slug template: %s", s.BuildConfig.PostLinkFormat)

	// build tag link template
	tpl, err = template.New("tag").Parse(s.BuildConfig.TagLinkFormat)
	if err != nil {
		zlog.Warn("posts: failed to parse tag link template", "err", err)
		return nil
	}
	ctx.tagLinkTemplate = tpl
	zlog.Debugf("load tag link template: %s", s.BuildConfig.TagLinkFormat)

	// prepare global template data
	ctx.templateData["site"] = s.Config.Site
	ctx.templateData["menu"] = s.Config.Menu
	// TODO: support authors list
	ctx.templateData["author"] = s.Config.Author[0]

	// update tag data
	var tagTemplateData []*models.TagLink
	for _, tagData := range s.Tags {
		ctx.updateTagLink(tagData.Tag)
		tagTemplateData = append(tagTemplateData, tagData.Tag)
	}
	ctx.templateData["tags"] = tagTemplateData
	// add pugo data
	ctx.templateData["pugo"] = map[string]interface{}{
		"Name":    constants.AppName(),
		"Version": constants.AppVersion(),
		"Github":  constants.AppGithubLink(),
	}
	ctx.templateData["server"] = map[string]interface{}{
		"Local": opt.IsLocalServer,
	}
	ctx.templateData["extension"] = s.Config.Extension

	themeConfig := s.Render.GetConfig()
	ctx.templateData["theme"] = map[string]interface{}{
		"EnableDarkMode":  themeConfig.EnableDarkMode,
		"ShowPuGoVersion": themeConfig.ShowPuGoVersion,
	}

	return ctx
}

func (ctx *Context) updateTagLink(t *models.TagLink) {
	data := map[string]interface{}{
		"Tag": t.Name,
	}
	var buf bytes.Buffer
	ctx.tagLinkTemplate.Execute(&buf, data)

	t.Link = buf.String()
	t.LocalFile = utils.FormatIndexHTML(buf.String())
}

// SetOutput sets the output for the given key.
func (ctx *Context) SetOutput(path, link string, buf *bytes.Buffer) *Context {
	ctx.outputs.Store(path, &models.OutputFile{
		Path: path,
		Link: link,
		Buf:  buf,
	})
	return ctx
}

func (ctx *Context) GetOutputs() []*models.OutputFile {
	var outputs []*models.OutputFile
	ctx.outputs.Range(func(key, value interface{}) bool {
		outputs = append(outputs, value.(*models.OutputFile))
		return true
	})
	sort.Slice(outputs, func(i, j int) bool {
		return outputs[i].Path < outputs[j].Path
	})
	return outputs
}

func (ctx *Context) GetRecordFiles() []*models.OutputFile {
	var result []*models.OutputFile
	ctx.allLinkFiles.Range(func(key, value interface{}) bool {
		result = append(result, &models.OutputFile{Link: key.(string), Path: value.(string)})
		return true
	})
	return result
}

func (ctx *Context) createPostLink(p *models.Post) (string, string, error) {
	slugData := map[string]interface{}{
		"Slug": p.Slug,
		"Date": map[string]interface{}{
			"Year":  int(p.Date().Year()),
			"Month": fmt.Sprintf("%02d", int(p.Date().Month())),
			"Day":   fmt.Sprintf("%02d", int(p.Date().Day())),
		},
	}
	var buf bytes.Buffer
	err := ctx.postSlugTemplate.Execute(&buf, slugData)
	return buf.String(), utils.FormatIndexHTML(buf.String()), err
}

func (ctx *Context) createTemplateData(data map[string]interface{}) map[string]interface{} {
	if data == nil {
		data = make(map[string]interface{})
	}
	for k, v := range ctx.templateData {
		// do not override existing data
		if _, ok := data[k]; !ok {
			data[k] = v
		}
	}
	return data
}

// outputLength returns the length of the output buffer
func (ctx *Context) getOutputCounter() int64 {
	return ctx.outputCounter.Load()
}

func (ctx *Context) appendCopyDir(srcDir, dstDir string) {
	ctx.copingDirs = append(ctx.copingDirs, &models.CopyDir{SrcDir: srcDir, DestDir: dstDir})
}

func (ctx *Context) recordLinkFile(link, file string) {
	ctx.allLinkFiles.Store(link, file)
	ctx.outputCounter.Inc()
}
