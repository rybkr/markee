package parser

import (
	"markee/internal/ast"
	"regexp"
	"strings"
)

type blockMatcher struct {
	priority     int
	name         string
	match        func(*Line) ast.Node
	canInterrupt func(ast.Node) bool
}

var blockMatchers = []blockMatcher{
	{priority: 1, name: "thematic_break", match: matchThematicBreak, canInterrupt: alwaysTrue},
	{priority: 2, name: "atx_heading", match: matchATXHeading, canInterrupt: alwaysTrue},
	{priority: 3, name: "fenced_code", match: matchFencedCodeBlock, canInterrupt: alwaysTrue},
	{priority: 4, name: "block_quote", match: matchBlockQuote, canInterrupt: alwaysTrue},
	{priority: 5, name: "indented_code", match: matchIndentedCodeBlock, canInterrupt: cannotInterruptParagraph},
	{priority: 10, name: "paragraph", match: matchParagraph, canInterrupt: alwaysTrue},
}

func alwaysTrue(n ast.Node) bool { return true }

func cannotInterruptParagraph(n ast.Node) bool {
	_, isParagraph := n.(*ast.Paragraph)
	return !isParagraph
}

func matchNewBlock(line *Line, currentTip ast.Node) ast.Node {
	if currentTip.Type() == ast.NodeCodeBlock {
		if fence := matchFencedCodeBlock(line); fence != nil {
			if fence.(*ast.CodeBlock).FenceChar == currentTip.(*ast.CodeBlock).FenceChar {
				line.ConsumeAll()
				currentTip.SetOpen(false)
				return nil
			}
		}
		line.Content = strings.Repeat(" ", line.Indent) + line.Content
		return nil
	}

	for _, matcher := range blockMatchers {
		if !matcher.canInterrupt(currentTip) {
			continue
		}
		if block := matcher.match(line); block != nil {
			if block.Type() == ast.NodeParagraph && currentTip.Type() == ast.NodeParagraph {
				return nil
			}
			return block
		}
	}
	return nil
}

// See: https://spec.commonmark.org/0.31.2/#thematic-breaks
func matchThematicBreak(line *Line) ast.Node {
	if line.Indent >= 4 {
		return nil
	}

	reAsterisk := regexp.MustCompile(`^[* \t]*\*[* \t]*\*[* \t]*\*[* \t]*$`)
	reHyphen := regexp.MustCompile(`^[\- \t]*-[\- \t]*-[\- \t]*-[\- \t]*$`)
	reUnderscore := regexp.MustCompile(`^[_ \t]*_[_ \t]*_[_ \t]*_[_ \t]*$`)

	if reAsterisk.MatchString(line.Content) || reHyphen.MatchString(line.Content) || reUnderscore.MatchString(line.Content) {
		line.ConsumeAll()
		return ast.NewThematicBreak()
	}

	return nil
}

// See: https://spec.commonmark.org/0.31.2/#atx-headings
func matchATXHeading(line *Line) ast.Node {
	if line.Indent >= 4 {
		return nil
	}

	// ^(#{1,6}) : 1-6 '#' characters
	// [ \t]*    : optional spaces/tabs
	// (?:#*)?   : optional trailing hashes (no space required when empty)
	// [ \t]*$   : optional trailing whitespace
	reATXHeadingEmpty := regexp.MustCompile(`^(#{1,6})[ \t]*(?:#*)?[ \t]*$`)

	if matches := reATXHeadingEmpty.FindStringSubmatch(line.Content); matches != nil {
		level := len(matches[1])
		line.ConsumeAll()
		return ast.NewHeading(level)
	}

	//      (#{1,6}) : opening sequence of 1-6 '#' characters
	//        [ \t]+ : space/tab after '#' sequence, needed if heading has content
	//      ([^#]*?) : capture content lazily as to not eat trailing hashes
	// (?:[ \t]+#*)? : optional closing '#' sequence, preceded by as least one space/tab
	//        [ \t]* : optional space/tab at the end
	reATXHeadingWithContent := regexp.MustCompile(`^(#{1,6})[ \t]+(.*?)(?:[ \t]+#*)?[ \t]*$`)

	if matches := reATXHeadingWithContent.FindStringSubmatch(line.Content); matches != nil {
		level := len(matches[1])
		line.Consume(len(matches[1]))
		line.ConsumeWhitespace()
		line.KeepUntil(len(matches[2]))
		return ast.NewHeading(level)
	}

	return nil
}

// See: https://spec.commonmark.org/0.31.2/#block-quotes
func matchBlockQuote(line *Line) ast.Node {
	if line.Indent >= 4 {
		return nil
	}
	if line.Peek(0) == '>' {
		line.Consume(1)
		if line.Peek(0) == ' ' {
			line.Consume(1)
		}
		return ast.NewBlockQuote()
	}
	return nil
}

// See: https://spec.commonmark.org/0.31.2/#paragraphs
func matchParagraph(line *Line) ast.Node {
	if !line.IsBlank && !line.IsEmpty() {
		return ast.NewParagraph()
	}
	return nil
}

// See: https://spec.commonmark.org/0.31.2/#indented-code-blocks
func matchIndentedCodeBlock(line *Line) ast.Node {
	if line.Indent >= 4 {
        line.Consume(4)
		return ast.NewCodeBlock(false)
	}
	return nil
}

// See: https://spec.commonmark.org/0.31.2/#fenced-code-blocks
func matchFencedCodeBlock(line *Line) ast.Node {
	if line.Indent >= 4 {
		return nil
	}

	// ^(`{3,}|~{3,}) : at least 3 backticks or tildes
	// [ \t]*         : optional spaces/tabs
	// ([^ \t`]*)     : info string (no spaces, tabs, or backticks)
	// [ \t]*$        : optional trailing whitespace
	reFence := regexp.MustCompile("^(`{3,}|~{3,})[ \\t]*([^ \\t`]*)[ \\t]*$")

	if matches := reFence.FindStringSubmatch(line.Content); matches != nil {
		fenceStr := matches[1]
		infoString := matches[2]

		codeBlock := ast.NewCodeBlock(true)
		codeBlock.FenceChar = fenceStr[0]
        codeBlock.FenceLen = len(fenceStr)
		codeBlock.Language = infoString

		line.ConsumeAll()
		return codeBlock
	}
	return nil
}

// See: https://spec.commonmark.org/0.31.2/#setext-headings
func matchSetextHeadingUnderline(line *Line, currentTip ast.Node) int {
    // Only check if we're in an open paragraph
    if currentTip.Type() != ast.NodeParagraph || !currentTip.IsOpen() {
        return 0
    }
    
    if line.Indent >= 4 {
        return 0
    }
    
    // Check for = underline (level 1)
    reEquals := regexp.MustCompile(`^=+[ \t]*$`)
    if reEquals.MatchString(line.Content) {
        return 1
    }
    
    // Check for - underline (level 2)
    reHyphens := regexp.MustCompile(`^-+[ \t]*$`)
    if reHyphens.MatchString(line.Content) {
        return 2
    }
    
    return 0
}
