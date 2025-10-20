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
