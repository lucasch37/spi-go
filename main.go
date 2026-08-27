package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	lexer := NewLexer(strings.TrimSpace(string(data)))
	parser := NewParser(lexer)
	interpreter := NewInterpreter(parser)

	interpreter.Interpret()
	fmt.Println(interpreter.GLOBALSCOPE)
}
