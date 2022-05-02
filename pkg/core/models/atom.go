// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Adapted from encoding/xml/read_test.go.
// Package atom defines XML data structures for an Atom feed.
// Modified from "golang.org/x/tools/blog/atom"
package models

import (
	"encoding/xml"
	"time"
)

type (
	AtomFeed struct {
		XMLName xml.Name     `xml:"http://www.w3.org/2005/Atom feed"`
		Title   string       `xml:"title"`
		ID      string       `xml:"id"`
		Link    []AtomLink   `xml:"link"`
		Updated AtomTimeStr  `xml:"updated"`
		Author  *AtomPerson  `xml:"author"`
		Entry   []*AtomEntry `xml:"entry"`
	}
	AtomEntry struct {
		Title     string      `xml:"title"`
		ID        string      `xml:"id"`
		Link      []AtomLink  `xml:"link"`
		Published AtomTimeStr `xml:"published"`
		Updated   AtomTimeStr `xml:"updated"`
		Author    *AtomPerson `xml:"author"`
		Summary   *AtomText   `xml:"summary"`
		Content   *AtomText   `xml:"content"`
	}
	AtomLink struct {
		Rel      string `xml:"rel,attr,omitempty"`
		Href     string `xml:"href,attr"`
		Type     string `xml:"type,attr,omitempty"`
		HrefLang string `xml:"hreflang,attr,omitempty"`
		Title    string `xml:"title,attr,omitempty"`
		Length   uint   `xml:"length,attr,omitempty"`
	}
	AtomPerson struct {
		Name     string `xml:"name"`
		URI      string `xml:"uri,omitempty"`
		Email    string `xml:"email,omitempty"`
		InnerXML string `xml:",innerxml"`
	}
	AtomText struct {
		Type string `xml:"type,attr"`
		Body string `xml:",chardata"`
	}
)

type AtomTimeStr string

func AtomTime(t time.Time) AtomTimeStr {
	return AtomTimeStr(t.Format("2006-01-02T15:04:05-07:00"))
}

// BuildAtom builds an Atom feed from the given source.
func BuildAtom(link string, posts []*Post, sc *SiteConfig) *AtomFeed {
	feed := AtomFeed{
		Title: sc.Title,
		// ID:      "tag:" + s.cfg.Hostname + ",2013:" + s.cfg.Hostname,
		Updated: AtomTime(posts[0].Date()),
		Link: []AtomLink{{
			Rel:  "self",
			Href: sc.FullURL(link),
		}},
	}
	for _, p := range posts {
		feed.Entry = append(feed.Entry, &AtomEntry{
			Title: p.Title,
			Link: []AtomLink{{
				Rel:  "alternate",
				Href: sc.FullURL(p.Link),
			}},
			Published: AtomTime(p.Date()),
			Updated:   AtomTime(p.Date()),
			Summary: &AtomText{
				Type: "html",
				Body: p.Brief(),
			},
			Content: &AtomText{
				Type: "html",
				Body: p.Content(),
			},
			Author: &AtomPerson{
				Name: p.Author.Name,
			},
		})
	}
	return &feed
}
