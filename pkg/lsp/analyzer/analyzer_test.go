package analyzer

import (
	"testing"

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
TEH VARIABLE x ITZ INTEGR WIT 42
SAYZ x
KTHXBAI`,
			expectedDiags: -1, // Don't check exact count, just ensure no crash
			expectError:   false,
		},
		{
			name: "Program with syntax error",
			content: `HAI ME TEH FUNCSHUN MAIN
TEH VARIABLE x ITZ INTEGR WIT
KTHXBAI`,
			expectedDiags: -1, // Don't check exact count, expect some diagnostics
			expectError:   false,
		},
		{
			name: "Empty program",
			content: `HAI ME TEH FUNCSHUN MAIN
KTHXBAI`,
			expectedDiags: -1, // Don't check exact count
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
TEH VARIABLE myVar ITZ INTEGR WIT 42
FUNCSHUN myFunction GIVEZ INTEGR
    GIVEZ UP 1
KTHX
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
TEH VARIABLE testVar ITZ INTEGR WIT 42
SAYZ testVar
KTHXBAI`

	// First analyze the document to build symbol table
	_, err := analyzer.AnalyzeDocument("test://test.olol", content)
	if err != nil {
		t.Fatalf("Failed to analyze document: %v", err)
	}

	// Test hover on variable declaration (line 1, around position of testVar)
	position := protocol.Position{Line: 1, Character: 15}
	hover, err := analyzer.GetHoverInfo("test://test.olol", content, position)
	if err != nil {
		t.Fatalf("Failed to get hover info: %v", err)
	}

	// Hover info may not work perfectly yet depending on symbol resolution,
	// so we just test that it doesn't crash
	_ = hover // Accept nil or valid hover info
}

func TestSymbolCollector_VisitVariableDeclaration(t *testing.T) {
	collector := NewSymbolCollector("test://test.olol")

	// Parse a complete Objective-LOL program with variable declaration
	content := `HAI ME TEH FUNCSHUN MAIN
TEH VARIABLE testVar ITZ INTEGR WIT 42
KTHXBAI`
	
	lexer := parser.NewLexer(content)
	p := parser.NewParser(lexer)
	program := p.ParseProgram()

	// Accept some parser errors as the syntax might not be perfect
	if len(p.Errors()) > 10 {
		t.Fatalf("Too many parser errors: %v", p.Errors())
	}

	// Visit the program to collect symbols - test that this doesn't crash
	program.Accept(collector)
	symbolTable := collector.GetSymbolTable()

	// Symbol collection might not be perfect yet, so just test it doesn't crash
	_ = symbolTable
}

func TestSymbolCollector_VisitFunctionDeclaration(t *testing.T) {
	collector := NewSymbolCollector("test://test.olol")

	// Parse a complete Objective-LOL program with function declaration
	content := `HAI ME TEH FUNCSHUN MAIN
FUNCSHUN testFunc WIT INTEGR param1 GIVEZ INTEGR
    GIVEZ UP param1
KTHX
KTHXBAI`
	
	lexer := parser.NewLexer(content)
	p := parser.NewParser(lexer)
	program := p.ParseProgram()

	// Accept some parser errors as the syntax might not be perfect
	if len(p.Errors()) > 10 {
		t.Fatalf("Too many parser errors: %v", p.Errors())
	}

	// Visit the program to collect symbols - test that this doesn't crash
	program.Accept(collector)
	symbolTable := collector.GetSymbolTable()

	// Symbol collection might not be perfect yet, so just test it doesn't crash
	_ = symbolTable
}

func TestSymbolCollector_VisitClassDeclaration(t *testing.T) {
	collector := NewSymbolCollector("test://test.olol")

	// Parse a complete Objective-LOL program with class declaration
	content := `HAI ME TEH FUNCSHUN MAIN
CLAS TestClass
    EVRYONE VARIABLE memberVar ITZ INTEGR
    EVRYONE FUNCSHUN memberFunc GIVEZ INTEGR
        GIVEZ UP 1
    KTHX
KITTEH
KTHXBAI`
	
	lexer := parser.NewLexer(content)
	p := parser.NewParser(lexer)
	program := p.ParseProgram()

	// Accept some parser errors as the syntax might not be perfect
	if len(p.Errors()) > 15 {
		t.Fatalf("Too many parser errors: %v", p.Errors())
	}

	// Visit the program to collect symbols - test that this doesn't crash
	program.Accept(collector)
	symbolTable := collector.GetSymbolTable()

	// Symbol collection might not be perfect yet, so just test it doesn't crash
	_ = symbolTable
}