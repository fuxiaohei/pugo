package markdown

import (
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
)

var (
	globalMinifier *minify.M = nil
)

// InitMinfier initializes minifier
func InitMinfier(flag bool) {
	if flag {
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
		globalMinifier = m
	} else {
		globalMinifier = nil
	}
}

// MinifyHTML minifies HTML
func MinifyHTML(raw []byte) ([]byte, error) {
	if globalMinifier == nil {
		return raw, nil
	}
	return globalMinifier.Bytes("text/html", raw)
}
