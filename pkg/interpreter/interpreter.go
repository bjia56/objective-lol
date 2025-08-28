package interpreter

import (
	"fmt"
	"strings"

	"github.com/bjia56/objective-lol/pkg/ast"
	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/stdlib"
	"github.com/bjia56/objective-lol/pkg/types"
)

// Interpreter implements the tree-walking interpreter for Objective-LOL
type Interpreter struct {
	runtime       *environment.RuntimeEnvironment
	environment   *environment.Environment
	currentClass  string                      // For tracking visibility context
	currentObject *environment.ObjectInstance // For tracking current object instance in method calls
}

// NewInterpreter creates a new interpreter instance
func NewInterpreter() *Interpreter {
	runtime := environment.NewRuntimeEnvironment()
	return &Interpreter{
		runtime:     runtime,
		environment: runtime.GlobalEnv,
	}
}

// Interpret executes the given AST
func (i *Interpreter) Interpret(program *ast.ProgramNode) error {
	_, err := program.Accept(i)
	return err
}

// VisitProgram executes the entire program
func (i *Interpreter) VisitProgram(node *ast.ProgramNode) (types.Value, error) {
	// First pass: process import statements
	for _, decl := range node.Declarations {
		switch n := decl.(type) {
		case *ast.ImportStatementNode:
			if _, err := i.VisitImportStatement(n); err != nil {
				return types.NOTHIN, err
			}
		}
	}

	// Second pass: declare all functions and classes
	for _, decl := range node.Declarations {
		switch n := decl.(type) {
		case *ast.FunctionDeclarationNode:
			if _, err := i.VisitFunctionDeclaration(n); err != nil {
				return types.NOTHIN, err
			}
		case *ast.ClassDeclarationNode:
			if _, err := i.VisitClassDeclaration(n); err != nil {
				return types.NOTHIN, err
			}
		}
	}

	// Third pass: execute variable declarations and other statements
	for _, decl := range node.Declarations {
		switch n := decl.(type) {
		case *ast.VariableDeclarationNode:
			if _, err := i.VisitVariableDeclaration(n); err != nil {
				return types.NOTHIN, err
			}
		}
	}

	// Look for and execute MAIN function
	if mainFunc, err := i.environment.GetFunction("MAIN"); err == nil {
		return i.callFunction(mainFunc, []types.Value{})
	}

	return types.NOTHIN, nil
}

// VisitImportStatement handles module import statements
func (i *Interpreter) VisitImportStatement(node *ast.ImportStatementNode) (types.Value, error) {
	moduleName := strings.ToUpper(node.ModuleName)

	// Load the requested module
	switch moduleName {
	case "STDIO":
		stdlib.RegisterSTDIO(i.runtime)
	case "MATH":
		stdlib.RegisterMATH(i.runtime)
	case "TIEM":
		stdlib.RegisterTIEM(i.runtime)
	default:
		return types.NOTHIN, fmt.Errorf("unknown module: %s", moduleName)
	}

	return types.NOTHIN, nil
}

// VisitVariableDeclaration handles variable declarations
func (i *Interpreter) VisitVariableDeclaration(node *ast.VariableDeclarationNode) (types.Value, error) {
	var value types.Value = types.NOTHIN
	var err error

	if node.Value != nil {
		value, err = node.Value.Accept(i)
		if err != nil {
			return types.NOTHIN, err
		}
	}

	err = i.environment.DefineVariable(
		strings.ToUpper(node.Name),
		strings.ToUpper(node.Type),
		value,
		node.IsLocked,
	)

	return types.NOTHIN, err
}

// VisitFunctionDeclaration handles function declarations
func (i *Interpreter) VisitFunctionDeclaration(node *ast.FunctionDeclarationNode) (types.Value, error) {
	function := &environment.Function{
		Name:        strings.ToUpper(node.Name),
		ReturnType:  strings.ToUpper(node.ReturnType),
		Parameters:  node.Parameters,
		Body:        node.Body,
		IsNative:    node.IsNative,
		IsShared:    node.IsShared,
		ParentClass: i.currentClass,
	}

	// Normalize parameter names and types
	for i := range function.Parameters {
		function.Parameters[i].Name = strings.ToUpper(function.Parameters[i].Name)
		function.Parameters[i].Type = strings.ToUpper(function.Parameters[i].Type)
	}

	return types.NOTHIN, i.environment.DefineFunction(function)
}

// VisitClassDeclaration handles class declarations
func (i *Interpreter) VisitClassDeclaration(node *ast.ClassDeclarationNode) (types.Value, error) {
	class := environment.NewClass(
		strings.ToUpper(node.Name),
		strings.ToUpper(node.ParentClass),
	)

	// Save current context
	oldClass := i.currentClass
	oldEnv := i.environment

	// Set class context for member processing
	i.currentClass = class.Name
	classEnv := environment.NewEnvironment(oldEnv)
	i.environment = classEnv

	// Process class members
	for _, member := range node.Members {
		if member.IsVariable {
			// Handle member variables
			var value types.Value = types.NOTHIN
			var err error

			if member.Variable.Value != nil {
				value, err = member.Variable.Value.Accept(i)
				if err != nil {
					return types.NOTHIN, err
				}
			}

			variable := &environment.Variable{
				Name:     strings.ToUpper(member.Variable.Name),
				Type:     strings.ToUpper(member.Variable.Type),
				Value:    value,
				IsLocked: member.Variable.IsLocked,
				IsPublic: member.IsPublic, // Use visibility from member declaration
			}

			// Add to appropriate collection based on visibility and sharing
			if member.IsShared {
				class.SharedVariables[variable.Name] = variable
			} else if member.IsPublic {
				class.PublicVariables[variable.Name] = variable
			} else {
				class.PrivateVariables[variable.Name] = variable
			}

		} else {
			// Handle member functions
			function := &environment.Function{
				Name:        strings.ToUpper(member.Function.Name),
				ReturnType:  strings.ToUpper(member.Function.ReturnType),
				Parameters:  member.Function.Parameters,
				Body:        member.Function.Body,
				IsNative:    member.Function.IsNative,
				IsShared:    &member.IsShared,
				ParentClass: class.Name,
			}

			// Normalize parameter names and types
			for i := range function.Parameters {
				function.Parameters[i].Name = strings.ToUpper(function.Parameters[i].Name)
				function.Parameters[i].Type = strings.ToUpper(function.Parameters[i].Type)
			}

			// Add to appropriate collection based on visibility and sharing
			if member.IsShared {
				class.SharedFunctions[function.Name] = function
			} else if member.IsPublic {
				class.PublicFunctions[function.Name] = function
			} else {
				class.PrivateFunctions[function.Name] = function
			}
		}
	}

	// Restore context
	i.currentClass = oldClass
	i.environment = oldEnv

	return types.NOTHIN, i.environment.DefineClass(class)
}

// VisitAssignment handles variable assignments
func (i *Interpreter) VisitAssignment(node *ast.AssignmentNode) (types.Value, error) {
	value, err := node.Value.Accept(i)
	if err != nil {
		return types.NOTHIN, err
	}

	switch target := node.Target.(type) {
	case *ast.IdentifierNode:
		// Simple assignment - check if it's a member variable in method context
		name := strings.ToUpper(target.Name)

		// If we're in a method call and the identifier refers to a member variable, treat as member assignment
		if i.currentObject != nil {
			if _, err := i.currentObject.GetMemberVariable(name, i.currentClass); err == nil {
				return value, i.currentObject.SetMemberVariable(name, value, i.currentClass)
			}
		}

		// Otherwise, treat as local variable assignment
		return value, i.environment.SetVariable(name, value)

	case *ast.MemberAccessNode:
		// Member assignment (obj IN member ITZ value)
		objectValue, err := target.Object.Accept(i)
		if err != nil {
			return types.NOTHIN, err
		}

		if objVal, ok := objectValue.(types.ObjectValue); ok {
			if obj, ok := objVal.Instance.(*environment.ObjectInstance); ok {
				return value, obj.SetMemberVariable(strings.ToUpper(target.Member), value, i.currentClass)
			}
		}
		return types.NOTHIN, fmt.Errorf("cannot assign to member of non-object")

	default:
		return types.NOTHIN, fmt.Errorf("invalid assignment target")
	}
}

// VisitIfStatement handles if statements
func (i *Interpreter) VisitIfStatement(node *ast.IfStatementNode) (types.Value, error) {
	condition, err := node.Condition.Accept(i)
	if err != nil {
		return types.NOTHIN, err
	}

	if condition.ToBool() == types.YEZ {
		return node.ThenBlock.Accept(i)
	} else if node.ElseBlock != nil {
		return node.ElseBlock.Accept(i)
	}

	return types.NOTHIN, nil
}

// VisitWhileStatement handles while loops
func (i *Interpreter) VisitWhileStatement(node *ast.WhileStatementNode) (types.Value, error) {
	for {
		condition, err := node.Condition.Accept(i)
		if err != nil {
			return types.NOTHIN, err
		}

		if condition.ToBool() != types.YEZ {
			break
		}

		_, err = node.Body.Accept(i)
		if err != nil {
			if ast.IsReturnValue(err) {
				return types.NOTHIN, err // Propagate return
			}
			return types.NOTHIN, err
		}
	}

	return types.NOTHIN, nil
}

// VisitReturnStatement handles return statements
func (i *Interpreter) VisitReturnStatement(node *ast.ReturnStatementNode) (types.Value, error) {
	var value types.Value = types.NOTHIN
	var err error

	if node.Value != nil {
		value, err = node.Value.Accept(i)
		if err != nil {
			return types.NOTHIN, err
		}
	}

	return types.NOTHIN, ast.ReturnValue{Value: value}
}

// VisitFunctionCall handles function calls
func (i *Interpreter) VisitFunctionCall(node *ast.FunctionCallNode) (types.Value, error) {
	// Evaluate arguments
	args := make([]types.Value, len(node.Arguments))
	for j, arg := range node.Arguments {
		val, err := arg.Accept(i)
		if err != nil {
			return types.NOTHIN, err
		}
		args[j] = val
	}

	switch funcNode := node.Function.(type) {
	case *ast.IdentifierNode:
		// Global function call
		functionName := strings.ToUpper(funcNode.Name)

		// Check local environment first
		if function, err := i.environment.GetFunction(functionName); err == nil {
			return i.callFunction(function, args)
		}

		// Check native functions
		if function, exists := i.runtime.GetNative(functionName); exists {
			return i.callFunction(function, args)
		}

		return types.NOTHIN, fmt.Errorf("undefined function '%s'", functionName)

	case *ast.MemberAccessNode:
		// Member function call (obj IN func WIT args)
		objectValue, err := funcNode.Object.Accept(i)
		if err != nil {
			return types.NOTHIN, err
		}

		if objVal, ok := objectValue.(types.ObjectValue); ok {
			if obj, ok := objVal.Instance.(*environment.ObjectInstance); ok {
				function, err := obj.GetMemberFunction(strings.ToUpper(funcNode.Member), i.currentClass, i.environment)
				if err != nil {
					return types.NOTHIN, err
				}
				return i.callMemberFunction(function, obj, args)
			}
		}
		return types.NOTHIN, fmt.Errorf("cannot call method on non-object")

	default:
		return types.NOTHIN, fmt.Errorf("invalid function call target")
	}
}

// callFunction executes a function with the given arguments
func (i *Interpreter) callFunction(function *environment.Function, args []types.Value) (types.Value, error) {
	// Check argument count
	if len(args) != len(function.Parameters) {
		return types.NOTHIN, fmt.Errorf("function '%s' expects %d arguments, got %d",
			function.Name, len(function.Parameters), len(args))
	}

	// Handle native functions
	if function.IsNative && function.NativeImpl != nil {
		argsCasted := make([]types.Value, len(args))
		for i, arg := range args {
			casted, err := arg.Cast(function.Parameters[i].Type)
			if err != nil {
				return types.NOTHIN, fmt.Errorf("cannot cast function argument %s to %s: %v",
					function.Parameters[i].Name, function.Parameters[i].Type, err)
			}
			argsCasted[i] = casted
		}
		return function.NativeImpl(argsCasted)
	}

	// Create new environment for function execution
	funcEnv := environment.NewEnvironment(i.environment)

	// Bind parameters
	for j, param := range function.Parameters {
		err := funcEnv.DefineVariable(param.Name, param.Type, args[j], false)
		if err != nil {
			return types.NOTHIN, err
		}
	}

	// Save and set environment
	oldEnv := i.environment
	i.environment = funcEnv

	// Execute function body
	var result types.Value = types.NOTHIN
	if body, ok := function.Body.(*ast.StatementBlockNode); ok {
		_, err := body.Accept(i)
		if err != nil {
			if ast.IsReturnValue(err) {
				result = ast.GetReturnValue(err)
			} else {
				i.environment = oldEnv
				return types.NOTHIN, err
			}
		}
	}

	// Restore environment
	i.environment = oldEnv

	// Cast result to return type if specified
	if function.ReturnType != "" {
		castedResult, err := result.Cast(function.ReturnType)
		if err != nil {
			return types.NOTHIN, fmt.Errorf("cannot cast return value to %s: %v", function.ReturnType, err)
		}
		result = castedResult
	}

	return result, nil
}

// callMemberFunction executes a member function
func (i *Interpreter) callMemberFunction(function *environment.Function, obj *environment.ObjectInstance, args []types.Value) (types.Value, error) {
	// Save current class and object context
	oldClass := i.currentClass
	oldObject := i.currentObject
	i.currentClass = obj.Class.Name
	i.currentObject = obj

	result, err := i.callFunction(function, args)

	// Restore context
	i.currentClass = oldClass
	i.currentObject = oldObject

	return result, err
}

// VisitMemberAccess handles member access
func (i *Interpreter) VisitMemberAccess(node *ast.MemberAccessNode) (types.Value, error) {
	objectValue, err := node.Object.Accept(i)
	if err != nil {
		return types.NOTHIN, err
	}

	if objVal, ok := objectValue.(types.ObjectValue); ok {
		if obj, ok := objVal.Instance.(*environment.ObjectInstance); ok {
			variable, err := obj.GetMemberVariable(strings.ToUpper(node.Member), i.currentClass)
			if err != nil {
				return types.NOTHIN, err
			}
			return variable.Value, nil
		}
	}

	return types.NOTHIN, fmt.Errorf("cannot access member of non-object")
}

// VisitBinaryOp handles binary operations
func (i *Interpreter) VisitBinaryOp(node *ast.BinaryOpNode) (types.Value, error) {
	left, err := node.Left.Accept(i)
	if err != nil {
		return types.NOTHIN, err
	}

	right, err := node.Right.Accept(i)
	if err != nil {
		return types.NOTHIN, err
	}

	switch strings.ToUpper(node.Operator) {
	case "MOAR": // Addition
		if leftNum, ok := left.(types.NumberValue); ok {
			if rightNum, ok := right.(types.NumberValue); ok {
				return leftNum.Add(rightNum), nil
			}
		}
		return types.NOTHIN, fmt.Errorf("operands must be numbers for addition")

	case "LES": // Subtraction
		if leftNum, ok := left.(types.NumberValue); ok {
			if rightNum, ok := right.(types.NumberValue); ok {
				return leftNum.Subtract(rightNum), nil
			}
		}
		return types.NOTHIN, fmt.Errorf("operands must be numbers for subtraction")

	case "TIEMZ": // Multiplication
		if leftNum, ok := left.(types.NumberValue); ok {
			if rightNum, ok := right.(types.NumberValue); ok {
				return leftNum.Multiply(rightNum), nil
			}
		}
		return types.NOTHIN, fmt.Errorf("operands must be numbers for multiplication")

	case "DIVIDEZ": // Division
		if leftNum, ok := left.(types.NumberValue); ok {
			if rightNum, ok := right.(types.NumberValue); ok {
				return leftNum.Divide(rightNum), nil
			}
		}
		return types.NOTHIN, fmt.Errorf("operands must be numbers for division")

	case "BIGGR THAN": // Greater than
		if leftNum, ok := left.(types.NumberValue); ok {
			if rightNum, ok := right.(types.NumberValue); ok {
				return leftNum.GreaterThan(rightNum), nil
			}
		}
		return types.NOTHIN, fmt.Errorf("operands must be numbers for comparison")

	case "SMALLR THAN": // Less than
		if leftNum, ok := left.(types.NumberValue); ok {
			if rightNum, ok := right.(types.NumberValue); ok {
				return leftNum.LessThan(rightNum), nil
			}
		}
		return types.NOTHIN, fmt.Errorf("operands must be numbers for comparison")

	case "SAEM AS": // Equal to
		return left.EqualTo(right)

	case "AN": // Logical AND
		return types.BoolValue(left.ToBool() == types.YEZ && right.ToBool() == types.YEZ), nil

	case "OR": // Logical OR
		return types.BoolValue(left.ToBool() == types.YEZ || right.ToBool() == types.YEZ), nil

	default:
		return types.NOTHIN, fmt.Errorf("unknown binary operator: %s", node.Operator)
	}
}

// VisitUnaryOp handles unary operations
func (i *Interpreter) VisitUnaryOp(node *ast.UnaryOpNode) (types.Value, error) {
	operand, err := node.Operand.Accept(i)
	if err != nil {
		return types.NOTHIN, err
	}

	switch strings.ToUpper(node.Operator) {
	case "NOT":
		return types.BoolValue(operand.ToBool() == types.NO), nil
	default:
		return types.NOTHIN, fmt.Errorf("unknown unary operator: %s", node.Operator)
	}
}

// VisitCast handles type casting
func (i *Interpreter) VisitCast(node *ast.CastNode) (types.Value, error) {
	value, err := node.Expression.Accept(i)
	if err != nil {
		return types.NOTHIN, err
	}

	return value.Cast(strings.ToUpper(node.TargetType))
}

// VisitLiteral handles literal values
func (i *Interpreter) VisitLiteral(node *ast.LiteralNode) (types.Value, error) {
	return node.Value, nil
}

// VisitIdentifier handles identifier references
func (i *Interpreter) VisitIdentifier(node *ast.IdentifierNode) (types.Value, error) {
	name := strings.ToUpper(node.Name)

	// If we're in a method call and the identifier might be a member variable, check current object first
	if i.currentObject != nil {
		if memberVar, err := i.currentObject.GetMemberVariable(name, i.currentClass); err == nil {
			return memberVar.Value, nil
		}
	}

	// Check local variables
	if variable, err := i.environment.GetVariable(name); err == nil {
		return variable.Value, nil
	}

	// Check if it's a zero-argument function call
	if function, err := i.environment.GetFunction(name); err == nil {
		return i.callFunction(function, []types.Value{})
	}

	// Check native functions
	if function, exists := i.runtime.GetNative(name); exists {
		return i.callFunction(function, []types.Value{})
	}

	return types.NOTHIN, fmt.Errorf("undefined variable or function '%s'", name)
}

// VisitObjectInstantiation handles object creation
func (i *Interpreter) VisitObjectInstantiation(node *ast.ObjectInstantiationNode) (types.Value, error) {
	var env *environment.Environment = i.environment

	// If source name is specified, get that environment
	if node.SourceName != "" {
		if sourceEnv, exists := i.runtime.GetSource(strings.ToUpper(node.SourceName)); exists {
			env = sourceEnv
		} else {
			return types.NOTHIN, fmt.Errorf("source '%s' not found", node.SourceName)
		}
	}

	instance, err := env.NewObjectInstance(strings.ToUpper(node.ClassName))
	if err != nil {
		return types.NOTHIN, err
	}
	return types.NewObjectValue(instance, strings.ToUpper(node.ClassName)), nil
}

// VisitStatementBlock handles statement blocks
func (i *Interpreter) VisitStatementBlock(node *ast.StatementBlockNode) (types.Value, error) {
	var result types.Value = types.NOTHIN

	for _, stmt := range node.Statements {
		val, err := stmt.Accept(i)
		if err != nil {
			return types.NOTHIN, err
		}
		result = val
	}

	return result, nil
}

// GetRuntime returns the runtime environment
func (i *Interpreter) GetRuntime() *environment.RuntimeEnvironment {
	return i.runtime
}
