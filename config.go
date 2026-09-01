package main

import (
	"flag"
	"os"
)

type Config struct {
	InputFile      string
	ShouldLogScope bool
}

func parseArgs() Config {
	scope := flag.Bool(
		"scope",
		false,
		"Print scope information",
	)

	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	return Config{
		InputFile:      flag.Arg(0),
		ShouldLogScope: *scope,
	}
}
