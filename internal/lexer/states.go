package lexer

// State functions implement the lexer's state machine. Each state function
// processes input at the current position and returns the next state function
// to execute or nil to halt lexing.
// This pattern enables context-aware tokenization, needed to differentiate
// between inline delimiters, block level delimiters, and code blocks.
type stateFunc func(*Lexer) stateFunc

func lexLineStart(l *Lexer) stateFunc {
	switch l.peek() {
	case 0:
		return lexEOF
	case '\n':
		return lexNewline
	case '#':
		return lexHeader
	case '>':
		return lexBlockquote
	case '-', '*', '+':
		return lexLineStartMarker
	case '`', '~':
		if isCodeFence(l.peekAhead(3)) {
			return lexCodeFence
		}
	}

	l.context = CtxInline
	return lexInline
}

func lexInline(l *Lexer) stateFunc {
	l.advanceWhile(isText)
	if l.pos > l.start {
		l.emit(TokenText)
	}

	switch l.peek() {
	case 0:
		return lexEOF
	case '\n':
		return lexNewline
	}

	return dispatchInlineDelimiter(l.peek())
}

func lexHeader(l *Lexer) stateFunc {
	// Headers have a max level of 6
	for l.peek() == '#' && l.pos-l.start < 6 {
		l.advance()
	}

	if isWhitespace(l.peek()) {
		l.emit(TokenHeader)
		l.advanceWhile(isWhitespace)
		l.ignore()
	} else {
		l.abort()
	}

	l.context = CtxInline
	return lexInline
}

func lexBlockquote(l *Lexer) stateFunc {
	l.advance()
	l.emit(TokenBlockquote)
	l.advanceWhile(isWhitespace)
	l.ignore()

	l.context = CtxInline
	return lexInline
}

func lexCodeFence(l *Lexer) stateFunc {
	l.advanceAhead(3)
	l.emit(TokenCodeFence)

	// TODO: (rybkr) Handle language specifiers such as "~~~py"
	l.advanceUntil(isNewline)
	l.ignore()

	l.context = CtxCodeBlock
	return lexCodeBlock
}

func lexCodeBlock(l *Lexer) stateFunc {
	// Check for closing fence at line start
	if l.column == 1 && isCodeFence(l.peekAhead(3)) {
		l.advanceAhead(3)
		l.emit(TokenCodeFence)
		l.advanceUntil(isNewline)
		l.ignore()

		l.context = CtxInline
		return lexInline
	}

	l.advanceWhile(isBenign)
	if l.pos > l.start {
		l.emit(TokenText)
	}

	switch l.peek() {
	case '\n':
		return lexNewline
	default:
		return lexEOF
	}
}

func lexLineStartMarker(l *Lexer) stateFunc {
	if isHorizontalRule(l.peekAhead(3)) {
		l.advanceAhead(3)
		l.emit(TokenHorizontalRule)
		l.advanceUntil(isNewline)
		l.ignore()

		l.context = CtxInline
		return lexInline
	}

	l.advance()
	if isWhitespace(l.peek()) {
		l.emit(TokenListMarker)
		l.advanceWhile(isWhitespace)
		l.ignore()
		l.context = CtxInline
		return lexInline
	}

	l.abort()
	l.context = CtxInline
	return lexInline
}

func lexStar(l *Lexer) stateFunc {
	for l.peek() == '*' && l.pos-l.start < 2 {
		l.advance()
	}

	l.emit(TokenStar)
	return lexInline
}

func lexUnderscore(l *Lexer) stateFunc {
	for l.peek() == '_' && l.pos-l.start < 2 {
		l.advance()
	}

	l.emit(TokenUnderscore)
	return lexInline
}

func lexBacktick(l *Lexer) stateFunc {
	l.advance()
	l.emit(TokenBacktick)

	for r := l.peek(); r != '`'; r = l.peek() {
		l.advance()
	}
	if l.pos > l.start {
		l.emit(TokenText)
	}
	l.advance()
	l.emit(TokenBacktick)

	return lexInline
}

func lexNewline(l *Lexer) stateFunc {
	l.advance()
	l.emit(TokenNewline)

	if l.context == CtxCodeBlock {
		return lexCodeBlock
	}

	return lexLineStart
}

func lexEOF(l *Lexer) stateFunc {
	l.emit(TokenEOF)
	return nil
}

func dispatchInlineDelimiter(r rune) stateFunc {
	switch r {
	case '*':
		return lexStar
	case '_':
		return lexUnderscore
	case '`':
		return lexBacktick
	default:
		return lexInline
	}
}
