package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/lucasch37/spi-go/interpreter"
	"github.com/lucasch37/spi-go/lexer"
	"github.com/lucasch37/spi-go/parser"
)

func main() {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	lexer := lexer.NewLexer(strings.TrimSpace(string(data)))
	parser := parser.NewParser(lexer)
	interpreter := interpreter.NewInterpreter(parser)

	interpreter.Interpret()
	fmt.Println(interpreter.GLOBALSCOPE)
}
