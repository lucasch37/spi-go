package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/lucasch37/spi-go/interpreter"
	"github.com/lucasch37/spi-go/lexer"
	"github.com/lucasch37/spi-go/parser"
	"github.com/lucasch37/spi-go/semantic"
)

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func main() {
	config := parseArgs()

	data, err := os.ReadFile(config.InputFile)
	if err != nil {
		fatal(err)
	}

	lexer := lexer.NewLexer(strings.TrimSpace(string(data)))
	parser, err := parser.NewParser(lexer)
	if err != nil {
		fatal(err)
	}

	tree, err := parser.Parse()
	if err != nil {
		fatal(err)
	}

	symTabBuilder := semantic.NewSemanticAnalyzer(config.ShouldLogScope)
	if err := symTabBuilder.Visit(tree); err != nil {
		fatal(err)
	}

	interpreter := interpreter.NewInterpreter(tree)
	if err := interpreter.Interpret(); err != nil {
		fatal(err)
	}

	for key, value := range interpreter.GLOBAL_SCOPE {
		fmt.Printf("%s = %s\n", key, value.String())
	}
}
