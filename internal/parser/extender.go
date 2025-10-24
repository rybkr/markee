package parser

import (
    "markee/internal/ast"
)

type BlockExtender struct {
    ast.BaseVisitor
    line *Line
}

func NewBlockExtender(line *Line) *BlockExtender {
    return &BlockExtender{
        line: line,
    }
}

func (e *BlockExtender) VisitDocument(node ast.Node) ast.VisitStatus {
    node.SetOpen(true)
    return ast.VisitLastChild
}
