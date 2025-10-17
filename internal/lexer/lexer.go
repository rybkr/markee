package lexer

// Lexer tokenizes markdown input using a state machine pattern.
// Each state function examines the input and returns the next state.
type Lexer struct {
	input   string  // The complete markdown input string
	start   int     // Start position of current token
	pos     int     // Current position in input
	line    int     // Current line number (1-indexed)
	column  int     // Current column number (1-indexed)
	tokens  []Token // Collected tokens
	context Context // Current parsing context
}

// New creates a new lexer for the given input string.
func New(input string) *Lexer {
	return &Lexer{
		input:   input,
		start:   0,
		pos:     0,
		line:    1,
		column:  1,
		tokens:  make([]Token, 0),
		context: CtxLineStart,
	}
}

// Tokenize runs the lexer state machine and returns all tokens.
// This is the main entry point for lexing.
func (l *Lexer) Tokenize() []Token {
	// Run the state machine starting from line start state
	for sf := lexLineStart; sf != nil; {
		sf = sf(l)
	}
	return l.tokens
}

// peek returns the next rune without consuming it.
func (l *Lexer) peek() rune {
	if l.pos >= len(l.input) {
		return 0
	}
	return rune(l.input[l.pos])
}

// next consumes and returns the next rune in input.
// It advances the position and updates line/column tracking.
func (l *Lexer) next() rune {
	r := l.peek()
	l.pos += 1

	// Track line and column
	if r == '\n' {
		l.line++
		l.column = 1
	} else {
		l.column++
	}

	return r
}

// advance consumes the next n runes in input.
// Equivalent to calling next n times and ignoring output.
func (l *Lexer) advance(n int) {
	for i := 0; i < n; i++ {
		l.next()
	}
}

// peekString checks if the input at current position starts with s.
// Does not consume any input.
func (l *Lexer) peekString(s string) bool {
	if l.pos+len(s) > len(l.input) {
		return false
	}
	return l.input[l.pos:l.pos+len(s)] == s
}

// emit creates a token from the text between start and pos.
// After emitting, start is moved to pos for the next token.
func (l *Lexer) emit(t TokenType) {
	l.tokens = append(l.tokens, Token{
		Type:   t,
		Value:  l.input[l.start:l.pos],
		Line:   l.line,
		Column: l.column - (l.pos - l.start),
	})
	l.start = l.pos
}

// ignore discards any text between start and pos.
// Useful for skipping whitespace that does not need to be a token.
func (l *Lexer) ignore() {
	l.start = l.pos
}

// abort resets pos to start, restarting parsing the current token.
// Because tokens never span multiple lines, we can faithfully restore the column.
func (l *Lexer) abort() {
	l.column -= l.pos - l.start
	l.pos = l.start
}

// skipWhitespace consumes and ignores horizontal whitespace.
func (l *Lexer) skipWhitespace() {
	for r := l.peek(); r == ' ' || r == '\t'; r = l.next() {
		l.ignore()
	}
}

// skipUntilEOL consumes and ignores everything until the next newline.
func (l *Lexer) skipUntilEOL() {
	for r := l.peek(); r != '\n' && r != 0; r = l.next() {
		l.ignore()
	}
}
