package lexer

import "testing"

func TestEOF(t *testing.T) {
    input := ""
    tokens := New(input).Tokenize()
    assertTokenAt(t, tokens, 0, TokenEOF, "", 1, 1)
}

func TestNewline(t *testing.T) {
    input := "\n"
    tokens := New(input).Tokenize()
    assertTokenAt(t, tokens, 0, TokenNewline, "\n", 2, 0)
    assertTokenAt(t, tokens, 1, TokenEOF, "", 2, 1)
}

func TestMultipleNewlines(t *testing.T) {
    input := "\n\n\n"
    tokens := New(input).Tokenize()
    assertTokenAt(t, tokens, 0, TokenNewline, "\n", 2, 0)
    assertTokenAt(t, tokens, 1, TokenNewline, "\n", 3, 0)
    assertTokenAt(t, tokens, 2, TokenNewline, "\n", 4, 0)
    assertTokenAt(t, tokens, 3, TokenEOF, "", 4, 1)
}

func TestBasicText(t *testing.T) {
    input := "markee"
    tokens := New(input).Tokenize()
    assertTokenAt(t, tokens, 0, TokenText, "markee", 1, 1)
    assertTokenAt(t, tokens, 1, TokenEOF, "", 1, 7)
}

func TestMultiWordText(t *testing.T) {
    input := "markee is the best"
    tokens := New(input).Tokenize()
    assertTokenAt(t, tokens, 0, TokenText, "markee is the best", 1, 1)
    assertTokenAt(t, tokens, 1, TokenEOF, "", 1, 19)
}

func TestMultiLineText(t *testing.T) {
    input := "markee\nis\nthe\nbest" 
    tokens := New(input).Tokenize()
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
    tokens := New(input).Tokenize()
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
    tokens := New(input).Tokenize()
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
    tokens := New(input).Tokenize()
    assertTokenAt(t, tokens, 0, TokenText, "#Header 1", 1, 1)
    assertTokenAt(t, tokens, 1, TokenEOF, "", 1, 10)
}

func TestHeaderTooManyLevels(t *testing.T) {
    input := "####### Header 7"
    tokens := New(input).Tokenize()
    assertTokenAt(t, tokens, 0, TokenText, "####### Header 7", 1, 1)
    assertTokenAt(t, tokens, 1, TokenEOF, "", 1, 17)
}

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
