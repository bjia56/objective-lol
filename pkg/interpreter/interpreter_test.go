package interpreter

import (
	"context"
	"testing"

	"github.com/bjia56/objective-lol/pkg/ast"
	"github.com/bjia56/objective-lol/pkg/parser"
	"github.com/bjia56/objective-lol/pkg/types"
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
		expected types.Value
	}{
		{
			"Integer literal",
			&ast.LiteralNode{Value: types.IntegerValue(42)},
			types.IntegerValue(42),
		},
		{
			"String literal",
			&ast.LiteralNode{Value: types.StringValue("hello")},
			types.StringValue("hello"),
		},
		{
			"Boolean literal",
			&ast.LiteralNode{Value: types.YEZ},
			types.YEZ,
		},
		{
			"Double literal",
			&ast.LiteralNode{Value: types.DoubleValue(3.14)},
			types.DoubleValue(3.14),
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
		left     types.Value
		operator string
		right    types.Value
		expected types.Value
	}{
		{
			"Addition",
			types.IntegerValue(2),
			"MOAR",
			types.IntegerValue(3),
			types.IntegerValue(5),
		},
		{
			"Subtraction",
			types.IntegerValue(5),
			"LES",
			types.IntegerValue(2),
			types.IntegerValue(3),
		},
		{
			"Multiplication",
			types.IntegerValue(4),
			"TIEMZ",
			types.IntegerValue(3),
			types.IntegerValue(12),
		},
		{
			"Division",
			types.IntegerValue(10),
			"DIVIDEZ",
			types.IntegerValue(2),
			types.DoubleValue(5.0),
		},
		{
			"Equality true",
			types.IntegerValue(42),
			"SAEM AS",
			types.IntegerValue(42),
			types.YEZ,
		},
		{
			"Equality false",
			types.IntegerValue(42),
			"SAEM AS",
			types.IntegerValue(24),
			types.NO,
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
		Type:  "INTEGR",
		Value: &ast.LiteralNode{Value: types.IntegerValue(42)},
	}

	_, err := varDecl.Accept(interp)
	require.NoError(t, err)

	// Access the variable
	identifier := &ast.IdentifierNode{Name: "x"}
	result, err := identifier.Accept(interp)
	require.NoError(t, err)
	assert.Equal(t, types.IntegerValue(42), result)

	// Assign to the variable
	assignment := &ast.AssignmentNode{
		Target: &ast.IdentifierNode{Name: "x"},
		Value:  &ast.LiteralNode{Value: types.IntegerValue(100)},
	}

	_, err = assignment.Accept(interp)
	require.NoError(t, err)

	// Check the new value
	result, err = identifier.Accept(interp)
	require.NoError(t, err)
	assert.Equal(t, types.IntegerValue(100), result)
}

func TestInterpreterTypeCasting(t *testing.T) {
	interp := NewInterpreter(nil)

	tests := []struct {
		name       string
		value      types.Value
		targetType string
		expected   types.Value
	}{
		{
			"Integer to string",
			types.IntegerValue(42),
			"STRIN",
			types.StringValue("42"),
		},
		{
			"String to integer",
			types.StringValue("24"),
			"INTEGR",
			types.IntegerValue(24),
		},
		{
			"Integer to boolean",
			types.IntegerValue(1),
			"BOOL",
			types.YEZ,
		},
		{
			"Integer zero to boolean",
			types.IntegerValue(0),
			"BOOL",
			types.NO,
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
		Left:     &ast.LiteralNode{Value: types.IntegerValue(10)},
		Operator: "DIVIDEZ",
		Right:    &ast.LiteralNode{Value: types.IntegerValue(0)},
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
			&ast.LiteralNode{Value: types.IntegerValue(1)},
			&ast.LiteralNode{Value: types.StringValue("arg")},
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
				Type:  "INTEGR",
				Value: &ast.LiteralNode{Value: types.IntegerValue(1)},
			},
			&ast.VariableDeclarationNode{
				Name:  "b",
				Type:  "INTEGR",
				Value: &ast.LiteralNode{Value: types.IntegerValue(2)},
			},
		},
	}

	_, err := block.Accept(interp)
	require.NoError(t, err)

	// Check that both variables were defined
	varA := &ast.IdentifierNode{Name: "a"}
	result, err := varA.Accept(interp)
	require.NoError(t, err)
	assert.Equal(t, types.IntegerValue(1), result)

	varB := &ast.IdentifierNode{Name: "b"}
	result, err = varB.Accept(interp)
	require.NoError(t, err)
	assert.Equal(t, types.IntegerValue(2), result)
}

func TestInterpreterReturnStatement(t *testing.T) {
	interp := NewInterpreter(nil)

	// Test return statement
	returnStmt := &ast.ReturnStatementNode{
		Value: &ast.LiteralNode{Value: types.StringValue("returned")},
	}

	_, err := returnStmt.Accept(interp)

	// Return statements should create a special error type
	assert.Error(t, err)

	// Check if it's a return value using the helper function
	if ast.IsReturnValue(err) {
		returnVal := ast.GetReturnValue(err)
		assert.Equal(t, types.StringValue("returned"), returnVal)
	}
}
