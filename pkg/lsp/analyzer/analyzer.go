package analyzer

import (
	"fmt"
	"strings"

	protocol "github.com/tliron/glsp/protocol_3_16"

	"github.com/bjia56/objective-lol/pkg/ast"
	"github.com/bjia56/objective-lol/pkg/parser"
)

// Analyzer provides semantic analysis for LSP features
type Analyzer struct {
	symbolCache map[string]*SymbolTable
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
	SymbolKindVariable SymbolKind = iota
	SymbolKindFunction
	SymbolKindClass
	SymbolKindParameter
	SymbolKindImport
)

// NewAnalyzer creates a new analyzer
func NewAnalyzer() *Analyzer {
	return &Analyzer{
		symbolCache: make(map[string]*SymbolTable),
	}
}

// AnalyzeDocument analyzes a document and returns diagnostics
func (a *Analyzer) AnalyzeDocument(uri, content string) ([]protocol.Diagnostic, error) {
	var diagnostics []protocol.Diagnostic

	// Parse the document
	lexer := parser.NewLexer(content)
	p := parser.NewParser(lexer)
	program := p.ParseProgram()

	// Check for parse errors
	if errors := p.Errors(); len(errors) > 0 {
		for _, errorMsg := range errors {
			diagnostic := protocol.Diagnostic{
				Range: protocol.Range{
					Start: protocol.Position{Line: 0, Character: 0},
					End:   protocol.Position{Line: 0, Character: 0},
				},
				Severity: &[]protocol.DiagnosticSeverity{protocol.DiagnosticSeverityError}[0],
				Source:   &[]string{"olol-lsp"}[0],
				Message:  errorMsg,
			}
			diagnostics = append(diagnostics, diagnostic)
		}
	}

	// Build symbol table
	symbolTable := a.buildSymbolTable(uri, program)
	a.symbolCache[uri] = symbolTable

	// TODO: Add semantic analysis (undefined variables, type checking, etc.)

	return diagnostics, nil
}

// GetHoverInfo provides hover information for a position
func (a *Analyzer) GetHoverInfo(uri, content string, position protocol.Position) (*protocol.Hover, error) {
	// Get symbol at position
	symbol := a.getSymbolAtPosition(uri, position)
	if symbol == nil {
		return nil, nil
	}

	// Create hover content
	var contents []string
	contents = append(contents, fmt.Sprintf("**%s** (%s)", symbol.Name, a.symbolKindString(symbol.Kind)))
	if symbol.Type != "" {
		contents = append(contents, fmt.Sprintf("Type: `%s`", symbol.Type))
	}

	return &protocol.Hover{
		Contents: protocol.MarkupContent{
			Kind:  protocol.MarkupKindMarkdown,
			Value: strings.Join(contents, "\n\n"),
		},
		Range: &symbol.Range,
	}, nil
}

// GetCompletions provides completion items for a position
func (a *Analyzer) GetCompletions(uri, content string, position protocol.Position) ([]protocol.CompletionItem, error) {
	var completions []protocol.CompletionItem

	// Get symbol table for the document
	symbolTable, exists := a.symbolCache[uri]
	if !exists {
		// If not cached, analyze the document first
		_, err := a.AnalyzeDocument(uri, content)
		if err != nil {
			return nil, err
		}
		symbolTable = a.symbolCache[uri]
	}

	// Add keyword completions
	keywords := []string{
		"HAI", "KTHXBAI", "TEH", "VARIABLE", "FUNCSHUN", "CLAS",
		"KITTEH", "OF", "ITZ", "WIT", "AN", "DIS", "IZ", "NOPE",
		"WHILE", "GIVEZ", "UP", "NEW", "AS", "I", "CAN", "HAS",
		"MAYB", "OOPS", "OOPSIE", "ALWAYZ", "INTEGR", "DUBBLE",
		"STRIN", "BOOL", "YEZ", "NO", "NOTHIN",
	}

	for _, keyword := range keywords {
		completions = append(completions, protocol.CompletionItem{
			Label: keyword,
			Kind:  &[]protocol.CompletionItemKind{protocol.CompletionItemKindKeyword}[0],
		})
	}

	// Add symbol completions
	for _, symbol := range symbolTable.Symbols {
		kind := a.symbolKindToCompletionKind(symbol.Kind)
		completions = append(completions, protocol.CompletionItem{
			Label:  symbol.Name,
			Kind:   &kind,
			Detail: &symbol.Type,
		})
	}

	return completions, nil
}

// GetDefinition provides definition location for a position
func (a *Analyzer) GetDefinition(uri, content string, position protocol.Position) (*protocol.Location, error) {
	// Get symbol at position
	symbol := a.getSymbolAtPosition(uri, position)
	if symbol == nil {
		return nil, nil
	}

	return &protocol.Location{
		URI:   uri,
		Range: symbol.Range,
	}, nil
}

// Helper methods

func (a *Analyzer) buildSymbolTable(uri string, program *ast.ProgramNode) *SymbolTable {
	visitor := NewSymbolCollector(uri)
	program.Accept(visitor)
	return visitor.GetSymbolTable()
}

func (a *Analyzer) getSymbolAtPosition(uri string, position protocol.Position) *Symbol {
	symbolTable, exists := a.symbolCache[uri]
	if !exists {
		return nil
	}

	// Find symbol that contains the position
	for _, symbol := range symbolTable.Symbols {
		if a.positionInRange(position, symbol.Range) {
			return &symbol
		}
	}

	return nil
}

func (a *Analyzer) positionInRange(position protocol.Position, rang protocol.Range) bool {
	if position.Line < rang.Start.Line || position.Line > rang.End.Line {
		return false
	}
	if position.Line == rang.Start.Line && position.Character < rang.Start.Character {
		return false
	}
	if position.Line == rang.End.Line && position.Character > rang.End.Character {
		return false
	}
	return true
}

func (a *Analyzer) symbolKindString(kind SymbolKind) string {
	switch kind {
	case SymbolKindVariable:
		return "Variable"
	case SymbolKindFunction:
		return "Function"
	case SymbolKindClass:
		return "Class"
	case SymbolKindParameter:
		return "Parameter"
	case SymbolKindImport:
		return "Import"
	default:
		return "Unknown"
	}
}

func (a *Analyzer) symbolKindToCompletionKind(kind SymbolKind) protocol.CompletionItemKind {
	switch kind {
	case SymbolKindVariable:
		return protocol.CompletionItemKindVariable
	case SymbolKindFunction:
		return protocol.CompletionItemKindFunction
	case SymbolKindClass:
		return protocol.CompletionItemKindClass
	case SymbolKindParameter:
		return protocol.CompletionItemKindVariable
	case SymbolKindImport:
		return protocol.CompletionItemKindModule
	default:
		return protocol.CompletionItemKindText
	}
}

// positionInfoToRange converts AST PositionInfo to LSP Range
func (a *Analyzer) positionInfoToRange(pos ast.PositionInfo, length int) protocol.Range {
	return protocol.Range{
		Start: protocol.Position{
			Line:      uint32(pos.Line - 1),      // LSP is 0-based, parser is 1-based
			Character: uint32(pos.Column - 1),    // LSP is 0-based, parser is 1-based
		},
		End: protocol.Position{
			Line:      uint32(pos.Line - 1),
			Character: uint32(pos.Column - 1 + length),
		},
	}
}