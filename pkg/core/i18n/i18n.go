package i18n

import (
	"fmt"
	"os"
	"path/filepath"
	"pugo/pkg/utils"
	"pugo/pkg/utils/zlog"
	"sort"
	"strings"
	"sync"
)

const (
	i18nLabelKey     = "i18n_label"
	i18nURLPrefixKey = "i18n_url_prefix"
)

type (
	I18n struct {
		defaultLang string
		langs       sync.Map
	}
	i18nLang struct {
		name      string
		localFile string
		values    map[string]*i18nText
		isCurrent bool
		isDefault bool
		urlPrefix string
	}
	i18nText struct {
		ID     string `json:"id,omitempty"`
		Text   string `json:"text,omitempty"`
		Plural string `json:"plural,omitempty"`
	}
)

func (i18n *I18n) Get(lang string) *i18nLang {
	if lang == "" {
		lang = i18n.defaultLang
	}
	langs := LangCode(lang)
	for _, l := range langs {
		il, ok := i18n.langs.Load(l)
		if ok {
			ilClone := il.(*i18nLang).clone()
			ilClone.isCurrent = true
			return ilClone
		}
	}
	return &i18nLang{
		name:      lang,
		isCurrent: true,
	}
}

// Languages returns all languages of the i18n instance.
func (in *I18n) Languages() []*i18nLang {
	var langs []*i18nLang
	in.langs.Range(func(key, value interface{}) bool {
		lang := value.(*i18nLang)
		langs = append(langs, lang)
		return true
	})
	sort.Slice(langs, func(i, j int) bool {
		return langs[i].Name() < langs[j].Name()
	})
	return langs
}

// Name returns the name of the language.
func (il *i18nLang) Name() string {
	return il.name
}

// Label returns the label of language.
func (il *i18nLang) Label() string {
	text, ok := il.values[i18nLabelKey]
	if !ok {
		return il.name
	}
	return text.Text
}

// Tr returns the translated text of the given key.
func (il *i18nLang) Tr(key string) string {
	text, ok := il.values[key]
	if !ok {
		return key
	}
	return text.Text
}

// Trf returns the formartted translated text of the given key.
func (il *i18nLang) Trf(key string, args ...interface{}) string {
	text, ok := il.values[key]
	if !ok {
		return key
	}
	return fmt.Sprintf(text.Text, args...)
}

// Plural returns the plural form of the given language.
func (il *i18nLang) Plural(key string, n int) string {
	text, ok := il.values[key]
	if !ok {
		return key
	}
	if n > 1 {
		return fmt.Sprintf(text.Plural, n)
	}
	return fmt.Sprintf(text.Text, n)
}

// TrLink returns the formatted link of language.
func (il *i18nLang) TrLink(link string) string {
	link = il.RawLink(link)
	// if is default, no need add i18n url prefix
	if il.isDefault {
		return link
	}
	link = filepath.Join(il.urlPrefix, link)
	link = "/" + strings.TrimPrefix(link, "/")
	return link
}

func (il *i18nLang) RawLink(link string) string {
	if il.urlPrefix == "" {
		return link
	}
	// if i18n is in current, trim i18n url prefix
	if il.isCurrent {
		link = strings.TrimPrefix(link, il.urlPrefix)
	}
	if link == "" {
		link = "/"
	}
	return link
}

func (il *i18nLang) clone() *i18nLang {
	lang := &i18nLang{
		name:      il.name,
		localFile: il.localFile,
		values:    make(map[string]*i18nText),
		isCurrent: il.isCurrent,
		isDefault: il.isDefault,
		urlPrefix: il.urlPrefix,
	}
	for k, txt := range il.values {
		lang.values[k] = txt
	}
	return lang
}

func newI18nLang(file string, data map[string]string) *i18nLang {
	basename := filepath.Base(file)
	ext := filepath.Ext(basename)
	lang := &i18nLang{
		name:      strings.TrimSuffix(basename, ext),
		localFile: file,
		values:    make(map[string]*i18nText),
	}
	for key, value := range data {
		isPlural := false
		if strings.HasSuffix(key, "_plural") {
			key = strings.TrimSuffix(key, "_plural")
			isPlural = true
		}
		text, ok := lang.values[key]
		if !ok {
			text = &i18nText{
				ID: key,
			}
			lang.values[key] = text
		}
		if isPlural {
			text.Plural = value
		} else {
			text.Text = value
		}
	}
	text, ok := lang.values[i18nURLPrefixKey]
	if ok {
		lang.urlPrefix = text.Text
	}
	return lang
}

// Load loads all i18n files from the given directory.
func Load(dir string, defaultLang string) (*I18n, error) {
	in := &I18n{
		defaultLang: defaultLang,
	}
	defaultLanguages := LangCode(defaultLang)
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// skip directory
		if info.IsDir() {
			return nil
		}
		data := make(map[string]string)
		if err := utils.LoadTOMLFile(path, &data); err != nil {
			return nil
		}
		lang := newI18nLang(path, data)
		if utils.Contains(defaultLanguages, lang.name) {
			lang.isDefault = true
		}
		in.langs.Store(lang.name, lang)
		zlog.Infof("load i18n lang ok: %s", lang.localFile)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return in, nil
}

// Empty returns an empty i18n instance.
func Empty() *I18n {
	return &I18n{}
}

// LangCode returns correct language code possibly
// en-US -> [en-US,en-us,en]
func LangCode(lang string) []string {
	languages := []string{lang} // [en-US]
	lower := strings.ToLower(lang)
	if lower != lang {
		languages = append(languages, lower) // use lowercase language code, [en-us]
	}
	if strings.Contains(lang, "-") {
		languages = append(languages, strings.Split(lang, "-")[0]) // use first word if en-US, [en]
	}
	return languages
}
