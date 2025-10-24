package parser

import (
	"markee/internal/ast"
)

type MatchVisitor struct {
	ast.BaseVisitor
	line     *Line
	matched  bool
	newBlock ast.Node
}

func NewMatchVisitor(line *Line) *MatchVisitor {
	return &MatchVisitor{
		line:     line,
		matched:  false,
		newBlock: nil,
	}
}

func (v *MatchVisitor) Matched() bool {
	return v.matched
}

func (v *MatchVisitor) NewBlock() ast.Node {
	return v.newBlock
}

func (v *MatchVisitor) VisitDocument(node ast.Node) ast.VisitStatus {
    node.SetOpen(true)
    return ast.VisitContinue
}

func (v *MatchVisitor) VisitBlockQuote(node ast.Node) ast.VisitStatus {
    if v.line.Peek(0) == '>' {
        node.SetOpen(true)
        v.line.Consume(1)
        v.line.ConsumeWhitespace()
        return ast.VisitContinue
    }
    return ast.VisitStop
}

func (v *MatchVisitor) VisitList(node ast.Node) ast.VisitStatus {
    node.SetOpen(true)
    return ast.VisitContinue
}

func (v *MatchVisitor) VisitListItem(node ast.Node) ast.VisitStatus {
    listItem := node.(*ast.ListItem)
    if v.line.IsBlank {
        return ast.VisitContinue
    }
    if v.line.Indent >= listItem.Indent {
        return ast.VisitContinue
    }
    return ast.VisitStop
}

func (v *MatchVisitor) VisitCodeBlock(node ast.Node) ast.VisitStatus {
    return ast.VisitStop
}
