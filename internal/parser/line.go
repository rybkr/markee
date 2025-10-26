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
	for indent, offset = 0, 0; offset < min(4, len(raw)); offset++ {
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

func (l *Line) Consume(n int) {
	if n > len(l.Content) {
		n = len(l.Content)
	}
	l.Content = l.Content[n:]
	l.Offset += n
}

func (l *Line) ConsumeWhitespace() int {
	consumed := 0
	for consumed < len(l.Content) {
		if l.Content[consumed] == ' ' {
			consumed++
		} else {
			break
		}
	}
	l.Consume(consumed)
	return consumed
}

func (l *Line) ConsumeAll() int {
    consumed := len(l.Content)
    l.Consume(consumed)
    return consumed
}

func (l *Line) KeepUntil(n int) {
    if n > len(l.Content) {
        return
    }
    l.Content = l.Content[:n]
}

func (l *Line) Peek(i int) byte {
	if i >= len(l.Content) {
		return 0
	}
	return l.Content[i]
}

func (l *Line) IsEmpty() bool {
	return len(l.Content) == 0
}

func stringIsBlank(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] != ' ' && s[i] != '\t' {
			return false
		}
	}
	return true
}
