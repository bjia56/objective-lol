package workspace

import (
	"testing"
)

func TestManager_DocumentLifecycle(t *testing.T) {
	manager := NewManager()
	uri := "test://test.olol"
	content := "HAI\nSAYZ \"Hello\"\nKTHXBYE"

	// Test opening a document
	err := manager.OpenDocument(uri, content)
	if err != nil {
		t.Fatalf("Failed to open document: %v", err)
	}

	// Test checking if document exists
	if !manager.HasDocument(uri) {
		t.Error("Expected document to exist after opening")
	}

	// Test getting document content
	retrievedContent, err := manager.GetDocument(uri)
	if err != nil {
		t.Fatalf("Failed to get document: %v", err)
	}
	if retrievedContent != content {
		t.Errorf("Expected content %q, got %q", content, retrievedContent)
	}

	// Test getting document version
	version, err := manager.GetDocumentVersion(uri)
	if err != nil {
		t.Fatalf("Failed to get document version: %v", err)
	}
	if version != 1 {
		t.Errorf("Expected version 1, got %d", version)
	}

	// Test updating document
	newContent := "HAI\nSAYZ \"Updated\"\nKTHXBYE"
	err = manager.UpdateDocument(uri, newContent)
	if err != nil {
		t.Fatalf("Failed to update document: %v", err)
	}

	// Check updated content and version
	retrievedContent, err = manager.GetDocument(uri)
	if err != nil {
		t.Fatalf("Failed to get updated document: %v", err)
	}
	if retrievedContent != newContent {
		t.Errorf("Expected updated content %q, got %q", newContent, retrievedContent)
	}

	version, err = manager.GetDocumentVersion(uri)
	if err != nil {
		t.Fatalf("Failed to get updated document version: %v", err)
	}
	if version != 2 {
		t.Errorf("Expected version 2 after update, got %d", version)
	}

	// Test closing document
	err = manager.CloseDocument(uri)
	if err != nil {
		t.Fatalf("Failed to close document: %v", err)
	}

	// Test that document no longer exists
	if manager.HasDocument(uri) {
		t.Error("Expected document to not exist after closing")
	}

	// Test getting non-existent document
	_, err = manager.GetDocument(uri)
	if err == nil {
		t.Error("Expected error when getting closed document")
	}
}

func TestManager_MultipleDocuments(t *testing.T) {
	manager := NewManager()

	documents := map[string]string{
		"test://file1.olol": "HAI\nSAYZ \"File 1\"\nKTHXBYE",
		"test://file2.olol": "HAI\nSAYZ \"File 2\"\nKTHXBYE",
		"test://file3.olol": "HAI\nSAYZ \"File 3\"\nKTHXBYE",
	}

	// Open multiple documents
	for uri, content := range documents {
		err := manager.OpenDocument(uri, content)
		if err != nil {
			t.Fatalf("Failed to open document %s: %v", uri, err)
		}
	}

	// Test GetAllDocuments
	allDocs := manager.GetAllDocuments()
	if len(allDocs) != len(documents) {
		t.Errorf("Expected %d documents, got %d", len(documents), len(allDocs))
	}

	// Verify each document exists and has correct content
	for uri, expectedContent := range documents {
		if !manager.HasDocument(uri) {
			t.Errorf("Document %s should exist", uri)
		}

		content, err := manager.GetDocument(uri)
		if err != nil {
			t.Errorf("Failed to get document %s: %v", uri, err)
		}
		if content != expectedContent {
			t.Errorf("Document %s: expected content %q, got %q", uri, expectedContent, content)
		}

		doc, exists := allDocs[uri]
		if !exists {
			t.Errorf("Document %s should exist in GetAllDocuments result", uri)
		}
		if doc.Content != expectedContent {
			t.Errorf("Document %s in GetAllDocuments: expected content %q, got %q", uri, expectedContent, doc.Content)
		}
	}

	// Close one document and verify others remain
	firstURI := "test://file1.olol"
	err := manager.CloseDocument(firstURI)
	if err != nil {
		t.Fatalf("Failed to close document %s: %v", firstURI, err)
	}

	if manager.HasDocument(firstURI) {
		t.Error("First document should not exist after closing")
	}

	// Verify other documents still exist
	for uri := range documents {
		if uri != firstURI && !manager.HasDocument(uri) {
			t.Errorf("Document %s should still exist", uri)
		}
	}
}

func TestManager_UpdateNonExistentDocument(t *testing.T) {
	manager := NewManager()
	uri := "test://nonexistent.olol"

	err := manager.UpdateDocument(uri, "some content")
	if err == nil {
		t.Error("Expected error when updating non-existent document")
	}
}

func TestManager_ThreadSafety(t *testing.T) {
	manager := NewManager()
	uri := "test://concurrent.olol"

	// Test concurrent operations
	done := make(chan bool, 2)

	// Goroutine 1: Open and close document repeatedly
	go func() {
		for i := 0; i < 100; i++ {
			manager.OpenDocument(uri, "content")
			manager.CloseDocument(uri)
		}
		done <- true
	}()

	// Goroutine 2: Check if document exists repeatedly
	go func() {
		for i := 0; i < 100; i++ {
			manager.HasDocument(uri)
			manager.GetAllDocuments()
		}
		done <- true
	}()

	// Wait for both goroutines to complete
	<-done
	<-done

	// If we get here without panicking, the test passes
}
