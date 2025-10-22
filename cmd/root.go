package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "markee",
	Short: "A CommonMark compliant markdown processor",
	Long:  "A markdown parser and renderer built from scratch in Go, supporting CommonMark specification.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(parseCmd)
}
