package main

import (
	"flag"
	"fmt"
	"log"
	"markee/internal/lexer"
    "markee/internal/parser"
	"os"
    "strings"
)

func main() {
	var inputFile string
	var lexMode bool
    var parseMode bool

	flag.StringVar(&inputFile, "i", "", "input file path")
	flag.StringVar(&inputFile, "input", "", "input file path")
	flag.BoolVar(&lexMode, "l", false, "print lexer tokens")
	flag.BoolVar(&lexMode, "lex", false, "print lexer tokens")
    flag.BoolVar(&parseMode, "p", false, "print AST")
    flag.BoolVar(&parseMode, "parse", false, "print AST")
	flag.Parse()

	// If no input flag is given, check the first positional argument
	if inputFile == "" && flag.NArg() > 0 {
		inputFile = flag.Arg(0)
	}

	// Early exit if no input file is given
	if inputFile == "" {
		log.Fatal("no input file provided")
	}

	data, err := os.ReadFile(inputFile)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("file %q does not exist", inputFile)
		} else {
			log.Fatalf("error checking file %q: %v", inputFile, err)
		}
	}

	if lexMode {
		for _, tok := range lexer.Tokenize(string(data)) {
			fmt.Printf("%s:%d:%d %d %q\n", inputFile, tok.Line, tok.Column, tok.Type, tok.Value)
		}
		os.Exit(0)
	}

    if parseMode {
       
        ast := parser.Parse(string(data))
        printAST(ast, 0)
        os.Exit(0)
    }
}

func printAST(node *parser.Node, depth int) {
    indent := strings.Repeat("  ", depth)

    fmt.Printf("%s%s", indent, node.Type.String())
    if node.Level > 0 {
        fmt.Printf(" (level %d)", node.Level)
    }
    if node.Value != "" {
        fmt.Printf(" %s", node.Value)
    }
    fmt.Printf("\n")

    for _, child := range node.Children {
        printAST(child, depth+1)
    }
}
