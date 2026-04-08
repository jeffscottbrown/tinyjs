package main

import (
	"fmt"
	"io"
	"os"

	"example.com/tinyjs/internal/compiler"
	"example.com/tinyjs/internal/parser"
)

func main() {
	src, err := readSource()
	if err != nil {
		exitf("read source: %v", err)
	}

	prog, err := parser.Parse(src)
	if err != nil {
		exitf("parse error: %v", err)
	}

	out, err := compiler.Compile(prog)
	if err != nil {
		exitf("compile error: %v", err)
	}

	fmt.Print(out)
}

func readSource() (string, error) {
	switch len(os.Args) {
	case 1:
		b, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", err
		}
		return string(b), nil
	case 2:
		b, err := os.ReadFile(os.Args[1])
		if err != nil {
			return "", err
		}
		return string(b), nil
	default:
		return "", fmt.Errorf("usage: tinyjs [file]")
	}
}

func exitf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
