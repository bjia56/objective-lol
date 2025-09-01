package integration

import (
	"bytes"
	"strings"
	"testing"

	"github.com/bjia56/objective-lol/pkg/interpreter"
	"github.com/bjia56/objective-lol/pkg/parser"
	"github.com/bjia56/objective-lol/pkg/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// executeProgram parses and executes a complete Objective-LOL program
func executeProgram(t *testing.T, program string) error {
	lexer := parser.NewLexer(program)
	p := parser.NewParser(lexer)
	ast := p.ParseProgram()

	// Check for parsing errors
	if errors := p.Errors(); len(errors) > 0 {
		t.Fatalf("Parser errors: %v", errors)
	}

	// Create interpreter - stdlib is already registered in NewInterpreter()
	interp := interpreter.NewInterpreter()

	// Execute the program
	return interp.Interpret(ast)
}

// executeProgramWithOutput parses and executes a program, capturing output
func executeProgramWithOutput(t *testing.T, program string) (string, error) {
	// Set up output capture
	var buf bytes.Buffer
	stdlib.SetOutput(&buf)
	defer stdlib.ResetToStandardStreams()

	err := executeProgram(t, program)
	return buf.String(), err
}

// executeProgramWithInputOutput parses and executes a program with input and captures output
func executeProgramWithInputOutput(t *testing.T, program string, input string) (string, error) {
	// Set up I/O capture
	var buf bytes.Buffer
	stdlib.SetOutput(&buf)
	stdlib.SetInput(strings.NewReader(input))
	defer stdlib.ResetToStandardStreams()

	err := executeProgram(t, program)
	return buf.String(), err
}

func TestBasicVariablesAndArithmetic(t *testing.T) {
	program := `
	BTW Basic variables and arithmetic test
	I CAN HAS TEST?
	
	HAI ME TEH FUNCSHUN MAIN
		I HAS A VARIABLE X TEH INTEGR ITZ 10
		I HAS A VARIABLE Y TEH INTEGR ITZ 20
		I HAS A VARIABLE RESULT TEH INTEGR ITZ X MOAR Y
		
		BTW Test floating point
		I HAS A VARIABLE PI TEH DUBBLE ITZ 3.14159
		I HAS A VARIABLE RADIUS TEH DUBBLE ITZ 5.0
		I HAS A VARIABLE AREA TEH DUBBLE ITZ PI TIEMZ RADIUS TIEMZ RADIUS
		
		BTW Verify calculations using assertions
		ASSERT WIT RESULT SAEM AS 30
		ASSERT WIT AREA BIGGR THAN 78.0
		ASSERT WIT AREA SMALLR THAN 79.0
	KTHXBAI
	`

	err := executeProgram(t, program)
	require.NoError(t, err, "Program should execute successfully with correct arithmetic")
}

func TestInputOutputOperations(t *testing.T) {
	program := `
	BTW Test input/output operations
	I CAN HAS STDIO?
	
	HAI ME TEH FUNCSHUN MAIN
		SAYZ WIT "Enter your name:"
		I HAS A VARIABLE NAME TEH STRIN ITZ GIMME
		SAY WIT "Hello, "
		SAYZ WIT NAME
		SAY WIT "You entered: "
		SAYZ WIT NAME
	KTHXBAI
	`

	input := "Alice\n"
	output, err := executeProgramWithInputOutput(t, program, input)
	require.NoError(t, err, "I/O program should execute successfully")

	expectedOutput := "Enter your name:\nHello, Alice\nYou entered: Alice\n"
	assert.Equal(t, expectedOutput, output, "Output should match expected I/O pattern")
}

func TestPrintingAndFormatting(t *testing.T) {
	program := `
	BTW Test various printing scenarios
	I CAN HAS STDIO?
	
	HAI ME TEH FUNCSHUN MAIN
		SAYZ WIT "=== Testing SAY and SAYZ ==="
		SAY WIT "This is on "
		SAY WIT "the same line. "
		SAYZ WIT "This starts a new line."
		
		SAYZ WIT "Numbers:"
		SAY WIT 42
		SAY WIT " and "
		SAYZ WIT 3.14159
		
		SAYZ WIT "Boolean values:"
		SAYZ WIT YEZ
		SAYZ WIT NO
	KTHXBAI
	`

	output, err := executeProgramWithOutput(t, program)
	require.NoError(t, err, "Print program should execute successfully")

	expectedLines := []string{
		"=== Testing SAY and SAYZ ===",
		"This is on the same line. This starts a new line.",
		"Numbers:",
		"42 and 3.14159",
		"Boolean values:",
		"YEZ",
		"NO",
	}

	for _, expectedLine := range expectedLines {
		assert.Contains(t, output, expectedLine, "Output should contain expected line")
	}
}

func TestMathStandardLibrary(t *testing.T) {
	program := `
	BTW Math standard library test
	I CAN HAS MATH?
	I CAN HAS STDIO?
	I CAN HAS TEST?
	
	HAI ME TEH FUNCSHUN MAIN
		I HAS A VARIABLE ABS_VAL TEH DUBBLE ITZ ABS WIT -42.5
		I HAS A VARIABLE MAX_VAL TEH DUBBLE ITZ MAX WIT 10.5 AN WIT 20.3
		I HAS A VARIABLE SQRT_VAL TEH DUBBLE ITZ SQRT WIT 16.0
		
		BTW Verify math results
		ASSERT WIT ABS_VAL SAEM AS 42.5
		ASSERT WIT MAX_VAL SAEM AS 20.3
		ASSERT WIT SQRT_VAL SAEM AS 4.0
		
		SAYZ WIT "Math library test passed!"
	KTHXBAI
	`

	output, err := executeProgramWithOutput(t, program)
	require.NoError(t, err, "Math program should execute successfully")
	assert.Contains(t, output, "Math library test passed!", "Should print success message")
}

func TestArrayOperations(t *testing.T) {
	program := `
	BTW Array operations test
	I CAN HAS STDIO?
	I CAN HAS TEST?
	
	HAI ME TEH FUNCSHUN MAIN
		I HAS A VARIABLE MY_ARRAY TEH BUKKIT ITZ NEW BUKKIT
		
		BTW Add some elements
		MY_ARRAY DO PUSH WIT 10
		MY_ARRAY DO PUSH WIT 20
		MY_ARRAY DO PUSH WIT 30
		
		BTW Test array operations
		I HAS A VARIABLE FIRST_ELEMENT TEH INTEGR ITZ MY_ARRAY DO AT WIT 0
		I HAS A VARIABLE ARRAY_SIZE TEH INTEGR ITZ MY_ARRAY SIZ
		
		ASSERT WIT FIRST_ELEMENT SAEM AS 10
		ASSERT WIT ARRAY_SIZE SAEM AS 3
		
		BTW Test contains
		I HAS A VARIABLE CONTAINS_20 TEH BOOL ITZ MY_ARRAY DO CONTAINS WIT 20
		ASSERT WIT CONTAINS_20 SAEM AS YEZ
		
		SAYZ WIT "Array operations test passed!"
	KTHXBAI
	`

	output, err := executeProgramWithOutput(t, program)
	require.NoError(t, err, "Array program should execute successfully")
	assert.Contains(t, output, "Array operations test passed!", "Should print success message")
}

func TestFunctionDefinitionAndCalling(t *testing.T) {
	program := `
	BTW Function definition and calling test
	I CAN HAS STDIO?
	I CAN HAS TEST?
	
	HAI ME TEH FUNCSHUN ADD_NUMBERS TEH INTEGR WIT X TEH INTEGR AN WIT Y TEH INTEGR
		GIVEZ X MOAR Y
	KTHXBAI
	
	HAI ME TEH FUNCSHUN FACTORIAL TEH INTEGR WIT N TEH INTEGR
		IZ N SMALLR THAN 2?
			GIVEZ 1
		KTHX
		GIVEZ N TIEMZ FACTORIAL WIT N LES 1
	KTHXBAI
	
	HAI ME TEH FUNCSHUN MAIN
		BTW Test the functions
		I HAS A VARIABLE SUM TEH INTEGR ITZ ADD_NUMBERS WIT 15 AN WIT 25
		I HAS A VARIABLE FACT5 TEH INTEGR ITZ FACTORIAL WIT 5
		
		ASSERT WIT SUM SAEM AS 40
		ASSERT WIT FACT5 SAEM AS 120
		
		SAYZ WIT "Function test passed!"
	KTHXBAI
	`

	output, err := executeProgramWithOutput(t, program)
	require.NoError(t, err, "Function program should execute successfully")
	assert.Contains(t, output, "Function test passed!", "Should print success message")
}

func TestClassDefinitionAndInheritance(t *testing.T) {
	program := `
	BTW Class definition and inheritance test
	I CAN HAS STDIO?
	I CAN HAS TEST?
	
	HAI ME TEH CLAS ANIMAL
		EVRYONE
		DIS TEH VARIABLE NAME TEH STRIN ITZ "Unknown"
		
		DIS TEH FUNCSHUN SET_NAME WIT NEW_NAME TEH STRIN
			NAME ITZ NEW_NAME
		KTHX
		
		DIS TEH FUNCSHUN GET_NAME TEH STRIN
			GIVEZ NAME
		KTHX
	KTHXBAI
	
	HAI ME TEH CLAS DOG KITTEH OF ANIMAL
		EVRYONE
		DIS TEH VARIABLE BREED TEH STRIN ITZ "Mixed"
		
		DIS TEH FUNCSHUN SET_BREED WIT NEW_BREED TEH STRIN
			BREED ITZ NEW_BREED
		KTHX
		
		DIS TEH FUNCSHUN GET_BREED TEH STRIN
			GIVEZ BREED
		KTHX
	KTHXBAI
	
	HAI ME TEH FUNCSHUN MAIN
		BTW Test the classes
		I HAS A VARIABLE MY_DOG TEH DOG ITZ NEW DOG
		MY_DOG DO SET_NAME WIT "Rex"
		MY_DOG DO SET_BREED WIT "German Shepherd"
		
		I HAS A VARIABLE DOG_NAME TEH STRIN ITZ MY_DOG DO GET_NAME
		I HAS A VARIABLE DOG_BREED TEH STRIN ITZ MY_DOG DO GET_BREED
		
		ASSERT WIT DOG_NAME SAEM AS "Rex"
		ASSERT WIT DOG_BREED SAEM AS "German Shepherd"
		
		SAYZ WIT "Class inheritance test passed!"
	KTHXBAI
	`

	output, err := executeProgramWithOutput(t, program)
	require.NoError(t, err, "Class program should execute successfully")
	assert.Contains(t, output, "Class inheritance test passed!", "Should print success message")
}

func TestControlFlow(t *testing.T) {
	program := `
	BTW Control flow test
	I CAN HAS STDIO?
	I CAN HAS TEST?
	
	HAI ME TEH FUNCSHUN MAIN
		I HAS A VARIABLE COUNTER TEH INTEGR ITZ 0
		I HAS A VARIABLE SUM TEH INTEGR ITZ 0
		
		BTW While loop test
		WHILE COUNTER SMALLR THAN 5
			SUM ITZ SUM MOAR COUNTER
			COUNTER ITZ COUNTER MOAR 1
		KTHX
		
		ASSERT WIT SUM SAEM AS 10
		ASSERT WIT COUNTER SAEM AS 5
		
		BTW Conditional test
		I HAS A VARIABLE MESSAGE TEH STRIN
		IZ SUM BIGGR THAN 5?
			MESSAGE ITZ "Large sum"
		NOPE
			MESSAGE ITZ "Small sum"
		KTHX
		
		ASSERT WIT MESSAGE SAEM AS "Large sum"
		
		SAYZ WIT "Control flow test passed!"
	KTHXBAI
	`

	output, err := executeProgramWithOutput(t, program)
	require.NoError(t, err, "Control flow program should execute successfully")
	assert.Contains(t, output, "Control flow test passed!", "Should print success message")
}

func TestExceptionHandling(t *testing.T) {
	program := `
	BTW Exception handling test
	I CAN HAS STDIO?
	I CAN HAS TEST?
	
	HAI ME TEH FUNCSHUN MAIN
		I HAS A VARIABLE EXCEPTION_CAUGHT TEH BOOL ITZ NO
		
		MAYB
			ASSERT WIT NO SAEM AS YEZ
			EXCEPTION_CAUGHT ITZ NO
		OOPSIE ERR
			EXCEPTION_CAUGHT ITZ YEZ
		KTHX
		
		ASSERT WIT EXCEPTION_CAUGHT SAEM AS YEZ
		
		SAYZ WIT "Exception handling test passed!"
	KTHXBAI
	`

	output, err := executeProgramWithOutput(t, program)
	require.NoError(t, err, "Exception program should execute successfully")
	assert.Contains(t, output, "Exception handling test passed!", "Should print success message")
}

func TestComplexProgram(t *testing.T) {
	program := `
	BTW Complex program - calculator with array history
	I CAN HAS MATH?
	I CAN HAS STDIO?
	I CAN HAS TEST?
	
	HAI ME TEH CLAS CALCULATOR
		EVRYONE
		DIS TEH VARIABLE HISTORY TEH BUKKIT
		DIS TEH VARIABLE CURRENT_VALUE TEH DUBBLE ITZ 0.0
		
		DIS TEH FUNCSHUN CALCULATOR
			HISTORY ITZ NEW BUKKIT
		KTHX
		
		DIS TEH FUNCSHUN ADD TEH DUBBLE WIT VALUE TEH DUBBLE
			CURRENT_VALUE ITZ CURRENT_VALUE MOAR VALUE
			HISTORY DO PUSH WIT VALUE
			GIVEZ CURRENT_VALUE
		KTHX
		
		DIS TEH FUNCSHUN MULTIPLY TEH DUBBLE WIT VALUE TEH DUBBLE
			CURRENT_VALUE ITZ CURRENT_VALUE TIEMZ VALUE
			HISTORY DO PUSH WIT VALUE
			GIVEZ CURRENT_VALUE
		KTHX
		
		DIS TEH FUNCSHUN GET_HISTORY_SIZE TEH INTEGR
			GIVEZ HISTORY SIZ
		KTHX
	KTHXBAI
	
	HAI ME TEH FUNCSHUN MAIN
		BTW Test the calculator: (5 + 3) * 2 = 16
		I HAS A VARIABLE CALC TEH CALCULATOR ITZ NEW CALCULATOR
		I HAS A VARIABLE RESULT1 TEH DUBBLE ITZ CALC DO ADD WIT 5.0
		I HAS A VARIABLE RESULT2 TEH DUBBLE ITZ CALC DO ADD WIT 3.0
		I HAS A VARIABLE FINAL_RESULT TEH DUBBLE ITZ CALC DO MULTIPLY WIT 2.0
		I HAS A VARIABLE HISTORY_SIZE TEH INTEGR ITZ CALC DO GET_HISTORY_SIZE
		
		ASSERT WIT RESULT1 SAEM AS 5.0
		ASSERT WIT RESULT2 SAEM AS 8.0
		ASSERT WIT FINAL_RESULT SAEM AS 16.0
		ASSERT WIT HISTORY_SIZE SAEM AS 3
		
		SAYZ WIT "Complex calculator test passed!"
	KTHXBAI
	`

	output, err := executeProgramWithOutput(t, program)
	require.NoError(t, err, "Complex program should execute successfully")
	assert.Contains(t, output, "Complex calculator test passed!", "Should print success message")
}

func TestOperatorPrecedenceAndParentheses(t *testing.T) {
	program := `
	BTW Test operator precedence and parentheses
	I CAN HAS STDIO?
	I CAN HAS TEST?
	
	HAI ME TEH FUNCSHUN MAIN
		I HAS A VARIABLE RESULT1 TEH INTEGR ITZ 2 MOAR 3 TIEMZ 4
		I HAS A VARIABLE RESULT2 TEH INTEGR ITZ (2 MOAR 3) TIEMZ 4
		I HAS A VARIABLE RESULT3 TEH DUBBLE ITZ 10.0 DIVIDEZ (2.0 MOAR 3.0)
		I HAS A VARIABLE RESULT4 TEH INTEGR ITZ ((2 MOAR 3) TIEMZ (4 MOAR 1))
		
		ASSERT WIT RESULT1 SAEM AS 14
		ASSERT WIT RESULT2 SAEM AS 20
		ASSERT WIT RESULT3 SAEM AS 2.0
		ASSERT WIT RESULT4 SAEM AS 25
		
		SAYZ WIT "Operator precedence test passed!"
	KTHXBAI
	`

	output, err := executeProgramWithOutput(t, program)
	require.NoError(t, err, "Precedence program should execute successfully")
	assert.Contains(t, output, "Operator precedence test passed!", "Should print success message")
}

func TestTypeCasting(t *testing.T) {
	program := `
	BTW Type casting test
	I CAN HAS STDIO?
	I CAN HAS TEST?
	
	HAI ME TEH FUNCSHUN MAIN
		I HAS A VARIABLE NUM TEH INTEGR ITZ 42
		I HAS A VARIABLE NUM_AS_STRING TEH STRIN ITZ NUM AS STRIN
		I HAS A VARIABLE NUM_AS_DOUBLE TEH DUBBLE ITZ NUM AS DUBBLE
		
		I HAS A VARIABLE STR TEH STRIN ITZ "123"
		I HAS A VARIABLE STR_AS_INT TEH INTEGR ITZ STR AS INTEGR
		
		I HAS A VARIABLE ZERO TEH INTEGR ITZ 0
		I HAS A VARIABLE ZERO_AS_BOOL TEH BOOL ITZ ZERO AS BOOL
		
		ASSERT WIT NUM_AS_STRING SAEM AS "42"
		ASSERT WIT NUM_AS_DOUBLE SAEM AS 42.0
		ASSERT WIT STR_AS_INT SAEM AS 123
		ASSERT WIT ZERO_AS_BOOL SAEM AS NO
		
		SAYZ WIT "Type casting test passed!"
	KTHXBAI
	`

	output, err := executeProgramWithOutput(t, program)
	require.NoError(t, err, "Type casting program should execute successfully")
	assert.Contains(t, output, "Type casting test passed!", "Should print success message")
}
