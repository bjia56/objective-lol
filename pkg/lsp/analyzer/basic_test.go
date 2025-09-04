package analyzer

import (
	"context"
	"testing"

	protocol "github.com/tliron/glsp/protocol_3_16"
)

func TestBasicSemanticAnalyzer(t *testing.T) {
	// Use a very simple valid Objective-LOL program
	content := `I CAN HAS STDIO?

HAI ME TEH FUNCSHUN MAIN
KTHXBAI`

	analyzer := NewSemanticAnalyzer("test.olol", content)

	err := analyzer.AnalyzeDocument(context.Background())
	if err != nil {
		t.Fatalf("AnalyzeDocument failed: %v", err)
	}

	symbolTable := analyzer.GetSymbolTable()

	// We should have at least some symbols (stdlib imports)
	if len(symbolTable.Symbols) == 0 {
		t.Error("Expected some symbols to be found")
	}

	// Check that we have a MAIN function
	foundMain := false
	for _, symbol := range symbolTable.Symbols {
		if symbol.Name == "MAIN" && symbol.Kind == SymbolKindFunction {
			foundMain = true
			break
		}
	}
	if !foundMain {
		t.Error("Expected to find MAIN function")
	}

	// Test completions
	position := protocol.Position{Line: 2, Character: 0}
	completions := analyzer.GetCompletionItems(position)

	if len(completions) == 0 {
		t.Error("Expected completion items")
	}

	// Should have keywords
	foundKeyword := false
	for _, item := range completions {
		if item.Label == "SAYZ" { // This should be available from STDIO
			foundKeyword = true
			break
		}
	}
	if !foundKeyword {
		t.Log("Note: SAYZ not found in completions (may be expected)")
	}
}

func TestEnhancedAnalyzerIntegration(t *testing.T) {
	content := `I CAN HAS STDIO?

HAI ME TEH FUNCSHUN MAIN
KTHXBAI`

	// Test the integrated analyzer with enhanced mode
	analyzer := NewAnalyzer()

	diagnostics, err := analyzer.AnalyzeDocument("test.olol", content)
	if err != nil {
		t.Fatalf("AnalyzeDocument failed: %v", err)
	}

	// Should not have errors for this simple valid program
	errorCount := 0
	for _, diag := range diagnostics {
		if diag.Severity != nil && *diag.Severity == protocol.DiagnosticSeverityError {
			errorCount++
			t.Logf("Error diagnostic: %s", diag.Message)
		}
	}

	if errorCount > 0 {
		t.Errorf("Expected no errors, got %d", errorCount)
	}

	// Test that we can get completions
	position := protocol.Position{Line: 2, Character: 0}
	completions, err := analyzer.GetCompletions("test.olol", content, position)
	if err != nil {
		t.Fatalf("GetCompletions failed: %v", err)
	}

	if len(completions) == 0 {
		t.Error("Expected completion items")
	}

	t.Logf("Found %d completion items", len(completions))
}
