package analyzer

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

const comprehensiveExample = `BTW Class documentation test - testing comment parsing for classes and methods
I CAN HAS STDIO?

BTW This is a simple calculator class
BTW It provides basic arithmetic operations
BTW @author Development Team
BTW @version 1.0
HAI ME TEH CLAS CALCULATOR
    EVRYONE
    DIS TEH VARIABLE RESULT TEH DUBBLE ITZ 0.0

    BTW This method adds two numbers together
    BTW @param x The first number to add
    BTW @param y The second number to add
    BTW @return The sum of the two numbers
    DIS TEH FUNCSHUN ADD TEH DUBBLE WIT X TEH DUBBLE AN WIT Y TEH DUBBLE
        RESULT ITZ X MOAR Y
        GIVEZ RESULT
    KTHX

    BTW This method subtracts the second number from the first
    BTW @param x The number to subtract from
    BTW @param y The number to subtract
    BTW @return The difference between the numbers
    DIS TEH FUNCSHUN SUBTRACT TEH DUBBLE WIT X TEH DUBBLE AN WIT Y TEH DUBBLE
        RESULT ITZ X LES Y
        GIVEZ RESULT
    KTHX

    BTW Reset the calculator result to zero
    DIS TEH FUNCSHUN RESET
        RESULT ITZ 0.0
    KTHX

    DIS TEH FUNCSHUN GET_RESULT TEH DUBBLE
        GIVEZ RESULT
    KTHX
KTHXBAI

BTW This class represents a simple counter
HAI ME TEH CLAS COUNTER
    EVRYONE
    DIS TEH VARIABLE COUNT TEH INTEGR ITZ 0

    BTW Increment the counter by 1
    DIS TEH FUNCSHUN INCREMENT
        COUNT ITZ COUNT MOAR 1
    KTHX

    BTW Get the current count value
    BTW @return The current count
    DIS TEH FUNCSHUN GET_COUNT TEH INTEGR
        GIVEZ COUNT
    KTHX
KTHXBAI

HAI ME TEH CLAS UNDOCUMENTED_CLASS
    EVRYONE
    DIS TEH VARIABLE VALUE TEH STRIN ITZ "test"

    DIS TEH FUNCSHUN DO_SOMETHING
        SAYZ WIT "This class and method have no documentation"
    KTHX
KTHXBAI

BTW Main function to test the documented classes
HAI ME TEH FUNCSHUN MAIN
    SAYZ WIT "=== Class Documentation Test ==="

    BTW Test documented calculator class
    I HAS A VARIABLE CALC TEH CALCULATOR ITZ NEW CALCULATOR

    I HAS A VARIABLE SUM TEH DUBBLE ITZ CALC DO ADD WIT 10.5 AN WIT 5.2
    SAY WIT "10.5 + 5.2 = "
    SAYZ WIT SUM

    I HAS A VARIABLE DIFF TEH DUBBLE ITZ CALC DO SUBTRACT WIT 20.0 AN WIT 7.5
    SAY WIT "20.0 - 7.5 = "
    SAYZ WIT DIFF

    CALC DO RESET
    I HAS A VARIABLE RESET_RESULT TEH DUBBLE ITZ CALC DO GET_RESULT
    SAY WIT "After reset: "
    SAYZ WIT RESET_RESULT

    BTW Test documented counter class
    I HAS A VARIABLE COUNTER TEH COUNTER ITZ NEW COUNTER

    SAY WIT "Initial count: "
    SAYZ WIT COUNTER DO GET_COUNT

    COUNTER DO INCREMENT
    COUNTER DO INCREMENT
    COUNTER DO INCREMENT

    SAY WIT "After 3 increments: "
    SAYZ WIT COUNTER DO GET_COUNT

    BTW Test undocumented class
    I HAS A VARIABLE UNDOC TEH UNDOCUMENTED_CLASS ITZ NEW UNDOCUMENTED_CLASS
    UNDOC DO DO_SOMETHING

    SAYZ WIT "=== Class Documentation Test Complete ==="
KTHXBAI`

// Test helper functions
func findSymbolByName(symbols []EnhancedSymbol, name string) *EnhancedSymbol {
	for i := range symbols {
		if strings.EqualFold(symbols[i].Name, name) {
			return &symbols[i]
		}
	}
	return nil
}

func findSymbolByNameAndKind(symbols []EnhancedSymbol, name string, kind SymbolKind) *EnhancedSymbol {
	for i := range symbols {
		if strings.EqualFold(symbols[i].Name, name) && symbols[i].Kind == kind {
			return &symbols[i]
		}
	}
	return nil
}

func findSymbolReferenceByName(symbols []EnhancedSymbol, name string) *EnhancedSymbol {
	// Find a symbol that represents a reference (has its position in its References slice)
	for i := range symbols {
		if strings.EqualFold(symbols[i].Name, name) && len(symbols[i].References) > 0 {
			// Check if this symbol's position matches one of its references
			for _, ref := range symbols[i].References {
				if ref.Line == symbols[i].Position.Line && ref.Column == symbols[i].Position.Column {
					return &symbols[i]
				}
			}
		}
	}
	return nil
}

// Removed unused function countSymbolsByKind

func countSymbolsByScope(symbols []EnhancedSymbol, scopeType ScopeType) int {
	count := 0
	for _, symbol := range symbols {
		if symbol.Scope == scopeType {
			count++
		}
	}
	return count
}

// Test fixtures using comprehensiveExample
func TestComprehensiveAnalysis_BasicStructure(t *testing.T) {
	analyzer := NewSemanticAnalyzer("test://comprehensive.olol", comprehensiveExample)

	err := analyzer.AnalyzeDocument(context.Background())
	require.NoError(t, err, "AnalyzeDocument should succeed")

	symbolTable := analyzer.GetSymbolTable()
	require.NotNil(t, symbolTable, "Symbol table should not be nil")

	// Basic structure validation
	assert.Greater(t, len(symbolTable.Symbols), 0, "Should have symbols")
	assert.Greater(t, len(symbolTable.Scopes), 0, "Should have scopes")

	// Check for main function
	mainFunc := findSymbolByNameAndKind(symbolTable.Symbols, "MAIN", SymbolKindFunction)
	require.NotNil(t, mainFunc, "Should find MAIN function")
	assert.Equal(t, ScopeTypeGlobal, mainFunc.Scope, "MAIN should be in global scope")

	t.Logf("Found %d symbols, %d scopes", len(symbolTable.Symbols), len(symbolTable.Scopes))
}

func TestComprehensiveAnalysis_ClassAnalysis(t *testing.T) {
	analyzer := NewSemanticAnalyzer("test://comprehensive.olol", comprehensiveExample)

	err := analyzer.AnalyzeDocument(context.Background())
	require.NoError(t, err)

	symbolTable := analyzer.GetSymbolTable()

	// Test class detection
	expectedClasses := []string{"CALCULATOR", "COUNTER", "UNDOCUMENTED_CLASS"}
	for _, className := range expectedClasses {
		classSymbol := findSymbolByNameAndKind(symbolTable.Symbols, className, SymbolKindClass)
		require.NotNil(t, classSymbol, "Should find class %s", className)
		assert.Equal(t, ScopeTypeGlobal, classSymbol.Scope, "Class %s should be in global scope", className)
		assert.Equal(t, className, classSymbol.Type, "Class type should match name")
	}

	// Test documented vs undocumented classes
	calculatorClass := findSymbolByName(symbolTable.Symbols, "CALCULATOR")
	undocumentedClass := findSymbolByName(symbolTable.Symbols, "UNDOCUMENTED_CLASS")

	assert.NotEmpty(t, calculatorClass.Documentation, "CALCULATOR should have documentation")
	assert.Empty(t, undocumentedClass.Documentation, "UNDOCUMENTED_CLASS should have no documentation")

	// Count class members
	calculatorMembers := 0
	counterMembers := 0
	for _, symbol := range symbolTable.Symbols {
		if symbol.ParentClass == "CALCULATOR" {
			calculatorMembers++
		}
		if symbol.ParentClass == "COUNTER" {
			counterMembers++
		}
	}

	assert.Greater(t, calculatorMembers, 0, "CALCULATOR should have members")
	assert.Greater(t, counterMembers, 0, "COUNTER should have members")

	t.Logf("CALCULATOR has %d members, COUNTER has %d members", calculatorMembers, counterMembers)
}

func TestComprehensiveAnalysis_ClassMembers(t *testing.T) {
	analyzer := NewSemanticAnalyzer("test://comprehensive.olol", comprehensiveExample)

	err := analyzer.AnalyzeDocument(context.Background())
	require.NoError(t, err)

	symbolTable := analyzer.GetSymbolTable()

	// Test CALCULATOR class members
	expectedCalculatorMembers := map[string]SymbolKind{
		"RESULT":     SymbolKindVariable,
		"ADD":        SymbolKindFunction,
		"SUBTRACT":   SymbolKindFunction,
		"RESET":      SymbolKindFunction,
		"GET_RESULT": SymbolKindFunction,
	}

	for memberName, expectedKind := range expectedCalculatorMembers {
		member := findSymbolByNameAndKind(symbolTable.Symbols, memberName, expectedKind)
		require.NotNil(t, member, "Should find CALCULATOR member %s", memberName)
		assert.Equal(t, "CALCULATOR", member.ParentClass, "%s should belong to CALCULATOR", memberName)
		assert.Equal(t, ScopeTypeClass, member.Scope, "%s should be in class scope", memberName)
	}

	// Test COUNTER class members
	expectedCounterMembers := map[string]SymbolKind{
		"COUNT":     SymbolKindVariable,
		"INCREMENT": SymbolKindFunction,
		"GET_COUNT": SymbolKindFunction,
	}

	for memberName, expectedKind := range expectedCounterMembers {
		member := findSymbolByNameAndKind(symbolTable.Symbols, memberName, expectedKind)
		require.NotNil(t, member, "Should find COUNTER member %s", memberName)
		assert.Equal(t, "COUNTER", member.ParentClass, "%s should belong to COUNTER", memberName)
		assert.Equal(t, ScopeTypeClass, member.Scope, "%s should be in class scope", memberName)
	}

	// Test member visibility
	resultVar := findSymbolByName(symbolTable.Symbols, "RESULT")
	assert.Equal(t, VisibilityPublic, resultVar.Visibility, "RESULT should be public (EVRYONE)")

	countVar := findSymbolByName(symbolTable.Symbols, "COUNT")
	assert.Equal(t, VisibilityPublic, countVar.Visibility, "COUNT should be public (EVRYONE)")
}

func TestComprehensiveAnalysis_FunctionDocumentation(t *testing.T) {
	analyzer := NewSemanticAnalyzer("test://comprehensive.olol", comprehensiveExample)

	err := analyzer.AnalyzeDocument(context.Background())
	require.NoError(t, err)

	symbolTable := analyzer.GetSymbolTable()

	// Test documented functions
	addFunc := findSymbolByNameAndKind(symbolTable.Symbols, "ADD", SymbolKindFunction)
	require.NotNil(t, addFunc, "Should find ADD function")
	assert.NotEmpty(t, addFunc.Documentation, "ADD should have documentation")
	assert.Contains(t, addFunc.Documentation, "@param", "ADD documentation should contain @param")
	assert.Contains(t, addFunc.Documentation, "@return", "ADD documentation should contain @return")

	subtractFunc := findSymbolByNameAndKind(symbolTable.Symbols, "SUBTRACT", SymbolKindFunction)
	require.NotNil(t, subtractFunc, "Should find SUBTRACT function")
	assert.NotEmpty(t, subtractFunc.Documentation, "SUBTRACT should have documentation")

	// Test undocumented function
	doSomethingFunc := findSymbolByNameAndKind(symbolTable.Symbols, "DO_SOMETHING", SymbolKindFunction)
	require.NotNil(t, doSomethingFunc, "Should find DO_SOMETHING function")
	assert.Empty(t, doSomethingFunc.Documentation, "DO_SOMETHING should have no documentation")

	// Test function return types
	assert.Equal(t, "DUBBLE", addFunc.Type, "ADD should return DUBBLE")
	assert.Equal(t, "DUBBLE", subtractFunc.Type, "SUBTRACT should return DUBBLE")
	assert.Equal(t, "", doSomethingFunc.Type, "DO_SOMETHING should have no return type specified")
}

func TestComprehensiveAnalysis_VariableAnalysis(t *testing.T) {
	analyzer := NewSemanticAnalyzer("test://comprehensive.olol", comprehensiveExample)

	err := analyzer.AnalyzeDocument(context.Background())
	require.NoError(t, err)

	symbolTable := analyzer.GetSymbolTable()

	// Test main function variables
	mainVars := []struct {
		name         string
		expectedType string
		scope        ScopeType
	}{
		{"CALC", "CALCULATOR", ScopeTypeFunction},
		{"SUM", "DUBBLE", ScopeTypeFunction},
		{"DIFF", "DUBBLE", ScopeTypeFunction},
		{"RESET_RESULT", "DUBBLE", ScopeTypeFunction},
		{"COUNTER", "COUNTER", ScopeTypeFunction},
		{"UNDOC", "UNDOCUMENTED_CLASS", ScopeTypeFunction},
	}

	for _, varTest := range mainVars {
		varSymbol := findSymbolByNameAndKind(symbolTable.Symbols, varTest.name, SymbolKindVariable)
		if assert.NotNil(t, varSymbol, "Should find variable %s", varTest.name) {
			assert.Equal(t, varTest.expectedType, varSymbol.Type, "Variable %s should have type %s", varTest.name, varTest.expectedType)
			assert.Equal(t, varTest.scope, varSymbol.Scope, "Variable %s should be in %v scope", varTest.name, varTest.scope)
		}
	}

	// Test class member variables
	resultVar := findSymbolByName(symbolTable.Symbols, "RESULT")
	require.NotNil(t, resultVar, "Should find RESULT variable")
	assert.Equal(t, "DUBBLE", resultVar.Type, "RESULT should be DUBBLE")
	assert.Equal(t, "CALCULATOR", resultVar.ParentClass, "RESULT should belong to CALCULATOR")

	countVar := findSymbolByName(symbolTable.Symbols, "COUNT")
	require.NotNil(t, countVar, "Should find COUNT variable")
	assert.Equal(t, "INTEGR", countVar.Type, "COUNT should be INTEGR")
	assert.Equal(t, "COUNTER", countVar.ParentClass, "COUNT should belong to COUNTER")
}

func TestComprehensiveAnalysis_FunctionCallReferences(t *testing.T) {
	analyzer := NewSemanticAnalyzer("test://comprehensive.olol", comprehensiveExample)

	err := analyzer.AnalyzeDocument(context.Background())
	require.NoError(t, err)

	symbolTable := analyzer.GetSymbolTable()

	// Should have tracked function call identifiers as reference symbols
	referenceSymbols := []EnhancedSymbol{}
	for _, symbol := range symbolTable.Symbols {
		if len(symbol.References) > 0 {
			// Check if this symbol's position matches one of its references (indicating it's a reference symbol)
			for _, ref := range symbol.References {
				if ref.Line == symbol.Position.Line && ref.Column == symbol.Position.Column {
					referenceSymbols = append(referenceSymbols, symbol)
					break
				}
			}
		}
	}
	assert.Greater(t, len(referenceSymbols), 0, "Should have tracked identifier references")

	// Test specific function calls that we know should be present as references
	expectedReferences := []string{
		"SAYZ", // SAYZ WIT - should definitely be present
		"SAY",  // SAY WIT - should be present  
	}

	for _, expectedRef := range expectedReferences {
		ref := findSymbolReferenceByName(symbolTable.Symbols, expectedRef)
		assert.NotNil(t, ref, "Should find reference symbol for %s", expectedRef)
		if ref != nil {
			assert.Greater(t, len(ref.References), 0, "Symbol %s should have references", expectedRef)
		}
	}

	t.Logf("Found %d reference symbols", len(referenceSymbols))
}

func TestComprehensiveAnalysis_ScopeManagement(t *testing.T) {
	analyzer := NewSemanticAnalyzer("test://comprehensive.olol", comprehensiveExample)

	err := analyzer.AnalyzeDocument(context.Background())
	require.NoError(t, err)

	symbolTable := analyzer.GetSymbolTable()

	// Test scope hierarchy
	scopes := symbolTable.Scopes
	assert.Greater(t, len(scopes), 0, "Should have scopes")

	// Should have global scope
	globalScope := false
	classScopeCount := 0
	functionScopeCount := 0

	for _, scope := range scopes {
		switch scope.Type {
		case ScopeTypeGlobal:
			globalScope = true
		case ScopeTypeClass:
			classScopeCount++
		case ScopeTypeFunction:
			functionScopeCount++
		}
	}

	assert.True(t, globalScope, "Should have global scope")
	assert.Equal(t, 3, classScopeCount, "Should have 3 class scopes")
	assert.Greater(t, functionScopeCount, 0, "Should have function scopes")

	// Test symbol distribution across scopes
	globalSymbols := countSymbolsByScope(symbolTable.Symbols, ScopeTypeGlobal)
	classSymbols := countSymbolsByScope(symbolTable.Symbols, ScopeTypeClass)
	functionSymbols := countSymbolsByScope(symbolTable.Symbols, ScopeTypeFunction)

	assert.Greater(t, globalSymbols, 0, "Should have global symbols")
	assert.Greater(t, classSymbols, 0, "Should have class symbols")
	assert.Greater(t, functionSymbols, 0, "Should have function symbols")

	t.Logf("Scopes: Global=%d, Class=%d, Function=%d", globalSymbols, classSymbols, functionSymbols)
}

func TestComprehensiveAnalysis_HoverInfo(t *testing.T) {
	analyzer := NewSemanticAnalyzer("test://comprehensive.olol", comprehensiveExample)

	err := analyzer.AnalyzeDocument(context.Background())
	require.NoError(t, err)

	// Test hover on class declaration
	// Line numbers are 1-based in the source, 0-based in LSP
	classPosition := protocol.Position{Line: 7, Character: 16} // Around "CALCULATOR" class declaration
	hover := analyzer.GetHoverInfo(classPosition)

	require.NotNil(t, hover, "Expected hover info for class declaration")
	if markup, ok := hover.Contents.(protocol.MarkupContent); ok {
		assert.Contains(t, markup.Value, "CALCULATOR", "Hover should contain class name")
		assert.Contains(t, markup.Value, "CLAS", "Hover should indicate it's a class")
		t.Logf("Class hover content: %s", markup.Value)
	}

	// Test hover on function declaration
	funcPosition := protocol.Position{Line: 15, Character: 21} // Around "ADD" function
	hover = analyzer.GetHoverInfo(funcPosition)

	require.NotNil(t, hover, "Expected hover info for function declaration")
	if markup, ok := hover.Contents.(protocol.MarkupContent); ok {
		assert.Contains(t, markup.Value, "ADD", "Hover should contain function name")
		t.Logf("Function hover content: %s", markup.Value)
	}

	// Test hover on variable declaration
	varPosition := protocol.Position{Line: 70, Character: 21} // Around "CALC" variable
	hover = analyzer.GetHoverInfo(varPosition)

	require.NotNil(t, hover, "Expected hover info for variable declaration")
	if markup, ok := hover.Contents.(protocol.MarkupContent); ok {
		assert.Contains(t, markup.Value, "CALC", "Hover should contain variable name")
		t.Logf("Variable hover content: %s", markup.Value)
	}

	// Test hover on implicit member access (RESULT within ADD function)
	memberAccessPosition := protocol.Position{Line: 16, Character: 8} // Around "RESULT" in "RESULT ITZ X MOAR Y"
	hover = analyzer.GetHoverInfo(memberAccessPosition)

	require.NotNil(t, hover, "Expected hover info for implicit member access")
	if markup, ok := hover.Contents.(protocol.MarkupContent); ok {
		assert.Contains(t, markup.Value, "RESULT", "Hover should contain member variable name")
		assert.Contains(t, markup.Value, "DUBBLE", "Hover should contain member variable type")
		t.Logf("Member access hover content: %s", markup.Value)
	}
}

func TestComprehensiveAnalysis_Completions(t *testing.T) {
	analyzer := NewSemanticAnalyzer("test://comprehensive.olol", comprehensiveExample)

	err := analyzer.AnalyzeDocument(context.Background())
	require.NoError(t, err)

	// Test completions inside MAIN function
	position := protocol.Position{Line: 72, Character: 4} // Inside MAIN function
	completions := analyzer.GetCompletionItems(position)

	assert.Greater(t, len(completions), 0, "Should have completion items")

	// Should have keywords
	hasKeywords := false
	hasVariables := false
	hasClasses := false

	for _, completion := range completions {
		if completion.Label == "SAYZ" || completion.Label == "TEH" || completion.Label == "GIVEZ" {
			hasKeywords = true
		}
		if completion.Label == "CALC" || completion.Label == "SUM" {
			hasVariables = true
		}
		if completion.Label == "CALCULATOR" || completion.Label == "COUNTER" {
			hasClasses = true
		}
	}

	assert.True(t, hasKeywords, "Should have keyword completions")
	t.Logf("Found %d completion items", len(completions))
	t.Logf("Has keywords: %v, variables: %v, classes: %v", hasKeywords, hasVariables, hasClasses)
}

func TestComprehensiveAnalysis_SymbolResolution(t *testing.T) {
	analyzer := NewSemanticAnalyzer("test://comprehensive.olol", comprehensiveExample)

	err := analyzer.AnalyzeDocument(context.Background())
	require.NoError(t, err)

	// Test symbol resolution at various positions
	testCases := []struct {
		name     string
		line     uint32
		char     uint32
		expected string
	}{
		{"Class declaration", 7, 16, "CALCULATOR"},
		{"Function declaration", 15, 21, "ADD"},
		{"Variable declaration", 70, 21, "CALC"},
		{"Member variable", 9, 21, "RESULT"},
		{"Implicit member access", 16, 8, "RESULT"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			position := protocol.Position{Line: tc.line, Character: tc.char}
			symbol := analyzer.ResolveSymbolAtPosition(position)

			require.NotNil(t, symbol, "Expected symbol at position %d:%d for %s", tc.line, tc.char, tc.name)
			assert.Contains(t, strings.ToUpper(symbol.Name), strings.ToUpper(tc.expected),
				"Symbol at %d:%d should be %s, got %s", tc.line, tc.char, tc.expected, symbol.Name)
			t.Logf("Resolved %s: %s (kind: %d)", tc.name, symbol.Name, int(symbol.Kind))
		})
	}
}

func TestComprehensiveAnalysis_TypeInference(t *testing.T) {
	analyzer := NewSemanticAnalyzer("test://comprehensive.olol", comprehensiveExample)

	err := analyzer.AnalyzeDocument(context.Background())
	require.NoError(t, err)

	symbolTable := analyzer.GetSymbolTable()

	// Test type inference for variables with object instantiation
	calcVar := findSymbolByName(symbolTable.Symbols, "CALC")
	require.NotNil(t, calcVar, "Should find CALC variable")
	assert.Equal(t, "CALCULATOR", calcVar.Type, "CALC should be inferred as CALCULATOR type")

	counterVar := findSymbolByName(symbolTable.Symbols, "COUNTER")
	require.NotNil(t, counterVar, "Should find COUNTER variable")
	assert.Equal(t, "COUNTER", counterVar.Type, "COUNTER should be inferred as COUNTER type")

	// Test type inference for method call results
	sumVar := findSymbolByName(symbolTable.Symbols, "SUM")
	require.NotNil(t, sumVar, "Should find SUM variable")
	assert.Equal(t, "DUBBLE", sumVar.Type, "SUM should be DUBBLE (return type of ADD)")

	diffVar := findSymbolByName(symbolTable.Symbols, "DIFF")
	require.NotNil(t, diffVar, "Should find DIFF variable")
	assert.Equal(t, "DUBBLE", diffVar.Type, "DIFF should be DUBBLE (return type of SUBTRACT)")
}

func TestComprehensiveAnalysis_ErrorRecovery(t *testing.T) {
	// Test analyzer behavior with partially malformed code
	malformedCode := comprehensiveExample + "\n\nMALFORMED SYNTAX HERE"

	analyzer := NewSemanticAnalyzer("test://malformed.olol", malformedCode)

	err := analyzer.AnalyzeDocument(context.Background())
	// Analyzer should handle parse errors gracefully
	assert.NoError(t, err, "Analyzer should handle parse errors gracefully")

	symbolTable := analyzer.GetSymbolTable()

	// Should still have some symbols from the valid parts
	assert.Greater(t, len(symbolTable.Symbols), 0, "Should still have symbols despite parse errors")

	// Check diagnostics for parse errors
	diagnostics := analyzer.GetDiagnostics()
	hasParseError := false
	for _, diag := range diagnostics {
		if strings.Contains(diag.Message, "parse") || strings.Contains(diag.Message, "syntax") {
			hasParseError = true
			break
		}
	}

	t.Logf("Found %d diagnostics, has parse error: %v", len(diagnostics), hasParseError)
}

// Benchmark test for performance with complex code
func BenchmarkComprehensiveAnalysis(b *testing.B) {
	for i := 0; i < b.N; i++ {
		analyzer := NewSemanticAnalyzer("test://benchmark.olol", comprehensiveExample)
		err := analyzer.AnalyzeDocument(context.Background())
		if err != nil {
			b.Fatalf("AnalyzeDocument failed: %v", err)
		}

		// Test some analyzer operations
		position := protocol.Position{Line: 365, Character: 4}
		_ = analyzer.GetCompletionItems(position)
		_ = analyzer.GetHoverInfo(position)
		_ = analyzer.ResolveSymbolAtPosition(position)
	}
}

func TestIdentifierRangeHover(t *testing.T) {
	analyzer := NewSemanticAnalyzer("test://comprehensive.olol", comprehensiveExample)

	err := analyzer.AnalyzeDocument(context.Background())
	if err != nil {
		t.Fatalf("AnalyzeDocument failed: %v", err)
	}

	// Test hover at different positions within "RESULT" identifier
	// "RESULT ITZ X MOAR Y" is on line 16 (0-based)
	testCases := []struct {
		name       string
		line       uint32
		char       uint32
		shouldWork bool
	}{
		{"R in RESULT", 16, 8, true},  // First character
		{"E in RESULT", 16, 9, true},  // Second character - should now work
		{"S in RESULT", 16, 10, true}, // Third character - should now work
		{"U in RESULT", 16, 11, true}, // Fourth character - should now work
		{"L in RESULT", 16, 12, true}, // Fifth character - should now work
		{"T in RESULT", 16, 13, true}, // Sixth character - should now work
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			position := protocol.Position{Line: tc.line, Character: tc.char}
			hover := analyzer.GetHoverInfo(position)

			if tc.shouldWork {
				if hover == nil {
					t.Errorf("Expected hover info at position %d:%d but got nil", tc.line, tc.char)
				} else {
					t.Logf("✓ Hover works at position %d:%d", tc.line, tc.char)
				}
			} else {
				if hover == nil {
					t.Logf("✗ No hover info at position %d:%d (expected)", tc.line, tc.char)
				} else {
					t.Logf("✓ Unexpected hover info at position %d:%d", tc.line, tc.char)
				}
			}
		})
	}
}
