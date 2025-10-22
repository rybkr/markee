package ast

type BaseNode struct {
    nodeType NodeType
    parent   Node
    children []Node
}

func New(t NodeType) BaseNode {
    return BaseNode{
        nodeType: t,
        children: make([]Node, 0),
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

func (n *BaseNode) Accept(v Visitor) {
    v.Visit(n)
}
