package main

import (
    "flag"
    "fmt"
    "log"
    "markee/internal/lexer"
    "os"
)

func main() {
    var inputFile string
    var lexMode bool

    flag.StringVar(&inputFile, "i", "", "input file path")
    flag.StringVar(&inputFile, "input", "", "input file path")
    flag.BoolVar(&lexMode, "l", false, "print lexer tokens")
    flag.BoolVar(&lexMode, "lex", false, "print lexer tokens")
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
        l := lexer.New(string(data))
        for _, tok := range l.Tokenize() {
            fmt.Printf("%s:%d:%d %s %q\n", inputFile, tok.Line, tok.Column, tok.Type.String(), tok.Value)
        }
        os.Exit(0)
    }
}
