package main

import (
	"os"

	"example.com/aryaLee/golang/subcommand/pkg/subcommands"
)

func main() {
	if len(os.Args) < 2 {
		subcommands.Usage()
	}

	cmd := os.Args[1]
	subcommands.Run(cmd)
}
