package lexer

type Context int

const (
    CtxInline Context = iota
    CtxLineStart
    CtxCodeBlock
)
