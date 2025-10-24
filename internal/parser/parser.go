package parser

import (
	"markee/internal/ast"
	"regexp"
)

var reNewline = regexp.MustCompile(`\n|\r|\r\n`)

func Parse(input string) *ast.Document {
	ctx := NewContext()
	lines := reNewline.Split(input, -1)

	for _, line := range lines {
		incorporateLine(ctx, NewLine(line))
	}
	closeUnmatchedBlocks(ctx)

	return ctx.Doc
}

func incorporateLine(ctx *Context, line *Line) {
}

func closeUnmatchedBlocks(ctx *Context) {
}

func handleUnmatched(ctx *Context, line *Line) {
	if ctx.Tip.Type() == ast.NodeParagraph {
		return
	}
	paragraph := ast.NewParagraph()
	ctx.AddChild(paragraph)
	ctx.SetTip(paragraph)
}
