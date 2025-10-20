package parser

//go:generate stringer -type=NodeType
type NodeType int

const (
    NodeDocument NodeType = iota
    NodeParagraph
    NodeHeader
    NodeBlockquote
    NodeCodeBlock
    NodeText
    NodeEmphasis
    NodeStrong
)

type Node struct {
    Type     NodeType
    Value    string
    Level    int
    Children []*Node
}
