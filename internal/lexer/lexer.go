package lexer

// Context represents where we are in the document structure.
// This allows the lexer to make context-sensitive decisions about what characters mean.
type Context int

const (
    CtxLineStart Context = iota // At the beginning of a line
    CtxInline                   // In the middle of a line
    CtxCodeBlock                // Inside a fenced code block
)

func isTokenChar(char rune) bool {
    return char == 0
}

// Lexer tokenizes markdown input using a state machine pattern.
// Each state function examines the input and returns the next state.
type Lexer struct {
    input   string  // The complete markdown input string
    start   int     // Start position of current token
    pos     int     // Current position in input
    width   int     // Width of last rune read
    line    int     // Current line number (1-indexed)
    column  int     // Current column number (1-indexed)
    tokens  []Token // Collected tokens
    context Context // Current parsing context
}

// stateFunc represents a state in the lexer state machine.
// Each state function processes input according the rules of its state.
// Returns either the next state function to execute or nil (indicating lexing is complete).
type stateFunc func(*Lexer) stateFunc

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
