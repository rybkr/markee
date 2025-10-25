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
	for _, child := range node.Children() {
		child.Accept(f)
	}
}

func (f *BlockFinalizer) VisitBlockQuote(node ast.Node) {
	for _, child := range node.Children() {
		child.Accept(f)
	}
}

func (f *BlockFinalizer) VisitList(node ast.Node) {
    if list, ok := node.(*ast.List); ok {
        for _, child := range list.Children() {
            if lastGrandChild := child.LastChild(); child.NextSibling() != nil && lastGrandChild.Type() == ast.NodeBlankLine {
                list.IsTight = false
                break
            }
            for _, grandChild := range child.Children() {
                if lastGreatGrandChild := grandChild.LastChild(); grandChild.NextSibling() != nil && lastGreatGrandChild.Type() == ast.NodeBlankLine {
                    list.IsTight = false
                    break
                }
            }
            if list.IsTight {
                break
            }
        }
    }

	for _, child := range  node.Children() {
		child.Accept(f)
	}
}
