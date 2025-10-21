package parser

//go:generate stringer -type=NodeType
type NodeType int

const (
	NodeDocument NodeType = iota

    NodeParagraph
    NodeHeader
    NodeCodeBlock
    NodeBlockquote
    NodeList
    NodeListItem
    NodeHorizontalRule

    NodeText
    NodeEmphasis
    NodeStrong
    NodeInlineCode
    NodeLineBreak
)

type Node struct {
	Type     NodeType
	Value    string
	Level    int
	Children []*Node
}

func NewNode(type NodeType) *Node {
    return &Node{
        Type:     type,
        Children: []*Node{},
    }
}

func (n *Node) AppendChild(child *Node) {
    n.Children = append(n.Children, child)
}
