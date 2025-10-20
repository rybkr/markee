package parser

import (
    "markee/internal/lexer"
)

type Context int

const (
    CtxBlock Context = iota
    CtxInline
    CtxCodeBlock
    CtxList
)

type Parser struct {
    tokens  []lexer.Token
    pos     int
    context Context
    stack   []*Node
}

func New(tokens []lexer.Token) *Parser {
	return &Parser{
		tokens:  tokens,
		pos:     0,
		context: CtxBlock,
		stack:   []*Node{},
	}
}

func (p *Parser) Parse() *Node {
    root := &Node{
        Type:     NodeDocument,
        Children: []*Node{},
    }
    p.push(root)

    state := 

    return doc
}
