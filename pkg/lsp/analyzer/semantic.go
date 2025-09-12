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
	scopeStack       []string // Track current scope IDs for proper nesting
	currentClass     string
	currentFile      string
	analysisErrors   []AnalysisError

	// Position tracking for enhanced IDE features
	positionToSymbol map[PositionKey]*EnhancedSymbol // Fast position-based lookup
}

// EnhancedSymbolTable represents symbols with enhanced metadata
type EnhancedSymbolTable struct {
	URI              string
	Symbols          []EnhancedSymbol
	Scopes           []ScopeInfo
	ImportedModules  []ModuleImport
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

// PositionKey represents a unique key for position-based lookups
type PositionKey struct {
	Line   int
	Column int
}

func joinDocs(docs []string) string {
	return strings.Join(docs, "\n")
}

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
		scopeStack:        []string{},
		analysisErrors:    []AnalysisError{},
		positionToSymbol:  make(map[PositionKey]*EnhancedSymbol),
	}
}

// Scope management methods for proper environment tracking

// pushScope creates a new scope and pushes it onto the scope stack
func (sa *SemanticAnalyzer) pushScope(scopeInfo ScopeInfo) {
	// Add scope to symbol table
	sa.symbolTable.Scopes = append(sa.symbolTable.Scopes, scopeInfo)

	// Push scope ID onto stack
	sa.scopeStack = append(sa.scopeStack, scopeInfo.ID)

	// Create new environment and push onto environment stack
	newEnv := environment.NewEnvironment(sa.environment)
	sa.environmentStack = append(sa.environmentStack, sa.environment)
	sa.environment = newEnv
}

// popScope removes the current scope from the scope stack
func (sa *SemanticAnalyzer) popScope() {
	if len(sa.scopeStack) > 0 {
		// Pop scope ID from stack
		sa.scopeStack = sa.scopeStack[:len(sa.scopeStack)-1]
	}

	// Restore previous environment
	if len(sa.environmentStack) > 0 {
		sa.environment = sa.environmentStack[len(sa.environmentStack)-1]
		sa.environmentStack = sa.environmentStack[:len(sa.environmentStack)-1]
	}
}

// getCurrentScopeID returns the current scope ID with proper nesting
func (sa *SemanticAnalyzer) getCurrentScopeID() string {
	if len(sa.scopeStack) == 0 {
		return "global"
	}
	return sa.scopeStack[len(sa.scopeStack)-1]
}

// getCurrentScopeType returns the current scope type
func (sa *SemanticAnalyzer) getCurrentScopeType() ScopeType {
	currentScopeID := sa.getCurrentScopeID()
	for _, scope := range sa.symbolTable.Scopes {
		if scope.ID == currentScopeID {
			return scope.Type
		}
	}
	return ScopeTypeGlobal
}

// addSymbolWithPosition adds a symbol and tracks its position for IDE features
func (sa *SemanticAnalyzer) addSymbolWithPosition(symbol EnhancedSymbol) {
	// Add to symbol table
	sa.symbolTable.Symbols = append(sa.symbolTable.Symbols, symbol)

	// Track position for fast lookup
	if symbol.Position.Line > 0 && symbol.Position.Column > 0 {
		key := PositionKey{Line: symbol.Position.Line, Column: symbol.Position.Column}
		// Get reference to the symbol we just added
		lastIndex := len(sa.symbolTable.Symbols) - 1
		sa.positionToSymbol[key] = &sa.symbolTable.Symbols[lastIndex]
	}
}

// Type inference methods

// inferExpressionType attempts to infer the type of an expression
func (sa *SemanticAnalyzer) inferExpressionType(expr ast.Node) string {
	if expr == nil {
		return "NOTHIN"
	}

	switch node := expr.(type) {
	case *ast.LiteralNode:
		// Literal values have known types
		return node.Value.Type()

	case *ast.IdentifierNode:
		// Look up the identifier's declared type
		name := strings.ToUpper(node.Name)
		currentScope := sa.findScopeByID(sa.getCurrentScopeID())
		symbol := sa.findSymbolByNameInScope(name, currentScope)
		if symbol != nil {
			return symbol.Type
		}
		return "unknown"

	case *ast.BinaryOpNode:
		// Infer type based on operation
		return sa.inferBinaryOpType(node)

	case *ast.UnaryOpNode:
		// Unary operations usually preserve operand type or return BOOL
		switch strings.ToUpper(node.Operator) {
		case "NOT":
			return "BOOL"
		default:
			return sa.inferExpressionType(node.Operand)
		}

	case *ast.CastNode:
		// Cast operations explicitly specify the target type
		return strings.ToUpper(node.TargetType)

	case *ast.FunctionCallNode:
		// Infer type from function return type
		return sa.inferFunctionCallType(node)

	case *ast.MemberAccessNode:
		// Infer type from member variable type
		return sa.inferMemberAccessType(node)

	case *ast.ObjectInstantiationNode:
		// Object instantiation returns the class type
		return strings.ToUpper(node.Class.Name)

	default:
		return "unknown"
	}
}

// inferBinaryOpType infers the result type of binary operations
func (sa *SemanticAnalyzer) inferBinaryOpType(node *ast.BinaryOpNode) string {
	leftType := sa.inferExpressionType(node.Left)
	rightType := sa.inferExpressionType(node.Right)

	switch strings.ToUpper(node.Operator) {
	case "MOAR", "LES", "TIEMZ", "DIVIDEZ":
		// Arithmetic operations: promote to higher precision type
		if leftType == "DUBBLE" || rightType == "DUBBLE" {
			return "DUBBLE"
		}
		if leftType == "INTEGR" || rightType == "INTEGR" {
			return "INTEGR"
		}
		return "INTEGR" // Default for arithmetic

	case "BIGGR THAN", "SMALLR THAN", "SAEM AS", "AN", "OR", "NOT":
		// Comparison and logical operations return BOOL
		return "BOOL"

	default:
		// For unknown operations, return the left operand type
		return leftType
	}
}

// inferFunctionCallType infers the return type of a function call
func (sa *SemanticAnalyzer) inferFunctionCallType(node *ast.FunctionCallNode) string {
	switch funcNode := node.Function.(type) {
	case *ast.IdentifierNode:
		// Global function call
		functionName := strings.ToUpper(funcNode.Name)
		currentScope := sa.findScopeByID(sa.getCurrentScopeID())
		symbol := sa.findSymbolByNameInScope(functionName, currentScope)
		if symbol != nil && symbol.Kind == SymbolKindFunction {
			return symbol.Type
		}

	case *ast.MemberAccessNode:
		// Method call - look up in class
		objectType := sa.inferExpressionType(funcNode.Object)
		methodName := strings.ToUpper(funcNode.Member.Name)

		// Find the class and method
		for _, symbol := range sa.symbolTable.Symbols {
			if symbol.Kind == SymbolKindFunction &&
				symbol.ParentClass == objectType &&
				strings.EqualFold(symbol.Name, methodName) {
				return symbol.Type
			}
		}
	}

	return "unknown"
}

// inferMemberAccessType infers the type of member access
func (sa *SemanticAnalyzer) inferMemberAccessType(node *ast.MemberAccessNode) string {
	objectType := sa.inferExpressionType(node.Object)
	memberName := strings.ToUpper(node.Member.Name)

	// Find the member variable in the class
	for _, symbol := range sa.symbolTable.Symbols {
		if symbol.Kind == SymbolKindVariable &&
			symbol.ParentClass == objectType &&
			strings.EqualFold(symbol.Name, memberName) {
			return symbol.Type
		}
	}

	return "unknown"
}

// AnalyzeDocument performs comprehensive semantic analysis
func (sa *SemanticAnalyzer) AnalyzeDocument(ctx context.Context) error {
	// Clear previous analysis results
	sa.symbolTable.Symbols = []EnhancedSymbol{}
	sa.symbolTable.Scopes = []ScopeInfo{}
	sa.symbolTable.ImportedModules = []ModuleImport{}
	sa.analysisErrors = []AnalysisError{}
	sa.positionToSymbol = make(map[PositionKey]*EnhancedSymbol)
	sa.scopeStack = []string{}
	sa.environmentStack = []*environment.Environment{}

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

	// Initialize global scope using new scope management
	globalScope := ScopeInfo{
		ID:       "global",
		Type:     ScopeTypeGlobal,
		Parent:   "",
		StartPos: ast.PositionInfo{Line: 1, Column: 1},
		EndPos:   ast.PositionInfo{Line: 999999, Column: 999999}, // Will be updated
	}
	sa.pushScope(globalScope)

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
			sa.addSymbolWithPosition(symbol)
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
			Documentation: joinDocs(function.Documentation),
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
			Documentation: joinDocs(class.Documentation),
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

	// Add function symbol to current scope first
	visibility := VisibilityPublic
	if strings.HasPrefix(functionName, "_") {
		visibility = VisibilityPrivate
	}

	funcSymbol := EnhancedSymbol{
		Documentation: joinDocs(node.Documentation),
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
	sa.addSymbolWithPosition(funcSymbol)

	// Create function scope for analyzing function body
	scopeID := fmt.Sprintf("func_%s_%d_%d", functionName, node.GetPosition().Line, node.GetPosition().Column)
	funcScope := ScopeInfo{
		ID:           scopeID,
		Type:         ScopeTypeFunction,
		Parent:       sa.getCurrentScopeID(),
		StartPos:     node.GetPosition(),
		EndPos:       ast.PositionInfo{}, // TODO: Calculate end position
		FunctionName: functionName,
	}
	sa.pushScope(funcScope)

	// Analyze function parameters in the new function scope
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
		sa.addSymbolWithPosition(paramSymbol)
	}

	// Analyze function body for local variables and nested symbols
	if node.Body != nil {
		sa.analyzeStatementBlock(node.Body)
	}

	// Pop function scope when done
	sa.popScope()
}

// analyzeClassDeclaration analyzes class declarations
func (sa *SemanticAnalyzer) analyzeClassDeclaration(node *ast.ClassDeclarationNode) {
	className := strings.ToUpper(node.Name)

	// Add class symbol to current scope first
	visibility := VisibilityPublic
	if strings.HasPrefix(className, "_") {
		visibility = VisibilityPrivate
	}

	classSymbol := EnhancedSymbol{
		Documentation: joinDocs(node.Documentation),
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
	sa.addSymbolWithPosition(classSymbol)

	// Create class scope for analyzing class members
	scopeID := fmt.Sprintf("class_%s", className)
	classScope := ScopeInfo{
		ID:        scopeID,
		Type:      ScopeTypeClass,
		Parent:    sa.getCurrentScopeID(),
		StartPos:  node.GetPosition(),
		EndPos:    ast.PositionInfo{}, // TODO: Calculate end position
		ClassName: className,
	}
	sa.pushScope(classScope)

	// Save current context and analyze class members
	oldClass := sa.currentClass
	sa.currentClass = className

	for _, member := range node.Members {
		if member.IsVariable {
			sa.analyzeClassMemberVariable(member)
		} else {
			sa.analyzeClassMemberFunction(member)
		}
	}

	sa.currentClass = oldClass
	sa.popScope()
}

// analyzeClassMemberVariable analyzes class member variables
func (sa *SemanticAnalyzer) analyzeClassMemberVariable(member *ast.ClassMemberNode) {
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

	// Extract type name from IdentifierNode, or use empty string if nil
	memberVarType := ""
	if member.Variable.Type != nil {
		memberVarType = strings.ToUpper(member.Variable.Type.Name)
	}

	memberSymbol := EnhancedSymbol{
		Documentation: joinDocs(member.Variable.Documentation),
		Name:          varName,
		Kind:          SymbolKindVariable,
		Type:          memberVarType,
		Position:      member.Variable.GetPosition(),
		Range:         sa.positionToRange(member.Variable.GetPosition(), len(member.Variable.Name)),
		Scope:         ScopeTypeClass,
		ScopeID:       sa.getCurrentScopeID(),
		Visibility:    visibility,
		ParentClass:   sa.currentClass,
		QualifiedName: sa.getQualifiedName(varName),
		IsShared:      member.IsShared,
	}
	sa.addSymbolWithPosition(memberSymbol)

	// If there's an explicit type, track it as an identifier reference for hover support
	if member.Variable.Type != nil {
		sa.trackIdentifierReference(member.Variable.Type)
	}

	// Analyze member variable initialization if present
	if member.Variable.Value != nil {
		sa.analyzeExpression(member.Variable.Value)
	}
}

// analyzeClassMemberFunction analyzes class member functions
func (sa *SemanticAnalyzer) analyzeClassMemberFunction(member *ast.ClassMemberNode) {
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
		Documentation: joinDocs(member.Function.Documentation),
		Name:          funcName,
		Kind:          SymbolKindFunction,
		Type:          strings.ToUpper(member.Function.ReturnType),
		Position:      member.Function.GetPosition(),
		Range:         sa.positionToRange(member.Function.GetPosition(), len(member.Function.Name)),
		Scope:         ScopeTypeClass,
		ScopeID:       sa.getCurrentScopeID(),
		Visibility:    visibility,
		ParentClass:   sa.currentClass,
		QualifiedName: sa.getQualifiedName(funcName),
		IsShared:      member.IsShared,
	}
	sa.addSymbolWithPosition(memberSymbol)

	// Analyze the member function body like a regular function
	sa.analyzeFunctionDeclaration(member.Function)
}

// analyzeVariableDeclaration analyzes variable declarations
func (sa *SemanticAnalyzer) analyzeVariableDeclaration(node *ast.VariableDeclarationNode) {
	varName := strings.ToUpper(node.Name)

	visibility := VisibilityPublic
	if strings.HasPrefix(varName, "_") {
		visibility = VisibilityPrivate
	}

	// Determine variable type - use declared type or infer from initialization
	varType := ""
	if node.Type != nil {
		varType = strings.ToUpper(node.Type.Name)
	}
	if node.Value != nil {
		// Analyze initialization expression first
		sa.analyzeExpression(node.Value)

		// If no explicit type was declared, infer it from the initialization
		if varType == "" {
			inferredType := sa.inferExpressionType(node.Value)
			if inferredType != "unknown" {
				varType = inferredType
			} else {
				varType = "NOTHIN" // Default for unresolved types
			}
		}
	} else if varType == "" {
		// No initialization and no explicit type - default to NOTHIN
		varType = "NOTHIN"
	}

	varSymbol := EnhancedSymbol{
		Documentation: joinDocs(node.Documentation),
		Name:          varName,
		Kind:          SymbolKindVariable,
		Type:          varType,
		Position:      node.GetPosition(),
		Range:         sa.positionToRange(node.GetPosition(), len(node.Name)),
		Scope:         sa.getCurrentScopeType(),
		ScopeID:       sa.getCurrentScopeID(),
		Visibility:    visibility,
		ParentClass:   sa.currentClass,
		QualifiedName: sa.getQualifiedName(varName),
	}
	sa.addSymbolWithPosition(varSymbol)

	// If there's an explicit type, track it as an identifier reference for hover support
	if node.Type != nil {
		sa.trackIdentifierReference(node.Type)
	}
}

// analyzeStatementBlock analyzes statements within a block (function body, etc.)
func (sa *SemanticAnalyzer) analyzeStatementBlock(block *ast.StatementBlockNode) {
	if block == nil || block.Statements == nil {
		return
	}

	// For function and class bodies, we don't need an additional block scope
	// since they already have their own scopes. For control structure blocks,
	// we do need a block scope.
	currentScopeType := sa.getCurrentScopeType()
	needsBlockScope := currentScopeType == ScopeTypeBlock || currentScopeType == ScopeTypeGlobal

	if needsBlockScope {
		// Create a block scope
		blockScopeID := fmt.Sprintf("block_%s_%p", sa.getCurrentScopeID(), block)
		blockScope := ScopeInfo{
			ID:       blockScopeID,
			Type:     ScopeTypeBlock,
			Parent:   sa.getCurrentScopeID(),
			StartPos: ast.PositionInfo{}, // TODO: Get actual position from block
			EndPos:   ast.PositionInfo{}, // TODO: Get actual end position
		}
		sa.pushScope(blockScope)
	}

	// Analyze each statement in the block
	for _, stmt := range block.Statements {
		sa.analyzeStatement(stmt)
	}

	if needsBlockScope {
		sa.popScope()
	}
}

// analyzeStatement analyzes individual statements
func (sa *SemanticAnalyzer) analyzeStatement(stmt ast.Node) {
	if stmt == nil {
		return
	}

	switch node := stmt.(type) {
	case *ast.VariableDeclarationNode:
		// Analyze the variable declaration
		sa.analyzeVariableDeclaration(node)

	case *ast.AssignmentNode:
		// Analyze assignment target and value expressions
		sa.analyzeExpression(node.Target)
		sa.analyzeExpression(node.Value)

	case *ast.IfStatementNode:
		// Analyze condition expression
		sa.analyzeExpression(node.Condition)
		// Analyze nested blocks with proper scope management
		sa.analyzeIfStatement(node)

	case *ast.WhileStatementNode:
		// Analyze while statement with proper scope management
		sa.analyzeWhileStatement(node)

	case *ast.TryStatementNode:
		// Analyze try statement with proper scope management
		sa.analyzeTryStatement(node)

	case *ast.StatementBlockNode:
		// Nested statement block
		sa.analyzeStatementBlock(node)

	case *ast.FunctionCallNode:
		// Track function call
		sa.analyzeFunctionCall(node)

	case *ast.ReturnStatementNode:
		// Analyze return value expression
		if node.Value != nil {
			sa.analyzeExpression(node.Value)
		}

	case *ast.ThrowStatementNode:
		// Analyze throw expression
		if node.Expression != nil {
			sa.analyzeExpression(node.Expression)
		}

	// Add more statement types as needed
	default:
		// For expressions and other nodes, recursively analyze expressions
		sa.analyzeExpression(stmt)
	}
}

// analyzeIfStatement analyzes if statements with proper scope management
func (sa *SemanticAnalyzer) analyzeIfStatement(node *ast.IfStatementNode) {
	// Analyze condition expression
	sa.analyzeExpression(node.Condition)

	// Analyze then block with its own scope
	if node.ThenBlock != nil {
		sa.analyzeStatementBlock(node.ThenBlock)
	}

	// Analyze else-if branches
	for _, elseIfBranch := range node.ElseIfBranches {
		sa.analyzeExpression(elseIfBranch.Condition)
		if elseIfBranch.Block != nil {
			sa.analyzeStatementBlock(elseIfBranch.Block)
		}
	}

	// Analyze else block
	if node.ElseBlock != nil {
		sa.analyzeStatementBlock(node.ElseBlock)
	}
}

// analyzeWhileStatement analyzes while statements with proper scope management
func (sa *SemanticAnalyzer) analyzeWhileStatement(node *ast.WhileStatementNode) {
	// Analyze condition expression
	sa.analyzeExpression(node.Condition)

	// Analyze loop body with its own scope
	if node.Body != nil {
		sa.analyzeStatementBlock(node.Body)
	}
}

// analyzeTryStatement analyzes try-catch-finally statements with proper scope management
func (sa *SemanticAnalyzer) analyzeTryStatement(node *ast.TryStatementNode) {
	// Analyze try block with its own scope
	if node.TryBody != nil {
		sa.analyzeStatementBlock(node.TryBody)
	}

	// Analyze catch block with its own scope and exception variable
	if node.CatchBody != nil {
		// Create scope for catch block
		catchScopeID := fmt.Sprintf("catch_%p", node)
		catchScope := ScopeInfo{
			ID:       catchScopeID,
			Type:     ScopeTypeBlock,
			Parent:   sa.getCurrentScopeID(),
			StartPos: ast.PositionInfo{}, // TODO: Get actual position
			EndPos:   ast.PositionInfo{}, // TODO: Get actual end position
		}
		sa.pushScope(catchScope)

		// Add catch variable if present
		if node.CatchVar != "" {
			catchVarSymbol := EnhancedSymbol{
				Name:          strings.ToUpper(node.CatchVar),
				Kind:          SymbolKindVariable,
				Type:          "STRIN",            // Exception message is always a string
				Position:      ast.PositionInfo{}, // TODO: Get position from parser
				Range:         protocol.Range{},
				Scope:         ScopeTypeBlock,
				ScopeID:       catchScopeID,
				Visibility:    VisibilityPrivate,
				QualifiedName: strings.ToUpper(node.CatchVar),
			}
			sa.addSymbolWithPosition(catchVarSymbol)
		}

		sa.analyzeStatementBlock(node.CatchBody)
		sa.popScope()
	}

	// Analyze finally block with its own scope
	if node.FinallyBody != nil {
		sa.analyzeStatementBlock(node.FinallyBody)
	}
}

// analyzeExpression recursively analyzes expressions to track identifier references
func (sa *SemanticAnalyzer) analyzeExpression(expr ast.Node) {
	if expr == nil {
		return
	}

	switch node := expr.(type) {
	case *ast.IdentifierNode:
		// This is a variable reference - track it
		sa.trackIdentifierReference(node)

	case *ast.AssignmentNode:
		// Analyze assignment target and value
		sa.analyzeExpression(node.Target)
		sa.analyzeExpression(node.Value)

	case *ast.BinaryOpNode:
		// Analyze left and right operands
		sa.analyzeExpression(node.Left)
		sa.analyzeExpression(node.Right)

	case *ast.UnaryOpNode:
		// Analyze operand
		sa.analyzeExpression(node.Operand)

	case *ast.CastNode:
		// Analyze cast expression
		sa.analyzeExpression(node.Expression)

	case *ast.FunctionCallNode:
		// Track function call and analyze arguments
		sa.analyzeFunctionCall(node)
		for _, arg := range node.Arguments {
			sa.analyzeExpression(arg)
		}

	case *ast.MemberAccessNode:
		// Use type-aware member reference tracking
		sa.trackMemberReference(node)

	case *ast.ObjectInstantiationNode:
		// Analyze constructor arguments
		for _, arg := range node.ConstructorArgs {
			sa.analyzeExpression(arg)
		}

		sa.trackIdentifierReference(node.Class)

	case *ast.LiteralNode:
		// Literals don't need analysis
		return

	case *ast.StatementBlockNode:
		// Handle statement blocks that appear in expressions
		sa.analyzeStatementBlock(node)

	// Add more expression types as needed
	default:
		// For other node types, no specific analysis needed
	}
}

// trackIdentifierReference tracks an identifier reference and adds it to symbol references
func (sa *SemanticAnalyzer) trackIdentifierReference(node *ast.IdentifierNode) {
	if node == nil {
		return
	}

	identifierName := strings.ToUpper(node.Name)
	position := node.GetPosition()

	// Find the symbol this identifier refers to
	currentScopeID := sa.getCurrentScopeID()
	scope := sa.findScopeByID(currentScopeID)
	symbol := sa.findSymbolByNameInScope(identifierName, scope)

	// Always create a reference symbol for each usage location
	// This ensures each position has the correct range for hover functionality
	var refSymbol EnhancedSymbol
	if symbol != nil {
		// Add this position as a reference to the existing symbol
		symbol.References = append(symbol.References, position)

		// Create a reference symbol with the resolved symbol's type and properties
		refSymbol = EnhancedSymbol{
			Name:          identifierName,
			Kind:          symbol.Kind,
			Type:          symbol.Type,
			Position:      position,
			Range:         sa.positionToRange(position, len(node.Name)),
			Scope:         sa.getCurrentScopeType(),
			ScopeID:       currentScopeID,
			Visibility:    symbol.Visibility,
			ParentClass:   symbol.ParentClass,
			QualifiedName: symbol.QualifiedName,
			References:    []ast.PositionInfo{position},
		}
	} else {
		// Create a "reference-only" symbol for unresolved identifiers
		// This allows hover to work even if we can't find the declaration
		refSymbol = EnhancedSymbol{
			Name:          identifierName,
			Kind:          SymbolKindVariable, // Assume variable for now
			Type:          "unknown",
			Position:      position,
			Range:         sa.positionToRange(position, len(node.Name)),
			Scope:         sa.getCurrentScopeType(),
			ScopeID:       currentScopeID,
			Visibility:    VisibilityPublic,
			QualifiedName: identifierName,
			References:    []ast.PositionInfo{position},
		}
	}
	sa.addSymbolWithPosition(refSymbol)
}

// trackMemberReference tracks a member access reference with type-aware resolution
func (sa *SemanticAnalyzer) trackMemberReference(node *ast.MemberAccessNode) {
	if node == nil || node.Member == nil {
		return
	}

	// First analyze the object to understand its type
	sa.analyzeExpression(node.Object)
	objectType := sa.inferExpressionType(node.Object)

	memberName := strings.ToUpper(node.Member.Name)
	position := node.Member.GetPosition()

	// Look for the member in the object's class scope
	classScope := sa.findClassScope(objectType)
	var symbol *EnhancedSymbol
	if classScope != nil {
		symbol = sa.findSymbolByNameInScope(memberName, classScope)
	}

	// Create a reference symbol for the member
	var refSymbol EnhancedSymbol
	if symbol != nil {
		// Add this position as a reference to the existing symbol
		symbol.References = append(symbol.References, position)

		// Create a reference symbol with the resolved symbol's type and properties
		refSymbol = EnhancedSymbol{
			Name:          memberName,
			Kind:          symbol.Kind,
			Type:          symbol.Type,
			Position:      position,
			Range:         sa.positionToRange(position, len(node.Member.Name)),
			Scope:         sa.getCurrentScopeType(),
			ScopeID:       sa.getCurrentScopeID(),
			Visibility:    symbol.Visibility,
			ParentClass:   symbol.ParentClass,
			QualifiedName: symbol.QualifiedName,
			References:    []ast.PositionInfo{position},
		}
	} else {
		// Create a "reference-only" symbol for unresolved members
		refSymbol = EnhancedSymbol{
			Name:          memberName,
			Kind:          SymbolKindVariable, // Default assumption
			Type:          "unknown",
			Position:      position,
			Range:         sa.positionToRange(position, len(node.Member.Name)),
			Scope:         sa.getCurrentScopeType(),
			ScopeID:       sa.getCurrentScopeID(),
			Visibility:    VisibilityPublic,
			ParentClass:   objectType,
			QualifiedName: fmt.Sprintf("%s.%s", objectType, memberName),
			References:    []ast.PositionInfo{position},
		}
	}
	sa.addSymbolWithPosition(refSymbol)
}

// findClassScope finds the scope for a given class name
func (sa *SemanticAnalyzer) findClassScope(className string) *ScopeInfo {
	if className == "" || className == "unknown" {
		return nil
	}

	// Look for a class scope with the given class name
	for i := range sa.symbolTable.Scopes {
		scope := &sa.symbolTable.Scopes[i]
		if scope.Type == ScopeTypeClass && strings.EqualFold(scope.ClassName, className) {
			return scope
		}
	}
	return nil
}

// findSymbolByNameInScope finds a symbol by name considering scope hierarchy
// This follows the same lookup rules as the interpreter
func (sa *SemanticAnalyzer) findSymbolByNameInScope(name string, currentScope *ScopeInfo) *EnhancedSymbol {
	if currentScope == nil {
		currentScope = &ScopeInfo{ID: "global", Type: ScopeTypeGlobal}
	}

	// Create a list of scopes to search, ordered from most specific to most general
	scopesToSearch := sa.buildScopeHierarchy(currentScope)

	// Search in each scope, starting from the most specific
	for _, scopeID := range scopesToSearch {
		// Find symbols in this specific scope
		for i := range sa.symbolTable.Symbols {
			symbol := &sa.symbolTable.Symbols[i]
			if strings.EqualFold(symbol.Name, name) && symbol.ScopeID == scopeID {
				// Check visibility rules
				if sa.isSymbolAccessibleFromScope(*symbol, currentScope) {
					return symbol
				}
			}
		}
	}

	return nil
}

// buildScopeHierarchy builds a list of scope IDs to search, ordered from most specific to most general
func (sa *SemanticAnalyzer) buildScopeHierarchy(currentScope *ScopeInfo) []string {
	var hierarchy []string

	// Start with current scope
	if currentScope != nil {
		hierarchy = append(hierarchy, currentScope.ID)

		// Walk up the parent chain
		parentID := currentScope.Parent
		for parentID != "" {
			hierarchy = append(hierarchy, parentID)

			// Find the parent scope
			parentScope := sa.findScopeByID(parentID)
			if parentScope == nil {
				break
			}
			parentID = parentScope.Parent
		}
	}

	// Always include global scope as the last resort
	if len(hierarchy) == 0 || hierarchy[len(hierarchy)-1] != "global" {
		hierarchy = append(hierarchy, "global")
	}

	return hierarchy
}

// isSymbolAccessibleFromScope checks if a symbol is accessible from a given scope
// This is an enhanced version that considers the interpreter's visibility rules
func (sa *SemanticAnalyzer) isSymbolAccessibleFromScope(symbol EnhancedSymbol, fromScope *ScopeInfo) bool {
	// Public symbols are always accessible
	if symbol.Visibility == VisibilityPublic {
		return true
	}

	// Private symbols - check scope rules
	if symbol.Visibility == VisibilityPrivate {
		// If it's a class member, check if we're in the same class
		if symbol.ParentClass != "" {
			return sa.isInSameClass(fromScope, symbol.ParentClass)
		}

		// For non-class members, check if we're in the same scope or a nested scope
		symbolScope := sa.findScopeByID(symbol.ScopeID)
		if symbolScope == nil {
			return false
		}

		// Check if fromScope is the same as or nested within symbolScope
		return fromScope.ID == symbolScope.ID || sa.isScopeNestedIn(fromScope, symbolScope)
	}

	// Shared symbols (class members) - accessible if we can access the class
	if symbol.Visibility == VisibilityShared {
		return symbol.ParentClass != ""
	}

	return true
}

// Helper methods

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
	// First, try fast lookup via position map
	key := PositionKey{Line: int(position.Line) + 1, Column: int(position.Character) + 1} // Convert from 0-based to 1-based
	if symbol, exists := sa.positionToSymbol[key]; exists {
		return symbol
	}

	// Fallback to range-based search
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

// GetAllSymbolsAtPosition returns all symbols that exist at the given position
func (sa *SemanticAnalyzer) GetAllSymbolsAtPosition(position protocol.Position) []*EnhancedSymbol {
	var results []*EnhancedSymbol

	// First, try fast lookup via position map
	key := PositionKey{Line: int(position.Line) + 1, Column: int(position.Character) + 1}
	if symbol, exists := sa.positionToSymbol[key]; exists {
		results = append(results, symbol)
	}

	// Also check for symbols in range (for broader matches)
	scope := sa.findScopeAtPosition(position)
	if scope == nil {
		scope = &ScopeInfo{ID: "global", Type: ScopeTypeGlobal}
	}

	for i := range sa.symbolTable.Symbols {
		symbol := &sa.symbolTable.Symbols[i]
		if sa.positionInRange(position, symbol.Range) {
			// Check if symbol is accessible from current scope
			if sa.isSymbolAccessible(*symbol, scope) {
				// Avoid duplicates
				found := false
				for _, existing := range results {
					if existing == symbol {
						found = true
						break
					}
				}
				if !found {
					results = append(results, symbol)
				}
			}
		}
	}

	return results
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
// This is a wrapper around the enhanced isSymbolAccessibleFromScope method
func (sa *SemanticAnalyzer) isSymbolAccessible(symbol EnhancedSymbol, currentScope *ScopeInfo) bool {
	return sa.isSymbolAccessibleFromScope(symbol, currentScope)
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
	// Look for symbol definitions at position
	symbol := sa.ResolveSymbolAtPosition(position)
	if symbol == nil {
		return nil
	}

	// Builtin types
	switch symbol.Name {
	case "STRIN":
		return &protocol.Hover{
			Contents: protocol.MarkupContent{
				Kind:  protocol.MarkupKindMarkdown,
				Value: "**STRIN**: Represents text strings.",
			},
			Range: &symbol.Range,
		}
	case "INTEGR":
		return &protocol.Hover{
			Contents: protocol.MarkupContent{
				Kind:  protocol.MarkupKindMarkdown,
				Value: "**INTEGR**: Represents integer numbers.",
			},
			Range: &symbol.Range,
		}
	case "DUBBLE":
		return &protocol.Hover{
			Contents: protocol.MarkupContent{
				Kind:  protocol.MarkupKindMarkdown,
				Value: "**DUBBLE**: Represents double-precision floating-point numbers.",
			},
			Range: &symbol.Range,
		}
	case "BOOL":
		return &protocol.Hover{
			Contents: protocol.MarkupContent{
				Kind:  protocol.MarkupKindMarkdown,
				Value: "**BOOL**: Represents boolean values (`YEZ` or `NO`).",
			},
			Range: &symbol.Range,
		}
	case "NOTHIN":
		return &protocol.Hover{
			Contents: protocol.MarkupContent{
				Kind:  protocol.MarkupKindMarkdown,
				Value: "**NOTHIN**: Represents the absence of a value.",
			},
			Range: &symbol.Range,
		}
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

// analyzeFunctionCall analyzes a function call and creates identifier reference symbols
func (sa *SemanticAnalyzer) analyzeFunctionCall(node *ast.FunctionCallNode) {
	if node == nil {
		return
	}

	// Extract function identifier and create reference symbol
	switch funcNode := node.Function.(type) {
	case *ast.IdentifierNode:
		// Global function call - track the identifier
		sa.trackIdentifierReference(funcNode)

	case *ast.MemberAccessNode:
		// Method call - use type-aware member reference tracking
		sa.trackMemberReference(funcNode)
	}

	// Analyze function arguments
	for _, arg := range node.Arguments {
		sa.analyzeExpression(arg)
	}
}
