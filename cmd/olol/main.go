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

	if err := cli.Run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
