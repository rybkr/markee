package parser

import (
    "markee/internal/lexer"
)

func (p *Parser) advance() lexer.Token {
    tok := p.peek()
    if tok.Type == lexer.TokenEOF {
        return tok
    }

    p.pos++
    return tok
}

func (p *Parser) peek() lexer.Token {
    if p.pos >= len(p.tokens) {
        return lexer.Token{Type: lexer.TokenEOF}
    }
    return p.tokens[p.pos]
}

func (p *Parser) push(node *Node) {
    p.stack = append(p.stack, node)
}

func (p *Parser) pop() *Node {
    if len(p.stack) == 0 {
        return nil
    }
    node := p.top()
    p.stack = p.stack[:len(p.stack)-1]
    return node
}

func (p *Parser) top() *Node {
    if len(p.stack) == 0 {
        return nil
    }
    return p.stack[len(p.stack)-1]
}

func (p *Parser) appendChild(node *Node) {
    if parent := p.top(); parent != nil {
        parent.Children = append(parent.Children, node)
    }
}
