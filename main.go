package main

import (
	"fmt"
	"os"

	"github.com/bjia56/objective-lol/cmd/olol"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <file.olol>\n", os.Args[0])
		os.Exit(1)
	}

	os.Exit(olol.Run(os.Args[1:]))
}
