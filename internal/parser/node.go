package parser

type NodeType int

const (
    NodeDocument NodeType = iota
    NodeParagraph
    NodeHeader
)

type Node struct {
    Type
    Value
    Level
    Children []*Node
}
