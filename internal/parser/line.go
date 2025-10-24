package parser

type Line struct {
	Literal string
	Content string
	Indent  int
	Offset  int
	IsBlank bool
}

func NewLine(raw string) *Line {
	var offset, indent int
	for indent, offset = 0, 0; offset < len(raw); offset++ {
		if raw[offset] == ' ' {
			indent++
		} else if raw[offset] == '\t' {
			indent = (indent + 4) - (indent % 4)
			if indent > 4 {
				break
			}
		} else {
			break
		}
	}

	content := raw[offset:]
	isBlank := stringIsBlank(content)

	return &Line{
		Literal: raw,
		Content: content,
		Indent:  indent,
		Offset:  offset,
		IsBlank: isBlank,
	}
}

func stringIsBlank(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] != ' ' && s[i] != '\t' {
			return false
		}
	}
	return true
}
