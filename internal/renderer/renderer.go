package renderer

type Renderer interface {
    Render(node *parser.Node) string

    renderChildren(node *parser.Node) string
}
