package ast

import (
	"github.com/bjia56/objective-lol/pkg/environment"
)

// PositionInfo represents position information in source code
type PositionInfo struct {
	Line   int
	Column int
}

// Node represents the base interface for all AST nodes
type Node interface {
	Accept(visitor Visitor) (environment.Value, error)
	GetPosition() PositionInfo
	SetPosition(pos PositionInfo)
}

// Visitor defines the visitor pattern for AST traversal
type Visitor interface {
	VisitProgram(node *ProgramNode) (environment.Value, error)
	VisitImportStatement(node *ImportStatementNode) (environment.Value, error)
	VisitVariableDeclaration(node *VariableDeclarationNode) (environment.Value, error)
	VisitFunctionDeclaration(node *FunctionDeclarationNode) (environment.Value, error)
	VisitClassDeclaration(node *ClassDeclarationNode) (environment.Value, error)
	VisitAssignment(node *AssignmentNode) (environment.Value, error)
	VisitIfStatement(node *IfStatementNode) (environment.Value, error)
	VisitWhileStatement(node *WhileStatementNode) (environment.Value, error)
	VisitReturnStatement(node *ReturnStatementNode) (environment.Value, error)
	VisitFunctionCall(node *FunctionCallNode) (environment.Value, error)
	VisitMemberAccess(node *MemberAccessNode) (environment.Value, error)
	VisitBinaryOp(node *BinaryOpNode) (environment.Value, error)
	VisitUnaryOp(node *UnaryOpNode) (environment.Value, error)
	VisitCast(node *CastNode) (environment.Value, error)
	VisitLiteral(node *LiteralNode) (environment.Value, error)
	VisitIdentifier(node *IdentifierNode) (environment.Value, error)
	VisitObjectInstantiation(node *ObjectInstantiationNode) (environment.Value, error)
	VisitStatementBlock(node *StatementBlockNode) (environment.Value, error)
	VisitTryStatement(node *TryStatementNode) (environment.Value, error)
	VisitThrowStatement(node *ThrowStatementNode) (environment.Value, error)
}

// ProgramNode represents the root of the AST
type ProgramNode struct {
	Declarations []Node
	Position     PositionInfo
}

func (n *ProgramNode) Accept(visitor Visitor) (environment.Value, error) {
	return visitor.VisitProgram(n)
}

func (n *ProgramNode) GetPosition() PositionInfo {
	return n.Position
}

func (n *ProgramNode) SetPosition(pos PositionInfo) {
	n.Position = pos
}

// ImportStatementNode represents module import statements
type ImportStatementNode struct {
	ModuleName   string
	Declarations []string // Specific declarations to import (empty means import all)
	IsFileImport bool     // Whether this imports a file (true) or built-in module (false)
	Position     PositionInfo
}

func (n *ImportStatementNode) Accept(visitor Visitor) (environment.Value, error) {
	return visitor.VisitImportStatement(n)
}

func (n *ImportStatementNode) GetPosition() PositionInfo {
	return n.Position
}

func (n *ImportStatementNode) SetPosition(pos PositionInfo) {
	n.Position = pos
}

// VariableDeclarationNode represents variable declarations
type VariableDeclarationNode struct {
	Name          string
	Type          string
	Value         Node
	IsLocked      bool
	Documentation []string // Documentation comments preceding the variable
	Position      PositionInfo
}

func (n *VariableDeclarationNode) Accept(visitor Visitor) (environment.Value, error) {
	return visitor.VisitVariableDeclaration(n)
}

func (n *VariableDeclarationNode) GetPosition() PositionInfo {
	return n.Position
}

func (n *VariableDeclarationNode) SetPosition(pos PositionInfo) {
	n.Position = pos
}

// FunctionDeclarationNode represents function declarations
type FunctionDeclarationNode struct {
	Name          string
	ReturnType    string
	Parameters    []environment.Parameter
	Body          *StatementBlockNode
	IsShared      *bool    // nil for global, true/false for class methods
	Documentation []string // Documentation comments preceding the function
	Position      PositionInfo
}

func (n *FunctionDeclarationNode) Accept(visitor Visitor) (environment.Value, error) {
	return visitor.VisitFunctionDeclaration(n)
}

func (n *FunctionDeclarationNode) GetPosition() PositionInfo {
	return n.Position
}

func (n *FunctionDeclarationNode) SetPosition(pos PositionInfo) {
	n.Position = pos
}

// ClassDeclarationNode represents class declarations
type ClassDeclarationNode struct {
	Name          string
	ParentClasses []string
	Members       []*ClassMemberNode
	Documentation []string
	Position      PositionInfo
}

func (n *ClassDeclarationNode) Accept(visitor Visitor) (environment.Value, error) {
	return visitor.VisitClassDeclaration(n)
}

func (n *ClassDeclarationNode) GetPosition() PositionInfo {
	return n.Position
}

func (n *ClassDeclarationNode) SetPosition(pos PositionInfo) {
	n.Position = pos
}

// ClassMemberNode represents class members (variables and functions)
type ClassMemberNode struct {
	IsPublic   bool
	IsShared   bool
	IsVariable bool
	Variable   *VariableDeclarationNode
	Function   *FunctionDeclarationNode
	Position   PositionInfo
}

func (n *ClassMemberNode) GetPosition() PositionInfo {
	return n.Position
}

func (n *ClassMemberNode) SetPosition(pos PositionInfo) {
	n.Position = pos
}

// AssignmentNode represents variable assignments
type AssignmentNode struct {
	Target   Node // IdentifierNode or MemberAccessNode
	Value    Node
	Position PositionInfo
}

func (n *AssignmentNode) Accept(visitor Visitor) (environment.Value, error) {
	return visitor.VisitAssignment(n)
}

func (n *AssignmentNode) GetPosition() PositionInfo {
	return n.Position
}

func (n *AssignmentNode) SetPosition(pos PositionInfo) {
	n.Position = pos
}

// ElseIfBranch represents a MEBBE (else-if) branch
type ElseIfBranch struct {
	Condition Node
	Block     *StatementBlockNode
	Position  PositionInfo
}

// IfStatementNode represents if statements
type IfStatementNode struct {
	Condition      Node
	ThenBlock      *StatementBlockNode
	ElseIfBranches []*ElseIfBranch
	ElseBlock      *StatementBlockNode
	Position       PositionInfo
}

func (n *IfStatementNode) Accept(visitor Visitor) (environment.Value, error) {
	return visitor.VisitIfStatement(n)
}

func (n *IfStatementNode) GetPosition() PositionInfo {
	return n.Position
}

func (n *IfStatementNode) SetPosition(pos PositionInfo) {
	n.Position = pos
}

// WhileStatementNode represents while loops
type WhileStatementNode struct {
	Condition Node
	Body      *StatementBlockNode
	Position  PositionInfo
}

func (n *WhileStatementNode) Accept(visitor Visitor) (environment.Value, error) {
	return visitor.VisitWhileStatement(n)
}

func (n *WhileStatementNode) GetPosition() PositionInfo {
	return n.Position
}

func (n *WhileStatementNode) SetPosition(pos PositionInfo) {
	n.Position = pos
}

// ReturnStatementNode represents return statements
type ReturnStatementNode struct {
	Value    Node // nil for "GIVEZ UP"
	Position PositionInfo
}

func (n *ReturnStatementNode) Accept(visitor Visitor) (environment.Value, error) {
	return visitor.VisitReturnStatement(n)
}

func (n *ReturnStatementNode) GetPosition() PositionInfo {
	return n.Position
}

func (n *ReturnStatementNode) SetPosition(pos PositionInfo) {
	n.Position = pos
}

// FunctionCallNode represents function calls
type FunctionCallNode struct {
	Function  Node // IdentifierNode or MemberAccessNode
	Arguments []Node
	Position  PositionInfo
}

func (n *FunctionCallNode) Accept(visitor Visitor) (environment.Value, error) {
	return visitor.VisitFunctionCall(n)
}

func (n *FunctionCallNode) GetPosition() PositionInfo {
	return n.Position
}

func (n *FunctionCallNode) SetPosition(pos PositionInfo) {
	n.Position = pos
}

// MemberAccessNode represents member access (obj IN member)
type MemberAccessNode struct {
	Object   Node
	Member   string
	Position PositionInfo
}

func (n *MemberAccessNode) Accept(visitor Visitor) (environment.Value, error) {
	return visitor.VisitMemberAccess(n)
}

func (n *MemberAccessNode) GetPosition() PositionInfo {
	return n.Position
}

func (n *MemberAccessNode) SetPosition(pos PositionInfo) {
	n.Position = pos
}

// BinaryOpNode represents binary operations
type BinaryOpNode struct {
	Left     Node
	Operator string
	Right    Node
	Position PositionInfo
}

func (n *BinaryOpNode) Accept(visitor Visitor) (environment.Value, error) {
	return visitor.VisitBinaryOp(n)
}

func (n *BinaryOpNode) GetPosition() PositionInfo {
	return n.Position
}

func (n *BinaryOpNode) SetPosition(pos PositionInfo) {
	n.Position = pos
}

// UnaryOpNode represents unary operations
type UnaryOpNode struct {
	Operator string
	Operand  Node
	Position PositionInfo
}

func (n *UnaryOpNode) Accept(visitor Visitor) (environment.Value, error) {
	return visitor.VisitUnaryOp(n)
}

func (n *UnaryOpNode) GetPosition() PositionInfo {
	return n.Position
}

func (n *UnaryOpNode) SetPosition(pos PositionInfo) {
	n.Position = pos
}

// CastNode represents type casting
type CastNode struct {
	Expression Node
	TargetType string
	Position   PositionInfo
}

func (n *CastNode) Accept(visitor Visitor) (environment.Value, error) {
	return visitor.VisitCast(n)
}

func (n *CastNode) GetPosition() PositionInfo {
	return n.Position
}

func (n *CastNode) SetPosition(pos PositionInfo) {
	n.Position = pos
}

// LiteralNode represents literal values
type LiteralNode struct {
	Value    environment.Value
	Position PositionInfo
}

func (n *LiteralNode) Accept(visitor Visitor) (environment.Value, error) {
	return visitor.VisitLiteral(n)
}

func (n *LiteralNode) GetPosition() PositionInfo {
	return n.Position
}

func (n *LiteralNode) SetPosition(pos PositionInfo) {
	n.Position = pos
}

// IdentifierNode represents identifiers
type IdentifierNode struct {
	Name     string
	Position PositionInfo
}

func (n *IdentifierNode) Accept(visitor Visitor) (environment.Value, error) {
	return visitor.VisitIdentifier(n)
}

func (n *IdentifierNode) GetPosition() PositionInfo {
	return n.Position
}

func (n *IdentifierNode) SetPosition(pos PositionInfo) {
	n.Position = pos
}

// ObjectInstantiationNode represents object creation
type ObjectInstantiationNode struct {
	ClassName       string
	ConstructorArgs []Node // arguments for constructor call
	Position        PositionInfo
}

func (n *ObjectInstantiationNode) Accept(visitor Visitor) (environment.Value, error) {
	return visitor.VisitObjectInstantiation(n)
}

func (n *ObjectInstantiationNode) GetPosition() PositionInfo {
	return n.Position
}

func (n *ObjectInstantiationNode) SetPosition(pos PositionInfo) {
	n.Position = pos
}

// StatementBlockNode represents a block of statements
type StatementBlockNode struct {
	Statements []Node
	Position   PositionInfo
}

func (n *StatementBlockNode) Accept(visitor Visitor) (environment.Value, error) {
	return visitor.VisitStatementBlock(n)
}

func (n *StatementBlockNode) GetPosition() PositionInfo {
	return n.Position
}

func (n *StatementBlockNode) SetPosition(pos PositionInfo) {
	n.Position = pos
}

// TryStatementNode represents try-catch-finally blocks
type TryStatementNode struct {
	TryBody     *StatementBlockNode
	CatchVar    string // Variable name to bind the exception message
	CatchBody   *StatementBlockNode
	FinallyBody *StatementBlockNode // Optional finally block
	Position    PositionInfo
}

func (n *TryStatementNode) Accept(visitor Visitor) (environment.Value, error) {
	return visitor.VisitTryStatement(n)
}

func (n *TryStatementNode) GetPosition() PositionInfo {
	return n.Position
}

func (n *TryStatementNode) SetPosition(pos PositionInfo) {
	n.Position = pos
}

// ThrowStatementNode represents throw statements
type ThrowStatementNode struct {
	Expression Node // Expression that evaluates to the error message (string)
	Position   PositionInfo
}

func (n *ThrowStatementNode) Accept(visitor Visitor) (environment.Value, error) {
	return visitor.VisitThrowStatement(n)
}

func (n *ThrowStatementNode) GetPosition() PositionInfo {
	return n.Position
}

func (n *ThrowStatementNode) SetPosition(pos PositionInfo) {
	n.Position = pos
}
