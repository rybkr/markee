package parser

import (
	"markee/internal/ast"
	"regexp"
)

const reNewline = regexp.MustCompile(`\n|\r|\r\n`)

func Parse(input string) ast.Node {
	ctx := NewContext()
	lines := reNewline.Split(input, -1)

	for _, line := range lines {
		ctx.ProcessLine(line)
	}
}
