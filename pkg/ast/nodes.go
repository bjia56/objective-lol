package ast

import (
	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/types"
)

// Node represents the base interface for all AST nodes
type Node interface {
	Accept(visitor Visitor) (types.Value, error)
}

// Visitor defines the visitor pattern for AST traversal
type Visitor interface {
	VisitProgram(node *ProgramNode) (types.Value, error)
	VisitVariableDeclaration(node *VariableDeclarationNode) (types.Value, error)
	VisitFunctionDeclaration(node *FunctionDeclarationNode) (types.Value, error)
	VisitClassDeclaration(node *ClassDeclarationNode) (types.Value, error)
	VisitAssignment(node *AssignmentNode) (types.Value, error)
	VisitIfStatement(node *IfStatementNode) (types.Value, error)
	VisitWhileStatement(node *WhileStatementNode) (types.Value, error)
	VisitReturnStatement(node *ReturnStatementNode) (types.Value, error)
	VisitFunctionCall(node *FunctionCallNode) (types.Value, error)
	VisitMemberAccess(node *MemberAccessNode) (types.Value, error)
	VisitBinaryOp(node *BinaryOpNode) (types.Value, error)
	VisitUnaryOp(node *UnaryOpNode) (types.Value, error)
	VisitCast(node *CastNode) (types.Value, error)
	VisitLiteral(node *LiteralNode) (types.Value, error)
	VisitIdentifier(node *IdentifierNode) (types.Value, error)
	VisitObjectInstantiation(node *ObjectInstantiationNode) (types.Value, error)
	VisitStatementBlock(node *StatementBlockNode) (types.Value, error)
}

// ProgramNode represents the root of the AST
type ProgramNode struct {
	Declarations []Node
}

func (n *ProgramNode) Accept(visitor Visitor) (types.Value, error) {
	return visitor.VisitProgram(n)
}

// VariableDeclarationNode represents variable declarations
type VariableDeclarationNode struct {
	Name     string
	Type     string
	Value    Node
	IsLocked bool
}

func (n *VariableDeclarationNode) Accept(visitor Visitor) (types.Value, error) {
	return visitor.VisitVariableDeclaration(n)
}

// FunctionDeclarationNode represents function declarations
type FunctionDeclarationNode struct {
	Name       string
	ReturnType string
	Parameters []environment.Parameter
	Body       *StatementBlockNode
	IsNative   bool
	IsShared   *bool // nil for global, true/false for class methods
}

func (n *FunctionDeclarationNode) Accept(visitor Visitor) (types.Value, error) {
	return visitor.VisitFunctionDeclaration(n)
}

// ClassDeclarationNode represents class declarations
type ClassDeclarationNode struct {
	Name        string
	ParentClass string
	Members     []*ClassMemberNode
}

func (n *ClassDeclarationNode) Accept(visitor Visitor) (types.Value, error) {
	return visitor.VisitClassDeclaration(n)
}

// ClassMemberNode represents class members (variables and functions)
type ClassMemberNode struct {
	IsPublic   bool
	IsShared   bool
	IsVariable bool
	Variable   *VariableDeclarationNode
	Function   *FunctionDeclarationNode
}

// AssignmentNode represents variable assignments
type AssignmentNode struct {
	Target Node // IdentifierNode or MemberAccessNode
	Value  Node
}

func (n *AssignmentNode) Accept(visitor Visitor) (types.Value, error) {
	return visitor.VisitAssignment(n)
}

// IfStatementNode represents if statements
type IfStatementNode struct {
	Condition Node
	ThenBlock *StatementBlockNode
	ElseBlock *StatementBlockNode
}

func (n *IfStatementNode) Accept(visitor Visitor) (types.Value, error) {
	return visitor.VisitIfStatement(n)
}

// WhileStatementNode represents while loops
type WhileStatementNode struct {
	Condition Node
	Body      *StatementBlockNode
}

func (n *WhileStatementNode) Accept(visitor Visitor) (types.Value, error) {
	return visitor.VisitWhileStatement(n)
}

// ReturnStatementNode represents return statements
type ReturnStatementNode struct {
	Value Node // nil for "GIVEZ UP"
}

func (n *ReturnStatementNode) Accept(visitor Visitor) (types.Value, error) {
	return visitor.VisitReturnStatement(n)
}

// FunctionCallNode represents function calls
type FunctionCallNode struct {
	Function  Node // IdentifierNode or MemberAccessNode
	Arguments []Node
}

func (n *FunctionCallNode) Accept(visitor Visitor) (types.Value, error) {
	return visitor.VisitFunctionCall(n)
}

// MemberAccessNode represents member access (obj IN member)
type MemberAccessNode struct {
	Object Node
	Member string
}

func (n *MemberAccessNode) Accept(visitor Visitor) (types.Value, error) {
	return visitor.VisitMemberAccess(n)
}

// BinaryOpNode represents binary operations
type BinaryOpNode struct {
	Left     Node
	Operator string
	Right    Node
}

func (n *BinaryOpNode) Accept(visitor Visitor) (types.Value, error) {
	return visitor.VisitBinaryOp(n)
}

// UnaryOpNode represents unary operations
type UnaryOpNode struct {
	Operator string
	Operand  Node
}

func (n *UnaryOpNode) Accept(visitor Visitor) (types.Value, error) {
	return visitor.VisitUnaryOp(n)
}

// CastNode represents type casting
type CastNode struct {
	Expression Node
	TargetType string
}

func (n *CastNode) Accept(visitor Visitor) (types.Value, error) {
	return visitor.VisitCast(n)
}

// LiteralNode represents literal values
type LiteralNode struct {
	Value types.Value
}

func (n *LiteralNode) Accept(visitor Visitor) (types.Value, error) {
	return visitor.VisitLiteral(n)
}

// IdentifierNode represents identifiers
type IdentifierNode struct {
	Name string
}

func (n *IdentifierNode) Accept(visitor Visitor) (types.Value, error) {
	return visitor.VisitIdentifier(n)
}

// ObjectInstantiationNode represents object creation
type ObjectInstantiationNode struct {
	ClassName  string
	SourceName string // optional, for "NEW TYPE IN SOURCE"
}

func (n *ObjectInstantiationNode) Accept(visitor Visitor) (types.Value, error) {
	return visitor.VisitObjectInstantiation(n)
}

// StatementBlockNode represents a block of statements
type StatementBlockNode struct {
	Statements []Node
}

func (n *StatementBlockNode) Accept(visitor Visitor) (types.Value, error) {
	return visitor.VisitStatementBlock(n)
}

// ReturnValue is used to control return flow in function execution
type ReturnValue struct {
	Value types.Value
}

func (r ReturnValue) Error() string {
	return "return"
}

// Helper function to check if an error is a return value
func IsReturnValue(err error) bool {
	_, ok := err.(ReturnValue)
	return ok
}

// Helper function to extract return value from error
func GetReturnValue(err error) types.Value {
	if retVal, ok := err.(ReturnValue); ok {
		return retVal.Value
	}
	return types.NOTHIN
}
