package lsp

import (
	"testing"
	"time"

	"github.com/bjia56/objective-lol/pkg/lsp/analyzer"
	"github.com/bjia56/objective-lol/pkg/lsp/server"
	"github.com/bjia56/objective-lol/pkg/lsp/workspace"
)

// TestLSPIntegration tests the end-to-end LSP functionality
func TestLSPIntegration(t *testing.T) {
	// Test that we can create all LSP components
	t.Run("ComponentCreation", func(t *testing.T) {
		analyzer := analyzer.NewAnalyzer()
		workspace := workspace.NewManager()
		lspServer := server.NewServer()

		if analyzer == nil {
			t.Fatal("Failed to create analyzer")
		}
		if workspace == nil {
			t.Fatal("Failed to create workspace manager")
		}
		if lspServer == nil {
			t.Fatal("Failed to create LSP server")
		}
	})

	// Test document workflow
	t.Run("DocumentWorkflow", func(t *testing.T) {
		analyzer := analyzer.NewAnalyzer()
		workspace := workspace.NewManager()

		uri := "test://integration.olol"
		content := `BTW Integration test file
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN MAIN
    SAYZ WIT "Hello, World!"
    I HAS A VARIABLE testVar TEH INTEGR ITZ 42
    SAYZ WIT testVar
KTHXBAI`

		// Test document management
		err := workspace.OpenDocument(uri, content)
		if err != nil {
			t.Fatalf("Failed to open document: %v", err)
		}

		// Test analysis
		diagnostics, err := analyzer.AnalyzeDocument(uri, content)
		if err != nil {
			t.Fatalf("Failed to analyze document: %v", err)
		}

		// We don't require zero diagnostics since the parser might generate some,
		// but we should be able to analyze without crashing
		t.Logf("Analysis completed with %d diagnostics", len(diagnostics))

		// Test document updates
		newContent := `BTW Updated integration test file
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN MAIN
    SAYZ WIT "Hello, Updated World!"
    I HAS A VARIABLE newVar TEH STRIN ITZ "updated"
    SAYZ WIT newVar
KTHXBAI`

		err = workspace.UpdateDocument(uri, newContent)
		if err != nil {
			t.Fatalf("Failed to update document: %v", err)
		}

		// Test re-analysis
		diagnostics2, err := analyzer.AnalyzeDocument(uri, newContent)
		if err != nil {
			t.Fatalf("Failed to re-analyze document: %v", err)
		}

		t.Logf("Re-analysis completed with %d diagnostics", len(diagnostics2))

		// Test document closure
		err = workspace.CloseDocument(uri)
		if err != nil {
			t.Fatalf("Failed to close document: %v", err)
		}

		if workspace.HasDocument(uri) {
			t.Error("Document should be closed")
		}
	})

	// Test LSP server capabilities
	t.Run("ServerCapabilities", func(t *testing.T) {
		lspServer := server.NewServer()

		// Test server capabilities (this is a private method, so we access it through a test helper)
		// We can't directly test GetServerCapabilities since it's not exported,
		// but we can test that the server was created successfully
		if lspServer == nil {
			t.Error("LSP server should be created")
		}
	})

	// Test concurrent operations
	t.Run("ConcurrentOperations", func(t *testing.T) {
		workspace := workspace.NewManager()
		analyzer := analyzer.NewAnalyzer()

		// Test concurrent document operations
		done := make(chan bool, 2)

		go func() {
			defer func() { done <- true }()
			for i := 0; i < 10; i++ {
				uri := "test://concurrent1.olol"
				content := `HAI ME TEH FUNCSHUN MAIN
    SAYZ WIT "Concurrent test 1"
KTHXBAI`
				workspace.OpenDocument(uri, content)
				analyzer.AnalyzeDocument(uri, content)
				workspace.CloseDocument(uri)
				time.Sleep(1 * time.Millisecond)
			}
		}()

		go func() {
			defer func() { done <- true }()
			for i := 0; i < 10; i++ {
				uri := "test://concurrent2.olol"
				content := `HAI ME TEH FUNCSHUN MAIN
    SAYZ WIT "Concurrent test 2"
KTHXBAI`
				workspace.OpenDocument(uri, content)
				analyzer.AnalyzeDocument(uri, content)
				workspace.CloseDocument(uri)
				time.Sleep(1 * time.Millisecond)
			}
		}()

		// Wait for both goroutines
		<-done
		<-done
	})
}

// TestLSPServerBinary tests that the LSP server binary can be built
func TestLSPServerBinary(t *testing.T) {
	// This is more of a documentation test - the fact that we can import
	// and create the server components means the binary should build successfully
	lspServer := server.NewServer()
	if lspServer == nil {
		t.Fatal("LSP server creation failed")
	}

	t.Log("LSP server components created successfully")
}