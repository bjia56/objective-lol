package analyzer

import (
	"context"
	"testing"
)

// Parse errors must be reported where the problem is, not at 1:1.
func TestParseErrorDiagnosticPositions(t *testing.T) {
	src := "I CAN HAS STDIO?\n\nHAI ME TEH FUNCSHUN MAIN\n    SAYZ WIT 1\nKTHXBAI\n\nGARBAGE HERE\n"
	sa := NewSemanticAnalyzer("file:///pos.olol", src)
	if err := sa.AnalyzeDocument(context.Background()); err != nil {
		t.Fatal(err)
	}
	diags := sa.GetDiagnostics()
	if len(diags) == 0 {
		t.Fatal("expected a diagnostic")
	}
	found := false
	for _, d := range diags {
		if d.Range.Start.Line == 6 { // "GARBAGE HERE" is line 7 (0-based 6)
			found = true
			if d.Range.End.Character <= d.Range.Start.Character {
				t.Errorf("empty range: %+v", d.Range)
			}
		}
	}
	if !found {
		t.Errorf("no diagnostic on line 7; got %+v", diags)
	}
}

func TestParseErrorPositionsMultiple(t *testing.T) {
	src := "HAI ME TEH FUNCSHUN MAIN\n    SAYZ WIT 1 MOAR\nKTHXBAI\n\nHAI ME TEH CLAS\n"
	sa := NewSemanticAnalyzer("file:///pos2.olol", src)
	_ = sa.AnalyzeDocument(context.Background())
	for _, d := range sa.GetDiagnostics() {
		if d.Range.Start.Line == 0 && d.Range.Start.Character == 0 {
			t.Errorf("diagnostic still at 0:0: %q", d.Message)
		}
		if int(d.Range.Start.Line) > 5 {
			t.Errorf("diagnostic beyond end of document: %+v", d.Range)
		}
	}
}
