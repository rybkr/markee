package parser

import (
    "markee/internal/lexer"
)

func (p *Parser) advance() {
    if p.pos < len(p.tokens)-1 {
        p.pos++
        p.current = p.tokens[p.pos]
    }
}

func (p *Parser) peek() {
    if p.pos < len(p.tokens)-1 {
        return p.tokens[p.pos+1]
    }
    return lexer.Token{Type: lexer.TokenEOF}
}
