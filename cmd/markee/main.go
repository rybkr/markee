package main

import (
    "flag"
    "log"
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

    // Early exit if no input file is given
    if inputFile == "" {
        log.Fatal("no input file provided")
    }

    // Check if the file exists
    if _, err := os.Stat(inputFile); err != nil {
        if os.IsNotExist(err) {
            log.Fatalf("file %q does not exist", inputFile)
        } else {
            log.Fatalf("error checking file %q: %v", inputFile, err)
        }
    }
}
