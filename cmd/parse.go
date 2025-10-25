package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"markee/internal/ast"
	"markee/internal/parser"
	"markee/internal/renderer"
)

var parseCmd = &cobra.Command{
	Use:   "parse [file]",
	Short: "Parse markdown input and display AST",
	Long:  "Parse markdown from a file or stdin and display the resulting AST structure.",
	Args:  cobra.MaximumNArgs(1),
	Run:   runParse,
}

func runParse(cmd *cobra.Command, args []string) {
	var input string

	if len(args) == 1 {
		// Read from file
		content, err := os.ReadFile(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			os.Exit(1)
		}
		input = string(content)
	} else {
		// Read from stdin
		content, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
			os.Exit(1)
		}
		input = string(content)
	}

	doc := parser.Parse(input)
    html := renderer.RenderHTML(doc)
    fmt.Println(html)
}

func printTree(node ast.Node, depth int) {
	indent := strings.Repeat("  ", depth)

	switch n := node.(type) {
	case *ast.Document:
		fmt.Printf("%s[Document]\n", indent)
	case *ast.Heading:
		fmt.Printf("%s[Heading level=%d]\n", indent, n.Level)
	case *ast.Paragraph:
		fmt.Printf("%s[Paragraph]\n", indent)
	case *ast.BlockQuote:
		fmt.Printf("%s[BlockQuote]\n", indent)
	case *ast.CodeBlock:
		ftype := "indented"
		if n.IsFenced {
			ftype = "fenced"
		}
		fmt.Printf("%s[CodeBlock %s info=%q]\n", indent, ftype, n.Language)
	case *ast.ThematicBreak:
		fmt.Printf("%s[ThematicBreak]\n", indent)
	case *ast.List:
		ltype := "unordered"
		if n.IsOrdered {
			ltype = "ordered"
		}
		fmt.Printf("%s[List %s tight=%v]\n", indent, ltype, n.IsTight)
	case *ast.ListItem:
		fmt.Printf("%s[ListItem]\n", indent)
	case *ast.Content:
		fmt.Printf("%s[Text] %q\n", indent, n.Literal)
	case *ast.Emphasis:
		fmt.Printf("%s[Emphasis]\n", indent)
	case *ast.Strong:
		fmt.Printf("%s[Strong]\n", indent)
	case *ast.CodeSpan:
		fmt.Printf("%s[Code] %q\n", indent, n.Literal)
	case *ast.Link:
		fmt.Printf("%s[Link dest=%q title=%q]\n", indent, n.Destination, n.Title)
	case *ast.Image:
		fmt.Printf("%s[Image dest=%q title=%q alt=%q]\n", indent, n.Destination, n.Title, n.AltText)
	case *ast.LineBreak:
		fmt.Printf("%s[LineBreak]\n", indent)
	case *ast.SoftBreak:
		fmt.Printf("%s[SoftBreak]\n", indent)
	case *ast.HTMLBlock:
		fmt.Printf("%s[HTMLBlock] %q\n", indent, n.Literal)
	case *ast.HTMLSpan:
		fmt.Printf("%s[HTMLInline] %q\n", indent, n.Literal)
	default:
		fmt.Printf("%s[Unknown: %T]\n", indent, node)
	}

	for _, child := range node.Children() {
		printTree(child, depth+1)
	}
}
