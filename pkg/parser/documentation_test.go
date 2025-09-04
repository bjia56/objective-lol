package parser

import (
	"testing"
	
	"github.com/bjia56/objective-lol/pkg/ast"
)

func TestFunctionDocumentationParsing(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		expectedFuncName string
		expectedDocs     []string
	}{
		{
			name: "Single line documentation",
			input: `BTW This function adds two numbers
HAI ME TEH FUNCSHUN ADD TEH INTEGR WIT X TEH INTEGR AN WIT Y TEH INTEGR
    GIVEZ X MOAR Y
KTHXBAI`,
			expectedFuncName: "ADD",
			expectedDocs:     []string{"This function adds two numbers"},
		},
		{
			name: "Multi-line documentation",
			input: `BTW This function calculates factorial
BTW @param n The number to calculate factorial for
BTW @return The factorial result
HAI ME TEH FUNCSHUN FACTORIAL TEH INTEGR WIT N TEH INTEGR
    IZ N SMALLR THAN 2?
        GIVEZ 1
    NOPE
        GIVEZ N TIEMZ FACTORIAL WIT N LES 1
    KTHX
KTHXBAI`,
			expectedFuncName: "FACTORIAL",
			expectedDocs: []string{
				"This function calculates factorial",
				"@param n The number to calculate factorial for",
				"@return The factorial result",
			},
		},
		{
			name: "No documentation",
			input: `HAI ME TEH FUNCSHUN SIMPLE
    SAYZ WIT "Hello"
KTHXBAI`,
			expectedFuncName: "SIMPLE",
			expectedDocs:     nil,
		},
		{
			name: "Documentation with empty lines",
			input: `BTW First line of docs
BTW 
BTW Third line after empty comment
HAI ME TEH FUNCSHUN WITH_EMPTY TEH STRIN
    GIVEZ "test"
KTHXBAI`,
			expectedFuncName: "WITH_EMPTY",
			expectedDocs: []string{
				"First line of docs",
				"",
				"Third line after empty comment",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			parser := NewParser(lexer)
			program := parser.ParseProgram()

			if len(parser.Errors()) > 0 {
				t.Fatalf("Parser errors: %v", parser.Errors())
			}

			if len(program.Declarations) != 1 {
				t.Fatalf("Expected 1 declaration, got %d", len(program.Declarations))
			}

			funcDecl, ok := program.Declarations[0].(*ast.FunctionDeclarationNode)
			if !ok {
				t.Fatalf("Expected FunctionDeclarationNode, got %T", program.Declarations[0])
			}

			if funcDecl.Name != tt.expectedFuncName {
				t.Errorf("Expected function name %s, got %s", tt.expectedFuncName, funcDecl.Name)
			}

			if len(funcDecl.Documentation) != len(tt.expectedDocs) {
				t.Errorf("Expected %d documentation lines, got %d", len(tt.expectedDocs), len(funcDecl.Documentation))
				t.Errorf("Got documentation: %v", funcDecl.Documentation)
			}

			for i, expectedDoc := range tt.expectedDocs {
				if i >= len(funcDecl.Documentation) {
					t.Errorf("Missing documentation line %d: expected %q", i, expectedDoc)
					continue
				}
				if funcDecl.Documentation[i] != expectedDoc {
					t.Errorf("Documentation line %d: expected %q, got %q", i, expectedDoc, funcDecl.Documentation[i])
				}
			}
		})
	}
}

func TestMultipleFunctionDocumentation(t *testing.T) {
	input := `BTW First function documentation
HAI ME TEH FUNCSHUN FIRST
    SAYZ WIT "First"
KTHXBAI

BTW Second function documentation
BTW This one has multiple lines
HAI ME TEH FUNCSHUN SECOND TEH STRIN
    GIVEZ "Second"
KTHXBAI

HAI ME TEH FUNCSHUN THIRD
    SAYZ WIT "No docs"
KTHXBAI`

	lexer := NewLexer(input)
	parser := NewParser(lexer)
	program := parser.ParseProgram()

	if len(parser.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", parser.Errors())
	}

	if len(program.Declarations) != 3 {
		t.Fatalf("Expected 3 declarations, got %d", len(program.Declarations))
	}

	// Check first function
	func1, ok := program.Declarations[0].(*ast.FunctionDeclarationNode)
	if !ok {
		t.Fatalf("Expected FunctionDeclarationNode, got %T", program.Declarations[0])
	}
	if func1.Name != "FIRST" {
		t.Errorf("Expected function name FIRST, got %s", func1.Name)
	}
	if len(func1.Documentation) != 1 || func1.Documentation[0] != "First function documentation" {
		t.Errorf("Expected documentation [\"First function documentation\"], got %v", func1.Documentation)
	}

	// Check second function
	func2, ok := program.Declarations[1].(*ast.FunctionDeclarationNode)
	if !ok {
		t.Fatalf("Expected FunctionDeclarationNode, got %T", program.Declarations[1])
	}
	if func2.Name != "SECOND" {
		t.Errorf("Expected function name SECOND, got %s", func2.Name)
	}
	expectedDocs := []string{"Second function documentation", "This one has multiple lines"}
	if len(func2.Documentation) != len(expectedDocs) {
		t.Errorf("Expected %d documentation lines, got %d", len(expectedDocs), len(func2.Documentation))
	}
	for i, expected := range expectedDocs {
		if i >= len(func2.Documentation) || func2.Documentation[i] != expected {
			t.Errorf("Documentation line %d: expected %q, got %q", i, expected, func2.Documentation[i])
		}
	}

	// Check third function (no documentation)
	func3, ok := program.Declarations[2].(*ast.FunctionDeclarationNode)
	if !ok {
		t.Fatalf("Expected FunctionDeclarationNode, got %T", program.Declarations[2])
	}
	if func3.Name != "THIRD" {
		t.Errorf("Expected function name THIRD, got %s", func3.Name)
	}
	if len(func3.Documentation) != 0 {
		t.Errorf("Expected no documentation, got %v", func3.Documentation)
	}
}

func TestClassDocumentationParsing(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		expectedClassName string
		expectedDocs     []string
	}{
		{
			name: "Single line class documentation",
			input: `BTW This is a simple class for demonstration
HAI ME TEH CLAS PERSON
    EVRYONE
    DIS TEH VARIABLE NAME TEH STRIN ITZ "Unknown"
KTHXBAI`,
			expectedClassName: "PERSON",
			expectedDocs:     []string{"This is a simple class for demonstration"},
		},
		{
			name: "Multi-line class documentation",
			input: `BTW This class represents a bank account
BTW It manages balance and transaction operations
BTW @author Development Team
BTW @version 1.0
HAI ME TEH CLAS BANK_ACCOUNT
    EVRYONE
    DIS TEH VARIABLE BALANCE TEH DUBBLE ITZ 0.0
KTHXBAI`,
			expectedClassName: "BANK_ACCOUNT",
			expectedDocs: []string{
				"This class represents a bank account",
				"It manages balance and transaction operations",
				"@author Development Team",
				"@version 1.0",
			},
		},
		{
			name: "Class with no documentation",
			input: `HAI ME TEH CLAS SIMPLE
    EVRYONE
    DIS TEH VARIABLE VALUE TEH INTEGR ITZ 0
KTHXBAI`,
			expectedClassName: "SIMPLE",
			expectedDocs:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			parser := NewParser(lexer)
			program := parser.ParseProgram()

			if len(parser.Errors()) > 0 {
				t.Fatalf("Parser errors: %v", parser.Errors())
			}

			if len(program.Declarations) != 1 {
				t.Fatalf("Expected 1 declaration, got %d", len(program.Declarations))
			}

			classDecl, ok := program.Declarations[0].(*ast.ClassDeclarationNode)
			if !ok {
				t.Fatalf("Expected ClassDeclarationNode, got %T", program.Declarations[0])
			}

			if classDecl.Name != tt.expectedClassName {
				t.Errorf("Expected class name %s, got %s", tt.expectedClassName, classDecl.Name)
			}

			if len(classDecl.Documentation) != len(tt.expectedDocs) {
				t.Errorf("Expected %d documentation lines, got %d", len(tt.expectedDocs), len(classDecl.Documentation))
				t.Errorf("Got documentation: %v", classDecl.Documentation)
			}

			for i, expectedDoc := range tt.expectedDocs {
				if i >= len(classDecl.Documentation) {
					t.Errorf("Missing documentation line %d: expected %q", i, expectedDoc)
					continue
				}
				if classDecl.Documentation[i] != expectedDoc {
					t.Errorf("Documentation line %d: expected %q, got %q", i, expectedDoc, classDecl.Documentation[i])
				}
			}
		})
	}
}

func TestClassMethodDocumentation(t *testing.T) {
	input := `BTW This class demonstrates method documentation
HAI ME TEH CLAS CALCULATOR
    EVRYONE
    DIS TEH VARIABLE RESULT TEH DUBBLE ITZ 0.0

    BTW This method adds two numbers
    BTW @param x The first number
    BTW @param y The second number
    BTW @return The sum of x and y
    DIS TEH FUNCSHUN ADD TEH DUBBLE WIT X TEH DUBBLE AN WIT Y TEH DUBBLE
        GIVEZ X MOAR Y
    KTHX

    BTW This method has no parameters
    DIS TEH FUNCSHUN RESET
        RESULT ITZ 0.0
    KTHX

    DIS TEH FUNCSHUN UNDOCUMENTED TEH DUBBLE
        GIVEZ RESULT
    KTHX
KTHXBAI`

	lexer := NewLexer(input)
	parser := NewParser(lexer)
	program := parser.ParseProgram()

	if len(parser.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", parser.Errors())
	}

	if len(program.Declarations) != 1 {
		t.Fatalf("Expected 1 declaration, got %d", len(program.Declarations))
	}

	classDecl, ok := program.Declarations[0].(*ast.ClassDeclarationNode)
	if !ok {
		t.Fatalf("Expected ClassDeclarationNode, got %T", program.Declarations[0])
	}

	// Check class documentation
	expectedClassDocs := []string{"This class demonstrates method documentation"}
	if len(classDecl.Documentation) != len(expectedClassDocs) {
		t.Errorf("Expected %d class documentation lines, got %d", len(expectedClassDocs), len(classDecl.Documentation))
	}
	for i, expectedDoc := range expectedClassDocs {
		if i >= len(classDecl.Documentation) || classDecl.Documentation[i] != expectedDoc {
			t.Errorf("Class documentation line %d: expected %q, got %q", i, expectedDoc, classDecl.Documentation[i])
		}
	}

	// Check method documentation
	if len(classDecl.Members) != 4 { // 1 variable + 3 methods
		t.Fatalf("Expected 4 class members, got %d", len(classDecl.Members))
	}

	// Check ADD method (should have documentation)
	addMethod := classDecl.Members[1] // Skip the variable member
	if addMethod.IsVariable {
		t.Fatalf("Expected method, got variable")
	}
	expectedAddDocs := []string{
		"This method adds two numbers",
		"@param x The first number", 
		"@param y The second number",
		"@return The sum of x and y",
	}
	if len(addMethod.Function.Documentation) != len(expectedAddDocs) {
		t.Errorf("Expected %d method documentation lines, got %d", len(expectedAddDocs), len(addMethod.Function.Documentation))
	}
	for i, expectedDoc := range expectedAddDocs {
		if i >= len(addMethod.Function.Documentation) || addMethod.Function.Documentation[i] != expectedDoc {
			t.Errorf("ADD method documentation line %d: expected %q, got %q", i, expectedDoc, addMethod.Function.Documentation[i])
		}
	}

	// Check RESET method (should have simple documentation)
	resetMethod := classDecl.Members[2]
	if resetMethod.IsVariable {
		t.Fatalf("Expected method, got variable")
	}
	expectedResetDocs := []string{"This method has no parameters"}
	if len(resetMethod.Function.Documentation) != len(expectedResetDocs) {
		t.Errorf("Expected %d RESET documentation lines, got %d", len(expectedResetDocs), len(resetMethod.Function.Documentation))
	}
	if len(resetMethod.Function.Documentation) > 0 && resetMethod.Function.Documentation[0] != expectedResetDocs[0] {
		t.Errorf("RESET method documentation: expected %q, got %q", expectedResetDocs[0], resetMethod.Function.Documentation[0])
	}

	// Check UNDOCUMENTED method (should have no documentation)
	undocumentedMethod := classDecl.Members[3]
	if undocumentedMethod.IsVariable {
		t.Fatalf("Expected method, got variable")
	}
	if len(undocumentedMethod.Function.Documentation) != 0 {
		t.Errorf("Expected no documentation for UNDOCUMENTED method, got %v", undocumentedMethod.Function.Documentation)
	}
}

func TestVariableDocumentationParsing(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedDocs []string
		expectedName string
	}{
		{
			name: "Global variable with single line documentation",
			input: `BTW This variable stores the application version
HAI ME TEH VARIABLE APP_VERSION TEH STRIN ITZ "1.0"`,
			expectedDocs: []string{"This variable stores the application version"},
			expectedName: "APP_VERSION",
		},
		{
			name: "Global variable with multi-line documentation",
			input: `BTW This variable stores the maximum retry count
BTW It is used to limit connection attempts
BTW @default 5
BTW @type INTEGR
HAI ME TEH VARIABLE MAX_RETRIES TEH INTEGR ITZ 5`,
			expectedDocs: []string{
				"This variable stores the maximum retry count",
				"It is used to limit connection attempts",
				"@default 5",
				"@type INTEGR",
			},
			expectedName: "MAX_RETRIES",
		},
		{
			name: "Variable with no documentation",
			input: `HAI ME TEH VARIABLE SIMPLE TEH STRIN ITZ "test"`,
			expectedDocs: nil,
			expectedName: "SIMPLE",
		},
		{
			name: "Locked variable with documentation",
			input: `BTW This constant should not be changed
HAI ME TEH LOCKD VARIABLE PI TEH DUBBLE ITZ 3.14159`,
			expectedDocs: []string{"This constant should not be changed"},
			expectedName: "PI",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			parser := NewParser(lexer)
			program := parser.ParseProgram()

			if len(parser.Errors()) > 0 {
				t.Fatalf("Parser errors: %v", parser.Errors())
			}

			if len(program.Declarations) != 1 {
				t.Fatalf("Expected 1 declaration, got %d", len(program.Declarations))
			}

			varDecl, ok := program.Declarations[0].(*ast.VariableDeclarationNode)
			if !ok {
				t.Fatalf("Expected VariableDeclarationNode, got %T", program.Declarations[0])
			}

			if varDecl.Name != tt.expectedName {
				t.Errorf("Expected variable name %s, got %s", tt.expectedName, varDecl.Name)
			}

			if len(varDecl.Documentation) != len(tt.expectedDocs) {
				t.Errorf("Expected %d documentation lines, got %d", len(tt.expectedDocs), len(varDecl.Documentation))
				t.Errorf("Got documentation: %v", varDecl.Documentation)
			}

			for i, expectedDoc := range tt.expectedDocs {
				if i >= len(varDecl.Documentation) {
					t.Errorf("Missing documentation line %d: expected %q", i, expectedDoc)
					continue
				}
				if varDecl.Documentation[i] != expectedDoc {
					t.Errorf("Documentation line %d: expected %q, got %q", i, expectedDoc, varDecl.Documentation[i])
				}
			}
		})
	}
}

func TestClassMemberVariableDocumentation(t *testing.T) {
	input := `BTW This class demonstrates variable documentation
HAI ME TEH CLAS DATA_CONTAINER
    EVRYONE
    
    BTW Public identifier for the container
    BTW @type STRIN
    DIS TEH VARIABLE ID TEH STRIN ITZ "default"

    MAHSELF
    BTW Private data storage
    BTW Contains the main data payload
    BTW @access private
    DIS TEH VARIABLE DATA TEH STRIN ITZ ""

    EVRYONE
    BTW Counter without initialization
    DIS TEH VARIABLE COUNT TEH INTEGR

    DIS TEH VARIABLE UNDOCUMENTED TEH BOOL ITZ NO
KTHXBAI`

	lexer := NewLexer(input)
	parser := NewParser(lexer)
	program := parser.ParseProgram()

	if len(parser.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", parser.Errors())
	}

	if len(program.Declarations) != 1 {
		t.Fatalf("Expected 1 declaration, got %d", len(program.Declarations))
	}

	classDecl, ok := program.Declarations[0].(*ast.ClassDeclarationNode)
	if !ok {
		t.Fatalf("Expected ClassDeclarationNode, got %T", program.Declarations[0])
	}

	// Check we have 4 member variables
	if len(classDecl.Members) != 4 {
		t.Fatalf("Expected 4 class members, got %d", len(classDecl.Members))
	}

	// Check ID variable (should have documentation)
	idVar := classDecl.Members[0]
	if !idVar.IsVariable {
		t.Fatalf("Expected variable member, got method")
	}
	expectedIdDocs := []string{"Public identifier for the container", "@type STRIN"}
	if len(idVar.Variable.Documentation) != len(expectedIdDocs) {
		t.Errorf("Expected %d ID documentation lines, got %d", len(expectedIdDocs), len(idVar.Variable.Documentation))
	}
	for i, expectedDoc := range expectedIdDocs {
		if i >= len(idVar.Variable.Documentation) || idVar.Variable.Documentation[i] != expectedDoc {
			t.Errorf("ID variable documentation line %d: expected %q, got %q", i, expectedDoc, idVar.Variable.Documentation[i])
		}
	}

	// Check DATA variable (should have multi-line documentation)
	dataVar := classDecl.Members[1]
	if !dataVar.IsVariable {
		t.Fatalf("Expected variable member, got method")
	}
	expectedDataDocs := []string{"Private data storage", "Contains the main data payload", "@access private"}
	if len(dataVar.Variable.Documentation) != len(expectedDataDocs) {
		t.Errorf("Expected %d DATA documentation lines, got %d", len(expectedDataDocs), len(dataVar.Variable.Documentation))
	}
	for i, expectedDoc := range expectedDataDocs {
		if i >= len(dataVar.Variable.Documentation) || dataVar.Variable.Documentation[i] != expectedDoc {
			t.Errorf("DATA variable documentation line %d: expected %q, got %q", i, expectedDoc, dataVar.Variable.Documentation[i])
		}
	}

	// Check COUNT variable (should have simple documentation)
	countVar := classDecl.Members[2]
	if !countVar.IsVariable {
		t.Fatalf("Expected variable member, got method")
	}
	expectedCountDocs := []string{"Counter without initialization"}
	if len(countVar.Variable.Documentation) != len(expectedCountDocs) {
		t.Errorf("Expected %d COUNT documentation lines, got %d", len(expectedCountDocs), len(countVar.Variable.Documentation))
	}
	if len(countVar.Variable.Documentation) > 0 && countVar.Variable.Documentation[0] != expectedCountDocs[0] {
		t.Errorf("COUNT variable documentation: expected %q, got %q", expectedCountDocs[0], countVar.Variable.Documentation[0])
	}

	// Check UNDOCUMENTED variable (should have no documentation)
	undocVar := classDecl.Members[3]
	if !undocVar.IsVariable {
		t.Fatalf("Expected variable member, got method")
	}
	if len(undocVar.Variable.Documentation) != 0 {
		t.Errorf("Expected no documentation for UNDOCUMENTED variable, got %v", undocVar.Variable.Documentation)
	}
}

func TestLocalVariableDocumentation(t *testing.T) {
	input := `HAI ME TEH FUNCSHUN TEST
    BTW Counter for loop iterations
    BTW Keeps track of how many loops we've done
    I HAS A VARIABLE LOOP_COUNT TEH INTEGR ITZ 0

    BTW Temporary storage for results
    I HAS A VARIABLE TEMP_RESULT TEH STRIN ITZ ""

    I HAS A VARIABLE UNDOCUMENTED_LOCAL TEH BOOL ITZ YEZ
KTHXBAI`

	lexer := NewLexer(input)
	parser := NewParser(lexer)
	program := parser.ParseProgram()

	if len(parser.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", parser.Errors())
	}

	if len(program.Declarations) != 1 {
		t.Fatalf("Expected 1 declaration, got %d", len(program.Declarations))
	}

	funcDecl, ok := program.Declarations[0].(*ast.FunctionDeclarationNode)
	if !ok {
		t.Fatalf("Expected FunctionDeclarationNode, got %T", program.Declarations[0])
	}

	if len(funcDecl.Body.Statements) != 3 {
		t.Fatalf("Expected 3 statements in function body, got %d", len(funcDecl.Body.Statements))
	}

	// Check first variable (LOOP_COUNT - should have multi-line documentation)
	loopCountVar, ok := funcDecl.Body.Statements[0].(*ast.VariableDeclarationNode)
	if !ok {
		t.Fatalf("Expected first statement to be VariableDeclarationNode, got %T", funcDecl.Body.Statements[0])
	}
	if loopCountVar.Name != "LOOP_COUNT" {
		t.Errorf("Expected variable name LOOP_COUNT, got %s", loopCountVar.Name)
	}
	expectedLoopCountDocs := []string{"Counter for loop iterations", "Keeps track of how many loops we've done"}
	if len(loopCountVar.Documentation) != len(expectedLoopCountDocs) {
		t.Errorf("Expected %d documentation lines for LOOP_COUNT, got %d", len(expectedLoopCountDocs), len(loopCountVar.Documentation))
	}
	for i, expectedDoc := range expectedLoopCountDocs {
		if i >= len(loopCountVar.Documentation) || loopCountVar.Documentation[i] != expectedDoc {
			t.Errorf("LOOP_COUNT documentation line %d: expected %q, got %q", i, expectedDoc, loopCountVar.Documentation[i])
		}
	}

	// Check second variable (TEMP_RESULT - should have single line documentation)
	tempVar, ok := funcDecl.Body.Statements[1].(*ast.VariableDeclarationNode)
	if !ok {
		t.Fatalf("Expected second statement to be VariableDeclarationNode, got %T", funcDecl.Body.Statements[1])
	}
	if tempVar.Name != "TEMP_RESULT" {
		t.Errorf("Expected variable name TEMP_RESULT, got %s", tempVar.Name)
	}
	expectedTempDocs := []string{"Temporary storage for results"}
	if len(tempVar.Documentation) != len(expectedTempDocs) {
		t.Errorf("Expected %d documentation lines for TEMP_RESULT, got %d", len(expectedTempDocs), len(tempVar.Documentation))
	}
	if len(tempVar.Documentation) > 0 && tempVar.Documentation[0] != expectedTempDocs[0] {
		t.Errorf("TEMP_RESULT documentation: expected %q, got %q", expectedTempDocs[0], tempVar.Documentation[0])
	}

	// Check third variable (UNDOCUMENTED_LOCAL - should have no documentation)
	undocVar, ok := funcDecl.Body.Statements[2].(*ast.VariableDeclarationNode)
	if !ok {
		t.Fatalf("Expected third statement to be VariableDeclarationNode, got %T", funcDecl.Body.Statements[2])
	}
	if undocVar.Name != "UNDOCUMENTED_LOCAL" {
		t.Errorf("Expected variable name UNDOCUMENTED_LOCAL, got %s", undocVar.Name)
	}
	if len(undocVar.Documentation) != 0 {
		t.Errorf("Expected no documentation for UNDOCUMENTED_LOCAL, got %v", undocVar.Documentation)
	}
}