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
