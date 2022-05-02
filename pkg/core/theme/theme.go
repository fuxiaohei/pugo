package theme

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"pugo/pkg/utils"
	"pugo/pkg/utils/zlog"
	"regexp"
	"sync"

	"github.com/BurntSushi/toml"
)

type Theme struct {
	Directory  string `toml:"directory"`
	ConfigFile string `toml:"config_file"`
}

var (
	reTemplateTag = regexp.MustCompile("{{ ?template \"([^\"]*)\" ?([^ ]*)? ?}}")
)

type namedTemplateFile struct {
	Name string
	Src  string
}

// Render renders the parsed data to static files.
type Render struct {
	dir        string
	configFile string
	config     *Config
	funcMap    template.FuncMap

	lock      sync.Mutex
	templates map[string]*template.Template
	cache     []*namedTemplateFile
}

func NewRender(cfg *Theme) (*Render, error) {
	r := &Render{
		dir:        cfg.Directory,
		configFile: cfg.ConfigFile,
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
	if err := r.loadTemplates(); err != nil {
		return err
	}
	return nil
}

func (r *Render) parseThemeConfig() error {
	r.config = NewDefaultConfig()

	// if no config file, use empty config
	if r.configFile == "" {
		return nil
	}
	configFile := filepath.Join(r.dir, r.configFile)

	fileBytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		zlog.Warnf("failed to load theme config file: %s", configFile)
		return err
	}

	if err = toml.Unmarshal(fileBytes, r.config); err != nil {
		zlog.Warnf("failed to parse theme config file: %s", configFile)
		return err
	}

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
			zlog.Warnf("failed to load template: %s, %s", path, err)
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

		zlog.Debugf("load template ok: %s", path)

		// release cache between twice render
		r.cache = nil

		return nil
	})

	r.templates = templates

	return err
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

func (r *Render) getTemplate(name string) *template.Template {
	r.lock.Lock()
	defer r.lock.Unlock()

	return r.templates[name]
}

// Execute executes template by name with data, and write into a Writer
func (r *Render) Execute(w io.Writer, name string, data interface{}) error {
	tpl := r.getTemplate(name)
	if tpl == nil {
		return fmt.Errorf("template '%s' is missing", name)
	}
	return tpl.ExecuteTemplate(w, name, data)
}

// GetIndexTemplate gets index template
func (r *Render) GetIndexTemplate() string {
	return r.config.IndexTemplate
}

// GetDir gets theme dir
func (r *Render) GetDir() string {
	return r.dir
}

// GetStaticDirs gets static directories
func (r *Render) GetStaticDirs() []string {
	return r.config.StaticDirs
}
