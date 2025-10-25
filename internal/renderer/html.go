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
    for _, child := range node.Children() {
        child.Accept(r)
    }
}

func (r *HTMLRenderer) VisitBlockQuote(node ast.Node) {
    r.output.WriteString("<blockquote>\n")
    for _, child := range node.Children() {
        child.Accept(r)
    }
    r.output.WriteString("</blockquote>\n")
}

func (r *HTMLRenderer) VisitCodeBlock(node ast.Node) {
    r.output.WriteString("<pre><code>")
    for _, child := range node.Children() {
        child.Accept(r)
    }
    r.output.WriteString("\n</code></pre>\n")
}

func (r *HTMLRenderer) VisitHeading(node ast.Node) {
    if heading, ok := node.(*ast.Heading); ok {
        r.output.WriteString(fmt.Sprintf("<h%d>", heading.Level))
        for _, child := range node.Children() {
            child.Accept(r)
        }
        r.output.WriteString(fmt.Sprintf("</h%d>\n", heading.Level))
    }
}

func (r *HTMLRenderer) VisitParagraph(node ast.Node) {
    r.output.WriteString("<p>")
    for _, child := range node.Children() {
        child.Accept(r)
    }
    r.output.WriteString("</p>\n")
}

func (r *HTMLRenderer) VisitContent(node ast.Node) {
    if content, ok := node.(*ast.Content); ok {
        r.output.WriteString(content.Literal)
        for _, child := range node.Children() {
            child.Accept(r)
        }
    }
}

func (r *HTMLRenderer) VisitThematicBreak(node ast.Node) {
    r.output.WriteString("<hr />\n")
}
