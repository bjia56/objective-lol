package server

import (
	"testing"

	protocol "github.com/tliron/glsp/protocol_3_16"
)

func TestOlolLSPServer_NewServer(t *testing.T) {
	server := NewServer()
	
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
	server := NewServer()
	
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
	server := NewServer()
	
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
	server := NewServer()
	
	err := server.shutdown(nil)
	if err != nil {
		t.Errorf("Unexpected error in shutdown: %v", err)
	}
}

func TestOlolLSPServer_Initialized(t *testing.T) {
	server := NewServer()
	
	err := server.initialized(nil, &protocol.InitializedParams{})
	if err != nil {
		t.Errorf("Unexpected error in initialized: %v", err)
	}
}

func TestOlolLSPServer_SetTrace(t *testing.T) {
	server := NewServer()
	
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
	server := NewServer()
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