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
)

var (
	postStartPrefix       = []byte("```toml\n")
	postBriefSeperator    = "<!--more-->"
	postDefaultDateLayout = "2006-01-02 15:04"
)

var (
	// ErrInvalidPostStartLine means the post start line is invalid.
	ErrInvalidPostStartLine = fmt.Errorf("invalid content start")
	// ErrInvalidPostDate means the post date is invalid, it must be format as postDefaultDateLayout
	ErrInvalidPostDate = fmt.Errorf("invalid content date, it must be format as " + postDefaultDateLayout)
)

// Post is the definition of a post.
type Post struct {
	Title        string   `toml:"title"`
	Slug         string   `toml:"slug"`
	Descripition string   `toml:"descripition"`
	Tags         []string `toml:"tags"`
	Date         string   `toml:"date"`
	Template     string   `toml:"template"`
	Draft        bool     `toml:"draft"`

	localFile    string
	rawContent   string
	briefContent string
}

func (p *Post) Content() string {
	return p.rawContent
}

func (p *Post) Brief() string {
	if p.briefContent == "" {
		// if brief content is empty, use raw content
		return p.rawContent
	}
	return p.briefContent
}

func (p *Post) Parse(file string, rawData []byte) error {
	p.localFile = file

	// parse post basic info via toml code section
	if !bytes.HasPrefix(rawData, postStartPrefix) {
		return ErrInvalidPostStartLine
	}
	rawData = bytes.TrimPrefix(rawData, postStartPrefix)
	rawDataSlice := bytes.SplitN(rawData, []byte("```"), 2)
	if err := toml.Unmarshal(bytes.TrimSpace(rawDataSlice[0]), p); err != nil {
		return err
	}

	// if date is empty, use file modified time
	if p.Date == "" {
		info, _ := os.Stat(p.localFile)
		p.Date = info.ModTime().Format(postDefaultDateLayout)
	}
	if _, err := time.Parse(postDefaultDateLayout, p.Date); err != nil {
		return ErrInvalidPostDate
	}

	p.rawContent = string(bytes.TrimSpace(rawDataSlice[1]))

	// parse brief content
	if strings.Contains(p.rawContent, postBriefSeperator) {
		p.briefContent = strings.SplitN(p.rawContent, postBriefSeperator, 2)[0]
	}

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

	return p, nil
}
