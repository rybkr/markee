package parser

import (
	"markee/internal/ast"
)

type BlockFinalizer struct {
	ast.BaseVisitor
}

func NewBlockFinalizer() *BlockFinalizer {
	return &BlockFinalizer{}
}

func (f *BlockFinalizer) VisitDocument(node ast.Node) {
    ast.WalkChildren(f, node)
}

func (f *BlockFinalizer) VisitBlockQuote(node ast.Node) {
    ast.WalkChildren(f, node)
}

func (f *BlockFinalizer) VisitList(node ast.Node) {
    ast.WalkChildren(f, node)
}

func (f *BlockFinalizer) VisitListItem(node ast.Node) {
    ast.WalkChildren(f, node)
}

func (f *BlockFinalizer) VisitParagraph(node ast.Node) {
    var content string
    for child := node.FirstChild(); child != nil; child = child.NextSibling() {
        if contentNode, ok := child.(*ast.Content); ok {
            if len(content) > 0 {
                content += "\n"
            }
            content += contentNode.Literal
        }
    }
    
    for child := node.FirstChild(); child != nil; {
        next := child.NextSibling()
        if _, ok := child.(*ast.Content); ok {
            node.RemoveChild(child)
        }
        child = next
    }
    
    if len(content) > 0 {
        ParseInlines(node, content)
    }
}

func (f *BlockFinalizer) VisitHeading(node ast.Node) {
    var content string
    for child := node.FirstChild(); child != nil; child = child.NextSibling() {
        if contentNode, ok := child.(*ast.Content); ok {
            if len(content) > 0 {
                content += "\n"
            }
            content += contentNode.Literal
        }
    }
    
    for child := node.FirstChild(); child != nil; {
        next := child.NextSibling()
        if _, ok := child.(*ast.Content); ok {
            node.RemoveChild(child)
        }
        child = next
    }
    
    if len(content) > 0 {
        ParseInlines(node, content)
    }
}

func (f *BlockFinalizer) VisitCodeBlock(node ast.Node) {
    ast.WalkChildren(f, node)
}

func (f *BlockFinalizer) VisitHTMLBlock(node ast.Node) {
    ast.WalkChildren(f, node)
}

func (f *BlockFinalizer) VisitThematicBreak(node ast.Node) {
    ast.WalkChildren(f, node)
}
