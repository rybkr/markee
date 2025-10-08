package lexer

type TokenType int

const (
	EOF TokenType = iota

	TEXT
)

type Token struct {
	Type     TokenType
	Value    string
    Line     int
    Column   int
}

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
