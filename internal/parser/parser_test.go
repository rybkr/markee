package parser

import "testing"

func assertNodeType(t *testing.T, node *Node, expectedType NodeType) {
	t.Helper()
	if node.Type != expectedType {
		t.Errorf("expected node type %v, got %v", expectedType, node.Type)
	}
}

func assertNodeValue(t *testing.T, node *Node, expectedValue string) {
	t.Helper()
	if node.Value != expectedValue {
		t.Errorf("expected node value %q, got %q", expectedValue, node.Value)
	}
}

func assertNodeLevel(t *testing.T, node *Node, expectedLevel int) {
	t.Helper()
	if node.Level != expectedLevel {
		t.Errorf("expected node level %d, got %d", expectedLevel, node.Level)
	}
}

func assertChildCount(t *testing.T, node *Node, expectedCount int) {
	t.Helper()
	if len(node.Children) != expectedCount {
		t.Errorf("expected %d children, got %d", expectedCount, len(node.Children))
	}
}

func assertChild(t *testing.T, parent *Node, index int, expectedType NodeType) *Node {
	t.Helper()
	if index >= len(parent.Children) {
		t.Fatalf("expected child at index %d, but only %d children exist", index, len(parent.Children))
	}
	child := parent.Children[index]
	assertNodeType(t, child, expectedType)
	return child
}

func TestEmptyDocument(t *testing.T) {
	input := ""
	ast := Parse(input)
	assertNodeType(t, ast, NodeDocument)
	assertChildCount(t, ast, 0)
}

func TestSimpleText(t *testing.T) {
	input := "Hello world"
	ast := Parse(input)
	assertNodeType(t, ast, NodeDocument)
	assertChildCount(t, ast, 1)

	para := assertChild(t, ast, 0, NodeParagraph)
	assertChildCount(t, para, 1)

	text := assertChild(t, para, 0, NodeText)
	assertNodeValue(t, text, "Hello world")
}

func TestMultipleParagraphs(t *testing.T) {
	input := "First paragraph\n\nSecond paragraph"
	ast := Parse(input)
	assertNodeType(t, ast, NodeDocument)
	assertChildCount(t, ast, 2)

	para1 := assertChild(t, ast, 0, NodeParagraph)
	assertChildCount(t, para1, 1)
	text1 := assertChild(t, para1, 0, NodeText)
	assertNodeValue(t, text1, "First paragraph")

	para2 := assertChild(t, ast, 1, NodeParagraph)
	assertChildCount(t, para2, 1)
	text2 := assertChild(t, para2, 0, NodeText)
	assertNodeValue(t, text2, "Second paragraph")
}

func TestHeader(t *testing.T) {
	input := "# Heading 1"
	ast := Parse(input)
	assertNodeType(t, ast, NodeDocument)
	assertChildCount(t, ast, 1)

	header := assertChild(t, ast, 0, NodeHeader)
	assertNodeLevel(t, header, 1)
	assertChildCount(t, header, 1)

	text := assertChild(t, header, 0, NodeText)
	assertNodeValue(t, text, "Heading 1")
}

func TestMultipleLevelHeaders(t *testing.T) {
	input := `# H1
## H2
### H3`
	ast := Parse(input)
	assertNodeType(t, ast, NodeDocument)
	assertChildCount(t, ast, 3)

	h1 := assertChild(t, ast, 0, NodeHeader)
	assertNodeLevel(t, h1, 1)
	text1 := assertChild(t, h1, 0, NodeText)
	assertNodeValue(t, text1, "H1")

	h2 := assertChild(t, ast, 1, NodeHeader)
	assertNodeLevel(t, h2, 2)
	text2 := assertChild(t, h2, 0, NodeText)
	assertNodeValue(t, text2, "H2")

	h3 := assertChild(t, ast, 2, NodeHeader)
	assertNodeLevel(t, h3, 3)
	text3 := assertChild(t, h3, 0, NodeText)
	assertNodeValue(t, text3, "H3")
}

func TestBlockquote(t *testing.T) {
	input := "> This is a quote"
	ast := Parse(input)
	assertNodeType(t, ast, NodeDocument)
	assertChildCount(t, ast, 1)

	quote := assertChild(t, ast, 0, NodeBlockquote)
	assertChildCount(t, quote, 1)

	text := assertChild(t, quote, 0, NodeText)
	assertNodeValue(t, text, "This is a quote")
}

func TestMultiLineBlockquote(t *testing.T) {
	input := `> First line
> Second line`
	ast := Parse(input)
	assertNodeType(t, ast, NodeDocument)
	assertChildCount(t, ast, 2)

	quote1 := assertChild(t, ast, 0, NodeBlockquote)
	text1 := assertChild(t, quote1, 0, NodeText)
	assertNodeValue(t, text1, "First line")

	quote2 := assertChild(t, ast, 1, NodeBlockquote)
	text2 := assertChild(t, quote2, 0, NodeText)
	assertNodeValue(t, text2, "Second line")
}

func TestCodeBlock(t *testing.T) {
	input := `~~~
code here
~~~`
	ast := Parse(input)
	assertNodeType(t, ast, NodeDocument)
	assertChildCount(t, ast, 1)

	code := assertChild(t, ast, 0, NodeCodeBlock)
	assertNodeValue(t, code, "code here\n")
}
