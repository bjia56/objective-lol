package parser

import (
	"testing"
	
	"github.com/bjia56/objective-lol/pkg/ast"
)

func TestFunctionDocumentationParsing(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		expectedFuncName string
		expectedDocs     []string
	}{
		{
			name: "Single line documentation",
			input: `BTW This function adds two numbers
HAI ME TEH FUNCSHUN ADD TEH INTEGR WIT X TEH INTEGR AN WIT Y TEH INTEGR
    GIVEZ X MOAR Y
KTHXBAI`,
			expectedFuncName: "ADD",
			expectedDocs:     []string{"This function adds two numbers"},
		},
		{
			name: "Multi-line documentation",
			input: `BTW This function calculates factorial
BTW @param n The number to calculate factorial for
BTW @return The factorial result
HAI ME TEH FUNCSHUN FACTORIAL TEH INTEGR WIT N TEH INTEGR
    IZ N SMALLR THAN 2?
        GIVEZ 1
    NOPE
        GIVEZ N TIEMZ FACTORIAL WIT N LES 1
    KTHX
KTHXBAI`,
			expectedFuncName: "FACTORIAL",
			expectedDocs: []string{
				"This function calculates factorial",
				"@param n The number to calculate factorial for",
				"@return The factorial result",
			},
		},
		{
			name: "No documentation",
			input: `HAI ME TEH FUNCSHUN SIMPLE
    SAYZ WIT "Hello"
KTHXBAI`,
			expectedFuncName: "SIMPLE",
			expectedDocs:     nil,
		},
		{
			name: "Documentation with empty lines",
			input: `BTW First line of docs
BTW 
BTW Third line after empty comment
HAI ME TEH FUNCSHUN WITH_EMPTY TEH STRIN
    GIVEZ "test"
KTHXBAI`,
			expectedFuncName: "WITH_EMPTY",
			expectedDocs: []string{
				"First line of docs",
				"",
				"Third line after empty comment",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			parser := NewParser(lexer)
			program := parser.ParseProgram()

			if len(parser.Errors()) > 0 {
				t.Fatalf("Parser errors: %v", parser.Errors())
			}

			if len(program.Declarations) != 1 {
				t.Fatalf("Expected 1 declaration, got %d", len(program.Declarations))
			}

			funcDecl, ok := program.Declarations[0].(*ast.FunctionDeclarationNode)
			if !ok {
				t.Fatalf("Expected FunctionDeclarationNode, got %T", program.Declarations[0])
			}

			if funcDecl.Name != tt.expectedFuncName {
				t.Errorf("Expected function name %s, got %s", tt.expectedFuncName, funcDecl.Name)
			}

			if len(funcDecl.Documentation) != len(tt.expectedDocs) {
				t.Errorf("Expected %d documentation lines, got %d", len(tt.expectedDocs), len(funcDecl.Documentation))
				t.Errorf("Got documentation: %v", funcDecl.Documentation)
			}

			for i, expectedDoc := range tt.expectedDocs {
				if i >= len(funcDecl.Documentation) {
					t.Errorf("Missing documentation line %d: expected %q", i, expectedDoc)
					continue
				}
				if funcDecl.Documentation[i] != expectedDoc {
					t.Errorf("Documentation line %d: expected %q, got %q", i, expectedDoc, funcDecl.Documentation[i])
				}
			}
		})
	}
}

func TestMultipleFunctionDocumentation(t *testing.T) {
	input := `BTW First function documentation
HAI ME TEH FUNCSHUN FIRST
    SAYZ WIT "First"
KTHXBAI

BTW Second function documentation
BTW This one has multiple lines
HAI ME TEH FUNCSHUN SECOND TEH STRIN
    GIVEZ "Second"
KTHXBAI

HAI ME TEH FUNCSHUN THIRD
    SAYZ WIT "No docs"
KTHXBAI`

	lexer := NewLexer(input)
	parser := NewParser(lexer)
	program := parser.ParseProgram()

	if len(parser.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", parser.Errors())
	}

	if len(program.Declarations) != 3 {
		t.Fatalf("Expected 3 declarations, got %d", len(program.Declarations))
	}

	// Check first function
	func1, ok := program.Declarations[0].(*ast.FunctionDeclarationNode)
	if !ok {
		t.Fatalf("Expected FunctionDeclarationNode, got %T", program.Declarations[0])
	}
	if func1.Name != "FIRST" {
		t.Errorf("Expected function name FIRST, got %s", func1.Name)
	}
	if len(func1.Documentation) != 1 || func1.Documentation[0] != "First function documentation" {
		t.Errorf("Expected documentation [\"First function documentation\"], got %v", func1.Documentation)
	}

	// Check second function
	func2, ok := program.Declarations[1].(*ast.FunctionDeclarationNode)
	if !ok {
		t.Fatalf("Expected FunctionDeclarationNode, got %T", program.Declarations[1])
	}
	if func2.Name != "SECOND" {
		t.Errorf("Expected function name SECOND, got %s", func2.Name)
	}
	expectedDocs := []string{"Second function documentation", "This one has multiple lines"}
	if len(func2.Documentation) != len(expectedDocs) {
		t.Errorf("Expected %d documentation lines, got %d", len(expectedDocs), len(func2.Documentation))
	}
	for i, expected := range expectedDocs {
		if i >= len(func2.Documentation) || func2.Documentation[i] != expected {
			t.Errorf("Documentation line %d: expected %q, got %q", i, expected, func2.Documentation[i])
		}
	}

	// Check third function (no documentation)
	func3, ok := program.Declarations[2].(*ast.FunctionDeclarationNode)
	if !ok {
		t.Fatalf("Expected FunctionDeclarationNode, got %T", program.Declarations[2])
	}
	if func3.Name != "THIRD" {
		t.Errorf("Expected function name THIRD, got %s", func3.Name)
	}
	if len(func3.Documentation) != 0 {
		t.Errorf("Expected no documentation, got %v", func3.Documentation)
	}
}