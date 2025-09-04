package parser

import (
	"testing"
	"reflect"
	"github.com/bjia56/objective-lol/pkg/ast"
)

func TestMultipleInheritanceParser(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name: "Single inheritance (backwards compatibility)",
			input: `HAI ME TEH CLAS CHILD KITTEH OF PARENT
			       KTHXBAI`,
			expected: []string{"PARENT"},
		},
		{
			name: "Two parent classes",
			input: `HAI ME TEH CLAS CHILD KITTEH OF PARENT1 AN OF PARENT2
			       KTHXBAI`,
			expected: []string{"PARENT1", "PARENT2"},
		},
		{
			name: "Three parent classes",
			input: `HAI ME TEH CLAS COMPLEX KITTEH OF A AN OF B AN OF C
			       KTHXBAI`,
			expected: []string{"A", "B", "C"},
		},
		{
			name: "Four parent classes",
			input: `HAI ME TEH CLAS MULTI KITTEH OF A AN OF B AN OF C AN OF D
			       KTHXBAI`,
			expected: []string{"A", "B", "C", "D"},
		},
		{
			name: "No inheritance",
			input: `HAI ME TEH CLAS STANDALONE
			       KTHXBAI`,
			expected: []string{},
		},
		{
			name: "Mixed case parent names",
			input: `HAI ME TEH CLAS CHILD KITTEH OF parent1 AN OF Parent2 AN OF PARENT3
			       KTHXBAI`,
			expected: []string{"PARENT1", "PARENT2", "PARENT3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			parser := NewParser(lexer)
			program := parser.ParseProgram()
			
			// Check for parser errors
			if len(parser.Errors()) > 0 {
				t.Fatalf("Parser errors: %v", parser.Errors())
			}
			
			if len(program.Declarations) != 1 {
				t.Fatalf("Expected 1 declaration, got %d", len(program.Declarations))
			}
			
			classDecl, ok := program.Declarations[0].(*ast.ClassDeclarationNode)
			if !ok {
				t.Fatalf("Expected ClassDeclarationNode, got %T", program.Declarations[0])
			}
			
			if !reflect.DeepEqual(classDecl.ParentClasses, tt.expected) {
				t.Errorf("Expected parent classes %v, got %v", tt.expected, classDecl.ParentClasses)
			}
		})
	}
}

func TestMultipleInheritanceParserErrors(t *testing.T) {
	errorTests := []struct {
		name  string
		input string
		expectedErrorCount int
	}{
		{
			name: "Missing OF after AN",
			input: `HAI ME TEH CLAS CHILD KITTEH OF PARENT1 AN PARENT2
			       KTHXBAI`,
			expectedErrorCount: 1,
		},
		{
			name: "Missing parent name after AN OF",
			input: `HAI ME TEH CLAS CHILD KITTEH OF PARENT1 AN OF
			       KTHXBAI`,
			expectedErrorCount: 1,
		},
		{
			name: "Missing OF after KITTEH",
			input: `HAI ME TEH CLAS CHILD KITTEH PARENT1
			       KTHXBAI`,
			expectedErrorCount: 1,
		},
	}

	for _, tt := range errorTests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			parser := NewParser(lexer)
			parser.ParseProgram()
			
			errors := parser.Errors()
			if len(errors) != tt.expectedErrorCount {
				t.Errorf("Expected %d errors, got %d: %v", tt.expectedErrorCount, len(errors), errors)
			}
		})
	}
}

func TestMultipleInheritanceWithMembers(t *testing.T) {
	input := `BTW Test class with multiple inheritance and members
	HAI ME TEH CLAS CHILD KITTEH OF PARENT1 AN OF PARENT2
		EVRYONE
		DIS TEH VARIABLE NAME TEH STRIN ITZ "Test"
		
		DIS TEH FUNCSHUN GET_NAME TEH STRIN
			GIVEZ NAME
		KTHX
	KTHXBAI`

	lexer := NewLexer(input)
	parser := NewParser(lexer)
	program := parser.ParseProgram()
	
	// Check for parser errors
	if len(parser.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", parser.Errors())
	}
	
	if len(program.Declarations) != 1 {
		t.Fatalf("Expected 1 declaration, got %d", len(program.Declarations))
	}
	
	classDecl, ok := program.Declarations[0].(*ast.ClassDeclarationNode)
	if !ok {
		t.Fatalf("Expected ClassDeclarationNode, got %T", program.Declarations[0])
	}
	
	expectedParents := []string{"PARENT1", "PARENT2"}
	if !reflect.DeepEqual(classDecl.ParentClasses, expectedParents) {
		t.Errorf("Expected parent classes %v, got %v", expectedParents, classDecl.ParentClasses)
	}
	
	// Check that class members are parsed correctly
	if len(classDecl.Members) != 2 {
		t.Errorf("Expected 2 class members, got %d", len(classDecl.Members))
	}
}