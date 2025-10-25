package ast

type Node interface {
	Type() NodeType
	Parent() Node
	SetParent(Node)

	PrevSibling() Node
	NextSibling() Node
    setPrevSibling(Node)
    setNextSibling(Node)

	FirstChild() Node
	LastChild() Node
	Children() []Node
	AddChild(Node)
	RemoveChild(Node)
	InsertAfter(oldNode, newNode Node)
	ReplaceChild(oldNode, newNode Node)

	IsOpen() bool
	SetOpen(bool)

	Accept(Visitor)
}

type BaseNode struct {
	parent      Node
	firstChild  Node
	lastChild   Node
	prevSibling Node
	nextSibling Node
	nodeType    NodeType
	isOpen      bool
}

func New(t NodeType) BaseNode {
	return BaseNode{
		parent:      nil,
		firstChild:  nil,
		lastChild:   nil,
		nextSibling: nil,
		prevSibling: nil,
		nodeType:    t,
		isOpen:      true,
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

func (n *BaseNode) PrevSibling() Node {
	return n.prevSibling
}

func (n *BaseNode) NextSibling() Node {
	return n.nextSibling
}

func (n *BaseNode) setPrevSibling(prev Node) {
    n.prevSibling = prev
}

func (n *BaseNode) setNextSibling(next Node) {
    n.nextSibling = next
}

func (n *BaseNode) FirstChild() Node {
	return n.firstChild
}

func (n *BaseNode) LastChild() Node {
	return n.lastChild
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

	if n.lastChild != nil {
		n.lastChild.setNextSibling(child)
		child.setPrevSibling(n.lastChild)
		child.setNextSibling(nil)
		n.lastChild = child
	} else {
		n.firstChild = child
		n.lastChild = child
		child.setPrevSibling(nil)
        child.setNextSibling(nil)
	}
}

func (n *BaseNode) RemoveChild(child Node) {
	if child.Parent() != n {
		return
	}

	if child.PrevSibling() != nil {
		child.PrevSibling().setNextSibling(child.NextSibling())
	} else {
		n.firstChild = child.NextSibling()
	}

	if child.NextSibling() != nil {
		child.NextSibling().setPrevSibling(child.PrevSibling())
	} else {
		n.lastChild = child.PrevSibling()
	}

	child.SetParent(nil)
	child.setPrevSibling(nil)
    child.setNextSibling(nil)
}

func (n *BaseNode) InsertAfter(oldNode, newNode Node) {
	if oldNode.Parent() != n {
		return
	}
	newNode.SetParent(n)

	newNode.setPrevSibling(oldNode)
	newNode.setNextSibling(oldNode.NextSibling())

	if oldNode.NextSibling() != nil {
		oldNode.NextSibling().setPrevSibling(newNode)
	} else {
		n.lastChild = newNode
	}

	oldNode.setNextSibling(newNode)
}

func (n *BaseNode) ReplaceChild(oldNode, newNode Node) {
	if oldNode.Parent() != n {
		return
	}

	newNode.SetParent(n)

	newNode.setPrevSibling(oldNode.PrevSibling())
	newNode.setNextSibling(oldNode.NextSibling())

	if oldNode.PrevSibling() != nil {
		oldNode.PrevSibling().setNextSibling(newNode)
	} else {
		n.firstChild = newNode
	}

	if oldNode.NextSibling() != nil {
		oldNode.NextSibling().setPrevSibling(newNode)
	} else {
		n.lastChild = newNode
	}

	oldNode.SetParent(nil)
	oldNode.setPrevSibling(nil)
	oldNode.setNextSibling(nil)
}

func (n *BaseNode) IsOpen() bool {
	return n.isOpen
}

func (n *BaseNode) SetOpen(isOpen bool) {
	n.isOpen = isOpen
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
	return t >= NodeCodeBlock && t <= NodeParagraph
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
