package parser

import (
    "markee/internal/lexer"
)

type parseFunc func(*Parser) parseFunc

func parseBlock(p *Parser) parseFunc {
    switch p.tokens[p.pos].Type {
    case lexer.TokenEOF:
        return nil
    case lexer.TokenHeader:
        return parseHeader
    }

    p.advance()
    return parseBlock
}

func parseHeader(p *Parser) parseFunc {
    tok := p.advance()
    node := &Node{
        Type:     NodeHeader,
        Level:    len(tok.Value),
        Children: []*Node{}
    }

    p.appendChild(node)
    p.push(node)
    parseInlineUntil(p, lexer.TokenNewline, lexer.TokenEOF)
    p.pop()

    return parseBlock
}

func parseInlineUntil(p *Parser, stops ...lexer.TokenType) {
    for {
        tok := p.tokens[p.pos]

        for _, stop := range stops {
            if tok.Type == stop {
                p.advance()
                return
            }
        }

        switch tok.Type {
        case lexer.TokenText:
            p.appendChild(&Node{
                Type:  NodeText,
                Value: tok.Value,
            })
            p.advance()
        case lexer.TokenEmphasisDelimiter:
            parseEmphasis(p)
        case lexer.TokenStrongDelimiter:
            parseStrong(p)
        default:
            p.advance()
        }
    }
}

func parseEmphasis(p *Parser) {
	p.advance()
	node := &Node{
		Type:     NodeEmphasis,
		Children: []*Node{},
	}
	
	p.appendChild(node)
	p.push(node)
	parseInlineUntil(p, lexer.TokenEmphasisDelimiter)
	p.pop()
}

func parseStrong(p *Parser) {
    p.advance()
	node := &Node{
		Type:     NodeStrong,
		Children: []*Node{},
	}
	
	p.appendChild(node)
	p.push(node)
	parseInlineUntil(p, lexer.TokenStrongDelimiter)
	p.pop()
}
