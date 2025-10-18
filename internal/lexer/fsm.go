package lexer

// stateFunc represents a state in the lexer state machine.
// Returns either the next state function to execute or nil (indicating lexing is complete).
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
		if l.peekString("```") || l.peekString("~~~") {
			return lexCodeFence
        }
	}

    l.context = CtxInline
    return lexInline
}

func lexInline(l *Lexer) stateFunc {
    for {
		switch l.peek() {
        case 0:
            l.emitIfPresent(TokenText)
            return lexEOF
		case '\n':
            l.emitIfPresent(TokenText)
            return lexNewline
		case '*':
            l.emitIfPresent(TokenText)
			return lexStar
		case '_':
            l.emitIfPresent(TokenText)
			return lexUnderscore
		case '`':
            l.emitIfPresent(TokenText)
			return lexBacktick
		case '[', ']', '(', ')':
            l.emitIfPresent(TokenText)
			return lexBracket
		default:
			l.next()
		}
	}
}

func lexHeader(l *Lexer) stateFunc {
	for l.peek() == '#' && l.pos-l.start < 6 {
		l.next()
	}

	// Headers must be followed by whitespace to be valid
	next := l.peek()
	if next != ' ' && next != '\t' {
		l.abort()
		return lexInline
	}

	l.emit(TokenHeader)
	l.skipWhitespace()

	l.context = CtxInline
	return lexInline
}

func lexBlockquote(l *Lexer) stateFunc {
	l.emitRune(TokenBlockquote)
	l.skipWhitespace()

	l.context = CtxInline
	return lexInline
}

func lexCodeFence(l *Lexer) stateFunc {
	l.advance(3)
	l.emit(TokenCodeFence)
	l.skipUntilEOL()

	l.context = CtxCodeBlock
	return lexNewline
}

func lexCodeBlockContent(l *Lexer) stateFunc {
	for {
		if l.column == 1 && (l.peekString("```") || l.peekString("~~~")) {
			l.advance(3)
			l.emit(TokenCodeFence)
			l.skipUntilEOL()
            l.context = CtxInline
            return lexInline
		}
        switch l.peek() {
        case 0:
            l.emitIfPresent(TokenText)
            return lexEOF
        case '\n':
            l.emitIfPresent(TokenText)
            return lexNewline
        default:
            l.next()
        }
	}
}

func lexLineStartMarker(l *Lexer) stateFunc {
	return nil
}

func lexStar(l *Lexer) stateFunc {
	return nil
}

func lexUnderscore(l *Lexer) stateFunc {
	return nil
}

func lexBacktick(l *Lexer) stateFunc {
	return nil
}

func lexBracket(l *Lexer) stateFunc {
	return nil
}

func lexNewline(l *Lexer) stateFunc {
    l.next()
    l.emit(TokenNewline)
    if l.context == CtxCodeBlock {
        return lexCodeBlockContent
    } else {
        return lexLineStart
    }
}

func lexEOF(l *Lexer) stateFunc {
	l.emit(TokenEOF)
	return nil
}
