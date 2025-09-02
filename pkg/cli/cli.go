package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/bjia56/objective-lol/pkg/interpreter"
	"github.com/bjia56/objective-lol/pkg/parser"
	"github.com/bjia56/objective-lol/pkg/stdlib"
)

func Run(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: olol <file.olol>")
	}

	filename := args[0]

	// Create interpreter with standard library
	interp := interpreter.NewInterpreter(
		stdlib.DefaultStdlibInitializers(),
		stdlib.DefaultGlobalInitializers()...,
	)

	// Execute the file
	return executeFile(interp, filename)
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

	// Set current file for relative import resolution
	absFilename, _ := filepath.Abs(filename)
	interp.SetCurrentFile(absFilename)

	// Execute the program
	if err := interp.Interpret(program); err != nil {
		return err
	}

	return nil
}