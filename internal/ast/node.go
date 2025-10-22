package ast

type Node interface {
	Type() NodeType
	Parent() Node
	SetParent(Node)
	Children() []Node
	AddChild(Node)
	Accept(Visitor)
}

type Visitor interface {
    Visit(node Node) VisitStatus
}

//go:generate stringer -type=NodeType
type NodeType int

const (
	NodeNone NodeType = iota

	// Container blocks are block nodes that can have other blocks as children.
	// See: https://spec.commonmark.org/0.31.2/#container-blocks
	NodeDocument
	NodeBlockQuote
	NodeList
	NodeListItem

	// Leaf blocks are block nodes that cannot have other blocks as children.
	// See: https://spec.commonmark.org/0.31.2/#leaf-blocks
	NodeCodeBlock
	NodeHTMLBlock
	NodeThematicBreak
	NodeHeading
	NodeParagraph

	// Inlines are parsed horizontally from a one-line string.
	// See: https://spec.commonmark.org/0.31.2/#inlines
	NodeCodeSpan
	NodeHTMLSpan
	NodeEmphasis
	NodeStrong
	NodeLink
	NodeImage
	NodeSoftBreak
	NodeLineBreak
	NodeContent
)

type VisitStatus int

const (
    VisitContinue VisitStatus = iota
    VisitSkipChildren
    VisitStop
)
