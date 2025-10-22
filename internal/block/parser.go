package block

import (
	"markee/internal/ast"
	"strings"
	"unicode"
)

// block.Parser handles block-level structure parsing.
// See: https://spec.commonmark.org/0.31.2/#phase-1-block-structure
type Parser struct {
}

func New() *Parser {
}

func (p *Parser) Parse(input string) *ast.Node {
    lines := reNewline.Split(text, -1)
    for _, line := range lines {
        p.incorporateLine(line)
    }
}

func (p *Parser) incorporateLine(line string) {
}
