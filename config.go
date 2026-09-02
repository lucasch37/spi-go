package main

import (
	"flag"
	"os"
)

type Config struct {
	InputFile      string
	ShouldLogScope bool
	ShouldLogStack bool
}

func parseArgs() Config {
	scope := flag.Bool(
		"scope",
		false,
		"Print scope information",
	)

	stack := flag.Bool(
		"stack",
		false,
		"Print call stack information",
	)

	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	return Config{
		InputFile:      flag.Arg(0),
		ShouldLogScope: *scope,
		ShouldLogStack: *stack,
	}
}
