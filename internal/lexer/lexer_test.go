package lexer

import "testing"

func assertTokenType(t *testing.T, token Token, expectedType TokenType) {
	t.Helper()
	if token.Type != expectedType {
		t.Errorf("expected token type %v, got %v", expectedType, token.Type)
	}
}

func assertTokenValue(t *testing.T, token Token, expectedValue string) {
	t.Helper()
	if token.Value != expectedValue {
		t.Errorf("expected token literal %q, got %q", expectedValue, token.Value)
	}
}

func assertTokenPos(t *testing.T, token Token, expectedLine, expectedColumn int) {
	t.Helper()
	if token.Line != expectedLine || token.Column != expectedColumn {
		t.Errorf("expected token position %d:%d, got %d:%d", expectedLine, expectedColumn, token.Line, token.Column)
	}
}

func assertTokenAt(t *testing.T, tokens []Token, index int, expectedType TokenType, expectedValue string, expectedLine, expectedColumn int) {
	t.Helper()
	if index >= len(tokens) {
		t.Fatalf("expected token at index %d, but tokens length is %d", index, len(tokens))
	}
	assertTokenType(t, tokens[index], expectedType)
	if expectedValue != "" {
		assertTokenValue(t, tokens[index], expectedValue)
	}
	assertTokenPos(t, tokens[index], expectedLine, expectedColumn)
}

func TestTokenTypeString(t *testing.T) {
	tests := []struct {
		value    TokenType
		expected string
	}{
		{TokenEOF, "TokenEOF"}, {TokenError, "TokenError"}, {TokenText, "TokenText"}, {TokenSpace, "TokenSpace"},
		{TokenNewline, "TokenNewline"}, {TokenHeader, "TokenHeader"}, {TokenCodeFence, "TokenCodeFence"},
		{TokenHorizontalRule, "TokenHorizontalRule"}, {TokenBlockquote, "TokenBlockquote"},
		{TokenListMarker, "TokenListMarker"}, {TokenBacktick, "TokenBacktick"}, {TokenUnderscore, "TokenUnderscore"},
		{TokenStar, "TokenStar"},
	}

	for _, tt := range tests {
		if got := tt.value.String(); got != tt.expected {
			t.Errorf("%v.String() = %q, want %q", tt.value, got, tt.expected)
		}
	}

	unknown := TokenType(9999)
	expected := "TokenType(9999)"
	if got := unknown.String(); got != expected {
		t.Errorf("Unexpected String() for unknown value: got %q, want %q", got, expected)
	}
}

func TestAdvanceTooFar(t *testing.T) {
	input := "markee"
	l := New(input)
	l.advanceUntil(isEOF)

	if r := l.advance(); r != 0 {
		t.Errorf("Expected rune %c, found %c", 0, r)
	}
	if r := l.advance(); r != 0 {
		t.Errorf("Expected rune %c, found %c", 0, r)
	}
}

func TestPeekAheadTooFar(t *testing.T) {
	input := "markee"
	l := New(input)

	if s := l.peekAhead(5); s != "marke" {
		t.Errorf("Expected string %s, found %s", "marke", s)
	}
	if s := l.peekAhead(6); s != "markee" {
		t.Errorf("Expected string %s, found %s", "markee", s)
	}
	if s := l.peekAhead(7); s != "markee" {
		t.Errorf("Expected string %s, found %s", "markee", s)
	}
}

func TestDispatchInvalidDelimiter(t *testing.T) {
	if dispatchInlineDelimiter('a') == nil {
		t.Errorf("Unexpected nil result from dispatch.")
	}
}

func TestEOF(t *testing.T) {
	input := ""
	tokens := Tokenize(input)
	assertTokenAt(t, tokens, 0, TokenEOF, "", 1, 1)
}

func TestEOFAgain(t *testing.T) {
	input := ""
	l := New(input)
	tokens := l.All()
	assertTokenAt(t, tokens, 0, TokenEOF, "", 1, 1)

	tok := l.Next()
	assertTokenType(t, tok, TokenEOF)
	tok = l.Next()
	assertTokenType(t, tok, TokenEOF)
}

func TestNewline(t *testing.T) {
	input := "\n"
	tokens := Tokenize(input)
	assertTokenAt(t, tokens, 0, TokenNewline, "\n", 2, 0)
	assertTokenAt(t, tokens, 1, TokenEOF, "", 2, 1)
}

func TestMultipleNewlines(t *testing.T) {
	input := "\n\n\n"
	tokens := Tokenize(input)
	assertTokenAt(t, tokens, 0, TokenNewline, "\n", 2, 0)
	assertTokenAt(t, tokens, 1, TokenNewline, "\n", 3, 0)
	assertTokenAt(t, tokens, 2, TokenNewline, "\n", 4, 0)
	assertTokenAt(t, tokens, 3, TokenEOF, "", 4, 1)
}

func TestBasicText(t *testing.T) {
	input := "markee"
	tokens := Tokenize(input)
	assertTokenAt(t, tokens, 0, TokenText, "markee", 1, 1)
	assertTokenAt(t, tokens, 1, TokenEOF, "", 1, 7)
}

func TestMultiWordText(t *testing.T) {
	input := "markee is the best"
	tokens := Tokenize(input)
	assertTokenAt(t, tokens, 0, TokenText, "markee is the best", 1, 1)
	assertTokenAt(t, tokens, 1, TokenEOF, "", 1, 19)
}

func TestMultiLineText(t *testing.T) {
	input := "markee\nis\nthe\nbest"
	tokens := Tokenize(input)
	assertTokenAt(t, tokens, 0, TokenText, "markee", 1, 1)
	assertTokenAt(t, tokens, 2, TokenText, "is", 2, 1)
	assertTokenAt(t, tokens, 4, TokenText, "the", 3, 1)
	assertTokenAt(t, tokens, 6, TokenText, "best", 4, 1)
	assertTokenAt(t, tokens, 1, TokenNewline, "\n", 2, 0)
	assertTokenAt(t, tokens, 3, TokenNewline, "\n", 3, 0)
	assertTokenAt(t, tokens, 5, TokenNewline, "\n", 4, 0)
	assertTokenAt(t, tokens, 7, TokenEOF, "", 4, 5)
}

func TestHeader(t *testing.T) {
	input := "# Header 1"
	tokens := Tokenize(input)
	assertTokenAt(t, tokens, 0, TokenHeader, "#", 1, 1)
	assertTokenAt(t, tokens, 1, TokenText, "Header 1", 1, 3)
	assertTokenAt(t, tokens, 2, TokenEOF, "", 1, 11)
}

func TestMultiLevelHeader(t *testing.T) {
	input := `# Header 1
## Header 2
### Header 3
#### Header 4
##### Header 5
###### Header 6`
	tokens := Tokenize(input)
	assertTokenAt(t, tokens, 0, TokenHeader, "#", 1, 1)
	assertTokenAt(t, tokens, 3, TokenHeader, "##", 2, 1)
	assertTokenAt(t, tokens, 6, TokenHeader, "###", 3, 1)
	assertTokenAt(t, tokens, 9, TokenHeader, "####", 4, 1)
	assertTokenAt(t, tokens, 12, TokenHeader, "#####", 5, 1)
	assertTokenAt(t, tokens, 15, TokenHeader, "######", 6, 1)
	assertTokenAt(t, tokens, 1, TokenText, "Header 1", 1, 3)
	assertTokenAt(t, tokens, 4, TokenText, "Header 2", 2, 4)
	assertTokenAt(t, tokens, 7, TokenText, "Header 3", 3, 5)
	assertTokenAt(t, tokens, 10, TokenText, "Header 4", 4, 6)
	assertTokenAt(t, tokens, 13, TokenText, "Header 5", 5, 7)
	assertTokenAt(t, tokens, 16, TokenText, "Header 6", 6, 8)
	assertTokenAt(t, tokens, 2, TokenNewline, "\n", 2, 0)
	assertTokenAt(t, tokens, 5, TokenNewline, "\n", 3, 0)
	assertTokenAt(t, tokens, 8, TokenNewline, "\n", 4, 0)
	assertTokenAt(t, tokens, 11, TokenNewline, "\n", 5, 0)
	assertTokenAt(t, tokens, 14, TokenNewline, "\n", 6, 0)
	assertTokenAt(t, tokens, 17, TokenEOF, "", 6, 16)
}

func TestHeaderNoSpace(t *testing.T) {
	input := "#Header 1"
	tokens := Tokenize(input)
	assertTokenAt(t, tokens, 0, TokenText, "#Header 1", 1, 1)
	assertTokenAt(t, tokens, 1, TokenEOF, "", 1, 10)
}

func TestHeaderTooManyLevels(t *testing.T) {
	input := "####### Header 7"
	tokens := Tokenize(input)
	assertTokenAt(t, tokens, 0, TokenText, "####### Header 7", 1, 1)
	assertTokenAt(t, tokens, 1, TokenEOF, "", 1, 17)
}

func TestBlockquote(t *testing.T) {
	input := "> Quote"
	tokens := Tokenize(input)
	assertTokenAt(t, tokens, 0, TokenBlockquote, ">", 1, 1)
	assertTokenAt(t, tokens, 1, TokenText, "Quote", 1, 3)
	assertTokenAt(t, tokens, 2, TokenEOF, "", 1, 8)
}

func TestBlockquoteNoSpace(t *testing.T) {
	input := ">Quote"
	tokens := Tokenize(input)
	assertTokenAt(t, tokens, 0, TokenBlockquote, ">", 1, 1)
	assertTokenAt(t, tokens, 1, TokenText, "Quote", 1, 2)
	assertTokenAt(t, tokens, 2, TokenEOF, "", 1, 7)
}

func TestMultiLineBlockquote(t *testing.T) {
	input := `> This is a quote
> continued
> ...`
	tokens := Tokenize(input)
	assertTokenAt(t, tokens, 0, TokenBlockquote, ">", 1, 1)
	assertTokenAt(t, tokens, 3, TokenBlockquote, ">", 2, 1)
	assertTokenAt(t, tokens, 6, TokenBlockquote, ">", 3, 1)
	assertTokenAt(t, tokens, 1, TokenText, "This is a quote", 1, 3)
	assertTokenAt(t, tokens, 4, TokenText, "continued", 2, 3)
	assertTokenAt(t, tokens, 7, TokenText, "...", 3, 3)
	assertTokenAt(t, tokens, 2, TokenNewline, "\n", 2, 0)
	assertTokenAt(t, tokens, 5, TokenNewline, "\n", 3, 0)
	assertTokenAt(t, tokens, 8, TokenEOF, "", 3, 6)
}

func TestCodeBlock(t *testing.T) {
	input := `~~~
#include <stdio.h>

int main() {
    printf("Hello, World!\n");
}
~~~`
	tokens := Tokenize(input)
	assertTokenAt(t, tokens, 0, TokenCodeFence, "~~~", 1, 1)
	assertTokenAt(t, tokens, 1, TokenNewline, "\n", 2, 0)
	assertTokenAt(t, tokens, 2, TokenText, "#include <stdio.h>", 2, 1)
	assertTokenAt(t, tokens, 3, TokenNewline, "\n", 3, 0)
	assertTokenAt(t, tokens, 4, TokenNewline, "\n", 4, 0)
	assertTokenAt(t, tokens, 5, TokenText, "int main() {", 4, 1)
	assertTokenAt(t, tokens, 6, TokenNewline, "\n", 5, 0)
	assertTokenAt(t, tokens, 7, TokenText, `    printf("Hello, World!\n");`, 5, 1)
	assertTokenAt(t, tokens, 8, TokenNewline, "\n", 6, 0)
	assertTokenAt(t, tokens, 9, TokenText, "}", 6, 1)
	assertTokenAt(t, tokens, 10, TokenNewline, "\n", 7, 0)
	assertTokenAt(t, tokens, 11, TokenCodeFence, "~~~", 7, 1)
	assertTokenAt(t, tokens, 12, TokenEOF, "", 7, 4)
}

func TestEmptyCodeBlock(t *testing.T) {
	input := `~~~
~~~`
	tokens := Tokenize(input)
	assertTokenAt(t, tokens, 0, TokenCodeFence, "~~~", 1, 1)
	assertTokenAt(t, tokens, 2, TokenCodeFence, "~~~", 2, 1)
	assertTokenAt(t, tokens, 1, TokenNewline, "\n", 2, 0)
	assertTokenAt(t, tokens, 3, TokenEOF, "", 2, 4)
}

func TestBacktickCodeBlock(t *testing.T) {
	input := "```\n#include <stdio.h>\n\nint main() {\n    printf(\"Hello, World!\\n\");\n}\n```"
	tokens := Tokenize(input)
	assertTokenAt(t, tokens, 0, TokenCodeFence, "```", 1, 1)
	assertTokenAt(t, tokens, 1, TokenNewline, "\n", 2, 0)
	assertTokenAt(t, tokens, 2, TokenText, "#include <stdio.h>", 2, 1)
	assertTokenAt(t, tokens, 3, TokenNewline, "\n", 3, 0)
	assertTokenAt(t, tokens, 4, TokenNewline, "\n", 4, 0)
	assertTokenAt(t, tokens, 5, TokenText, "int main() {", 4, 1)
	assertTokenAt(t, tokens, 6, TokenNewline, "\n", 5, 0)
	assertTokenAt(t, tokens, 7, TokenText, "    printf(\"Hello, World!\\n\");", 5, 1)
	assertTokenAt(t, tokens, 8, TokenNewline, "\n", 6, 0)
	assertTokenAt(t, tokens, 9, TokenText, "}", 6, 1)
	assertTokenAt(t, tokens, 10, TokenNewline, "\n", 7, 0)
	assertTokenAt(t, tokens, 11, TokenCodeFence, "```", 7, 1)
	assertTokenAt(t, tokens, 12, TokenEOF, "", 7, 4)
}

func TestIncompleteCodeBlock(t *testing.T) {
	input := `~~~
code`
	tokens := Tokenize(input)
	assertTokenAt(t, tokens, 0, TokenCodeFence, "~~~", 1, 1)
	assertTokenAt(t, tokens, 1, TokenNewline, "\n", 2, 0)
	assertTokenAt(t, tokens, 2, TokenText, "code", 2, 1)
	assertTokenAt(t, tokens, 3, TokenEOF, "", 2, 5)
}

func TestHorizontalRule(t *testing.T) {
	input := `Hello
---
World`
	tokens := Tokenize(input)
	assertTokenAt(t, tokens, 0, TokenText, "Hello", 1, 1)
	assertTokenAt(t, tokens, 1, TokenNewline, "\n", 2, 0)
	assertTokenAt(t, tokens, 2, TokenHorizontalRule, "---", 2, 1)
	assertTokenAt(t, tokens, 3, TokenNewline, "\n", 3, 0)
	assertTokenAt(t, tokens, 4, TokenText, "World", 3, 1)
	assertTokenAt(t, tokens, 5, TokenEOF, "", 3, 6)
}

func TestStarHorizontalRule(t *testing.T) {
	input := `Hello
***
World`
	tokens := Tokenize(input)
	assertTokenAt(t, tokens, 0, TokenText, "Hello", 1, 1)
	assertTokenAt(t, tokens, 1, TokenNewline, "\n", 2, 0)
	assertTokenAt(t, tokens, 2, TokenHorizontalRule, "***", 2, 1)
	assertTokenAt(t, tokens, 3, TokenNewline, "\n", 3, 0)
	assertTokenAt(t, tokens, 4, TokenText, "World", 3, 1)
	assertTokenAt(t, tokens, 5, TokenEOF, "", 3, 6)
}

func TestListMarker(t *testing.T) {
	input := "- Item 1"
	tokens := Tokenize(input)
	assertTokenAt(t, tokens, 0, TokenListMarker, "-", 1, 1)
	assertTokenAt(t, tokens, 1, TokenText, "Item 1", 1, 3)
	assertTokenAt(t, tokens, 2, TokenEOF, "", 1, 9)
}

func TestMultipleListItems(t *testing.T) {
	input := `- First
* Second
+ Third`
	tokens := Tokenize(input)
	assertTokenAt(t, tokens, 0, TokenListMarker, "-", 1, 1)
	assertTokenAt(t, tokens, 1, TokenText, "First", 1, 3)
	assertTokenAt(t, tokens, 2, TokenNewline, "\n", 2, 0)
	assertTokenAt(t, tokens, 3, TokenListMarker, "*", 2, 1)
	assertTokenAt(t, tokens, 4, TokenText, "Second", 2, 3)
	assertTokenAt(t, tokens, 5, TokenNewline, "\n", 3, 0)
	assertTokenAt(t, tokens, 6, TokenListMarker, "+", 3, 1)
	assertTokenAt(t, tokens, 7, TokenText, "Third", 3, 3)
	assertTokenAt(t, tokens, 8, TokenEOF, "", 3, 8)
}

func TestListMarkerNoSpace(t *testing.T) {
	input := "-item"
	tokens := Tokenize(input)
	assertTokenAt(t, tokens, 0, TokenText, "-item", 1, 1)
	assertTokenAt(t, tokens, 1, TokenEOF, "", 1, 6)
}

func TestEmphasisStar(t *testing.T) {
	input := "*italic*"
	tokens := Tokenize(input)
	assertTokenAt(t, tokens, 0, TokenStar, "*", 1, 1)
	assertTokenAt(t, tokens, 1, TokenText, "italic", 1, 2)
	assertTokenAt(t, tokens, 2, TokenStar, "*", 1, 8)
	assertTokenAt(t, tokens, 3, TokenEOF, "", 1, 9)
}

func TestEmphasisUnderscore(t *testing.T) {
	input := "_italic_"
	tokens := Tokenize(input)
	assertTokenAt(t, tokens, 0, TokenUnderscore, "_", 1, 1)
	assertTokenAt(t, tokens, 1, TokenText, "italic", 1, 2)
	assertTokenAt(t, tokens, 2, TokenUnderscore, "_", 1, 8)
	assertTokenAt(t, tokens, 3, TokenEOF, "", 1, 9)
}

func TestStrongStar(t *testing.T) {
	input := "**bold**"
	tokens := Tokenize(input)
	assertTokenAt(t, tokens, 0, TokenStar, "**", 1, 1)
	assertTokenAt(t, tokens, 1, TokenText, "bold", 1, 3)
	assertTokenAt(t, tokens, 2, TokenStar, "**", 1, 7)
	assertTokenAt(t, tokens, 3, TokenEOF, "", 1, 9)
}

func TestStrongUnderscore(t *testing.T) {
	input := "__bold__"
	tokens := Tokenize(input)
	assertTokenAt(t, tokens, 0, TokenUnderscore, "__", 1, 1)
	assertTokenAt(t, tokens, 1, TokenText, "bold", 1, 3)
	assertTokenAt(t, tokens, 2, TokenUnderscore, "__", 1, 7)
	assertTokenAt(t, tokens, 3, TokenEOF, "", 1, 9)
}

func TestMixedEmphasis(t *testing.T) {
	input := "Some *italic* and **bold** text"
	tokens := Tokenize(input)
	assertTokenAt(t, tokens, 0, TokenText, "Some ", 1, 1)
	assertTokenAt(t, tokens, 1, TokenStar, "*", 1, 6)
	assertTokenAt(t, tokens, 2, TokenText, "italic", 1, 7)
	assertTokenAt(t, tokens, 3, TokenStar, "*", 1, 13)
	assertTokenAt(t, tokens, 4, TokenText, " and ", 1, 14)
	assertTokenAt(t, tokens, 5, TokenStar, "**", 1, 19)
	assertTokenAt(t, tokens, 6, TokenText, "bold", 1, 21)
	assertTokenAt(t, tokens, 7, TokenStar, "**", 1, 25)
	assertTokenAt(t, tokens, 8, TokenText, " text", 1, 27)
	assertTokenAt(t, tokens, 9, TokenEOF, "", 1, 32)
}

func TestBacktick(t *testing.T) {
	input := "This is `code`."
	tokens := Tokenize(input)
	assertTokenAt(t, tokens, 0, TokenText, "This is ", 1, 1)
	assertTokenAt(t, tokens, 1, TokenBacktick, "`", 1, 9)
	assertTokenAt(t, tokens, 2, TokenText, "code", 1, 10)
	assertTokenAt(t, tokens, 3, TokenBacktick, "`", 1, 14)
	assertTokenAt(t, tokens, 4, TokenText, ".", 1, 15)
	assertTokenAt(t, tokens, 5, TokenEOF, "", 1, 16)
}
