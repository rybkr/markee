package renderer

import (
    "markee/internal/parser"
)

type Renderer interface {
    Render(node *parser.Node) string

    renderChildren(node *parser.Node) string
}

func RenderHTML(document *parser.Node) string {
    return NewHTMLRenderer().Render(document)
}
