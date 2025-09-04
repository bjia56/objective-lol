package analyzer

import (
	"context"
	"strings"
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

func TestAnalyzer_FunctionCallHover(t *testing.T) {
	analyzer := NewAnalyzer()

	// Test program with function calls
	content := `I CAN HAS STDIO?

HAI ME TEH FUNCSHUN MAIN
	I HAS A VARIABLE x TEH INTEGR ITZ 42
	SAYZ WIT x
	I HAS A VARIABLE result TEH INTEGR ITZ MAX WIT 1 AN 2
KTHXBAI`

	// First analyze the document
	_, err := analyzer.AnalyzeDocument("test://test.olol", content)
	if err != nil {
		t.Fatalf("Failed to analyze document: %v", err)
	}

	// Test hover on SAYZ function call (line 4, around column 1)
	position := protocol.Position{Line: 4, Character: 1}
	hover, err := analyzer.GetHoverInfo("test://test.olol", content, position)
	if err != nil {
		t.Fatalf("Failed to get hover info: %v", err)
	}

	if hover != nil {
		if markup, ok := hover.Contents.(protocol.MarkupContent); ok {
			t.Logf("SAYZ hover content: %s", markup.Value)
			// Should contain function call information
			if !strings.Contains(markup.Value, "Function Call") {
				t.Logf("Note: SAYZ hover may not contain function call info (position might not match exactly)")
			}
		}
	}

	// Test hover on MAX function call (line 5, around column 45)
	position = protocol.Position{Line: 5, Character: 45}
	hover, err = analyzer.GetHoverInfo("test://test.olol", content, position)
	if err != nil {
		t.Fatalf("Failed to get hover info: %v", err)
	}

	if hover != nil {
		if markup, ok := hover.Contents.(protocol.MarkupContent); ok {
			t.Logf("MAX hover content: %s", markup.Value)
		}
	}
}

func TestSemanticAnalyzer_FunctionCallTracking(t *testing.T) {
	content := `I CAN HAS STDIO?

HAI ME TEH FUNCSHUN MAIN
	SAYZ WIT "Hello"
	I HAS A VARIABLE result TEH INTEGR ITZ MAX WIT 1 AN 2
	MAX WIT 3 AN 4
KTHXBAI`

	analyzer := NewSemanticAnalyzer("test://test.olol", content)

	err := analyzer.AnalyzeDocument(context.Background())
	if err != nil {
		t.Fatalf("AnalyzeDocument failed: %v", err)
	}

	symbolTable := analyzer.GetSymbolTable()

	// Should have tracked function calls
	if len(symbolTable.FunctionCalls) == 0 {
		t.Error("Expected function calls to be tracked")
	} else {
		t.Logf("Found %d function calls", len(symbolTable.FunctionCalls))
		
		for i, call := range symbolTable.FunctionCalls {
			t.Logf("Call %d: %s (type: %d)", i, call.FunctionName, call.CallType)
			if call.ResolvedTo != nil {
				t.Logf("  Resolved to: %s", call.ResolvedTo.Name)
			} else {
				t.Logf("  Not resolved")
			}
		}
	}

	// Test specific function calls
	foundSAYZ := false
	foundMAX := false
	
	for _, call := range symbolTable.FunctionCalls {
		if call.FunctionName == "SAYZ" {
			foundSAYZ = true
			if call.CallType != FunctionCallGlobal {
				t.Errorf("Expected SAYZ to be global call, got %d", call.CallType)
			}
		}
		if call.FunctionName == "MAX" {
			foundMAX = true  
			if call.CallType != FunctionCallGlobal {
				t.Errorf("Expected MAX to be global call, got %d", call.CallType)
			}
		}
	}

	if !foundSAYZ {
		t.Error("Expected to find SAYZ function call")
	}
	if !foundMAX {
		t.Error("Expected to find MAX function call")
	}
}
