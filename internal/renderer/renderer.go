package renderer

import (
    "markee/internal/ast"
)

type Renderer interface {
    Render(ast.Document) string
}
