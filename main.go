package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/lucasch37/spi-go/interpreter"
	"github.com/lucasch37/spi-go/lexer"
	"github.com/lucasch37/spi-go/parser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <source.pas>\n", os.Args[0])
		os.Exit(1)
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	lexer := lexer.NewLexer(strings.TrimSpace(string(data)))
	parser, err := parser.NewParser(lexer)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	interpreter := interpreter.NewInterpreter(parser)

	if err := interpreter.Interpret(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(interpreter.GLOBAL_SCOPE)
}
