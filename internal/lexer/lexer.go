package lexer

type Context int

const (
	CtxLineStart Context = iota
	CtxInline
	CtxCodeBlock
)

type Lexer struct {
	input   string
	start   int
	pos     int
	line    int
	column  int
	tokens  []Token
	state   stateFunc
	context Context
}

func New(input string) *Lexer {
	return &Lexer{
		input:   input,
		start:   0,
		pos:     0,
		line:    1,
		column:  1,
		tokens:  make([]Token, 0),
		state:   lexLineStart,
		context: CtxLineStart,
	}
}

func (l *Lexer) Next() Token {
	if len(l.tokens) > 0 {
		tok := l.tokens[0]
		l.tokens = l.tokens[1:]
		return tok
	}

	// If we've reach EOF, keep returning EOF
	if l.state == nil {
		return Token{Type: TokenEOF}
	}

	for len(l.tokens) == 0 && l.state != nil {
		l.state = l.state(l)
	}

	if len(l.tokens) > 0 {
		tok := l.tokens[0]
		l.tokens = l.tokens[1:]
		return tok
	}

	return Token{Type: TokenEOF}
}

func (l *Lexer) All() []Token {
	var toks []Token
	for l.state != nil {
		tok := l.Next()
		toks = append(toks, tok)
	}
	return toks
}

func Tokenize(input string) []Token {
	return New(input).All()
}
