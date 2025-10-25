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

	return ctx.Doc
}

// incorporateLine handles line-by-line block parsing logic.
// See: https://spec.commonmark.org/0.31.2/#phase-1-block-structure
func incorporateLine(ctx *Context, line *Line) {
	// First we need to check which blocks the line "extends".  To extend a block, the line must
	// meet a requirement imposed by the block's type. For example, to extend a BlockQuote, the
	// line must begin with '>'.
	// Here, while we check which blocks are extended by the line, we can also consume relevant
	// tokens from the line. For example, if a BlockQuote is extended by the line "> continued...",
	// then the node should be set to open the the line's content should become "continued...".
	// We cannot close unmatched blocks yet because we may have a lazy continuation line.
	extender := NewBlockExtender(line)
	ctx.Doc.Accept(extender)
    lastMatched := extender.LastMatch()
    ctx.SetTip(lastMatched)


	// Second, we should check the line for any tokens that would create new blocks. For example,
	// if after consuming the extension markers, the line still starts with '>', then we should
	// create a new BlockQuote node.
	// This logic differs from the extension logic because we need to consider the precedence of
	// block nodes as defined by CommonMark, rather than the order in which they appear in the AST.
	// If we find a match, we should close unmatched blocks from the previous step.
	var newBlock ast.Node
	newBlock = matchNewBlock(line, ctx.Tip)

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
		ctx.SetTip(lastMatched)
	}

	// Next, we look at the rest fo the line and incorporate the content into the last open block.
    if !line.IsEmpty() && ctx.Tip.Type() != ast.NodeCodeBlock {
        content := ast.NewContent(strings.TrimSpace(line.Content))
        ctx.Tip.AddChild(content)
    } else if !line.IsEmpty() && ctx.Tip.Type() == ast.NodeCodeBlock {
        // Code blocks accept content differently - preserve exact content
        content := ast.NewContent(line.Content)
        ctx.Tip.AddChild(content)
    }
}
