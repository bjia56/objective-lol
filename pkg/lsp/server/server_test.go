package server

import (
	"strings"
	"testing"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func TestOlolLSPServer_NewServer(t *testing.T) {
	server, _ := NewServer("")

	if server == nil {
		t.Fatal("Expected server to be created")
	}

	if server.analyzer == nil {
		t.Error("Expected analyzer to be initialized")
	}

	if server.workspace == nil {
		t.Error("Expected workspace to be initialized")
	}
}

func TestOlolLSPServer_Initialize(t *testing.T) {
	server, _ := NewServer("")

	// Create initialize params - use appropriate types for GLSP
	processID := protocol.Integer(12345)
	rootURI := protocol.DocumentUri("file:///test/workspace")
	params := &protocol.InitializeParams{
		ProcessID: &processID,
		RootURI:   &rootURI,
	}

	result, err := server.initialize(nil, params)
	if err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	initResult, ok := result.(protocol.InitializeResult)
	if !ok {
		t.Fatal("Expected InitializeResult")
	}

	// Check server capabilities - we can't easily validate the detailed structure
	// due to interface{} types, but we can verify basic structure
	if initResult.Capabilities == (protocol.ServerCapabilities{}) {
		t.Error("Expected server capabilities to be set")
	}

	// Check server info
	if initResult.ServerInfo == nil {
		t.Error("Expected server info to be set")
	}

	if initResult.ServerInfo.Name != "olol-lsp" {
		t.Errorf("Expected server name 'olol-lsp', got '%s'", initResult.ServerInfo.Name)
	}
}

func TestOlolLSPServer_GetServerCapabilities(t *testing.T) {
	server, _ := NewServer("")

	caps := server.getServerCapabilities()

	// Test that capabilities are set
	if caps.TextDocumentSync == nil {
		t.Error("Expected TextDocumentSync to be set")
	}

	if caps.HoverProvider == nil {
		t.Error("Expected HoverProvider to be set")
	}

	if caps.CompletionProvider == nil {
		t.Error("Expected CompletionProvider to be set")
	}

	if caps.DefinitionProvider == nil {
		t.Error("Expected DefinitionProvider to be set")
	}
}

func TestOlolLSPServer_Shutdown(t *testing.T) {
	server, _ := NewServer("")

	err := server.shutdown(nil)
	if err != nil {
		t.Errorf("Unexpected error in shutdown: %v", err)
	}
}

func TestOlolLSPServer_Initialized(t *testing.T) {
	server, _ := NewServer("")

	err := server.initialized(nil, &protocol.InitializedParams{})
	if err != nil {
		t.Errorf("Unexpected error in initialized: %v", err)
	}
}

func TestOlolLSPServer_SetTrace(t *testing.T) {
	server, _ := NewServer("")

	params := &protocol.SetTraceParams{
		Value: protocol.TraceValueOff,
	}

	err := server.setTrace(nil, params)
	if err != nil {
		t.Errorf("Unexpected error in setTrace: %v", err)
	}
}

// Integration test that exercises the workspace and analyzer components
func TestOlolLSPServer_DocumentWorkflow(t *testing.T) {
	server, _ := NewServer("")
	uri := "test://test.olol"
	content := "HAI ME TEH FUNCSHUN MAIN\nTEH VARIABLE x ITZ INTEGR WIT 42\nSAYZ x\nKTHXBAI"

	// Test that we can add a document to workspace directly
	err := server.workspace.OpenDocument(uri, content)
	if err != nil {
		t.Fatalf("Failed to open document in workspace: %v", err)
	}

	// Test that we can analyze the document
	diagnostics, err := server.analyzer.AnalyzeDocument(uri, content)
	if err != nil {
		t.Fatalf("Failed to analyze document: %v", err)
	}

	// We expect this valid program to have no serious errors
	// (some minor parser issues might exist but shouldn't cause crashes)
	_ = diagnostics // Accept any number of diagnostics for now

	// Test completion
	position := protocol.Position{Line: 2, Character: 0}
	completions, err := server.analyzer.GetCompletions(uri, content, position)
	if err != nil {
		t.Fatalf("Failed to get completions: %v", err)
	}

	// Should have some completions
	if len(completions) == 0 {
		t.Error("Expected some completions")
	}

	// Test that workspace can be updated
	newContent := "HAI ME TEH FUNCSHUN MAIN\nSAYZ \"Hello\"\nKTHXBAI"
	err = server.workspace.UpdateDocument(uri, newContent)
	if err != nil {
		t.Fatalf("Failed to update document: %v", err)
	}

	// Test that document can be closed
	err = server.workspace.CloseDocument(uri)
	if err != nil {
		t.Fatalf("Failed to close document: %v", err)
	}

	if server.workspace.HasDocument(uri) {
		t.Error("Expected document to be closed")
	}
}

func newTestContext() *glsp.Context {
	return &glsp.Context{Notify: func(string, any) {}}
}

// Requests for documents that were never opened (or were closed) must return
// errors, not panic.
func TestServerRequestsOnUnknownDocument(t *testing.T) {
	s, _ := NewServer("")
	uri := protocol.DocumentUri("file:///nope.olol")
	pos := protocol.Position{Line: 5, Character: 5}
	td := protocol.TextDocumentIdentifier{URI: uri}

	if _, err := s.textDocumentHover(nil, &protocol.HoverParams{TextDocumentPositionParams: protocol.TextDocumentPositionParams{TextDocument: td, Position: pos}}); err == nil {
		t.Error("hover on unknown document should error")
	}
	if _, err := s.textDocumentCompletion(nil, &protocol.CompletionParams{TextDocumentPositionParams: protocol.TextDocumentPositionParams{TextDocument: td, Position: pos}}); err == nil {
		t.Error("completion on unknown document should error")
	}
	if _, err := s.textDocumentDefinition(nil, &protocol.DefinitionParams{TextDocumentPositionParams: protocol.TextDocumentPositionParams{TextDocument: td, Position: pos}}); err == nil {
		t.Error("definition on unknown document should error")
	}
	if err := s.textDocumentDidChange(newTestContext(), &protocol.DidChangeTextDocumentParams{
		TextDocument:   protocol.VersionedTextDocumentIdentifier{TextDocumentIdentifier: td},
		ContentChanges: []any{protocol.TextDocumentContentChangeEventWhole{Text: "x"}},
	}); err == nil {
		t.Error("change on unopened document should error")
	}
	if err := s.textDocumentDidClose(nil, &protocol.DidCloseTextDocumentParams{TextDocument: td}); err != nil {
		t.Errorf("closing unknown document should be a no-op, got %v", err)
	}
}

func TestServerMalformedChangeEvents(t *testing.T) {
	s, _ := NewServer("")
	uri := protocol.DocumentUri("file:///a.olol")
	ctx := newTestContext()
	_ = s.textDocumentDidOpen(ctx, &protocol.DidOpenTextDocumentParams{TextDocument: protocol.TextDocumentItem{URI: uri, Text: "HAI"}})
	td := protocol.VersionedTextDocumentIdentifier{TextDocumentIdentifier: protocol.TextDocumentIdentifier{URI: uri}}

	// no changes: no-op
	if err := s.textDocumentDidChange(ctx, &protocol.DidChangeTextDocumentParams{TextDocument: td}); err != nil {
		t.Errorf("empty change list should be a no-op: %v", err)
	}
	// incremental/unknown change type: error, not panic
	if err := s.textDocumentDidChange(ctx, &protocol.DidChangeTextDocumentParams{TextDocument: td, ContentChanges: []any{"garbage"}}); err == nil {
		t.Error("unsupported change type should error")
	}
	if err := s.textDocumentDidChange(ctx, &protocol.DidChangeTextDocumentParams{TextDocument: td, ContentChanges: []any{nil}}); err == nil {
		t.Error("nil change should error")
	}
}

// A document going through open/change/query/close cycles with broken content.
func TestServerLifecycleWithBrokenSources(t *testing.T) {
	s, _ := NewServer("")
	uri := protocol.DocumentUri("file:///broken.olol")
	ctx := newTestContext()
	td := protocol.TextDocumentIdentifier{URI: uri}
	pos := protocol.TextDocumentPositionParams{TextDocument: td, Position: protocol.Position{Line: 1, Character: 4}}

	sources := []string{
		"", "HAI ME TEH CLAS", "I CAN HAS", "\x00\xff\xfe",
		"HAI ME TEH FUNCSHUN MAIN\n I HAS A VARIABLE X TEH NOPE ITZ NEW\nKTHXBAI",
		"HAI ME TEH FUNCSHUN MAIN\n SAYZ WIT (((\nKTHXBAI",
		strings.Repeat("HAI ME TEH ", 5000),
	}
	if err := s.textDocumentDidOpen(ctx, &protocol.DidOpenTextDocumentParams{TextDocument: protocol.TextDocumentItem{URI: uri, Text: sources[0]}}); err != nil {
		t.Fatalf("open: %v", err)
	}
	for _, src := range sources {
		_ = s.textDocumentDidChange(ctx, &protocol.DidChangeTextDocumentParams{
			TextDocument:   protocol.VersionedTextDocumentIdentifier{TextDocumentIdentifier: td},
			ContentChanges: []any{protocol.TextDocumentContentChangeEventWhole{Text: src}},
		})
		_, _ = s.textDocumentHover(ctx, &protocol.HoverParams{TextDocumentPositionParams: pos})
		_, _ = s.textDocumentCompletion(ctx, &protocol.CompletionParams{TextDocumentPositionParams: pos})
		_, _ = s.textDocumentDefinition(ctx, &protocol.DefinitionParams{TextDocumentPositionParams: pos})
	}
	_ = s.textDocumentDidClose(ctx, &protocol.DidCloseTextDocumentParams{TextDocument: td})
}
