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
	case lexer.TokenBlockquote:
		return parseBlockquote
	case lexer.TokenCodeFence:
		return parseCodeBlock
	case lexer.TokenText, lexer.TokenStar, lexer.TokenUnderscore:
		return parseParagraph
	default:
		p.advance()
		return parseBlock
	}
}

func parseHeader(p *Parser) stateFunc {
	tok := p.advance()
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

func parseBlockquote(p *Parser) stateFunc {
	p.advance()
	node := &Node{
		Type:     NodeBlockquote,
		Children: []*Node{},
	}

	p.appendChild(node)
	p.push(node)
	parseInlineUntil(p, lexer.TokenNewline, lexer.TokenEOF)
	p.pop()

	return parseBlock
}

func parseCodeBlock(p *Parser) stateFunc {
	p.advance()
	node := &Node{
		Type:     NodeCodeBlock,
		Children: []*Node{},
	}

	if p.peek().Type == lexer.TokenNewline {
		p.advance()
	}

	var content string
	for p.peek().Type != lexer.TokenCodeFence && p.peek().Type != lexer.TokenEOF {
		tok := p.peek()
		if tok.Type == lexer.TokenText {
			content += tok.Value
		} else if tok.Type == lexer.TokenNewline {
			content += tok.Value
		}
		p.advance()
	}
	p.advance()

	node.Value = content
	p.appendChild(node)
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
		tok := p.peek()
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
		case lexer.TokenStar, lexer.TokenUnderscore:
			switch len(tok.Value) {
			case 1:
				parseEmphasis(p)
			case 2:
				parseStrong(p)
			}
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
	parseInlineUntil(p, lexer.TokenStar, lexer.TokenUnderscore, lexer.TokenEOF)
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
	parseInlineUntil(p, lexer.TokenStar, lexer.TokenUnderscore, lexer.TokenEOF)
	p.pop()
}
