package parser

import (
	"markee/internal/ast"
	"regexp"
	"strings"
)

func matchNewBlock(line *Line) ast.Node {
	if thematicBreak := matchThematicBreak(line); thematicBreak != nil {
		line.ConsumeAll()
		return thematicBreak
	}

	if heading := matchATXHeading(line); heading != nil {
		line.Consume(heading.Level)
		line.ConsumeWhitespace()
		return heading
	}

    if blockQuote := matchBlockQuote(line); blockQuote != nil {
        line.Consume(1)
        line.ConsumeWhitespace()
        return blockQuote
    }

    if paragraph := matchParagraph(line); paragraph != nil {
        return paragraph
    }

    if blankLine := matchBlankLine(line); blankLine != nil {
        return blankLine
    }

	return nil
}

// See: https://spec.commonmark.org/0.31.2/#thematic-breaks
func matchThematicBreak(line *Line) *ast.ThematicBreak {
	// ([*\-_ \t]{3,}) : capture at least three instances of a thematic character
	reThematicBreak := regexp.MustCompile(`^([*\-_ \t]{3,})$`)
	matches := reThematicBreak.FindStringSubmatch(line.Content)
	if matches == nil {
		return nil
	}

	content := strings.ReplaceAll(matches[1], " ", "")
	content = strings.ReplaceAll(content, "\t", "")
	thematicChar := content[0]
	for i := 0; i < len(content); i++ {
		if content[i] != thematicChar {
			return nil
		}
	}

	return ast.NewThematicBreak()
}

// See: https://spec.commonmark.org/0.31.2/#atx-headings
func matchATXHeading(line *Line) *ast.Heading {
	//      (#{1,6}) : opening sequence of 1-6 '#' characters
	//        [ \t]+ : space/tab after '#' sequence, needed if heading has content
	//      ([^#]*?) : capture content lazily as to not eat trailing hashes
	// (?:[ \t]+#*)? : optional closing '#' sequence, preceded by as least one space/tab
	//        [ \t]* : optional space/tab at the end
	reATXHeading := regexp.MustCompile(`^(#{1,6})[ \t]+([^#]*?)(?:[ \t]+#*)?[ \t]*$`)
	if matches := reATXHeading.FindStringSubmatch(line.Content); matches != nil {
		level := len(matches[1])
		return ast.NewHeading(level)
	}

    // Alternative pattern for empty headings, does not need space after `#` sequence
    reATXHeadingEmpty := regexp.MustCompile(`^(#{1,6})[ \t]*(?:[ \t]+#*)?[ \t]*$`)
    if matches := reATXHeadingEmpty.FindStringSubmatch(line.Content); matches != nil {
        line.ConsumeAll()
        level := len(matches[1])
        return ast.NewHeading(level)
    }
    
	return nil
}

// See: https://spec.commonmark.org/0.31.2/#block-quotes
func matchBlockQuote(line *Line) *ast.BlockQuote {
    //      > : block quote marker
    // [ \t]? : optional single space or tab after '>'
    //   (.*) : capture content of the line
    reBlockQuote := regexp.MustCompile(`^>[ \t]?(.*)$`)
    if matches := reBlockQuote.FindStringSubmatch(line.Content); matches != nil {
        return ast.NewBlockQuote()
    }
    return nil
}

// See: https://spec.commonmark.org/0.31.2/#paragraphs
func matchParagraph(line *Line) *ast.Paragraph {
    if !line.IsBlank {
        return ast.NewParagraph()
    }
    return nil
}

// See: https://spec.commonmark.org/0.31.2/#blank-lines
func matchBlankLine(line *Line) *ast.BlankLine {
    if line.IsBlank {
        return ast.NewBlankLine()
    }
    return nil
}
