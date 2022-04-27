package converter

import (
	"bytes"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

var (
	globalMarkdown goldmark.Markdown = nil
)

// GetMarkdown returns the global markdown instance.
func GetMarkdown() goldmark.Markdown {
	if globalMarkdown == nil {
		globalMarkdown = NewMarkdown()
	}
	return globalMarkdown
}

// NewMarkdown returns a new goldmark.Markdown instance.
func NewMarkdown() goldmark.Markdown {
	return goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithASTTransformers(
				util.Prioritized(newAstTransformer(), 10000),
			),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)
}

type astTransformer struct {
	linkProtocols []string
}

func newAstTransformer() *astTransformer {
	return &astTransformer{
		linkProtocols: []string{"http://", "https://", "//"},
	}
}

// Transform transforms the given AST tree.
func (g *astTransformer) Transform(node *ast.Document, reader text.Reader, pc parser.Context) {
	_ = ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		switch v := n.(type) {
		case *ast.Link:
			for _, p := range g.linkProtocols {
				if bytes.HasPrefix(v.Destination, []byte(p)) {
					v.SetAttributeString("target", []byte("_blank"))
					break
				}
			}
		}
		return ast.WalkContinue, nil
	})
}
