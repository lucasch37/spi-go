package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/lucasch37/spi-go/internal/interpreter"
	"github.com/lucasch37/spi-go/internal/lexer"
	"github.com/lucasch37/spi-go/internal/parser"
	"github.com/lucasch37/spi-go/internal/semantic"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	config := parseArgs()

	data, err := os.ReadFile(config.InputFile)
	if err != nil {
		return err
	}

	lexer := lexer.NewLexer(strings.TrimSpace(string(data)))

	parser, err := parser.NewParser(lexer)
	if err != nil {
		return err
	}

	tree, err := parser.Parse()
	if err != nil {
		return err
	}

	analyzer := semantic.NewSemanticAnalyzer(config.ShouldLogScope)
	if err := analyzer.Analyze(tree); err != nil {
		return err
	}

	interpreter := interpreter.NewInterpreter(tree, config.ShouldLogStack)

	return interpreter.Interpret()
}
