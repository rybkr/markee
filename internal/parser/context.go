package parser

import (
    "markee/internal/ast"
)

type Context struct {
    Doc *ast.Document
    Tip ast.Node
}

func NewContext() *Context {
    doc := ast.NewDocument()
    return &Context{
        Doc: doc,
        Tip: doc,
    }
}

func (c *Context) AddChild(node ast.Node) {
    c.Tip.AddChild(node)
}

func (c *Context) SetTip(node ast.Node) {
    c.Tip = node
}

func (c *Context) CloseUnmatchedBlocks(lastMatched ast.Node) {
    for c.Tip != lastMatched && c.Tip != nil {
        c.CloseBlock()
    }
}

func (c *Context) CloseBlock() {
    if c.Tip.Parent() != nil {
        c.Tip = c.Tip.Parent()
    }
}

func (c *Context) GetOpenBlocks() []ast.Node {
    blocks := make([]ast.Node, 0)
    current := c.Tip
    for current != nil {
        blocks = append([]ast.Node{current}, blocks...)
        current = current.Parent()
    }
    return blocks
}
