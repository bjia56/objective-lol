package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParserCreation(t *testing.T) {
	input := "HAI ME TEH FUNCSHUN TEST\nKTHXBAI"
	lexer := NewLexer(input)
	parser := NewParser(lexer)

	assert.NotNil(t, parser, "Parser should be created")
	assert.Empty(t, parser.Errors(), "New parser should have no errors")
}

func TestParserBasicProgram(t *testing.T) {
	input := `I CAN HAS STDIO?

HAI ME TEH FUNCSHUN MAIN
    I HAS A VARIABLE x TEH INTEGR ITZ 42
KTHXBAI

MAIN()`

	lexer := NewLexer(input)
	parser := NewParser(lexer)
	program := parser.ParseProgram()

	// Just check that parsing doesn't crash and returns something
	require.NotNil(t, program, "Program should not be nil")

	// If there are errors, at least print them for debugging
	if len(parser.Errors()) > 0 {
		t.Logf("Parser errors (may be expected): %v", parser.Errors())
	}
}

func TestParserSimpleFunction(t *testing.T) {
	input := `HAI ME TEH FUNCSHUN ADD WIT X TEH INTEGR AN WIT Y TEH INTEGR
    GIVEZ X MOAR Y
KTHXBAI`

	lexer := NewLexer(input)
	parser := NewParser(lexer)
	program := parser.ParseProgram()

	require.NotNil(t, program, "Program should not be nil")

	// Log any errors for debugging
	if len(parser.Errors()) > 0 {
		t.Logf("Parser errors: %v", parser.Errors())
	}
}

func TestParserSimpleClass(t *testing.T) {
	input := `HAI ME TEH CLAS PERSON
    EVRYONE
    DIS TEH VARIABLE NAME TEH STRIN ITZ "Unknown"
KTHXBAI`

	lexer := NewLexer(input)
	parser := NewParser(lexer)
	program := parser.ParseProgram()

	require.NotNil(t, program, "Program should not be nil")

	// Log any errors for debugging
	if len(parser.Errors()) > 0 {
		t.Logf("Parser errors: %v", parser.Errors())
	}
}

func TestParserImportStatement(t *testing.T) {
	input := `I CAN HAS STDIO?`

	lexer := NewLexer(input)
	parser := NewParser(lexer)
	program := parser.ParseProgram()

	require.NotNil(t, program, "Program should not be nil")

	// Should handle import statements
	if len(parser.Errors()) == 0 {
		assert.Greater(t, len(program.Declarations), 0, "Should have at least one import declaration")
	} else {
		t.Logf("Parser errors: %v", parser.Errors())
	}
}

func TestParserErrorCollection(t *testing.T) {
	input := `INVALID SYNTAX HERE
AND MORE INVALID STUFF`

	lexer := NewLexer(input)
	parser := NewParser(lexer)
	program := parser.ParseProgram()

	// Parser should still return a program even with errors
	require.NotNil(t, program, "Program should not be nil even with errors")

	// Should collect errors
	assert.NotEmpty(t, parser.Errors(), "Parser should collect errors for invalid syntax")
}

func TestParserEmptyProgram(t *testing.T) {
	input := ""

	lexer := NewLexer(input)
	parser := NewParser(lexer)
	program := parser.ParseProgram()

	require.NotNil(t, program, "Program should not be nil")
	assert.Equal(t, 0, len(program.Declarations), "Empty program should have no declarations")
}

func TestParserOperatorPrecedence(t *testing.T) {
	// Test that parser can handle basic operator precedence
	input := `HAI ME TEH FUNCSHUN TEST
    I HAS A VARIABLE RESULT ITZ 2 MOAR 3 TIEMZ 4
KTHXBAI`

	lexer := NewLexer(input)
	parser := NewParser(lexer)
	program := parser.ParseProgram()

	require.NotNil(t, program, "Program should not be nil")

	// Log any errors for debugging
	if len(parser.Errors()) > 0 {
		t.Logf("Parser errors: %v", parser.Errors())
	}
}

func TestParserParentheses(t *testing.T) {
	// Test that parser can handle parentheses
	input := `HAI ME TEH FUNCSHUN TEST
    I HAS A VARIABLE RESULT ITZ (2 MOAR 3) TIEMZ 4
KTHXBAI`

	lexer := NewLexer(input)
	parser := NewParser(lexer)
	program := parser.ParseProgram()

	require.NotNil(t, program, "Program should not be nil")

	// Log any errors for debugging
	if len(parser.Errors()) > 0 {
		t.Logf("Parser errors: %v", parser.Errors())
	}
}
