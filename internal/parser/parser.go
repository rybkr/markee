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
	stack   []*Node
	state   stateFunc
	context Context
}

func New(tokens []lexer.Token) *Parser {
	return &Parser{
		tokens:  tokens,
		pos:     0,
		stack:   []*Node{},
		state:   parseBlock,
		context: CtxBlock,
	}
}

func (p *Parser) All() *Node {
	root := &Node{
		Type:     NodeDocument,
		Children: []*Node{},
	}
	p.push(root)

	for p.state != nil {
		p.state = p.state(p)
	}

	p.pop()
	return root
}

func Parse(input string) *Node {
	return New(lexer.Tokenize(input)).All()
}
