package models

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"pugo/pkg/core/constants"
	"pugo/pkg/ext/markdown"
	"pugo/pkg/utils/zlog"
	"sort"
	"time"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

// Post is the definition of a post.
type Post struct {
	Title        string   `toml:"title" yaml:"title"`
	Slug         string   `toml:"slug" yaml:"slug"`
	Descripition string   `toml:"description" yaml:"description"`
	Tags         []string `toml:"tags" yaml:"tags"`
	DateString   string   `toml:"date" yaml:"date"`
	Template     string   `toml:"template" yaml:"template"`
	Draft        bool     `toml:"draft" yaml:"draft"`
	Comment      bool     `toml:"comment" yaml:"comment"`
	AuthorName   string   `toml:"author" yaml:"author"`
	Language     string   `toml:"language" yaml:"language"`

	Author   *Author    `toml:"-" yaml:"-"`
	Link     string     `toml:"-" yaml:"-"`
	TagLinks []*TagLink `toml:"-" yaml:"-"`

	localFile   string
	rawContent  []byte
	htmlContent string
	rawBrief    []byte
	htmlBrief   string
	dateTime    time.Time
}

// NewPostFromFile returns a new post from file.
func NewPostFromFile(path string) (*Post, error) {
	p, err := parseContentBase(path)
	if err != nil {
		return nil, err
	}

	// fix slug empty
	// use title as slug if slug is empty
	if p.Slug == "" {
		p.Slug = url.PathEscape(p.Title)
	}

	// fix empty template
	if p.Template == "" {
		p.Template = constants.PostTemplate
	}

	return p, nil
}

func parseContentBase(path string) (*Post, error) {
	rawData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	// trim space lines
	rawData = bytes.TrimSpace(rawData)
	p := &Post{
		Draft:   false,
		Comment: true,
	}
	if err = p.Parse(path, rawData); err != nil {
		return nil, err
	}
	return p, nil
}

func (p *Post) Parse(file string, rawData []byte) error {
	p.localFile = file

	// parse post meta info via code block section
	contentBytes, err := p.parseMeta(rawData)
	if err != nil {
		return err
	}
	p.rawContent = bytes.TrimSpace(contentBytes)

	// parse date time
	if err := p.parseDate(); err != nil {
		return err
	}

	// parse brief content
	seperator := constants.PostBriefSeperator()
	if bytes.Contains(p.rawContent, seperator) {
		p.rawBrief = bytes.SplitN(p.rawContent, seperator, 2)[0]
	}

	return nil
}

func (p *Post) Content() string {
	return p.htmlContent
}

func (p *Post) Brief() string {
	if p.htmlBrief == "" {
		// if brief content is empty, use raw content
		return p.htmlContent
	}
	return p.htmlBrief
}

// Date returns the post date as time.Time
func (p *Post) Date() time.Time {
	return p.dateTime
}

// LocalFile returns the local file path of the post.
func (p *Post) LocalFile() string {
	return p.localFile
}

func (p *Post) parseMeta(rawData []byte) ([]byte, error) {
	separators := constants.PostMetaSeperators()
	for _, seperator := range separators {
		// find prefix
		if !bytes.HasPrefix(rawData, seperator.StartChars) {
			continue
		}
		rawData = bytes.TrimPrefix(rawData, seperator.StartChars)
		rawDataSlice := bytes.SplitN(rawData, seperator.EndChars, 2)

		// parse as toml
		if seperator.MetaType == "toml" {
			if err := toml.Unmarshal(bytes.TrimSpace(rawDataSlice[0]), p); err != nil {
				return nil, err
			}
			return bytes.TrimSpace(rawDataSlice[1]), nil
		}

		// parse as yaml
		if seperator.MetaType == "yaml" {
			if err := yaml.Unmarshal(bytes.TrimSpace(rawDataSlice[0]), p); err != nil {
				return nil, err
			}
			return rawDataSlice[1], nil
		}
	}
	return nil, constants.ErrInvalidPostStartLine
}

func (p *Post) parseDate() error {
	dateLayouts := constants.PostDateLayouts()
	// if date is empty, use file modified time
	if p.DateString == "" {
		info, _ := os.Stat(p.localFile)
		p.DateString = info.ModTime().Format(dateLayouts[0])
	}
	for _, layout := range dateLayouts {
		dt, err := time.Parse(layout, p.DateString)
		if err != nil {
			continue
		}
		p.dateTime = dt
		return nil
	}
	return constants.ErrInvalidPostDate
}

// Convert converts post markdown content to html content.
func (p *Post) Convert(fn markdown.ConvertFunc) error {
	if fn == nil {
		return errors.New("converter is nil")
	}
	buf := bytes.NewBuffer(nil)
	if len(p.rawBrief) > 0 {
		if err := fn(p.rawBrief, buf); err != nil {
			return err
		}
		p.htmlBrief = buf.String()
		buf.Reset()
	}
	if err := fn(p.rawContent, buf); err != nil {
		return err
	}
	p.htmlContent = buf.String()
	return nil
}

// LoadPosts loads posts from content/posts directory.
func LoadPosts(withDrafts bool) ([]*Post, error) {
	var posts []*Post
	err := filepath.Walk(constants.ContentPostsDir, func(path string, info os.FileInfo, err error) error {
		// skip directory
		if info.IsDir() {
			return nil
		}

		// only process markdown files
		if filepath.Ext(path) != ".md" {
			return nil
		}

		post, err := NewPostFromFile(path)
		if err != nil {
			zlog.Warnf("failed to load post: %s, %s", path, err)
			return nil
		}
		if post.Draft && !withDrafts {
			zlog.Warnf("skip draft post: %s", path)
			return nil
		}

		// save post into parsed data
		posts = append(posts, post)
		zlog.Infof("load post ok: %s", path)

		return nil
	})

	if err != nil {
		return nil, err
	}
	sort.Slice(posts, func(i, j int) bool {
		// order by date desc
		return posts[i].Date().Unix() > posts[j].Date().Unix()
	})
	return posts, nil
}
