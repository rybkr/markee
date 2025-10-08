package lexer

type Lexer struct {
    input    string
    position int
}

func New(input string) *Lexer {
    return &Lexer{
        input: input,
        position: 0,
    }
}

func (l *Lexer) Tokenize() []Token {
    var tokens []Token

    for l.position < len(l.input) {
        tokens = append(tokens, *l.nextToken())
    }

    return tokens
}

func (l *Lexer) nextToken() *Token {
    readChar := l.peek()

    switch {
    case readChar == 0:
        l.advance()
        return &Token{Type: EOF, Value: ""}

    default:
        text :=  ""
        for !isTokenChar(l.peek()) {
            text += string(l.advance())
        }
        return &Token{Type: TEXT, Value: text}
    }
}

func (l *Lexer) peek() rune {
    if l.position >= len(l.input) {
        return 0
    }
    return rune(l.input[l.position])
}

func (l *Lexer) advance() rune {
    char := l.peek()
    l.position++
    return char
}
