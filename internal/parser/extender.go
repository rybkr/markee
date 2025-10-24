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

func (e *BlockExtender) VisitDocument(node ast.Node) ast.VisitStatus {
	node.SetOpen(true)
    e.lastMatch = node
	return ast.VisitLastChild
}

func (e *BlockExtender) VisitParagraph(node ast.Node) ast.VisitStatus {
    if !e.line.IsBlank {
        node.SetOpen(true)
        e.lastMatch = node
        return ast.VisitLastChild
    }
    return ast.VisitStop
}
