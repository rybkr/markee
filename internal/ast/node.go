package ast

type Node interface {
	Type() NodeType
	IsLeaf() bool
	IsContainer() bool
	IsBlock() bool
	IsInline() bool
	Parent() Node
	SetParent(Node)
	Children() []Node
    LastChild() Node
	AddChild(Node)
	IsOpen() bool
	SetOpen(bool)
	Accept(Visitor)
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
		isOpen:   false,
	}
}

func (n *BaseNode) Type() NodeType {
	return n.nodeType
}

func (n *BaseNode) IsLeaf() bool {
	return n.Type() >= nodeLeafStart && n.Type() <= nodeLeafEnd
}

func (n *BaseNode) IsContainer() bool {
	return n.Type() >= nodeContainerStart && n.Type() <= nodeContainerEnd
}

func (n *BaseNode) IsBlock() bool {
	return n.Type() >= nodeBlockStart && n.Type() <= nodeBlockEnd
}

func (n *BaseNode) IsInline() bool {
	return n.Type() >= nodeInlineStart && n.Type() <= nodeInlineEnd
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

func (n *BaseNode) LastChild() Node {
    if children := n.Children(); len(children) != 0 {
        return children[len(children)-1]
    }
    return nil
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

func (n *BaseNode) Accept(v Visitor) { return }

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
    NodeBlankLine

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
	nodeLeafEnd        = NodeBlankLine
	nodeBlockStart     = NodeDocument
	nodeBlockEnd       = NodeBlankLine
	nodeInlineStart    = NodeCodeSpan
	nodeInlineEnd      = NodeContent
)
