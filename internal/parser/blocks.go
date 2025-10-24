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
	//        [ \t]* : optional space/tab after '#' sequence, needed if heading has content
	//         (.*?) : capture content lazily as to not eat trailing hashes
	// (?:[ \t]+#*)? : optional closing '#' sequence, preceded by as least one space/tab
	//        [ \t]* : optional space/tab at the end
	reATXHeading := regexp.MustCompile(`^(#{1,6})[ \t]*(.*?)(?:[ \t]+#*)?[ \t]*$`)
	if matches := reATXHeading.FindStringSubmatch(line.Content); matches != nil {
		level := len(matches[1])
		heading := ast.NewHeading(level)
		content := ast.NewContent(strings.TrimSpace(matches[2]))
		heading.AddChild(content)
		return heading
	}
	return nil
}
