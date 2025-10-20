package lexer

//go:generate stringer -type=TokenType
type TokenType int

const (
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
	TokenBacktick // `
	TokenEmphasis // * or _
	TokenStrong   // ** or __

	// Links and images
	TokenBracketOpen  // [
	TokenBracketClose // ]
	TokenParenOpen    // (
	TokenParenClose   // )
	TokenBang         // !
)

type Token struct {
	Type   TokenType
	Value  string
	Line   int
	Column int
}
