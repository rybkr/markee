package renderer

import (
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

func (r *HTMLRenderer) VisitDocument(node ast.Node) ast.VisitStatus {
    return ast.VisitChildrenDFS
}

func (r *HTMLRenderer) VisitBlockQuote(node ast.Node) ast.VisitStatus {
    r.output.WriteString("<blockquote>")
    return ast.VisitChildrenDFS
}

func (r *HTMLRenderer) VisitParagraph(node ast.Node) ast.VisitStatus {
    r.output.WriteString("<p>")
    return ast.VisitChildrenDFS
}

func (r *HTMLRenderer) VisitContent(node ast.Node) ast.VisitStatus {
    if content, ok := node.(*ast.Content); ok {
        r.output.WriteString(content.Literal)
        return ast.VisitChildrenDFS
    }
    return ast.VisitStop
}
