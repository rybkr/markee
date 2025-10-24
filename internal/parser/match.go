package parser

import (
    "markee/internal/ast"
    "regexp"
    "strings"
)

func matchNewBlock(line *Line) ast.Node {
    if heading := matchATXHeading(line); heading != nil {
        line.Consume(heading.Level)
        line.ConsumeWhitespace()
        return heading
    }
    return nil
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
