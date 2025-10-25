package ast

type Node interface {
	Type() NodeType
	Parent() Node
	SetParent(Node)
	NextSibling() Node
	SetSibling(Node)
	FirstChild() Node
    SetFirstChild(Node)
	LastChild() Node
    SetLastChild(Node)
    Children() []Node
	AddChild(Node)
	Accept(Visitor)
}

type BaseNode struct {
	parent      Node
	firstChild  Node
	lastChild   Node
	nextSibling Node
	nodeType    NodeType
}

func New(t NodeType) BaseNode {
	return BaseNode{
		parent:      nil,
		firstChild:  nil,
		lastChild:   nil,
		nextSibling: nil,
		nodeType:    t,
	}
}

func (n *BaseNode) Type() NodeType {
	return n.nodeType
}

func (n *BaseNode) Parent() Node {
	return n.parent
}

func (n *BaseNode) SetParent(node Node) {
	n.parent = node
}

func (n *BaseNode) NextSibling() Node {
	return n.nextSibling
}

func (n *BaseNode) SetSibling(sibling Node) {
    sibling.SetParent(n.Parent())
    if nextSibling := n.NextSibling(); nextSibling != nil {
        sibling.SetSibling(nextSibling)
    } else {
        n.Parent().SetLastChild(sibling)
    }
	n.nextSibling = sibling
}

func (n *BaseNode) FirstChild() Node {
	return n.firstChild
}

func (n *BaseNode) SetFirstChild(child Node) {
    n.firstChild = child
}

func (n *BaseNode) LastChild() Node {
	return n.lastChild
}

func (n *BaseNode) SetLastChild(child Node) {
    n.lastChild = child
}

func (n *BaseNode) Children() []Node {
    children := make([]Node, 0)
    for child := n.FirstChild(); child != nil; child = child.NextSibling() {
        children = append(children, child)
    }
    return children
}

func (n *BaseNode) AddChild(child Node) {
	child.SetParent(n)
	if lastChild := n.LastChild(); lastChild != nil {
		lastChild.SetSibling(child)
	} else {
        n.SetLastChild(child)
    }
    if firstChild := n.FirstChild(); firstChild == nil {
        n.SetFirstChild(child)
    }
}

func (n *BaseNode) Accept(v Visitor) {
    // Default no-op
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

func (t NodeType) IsLeaf() bool {
	return t >= NodeCodeBlock && t <= NodeBlankLine
}

func (t NodeType) IsContainer() bool {
	return t >= NodeDocument && t <= NodeListItem
}

func (t NodeType) IsBlock() bool {
	return t.IsLeaf() || t.IsContainer()
}

func (t NodeType) IsInline() bool {
	return !t.IsBlock()
}
