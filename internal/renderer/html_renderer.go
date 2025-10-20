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
        return r.renderChildren(node)
    case parser.NodeParagraph:
        return "<p>" + r.renderChildren(node) + "</p>\n"
    case parser.NodeHeader:
        return r.renderHeader(node)
    case parser.NodeBlockquote:
        return "<blockquote>" + r.renderChildren(node) + "</blockquote>\n"
    case parser.NodeCodeBlock:
        return "<code>" + node.Value + "</code>\n"
    default:
        return node.Value
    }
}

func (r *HTMLRenderer) renderChildren(node *parser.Node) string {
    var result strings.Builder
    for _, child := range node.Children {
        result.WriteString(r.Render(child))
    }
    return result.String()
}

func (r *HTMLRenderer) renderHeader(node *parser.Node) string {
    tag := "h" + string(rune('0'+node.Level))
    return "<" + tag + ">" + r.renderChildren(node) + "</" + tag + ">\n"
}
