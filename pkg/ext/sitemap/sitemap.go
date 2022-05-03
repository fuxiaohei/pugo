/*
Copyright © 2016-2017 Janne Snabb snabb AT epipe.com

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
"Software"), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be included
in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/
package sitemap

import (
	"bytes"
	"encoding/xml"
	"io"
	"path/filepath"
	"pugo/pkg/core/models"
	"pugo/pkg/utils/zlog"
	"strings"
	"time"
)

// copy from https://github.com/snabb/sitemap/blob/master/sitemap.go

// ChangeFreq specifies change frequency of a sitemap entry. It is just a string.
type ChangeFreq string

// Feel free to use these constants for ChangeFreq (or you can just supply
// a string directly).
const (
	Always  ChangeFreq = "always"
	Hourly  ChangeFreq = "hourly"
	Daily   ChangeFreq = "daily"
	Weekly  ChangeFreq = "weekly"
	Monthly ChangeFreq = "monthly"
	Yearly  ChangeFreq = "yearly"
	Never   ChangeFreq = "never"
)

// URL entry in sitemap or sitemap index. LastMod is a pointer
// to time.Time because omitempty does not work otherwise. Loc is the
// only mandatory item. ChangeFreq and Priority must be left empty when
// using with a sitemap index.
type URL struct {
	Loc        string     `xml:"loc"`
	LastMod    *time.Time `xml:"lastmod,omitempty"`
	ChangeFreq ChangeFreq `xml:"changefreq,omitempty"`
	Priority   float32    `xml:"priority,omitempty"`
}

// Sitemap represents a complete sitemap which can be marshaled to XML.
// New instances must be created with New() in order to set the xmlns
// attribute correctly. Minify can be set to make the output less human
// readable.
type Sitemap struct {
	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr"`
	URLs    []*URL   `xml:"url"`
	baseURL string
}

var (
	globalSiteMap *Sitemap = nil
)

// Init initializes the global sitemap.
func Init(cfg *Config, baseURL string) {
	if cfg != nil && cfg.Enabled {
		globalSiteMap = NewSiteMap(baseURL)
		return
	}
	globalSiteMap = nil
	return
}

// Add adds an URL to the global sitemap.
func Add(u *URL) {
	if globalSiteMap != nil {
		globalSiteMap.Add(u)
	}
}

// Write writes XML encoded sitemap to given io.Writer.
func Write(w io.Writer) error {
	if globalSiteMap == nil {
		return nil
	}
	return globalSiteMap.Write(w)
}

// NewSiteMap returns a new Sitemap.
func NewSiteMap(baseURL string) *Sitemap {
	return &Sitemap{
		Xmlns:   "http://www.sitemaps.org/schemas/sitemap/0.9",
		baseURL: baseURL,
	}
}

func (s *Sitemap) fullLoc(loc string) string {
	return strings.TrimSuffix(s.baseURL, "/") + "/" + strings.TrimPrefix(loc, "/")
}

// Add adds an URL to a Sitemap.
func (s *Sitemap) Add(u *URL) {
	u.Loc = s.fullLoc(u.Loc)
	s.URLs = append(s.URLs, u)
}

// WriteTo writes XML encoded sitemap to given io.Writer.
func (s *Sitemap) Write(w io.Writer) error {
	_, err := w.Write([]byte(xml.Header))
	if err != nil {
		return err
	}
	en := xml.NewEncoder(w)
	if err = en.Encode(s); err != nil {
		return err
	}
	_, err = w.Write([]byte{'\n'})
	return nil
}

// Render renders the sitemap to a string.
func Render(cfg *Config, outputDir string) (*models.OutputFile, error) {
	if cfg == nil || !cfg.Enabled {
		return nil, nil
	}
	buf := bytes.NewBuffer(nil)
	if err := Write(buf); err != nil {
		zlog.Warnf("failed to marshal sitemap: %s", err)
		return nil, err
	}
	dstFile := filepath.Join(outputDir, cfg.Link)
	return &models.OutputFile{
		Path: dstFile,
		Buf:  buf,
	}, nil
}
