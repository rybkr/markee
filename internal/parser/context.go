package parser

import (
	"markee/internal/ast"
)

type Context struct {
	Doc *ast.Document
	Tip ast.Node
	Pos int
}

func NewContext() *Context {
	doc := ast.NewDocument()
	return &Context{
		Doc: doc,
		Tip: doc,
		Pos: 0,
	}
}

func (c *Context) AddChild(node ast.Node) {
	c.Tip.AddChild(node)
}

func (c *Context) SetTip(node ast.Node) {
	c.Tip = node
}

func (c *Context) CloseBlock() {
	if c.Tip.Parent() != nil {
		c.Tip = c.Tip.Parent()
	}
}
