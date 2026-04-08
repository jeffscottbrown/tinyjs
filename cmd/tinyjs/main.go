package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"example.com/tinyjs/internal/compiler"
	"example.com/tinyjs/internal/parser"
	"example.com/tinyjs/internal/runner"
)

func main() {
	showSource := flag.Bool("s", false, "Show LLVM IR instead of running")
	outputFile := flag.String("o", "", "Specify output file")

	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "Usage: tinyjs [-s] <file.tjs>")
		os.Exit(1)
	}

	arg := flag.Arg(0)

	content, err := os.ReadFile(arg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading file: %v\n", err)
		os.Exit(1)
	}

	parser, err := parser.New()
	program, err := parser.ParseString("demo.tjs", string(content))
	if err != nil {
		log.Fatal(err)
	}
	compiler := compiler.MustNew()
	irText, err := compiler.GenerateIR(program)
	if err != nil {
		log.Fatal(err)
	}

	if *showSource {
		fmt.Println(irText)
	} else if *outputFile != "" {
		if err := compiler.Build(irText, *outputFile); err != nil {
			log.Fatal(err)
		}
	} else {
		out, err := runner.RunIR(irText)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(out)
	}
}
