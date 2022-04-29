package generator

import (
	"bytes"
	"fmt"
	"html/template"
	"pugo/pkg/constants"
	"pugo/pkg/models"
	"pugo/pkg/utils"
	"pugo/pkg/zlog"
	"sort"
	"sync"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
	"go.uber.org/atomic"
)

type Context struct {
	outputs       sync.Map
	outputCounter *atomic.Int64

	templateData map[string]interface{}

	copingDirs []*models.CopyDir

	postSlugTemplate *template.Template
	tagLinkTemplate  *template.Template

	sitemap *models.Sitemap

	minifier *minify.M
}

func NewContext(s *models.SiteData, opt *Option) *Context {
	ctx := &Context{
		templateData:  map[string]interface{}{},
		copingDirs:    make([]*models.CopyDir, 0, len(s.BuildConfig.StaticAssetsDir)),
		outputCounter: atomic.NewInt64(0),
		sitemap:       models.NewSiteMap(s.Config.Site.Base),
	}

	for _, dir := range s.BuildConfig.StaticAssetsDir {
		ctx.copingDirs = append(ctx.copingDirs, &models.CopyDir{
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
		"local": opt.IsLocalServer,
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
func (ctx *Context) SetOutput(path string, buf *bytes.Buffer) *Context {
	ctx.outputs.Store(path, buf)
	return ctx
}

type outputFile struct {
	Path string
	Buf  *bytes.Buffer
}

func (ctx *Context) GetOutputs() []*outputFile {
	outputs := make(map[string]*bytes.Buffer)
	ctx.outputs.Range(func(key, value interface{}) bool {
		outputs[key.(string)] = value.(*bytes.Buffer)
		return true
	})
	result := make([]*outputFile, 0, len(outputs))
	for key, buf := range outputs {
		result = append(result, &outputFile{Path: key, Buf: buf})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Path < result[j].Path
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

func (ctx *Context) incrOutputCounter(delta int64) int64 {
	return ctx.outputCounter.Add(delta)
}

func (ctx *Context) addSitemap(url *models.SitemapURL) {
	ctx.sitemap.Add(url)
}

func (ctx *Context) getSitemap() *models.Sitemap {
	return ctx.sitemap
}

func (ctx *Context) appendCopyDir(srcDir, dstDir string) {
	ctx.copingDirs = append(ctx.copingDirs, &models.CopyDir{SrcDir: srcDir, DestDir: dstDir})
}

func (ctx *Context) MinifyHTML(raw []byte) ([]byte, error) {
	if ctx.minifier == nil {
		m := minify.New()
		m.Add("text/html", &html.Minifier{
			KeepComments:            false,
			KeepConditionalComments: true,
			KeepDefaultAttrVals:     true,
			KeepDocumentTags:        true,
			KeepEndTags:             false,
			KeepQuotes:              true,
			KeepWhitespace:          false,
		})
		ctx.minifier = m
	}
	return ctx.minifier.Bytes("text/html", raw)
}
