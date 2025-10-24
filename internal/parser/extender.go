package parser

import (
	"markee/internal/ast"
)

type BlockExtender struct {
	ast.BaseVisitor
	line      *Line
	lastMatch ast.Node
}

func NewBlockExtender(line *Line) *BlockExtender {
	return &BlockExtender{
		line:      line,
		lastMatch: nil,
	}
}

func (e *BlockExtender) LastMatch() ast.Node {
    return e.lastMatch
}

func (e *BlockExtender) VisitDocument(node ast.Node) {
	node.SetOpen(true)
    e.lastMatch = node
    if child := node.LastChild(); child != nil {
        child.Accept(e)
    }
}

func (e *BlockExtender) VisitParagraph(node ast.Node) {
    if !e.line.IsBlank {
        node.SetOpen(true)
        e.lastMatch = node
        if child := node.LastChild(); child != nil {
            child.Accept(e)
        }
    }
}
