package parser

//go:generate stringer -type=NodeType
type NodeType int

const (
    NodeDocument NodeType = iota
    NodeParagraph
    NodeHeader
)

type Node struct {
    Type     NodeType
    Value    string
    Level    int
    Children []*Node
}
