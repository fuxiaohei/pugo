package builder

import (
	"bytes"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

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
