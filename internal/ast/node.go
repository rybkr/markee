package ast

type Node interface {
	Type() NodeType
	Parent() Node
	SetParent(Node)
	Children() []Node
	AddChild(Node)
    IsOpen() bool
    SetOpen(bool)
	Accept(Visitor) VisitStatus
}

type BaseNode struct {
	nodeType NodeType
	parent   Node
	children []Node
    isOpen   bool
}

func New(t NodeType) BaseNode {
	return BaseNode{
		nodeType: t,
		parent:   nil,
		children: make([]Node, 0),
        isOpen:   true,
	}
}

func (n *BaseNode) Type() NodeType {
	return n.nodeType
}

func (n *BaseNode) Parent() Node {
	return n.parent
}

func (n *BaseNode) SetParent(parent Node) {
	n.parent = parent
}

func (n *BaseNode) Children() []Node {
	return n.children
}

func (n *BaseNode) AddChild(child Node) {
	child.SetParent(n)
	n.children = append(n.children, child)
}

func (n *BaseNode) IsOpen() bool {
    return n.isOpen
}

func (n *BaseNode) SetOpen(isOpen bool) {
    n.isOpen = isOpen
}

func (n *BaseNode) Accept(v Visitor) VisitStatus {
	var status VisitStatus = VisitStop
	switch n.Type() {
	case NodeDocument:
		status = v.VisitDocument(n)
	case NodeBlockQuote:
		status = v.VisitBlockQuote(n)
	case NodeList:
		status = v.VisitList(n)
	case NodeListItem:
		status = v.VisitListItem(n)
	case NodeCodeBlock:
		status = v.VisitCodeBlock(n)
	case NodeHTMLBlock:
		status = v.VisitHTMLBlock(n)
	case NodeThematicBreak:
		status = v.VisitThematicBreak(n)
	case NodeHeading:
		status = v.VisitHeading(n)
	case NodeParagraph:
		status = v.VisitParagraph(n)
	case NodeCodeSpan:
		status = v.VisitCodeSpan(n)
	case NodeHTMLSpan:
		status = v.VisitHTMLSpan(n)
	case NodeEmphasis:
		status = v.VisitEmphasis(n)
	case NodeStrong:
		status = v.VisitStrong(n)
	case NodeLink:
		status = v.VisitLink(n)
	case NodeImage:
		status = v.VisitImage(n)
	case NodeSoftBreak:
		status = v.VisitSoftBreak(n)
	case NodeLineBreak:
		status = v.VisitLineBreak(n)
	case NodeContent:
		status = v.VisitContent(n)
	}

	switch status {
	case VisitStop:
		return VisitStop
	case VisitSkipChildren:
		return VisitContinue
	case VisitContinue:
		for _, child := range n.Children() {
			if child.Accept(v) == VisitStop {
				return VisitStop
			}
		}
		return VisitContinue
    default:
        return VisitStop
	}
}

type NodeType int

const (
	// Container blocks are block nodes that can have other blocks as children.
	// See: https://spec.commonmark.org/0.31.2/#container-blocks
	NodeDocument NodeType = iota
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

const (
	nodeContainerStart = NodeDocument
	nodeContainerEnd   = NodeListItem
	nodeLeafStart      = NodeCodeBlock
	nodeLeafEnd        = NodeParagraph
	nodeBlockStart     = NodeDocument
	nodeBlockEnd       = NodeParagraph
	nodeInlineStart    = NodeCodeSpan
	nodeInlineEnd      = NodeContent
)
