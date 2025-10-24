package ast

type CodeSpan struct {
	BaseNode
	Literal string
}

func NewCodeSpan(literal string) *CodeSpan {
	return &CodeSpan{
		BaseNode: New(NodeCodeSpan),
		Literal:  literal,
	}
}

func (c *CodeSpan) Accept(v Visitor) {
    v.VisitCodeSpan(c)
}

type HTMLSpan struct {
	BaseNode
	Literal string
}

func NewHTMLSpan(literal string) *HTMLSpan {
	return &HTMLSpan{
		BaseNode: New(NodeHTMLSpan),
		Literal:  literal,
	}
}

func (h *HTMLSpan) Accept(v Visitor) {
    v.VisitHTMLSpan(h)
}

type Emphasis struct{ BaseNode }

func NewEmphasis() *Emphasis {
	return &Emphasis{
		BaseNode: New(NodeEmphasis),
	}
}

func (e *Emphasis) Accept(v Visitor) {
    v.VisitEmphasis(e)
}

type Strong struct{ BaseNode }

func NewStrong() *Strong {
	return &Strong{
		BaseNode: New(NodeStrong),
	}
}

func (s *Strong) Accept(v Visitor) {
    v.VisitStrong(s)
}

type Link struct {
	BaseNode
	Destination string
	Title       string
}

func NewLink(destination, title string) *Link {
	return &Link{
		BaseNode:    New(NodeLink),
		Destination: destination,
		Title:       title,
	}
}

func (l *Link) Accept(v Visitor) {
    v.VisitLink(l)
}

type Image struct {
	BaseNode
	Destination string
	Title       string
	AltText     string
}

func NewImage(destination, title, altText string) *Image {
	return &Image{
		BaseNode:    New(NodeImage),
		Destination: destination,
		Title:       title,
		AltText:     altText,
	}
}

func (i *Image) Accept(v Visitor) {
    v.VisitImage(i)
}

type SoftBreak struct{ BaseNode }

func NewSoftBreak() *SoftBreak {
	return &SoftBreak{
		BaseNode: New(NodeSoftBreak),
	}
}

func (s *SoftBreak) Accept(v Visitor) {
    v.VisitSoftBreak(s)
}

type LineBreak struct{ BaseNode }

func NewLineBreak() *LineBreak {
	return &LineBreak{
		BaseNode: New(NodeLineBreak),
	}
}

func (l *LineBreak) Accept(v Visitor) {
    v.VisitLineBreak(l)
}

type Content struct {
	BaseNode
	Literal string
}

func NewContent(literal string) *Content {
	return &Content{
		BaseNode: New(NodeContent),
		Literal:  literal,
	}
}

func (c *Content) Accept(v Visitor) {
    v.VisitContent(c)
}
