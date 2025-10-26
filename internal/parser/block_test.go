package parser

import (
	"markee/internal/ast"
	"testing"
)

func assertNodeType(t *testing.T, node ast.Node, expectedType ast.NodeType) {
	t.Helper()
	if node.Type() != expectedType {
		t.Errorf("expected node type %v, got %v", expectedType, node.Type())
	}
}

func assertNodeTypeNot(t *testing.T, node ast.Node, notType ast.NodeType) {
	t.Helper()
	if node.Type() == notType {
		t.Errorf("expected node type not to be %v, got %v", notType, node.Type())
	}
}

func assertChildCount(t *testing.T, node ast.Node, expectedCount int) {
	t.Helper()
	if len(node.Children()) != expectedCount {
		t.Errorf("expected %d children, got %d", expectedCount, len(node.Children()))
	}
}

func assertChild(t *testing.T, parent ast.Node, index int, expectedType ast.NodeType) ast.Node {
	t.Helper()
	if index >= len(parent.Children()) {
		t.Fatalf("expected child at index %d, but only %d children exist", index, len(parent.Children()))
	}
	child := parent.Children()[index]
	assertNodeType(t, child, expectedType)
	return child
}

func assertChildNot(t *testing.T, parent ast.Node, index int, expectedType ast.NodeType) ast.Node {
	t.Helper()
	if index >= len(parent.Children()) {
		t.Fatalf("expected child at index %d, but only %d children exist", index, len(parent.Children()))
	}
	child := parent.Children()[index]
	assertNodeTypeNot(t, child, expectedType)
	return child
}

func assertContent(t *testing.T, node ast.Node, expectedLiteral string) {
	t.Helper()
	if content, ok := node.(*ast.Content); ok {
		if content.Literal != expectedLiteral {
			t.Errorf("expected node literal %q, got %q", expectedLiteral, content.Literal)
		}
	} else {
		assertNodeType(t, node, ast.NodeContent)
	}
}

func TestEmptyInput(t *testing.T) {
	input := ""
	doc := Parse(input)
	assertNodeType(t, doc, ast.NodeDocument)
	assertChildCount(t, doc, 0)
}

func TestATXHeading(t *testing.T) {
	input := `# foo
## foo
### foo
#### foo
##### foo
###### foo`
	doc := Parse(input)
	assertNodeType(t, doc, ast.NodeDocument)
	assertChildCount(t, doc, 6)

	h1 := assertChild(t, doc, 0, ast.NodeHeading)
	assertChildCount(t, h1, 1)
	h1Content := assertChild(t, h1, 0, ast.NodeContent)
	assertContent(t, h1Content, "foo")

	h2 := assertChild(t, doc, 1, ast.NodeHeading)
	assertChildCount(t, h2, 1)
	h2Content := assertChild(t, h2, 0, ast.NodeContent)
	assertContent(t, h2Content, "foo")

	h3 := assertChild(t, doc, 2, ast.NodeHeading)
	assertChildCount(t, h3, 1)
	h3Content := assertChild(t, h3, 0, ast.NodeContent)
	assertContent(t, h3Content, "foo")

	h4 := assertChild(t, doc, 3, ast.NodeHeading)
	assertChildCount(t, h4, 1)
	h4Content := assertChild(t, h4, 0, ast.NodeContent)
	assertContent(t, h4Content, "foo")

	h5 := assertChild(t, doc, 4, ast.NodeHeading)
	assertChildCount(t, h5, 1)
	h5Content := assertChild(t, h5, 0, ast.NodeContent)
	assertContent(t, h5Content, "foo")

	h6 := assertChild(t, doc, 5, ast.NodeHeading)
	assertChildCount(t, h6, 1)
	h6Content := assertChild(t, h6, 0, ast.NodeContent)
	assertContent(t, h6Content, "foo")
}

func TestATXHeadingEmpty(t *testing.T) {
	input := `#
##
###
####
#####
######`
	doc := Parse(input)
	assertNodeType(t, doc, ast.NodeDocument)
	assertChildCount(t, doc, 6)

	h1 := assertChild(t, doc, 0, ast.NodeHeading)
	assertChildCount(t, h1, 0)

	h2 := assertChild(t, doc, 1, ast.NodeHeading)
	assertChildCount(t, h2, 0)

	h3 := assertChild(t, doc, 2, ast.NodeHeading)
	assertChildCount(t, h3, 0)

	h4 := assertChild(t, doc, 3, ast.NodeHeading)
	assertChildCount(t, h4, 0)

	h5 := assertChild(t, doc, 4, ast.NodeHeading)
	assertChildCount(t, h5, 0)

	h6 := assertChild(t, doc, 5, ast.NodeHeading)
	assertChildCount(t, h6, 0)
}

func TestATXHeadingTooLong(t *testing.T) {
	input := "####### foo"
	doc := Parse(input)
	assertNodeType(t, doc, ast.NodeDocument)
	assertChildCount(t, doc, 1)
	paragraph := assertChild(t, doc, 0, ast.NodeParagraph)
	assertChildCount(t, paragraph, 1)
	content := assertChild(t, paragraph, 0, ast.NodeContent)
	assertChildCount(t, content, 0)
	assertContent(t, content, "####### foo")
}

func TestATXHeadingNoSpace(t *testing.T) {
	input := "#hashtag"
	doc := Parse(input)
	assertNodeType(t, doc, ast.NodeDocument)
	assertChildCount(t, doc, 1)
	paragraph := assertChild(t, doc, 0, ast.NodeParagraph)
	assertChildCount(t, paragraph, 1)
	content := assertChild(t, paragraph, 0, ast.NodeContent)
	assertChildCount(t, content, 0)
	assertContent(t, content, "#hashtag")
}

func TestATXHeadingExtraSpace(t *testing.T) {
	input := "#                 foo                     "
	doc := Parse(input)
	assertNodeType(t, doc, ast.NodeDocument)
	assertChildCount(t, doc, 1)

	h1 := assertChild(t, doc, 0, ast.NodeHeading)
	assertChildCount(t, h1, 1)
	h1Content := assertChild(t, h1, 0, ast.NodeContent)
	assertContent(t, h1Content, "foo")
}

func TestATXHeadingIndent(t *testing.T) {
	input := `   # foo
  ## foo
 ### foo
   #### foo
  ##### foo
 ###### foo`
	doc := Parse(input)
	assertNodeType(t, doc, ast.NodeDocument)
	assertChildCount(t, doc, 6)

	h1 := assertChild(t, doc, 0, ast.NodeHeading)
	assertChildCount(t, h1, 1)
	h1Content := assertChild(t, h1, 0, ast.NodeContent)
	assertContent(t, h1Content, "foo")

	h2 := assertChild(t, doc, 1, ast.NodeHeading)
	assertChildCount(t, h2, 1)
	h2Content := assertChild(t, h2, 0, ast.NodeContent)
	assertContent(t, h2Content, "foo")

	h3 := assertChild(t, doc, 2, ast.NodeHeading)
	assertChildCount(t, h3, 1)
	h3Content := assertChild(t, h3, 0, ast.NodeContent)
	assertContent(t, h3Content, "foo")

	h4 := assertChild(t, doc, 3, ast.NodeHeading)
	assertChildCount(t, h4, 1)
	h4Content := assertChild(t, h4, 0, ast.NodeContent)
	assertContent(t, h4Content, "foo")

	h5 := assertChild(t, doc, 4, ast.NodeHeading)
	assertChildCount(t, h5, 1)
	h5Content := assertChild(t, h5, 0, ast.NodeContent)
	assertContent(t, h5Content, "foo")

	h6 := assertChild(t, doc, 5, ast.NodeHeading)
	assertChildCount(t, h6, 1)
	h6Content := assertChild(t, h6, 0, ast.NodeContent)
	assertContent(t, h6Content, "foo")
}

func TestATXHeadingTooMuchIndent(t *testing.T) {
	input := "    # foo"
	doc := Parse(input)
	assertNodeType(t, doc, ast.NodeDocument)
	assertChildCount(t, doc, 1)
	assertChildNot(t, doc, 0, ast.NodeHeading)
}

func TestATXHeadingClosingSequence(t *testing.T) {
	input := `# foo #
## foo ##
### foo #
#### foo ######################
##### foo #####
###### foo ######        `
	doc := Parse(input)
	assertNodeType(t, doc, ast.NodeDocument)
	assertChildCount(t, doc, 6)

	h1 := assertChild(t, doc, 0, ast.NodeHeading)
	assertChildCount(t, h1, 1)
	h1Content := assertChild(t, h1, 0, ast.NodeContent)
	assertContent(t, h1Content, "foo")

	h2 := assertChild(t, doc, 1, ast.NodeHeading)
	assertChildCount(t, h2, 1)
	h2Content := assertChild(t, h2, 0, ast.NodeContent)
	assertContent(t, h2Content, "foo")

	h3 := assertChild(t, doc, 2, ast.NodeHeading)
	assertChildCount(t, h3, 1)
	h3Content := assertChild(t, h3, 0, ast.NodeContent)
	assertContent(t, h3Content, "foo")

	h4 := assertChild(t, doc, 3, ast.NodeHeading)
	assertChildCount(t, h4, 1)
	h4Content := assertChild(t, h4, 0, ast.NodeContent)
	assertContent(t, h4Content, "foo")

	h5 := assertChild(t, doc, 4, ast.NodeHeading)
	assertChildCount(t, h5, 1)
	h5Content := assertChild(t, h5, 0, ast.NodeContent)
	assertContent(t, h5Content, "foo")

	h6 := assertChild(t, doc, 5, ast.NodeHeading)
	assertChildCount(t, h6, 1)
	h6Content := assertChild(t, h6, 0, ast.NodeContent)
	assertContent(t, h6Content, "foo")
}

func TestATXHeadingClosingSequenceEmpty(t *testing.T) {
	input := `# #
## ###
### #
#### ################
##### #
######      #######       `
	doc := Parse(input)
	assertNodeType(t, doc, ast.NodeDocument)
	assertChildCount(t, doc, 6)

	h1 := assertChild(t, doc, 0, ast.NodeHeading)
	assertChildCount(t, h1, 0)

	h2 := assertChild(t, doc, 1, ast.NodeHeading)
	assertChildCount(t, h2, 0)

	h3 := assertChild(t, doc, 2, ast.NodeHeading)
	assertChildCount(t, h3, 0)

	h4 := assertChild(t, doc, 3, ast.NodeHeading)
	assertChildCount(t, h4, 0)

	h5 := assertChild(t, doc, 4, ast.NodeHeading)
	assertChildCount(t, h5, 0)

	h6 := assertChild(t, doc, 5, ast.NodeHeading)
	assertChildCount(t, h6, 0)
}

func TestATXHeadingTrickClosingSequence(t *testing.T) {
	input := "# foo # foo"
	doc := Parse(input)
	assertNodeType(t, doc, ast.NodeDocument)
	assertChildCount(t, doc, 1)
	h1 := assertChild(t, doc, 0, ast.NodeHeading)
	assertChildCount(t, h1, 1)
	h1Content := assertChild(t, h1, 0, ast.NodeContent)
	assertContent(t, h1Content, "foo # foo")
}

func TestThematicBreaks(t *testing.T) {
	input := `***
---
___`
	doc := Parse(input)
	assertNodeType(t, doc, ast.NodeDocument)
	assertChildCount(t, doc, 3)
	assertChild(t, doc, 0, ast.NodeThematicBreak)
	assertChild(t, doc, 1, ast.NodeThematicBreak)
	assertChild(t, doc, 2, ast.NodeThematicBreak)
}

func TestThematicBreakTooFewCharacters(t *testing.T) {
	input := `**
--
__`
	doc := Parse(input)
	assertNodeType(t, doc, ast.NodeDocument)
	assertChildNot(t, doc, 0, ast.NodeThematicBreak)
}

func TestThematicBreakIndent(t *testing.T) {
	input := `  ***
   ---
 ___`
	doc := Parse(input)
	assertNodeType(t, doc, ast.NodeDocument)
	assertChildCount(t, doc, 3)
	assertChild(t, doc, 0, ast.NodeThematicBreak)
	assertChild(t, doc, 1, ast.NodeThematicBreak)
	assertChild(t, doc, 2, ast.NodeThematicBreak)
}

func TestThematicBreakTooMuchIndent(t *testing.T) {
	input := "    ***"
	doc := Parse(input)
	assertNodeType(t, doc, ast.NodeDocument)
	assertChildNot(t, doc, 0, ast.NodeThematicBreak)
}

func TestThematicBreakMoreCharacters(t *testing.T) {
	input := `****
- - - - 
__ _               _       _           __              `
	doc := Parse(input)
	assertNodeType(t, doc, ast.NodeDocument)
	assertChild(t, doc, 0, ast.NodeThematicBreak)
	assertChild(t, doc, 1, ast.NodeThematicBreak)
	assertChild(t, doc, 2, ast.NodeThematicBreak)
}

func TestThematicBreakInvalidCharacters(t *testing.T) {
	input := "*-_"
	doc := Parse(input)
	assertNodeType(t, doc, ast.NodeDocument)
	assertChildNot(t, doc, 0, ast.NodeThematicBreak)
}

func TestThematicBreakInvalidCharacters2(t *testing.T) {
	input := "***8"
	doc := Parse(input)
	assertNodeType(t, doc, ast.NodeDocument)
	assertChildNot(t, doc, 0, ast.NodeThematicBreak)
}

func TestParagraph(t *testing.T) {
	input := "foo"
	doc := Parse(input)
	assertNodeType(t, doc, ast.NodeDocument)
	assertChildCount(t, doc, 1)
	paragraph := assertChild(t, doc, 0, ast.NodeParagraph)
	assertChildCount(t, paragraph, 1)
	content := assertChild(t, paragraph, 0, ast.NodeContent)
	assertChildCount(t, content, 0)
	assertContent(t, content, "foo")
}

func TestMultipleParagraphs(t *testing.T) {
	input := `foo

bar



baz`
	doc := Parse(input)
	assertNodeType(t, doc, ast.NodeDocument)
	assertChildCount(t, doc, 3)

	p1 := assertChild(t, doc, 0, ast.NodeParagraph)
	assertChildCount(t, p1, 1)
	c1 := assertChild(t, p1, 0, ast.NodeContent)
	assertChildCount(t, c1, 0)
	assertContent(t, c1, "foo")

	p2 := assertChild(t, doc, 1, ast.NodeParagraph)
	assertChildCount(t, p2, 1)
	c2 := assertChild(t, p2, 0, ast.NodeContent)
	assertChildCount(t, c2, 0)
	assertContent(t, c2, "bar")

	p3 := assertChild(t, doc, 2, ast.NodeParagraph)
	assertChildCount(t, p3, 1)
	c3 := assertChild(t, p3, 0, ast.NodeContent)
	assertChildCount(t, c3, 0)
	assertContent(t, c3, "baz")
}

func TestFencedCodeBlock(t *testing.T) {
    input := `~~~
<
 >
~~~`
    doc := Parse(input)
    assertNodeType(t, doc, ast.NodeDocument)
    assertChildCount(t, doc, 1)
    codeBlock := assertChild(t, doc, 0, ast.NodeCodeBlock)
    assertChildCount(t, codeBlock, 2)
    line1 := assertChild(t, codeBlock, 0, ast.NodeContent)
    assertContent(t, line1, "<")
    line2 := assertChild(t, codeBlock, 1, ast.NodeContent)
    assertContent(t, line2, " >")
}

func TestFencedCodeBlockBacktick(t *testing.T) {
    input := "```\n<\n >\n```"
    doc := Parse(input)
    assertNodeType(t, doc, ast.NodeDocument)
    assertChildCount(t, doc, 1)
    codeBlock := assertChild(t, doc, 0, ast.NodeCodeBlock)
    assertChildCount(t, codeBlock, 2)
    line1 := assertChild(t, codeBlock, 0, ast.NodeContent)
    assertContent(t, line1, "<")
    line2 := assertChild(t, codeBlock, 1, ast.NodeContent)
    assertContent(t, line2, " >")
}
