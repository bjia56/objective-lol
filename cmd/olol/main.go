package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/bjia56/objective-lol/pkg/interpreter"
	"github.com/bjia56/objective-lol/pkg/parser"
	"github.com/bjia56/objective-lol/pkg/stdlib"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <file.olol>\n", os.Args[0])
		os.Exit(1)
	}

	filename := os.Args[1]

	// Create interpreter
	interp := interpreter.NewInterpreter()

	// Register standard library functions
	stdlib.RegisterSTDIO(interp.GetRuntime())
	stdlib.RegisterMATH(interp.GetRuntime())
	stdlib.RegisterTIEM(interp.GetRuntime())

	// Execute the file
	if err := executeFile(interp, filename); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func executeFile(interp *interpreter.Interpreter, filename string) error {
	// Check if file exists and has .olol extension
	if !strings.HasSuffix(strings.ToLower(filename), ".olol") {
		return fmt.Errorf("file must have .olol extension")
	}

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", filename)
	}

	// Read the source file
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	// Parse the source code
	lexer := parser.NewLexer(string(content))
	p := parser.NewParser(lexer)
	program := p.ParseProgram()

	// Check for parsing errors
	if errors := p.Errors(); len(errors) > 0 {
		for _, err := range errors {
			fmt.Fprintf(os.Stderr, "Parse error: %s\n", err)
		}
		return fmt.Errorf("parsing failed with %d errors", len(errors))
	}

	// Execute the program
	if err := interp.Interpret(program); err != nil {
		return err
	}

	return nil
}
