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
	e.lastMatch = node
	if child := node.LastChild(); child != nil {
		child.Accept(e)
	}
}

func (e *BlockExtender) VisitBlockQuote(node ast.Node) {
	if e.line.Indent < 4 && e.line.Peek(0) == '>' {
		e.line.Consume(1)
		e.line.ConsumeWhitespace()
		e.lastMatch = node
		if child := node.LastChild(); child != nil {
			child.Accept(e)
		}
	}
}

func (e *BlockExtender) VisitList(node ast.Node) {
	e.lastMatch = node
	if child := node.LastChild(); child != nil {
		child.Accept(e)
	}
}

func (e *BlockExtender) VisitListItem(node ast.Node) {
}

func (e *BlockExtender) VisitCodeBlock(node ast.Node) {
}

func (e *BlockExtender) VisitHTMLBlock(node ast.Node) {
}

func (e *BlockExtender) VisitThematicBreak(node ast.Node) {
}

func (e *BlockExtender) VisitHeading(node ast.Node) {
}

func (e *BlockExtender) VisitParagraph(node ast.Node) {
	if !e.line.IsBlank {
		e.lastMatch = node
		if child := node.LastChild(); child != nil {
			child.Accept(e)
		}
	}
}
