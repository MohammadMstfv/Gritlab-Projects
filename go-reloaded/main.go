package main

import (
	"fmt"
	"os"
)

func main() {
	// Check that the correct number of command-line arguments are provided. and if not say how it is supposed to work
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "Usage: go run . <input file> <output file>")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]

	data, err := os.ReadFile(inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", inputPath, err)
		os.Exit(1)
	}

	result := Transform(string(data))

	if err := os.WriteFile(outputPath, []byte(result), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", outputPath, err)
		os.Exit(1)
	}
}
