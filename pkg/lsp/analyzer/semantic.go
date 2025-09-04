package analyzer

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	protocol "github.com/tliron/glsp/protocol_3_16"

	"github.com/bjia56/objective-lol/pkg/ast"
	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/modules"
	"github.com/bjia56/objective-lol/pkg/parser"
	"github.com/bjia56/objective-lol/pkg/stdlib"
)

// SemanticAnalyzer provides enhanced semantic analysis with environment awareness
type SemanticAnalyzer struct {
	uri               string
	content           string
	symbolTable       *EnhancedSymbolTable
	moduleResolver    *modules.ModuleResolver
	stdlibInitializer map[string]func(*environment.Environment, ...string) error

	// Analysis state
	environment      *environment.Environment
	environmentStack []*environment.Environment
	currentClass     string
	currentFile      string
	analysisErrors   []AnalysisError
}

// EnhancedSymbolTable represents symbols with enhanced metadata
type EnhancedSymbolTable struct {
	URI              string
	Symbols          []EnhancedSymbol
	Scopes           []ScopeInfo
	ImportedModules  []ModuleImport
	FunctionCalls    []FunctionCall
	DiagnosticsCache []protocol.Diagnostic
}

// EnhancedSymbol represents a symbol with full semantic context
type EnhancedSymbol struct {
	Name          string
	Kind          SymbolKind
	Type          string
	Position      ast.PositionInfo
	Range         protocol.Range
	Scope         ScopeType
	ScopeID       string
	Visibility    VisibilityType
	ParentClass   string
	QualifiedName string
	IsShared      bool
	SourceModule  string
	References    []ast.PositionInfo
	Documentation string
}

// ScopeInfo represents scope information
type ScopeInfo struct {
	ID           string
	Type         ScopeType
	Parent       string
	StartPos     ast.PositionInfo
	EndPos       ast.PositionInfo
	ClassName    string // For class scopes
	FunctionName string // For function scopes
}

// ModuleImport represents an imported module
type ModuleImport struct {
	ModuleName      string
	FilePath        string
	IsFileImport    bool
	ImportedSymbols []string
	Position        ast.PositionInfo
}

// ScopeType represents different types of scopes
type ScopeType int

const (
	ScopeTypeGlobal ScopeType = iota
	ScopeTypeFunction
	ScopeTypeBlock
	ScopeTypeClass
	ScopeTypeModule
)

// VisibilityType represents symbol visibility
type VisibilityType int

const (
	VisibilityPublic VisibilityType = iota
	VisibilityPrivate
	VisibilityShared
)

// FunctionCall represents a function call site with resolution info
type FunctionCall struct {
	CallSite     ast.PositionInfo   // Position of the function call
	FunctionName string             // Name of the called function
	CallType     FunctionCallType   // Type of call (global, method, etc.)
	Arguments    []ast.PositionInfo // Positions of arguments
	ResolvedTo   *EnhancedSymbol    // Resolved function symbol (if found)
	ObjectType   string             // For method calls, type of the object
	Range        protocol.Range     // LSP range for the call
}

// FunctionCallType represents different types of function calls
type FunctionCallType int

const (
	FunctionCallGlobal FunctionCallType = iota
	FunctionCallMethod
	FunctionCallBuiltin
)

// AnalysisError represents semantic analysis errors
type AnalysisError struct {
	Type     ErrorType
	Message  string
	Position ast.PositionInfo
	Severity protocol.DiagnosticSeverity
	Code     string
}

// ErrorType represents different types of analysis errors
type ErrorType int

const (
	ErrorTypeUndefinedSymbol ErrorType = iota
	ErrorTypeTypeError
	ErrorTypeVisibilityError
	ErrorTypeModuleError
	ErrorTypeCircularImport
	ErrorTypeDuplicateDeclaration
)

// NewSemanticAnalyzer creates a new semantic analyzer
func NewSemanticAnalyzer(uri, content string) *SemanticAnalyzer {
	// Use current working directory as base for module resolution
	workingDir, _ := filepath.Abs(".")

	// Get stdlib initializers (we'll need to create a way to access these)
	stdlibInit := make(map[string]func(*environment.Environment, ...string) error)
	for name, init := range stdlib.DefaultStdlibInitializers() {
		stdlibInit[name] = init
	}

	return &SemanticAnalyzer{
		uri:     uri,
		content: content,
		symbolTable: &EnhancedSymbolTable{
			URI:             uri,
			Symbols:         []EnhancedSymbol{},
			Scopes:          []ScopeInfo{},
			ImportedModules: []ModuleImport{},
		},
		moduleResolver:    modules.NewModuleResolver(workingDir),
		stdlibInitializer: stdlibInit,
		environment:       environment.NewEnvironment(nil),
		environmentStack:  []*environment.Environment{},
		analysisErrors:    []AnalysisError{},
	}
}

// AnalyzeDocument performs comprehensive semantic analysis
func (sa *SemanticAnalyzer) AnalyzeDocument(ctx context.Context) error {
	// Clear previous analysis results
	sa.symbolTable.Symbols = []EnhancedSymbol{}
	sa.symbolTable.Scopes = []ScopeInfo{}
	sa.symbolTable.ImportedModules = []ModuleImport{}
	sa.analysisErrors = []AnalysisError{}

	// Parse the document
	lexer := parser.NewLexer(sa.content)
	p := parser.NewParser(lexer)
	program := p.ParseProgram()

	// Check for parse errors first
	if errors := p.Errors(); len(errors) > 0 {
		for _, errorMsg := range errors {
			sa.analysisErrors = append(sa.analysisErrors, AnalysisError{
				Type:     ErrorTypeUndefinedSymbol, // Generic for parse errors
				Message:  errorMsg,
				Position: ast.PositionInfo{Line: 1, Column: 1}, // TODO: Get actual position
				Severity: protocol.DiagnosticSeverityError,
				Code:     "parse_error",
			})
		}
		// Don't return error - let parse errors be reported as diagnostics
		// Continue with analysis of what we can parse
	}

	// If program is nil due to complete parse failure, we can still provide diagnostics
	if program == nil {
		// Initialize stdlib symbols and return - we can't do much more without a valid AST
		sa.initializeStdlibSymbols()
		return nil
	}

	// Set current file for module resolution
	if sa.uri != "" {
		// Convert URI to file path for module resolution
		if path, found := strings.CutPrefix(sa.uri, "file://"); found {
			sa.currentFile = path
		} else {
			sa.currentFile = sa.uri
		}
	}

	// Initialize stdlib symbols
	sa.initializeStdlibSymbols()

	// Add global scope
	globalScope := ScopeInfo{
		ID:       "global",
		Type:     ScopeTypeGlobal,
		Parent:   "",
		StartPos: ast.PositionInfo{Line: 1, Column: 1},
		EndPos:   ast.PositionInfo{Line: 999999, Column: 999999}, // Will be updated
	}
	sa.symbolTable.Scopes = append(sa.symbolTable.Scopes, globalScope)

	// Perform multi-pass analysis similar to interpreter
	if err := sa.analyzeProgram(program); err != nil {
		return err
	}

	return nil
}

// initializeStdlibSymbols adds stdlib symbols to the analysis
func (sa *SemanticAnalyzer) initializeStdlibSymbols() {
	// Add global stdlib symbols
	for _, globalInit := range stdlib.DefaultGlobalInitializers() {
		for _, decl := range stdlib.GetStdlibDefinitions(globalInit) {
			symbol := EnhancedSymbol{
				Name:          decl.Name,
				Kind:          sa.stdlibKindToSymbolKind(decl.Kind),
				Type:          decl.Type,
				Position:      ast.PositionInfo{}, // Stdlib symbols have no position
				Range:         protocol.Range{},
				Scope:         ScopeTypeGlobal,
				ScopeID:       "global",
				Visibility:    VisibilityPublic,
				QualifiedName: decl.Name,
				SourceModule:  "stdlib",
			}
			sa.symbolTable.Symbols = append(sa.symbolTable.Symbols, symbol)
		}
	}
}

// stdlibKindToSymbolKind converts stdlib definition kind to symbol kind
func (sa *SemanticAnalyzer) stdlibKindToSymbolKind(kind stdlib.StdlibDefinitionKind) SymbolKind {
	switch kind {
	case stdlib.StdlibDefinitionKindFunction:
		return SymbolKindFunction
	case stdlib.StdlibDefinitionKindClass:
		return SymbolKindClass
	default:
		return SymbolKindUnknown
	}
}

// analyzeProgram performs the main program analysis
func (sa *SemanticAnalyzer) analyzeProgram(program *ast.ProgramNode) error {
	// Pass 1: Collect imports and process them
	for _, decl := range program.Declarations {
		if importNode, ok := decl.(*ast.ImportStatementNode); ok {
			if err := sa.analyzeImport(importNode); err != nil {
				// Add as analysis error but continue
				sa.analysisErrors = append(sa.analysisErrors, AnalysisError{
					Type:     ErrorTypeModuleError,
					Message:  err.Error(),
					Position: importNode.GetPosition(),
					Severity: protocol.DiagnosticSeverityError,
					Code:     "import_error",
				})
			}
		}
	}

	// Pass 2: Collect function and class declarations
	for _, decl := range program.Declarations {
		switch node := decl.(type) {
		case *ast.FunctionDeclarationNode:
			sa.analyzeFunctionDeclaration(node)
		case *ast.ClassDeclarationNode:
			sa.analyzeClassDeclaration(node)
		}
	}

	// Pass 3: Analyze variable declarations and other statements
	for _, decl := range program.Declarations {
		switch node := decl.(type) {
		case *ast.VariableDeclarationNode:
			sa.analyzeVariableDeclaration(node)
		}
	}

	return nil
}

// GetDiagnostics returns analysis errors as LSP diagnostics
func (sa *SemanticAnalyzer) GetDiagnostics() []protocol.Diagnostic {
	var diagnostics []protocol.Diagnostic

	for _, err := range sa.analysisErrors {
		diagnostic := protocol.Diagnostic{
			Range: protocol.Range{
				Start: protocol.Position{
					Line:      uint32(err.Position.Line - 1),   // LSP is 0-based
					Character: uint32(err.Position.Column - 1), // LSP is 0-based
				},
				End: protocol.Position{
					Line:      uint32(err.Position.Line - 1),
					Character: uint32(err.Position.Column - 1 + len(err.Message)),
				},
			},
			Severity: &err.Severity,
			Source:   &[]string{"olol-semantic"}[0],
			Message:  err.Message,
			// Code:     err.Code, // Omitted for now due to type complexity
		}
		diagnostics = append(diagnostics, diagnostic)
	}

	return diagnostics
}

// GetSymbolTable returns the enhanced symbol table
func (sa *SemanticAnalyzer) GetSymbolTable() *EnhancedSymbolTable {
	return sa.symbolTable
}

// analyzeImport analyzes import statements
func (sa *SemanticAnalyzer) analyzeImport(node *ast.ImportStatementNode) error {
	moduleImport := ModuleImport{
		ModuleName:      node.ModuleName,
		IsFileImport:    node.IsFileImport,
		ImportedSymbols: node.Declarations,
		Position:        node.GetPosition(),
	}

	if node.IsFileImport {
		// Handle file import
		importingDir := ""
		if sa.currentFile != "" {
			importingDir = filepath.Dir(sa.currentFile)
		}

		moduleAST, resolvedPath, err := sa.moduleResolver.LoadModuleFromWithPath(node.ModuleName, importingDir)
		if err != nil {
			return fmt.Errorf("failed to load module %s: %v", node.ModuleName, err)
		}

		moduleImport.FilePath = resolvedPath

		// Analyze the imported module to get its symbols
		if err := sa.analyzeImportedModule(moduleAST, resolvedPath, node.Declarations); err != nil {
			return fmt.Errorf("failed to analyze imported module %s: %v", node.ModuleName, err)
		}

	} else {
		// Handle built-in module import
		if moduleInit, exists := sa.stdlibInitializer[strings.ToUpper(node.ModuleName)]; exists {
			// Create a temporary environment to get the module symbols
			tempEnv := environment.NewEnvironment(nil)
			if err := moduleInit(tempEnv, node.Declarations...); err != nil {
				return fmt.Errorf("failed to initialize stdlib module %s: %v", node.ModuleName, err)
			}

			// Extract symbols from the temporary environment
			sa.extractSymbolsFromEnvironment(tempEnv, "stdlib:"+node.ModuleName, node.Declarations)
		} else {
			return fmt.Errorf("unknown built-in module: %s", node.ModuleName)
		}
	}

	sa.symbolTable.ImportedModules = append(sa.symbolTable.ImportedModules, moduleImport)
	return nil
}

// analyzeImportedModule analyzes a module and extracts its symbols
func (sa *SemanticAnalyzer) analyzeImportedModule(moduleAST *ast.ProgramNode, modulePath string, declarations []string) error {
	// Create a separate analyzer for the imported module
	moduleAnalyzer := NewSemanticAnalyzer(modulePath, "")
	moduleAnalyzer.currentFile = modulePath

	// Analyze the module declarations
	for _, decl := range moduleAST.Declarations {
		switch node := decl.(type) {
		case *ast.FunctionDeclarationNode:
			moduleAnalyzer.analyzeFunctionDeclaration(node)
		case *ast.ClassDeclarationNode:
			moduleAnalyzer.analyzeClassDeclaration(node)
		case *ast.VariableDeclarationNode:
			moduleAnalyzer.analyzeVariableDeclaration(node)
		}
	}

	// Import symbols from the module analyzer
	sa.importSymbolsFromModule(moduleAnalyzer.symbolTable, modulePath, declarations)

	return nil
}

// importSymbolsFromModule imports symbols from another module's symbol table
func (sa *SemanticAnalyzer) importSymbolsFromModule(moduleTable *EnhancedSymbolTable, modulePath string, declarations []string) {
	for _, symbol := range moduleTable.Symbols {
		// Skip private symbols (those starting with _)
		if strings.HasPrefix(symbol.Name, "_") {
			continue
		}

		// If specific declarations are requested, only import those
		if len(declarations) > 0 {
			found := false
			for _, decl := range declarations {
				if strings.EqualFold(symbol.Name, decl) {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// Create imported symbol
		importedSymbol := symbol
		importedSymbol.SourceModule = modulePath
		importedSymbol.Scope = ScopeTypeGlobal
		importedSymbol.ScopeID = "global"

		sa.symbolTable.Symbols = append(sa.symbolTable.Symbols, importedSymbol)
	}
}

// extractSymbolsFromEnvironment extracts symbols from an environment
func (sa *SemanticAnalyzer) extractSymbolsFromEnvironment(env *environment.Environment, sourceModule string, declarations []string) {
	// Extract functions
	for name, function := range env.GetAllFunctions() {
		if len(declarations) > 0 {
			found := false
			for _, decl := range declarations {
				if strings.EqualFold(name, decl) {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		symbol := EnhancedSymbol{
			Name:          name,
			Kind:          SymbolKindFunction,
			Type:          function.ReturnType,
			Position:      ast.PositionInfo{}, // Stdlib symbols have no position
			Range:         protocol.Range{},
			Scope:         ScopeTypeGlobal,
			ScopeID:       "global",
			Visibility:    VisibilityPublic,
			QualifiedName: name,
			SourceModule:  sourceModule,
		}
		sa.symbolTable.Symbols = append(sa.symbolTable.Symbols, symbol)
	}

	// Extract classes
	for name, class := range env.GetAllClasses() {
		if len(declarations) > 0 {
			found := false
			for _, decl := range declarations {
				if strings.EqualFold(name, decl) {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		symbol := EnhancedSymbol{
			Name:          name,
			Kind:          SymbolKindClass,
			Type:          class.Name,
			Position:      ast.PositionInfo{},
			Range:         protocol.Range{},
			Scope:         ScopeTypeGlobal,
			ScopeID:       "global",
			Visibility:    VisibilityPublic,
			QualifiedName: name,
			SourceModule:  sourceModule,
		}
		sa.symbolTable.Symbols = append(sa.symbolTable.Symbols, symbol)
	}
}

// analyzeFunctionDeclaration analyzes function declarations
func (sa *SemanticAnalyzer) analyzeFunctionDeclaration(node *ast.FunctionDeclarationNode) {
	functionName := strings.ToUpper(node.Name)

	// Create function scope
	scopeID := fmt.Sprintf("func_%s_%d_%d", functionName, node.GetPosition().Line, node.GetPosition().Column)
	funcScope := ScopeInfo{
		ID:           scopeID,
		Type:         ScopeTypeFunction,
		Parent:       sa.getCurrentScopeID(),
		StartPos:     node.GetPosition(),
		EndPos:       ast.PositionInfo{}, // TODO: Calculate end position
		FunctionName: functionName,
	}
	sa.symbolTable.Scopes = append(sa.symbolTable.Scopes, funcScope)

	// Add function symbol
	visibility := VisibilityPublic
	if strings.HasPrefix(functionName, "_") {
		visibility = VisibilityPrivate
	}

	funcSymbol := EnhancedSymbol{
		Name:          functionName,
		Kind:          SymbolKindFunction,
		Type:          strings.ToUpper(node.ReturnType),
		Position:      node.GetPosition(),
		Range:         sa.positionToRange(node.GetPosition(), len(node.Name)),
		Scope:         sa.getCurrentScopeType(),
		ScopeID:       sa.getCurrentScopeID(),
		Visibility:    visibility,
		ParentClass:   sa.currentClass,
		QualifiedName: sa.getQualifiedName(functionName),
		IsShared:      node.IsShared != nil && *node.IsShared,
	}
	sa.symbolTable.Symbols = append(sa.symbolTable.Symbols, funcSymbol)

	// Analyze function parameters
	for _, param := range node.Parameters {
		paramSymbol := EnhancedSymbol{
			Name:          strings.ToUpper(param.Name),
			Kind:          SymbolKindParameter,
			Type:          strings.ToUpper(param.Type),
			Position:      ast.PositionInfo{}, // TODO: Get parameter positions from parser
			Range:         protocol.Range{},
			Scope:         ScopeTypeFunction,
			ScopeID:       scopeID,
			Visibility:    VisibilityPrivate,
			QualifiedName: strings.ToUpper(param.Name),
		}
		sa.symbolTable.Symbols = append(sa.symbolTable.Symbols, paramSymbol)
	}

	// Analyze function body for local variables and nested symbols
	if node.Body != nil {
		sa.analyzeStatementBlock(node.Body, scopeID)
	}
}

// analyzeClassDeclaration analyzes class declarations
func (sa *SemanticAnalyzer) analyzeClassDeclaration(node *ast.ClassDeclarationNode) {
	className := strings.ToUpper(node.Name)

	// Create class scope
	scopeID := fmt.Sprintf("class_%s_%d_%d", className, node.GetPosition().Line, node.GetPosition().Column)
	classScope := ScopeInfo{
		ID:        scopeID,
		Type:      ScopeTypeClass,
		Parent:    sa.getCurrentScopeID(),
		StartPos:  node.GetPosition(),
		EndPos:    ast.PositionInfo{}, // TODO: Calculate end position
		ClassName: className,
	}
	sa.symbolTable.Scopes = append(sa.symbolTable.Scopes, classScope)

	// Add class symbol
	visibility := VisibilityPublic
	if strings.HasPrefix(className, "_") {
		visibility = VisibilityPrivate
	}

	classSymbol := EnhancedSymbol{
		Name:          className,
		Kind:          SymbolKindClass,
		Type:          className,
		Position:      node.GetPosition(),
		Range:         sa.positionToRange(node.GetPosition(), len(node.Name)),
		Scope:         sa.getCurrentScopeType(),
		ScopeID:       sa.getCurrentScopeID(),
		Visibility:    visibility,
		QualifiedName: sa.getQualifiedName(className),
	}
	sa.symbolTable.Symbols = append(sa.symbolTable.Symbols, classSymbol)

	// Save current context and analyze class members
	oldClass := sa.currentClass
	sa.currentClass = className

	for _, member := range node.Members {
		if member.IsVariable {
			sa.analyzeClassMemberVariable(member, scopeID)
		} else {
			sa.analyzeClassMemberFunction(member, scopeID)
		}
	}

	sa.currentClass = oldClass
}

// analyzeClassMemberVariable analyzes class member variables
func (sa *SemanticAnalyzer) analyzeClassMemberVariable(member *ast.ClassMemberNode, classScopeID string) {
	if member.Variable == nil {
		return
	}

	varName := strings.ToUpper(member.Variable.Name)
	visibility := VisibilityPrivate
	if member.IsPublic {
		visibility = VisibilityPublic
	}
	if member.IsShared {
		visibility = VisibilityShared
	}

	memberSymbol := EnhancedSymbol{
		Name:          varName,
		Kind:          SymbolKindVariable,
		Type:          strings.ToUpper(member.Variable.Type),
		Position:      member.Variable.GetPosition(),
		Range:         sa.positionToRange(member.Variable.GetPosition(), len(member.Variable.Name)),
		Scope:         ScopeTypeClass,
		ScopeID:       classScopeID,
		Visibility:    visibility,
		ParentClass:   sa.currentClass,
		QualifiedName: sa.getQualifiedName(varName),
		IsShared:      member.IsShared,
	}
	sa.symbolTable.Symbols = append(sa.symbolTable.Symbols, memberSymbol)
}

// analyzeClassMemberFunction analyzes class member functions
func (sa *SemanticAnalyzer) analyzeClassMemberFunction(member *ast.ClassMemberNode, classScopeID string) {
	if member.Function == nil {
		return
	}

	funcName := strings.ToUpper(member.Function.Name)
	visibility := VisibilityPrivate
	if member.IsPublic {
		visibility = VisibilityPublic
	}
	if member.IsShared {
		visibility = VisibilityShared
	}

	memberSymbol := EnhancedSymbol{
		Name:          funcName,
		Kind:          SymbolKindFunction,
		Type:          strings.ToUpper(member.Function.ReturnType),
		Position:      member.Function.GetPosition(),
		Range:         sa.positionToRange(member.Function.GetPosition(), len(member.Function.Name)),
		Scope:         ScopeTypeClass,
		ScopeID:       classScopeID,
		Visibility:    visibility,
		ParentClass:   sa.currentClass,
		QualifiedName: sa.getQualifiedName(funcName),
		IsShared:      member.IsShared,
	}
	sa.symbolTable.Symbols = append(sa.symbolTable.Symbols, memberSymbol)
}

// analyzeVariableDeclaration analyzes variable declarations
func (sa *SemanticAnalyzer) analyzeVariableDeclaration(node *ast.VariableDeclarationNode) {
	sa.analyzeVariableDeclarationWithScope(node, sa.getCurrentScopeType(), sa.getCurrentScopeID())
}

// analyzeVariableDeclarationWithScope analyzes variable declarations with explicit scope
func (sa *SemanticAnalyzer) analyzeVariableDeclarationWithScope(node *ast.VariableDeclarationNode, scopeType ScopeType, scopeID string) {
	varName := strings.ToUpper(node.Name)

	visibility := VisibilityPublic
	if strings.HasPrefix(varName, "_") {
		visibility = VisibilityPrivate
	}

	varSymbol := EnhancedSymbol{
		Name:          varName,
		Kind:          SymbolKindVariable,
		Type:          strings.ToUpper(node.Type),
		Position:      node.GetPosition(),
		Range:         sa.positionToRange(node.GetPosition(), len(node.Name)),
		Scope:         scopeType,
		ScopeID:       scopeID,
		Visibility:    visibility,
		ParentClass:   sa.currentClass,
		QualifiedName: sa.getQualifiedName(varName),
	}
	sa.symbolTable.Symbols = append(sa.symbolTable.Symbols, varSymbol)
}

// analyzeStatementBlock analyzes statements within a block (function body, etc.)
func (sa *SemanticAnalyzer) analyzeStatementBlock(block *ast.StatementBlockNode, parentScopeID string) {
	if block == nil || block.Statements == nil {
		return
	}

	// Create a block scope
	blockScopeID := fmt.Sprintf("block_%s_%p", parentScopeID, block)
	blockScope := ScopeInfo{
		ID:       blockScopeID,
		Type:     ScopeTypeBlock,
		Parent:   parentScopeID,
		StartPos: ast.PositionInfo{}, // TODO: Get actual position from block
		EndPos:   ast.PositionInfo{}, // TODO: Get actual end position
	}
	sa.symbolTable.Scopes = append(sa.symbolTable.Scopes, blockScope)

	// Analyze each statement in the block
	for _, stmt := range block.Statements {
		sa.analyzeStatement(stmt, blockScopeID)
	}
}

// analyzeStatement analyzes individual statements
func (sa *SemanticAnalyzer) analyzeStatement(stmt ast.Node, scopeID string) {
	if stmt == nil {
		return
	}

	switch node := stmt.(type) {
	case *ast.VariableDeclarationNode:
		// Analyze the variable with the current block scope
		sa.analyzeVariableDeclarationWithScope(node, ScopeTypeBlock, scopeID)

	case *ast.AssignmentNode:
		// Could analyze assignment targets here

	case *ast.IfStatementNode:
		// Analyze nested blocks
		if node.ThenBlock != nil {
			sa.analyzeStatementBlock(node.ThenBlock, scopeID)
		}
		if node.ElseBlock != nil {
			sa.analyzeStatementBlock(node.ElseBlock, scopeID)
		}

	case *ast.WhileStatementNode:
		// Analyze loop body
		if node.Body != nil {
			sa.analyzeStatementBlock(node.Body, scopeID)
		}

	case *ast.TryStatementNode:
		// Analyze try/catch/finally blocks
		if node.TryBody != nil {
			sa.analyzeStatementBlock(node.TryBody, scopeID)
		}
		if node.CatchBody != nil {
			sa.analyzeStatementBlock(node.CatchBody, scopeID)
		}
		if node.FinallyBody != nil {
			sa.analyzeStatementBlock(node.FinallyBody, scopeID)
		}

	case *ast.StatementBlockNode:
		// Nested statement block
		sa.analyzeStatementBlock(node, scopeID)

	case *ast.FunctionCallNode:
		// Track function call
		sa.analyzeFunctionCall(node)

	// Add more statement types as needed
	default:
		// For expressions and other nodes, recursively check for function calls
		sa.analyzeNodeForFunctionCalls(stmt)
	}
}

// Helper methods

func (sa *SemanticAnalyzer) getCurrentScopeID() string {
	if len(sa.symbolTable.Scopes) == 0 {
		return "global"
	}
	return sa.symbolTable.Scopes[len(sa.symbolTable.Scopes)-1].ID
}

func (sa *SemanticAnalyzer) getCurrentScopeType() ScopeType {
	if len(sa.symbolTable.Scopes) == 0 {
		return ScopeTypeGlobal
	}
	return sa.symbolTable.Scopes[len(sa.symbolTable.Scopes)-1].Type
}

func (sa *SemanticAnalyzer) getQualifiedName(name string) string {
	if sa.currentClass != "" {
		return fmt.Sprintf("%s.%s", sa.currentClass, name)
	}
	return name
}

func (sa *SemanticAnalyzer) positionToRange(pos ast.PositionInfo, length int) protocol.Range {
	return protocol.Range{
		Start: protocol.Position{
			Line:      uint32(pos.Line - 1),   // LSP is 0-based, parser is 1-based
			Character: uint32(pos.Column - 1), // LSP is 0-based, parser is 1-based
		},
		End: protocol.Position{
			Line:      uint32(pos.Line - 1),
			Character: uint32(pos.Column - 1 + length),
		},
	}
}

// Scoped symbol resolution methods

// ResolveSymbolAtPosition resolves a symbol at a given position with scope awareness
func (sa *SemanticAnalyzer) ResolveSymbolAtPosition(position protocol.Position) *EnhancedSymbol {
	// Find the scope at the given position
	scope := sa.findScopeAtPosition(position)
	if scope == nil {
		scope = &ScopeInfo{ID: "global", Type: ScopeTypeGlobal}
	}

	// Find symbol that contains the position
	for _, symbol := range sa.symbolTable.Symbols {
		if sa.positionInRange(position, symbol.Range) {
			// Check if symbol is accessible from current scope
			if sa.isSymbolAccessible(symbol, scope) {
				return &symbol
			}
		}
	}

	return nil
}

// FindSymbolsByName finds all symbols with the given name, respecting scope visibility
func (sa *SemanticAnalyzer) FindSymbolsByName(name string, position protocol.Position) []EnhancedSymbol {
	var results []EnhancedSymbol
	scope := sa.findScopeAtPosition(position)
	if scope == nil {
		scope = &ScopeInfo{ID: "global", Type: ScopeTypeGlobal}
	}

	for _, symbol := range sa.symbolTable.Symbols {
		if strings.EqualFold(symbol.Name, name) && sa.isSymbolAccessible(symbol, scope) {
			results = append(results, symbol)
		}
	}

	return results
}

// GetAccessibleSymbols returns all symbols accessible from a given position
func (sa *SemanticAnalyzer) GetAccessibleSymbols(position protocol.Position) []EnhancedSymbol {
	var results []EnhancedSymbol
	scope := sa.findScopeAtPosition(position)
	if scope == nil {
		scope = &ScopeInfo{ID: "global", Type: ScopeTypeGlobal}
	}

	for _, symbol := range sa.symbolTable.Symbols {
		if sa.isSymbolAccessible(symbol, scope) {
			results = append(results, symbol)
		}
	}

	return results
}

// GetClassMembersAccessible returns class members accessible from a given context
func (sa *SemanticAnalyzer) GetClassMembersAccessible(className string, position protocol.Position) []EnhancedSymbol {
	var results []EnhancedSymbol
	scope := sa.findScopeAtPosition(position)

	for _, symbol := range sa.symbolTable.Symbols {
		if symbol.ParentClass == className {
			if sa.isClassMemberAccessible(symbol, scope, className) {
				results = append(results, symbol)
			}
		}
	}

	return results
}

// findScopeAtPosition finds the most specific scope that contains the given position
func (sa *SemanticAnalyzer) findScopeAtPosition(position protocol.Position) *ScopeInfo {
	var mostSpecificScope *ScopeInfo

	for i := range sa.symbolTable.Scopes {
		scope := &sa.symbolTable.Scopes[i]
		if sa.positionInScope(position, scope) {
			// Find the most nested scope
			if mostSpecificScope == nil || sa.isScopeNestedIn(scope, mostSpecificScope) {
				mostSpecificScope = scope
			}
		}
	}

	return mostSpecificScope
}

// positionInScope checks if a position is within a scope
func (sa *SemanticAnalyzer) positionInScope(position protocol.Position, scope *ScopeInfo) bool {
	// For now, we'll use a simple line-based check
	// TODO: Implement proper position range checking when parser provides end positions
	startLine := uint32(scope.StartPos.Line - 1) // Convert to 0-based

	// For global scope, always return true
	if scope.Type == ScopeTypeGlobal {
		return true
	}

	// For other scopes, check if position is after the start
	return position.Line >= startLine
}

// isScopeNestedIn checks if scope1 is nested within scope2
func (sa *SemanticAnalyzer) isScopeNestedIn(scope1, scope2 *ScopeInfo) bool {
	// Check if scope1 has scope2 as an ancestor
	current := scope1
	for current != nil && current.Parent != "" {
		if current.Parent == scope2.ID {
			return true
		}
		// Find parent scope
		current = sa.findScopeByID(current.Parent)
	}
	return false
}

// findScopeByID finds a scope by its ID
func (sa *SemanticAnalyzer) findScopeByID(id string) *ScopeInfo {
	for i := range sa.symbolTable.Scopes {
		if sa.symbolTable.Scopes[i].ID == id {
			return &sa.symbolTable.Scopes[i]
		}
	}
	return nil
}

// isSymbolAccessible checks if a symbol is accessible from a given scope
func (sa *SemanticAnalyzer) isSymbolAccessible(symbol EnhancedSymbol, currentScope *ScopeInfo) bool {
	symbolScope := sa.findScopeByID(symbol.ScopeID)
	if symbolScope == nil {
		// If we can't find the symbol's scope, assume it's global and accessible
		return true
	}

	switch symbol.Visibility {
	case VisibilityPublic:
		return true // Public symbols are always accessible
	case VisibilityPrivate:
		// Private symbols are only accessible within the same scope or class
		if symbol.ParentClass != "" {
			return sa.isInSameClass(currentScope, symbol.ParentClass)
		}
		return currentScope.ID == symbolScope.ID || sa.isScopeNestedIn(currentScope, symbolScope)
	case VisibilityShared:
		// Shared class members are accessible through the class
		return symbol.ParentClass != ""
	default:
		return true
	}
}

// isClassMemberAccessible checks if a class member is accessible from the current context
func (sa *SemanticAnalyzer) isClassMemberAccessible(symbol EnhancedSymbol, currentScope *ScopeInfo, targetClass string) bool {
	switch symbol.Visibility {
	case VisibilityPublic:
		return true
	case VisibilityPrivate:
		// Private members only accessible from within the same class
		return sa.isInSameClass(currentScope, targetClass)
	case VisibilityShared:
		// Shared members accessible if we can access the class
		return true
	default:
		return true
	}
}

// isInSameClass checks if the current scope is within the specified class
func (sa *SemanticAnalyzer) isInSameClass(currentScope *ScopeInfo, className string) bool {
	if currentScope == nil {
		return false
	}

	// Check if current scope is a class scope with the same name
	if currentScope.Type == ScopeTypeClass && currentScope.ClassName == className {
		return true
	}

	// Check if current scope is nested within the class
	current := currentScope
	for current != nil && current.Parent != "" {
		parent := sa.findScopeByID(current.Parent)
		if parent != nil && parent.Type == ScopeTypeClass && parent.ClassName == className {
			return true
		}
		current = parent
	}

	return false
}

// positionInRange checks if a position is within a range
func (sa *SemanticAnalyzer) positionInRange(position protocol.Position, rang protocol.Range) bool {
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

// Context-sensitive completion and hover features

// GetCompletionItems provides context-aware completion items
func (sa *SemanticAnalyzer) GetCompletionItems(position protocol.Position) []protocol.CompletionItem {
	var completions []protocol.CompletionItem

	// Get accessible symbols from current position
	accessibleSymbols := sa.GetAccessibleSymbols(position)

	// Add language keywords
	keywords := []string{
		"HAI", "KTHXBAI", "TEH", "VARIABLE", "FUNCSHUN", "CLAS",
		"KITTEH", "OF", "ITZ", "WIT", "AN", "DIS", "IZ", "NOPE",
		"WHILE", "GIVEZ", "UP", "NEW", "AS", "I", "CAN", "HAS",
		"MAYB", "OOPS", "OOPSIE", "ALWAYZ", "INTEGR", "DUBBLE",
		"STRIN", "BOOL", "YEZ", "NO", "NOTHIN", "BUKKIT", "BASKIT",
		"EVRYONE", "MAHSELF", "MOAR", "LES", "TIEMZ", "DIVIDEZ",
		"BIGGR", "THAN", "SMALLR", "SAEM", "OR", "NOT", "DO",
	}

	for _, keyword := range keywords {
		completions = append(completions, protocol.CompletionItem{
			Label:      keyword,
			Kind:       &[]protocol.CompletionItemKind{protocol.CompletionItemKindKeyword}[0],
			SortText:   &[]string{"0_" + keyword}[0], // Keywords get priority
			InsertText: &keyword,
		})
	}

	// Add accessible symbols with context-aware information
	for _, symbol := range accessibleSymbols {
		kind := sa.symbolKindToCompletionKind(symbol.Kind)
		detail := sa.buildSymbolDetail(symbol)
		documentation := sa.buildSymbolDocumentation(symbol)

		completion := protocol.CompletionItem{
			Label:  symbol.Name,
			Kind:   &kind,
			Detail: &detail,
			Documentation: &protocol.MarkupContent{
				Kind:  protocol.MarkupKindMarkdown,
				Value: documentation,
			},
			SortText:   &[]string{sa.getSymbolSortText(symbol)}[0],
			InsertText: &symbol.Name,
		}

		// Add snippet for functions with parameters
		if symbol.Kind == SymbolKindFunction {
			snippet := sa.buildFunctionSnippet(symbol)
			if snippet != "" {
				completion.InsertText = &snippet
				insertFormat := protocol.InsertTextFormatSnippet
				completion.InsertTextFormat = &insertFormat
			}
		}

		completions = append(completions, completion)
	}

	return completions
}

// GetHoverInfo provides enhanced hover information with context
func (sa *SemanticAnalyzer) GetHoverInfo(position protocol.Position) *protocol.Hover {
	// First, check if this is a function call
	if functionCall := sa.findFunctionCallAtPosition(position); functionCall != nil {
		return sa.buildFunctionCallHover(functionCall)
	}

	// Otherwise, look for symbol definitions
	symbol := sa.ResolveSymbolAtPosition(position)
	if symbol == nil {
		return nil
	}

	// Build comprehensive hover content
	var contents []string

	// Symbol signature
	signature := sa.buildSymbolSignature(*symbol)
	contents = append(contents, fmt.Sprintf("```olol\n%s\n```", signature))

	// Visibility and scope information
	visibilityInfo := sa.buildVisibilityInfo(*symbol)
	if visibilityInfo != "" {
		contents = append(contents, visibilityInfo)
	}

	// Module source information
	if symbol.SourceModule != "" {
		contents = append(contents, fmt.Sprintf("**Module:** `%s`", symbol.SourceModule))
	}

	// Documentation
	if symbol.Documentation != "" {
		contents = append(contents, "---")
		contents = append(contents, symbol.Documentation)
	}

	return &protocol.Hover{
		Contents: protocol.MarkupContent{
			Kind:  protocol.MarkupKindMarkdown,
			Value: strings.Join(contents, "\n\n"),
		},
		Range: &symbol.Range,
	}
}

// GetDefinitionLocation provides enhanced go-to-definition
func (sa *SemanticAnalyzer) GetDefinitionLocation(position protocol.Position) *protocol.Location {
	symbol := sa.ResolveSymbolAtPosition(position)
	if symbol == nil {
		return nil
	}

	// For imported symbols, try to provide the location in the source module
	uri := sa.uri
	if symbol.SourceModule != "" && symbol.SourceModule != "stdlib" {
		// Convert module path to URI
		if !strings.HasPrefix(symbol.SourceModule, "file://") {
			uri = "file://" + symbol.SourceModule
		} else {
			uri = symbol.SourceModule
		}
	}

	return &protocol.Location{
		URI:   uri,
		Range: symbol.Range,
	}
}

// Helper methods for completion and hover

func (sa *SemanticAnalyzer) symbolKindToCompletionKind(kind SymbolKind) protocol.CompletionItemKind {
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

func (sa *SemanticAnalyzer) buildSymbolDetail(symbol EnhancedSymbol) string {
	var parts []string

	switch symbol.Kind {
	case SymbolKindFunction:
		parts = append(parts, fmt.Sprintf("(%s)", sa.symbolKindToString(symbol.Kind)))
		if symbol.Type != "" && symbol.Type != "NOTHIN" {
			parts = append(parts, fmt.Sprintf("→ %s", symbol.Type))
		}
	case SymbolKindClass:
		parts = append(parts, fmt.Sprintf("(%s)", sa.symbolKindToString(symbol.Kind)))
	case SymbolKindVariable:
		parts = append(parts, fmt.Sprintf("(%s) %s", sa.symbolKindToString(symbol.Kind), symbol.Type))
	default:
		if symbol.Type != "" {
			parts = append(parts, symbol.Type)
		}
	}

	return strings.Join(parts, " ")
}

func (sa *SemanticAnalyzer) buildSymbolDocumentation(symbol EnhancedSymbol) string {
	var docs []string

	// Basic symbol info
	docs = append(docs, fmt.Sprintf("**%s** (%s)", symbol.Name, sa.symbolKindToString(symbol.Kind)))

	// Qualified name if different from simple name
	if symbol.QualifiedName != symbol.Name {
		docs = append(docs, fmt.Sprintf("Qualified name: `%s`", symbol.QualifiedName))
	}

	// Scope information
	if symbol.ScopeID != "global" {
		docs = append(docs, fmt.Sprintf("Scope: `%s`", symbol.ScopeID))
	}

	return strings.Join(docs, "\n")
}

func (sa *SemanticAnalyzer) getSymbolSortText(symbol EnhancedSymbol) string {
	// Priority ordering: local variables > parameters > functions > classes > imports > stdlib
	priority := "5"

	switch symbol.Kind {
	case SymbolKindVariable:
		if symbol.Scope == ScopeTypeFunction || symbol.Scope == ScopeTypeBlock {
			priority = "1" // Local variables first
		} else {
			priority = "3" // Global variables
		}
	case SymbolKindParameter:
		priority = "2" // Parameters second
	case SymbolKindFunction:
		if symbol.SourceModule == "stdlib" {
			priority = "6" // Stdlib functions lower priority
		} else {
			priority = "3" // User functions
		}
	case SymbolKindClass:
		priority = "4"
	case SymbolKindImport:
		priority = "7"
	}

	return fmt.Sprintf("%s_%s", priority, symbol.Name)
}

func (sa *SemanticAnalyzer) buildFunctionSnippet(symbol EnhancedSymbol) string {
	// TODO: This would require parameter information from the symbol table
	// For now, just return the function name
	return symbol.Name
}

func (sa *SemanticAnalyzer) buildSymbolSignature(symbol EnhancedSymbol) string {
	switch symbol.Kind {
	case SymbolKindFunction:
		// TODO: Build proper function signature with parameters
		signature := fmt.Sprintf("FUNCSHUN %s", symbol.Name)
		if symbol.Type != "" && symbol.Type != "NOTHIN" {
			signature += fmt.Sprintf(" GIVEZ %s", symbol.Type)
		}
		return signature
	case SymbolKindClass:
		return fmt.Sprintf("CLAS %s", symbol.Name)
	case SymbolKindVariable:
		return fmt.Sprintf("VARIABLE %s TEH %s", symbol.Name, symbol.Type)
	default:
		return symbol.Name
	}
}

func (sa *SemanticAnalyzer) buildVisibilityInfo(symbol EnhancedSymbol) string {
	var parts []string

	switch symbol.Visibility {
	case VisibilityPublic:
		parts = append(parts, "**Visibility:** Public")
	case VisibilityPrivate:
		parts = append(parts, "**Visibility:** Private")
	case VisibilityShared:
		parts = append(parts, "**Visibility:** Shared")
	}

	if symbol.IsShared {
		parts = append(parts, "**Shared member**")
	}

	if symbol.ParentClass != "" {
		parts = append(parts, fmt.Sprintf("**Class:** `%s`", symbol.ParentClass))
	}

	return strings.Join(parts, " • ")
}

func (sa *SemanticAnalyzer) symbolKindToString(kind SymbolKind) string {
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

// Function call analysis methods

// analyzeFunctionCall analyzes a function call and tracks it
func (sa *SemanticAnalyzer) analyzeFunctionCall(node *ast.FunctionCallNode) {
	if node == nil {
		return
	}

	var functionName string
	var callType FunctionCallType
	var objectType string

	// Determine the type of function call and extract function name
	switch funcNode := node.Function.(type) {
	case *ast.IdentifierNode:
		// Global function call
		functionName = strings.ToUpper(funcNode.Name)
		callType = FunctionCallGlobal

	case *ast.MemberAccessNode:
		// Method call (obj DO method)
		functionName = strings.ToUpper(funcNode.Member)
		callType = FunctionCallMethod

		// Try to determine object type if possible
		// For now, we'll leave this as a placeholder
		objectType = "unknown"

	default:
		return // Unsupported function call type
	}

	// Create function call record
	functionCall := FunctionCall{
		CallSite:     node.GetPosition(),
		FunctionName: functionName,
		CallType:     callType,
		ObjectType:   objectType,
		Range:        sa.positionToRange(node.GetPosition(), len(functionName)),
	}

	// Collect argument positions
	for _, arg := range node.Arguments {
		if arg != nil {
			functionCall.Arguments = append(functionCall.Arguments, arg.GetPosition())
		}
	}

	// Try to resolve the function
	functionCall.ResolvedTo = sa.resolveFunctionSymbol(functionName, callType)

	// Add to function calls list
	sa.symbolTable.FunctionCalls = append(sa.symbolTable.FunctionCalls, functionCall)
}

// analyzeNodeForFunctionCalls recursively searches for function calls in expressions
func (sa *SemanticAnalyzer) analyzeNodeForFunctionCalls(node ast.Node) {
	if node == nil {
		return
	}

	switch n := node.(type) {
	case *ast.FunctionCallNode:
		// Direct function call
		sa.analyzeFunctionCall(n)

	case *ast.BinaryOpNode:
		// Check left and right operands
		sa.analyzeNodeForFunctionCalls(n.Left)
		sa.analyzeNodeForFunctionCalls(n.Right)

	case *ast.UnaryOpNode:
		// Check operand
		sa.analyzeNodeForFunctionCalls(n.Operand)

	case *ast.AssignmentNode:
		// Check value expression
		sa.analyzeNodeForFunctionCalls(n.Value)

	case *ast.IfStatementNode:
		// Check condition
		sa.analyzeNodeForFunctionCalls(n.Condition)

	case *ast.WhileStatementNode:
		// Check condition
		sa.analyzeNodeForFunctionCalls(n.Condition)

	case *ast.ReturnStatementNode:
		// Check return value
		sa.analyzeNodeForFunctionCalls(n.Value)

		// Add more node types as needed for comprehensive coverage
	}
}

// resolveFunctionSymbol attempts to resolve a function call to its symbol definition
func (sa *SemanticAnalyzer) resolveFunctionSymbol(functionName string, callType FunctionCallType) *EnhancedSymbol {
	// Search through all symbols for a matching function
	for i := range sa.symbolTable.Symbols {
		symbol := &sa.symbolTable.Symbols[i]

		if symbol.Kind == SymbolKindFunction &&
			strings.ToUpper(symbol.Name) == functionName {

			// For global calls, prefer functions in current module or global scope
			if callType == FunctionCallGlobal {
				if symbol.Scope == ScopeTypeGlobal || symbol.SourceModule == sa.uri {
					return symbol
				}
			}

			// For method calls, we'd need more sophisticated resolution
			// based on the object type, but for now return any match
			if callType == FunctionCallMethod {
				return symbol
			}
		}
	}

	// Check if it's a builtin function
	if sa.isBuiltinFunction(functionName) {
		// Create a temporary symbol for builtin functions
		return &EnhancedSymbol{
			Name:          functionName,
			Kind:          SymbolKindFunction,
			Type:          "builtin",
			Scope:         ScopeTypeGlobal,
			Visibility:    VisibilityPublic,
			Documentation: sa.getBuiltinDocumentation(functionName),
		}
	}

	return nil // Function not found
}

// isBuiltinFunction checks if a function name is a builtin
func (sa *SemanticAnalyzer) isBuiltinFunction(functionName string) bool {
	// Check standard library functions
	builtins := map[string]bool{
		"SAYZ":   true, // STDIO
		"SAY":    true,
		"GIMME":  true,
		"ABS":    true, // MATH
		"MAX":    true,
		"MIN":    true,
		"SQRT":   true,
		"POW":    true,
		"SIN":    true,
		"COS":    true,
		"RANDOM": true,
		"SLEEP":  true, // TIME
		"ASSERT": true, // TEST
		"LEN":    true, // STRING
		"CONCAT": true,
	}

	return builtins[functionName]
}

// getBuiltinDocumentation returns documentation for builtin functions
func (sa *SemanticAnalyzer) getBuiltinDocumentation(functionName string) string {
	docs := map[string]string{
		"SAYZ":   "Prints a value followed by a newline",
		"SAY":    "Prints a value without a newline",
		"GIMME":  "Reads input from the user",
		"ABS":    "Returns the absolute value of a number",
		"MAX":    "Returns the maximum of two numbers",
		"MIN":    "Returns the minimum of two numbers",
		"SQRT":   "Returns the square root of a number",
		"POW":    "Returns base raised to the power of exponent",
		"SIN":    "Returns the sine of an angle in radians",
		"COS":    "Returns the cosine of an angle in radians",
		"RANDOM": "Returns a random number between 0 and 1",
		"SLEEP":  "Pauses execution for the specified number of milliseconds",
		"ASSERT": "Asserts that a condition is true, throws exception if false",
		"LEN":    "Returns the length of a string",
		"CONCAT": "Concatenates two strings",
	}

	if doc, exists := docs[functionName]; exists {
		return doc
	}
	return "Built-in function"
}

// findFunctionCallAtPosition finds a function call at the given position
func (sa *SemanticAnalyzer) findFunctionCallAtPosition(position protocol.Position) *FunctionCall {
	for i := range sa.symbolTable.FunctionCalls {
		call := &sa.symbolTable.FunctionCalls[i]
		if sa.positionInRange(position, call.Range) {
			return call
		}
	}
	return nil
}

// buildFunctionCallHover builds hover information for a function call
func (sa *SemanticAnalyzer) buildFunctionCallHover(call *FunctionCall) *protocol.Hover {
	var contents []string

	// Function call header
	callTypeStr := sa.functionCallTypeToString(call.CallType)
	contents = append(contents, fmt.Sprintf("**%s Call**", callTypeStr))

	// Function signature
	resolvedSymbol := call.ResolvedTo
	if resolvedSymbol == nil {
		// Try a more comprehensive lookup in outer scopes
		resolvedSymbol = sa.lookupFunctionInOuterScopes(call.FunctionName, call.CallType)
	}

	if resolvedSymbol != nil {
		signature := sa.buildSymbolSignature(*resolvedSymbol)
		contents = append(contents, fmt.Sprintf("```olol\n%s\n```", signature))
	} else {
		// Still unresolved function call
		contents = append(contents, fmt.Sprintf("```olol\n%s(?)\n```", call.FunctionName))
		contents = append(contents, "*⚠️ Function not found*")
	}

	// Call details
	if call.CallType == FunctionCallMethod && call.ObjectType != "unknown" {
		contents = append(contents, fmt.Sprintf("**Object Type:** `%s`", call.ObjectType))
	}

	if len(call.Arguments) > 0 {
		contents = append(contents, fmt.Sprintf("**Arguments:** %d", len(call.Arguments)))
	}

	// Function documentation (if resolved)
	if resolvedSymbol != nil && resolvedSymbol.Documentation != "" {
		contents = append(contents, "---")
		contents = append(contents, resolvedSymbol.Documentation)
	}

	return &protocol.Hover{
		Contents: protocol.MarkupContent{
			Kind:  protocol.MarkupKindMarkdown,
			Value: strings.Join(contents, "\n\n"),
		},
		Range: &call.Range,
	}
}

// functionCallTypeToString converts function call type to string
func (sa *SemanticAnalyzer) functionCallTypeToString(callType FunctionCallType) string {
	switch callType {
	case FunctionCallGlobal:
		return "Global Function"
	case FunctionCallMethod:
		return "Method"
	case FunctionCallBuiltin:
		return "Built-in Function"
	default:
		return "Function"
	}
}

// lookupFunctionInOuterScopes performs a comprehensive function lookup in outer scopes
func (sa *SemanticAnalyzer) lookupFunctionInOuterScopes(functionName string, callType FunctionCallType) *EnhancedSymbol {
	// First, try case-insensitive matching for all function symbols
	for i := range sa.symbolTable.Symbols {
		symbol := &sa.symbolTable.Symbols[i]

		if symbol.Kind == SymbolKindFunction &&
			strings.EqualFold(symbol.Name, functionName) {
			return symbol
		}
	}

	// Try looking in imported modules
	for _, moduleImport := range sa.symbolTable.ImportedModules {
		if moduleSymbol := sa.lookupInImportedModule(functionName, moduleImport); moduleSymbol != nil {
			return moduleSymbol
		}
	}

	// Try looking in the runtime environment if available
	if sa.environment != nil {
		if environmentSymbol := sa.lookupInEnvironment(functionName, callType); environmentSymbol != nil {
			return environmentSymbol
		}
	}

	// Check if it's a builtin that wasn't caught earlier (with case variations)
	upperFunctionName := strings.ToUpper(functionName)
	if sa.isBuiltinFunction(upperFunctionName) {
		return &EnhancedSymbol{
			Name:          upperFunctionName,
			Kind:          SymbolKindFunction,
			Type:          "builtin",
			Scope:         ScopeTypeGlobal,
			Visibility:    VisibilityPublic,
			Documentation: sa.getBuiltinDocumentation(upperFunctionName),
		}
	}

	return nil // Still not found
}

// lookupInImportedModule looks up a function in imported modules
func (sa *SemanticAnalyzer) lookupInImportedModule(functionName string, moduleImport ModuleImport) *EnhancedSymbol {
	// Look for functions from this specific module
	for i := range sa.symbolTable.Symbols {
		symbol := &sa.symbolTable.Symbols[i]

		if symbol.Kind == SymbolKindFunction &&
			strings.EqualFold(symbol.Name, functionName) &&
			symbol.SourceModule == moduleImport.FilePath {
			return symbol
		}
	}
	return nil
}

// lookupInEnvironment looks up a function in the runtime environment
func (sa *SemanticAnalyzer) lookupInEnvironment(functionName string, _ FunctionCallType) *EnhancedSymbol {
	if sa.environment == nil {
		return nil
	}

	// Try to get the function from the environment
	upperFunctionName := strings.ToUpper(functionName)
	if function, err := sa.environment.GetFunction(upperFunctionName); err == nil {
		// Create a symbol from the environment function
		return &EnhancedSymbol{
			Name:          upperFunctionName,
			Kind:          SymbolKindFunction,
			Type:          "function", // Could inspect function.Type for more details
			Scope:         ScopeTypeGlobal,
			Visibility:    VisibilityPublic,
			Documentation: fmt.Sprintf("Function from environment: %s", function.Name),
		}
	}

	return nil
}
