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
}
