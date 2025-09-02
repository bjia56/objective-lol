package main

import (
	"fmt"
	"os"

	"github.com/bjia56/objective-lol/pkg/cli"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <file.olol>\n", os.Args[0])
		os.Exit(1)
	}

	os.Exit(cli.Run(os.Args[1:]))
}
