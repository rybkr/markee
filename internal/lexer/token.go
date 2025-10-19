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
	TokenHeader         // # ## ### etc - stores level in token metadata
	TokenCodeFence      // ``` or ~~~ at line start
	TokenHorizontalRule // --- or *** at line start
	TokenBlockquote     // > at line start
	TokenListMarker     // - * + for unordered, or N. N) for ordered

	// Inline delimiters (recognized mid-line)
	TokenBacktick   // ` ```
	TokenStar       // * ** ***
	TokenUnderscore // _ __ ___

	// Links and images
	TokenBracketOpen  // [
	TokenBracketClose // ]
	TokenParenOpen    // (
	TokenParenClose   // )
	TokenBang         // !
)

// Token represents a single lexical unit from the input.
// Line and Column are 1-indexed for compatible error messages.
type Token struct {
	Type   TokenType
	Value  string // The actual text matched
	Line   int    // Line number where token starts (1-indexed)
	Column int    // Column where token starts (1-indexed)
}

func (t TokenType) String() string {
	switch t {
	case TokenEOF:
		return "EOF"
	case TokenText:
		return "Text"
	case TokenSpace:
		return "Space"
	case TokenNewline:
		return "Newline"
	case TokenHeader:
		return "Header"
	case TokenCodeFence:
		return "CodeFence"
	case TokenHorizontalRule:
		return "HorizontalRule"
	case TokenBlockquote:
		return "Blockquote"
	case TokenListMarker:
		return "ListMarker"
	case TokenBacktick:
		return "Backtick"
	case TokenStar:
		return "Star"
	case TokenUnderscore:
		return "Underscore"
	case TokenBracketOpen:
		return "BracketOpen"
	case TokenBracketClose:
		return "BracketClose"
	case TokenParenOpen:
		return "ParenOpen"
	case TokenParenClose:
		return "ParenClose"
	case TokenBang:
		return "Bang"
	default:
		return "Error"
	}
}
