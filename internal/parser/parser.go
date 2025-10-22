package parser

import(
    "markee/internal/ast"
)

func Parse(input string) ast.Node {
    ctx := NewContext(input)
    ParseBlocks(ctx)
    return ctx.Doc
}
