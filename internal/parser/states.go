package parser

import (
    "markee/internal/lexer"
)

type stateFunc func(*Parser) stateFunc

func parseBlock(p *Parser) stateFunc {
    switch p.peek().Type {
    case lexer.TokenEOF:
        return nil
    case lexer.TokenHeader:
        return parseHeader
    case lexer.TokenText, lexer.TokenEmphasis, lexer.TokenStrong:
        return parseParagraph
    default:
        p.advance()
        return parseBlock
    }
}

func parseHeader(p *Parser) stateFunc {
    tok := p.tokens[p.pos]
    node := &Node{
        Type:     NodeHeader,
        Level:    len(tok.Value),
        Children: []*Node{},
    }

    p.appendChild(node)
    p.push(node)
    parseInlineUntil(p, lexer.TokenNewline, lexer.TokenEOF)
    p.pop()

    return parseBlock
}

func parseParagraph(p *Parser) stateFunc {
    node := &Node{
        Type:     NodeParagraph,
        Children: []*Node{},
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
        case lexer.TokenEmphasis:
            parseEmphasis(p)
        case lexer.TokenStrong:
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
	parseInlineUntil(p, lexer.TokenEmphasis)
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
	parseInlineUntil(p, lexer.TokenStrong)
	p.pop()
}
