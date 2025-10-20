package lexer

//go:generate stringer -type=TokenType
type TokenType int

const (
	TokenEOF TokenType = iota
	TokenError

	// Content tokens - the actual text and spacing
	TokenText    // Any regular text content
	TokenNewline // Line break

	// Block-level markers (recognized only at the start of line)
	TokenHeader         // # ## ### etc - stores level in token metadata
	TokenCodeFence      // ``` or ~~~ at line start
	TokenHorizontalRule // --- or *** at line start
	TokenBlockquote     // > at line start
	TokenListMarker     // - * + for unordered, or N. N) for ordered

	// Inline delimiters (recognized mid-line)
	TokenBacktick   // `
	TokenStar       // * or **
	TokenUnderscore // _ or __
)

type Token struct {
	Type   TokenType
	Value  string
	Line   int
	Column int
}
