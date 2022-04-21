package model

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/yuin/goldmark"
	"gopkg.in/yaml.v3"
)

const (
	// DefaultPostTemplate is the default post template.
	DefaultPostTemplate = "post.html"
	// DefaultPostListTemplate is the default post list template.
	DefaultPostListTemplate = "post-list.html"
)

type postMetaSeperator struct {
	StartChars []byte
	EndChars   []byte
	MetaType   string
}

var (
	defaultPostMetaSeperator = []postMetaSeperator{
		{
			StartChars: []byte("---\n"),
			EndChars:   []byte("---\n"),
			MetaType:   "yaml",
		},
		{
			StartChars: []byte("```toml\n"),
			EndChars:   []byte("```\n"),
			MetaType:   "toml",
		},
	}
	defaultPostDateLayouts = []string{
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006-01-02",
	}
	defaultPostBriefSeperator = []byte("<!--more-->")
)

var (
	// ErrInvalidPostStartLine means the post start line is invalid.
	ErrInvalidPostStartLine = fmt.Errorf("invalid content start")
	// ErrInvalidPostDate means the post date is invalid, it must be format as postDefaultDateLayout
	ErrInvalidPostDate = fmt.Errorf("invalid content date, it must be format as " + strings.Join(defaultPostDateLayouts, " or "))
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
	AuthorName   string   `toml:"author" yaml:"author"`

	Author   *Author   `toml:"-" yaml:"-"`
	Link     string    `toml:"-" yaml:"-"`
	TagLinks []TagLink `toml:"-" yaml:"-"`

	localFile   string
	rawContent  []byte
	htmlContent string
	rawBrief    []byte
	htmlBrief   string
	dateTime    time.Time
}

// TagLink returns the tag link of the post.
type TagLink struct {
	Name string
	Link string
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
	for _, seperator := range defaultPostMetaSeperator {
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
	return nil, ErrInvalidPostStartLine
}

func (p *Post) parseDate() error {
	// if date is empty, use file modified time
	if p.DateString == "" {
		info, _ := os.Stat(p.localFile)
		p.DateString = info.ModTime().Format(defaultPostDateLayouts[0])
	}
	for _, layout := range defaultPostDateLayouts {
		dt, err := time.Parse(layout, p.DateString)
		if err != nil {
			continue
		}
		p.dateTime = dt
		return nil
	}
	return ErrInvalidPostDate
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
	if bytes.Contains(p.rawContent, defaultPostBriefSeperator) {
		p.rawBrief = bytes.SplitN(p.rawContent, defaultPostBriefSeperator, 2)[0]
	}

	return nil
}

// Convert converts post markdown content to html content.
func (p *Post) Convert(md goldmark.Markdown) error {
	buf := bytes.NewBuffer(nil)
	if len(p.rawBrief) > 0 {
		if err := md.Convert(p.rawBrief, buf); err != nil {
			return err
		}
		p.htmlBrief = buf.String()
		buf.Reset()
	}
	if err := md.Convert(p.rawContent, buf); err != nil {
		return err
	}
	p.htmlContent = buf.String()
	return nil
}

// NewPostFromFile returns a new post from file.
func NewPostFromFile(path string) (*Post, error) {
	rawData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	// trim space lines
	rawData = bytes.TrimSpace(rawData)
	p := new(Post)
	if err = p.Parse(path, rawData); err != nil {
		return nil, err
	}

	// fix slug empty
	// use title as slug if slug is empty
	if p.Slug == "" {
		p.Slug = url.PathEscape(p.Title)
	}

	// fix empty template
	if p.Template == "" {
		p.Template = DefaultPostTemplate
	}

	return p, nil
}
