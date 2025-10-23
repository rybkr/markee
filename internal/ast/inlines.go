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

type Emphasis struct{ BaseNode }

func NewEmphasis() *Emphasis {
	return &Emphasis{
		BaseNode: New(NodeEmphasis),
	}
}

type Strong struct{ BaseNode }

func NewStrong() *Strong {
	return &Strong{
		BaseNode: New(NodeStrong),
	}
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

type SoftBreak struct{ BaseNode }

func NewSoftBreak() *SoftBreak {
	return &SoftBreak{
		BaseNode: New(NodeSoftBreak),
	}
}

type LineBreak struct{ BaseNode }

func NewLineBreak() *LineBreak {
	return &LineBreak{
		BaseNode: New(NodeLineBreak),
	}
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
