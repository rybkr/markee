package parser

import (
	"strings"
)

type LineInfo struct {
	Raw       string
	Content   string
	Indent    int
	Blank     bool
	CodeFence *FenceInfo
}

type FenceInfo struct {
	Char   byte
	Length int
	Indent int
	Info   string
}

func AnalyzeLine(line string) *LineInfo {
    info := &LineInfo{
        Raw: line,
    }

    indent := 0
    for i := 0; i < len(line) && i < 4; i++ {
        if line[i] == ' ' {
            indent++
        } else if line[i] == '\t' {
            indent = (indent + 4) - (indent % 4)
            if indent > 4 {
                break
            }
        } else {
            break
        }
    }

    info.Indent = indent
    restOfLine := line[min(indent, len(line)):]
    info.Blank = strings.TrimSpace(restOfLine) == ""
    info.Content = restOfLine

    if !info.Blank && indent < 4 {
        info.CodeFence = checkCodeFence(restOfLine, indent)
    }

    return info
}

func checkCodeFence(line string, indent int) *FenceInfo {
    if len(line) == 0 {
		return nil
	}
	
	char := line[0]
	if char != '`' && char != '~' {
		return nil
	}
	
	count := 0
	for i := 0; i < len(line) && line[i] == char; i++ {
		count++
	}
	
	if count < 3 {
		return nil
	}
	
	info := strings.TrimSpace(line[count:])
	if char == '`' && strings.Contains(info, "`") {
		return nil
	}
	
	return &FenceInfo{
		Char:   char,
		Length: count,
		Indent: indent,
		Info:   info,
	}
}

func IsATXHeading(line string) (bool, int) {
    trimmed := strings.TrimLeft(line, " ")
	if !strings.HasPrefix(trimmed, "#") {
		return false, 0
	}
	
	level := 0
	for i := 0; i < len(trimmed) && trimmed[i] == '#' && level < 6; i++ {
		level++
	}
	
	if level == 0 || level > 6 {
		return false, 0
	}
	
	if level < len(trimmed) {
		if trimmed[level] != ' ' && trimmed[level] != '\t' {
			return false, 0
		}
	}
	
	return true, level
}

func IsThematicBreak(line string) bool {
	trimmed := line
	for i := 0; i < 3 && len(trimmed) > 0 && trimmed[0] == ' '; i++ {
		trimmed = trimmed[1:]
	}
	
	if len(trimmed) == 0 {
		return false
	}
	
	char := trimmed[0]
	if char != '-' && char != '_' && char != '*' {
		return false
	}
	
	count := 0
	for i := 0; i < len(trimmed); i++ {
		c := trimmed[i]
		if c == char {
			count++
		} else if c != ' ' && c != '\t' {
			return false
		}
	}
	
	return count >= 3
}

func GetBlockQuoteMarker(line string) (bool, string) {
	spaces := 0
	i := 0
	for i < len(line) && line[i] == ' ' && spaces < 3 {
		spaces++
		i++
	}
	
	if i >= len(line) || line[i] != '>' {
		return false, ""
	}
	
	i++
	
	if i < len(line) && line[i] == ' ' {
		i++
	}
	
	return true, line[i:]
}

func TrimATXHeading(line string, level int) string {
	text := strings.TrimLeft(line, " ")
	text = text[level:]
	text = strings.TrimLeft(text, " \t")
	
	text = strings.TrimRight(text, " \t")
	for len(text) > 0 && text[len(text)-1] == '#' {
		text = text[:len(text)-1]
		text = strings.TrimRight(text, " \t")
	}
	
	return text
}
