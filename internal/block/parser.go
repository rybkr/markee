package block

import (
	"markee/internal/ast"
	"strings"
	"unicode"
)

type Parser struct {
	lines   []string
	line    int
	offset  int
	indent  int
	root    *Node
	current *Node
}

func New(input string) *Parser {
	lines := strings.Split(input, "\n")
	root := ast.New(NodeDocument)

	return &Parser{
		lines:   lines,
		line:    0,
		root:    root,
		current: root,
	}
}

func (p *Parser) incorporateLine(line string) {
    p.offset = 0
    p.indent = 0

    for p.offset < len(line) && line[p.offset] == ' ' {
        p.indent++
        p.offset++
        if p.indent >= 
    }
}
