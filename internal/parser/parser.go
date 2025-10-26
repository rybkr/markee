package parser

import (
	"bufio"
	"markee/internal/ast"
	"strings"
)

// Parse is the main parsing entry point, turns raw strings into an AST.
// See: https://spec.commonmark.org/0.31.2/#appendix-a-parsing-strategy
func Parse(input string) *ast.Document {
	ctx := NewContext()

	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		line := NewLine(scanner.Text())
		incorporateLine(ctx, line)
	}

	ctx.CloseUnmatchedBlocks(ctx.Doc)

	finalizer := NewBlockFinalizer()
	ctx.Doc.Accept(finalizer)

	return ctx.Doc
}

// incorporateLine handles line-by-line block parsing logic.
// See: https://spec.commonmark.org/0.31.2/#phase-1-block-structure
func incorporateLine(ctx *Context, line *Line) {
	extender := NewBlockExtender(line)
	ctx.Doc.Accept(extender)
	lastMatched := extender.LastMatch()

    if line.IsBlank {
        if ctx.Tip.Type() == ast.NodeParagraph {
            ctx.Tip.SetOpen(false)
            ctx.SetTip(ctx.Tip.Parent())
        }
        ctx.SetTip(lastMatched)
        return
    }

	newBlock := matchNewBlock(line, ctx.Tip)

	if newBlock != nil {
		ctx.CloseUnmatchedBlocks(lastMatched)
		ctx.AddChild(newBlock)
		ctx.SetTip(newBlock)

		if !newBlock.Type().IsLeaf() {
			for {
				nextBlock := matchNewBlock(line, ctx.Tip)
				if nextBlock == nil {
					break
				}
				ctx.AddChild(nextBlock)
				ctx.SetTip(nextBlock)

				if nextBlock.Type().IsLeaf() {
					break
				}
			}
		}
	} else {
		if ctx.Tip.Type() == ast.NodeParagraph && lastMatched != ctx.Tip {
			ctx.SetTip(ctx.Tip)
		} else {
			ctx.CloseUnmatchedBlocks(lastMatched)
			ctx.SetTip(lastMatched)
		}
	}

	if !line.IsEmpty() {
		if ctx.Tip.Type() == ast.NodeCodeBlock {
			content := ast.NewContent(line.Content)
			ctx.Tip.AddChild(content)
		} else if ctx.Tip.Type() != ast.NodeDocument {
			content := ast.NewContent(strings.TrimSpace(line.Content))
			ctx.Tip.AddChild(content)
		}
	}
}
