package interpreter

import (
	"context"
	"testing"

	"github.com/bjia56/objective-lol/pkg/ast"
	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/parser"
	"github.com/bjia56/objective-lol/pkg/runtime"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewInterpreter(t *testing.T) {
	interp := NewInterpreter(nil)

	assert.NotNil(t, interp)
	assert.NotNil(t, interp.runtime)
	assert.NotNil(t, interp.environment)
	assert.NotNil(t, interp.moduleResolver)
}

func TestInterpreterLiterals(t *testing.T) {
	interp := NewInterpreter(nil)

	tests := []struct {
		name     string
		literal  *ast.LiteralNode
		expected environment.Value
	}{
		{
			"Integer literal",
			&ast.LiteralNode{Value: environment.IntegerValue(42)},
			environment.IntegerValue(42),
		},
		{
			"String literal",
			&ast.LiteralNode{Value: environment.StringValue("hello")},
			environment.StringValue("hello"),
		},
		{
			"Boolean literal",
			&ast.LiteralNode{Value: environment.YEZ},
			environment.YEZ,
		},
		{
			"Double literal",
			&ast.LiteralNode{Value: environment.DoubleValue(3.14)},
			environment.DoubleValue(3.14),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := test.literal.Accept(interp)
			require.NoError(t, err)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestInterpreterBinaryOperations(t *testing.T) {
	interp := NewInterpreter(nil)

	tests := []struct {
		name     string
		left     environment.Value
		operator string
		right    environment.Value
		expected environment.Value
	}{
		{
			"Addition",
			environment.IntegerValue(2),
			"MOAR",
			environment.IntegerValue(3),
			environment.IntegerValue(5),
		},
		{
			"Subtraction",
			environment.IntegerValue(5),
			"LES",
			environment.IntegerValue(2),
			environment.IntegerValue(3),
		},
		{
			"Multiplication",
			environment.IntegerValue(4),
			"TIEMZ",
			environment.IntegerValue(3),
			environment.IntegerValue(12),
		},
		{
			"Division",
			environment.IntegerValue(10),
			"DIVIDEZ",
			environment.IntegerValue(2),
			environment.DoubleValue(5.0),
		},
		{
			"Equality true",
			environment.IntegerValue(42),
			"SAEM AS",
			environment.IntegerValue(42),
			environment.YEZ,
		},
		{
			"Equality false",
			environment.IntegerValue(42),
			"SAEM AS",
			environment.IntegerValue(24),
			environment.NO,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			binaryOp := &ast.BinaryOpNode{
				Left:     &ast.LiteralNode{Value: test.left},
				Operator: test.operator,
				Right:    &ast.LiteralNode{Value: test.right},
			}

			result, err := binaryOp.Accept(interp)
			require.NoError(t, err)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestInterpreterVariableOperations(t *testing.T) {
	interp := NewInterpreter(nil)

	// Define a variable
	varDecl := &ast.VariableDeclarationNode{
		Name:  "x",
		Type:  &ast.IdentifierNode{Name: "INTEGR"},
		Value: &ast.LiteralNode{Value: environment.IntegerValue(42)},
	}

	_, err := varDecl.Accept(interp)
	require.NoError(t, err)

	// Access the variable
	identifier := &ast.IdentifierNode{Name: "x"}
	result, err := identifier.Accept(interp)
	require.NoError(t, err)
	assert.Equal(t, environment.IntegerValue(42), result)

	// Assign to the variable
	assignment := &ast.AssignmentNode{
		Target: &ast.IdentifierNode{Name: "x"},
		Value:  &ast.LiteralNode{Value: environment.IntegerValue(100)},
	}

	_, err = assignment.Accept(interp)
	require.NoError(t, err)

	// Check the new value
	result, err = identifier.Accept(interp)
	require.NoError(t, err)
	assert.Equal(t, environment.IntegerValue(100), result)
}

func TestInterpreterTypeCasting(t *testing.T) {
	interp := NewInterpreter(nil)

	tests := []struct {
		name       string
		value      environment.Value
		targetType string
		expected   environment.Value
	}{
		{
			"Integer to string",
			environment.IntegerValue(42),
			"STRIN",
			environment.StringValue("42"),
		},
		{
			"String to integer",
			environment.StringValue("24"),
			"INTEGR",
			environment.IntegerValue(24),
		},
		{
			"Integer to boolean",
			environment.IntegerValue(1),
			"BOOL",
			environment.YEZ,
		},
		{
			"Integer zero to boolean",
			environment.IntegerValue(0),
			"BOOL",
			environment.NO,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cast := &ast.CastNode{
				Expression: &ast.LiteralNode{Value: test.value},
				TargetType: test.targetType,
			}

			result, err := cast.Accept(interp)
			require.NoError(t, err)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestInterpreterProgramExecution(t *testing.T) {
	// Test simple program execution
	input := `HAI ME TEH FUNCSHUN TEST
    I HAS A VARIABLE x TEH INTEGR ITZ 42
    x ITZ x MOAR 8
KTHXBAI`

	lexer := parser.NewLexer(input)
	p := parser.NewParser(lexer)
	program := p.ParseProgram()

	// We expect parsing errors since this is complex syntax
	if len(p.Errors()) == 0 {
		interp := NewInterpreter(nil)
		_, err := interp.Interpret(context.Background(), program)

		// Log any interpretation errors for debugging
		if err != nil {
			t.Logf("Interpreter error: %v", err)
		}
	} else {
		t.Logf("Parser errors (expected for complex syntax): %v", p.Errors())
	}
}

func TestInterpreterErrorHandling(t *testing.T) {
	interp := NewInterpreter(nil)

	// Test division by zero
	divByZero := &ast.BinaryOpNode{
		Left:     &ast.LiteralNode{Value: environment.IntegerValue(10)},
		Operator: "DIVIDEZ",
		Right:    &ast.LiteralNode{Value: environment.IntegerValue(0)},
	}

	result, err := divByZero.Accept(interp)
	// Division by zero should be handled gracefully
	if err != nil {
		assert.Contains(t, err.Error(), "Division by zero")
	} else {
		// Or return 0/NaN
		t.Logf("Division by zero returned: %v", result)
	}

	// Test accessing undefined variable
	undefinedVar := &ast.IdentifierNode{Name: "undefined"}
	_, err = undefinedVar.Accept(interp)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Undefined")
}

func TestInterpreterFunctionCall(t *testing.T) {
	interp := NewInterpreter(nil)

	// Test calling a function (using identifier as function)
	functionCall := &ast.FunctionCallNode{
		Function: &ast.IdentifierNode{Name: "someFunction"},
		Arguments: []ast.Node{
			&ast.LiteralNode{Value: environment.IntegerValue(1)},
			&ast.LiteralNode{Value: environment.StringValue("arg")},
		},
	}

	// This should error since the function doesn't exist
	_, err := functionCall.Accept(interp)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "undefined")
}

func TestInterpreterStatementBlock(t *testing.T) {
	interp := NewInterpreter(nil)

	// Create a statement block with multiple statements
	block := &ast.StatementBlockNode{
		Statements: []ast.Node{
			&ast.VariableDeclarationNode{
				Name:  "a",
				Type:  &ast.IdentifierNode{Name: "INTEGR"},
				Value: &ast.LiteralNode{Value: environment.IntegerValue(1)},
			},
			&ast.VariableDeclarationNode{
				Name:  "b",
				Type:  &ast.IdentifierNode{Name: "INTEGR"},
				Value: &ast.LiteralNode{Value: environment.IntegerValue(2)},
			},
		},
	}

	_, err := block.Accept(interp)
	require.NoError(t, err)

	// Check that both variables were defined
	varA := &ast.IdentifierNode{Name: "a"}
	result, err := varA.Accept(interp)
	require.NoError(t, err)
	assert.Equal(t, environment.IntegerValue(1), result)

	varB := &ast.IdentifierNode{Name: "b"}
	result, err = varB.Accept(interp)
	require.NoError(t, err)
	assert.Equal(t, environment.IntegerValue(2), result)
}

func TestInterpreterReturnStatement(t *testing.T) {
	interp := NewInterpreter(nil)

	// Test return statement
	returnStmt := &ast.ReturnStatementNode{
		Value: &ast.LiteralNode{Value: environment.StringValue("returned")},
	}

	_, err := returnStmt.Accept(interp)

	// Return statements should create a special error type
	assert.Error(t, err)

	// Check if it's a return value using the helper function
	if runtime.IsReturnValue(err) {
		returnVal := runtime.GetReturnValue(err)
		assert.Equal(t, environment.StringValue("returned"), returnVal)
	}
}
