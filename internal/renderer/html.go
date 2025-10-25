package renderer

import (
    "fmt"
    "markee/internal/ast"
    "strings"
)

type HTMLRenderer struct {
    Renderer
    ast.BaseVisitor
    output strings.Builder
}

func NewHTMLRenderer() *HTMLRenderer {
    return &HTMLRenderer{}
}

func RenderHTML(doc *ast.Document) string {
    return NewHTMLRenderer().Render(doc)
}

func (r *HTMLRenderer) Render(doc *ast.Document) string {
    doc.Accept(r)
    return r.output.String()
}

func (r *HTMLRenderer) VisitDocument(node ast.Node) {
    ast.WalkChildren(r, node)
}

func (r *HTMLRenderer) VisitBlockQuote(node ast.Node) {
    r.output.WriteString("<blockquote>\n")
    ast.WalkChildren(r, node)
    r.output.WriteString("</blockquote>\n")
}

func (r *HTMLRenderer) VisitCodeBlock(node ast.Node) {
    codeBlock := node.(*ast.CodeBlock)
    
    r.output.WriteString("<pre><code")
    if codeBlock.Language != "" {
        r.output.WriteString(fmt.Sprintf(" class=\"language-%s\"", codeBlock.Language))
    }
    r.output.WriteString(">")
    
    for child := node.FirstChild(); child != nil; child = child.NextSibling() {
        if content, ok := child.(*ast.Content); ok {
            r.output.WriteString(escapeHTML(content.Literal))
            if child.NextSibling() != nil {
                r.output.WriteString("\n")
            }
        }
    }
    
    r.output.WriteString("</code></pre>\n")
}

func (r *HTMLRenderer) VisitHeading(node ast.Node) {
    heading := node.(*ast.Heading)
    r.output.WriteString(fmt.Sprintf("<h%d>", heading.Level))
    ast.WalkChildren(r, node)
    r.output.WriteString(fmt.Sprintf("</h%d>\n", heading.Level))
}

func (r *HTMLRenderer) VisitParagraph(node ast.Node) {
    r.output.WriteString("<p>")
    ast.WalkChildren(r, node)
    r.output.WriteString("</p>\n")
}

func (r *HTMLRenderer) VisitThematicBreak(node ast.Node) {
    r.output.WriteString("<hr />\n")
}

func (r *HTMLRenderer) VisitList(node ast.Node) {
    list := node.(*ast.List)
    tag := "ul"
    if list.IsOrdered {
        tag = "ol"
    }
    
    r.output.WriteString(fmt.Sprintf("<%s>\n", tag))
    ast.WalkChildren(r, node)
    r.output.WriteString(fmt.Sprintf("</%s>\n", tag))
}

func (r *HTMLRenderer) VisitListItem(node ast.Node) {
    r.output.WriteString("<li>")
    ast.WalkChildren(r, node)
    r.output.WriteString("</li>\n")
}

func (r *HTMLRenderer) VisitText(node ast.Node) {
    if text, ok := node.(*ast.Content); ok {
        r.output.WriteString(escapeHTML(text.Literal))
    }
}

func (r *HTMLRenderer) VisitStrong(node ast.Node) {
    r.output.WriteString("<strong>")
    ast.WalkChildren(r, node)
    r.output.WriteString("</strong>")
}

func (r *HTMLRenderer) VisitEmph(node ast.Node) {
    r.output.WriteString("<em>")
    ast.WalkChildren(r, node)
    r.output.WriteString("</em>")
}

func (r *HTMLRenderer) VisitCodeSpan(node ast.Node) {
    if code, ok := node.(*ast.CodeSpan); ok {
        r.output.WriteString("<code>")
        r.output.WriteString(escapeHTML(code.Literal))
        r.output.WriteString("</code>")
    }
}

func (r *HTMLRenderer) VisitLink(node ast.Node) {
    if link, ok := node.(*ast.Link); ok {
        r.output.WriteString(fmt.Sprintf("<a href=\"%s\"", escapeAttribute(link.Destination)))
        if link.Title != "" {
            r.output.WriteString(fmt.Sprintf(" title=\"%s\"", escapeAttribute(link.Title)))
        }
        r.output.WriteString(">")
        ast.WalkChildren(r, node)
        r.output.WriteString("</a>")
    }
}

func (r *HTMLRenderer) VisitImage(node ast.Node) {
    if image, ok := node.(*ast.Image); ok {
        r.output.WriteString(fmt.Sprintf("<img src=\"%s\" alt=\"%s\"",
            escapeAttribute(image.Destination),
            escapeAttribute(image.AltText)))
        if image.Title != "" {
            r.output.WriteString(fmt.Sprintf(" title=\"%s\"", escapeAttribute(image.Title)))
        }
        r.output.WriteString(" />")
    }
}

func (r *HTMLRenderer) VisitSoftBreak(node ast.Node) {
    r.output.WriteString("\n")
}

func (r *HTMLRenderer) VisitHardBreak(node ast.Node) {
    r.output.WriteString("<br />\n")
}

func (r *HTMLRenderer) VisitContent(node ast.Node) {
    if content, ok := node.(*ast.Content); ok {
        r.output.WriteString(escapeHTML(content.Literal))
    }
}

func escapeHTML(s string) string {
    s = strings.ReplaceAll(s, "&", "&amp;")
    s = strings.ReplaceAll(s, "<", "&lt;")
    s = strings.ReplaceAll(s, ">", "&gt;")
    s = strings.ReplaceAll(s, "\"", "&quot;")
    return s
}

func escapeAttribute(s string) string {
    s = strings.ReplaceAll(s, "&", "&amp;")
    s = strings.ReplaceAll(s, "\"", "&quot;")
    return s
}
