package parser

import (
    "markee/internal/block"
)

// Parser orchestrates block and inline parsing to produce a complete AST.
// Phase 1: Block structure parsing (containers and leaf blocks)
// Phase 2: Inline parsing (emphasis, links, code spans, etc.)
// See: https://spec.commonmark.org/0.31.2/#appendix-a-parsing-strategy
type Parser struct {
	blockParser  *block.Parser
}

func New() *Parser {
    p := &Parser{
        blockParser: block.New(),    
    }
}

func (p *Parser) Parse(input string) *Node {
    doc := p.blockParser.Parse()
}
