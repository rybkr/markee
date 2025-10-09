package lexer

// TokenType represents the type of token recognized by the lexer.
type TokenType int

const (
    // Special tokens
	TokenEOF TokenType = iota
    TokenError

    // Content tokens - the actual text and spacing
	TokenText    // Any regular text content
    TokenSpace   // Horizontal whitespace
    TokenNewline // Line break

    // Block-level markers (recognized only at the start of line)
    TokenHeader     // # ## ### etc - stores level in token metadata
    TokenCodeFence  // ``` or ~~~ at line start
    TokenBlockquote // > at line start
    TokenListMarker // - * + for unordered, or N. N) for ordered

    // Inline delimiters (recognized mid-line)
    TokenBacktick   // ` ```
    TokenStar       // * ** ***
    TokenUnderscore // _ __ ___

    // Links and images
    TokenBracketOpen  // [
    TokenBracketClode // ]
    TokenParenOpen    // (
    TokenParenClose   // )
    TokenBang         // !
)

// Token represents a single lexical unit from the input.
// Line and Column are 1-indexed for compatible error messages.
type Token struct {
	Type     TokenType
	Value    string    // The actual text matched
    Line     int       // Line number where token starts (1-indexed)
    Column   int       // Column where token starts (1-indexed)
}

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

func (t TokenType) String() string {
    switch t {
    case EOF:
        return "EOF"
    case TEXT:
        return "TEXT"
    default:
        return "UNKNOWN"
    }
}
