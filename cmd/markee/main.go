package main

import (
    "fmt"
	"github.com/spf13/cobra"
	"markee/internal/lexer"
	"markee/internal/parser"
	"markee/internal/renderer"
	"os"
	"path/filepath"
	"strings"
)

var (
    outputFile string
    lexMode    bool
    parseMode  bool
    renderMode bool
)

var rootCmd = &cobra.Command{
    Use:   "markee [input]",
    Short: "A markdown processor",
    Long:  `Markee is a markdown processor that can lex, parse, and render markdown files.

By default, markee renders markdown to HTML and prints to stdout.
You can specify an output file to write the result instead.`,
    Args:  cobra.ExactArgs(1),
    RunE:  run,
    Example: `  # Render to HTML (stdout)
  markee input.md

  # Render to HTML file
  markee input.md -o output.html

  # Show lexer tokens
  markee input.md --lex

  # Show AST
  markee input.md --parse`,
}

func init() {
    rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "output file path")
    rootCmd.Flags().BoolVarP(&lexMode, "lex", "l", false, "print lexer tokens")
    rootCmd.Flags().BoolVarP(&parseMode, "parse", "p", false, "print AST")
    rootCmd.Flags().BoolVar(&renderMode, "render", false, "render to output format (default)")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	inputFile := args[0]

	data, err := os.ReadFile(inputFile)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file %q does not exist", inputFile)
		}
		return fmt.Errorf("error reading file %q: %w", inputFile, err)
	}

	input := string(data)

	modeCount := 0
	if lexMode {
		modeCount++
	}
	if parseMode {
		modeCount++
	}
	if renderMode {
		modeCount++
	}
	if modeCount > 1 {
		return fmt.Errorf("cannot specify multiple modes (--lex, --parse, --render)")
	}

	output := os.Stdout
	if outputFile != "" {
		f, err := os.Create(outputFile)
		if err != nil {
			return fmt.Errorf("error creating output file %q: %w", outputFile, err)
		}
		defer f.Close()
		output = f
	}

	switch {
	case lexMode:
		return runLexMode(input, inputFile, output)
	case parseMode:
		return runParseMode(input, output)
	default:
		return runRenderMode(input, outputFile, output)
	}
}

func runLexMode(input, inputFile string, output *os.File) error {
	tokens := lexer.Tokenize(input)
	for _, tok := range tokens {
		fmt.Fprintf(output, "%s:%d:%d %d %q\n", 
			inputFile, tok.Line, tok.Column, tok.Type, tok.Value)
	}
	return nil
}

func runParseMode(input string, output *os.File) error {
	ast := parser.Parse(input)
	printAST(ast, 0, output)
	return nil
}

func runRenderMode(input, outputFile string, output *os.File) error {
	ast := parser.Parse(input)

	format := "html"
	if outputFile != "" {
		ext := strings.ToLower(filepath.Ext(outputFile))
		switch ext {
		case ".html", ".htm":
			format = "html"
		case ".txt", ".text":
			format = "text"
		case ".md", ".markdown":
			format = "markdown"
		default:
			return fmt.Errorf("unsupported output format: %s (supported: .html, .txt)", ext)
		}
	}

	var result string
	switch format {
	case "html":
		result = renderer.RenderHTML(ast)
	default:
		return fmt.Errorf("unknown format: %s", format)
	}

	fmt.Fprint(output, result)
	return nil
}

func printAST(node *parser.Node, depth int, output *os.File) {
	indent := strings.Repeat("  ", depth)
	fmt.Fprintf(output, "%s%s", indent, node.Type.String())
	
	if node.Level > 0 {
		fmt.Fprintf(output, " (level %d)", node.Level)
	}
	
	if node.Value != "" {
		val := node.Value
		if len(val) > 50 {
			val = val[:47] + "..."
		}
		fmt.Fprintf(output, " %q", val)
	}
	
	fmt.Fprintln(output)
	
	for _, child := range node.Children {
		printAST(child, depth+1, output)
	}
}
