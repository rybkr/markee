package main

import (
    "flag"
    "fmt"
    "os"
)

func main() {
    var inputFile string

    flag.StringVar(&inputFile, "i", "", "input file path")
    flag.StringVar(&inputFile, "input", "", "input file path")
    flag.Parse()

    // If no input flag is given, check the first positional argument
    if inputFile == "" && flag.NArg() > 0 {
        inputFile = flag.Arg(0)
    }

    if inputFile == "" {
        fmt.Fprintln(os.Stderr, "no input file provided")
        os.Exit(1)
    }
}
