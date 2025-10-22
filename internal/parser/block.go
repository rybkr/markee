package parser

import (
	"markee/internal/ast"
)

func ParseBlocks(ctx *Context) {
	for ctx.HasMoreLines() {
		processLine(ctx)
		ctx.NextLine()
	}
	finalizeDocument(ctx)
}

func processLine(ctx *Context) {
	line := ctx.CurrentLine()
	lineInfo := AnalyzeLine(line)
	containersMatched := checkOpenContainers(ctx, lineInfo)
	closeUnmatchedBlocks(ctx, containersMatched)
	tryNewBlocks(ctx, lineInfo)
    addTextContent(ctx, lineInfo.Content)
}

func checkOpenContainers(ctx *Context, lineInfo *LineInfo) int {
	matched := 0

	for i := 0; i < len(ctx.OpenBlocks); i++ {
		block := ctx.OpenBlocks[i]

		switch block.Type() {
		case ast.NodeDocument:
			matched++

		case ast.NodeBlockQuote:
			if hasMarker, remaining := GetBlockQuoteMarker(lineInfo.Content); hasMarker {
				lineInfo.Content = remaining
				matched++
			} else {
				return matched
			}

		case ast.NodeList:
			matched++

		case ast.NodeListItem:
			matched++

		default:
			matched++
		}
	}

	return matched
}

func closeUnmatchedBlocks(ctx *Context, keepCount int) {
	for len(ctx.OpenBlocks) > keepCount {
		lastBlock := ctx.OpenBlocks[len(ctx.OpenBlocks)-1]
		finalizeBlock(lastBlock)
		ctx.OpenBlocks = ctx.OpenBlocks[:len(ctx.OpenBlocks)-1]
		if len(ctx.OpenBlocks) > 0 {
			ctx.Tip = ctx.OpenBlocks[len(ctx.OpenBlocks)-1]
		}
	}
}

func tryNewBlocks(ctx *Context, lineInfo *LineInfo) {
	if lineInfo.Blank {
		return
	}
	content := lineInfo.Content

	if IsThematicBreak(content) {
		tb := ast.NewThematicBreak()
		ctx.AddChild(tb)
		finalizeBlock(tb)
		return
	}

	if isHeading, level := IsATXHeading(content); isHeading {
		heading := ast.NewHeading(level)
		ctx.AddChild(heading)
		ctx.OpenBlocks = append(ctx.OpenBlocks, heading)
		ctx.Tip = heading

		text := TrimATXHeading(content, level)
		if text != "" {
			textNode := ast.NewContent(text)
			heading.AddChild(textNode)
		}

		finalizeBlock(heading)
		ctx.OpenBlocks = ctx.OpenBlocks[:len(ctx.OpenBlocks)-1]
		ctx.Tip = ctx.OpenBlocks[len(ctx.OpenBlocks)-1]
		return
	}

	if lineInfo.CodeFence != nil && lineInfo.Indent < 4 {
		fence := lineInfo.CodeFence
		codeBlock := ast.NewCodeBlock(true)
		codeBlock.Language = fence.Info
		ctx.AddChild(codeBlock)
		ctx.OpenBlocks = append(ctx.OpenBlocks, codeBlock)
		ctx.Tip = codeBlock
		ctx.CurrentFence = fence
		return
	}

	if hasMarker, remaining := GetBlockQuoteMarker(content); hasMarker {
		bq := ast.NewBlockQuote()
		ctx.AddChild(bq)
		ctx.OpenBlocks = append(ctx.OpenBlocks, bq)
		ctx.Tip = bq
		lineInfo.Content = remaining
		if remaining != "" {
			tryNewBlocks(ctx, lineInfo)
		}
		return
	}

	if lineInfo.Indent >= 4 {
		if ctx.Tip.Type() == ast.NodeParagraph {
			return
		}

		codeBlock := ast.NewCodeBlock(false)
		ctx.AddChild(codeBlock)
		ctx.OpenBlocks = append(ctx.OpenBlocks, codeBlock)
		ctx.Tip = codeBlock
		codeBlock.Literal = content[4:]
		return
	}

	if ctx.Tip.Type() != ast.NodeParagraph {
		para := ast.NewParagraph()
		ctx.AddChild(para)
		ctx.OpenBlocks = append(ctx.OpenBlocks, para)
		ctx.Tip = para
	}
}

func addTextContent(ctx *Context, content string) {
	if content == "" {
		return
	}

	switch ctx.Tip.Type() {
	case ast.NodeParagraph:
		textNode := ast.NewContent(content)
		ctx.Tip.AddChild(textNode)

	case ast.NodeCodeBlock:
		cb := ctx.Tip.(*ast.CodeBlock)
		if cb.IsFenced && ctx.CurrentFence != nil {
			fenceInfo := checkCodeFence(content, 0)
			if fenceInfo != nil && fenceInfo.Char == ctx.CurrentFence.Char &&
				fenceInfo.Length >= ctx.CurrentFence.Length {
				finalizeBlock(cb)
				ctx.OpenBlocks = ctx.OpenBlocks[:len(ctx.OpenBlocks)-1]
				ctx.Tip = ctx.OpenBlocks[len(ctx.OpenBlocks)-1]
				ctx.CurrentFence = nil
				return
			}
		}

		if cb.Literal != "" {
			cb.Literal += "\n"
		}
		cb.Literal += content
	}
}

func finalizeBlock(block ast.Node) {
	switch block.Type() {
	case ast.NodeCodeBlock:
		cb := block.(*ast.CodeBlock)
		for len(cb.Literal) > 0 && cb.Literal[len(cb.Literal)-1] == '\n' {
			cb.Literal = cb.Literal[:len(cb.Literal)-1]
		}
	}
}

func finalizeDocument(ctx *Context) {
	for ctx.Tip.Type() != ast.NodeDocument {
		ctx.CloseBlock()
	}
}
