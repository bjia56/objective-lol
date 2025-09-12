package server

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"

	"github.com/bjia56/objective-lol/pkg/lsp/analyzer"
	"github.com/bjia56/objective-lol/pkg/lsp/workspace"
)

// OlolLSPServer represents the Objective-LOL Language Server
type OlolLSPServer struct {
	analyzer  *analyzer.Analyzer
	workspace *workspace.Manager
	logFile   *os.File
}

// NewServer creates a new Objective-LOL LSP server
func NewServer(logFileName string) (*OlolLSPServer, error) {
	var logFile *os.File
	var err error
	if logFileName != "" {
		logFile, err = os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}
	}
	return &OlolLSPServer{
		analyzer:  analyzer.NewAnalyzer(),
		workspace: workspace.NewManager(),
		logFile:   logFile,
	}, nil
}

// log logs a request with timestamp
func (s *OlolLSPServer) log(method string, args ...interface{}) {
	if s.logFile != nil {
		timestamp := time.Now().Format(time.RFC3339)
		if len(args) > 0 {
			if jsonData, err := json.Marshal(args); err == nil {
				fmt.Fprintf(s.logFile, "%s: %s %s\n", timestamp, method, string(jsonData))
			} else {
				fmt.Fprintf(s.logFile, "%s: %s\n", timestamp, method)
			}
		} else {
			fmt.Fprintf(s.logFile, "%s: %s\n", timestamp, method)
		}
	}
}

// Start initializes and starts the LSP server
func (s *OlolLSPServer) Start() error {
	// Create GLSP handler using the correct protocol structure
	handler := protocol.Handler{
		Initialize:             s.initialize,
		Initialized:            s.initialized,
		Shutdown:               s.shutdown,
		SetTrace:               s.setTrace,
		TextDocumentDidOpen:    s.textDocumentDidOpen,
		TextDocumentDidChange:  s.textDocumentDidChange,
		TextDocumentDidClose:   s.textDocumentDidClose,
		TextDocumentHover:      s.textDocumentHover,
		TextDocumentCompletion: s.textDocumentCompletion,
		TextDocumentDefinition: s.textDocumentDefinition,
	}

	// Create and start the server
	lspServer := server.NewServer(&handler, "olol-lsp", false)
	return lspServer.RunStdio()
}

// Initialize handles the initialize request
func (s *OlolLSPServer) initialize(context *glsp.Context, params *protocol.InitializeParams) (interface{}, error) {
	s.log("initialize", params)

	// Define server capabilities
	capabilities := s.getServerCapabilities()

	return protocol.InitializeResult{
		Capabilities: capabilities,
		ServerInfo: &protocol.InitializeResultServerInfo{
			Name:    "olol-lsp",
			Version: &[]string{"0.1.0"}[0],
		},
	}, nil
}

// Initialized handles the initialized notification
func (s *OlolLSPServer) initialized(context *glsp.Context, params *protocol.InitializedParams) error {
	s.log("initialized", params)

	return nil
}

// Shutdown handles the shutdown request
func (s *OlolLSPServer) shutdown(context *glsp.Context) error {
	s.log("shutdown")

	return nil
}

// SetTrace handles the setTrace notification
func (s *OlolLSPServer) setTrace(context *glsp.Context, params *protocol.SetTraceParams) error {
	s.log("setTrace", params)

	return nil
}

// getServerCapabilities returns the server's capabilities
func (s *OlolLSPServer) getServerCapabilities() protocol.ServerCapabilities {
	return protocol.ServerCapabilities{
		// Text document sync
		TextDocumentSync: protocol.TextDocumentSyncOptions{
			OpenClose: &[]bool{true}[0],
			Change:    &[]protocol.TextDocumentSyncKind{protocol.TextDocumentSyncKindFull}[0],
		},

		// Hover support
		HoverProvider: &[]bool{true}[0],

		// Completion support
		CompletionProvider: &protocol.CompletionOptions{
			ResolveProvider: &[]bool{false}[0],
			TriggerCharacters: []string{
				" ", // For keyword completion
			},
		},

		// Definition support
		DefinitionProvider: &[]bool{true}[0],

		// Future capabilities can be added here
	}
}

// Text Document handlers

func (s *OlolLSPServer) textDocumentDidOpen(context *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	s.log("textDocumentDidOpen", params)

	uri := params.TextDocument.URI
	content := params.TextDocument.Text

	// Store document in workspace
	err := s.workspace.OpenDocument(uri, content)
	if err != nil {
		return fmt.Errorf("failed to open document: %w", err)
	}

	// Analyze document and send diagnostics
	diagnostics, err := s.analyzer.AnalyzeDocument(uri, content)
	if err != nil {
		return fmt.Errorf("failed to analyze document: %w", err)
	}

	// Send diagnostics to client
	context.Notify(protocol.ServerTextDocumentPublishDiagnostics, &protocol.PublishDiagnosticsParams{
		URI:         uri,
		Diagnostics: diagnostics,
	})

	return nil
}

func (s *OlolLSPServer) textDocumentDidChange(context *glsp.Context, params *protocol.DidChangeTextDocumentParams) (err error) {
	s.log("textDocumentDidChange", params)

	uri := params.TextDocument.URI

	if len(params.ContentChanges) == 0 {
		return nil
	}

	// For full document sync, we take the last change
	change := params.ContentChanges[len(params.ContentChanges)-1]

	// Cast to access Text field
	var content string
	if changeEvent, ok := change.(protocol.TextDocumentContentChangeEventWhole); ok {
		content = changeEvent.Text
	} else {
		return fmt.Errorf("unexpected change event type")
	}

	// Update document in workspace
	err = s.workspace.UpdateDocument(uri, content)
	if err != nil {
		return fmt.Errorf("failed to update document: %w", err)
	}

	// Re-analyze document and send diagnostics
	diagnostics, err := s.analyzer.AnalyzeDocument(uri, content)
	if err != nil {
		return fmt.Errorf("failed to analyze document: %w", err)
	}

	// Send diagnostics to client
	context.Notify(protocol.ServerTextDocumentPublishDiagnostics, &protocol.PublishDiagnosticsParams{
		URI:         uri,
		Diagnostics: diagnostics,
	})

	return nil
}

func (s *OlolLSPServer) textDocumentDidClose(context *glsp.Context, params *protocol.DidCloseTextDocumentParams) error {
	s.log("textDocumentDidClose", params)

	uri := params.TextDocument.URI
	return s.workspace.CloseDocument(uri)
}

func (s *OlolLSPServer) textDocumentHover(context *glsp.Context, params *protocol.HoverParams) (*protocol.Hover, error) {
	s.log("textDocumentHover", params)

	uri := params.TextDocument.URI
	position := params.Position

	// Get document content
	content, err := s.workspace.GetDocument(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to get document: %w", err)
	}

	// Get hover information
	hoverInfo, err := s.analyzer.GetHoverInfo(uri, content, position)
	if err != nil {
		return nil, err
	}

	return hoverInfo, nil
}

func (s *OlolLSPServer) textDocumentCompletion(context *glsp.Context, params *protocol.CompletionParams) (interface{}, error) {
	s.log("textDocumentCompletion", params)

	uri := params.TextDocument.URI
	position := params.Position

	// Get document content
	content, err := s.workspace.GetDocument(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to get document: %w", err)
	}

	// Get completion items
	completions, err := s.analyzer.GetCompletions(uri, content, position)
	if err != nil {
		return nil, err
	}

	return protocol.CompletionList{
		IsIncomplete: false,
		Items:        completions,
	}, nil
}

func (s *OlolLSPServer) textDocumentDefinition(context *glsp.Context, params *protocol.DefinitionParams) (interface{}, error) {
	s.log("textDocumentDefinition", params)

	uri := params.TextDocument.URI
	position := params.Position

	// Get document content
	content, err := s.workspace.GetDocument(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to get document: %w", err)
	}

	// Get definition location
	location, err := s.analyzer.GetDefinition(uri, content, position)
	if err != nil {
		return nil, err
	}

	if location == nil {
		return nil, nil // No definition found
	}

	return location, nil
}
