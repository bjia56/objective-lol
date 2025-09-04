package analyzer

import (
	protocol "github.com/tliron/glsp/protocol_3_16"

	"github.com/bjia56/objective-lol/pkg/ast"
	"github.com/bjia56/objective-lol/pkg/stdlib"
	"github.com/bjia56/objective-lol/pkg/types"
)

// SymbolCollector implements the visitor pattern to collect symbols from AST
type SymbolCollector struct {
	uri         string
	symbolTable *SymbolTable
}

// NewSymbolCollector creates a new symbol collector
func NewSymbolCollector(uri string) *SymbolCollector {
	return &SymbolCollector{
		uri: uri,
		symbolTable: &SymbolTable{
			URI:     uri,
			Symbols: []Symbol{},
		},
	}
}

// GetSymbolTable returns the collected symbol table
func (sc *SymbolCollector) GetSymbolTable() *SymbolTable {
	return sc.symbolTable
}

// addSymbol adds a symbol to the table
func (sc *SymbolCollector) addSymbol(name string, kind SymbolKind, symbolType string, pos ast.PositionInfo) {
	for _, existingSymbol := range sc.symbolTable.Symbols {
		if existingSymbol.Name == name {
			// Symbol already exists, let's make a copy of it with updated position
			// We assume that the first occurrence is the declaration where we know the type
			newSymbol := existingSymbol
			newSymbol.Position = pos
			newSymbol.Range = sc.positionInfoToRange(pos, len(name))
			sc.symbolTable.Symbols = append(sc.symbolTable.Symbols, newSymbol)
			return
		}
	}

	if symbolType == "" {
		symbolType = "NOTHIN"
	}

	symbol := Symbol{
		Name:     name,
		Kind:     kind,
		Type:     symbolType,
		Position: pos,
		Range:    sc.positionInfoToRange(pos, len(name)),
	}
	sc.symbolTable.Symbols = append(sc.symbolTable.Symbols, symbol)
}

// Visitor interface implementation

func (sc *SymbolCollector) VisitProgram(node *ast.ProgramNode) (types.Value, error) {
	// Define stdlib symbols
	for _, globalInitializers := range stdlib.DefaultGlobalInitializers() {
		for _, decl := range stdlib.GetStdlibDefinitions(globalInitializers) {
			switch decl.Kind {
			case stdlib.StdlibDefinitionKindFunction:
				sc.addSymbol(decl.Name, SymbolKindFunction, decl.Type, ast.PositionInfo{})
			case stdlib.StdlibDefinitionKindClass:
				sc.addSymbol(decl.Name, SymbolKindClass, decl.Type, ast.PositionInfo{})
			}
		}
	}

	for _, decl := range node.Declarations {
		decl.Accept(sc)
	}
	return nil, nil
}

func (sc *SymbolCollector) VisitImportStatement(node *ast.ImportStatementNode) (types.Value, error) {
	sc.addSymbol(node.ModuleName, SymbolKindImport, "module", node.GetPosition())

	if !node.IsFileImport {
		if moduleInitializer, exists := stdlib.DefaultStdlibInitializers()[node.ModuleName]; exists {
			for _, decl := range stdlib.GetStdlibDefinitions(moduleInitializer) {
				switch decl.Kind {
				case stdlib.StdlibDefinitionKindFunction:
					sc.addSymbol(decl.Name, SymbolKindFunction, decl.Type, ast.PositionInfo{})
				case stdlib.StdlibDefinitionKindClass:
					sc.addSymbol(decl.Name, SymbolKindClass, decl.Type, ast.PositionInfo{})
				}
			}
		}
	}

	return nil, nil
}

func (sc *SymbolCollector) VisitVariableDeclaration(node *ast.VariableDeclarationNode) (types.Value, error) {
	sc.addSymbol(node.Name, SymbolKindVariable, node.Type, node.GetPosition())
	if node.Value != nil {
		node.Value.Accept(sc)
	}
	return nil, nil
}

func (sc *SymbolCollector) VisitFunctionDeclaration(node *ast.FunctionDeclarationNode) (types.Value, error) {
	sc.addSymbol(node.Name, SymbolKindFunction, node.ReturnType, node.GetPosition())

	// Add parameters as symbols
	for _, param := range node.Parameters {
		// Note: Parameters don't have position info in current implementation
		// We would need to enhance the parser to track parameter positions
		sc.addSymbol(param.Name, SymbolKindParameter, param.Type, ast.PositionInfo{})
	}

	if node.Body != nil {
		node.Body.Accept(sc)
	}
	return nil, nil
}

func (sc *SymbolCollector) VisitClassDeclaration(node *ast.ClassDeclarationNode) (types.Value, error) {
	sc.addSymbol(node.Name, SymbolKindClass, node.Name, node.GetPosition())

	// Visit class members
	for _, member := range node.Members {
		if member.IsVariable && member.Variable != nil {
			member.Variable.Accept(sc)
		}
		if !member.IsVariable && member.Function != nil {
			member.Function.Accept(sc)
		}
	}
	return nil, nil
}

func (sc *SymbolCollector) VisitAssignment(node *ast.AssignmentNode) (types.Value, error) {
	if node.Value != nil {
		node.Value.Accept(sc)
	}
	return nil, nil
}

func (sc *SymbolCollector) VisitIfStatement(node *ast.IfStatementNode) (types.Value, error) {
	if node.Condition != nil {
		node.Condition.Accept(sc)
	}
	if node.ThenBlock != nil {
		node.ThenBlock.Accept(sc)
	}
	if node.ElseBlock != nil {
		node.ElseBlock.Accept(sc)
	}
	return nil, nil
}

func (sc *SymbolCollector) VisitWhileStatement(node *ast.WhileStatementNode) (types.Value, error) {
	if node.Condition != nil {
		node.Condition.Accept(sc)
	}
	if node.Body != nil {
		node.Body.Accept(sc)
	}
	return nil, nil
}

func (sc *SymbolCollector) VisitReturnStatement(node *ast.ReturnStatementNode) (types.Value, error) {
	if node.Value != nil {
		node.Value.Accept(sc)
	}
	return nil, nil
}

func (sc *SymbolCollector) VisitFunctionCall(node *ast.FunctionCallNode) (types.Value, error) {
	// Visit the function being called (could be an identifier or member access)
	node.Function.Accept(sc)

	if node.Arguments != nil {
		for _, arg := range node.Arguments {
			if arg != nil {
				arg.Accept(sc)
			}
		}
	}
	return nil, nil
}

func (sc *SymbolCollector) VisitMemberAccess(node *ast.MemberAccessNode) (types.Value, error) {
	if node.Object != nil {
		node.Object.Accept(sc)
	}
	return nil, nil
}

func (sc *SymbolCollector) VisitBinaryOp(node *ast.BinaryOpNode) (types.Value, error) {
	if node.Left != nil {
		node.Left.Accept(sc)
	}
	if node.Right != nil {
		node.Right.Accept(sc)
	}
	return nil, nil
}

func (sc *SymbolCollector) VisitUnaryOp(node *ast.UnaryOpNode) (types.Value, error) {
	if node.Operand != nil {
		node.Operand.Accept(sc)
	}
	return nil, nil
}

func (sc *SymbolCollector) VisitCast(node *ast.CastNode) (types.Value, error) {
	if node.Expression != nil {
		node.Expression.Accept(sc)
	}
	return nil, nil
}

func (sc *SymbolCollector) VisitLiteral(node *ast.LiteralNode) (types.Value, error) {
	// Literals don't create symbols
	return nil, nil
}

func (sc *SymbolCollector) VisitIdentifier(node *ast.IdentifierNode) (types.Value, error) {
	sc.addSymbol(node.Name, SymbolKindUnknown, "", node.GetPosition())
	return nil, nil
}

func (sc *SymbolCollector) VisitObjectInstantiation(node *ast.ObjectInstantiationNode) (types.Value, error) {
	if node.ConstructorArgs != nil {
		for _, arg := range node.ConstructorArgs {
			if arg != nil {
				arg.Accept(sc)
			}
		}
	}
	return nil, nil
}

func (sc *SymbolCollector) VisitStatementBlock(node *ast.StatementBlockNode) (types.Value, error) {
	if node.Statements != nil {
		for _, stmt := range node.Statements {
			if stmt != nil {
				stmt.Accept(sc)
			}
		}
	}
	return nil, nil
}

func (sc *SymbolCollector) VisitTryStatement(node *ast.TryStatementNode) (types.Value, error) {
	if node.TryBody != nil {
		node.TryBody.Accept(sc)
	}
	if node.CatchBody != nil {
		node.CatchBody.Accept(sc)
	}
	if node.FinallyBody != nil {
		node.FinallyBody.Accept(sc)
	}
	return nil, nil
}

func (sc *SymbolCollector) VisitThrowStatement(node *ast.ThrowStatementNode) (types.Value, error) {
	if node.Expression != nil {
		node.Expression.Accept(sc)
	}
	return nil, nil
}

// Helper method to convert position info to LSP range
func (sc *SymbolCollector) positionInfoToRange(pos ast.PositionInfo, length int) protocol.Range {
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
