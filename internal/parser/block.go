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
	
	// CRITICAL: Handle fenced code blocks first (they have highest precedence)
	if isInFencedCodeBlock(ctx) {
		handleCodeBlockLine(ctx, lineInfo)
		return
	}
	
	// Handle container blocks (blockquotes, lists, etc.)
	containersMatched := checkOpenContainers(ctx, lineInfo)
	closeUnmatchedBlocks(ctx, containersMatched)
	
	// Try to start new blocks or add content to existing ones
	if !lineInfo.Blank {
		tryNewBlocks(ctx, lineInfo)
	}
}

// Check if we're currently inside a fenced code block
func isInFencedCodeBlock(ctx *Context) bool {
	if ctx.Tip.Type() != ast.NodeCodeBlock {
		return false
	}
	cb := ctx.Tip.(*ast.CodeBlock)
	return cb.IsFenced && ctx.CurrentFence != nil
}

// Handle a line when inside a fenced code block
func handleCodeBlockLine(ctx *Context, lineInfo *LineInfo) {
	cb := ctx.Tip.(*ast.CodeBlock)
	
	// Check if this line closes the fence
	if lineInfo.CodeFence != nil && 
	   lineInfo.CodeFence.Char == ctx.CurrentFence.Char &&
	   lineInfo.CodeFence.Length >= ctx.CurrentFence.Length &&
	   lineInfo.Indent < 4 {
		// This is the closing fence - close the block
		finalizeBlock(cb)
		ctx.OpenBlocks = ctx.OpenBlocks[:len(ctx.OpenBlocks)-1]
		if len(ctx.OpenBlocks) > 0 {
			ctx.Tip = ctx.OpenBlocks[len(ctx.OpenBlocks)-1]
		}
		ctx.CurrentFence = nil
		return
	}
	
	// Not a closing fence - add the entire line as code content
	// Remove up to 4 spaces of indentation if present
	content := lineInfo.Raw
	removed := 0
	for i := 0; i < len(content) && removed < 4 && content[i] == ' '; i++ {
		removed++
	}
	if removed > 0 {
		content = content[removed:]
	}
	
	if cb.Literal != "" {
		cb.Literal += "\n"
	}
	cb.Literal += content
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
	content := lineInfo.Content

	// Thematic break
	if IsThematicBreak(content) {
		tb := ast.NewThematicBreak()
		ctx.AddChild(tb)
		finalizeBlock(tb)
		return
	}

	// ATX Heading
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

	// Fenced code block (opening fence)
	if lineInfo.CodeFence != nil && lineInfo.Indent < 4 {
		fence := lineInfo.CodeFence
		codeBlock := ast.NewCodeBlock(true)
		codeBlock.Language = fence.Info
		ctx.AddChild(codeBlock)
		ctx.OpenBlocks = append(ctx.OpenBlocks, codeBlock)
		ctx.Tip = codeBlock
		ctx.CurrentFence = fence
		// Important: Don't add the fence line itself as content
		return
	}

	// Block quote
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

	// Indented code block
	if lineInfo.Indent >= 4 {
		// Can't interrupt a paragraph
		if ctx.Tip.Type() == ast.NodeParagraph {
			// Add to paragraph instead
			textNode := ast.NewContent(content)
			ctx.Tip.AddChild(textNode)
			return
		}

		codeBlock := ast.NewCodeBlock(false)
		ctx.AddChild(codeBlock)
		ctx.OpenBlocks = append(ctx.OpenBlocks, codeBlock)
		ctx.Tip = codeBlock
		
		// Remove the 4-space indent
		if len(content) >= 4 {
			codeBlock.Literal = content[4:]
		}
		return
	}

	// Paragraph (default case for text content)
	if ctx.Tip.Type() != ast.NodeParagraph {
		para := ast.NewParagraph()
		ctx.AddChild(para)
		ctx.OpenBlocks = append(ctx.OpenBlocks, para)
		ctx.Tip = para
	}
	
	// Add content to paragraph
	if content != "" {
		textNode := ast.NewContent(content)
		ctx.Tip.AddChild(textNode)
	}
}

func finalizeBlock(block ast.Node) {
	switch block.Type() {
	case ast.NodeCodeBlock:
		cb := block.(*ast.CodeBlock)
		// Trim trailing newlines
		for len(cb.Literal) > 0 && cb.Literal[len(cb.Literal)-1] == '\n' {
			cb.Literal = cb.Literal[:len(cb.Literal)-1]
		}
	}
}

func finalizeDocument(ctx *Context) {
	for ctx.Tip.Type() != ast.NodeDocument {
		lastBlock := ctx.OpenBlocks[len(ctx.OpenBlocks)-1]
		finalizeBlock(lastBlock)
		ctx.OpenBlocks = ctx.OpenBlocks[:len(ctx.OpenBlocks)-1]
		if len(ctx.OpenBlocks) > 0 {
			ctx.Tip = ctx.OpenBlocks[len(ctx.OpenBlocks)-1]
		}
	}
}
