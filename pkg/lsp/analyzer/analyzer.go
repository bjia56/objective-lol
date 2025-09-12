package analyzer

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"

	lru "github.com/hashicorp/golang-lru"
	protocol "github.com/tliron/glsp/protocol_3_16"

	"github.com/bjia56/objective-lol/pkg/ast"
)

// Analyzer provides semantic analysis for LSP features
type Analyzer struct {
	semanticCache     map[string]*lru.Cache // Cache for enhanced semantic analyzers
	semanticCacheLock sync.Mutex            // Mutex for thread-safe access to semanticCache
}

// SymbolTable represents symbols in a document
type SymbolTable struct {
	Symbols []Symbol
	URI     string
}

// Symbol represents a symbol (variable, function, class, etc.)
type Symbol struct {
	Name     string
	Kind     SymbolKind
	Type     string
	Position ast.PositionInfo
	Range    protocol.Range
}

// SymbolKind represents the kind of symbol
type SymbolKind int

const (
	SymbolKindUnknown SymbolKind = iota
	SymbolKindVariable
	SymbolKindFunction
	SymbolKindClass
	SymbolKindParameter
	SymbolKindImport
)

func hashContent(content string) string {
	sum := sha256.Sum256([]byte(content))
	return hex.EncodeToString(sum[:])
}

// NewAnalyzer creates a new analyzer
func NewAnalyzer() *Analyzer {
	return &Analyzer{
		semanticCache: make(map[string]*lru.Cache),
	}
}

func (a *Analyzer) getOrCreateSemanticAnalyzer(uri, content string) (*SemanticAnalyzer, error) {
	a.semanticCacheLock.Lock()
	defer a.semanticCacheLock.Unlock()

	contentKey := hashContent(content)

	// Get or create LRU cache for this URI
	lruCache, exists := a.semanticCache[uri]
	if !exists {
		var err error
		lruCache, err = lru.New(5) // e.g., cache up to 5 content versions per URI
		if err != nil {
			return nil, err
		}
		a.semanticCache[uri] = lruCache
	}

	if analyzer, ok := lruCache.Get(contentKey); ok {
		return analyzer.(*SemanticAnalyzer), nil
	}

	semanticAnalyzer := NewSemanticAnalyzer(uri, content)
	ctx := context.Background()
	if err := semanticAnalyzer.AnalyzeDocument(ctx); err != nil {
		return nil, fmt.Errorf("semantic analysis failed: %v", err)
	}

	lruCache.Add(contentKey, semanticAnalyzer)
	return semanticAnalyzer, nil
}

// AnalyzeDocument analyzes a document and returns diagnostics
func (a *Analyzer) AnalyzeDocument(uri, content string) ([]protocol.Diagnostic, error) {
	semanticAnalyzer, err := a.getOrCreateSemanticAnalyzer(uri, content)
	if err != nil {
		return nil, err
	}

	// Return diagnostics from enhanced analyzer
	return semanticAnalyzer.GetDiagnostics(), nil
}

// GetHoverInfo provides hover information for a position
func (a *Analyzer) GetHoverInfo(uri, content string, position protocol.Position) (*protocol.Hover, error) {
	semanticAnalyzer, err := a.getOrCreateSemanticAnalyzer(uri, content)
	if err != nil {
		return nil, err
	}
	return semanticAnalyzer.GetHoverInfo(position), nil
}

// GetCompletions provides completion items for a position
func (a *Analyzer) GetCompletions(uri, content string, position protocol.Position) ([]protocol.CompletionItem, error) {
	semanticAnalyzer, err := a.getOrCreateSemanticAnalyzer(uri, content)
	if err != nil {
		return nil, err
	}
	return semanticAnalyzer.GetCompletionItems(position), nil
}

// GetDefinition provides definition location for a position
func (a *Analyzer) GetDefinition(uri, content string, position protocol.Position) (*protocol.Location, error) {
	semanticAnalyzer, err := a.getOrCreateSemanticAnalyzer(uri, content)
	if err != nil {
		return nil, err
	}
	return semanticAnalyzer.GetDefinitionLocation(position), nil
}
