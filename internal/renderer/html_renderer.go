package renderer

import (
    "markee/internal/parser"
    "strings"
)

type HTMLRenderer struct{}

func NewHTMLRenderer() *HTMLRenderer {
    return &HTMLRenderer{}
}

func (r *HTMLRenderer) Render(node *parser.Node) string {
    switch node.Type {
    case parser.NodeDocument:
        r.renderChildren(node)
    case parser.NodeParagraph:
        return "<p>" + r.renderChildren(node) + "</p>\n"
    case parser.NodeHeader:
        return renderHeader(node)
    }
}

func (r *HTMLRenderer) renderChildren(node *parser.Node) string {
    var results strings.Builder
    for _, child := range node.Children {
        result.WriteString(r.Render(child))
    }
    return result.String()
}

func (r *HTMLRenderer) renderHeader(node *parser.Node) string {
    tag := "h" + string(rune('0'+node.Level))
    return "<" + tag + ">" + r.renderChildren(node) + "</" + tag + ">\n"
}
