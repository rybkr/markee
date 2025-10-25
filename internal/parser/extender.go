package parser

import (
	"markee/internal/ast"
)

type BlockExtender struct {
	ast.BaseVisitor
	line       *Line
	lastMatch  ast.Node
	allMatched []ast.Node
}

func NewBlockExtender(line *Line) *BlockExtender {
	return &BlockExtender{
		line:       line,
		lastMatch:  nil,
		allMatched: make([]ast.Node, 0),
	}
}

func (e *BlockExtender) LastMatch() ast.Node {
	return e.lastMatch
}

func (e *BlockExtender) AllMatched() []ast.Node {
    return e.allMatched
}

func (e *BlockExtender) VisitDocument(node ast.Node) {
	e.lastMatch = node
    e.allMatched = append(e.allMatched, node)
    ast.WalkLastChild(e, node)
}

func (e *BlockExtender) VisitBlockQuote(node ast.Node) {
	if e.line.Indent < 4 && e.line.Peek(0) == '>' {
		e.line.Consume(1)
        if len(e.line.Content) > 0 && e.line.Content[0] == ' ' {
            e.line.Consume(1)
        }
		e.lastMatch = node
        e.allMatched = append(e.allMatched, node)
        ast.WalkLastChild(e, node)
	}
}

func (e *BlockExtender) VisitList(node ast.Node) {
	e.lastMatch = node
    e.allMatched = append(e.allMatched, node)
    ast.WalkLastChild(e, node)
}

func (e *BlockExtender) VisitListItem(node ast.Node) {
    e.lastMatch = node
    e.allMatched = append(e.allMatched, node)
    ast.WalkLastChild(e, node)
}

func (e *BlockExtender) VisitCodeBlock(node ast.Node) {
    codeBlock, ok := node.(*ast.CodeBlock)
    if !ok {
        return
    }

    if codeBlock.IsFenced {
        if e.isClosingFence(codeBlock) {
            e.line.ConsumeAll()
            return
        }
        e.lastMatch = node
        e.allMatched = append(e.allMatched, node)
    } else {
        if e.line.IsBlank || e.line.Indent >= 4 {
            e.lastMatch = node
            e.allMatched = append(e.allMatched, node)
        }
	}
}

func (e *BlockExtender) isClosingFence(codeBlock *ast.CodeBlock) bool {
    fenceCount := 0
    for fenceCount < len(e.line.Content) && e.line.Content[fenceCount] == codeBlock.FenceChar {
        fenceCount++
    }
    
    if fenceCount < 3 {
        return false
    }
    return true
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
        e.allMatched = append(e.allMatched, node)
        ast.WalkLastChild(e, node)
	}
}
