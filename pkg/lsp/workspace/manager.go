package workspace

import (
	"fmt"
	"sync"
)

// Manager handles document lifecycle and caching for the LSP server
type Manager struct {
	documents map[string]*Document
	mutex     sync.RWMutex
}

// Document represents an open document in the workspace
type Document struct {
	URI     string
	Content string
	Version int
}

// NewManager creates a new workspace manager
func NewManager() *Manager {
	return &Manager{
		documents: make(map[string]*Document),
	}
}

// OpenDocument opens a document in the workspace
func (m *Manager) OpenDocument(uri, content string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.documents[uri] = &Document{
		URI:     uri,
		Content: content,
		Version: 1,
	}

	return nil
}

// UpdateDocument updates a document's content
func (m *Manager) UpdateDocument(uri, content string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	doc, exists := m.documents[uri]
	if !exists {
		return fmt.Errorf("document not found: %s", uri)
	}

	doc.Content = content
	doc.Version++

	return nil
}

// CloseDocument closes a document in the workspace
func (m *Manager) CloseDocument(uri string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.documents, uri)
	return nil
}

// GetDocument retrieves a document's content
func (m *Manager) GetDocument(uri string) (string, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	doc, exists := m.documents[uri]
	if !exists {
		return "", fmt.Errorf("document not found: %s", uri)
	}

	return doc.Content, nil
}

// GetDocumentVersion retrieves a document's version
func (m *Manager) GetDocumentVersion(uri string) (int, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	doc, exists := m.documents[uri]
	if !exists {
		return 0, fmt.Errorf("document not found: %s", uri)
	}

	return doc.Version, nil
}

// GetAllDocuments returns all open documents
func (m *Manager) GetAllDocuments() map[string]*Document {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// Create a copy to avoid race conditions
	result := make(map[string]*Document)
	for uri, doc := range m.documents {
		result[uri] = &Document{
			URI:     doc.URI,
			Content: doc.Content,
			Version: doc.Version,
		}
	}

	return result
}

// HasDocument checks if a document is open in the workspace
func (m *Manager) HasDocument(uri string) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	_, exists := m.documents[uri]
	return exists
}
