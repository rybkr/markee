package lexer

// Context represents where we are in the document structure.
// This allows the lexer to make context-sensitive decisions about what characters mean.
type Context int

const (
    CtxLineStart Context = iota // At the beginning of a line
    CtxInline                   // In the middle of a line
    CtxCodeBlock                // Inside a fenced code block
)

// stateFunc represents a state in the lexer state machine.
// Each state function processes input according the rules of its state.
// Returns either the next state function to execute or nil (indicating lexing is complete).
type stateFunc func(*Lexer) stateFunc

// lexLineStart handles the start of a line.
// Here, block level syntax can appear (headers, lists, code fences, etc).
func lexLineStart(l *Lexer) stateFunc {
    switch l.peek() {
    case '#':
        return lexHeader
    case '>':
        return lexBlockquote
    case '`':
        if l.peekString("```") {
            return lexCodeFence
        } else {
            l.context = CtxInline
            return lexInline
        }
    case '-', '*', '+':
        // Could be a list marker or horizontal rule or emphasis
        return lexLineStartMarker
    case '\n':
        l.next()
        l.emit(TokenNewline)
        return lexLineStart
    case 0:
        return lexEOF
    default:
        // Regular text, switch to inline mode
        l.context = CtxInline
        return lexInline
    }
}

// lexInline handles inline content and inline markdown syntax.
// (emphasis, code spans, links, etc).
func lexInline(l *Lexer) stateFunc {
    // Consume text until we hit a special character
    for {
        switch l.peek() {
        case '*':
            if l.pos > l.start {
                l.emit(TokenText)
            }
            return lexStar
        case '_':
            if l.pos > l.start {
                l.emit(TokenText)
            }
            return lexUnderscore
        case '`':
			if l.pos > l.start {
				l.emit(TokenText)
			}
			return lexBacktick
        case '[', ']', '(', ')':
			if l.pos > l.start {
				l.emit(TokenText)
			}
			return lexBracket
        case '\n':
			if l.pos > l.start {
				l.emit(TokenText)
			}
            l.next()
			l.emit(TokenNewline)
			l.context = CtxLineStart
			return lexLineStart
        case 0:
			if l.pos > l.start {
				l.emit(TokenText)
			}
			return lexEOF
		default:
			l.next()
        }
    }
}

// lexHeader recognizes header markers (# ## ### etc) at line start.
// Headers must be at the start of a line and followed by a space.
func lexHeader(l *Lexer) stateFunc {
    for l.peek() == '#' && l.pos - l.start < 6 {
        l.next()
    }

    // Headers must be followed by whitespace to be valid
    next := l.peek()
    if next != ' ' && next != '\t' {
        // Not a valid header, abort lexHeader
        l.abort()
        return lexInline
    }

    l.emit(TokenHeader)

    for l.peek() == ' ' || l.peek() == '\t' {
        l.next()
        l.ignore()
    }

    l.context = CtxInline
    return lexInline
}

// lexBlockquote
func lexBlockquote(l *Lexer) stateFunc {
    return nil
}

// lexCodeFence
func lexCodeFence(l *Lexer) stateFunc {
    return nil
}

// lexLineStartMarker
func lexLineStartMarker(l *Lexer) stateFunc {
    return nil
}

// lexStar
func lexStar(l *Lexer) stateFunc {
    return nil
}

// lexUnderscore
func lexUnderscore(l *Lexer) stateFunc {
    return nil
}

// lexBacktick
func lexBacktick(l *Lexer) stateFunc {
    return nil
}

// lexBracket
func lexBracket(l *Lexer) stateFunc {
    return nil
}

// lexEOF
func lexEOF(l *Lexer) stateFunc {
    l.emit(TokenEOF)
    return nil
}
