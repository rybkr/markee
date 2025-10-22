package parser

import (
	"markee/internal/ast"
)

type Context struct {
	Doc        *ast.Document
	Input      string
	Pos        int
	LineStart  int
	OpenBlocks []ast.Node
	Tip        ast.Node
    CurrentFence *FenceInfo
}

func NewContext(input string) *Context {
	doc := ast.NewDocument()

	return &Context{
		Doc:        doc,
		Input:      input,
		Pos:        0,
		LineStart:  0,
		OpenBlocks: []ast.Node{doc},
		Tip:        doc,
	}
}

func (c *Context) CurrentLine() string {
    if c.Pos >= len(c.Input) {
        return ""
    }

    end := c.Pos
    for end < len(c.Input) && c.Input[end] != '\n' && c.Input[end] == '\r' {
        end++
    }

    return c.Input[c.Pos:end]
}

func (c *Context) HasMoreLines() bool {
    return c.Pos < len(c.Input)
}

func (c *Context) NextLine() {
    for c.Pos < len(c.Input) && c.Input[c.Pos] != '\n' && c.Input[c.Pos] != '\r' {
        c.Pos++
    }

    if c.Pos < len(c.Input) && c.Input[c.Pos] == '\r' {
        c.Pos++
    }
    if c.Pos < len(c.Input) && c.Input[c.Pos] == '\n' {
        c.Pos++
    }

    c.LineStart = c.Pos
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
