package parser

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLexerBasicTokens(t *testing.T) {
	input := `HAI 1.0
ME TEH VARIABLE ITZ "hello world"
KTHXBAI`

	lexer := NewLexer(input)

	expectedTokens := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{HAI, "HAI"},
		{DOUBLE, "1.0"},
		{NEWLINE, "\\n"},
		{ME, "ME"},
		{TEH, "TEH"},
		{VARIABLE, "VARIABLE"},
		{ITZ, "ITZ"},
		{STRING, "hello world"},
		{NEWLINE, "\\n"},
		{KTHXBAI, "KTHXBAI"},
		{EOF, ""},
	}

	for i, expectedToken := range expectedTokens {
		token, err := lexer.NextToken()
		require.NoError(t, err, "Error at token %d", i)

		assert.Equal(t, expectedToken.expectedType, token.Type,
			"Token %d: expected type %v, got %v", i, expectedToken.expectedType, token.Type)
		assert.Equal(t, expectedToken.expectedLiteral, token.Literal,
			"Token %d: expected literal %q, got %q", i, expectedToken.expectedLiteral, token.Literal)
	}
}

func TestLexerKeywords(t *testing.T) {
	keywords := map[string]TokenType{
		"HAI":      HAI,
		"KTHXBAI":  KTHXBAI,
		"ME":       ME,
		"TEH":      TEH,
		"VARIABLE": VARIABLE,
		"FUNCSHUN": FUNCSHUN,
		"CLAS":     CLAS,
		"KITTEH":   KITTEH,
		"OF":       OF,
		"ITZ":      ITZ,
		"WIT":      WIT,
		"AN":       AN,
		"DIS":      DIS,
		"IZ":       IZ,
		"NOPE":     NOPE,
		"WHILE":    WHILE,
		"GIVEZ":    GIVEZ,
		"UP":       UP,
		"NEW":      NEW,
		"AS":       AS,
		"SAEM":     SAEM,
		"MOAR":     MOAR,
		"LES":      LES,
		"TIEMZ":    TIEMZ,
		"DIVIDEZ":  DIVIDEZ,
		"BIGGR":    BIGGR,
		"SMALLR":   SMALLR,
		"THAN":     THAN,
		"OR":       OR,
		"YEZ":      YEZ,
		"NO":       NO,
		"NOTHIN":   NOTHIN,
		"INTEGR":   INTEGR_TYPE,
		"DUBBLE":   DUBBLE_TYPE,
		"STRIN":    STRIN_TYPE,
		"BOOL":     BOOL_TYPE,
	}

	for keyword, expectedType := range keywords {
		t.Run(keyword, func(t *testing.T) {
			lexer := NewLexer(keyword)
			token, err := lexer.NextToken()
			require.NoError(t, err)

			assert.Equal(t, expectedType, token.Type)
			assert.Equal(t, keyword, token.Literal)
		})
	}
}

func TestLexerCaseInsensitivity(t *testing.T) {
	testCases := []struct {
		input    string
		expected TokenType
	}{
		{"hai", HAI},
		{"HAI", HAI},
		{"Hai", HAI},
		{"HaI", HAI},
		{"variable", VARIABLE},
		{"VARIABLE", VARIABLE},
		{"Variable", VARIABLE},
		{"funcshun", FUNCSHUN},
		{"FUNCSHUN", FUNCSHUN},
		{"FuncShun", FUNCSHUN},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			lexer := NewLexer(tc.input)
			token, err := lexer.NextToken()
			require.NoError(t, err)
			assert.Equal(t, tc.expected, token.Type)
			// All keywords should be normalized to uppercase
			assert.Equal(t, strings.ToUpper(tc.input), token.Literal) // Keywords are uppercased
		})
	}
}

func TestLexerLiterals(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []struct {
			tokenType TokenType
			literal   string
		}
	}{
		{
			"Integer literals",
			"42 0 -123",
			[]struct {
				tokenType TokenType
				literal   string
			}{
				{INTEGER, "42"},
				{INTEGER, "0"},
				{INTEGER, "-123"}, // Negative numbers are parsed as single tokens
			},
		},
		{
			"Double literals",
			"3.14 0.0 123.456",
			[]struct {
				tokenType TokenType
				literal   string
			}{
				{DOUBLE, "3.14"},
				{DOUBLE, "0.0"},
				{DOUBLE, "123.456"},
			},
		},
		{
			"String literals",
			`"hello" "world with spaces" ""`,
			[]struct {
				tokenType TokenType
				literal   string
			}{
				{STRING, "hello"},
				{STRING, "world with spaces"},
				{STRING, ""},
			},
		},
		{
			"Boolean literals",
			"YEZ NO",
			[]struct {
				tokenType TokenType
				literal   string
			}{
				{YEZ, "YEZ"},
				{NO, "NO"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lexer := NewLexer(test.input)

			for i, expected := range test.expected {
				token, err := lexer.NextToken()
				require.NoError(t, err, "Error at token %d", i)

				assert.Equal(t, expected.tokenType, token.Type,
					"Token %d: expected type %v, got %v", i, expected.tokenType, token.Type)
				assert.Equal(t, expected.literal, token.Literal,
					"Token %d: expected literal %q, got %q", i, expected.literal, token.Literal)
			}
		})
	}
}

func TestLexerIdentifiers(t *testing.T) {
	input := "myVariable _private variable123 test_var"
	lexer := NewLexer(input)

	expectedIdentifiers := []string{"MYVARIABLE", "_PRIVATE", "VARIABLE123", "TEST_VAR"}

	for i, expected := range expectedIdentifiers {
		token, err := lexer.NextToken()
		require.NoError(t, err, "Error at token %d", i)

		assert.Equal(t, IDENTIFIER, token.Type)
		assert.Equal(t, expected, token.Literal)
	}
}

func TestLexerOperators(t *testing.T) {
	input := "( ) ?"
	lexer := NewLexer(input)

	expectedTokens := []struct {
		tokenType TokenType
		literal   string
	}{
		{LPAREN, "("},
		{RPAREN, ")"},
		{QUESTION, "?"},
	}

	for i, expected := range expectedTokens {
		token, err := lexer.NextToken()
		require.NoError(t, err, "Error at token %d", i)

		assert.Equal(t, expected.tokenType, token.Type)
		assert.Equal(t, expected.literal, token.Literal)
	}
}

func TestLexerComments(t *testing.T) {
	input := `HAI 1.0
BTW This is a comment
ME TEH VARIABLE ITZ 42`

	lexer := NewLexer(input)

	// Comments should be ignored, let's see what tokens we actually get
	var tokens []TokenType
	for {
		token, err := lexer.NextToken()
		require.NoError(t, err)
		if token.Type == EOF {
			break
		}
		tokens = append(tokens, token.Type)
	}

	// Check that BTW is not in the tokens (comment should be skipped)
	for _, tokenType := range tokens {
		assert.NotEqual(t, BTW, tokenType, "BTW comment should be filtered out")
	}

	// Check that we have the main tokens (ignoring the exact sequence since comments are filtered)
	expectedTokens := []TokenType{HAI, DOUBLE, NEWLINE, ME, TEH, VARIABLE, ITZ, INTEGER}
	assert.Contains(t, tokens, HAI)
	assert.Contains(t, tokens, DOUBLE)
	assert.Contains(t, tokens, ME)
	assert.Contains(t, tokens, TEH)
	assert.Contains(t, tokens, VARIABLE)
	assert.Contains(t, tokens, ITZ)
	assert.Contains(t, tokens, INTEGER)

	// Should have at least as many tokens as expected (might have extra newlines)
	assert.GreaterOrEqual(t, len(tokens), len(expectedTokens)-1)
}

func TestLexerWhitespace(t *testing.T) {
	input := "HAI   1.0\n\nME\t\tTEH"
	lexer := NewLexer(input)

	expectedTokens := []TokenType{HAI, DOUBLE, NEWLINE, NEWLINE, ME, TEH}

	for i, expectedType := range expectedTokens {
		token, err := lexer.NextToken()
		require.NoError(t, err, "Error at token %d", i)
		assert.Equal(t, expectedType, token.Type)
	}
}

func TestLexerMultipleStatements(t *testing.T) {
	input := `HAI 1.0
ME TEH x ITZ 42
ME TEH y ITZ "hello"
x MOAR y
KTHXBAI`

	lexer := NewLexer(input)
	tokens := []Token{}

	for {
		token, err := lexer.NextToken()
		require.NoError(t, err)

		tokens = append(tokens, token)

		if token.Type == EOF {
			break
		}
	}

	// Verify we got all expected tokens
	assert.Greater(t, len(tokens), 10, "Should have more than 10 tokens")

	// Check first and last tokens
	assert.Equal(t, HAI, tokens[0].Type)
	assert.Equal(t, EOF, tokens[len(tokens)-1].Type)
}

func TestLexerErrorHandling(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"Unterminated string", `"hello world`},
		{"Invalid character", "HAI 1.0\n@#$%"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lexer := NewLexer(test.input)

			// Should eventually hit an error or illegal token
			foundError := false
			for i := 0; i < 20; i++ { // Prevent infinite loop
				token, err := lexer.NextToken()
				if err != nil || token.Type == ILLEGAL {
					foundError = true
					break
				}
				if token.Type == EOF {
					break
				}
			}

			assert.True(t, foundError, "Expected to encounter an error or illegal token")
		})
	}
}
