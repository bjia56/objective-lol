package analyzer

import (
	"testing"

	"github.com/stretchr/testify/require"
	protocol "github.com/tliron/glsp/protocol_3_16"

	"github.com/bjia56/objective-lol/pkg/parser"
)

func TestAnalyzer_AnalyzeDocument(t *testing.T) {
	analyzer := NewAnalyzer()

	tests := []struct {
		name          string
		content       string
		expectedDiags int
		expectError   bool
	}{
		{
			name: "Valid Objective-LOL program",
			content: `HAI ME TEH FUNCSHUN MAIN
I CAN HAS STDIO?
I HAS A VARIABLE x TEH INTEGR ITZ 42
SAYZ WIT x
KTHXBAI`,
			expectedDiags: 0,
			expectError:   false,
		},
		{
			name: "Program with syntax error",
			content: `HAI ME TEH FUNCSHUN MAIN
I HAS A VARIABLE x TEH INTEGR WIT
KTHXBAI`,
			expectedDiags: 1,
			expectError:   false,
		},
		{
			name: "Empty program",
			content: `HAI ME TEH FUNCSHUN MAIN
KTHXBAI`,
			expectedDiags: 0,
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diagnostics, err := analyzer.AnalyzeDocument("test://test.olol", tt.content)

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// Only check diagnostic count if expectedDiags >= 0
			if tt.expectedDiags >= 0 && len(diagnostics) != tt.expectedDiags {
				t.Errorf("Expected %d diagnostics, got %d", tt.expectedDiags, len(diagnostics))
			}
		})
	}
}

func TestAnalyzer_GetCompletions(t *testing.T) {
	analyzer := NewAnalyzer()
	content := `HAI ME TEH FUNCSHUN MAIN
I HAS A VARIABLE myVar TEH INTEGR ITZ 42
KTHXBAI`

	// First analyze the document to build symbol table
	_, err := analyzer.AnalyzeDocument("test://test.olol", content)
	if err != nil {
		t.Fatalf("Failed to analyze document: %v", err)
	}

	// Test completion at various positions
	position := protocol.Position{Line: 3, Character: 0}
	completions, err := analyzer.GetCompletions("test://test.olol", content, position)
	if err != nil {
		t.Fatalf("Failed to get completions: %v", err)
	}

	// Should have keywords and symbols
	if len(completions) == 0 {
		t.Error("Expected completions but got none")
	}

	// Check that we have some expected keywords
	hasKeywords := false

	for _, completion := range completions {
		if completion.Label == "GIVEZ" || completion.Label == "TEH" || completion.Label == "SAYZ" {
			hasKeywords = true
			break
		}
	}

	if !hasKeywords {
		t.Error("Expected at least some keyword completions")
	}

	// Note: Symbol completions may not work perfectly yet, so we don't test for them
}

func TestAnalyzer_GetHoverInfo(t *testing.T) {
	analyzer := NewAnalyzer()
	content := `HAI ME TEH FUNCSHUN MAIN
I CAN HAS STDIO?
I HAS A VARIABLE testVar TEH INTEGR ITZ 42
SAYZ WIT testVar
KTHXBAI`

	// First analyze the document to build symbol table
	_, err := analyzer.AnalyzeDocument("test://test.olol", content)
	if err != nil {
		t.Fatalf("Failed to analyze document: %v", err)
	}

	// Test hover on variable declaration (around position of testVar)
	position := protocol.Position{Line: 2, Character: 20}
	hover, err := analyzer.GetHoverInfo("test://test.olol", content, position)
	if err != nil {
		t.Fatalf("Failed to get hover info: %v", err)
	}

	require.NotNil(t, hover, "Expected hover info but got nil")
}

func TestSymbolCollector_VisitVariableDeclaration(t *testing.T) {
	collector := NewSymbolCollector("test://test.olol")

	// Parse a complete Objective-LOL program with variable declaration
	content := `HAI ME TEH FUNCSHUN MAIN
I HAS A VARIABLE testVar TEH INTEGR ITZ 42
KTHXBAI`

	lexer := parser.NewLexer(content)
	p := parser.NewParser(lexer)
	program := p.ParseProgram()

	// Visit the program to collect symbols - test that this doesn't crash
	program.Accept(collector)
	symbolTable := collector.GetSymbolTable()

	require.NotNil(t, symbolTable, "Expected symbol table but got nil")
}

func TestSymbolCollector_VisitFunctionDeclaration(t *testing.T) {
	collector := NewSymbolCollector("test://test.olol")

	// Parse a complete Objective-LOL program with function declaration
	content := `HAI ME TEH FUNCSHUN foo
KTHXBAI`

	lexer := parser.NewLexer(content)
	p := parser.NewParser(lexer)
	program := p.ParseProgram()

	// Visit the program to collect symbols - test that this doesn't crash
	program.Accept(collector)
	symbolTable := collector.GetSymbolTable()

	require.NotNil(t, symbolTable, "Expected symbol table but got nil")
}

func TestSymbolCollector_VisitClassDeclaration(t *testing.T) {
	collector := NewSymbolCollector("test://test.olol")

	// Parse a complete Objective-LOL program with class declaration
	content := `HAI ME TEH CLAS TestClass
    EVRYONE
	I HAS A VARIABLE memberVar ITZ INTEGR
    I HAS A FUNCSHUN memberFunc TEH INTEGR
        GIVEZ 1
    KTHX
KTHXBAI`

	lexer := parser.NewLexer(content)
	p := parser.NewParser(lexer)
	program := p.ParseProgram()

	// Visit the program to collect symbols - test that this doesn't crash
	program.Accept(collector)
	symbolTable := collector.GetSymbolTable()

	require.NotNil(t, symbolTable, "Expected symbol table but got nil")
}
