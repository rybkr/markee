package lexer

func (l *Lexer) advance() rune {
	r := l.peek()
	if r == 0 {
		return 0
	}

	l.pos++
	l.column++

	if r == '\n' {
		l.line++
		l.column = 1
	}

	return r
}

func (l *Lexer) advanceAhead(n int) {
	for i := 0; i < n; i++ {
		l.advance()
	}
}

func (l *Lexer) peek() rune {
	if l.pos >= len(l.input) {
		return 0
	}
	return rune(l.input[l.pos])
}

func (l *Lexer) peekAhead(n int) string {
	if l.pos+n > len(l.input) {
		return l.input[l.pos:]
	}
	return l.input[l.pos : l.pos+n]
}

func (l *Lexer) emit(t TokenType) {
	l.tokens = append(l.tokens, Token{
		Type:   t,
		Value:  l.input[l.start:l.pos],
		Line:   l.line,
		Column: l.column - (l.pos - l.start),
	})
	l.start = l.pos
}

func (l *Lexer) abort() {
	// We need not update l.line, tokens cannot span multiple lines.
	l.column -= (l.pos - l.start)
	l.pos = l.start
}

func (l *Lexer) ignore() {
	l.start = l.pos
}

type predicateFunc func(rune) bool

// advanceWhile consumes runes while predicate returns true.
// Returns count of runes consumed.
func (l *Lexer) advanceWhile(predicate predicateFunc) int {
	count := 0
	for predicate(l.peek()) {
		l.advance()
		count++
	}
	return count
}

// advanceUntil consumes runes until predicate returns true.
// Stops before consuming the matching rune.
// Returns count of runes consumed.
func (l *Lexer) advanceUntil(predicate predicateFunc) int {
	count := 0
	for r := l.peek(); r != 0 && !predicate(r); r = l.peek() {
		l.advance()
		count++
	}
	return count
}

var isWhitespace predicateFunc = func(r rune) bool {
	return r == ' ' || r == '\t'
}

var isNewline predicateFunc = func(r rune) bool {
	return r == '\n'
}

var isEOF predicateFunc = func(r rune) bool {
    return r == 0
}

var isInlineDelimiter predicateFunc = func(r rune) bool {
	return r == '*' || r == '_' || r == '`' || r == '[' ||
		r == ']' || r == '(' || r == ')' || r == '!'
}

var isBenign predicateFunc = func(r rune) bool {
	return r != '\n' && r != 0
}

var isText predicateFunc = func(r rune) bool {
	return isBenign(r) && !isInlineDelimiter(r)
}

func isCodeFence(s string) bool {
	return s == "```" || s == "~~~"
}

func isHorizontalRule(s string) bool {
	return s == "---" || s == "***"
}
