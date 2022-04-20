package builder

import (
	"fmt"
	"haisite/internal/utils"
	"haisite/internal/zlog"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/BurntSushi/toml"
)

var (
	reTemplateTag = regexp.MustCompile("{{ ?template \"([^\"]*)\" ?([^ ]*)? ?}}")
)

// ThemeConfig is the theme config.
type ThemeConfig struct {
	Name          string   `toml:"name"`
	IndexTemplate string   `toml:"index_template"`
	Extension     []string `toml:"extension"`
}

// NewDefaultThemeConfig returns default theme config
func NewDefaultThemeConfig() *ThemeConfig {
	return &ThemeConfig{
		Name:          "theme",
		Extension:     []string{".html"},
		IndexTemplate: "post-list.html",
	}
}

// Render renders the parsed data to static files.
type Render struct {
	dir        string
	configFile string
	config     *ThemeConfig
	funcMap    template.FuncMap

	lock      sync.Mutex
	templates map[string]*template.Template
	cache     []*namedTemplateFile
}

func NewRender(dir, cfgFile string) (*Render, error) {
	r := &Render{
		dir:        dir,
		configFile: cfgFile,
		funcMap:    make(template.FuncMap),
	}
	r.initDefaultFuncMap()
	return r, r.Parse()
}

func (r *Render) initDefaultFuncMap() {
	r.funcMap["HTML"] = func(v interface{}) template.HTML {
		if str, ok := v.(string); ok {
			return template.HTML(str)
		}
		if b, ok := v.([]byte); ok {
			return template.HTML(string(b))
		}
		return template.HTML(fmt.Sprintf("%v", v))
	}
}

// Parse parses theme config and template files.
func (r *Render) Parse() error {
	if err := r.parseThemeConfig(); err != nil {
		return err
	}
	zlog.Info("theme: parse theme config", "dir", r.dir, "name", r.config.Name)
	if err := r.loadTemplates(); err != nil {
		return err
	}
	return nil
}

func (r *Render) parseThemeConfig() error {
	r.config = NewDefaultThemeConfig()

	// if no config file, use empty config
	if r.configFile == "" {
		return nil
	}
	configFile := filepath.Join(r.dir, r.configFile)

	fileBytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		zlog.Warn("theme: failed to parse theme config", "err", err, "file", configFile)
		return err
	}

	var config ThemeConfig
	if err = toml.Unmarshal(fileBytes, &config); err != nil {
		zlog.Warn("theme: failed to parse theme config", "err", err, "file", configFile)
		return err
	}
	r.config = &config

	return nil
}

func (r *Render) loadTemplates() error {
	r.lock.Lock()
	defer r.lock.Unlock()

	templates := make(map[string]*template.Template, len(r.templates))
	r.cache = make([]*namedTemplateFile, 0, len(r.templates))

	err := filepath.Walk(r.dir, func(path string, fi os.FileInfo, err error) error {
		if fi.IsDir() {
			return nil
		}
		ext := filepath.Ext(path)
		if !utils.Contains(r.config.Extension, ext) {
			return nil
		}

		tpl, err := filepath.Rel(r.dir, path) // get relative path
		if err != nil {
			return err
		}
		if err := r.loadOneTemplate(path, tpl); err != nil {
			zlog.Warn("theme: failed to load template", "err", err, "file", path)
			return err
		}

		var (
			baseTmpl    *template.Template
			currentTmpl *template.Template
		)

		for i, nt := range r.cache {
			if i == 0 {
				baseTmpl = template.New(nt.Name)
				currentTmpl = baseTmpl
			} else {
				currentTmpl = baseTmpl.New(nt.Name)
			}

			if _, err := currentTmpl.Funcs(r.funcMap).Parse(nt.Src); err != nil {
				return err
			}
			i++
		}
		templates[tpl] = baseTmpl

		zlog.Info("theme: load template ok", "file", path)

		// release cache between twice render
		r.cache = nil

		return nil
	})

	r.templates = templates

	return err
}

type namedTemplateFile struct {
	Name string
	Src  string
}

func (r *Render) loadOneTemplate(path, rel string) error {
	// already loaded in cache
	for _, t := range r.cache {
		if t.Name == rel {
			return nil
		}
	}

	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	tpl := &namedTemplateFile{
		Name: rel,
		Src:  string(fileBytes),
	}
	r.cache = append(r.cache, tpl)

	// parse {{template}} including sub-template
	for _, raw := range reTemplateTag.FindAllString(tpl.Src, -1) {
		parsed := reTemplateTag.FindStringSubmatch(raw)
		tplPath := parsed[1]
		tplExt := filepath.Ext(tplPath)
		if !utils.Contains(r.config.Extension, tplExt) {
			continue
		}
		fullPath := filepath.Join(r.dir, tplPath)
		if err := r.loadOneTemplate(fullPath, tplPath); err != nil {
			return err
		}
	}

	return nil
}

// Template gets template by name
func (r *Render) GetTemplate(name string) *template.Template {
	r.lock.Lock()
	defer r.lock.Unlock()

	return r.templates[name]
}

// Execute executes template by name with data, and write into a Writer
func (r *Render) Execute(w io.Writer, name string, data interface{}) error {
	tpl := r.GetTemplate(name)
	if tpl == nil {
		return fmt.Errorf("template '%s' is missing", name)
	}
	return tpl.ExecuteTemplate(w, name, data)
}
