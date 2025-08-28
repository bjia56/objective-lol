package parser

import (
	"fmt"
	"strconv"
	"strings"
)

// TokenType represents the type of a lexical token
type TokenType int

const (
	// Special tokens
	ILLEGAL TokenType = iota
	EOF
	NEWLINE

	// Literals
	IDENTIFIER
	STRING
	INTEGER
	DOUBLE
	BOOLEAN

	// Keywords
	HAI
	ME
	TEH
	VARIABLE
	FUNCSHUN
	NATIV
	CLAS
	KITTEH
	OF
	LOCKD
	SHARD
	KTHXBAI
	KTHX
	ITZ
	WIT
	AN
	DIS
	EVRYONE
	MAHSELF
	IZ
	NOPE
	WHILE
	GIVEZ
	UP
	NEW
	IN
	AS
	BTW
	I
	HAS
	A
	DO

	// Operators
	MOAR    // +
	LES     // -
	TIEMZ   // *
	DIVIDEZ // /
	BIGGR   // >
	THAN    // than (used with BIGGR/SMALLR)
	SMALLR  // <
	SAEM    // ==
	OR      // ||

	// Types
	BOOL_TYPE
	INTEGR_TYPE
	DUBBLE_TYPE
	STRIN_TYPE

	// Punctuation
	QUESTION
	LPAREN  // (
	RPAREN  // )

	// Special values
	YEZ
	NO
	NOTHIN
)

// String returns the string representation of a TokenType
func (t TokenType) String() string {
	switch t {
	case ILLEGAL:
		return "ILLEGAL"
	case EOF:
		return "EOF"
	case NEWLINE:
		return "NEWLINE"
	case IDENTIFIER:
		return "IDENTIFIER"
	case STRING:
		return "STRING"
	case INTEGER:
		return "INTEGER"
	case DOUBLE:
		return "DOUBLE"
	case BOOLEAN:
		return "BOOLEAN"
	case HAI:
		return "HAI"
	case ME:
		return "ME"
	case TEH:
		return "TEH"
	case VARIABLE:
		return "VARIABLE"
	case FUNCSHUN:
		return "FUNCSHUN"
	case NATIV:
		return "NATIV"
	case CLAS:
		return "CLAS"
	case KITTEH:
		return "KITTEH"
	case OF:
		return "OF"
	case LOCKD:
		return "LOCKD"
	case SHARD:
		return "SHARD"
	case KTHXBAI:
		return "KTHXBAI"
	case KTHX:
		return "KTHX"
	case ITZ:
		return "ITZ"
	case WIT:
		return "WIT"
	case AN:
		return "AN"
	case DIS:
		return "DIS"
	case EVRYONE:
		return "EVRYONE"
	case MAHSELF:
		return "MAHSELF"
	case IZ:
		return "IZ"
	case NOPE:
		return "NOPE"
	case WHILE:
		return "WHILE"
	case GIVEZ:
		return "GIVEZ"
	case UP:
		return "UP"
	case NEW:
		return "NEW"
	case IN:
		return "IN"
	case AS:
		return "AS"
	case BTW:
		return "BTW"
	case I:
		return "I"
	case HAS:
		return "HAS"
	case A:
		return "A"
	case DO:
		return "DO"
	case MOAR:
		return "MOAR"
	case LES:
		return "LES"
	case TIEMZ:
		return "TIEMZ"
	case DIVIDEZ:
		return "DIVIDEZ"
	case BIGGR:
		return "BIGGR"
	case THAN:
		return "THAN"
	case SMALLR:
		return "SMALLR"
	case SAEM:
		return "SAEM"
	case OR:
		return "OR"
	case BOOL_TYPE:
		return "BOOL"
	case INTEGR_TYPE:
		return "INTEGR"
	case DUBBLE_TYPE:
		return "DUBBLE"
	case STRIN_TYPE:
		return "STRIN"
	case QUESTION:
		return "?"
	case LPAREN:
		return "("
	case RPAREN:
		return ")"
	case YEZ:
		return "YEZ"
	case NO:
		return "NO"
	case NOTHIN:
		return "NOTHIN"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", int(t))
	}
}

// Token represents a lexical token
type Token struct {
	Type     TokenType
	Literal  string
	Position int
	Line     int
	Column   int
}

// Lexer tokenizes Objective-LOL source code
type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
	line         int  // current line number
	column       int  // current column number
}

// NewLexer creates a new lexer instance
func NewLexer(input string) *Lexer {
	l := &Lexer{
		input:  input,
		line:   1,
		column: 0,
	}
	l.readChar()
	return l
}

// Keywords maps keyword strings to their token types
var keywords = map[string]TokenType{
	"HAI":      HAI,
	"ME":       ME,
	"TEH":      TEH,
	"VARIABLE": VARIABLE,
	"FUNCSHUN": FUNCSHUN,
	"NATIV":    NATIV,
	"CLAS":     CLAS,
	"KITTEH":   KITTEH,
	"OF":       OF,
	"LOCKD":    LOCKD,
	"SHARD":    SHARD,
	"KTHXBAI":  KTHXBAI,
	"KTHX":     KTHX,
	"ITZ":      ITZ,
	"WIT":      WIT,
	"AN":       AN,
	"DIS":      DIS,
	"EVRYONE":  EVRYONE,
	"MAHSELF":  MAHSELF,
	"IZ":       IZ,
	"NOPE":     NOPE,
	"WHILE":    WHILE,
	"GIVEZ":    GIVEZ,
	"UP":       UP,
	"NEW":      NEW,
	"IN":       IN,
	"AS":       AS,
	"BTW":      BTW,
	"MOAR":     MOAR,
	"LES":      LES,
	"TIEMZ":    TIEMZ,
	"DIVIDEZ":  DIVIDEZ,
	"BIGGR":    BIGGR,
	"THAN":     THAN,
	"SMALLR":   SMALLR,
	"SAEM":     SAEM,
	"OR":       OR,
	"BOOL":     BOOL_TYPE,
	"INTEGR":   INTEGR_TYPE,
	"DUBBLE":   DUBBLE_TYPE,
	"STRIN":    STRIN_TYPE,
	"YEZ":      YEZ,
	"NO":       NO,
	"I":        I,
	"HAS":      HAS,
	"A":        A,
	"DO":       DO,
	"NOTHIN":   NOTHIN,
}

// readChar reads the next character and advances position
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // ASCII NUL character represents EOF
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++

	if l.ch == '\n' {
		l.line++
		l.column = 0
	} else {
		l.column++
	}
}

// peekChar returns the next character without advancing position
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// skipWhitespace skips whitespace characters
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' {
		l.readChar()
	}
}

// readIdentifier reads an identifier or keyword
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' {
		l.readChar()
	}
	return l.input[position:l.position]
}

// readString reads a string literal
func (l *Lexer) readString() (string, error) {
	position := l.position + 1 // skip opening quote
	for {
		l.readChar()
		if l.ch == '"' {
			break
		}
		if l.ch == 0 {
			return "", fmt.Errorf("unterminated string literal at line %d", l.line)
		}
		// Handle escape sequences
		if l.ch == '\\' {
			l.readChar() // skip the escaped character
		}
	}
	return l.input[position:l.position], nil
}

// readNumber reads a numeric literal (integer or double)
func (l *Lexer) readNumber() (string, TokenType) {
	position := l.position
	tokenType := INTEGER

	// Handle negative sign
	if l.ch == '-' {
		l.readChar() // consume '-'
	}

	// Handle hex numbers (0X prefix)
	if l.ch == '0' && (l.peekChar() == 'X' || l.peekChar() == 'x') {
		l.readChar() // consume '0'
		l.readChar() // consume 'X'
		for isHexDigit(l.ch) {
			l.readChar()
		}
		return l.input[position:l.position], INTEGER
	}

	// Read digits
	for isDigit(l.ch) {
		l.readChar()
	}

	// Check for decimal point
	if l.ch == '.' && isDigit(l.peekChar()) {
		tokenType = DOUBLE
		l.readChar() // consume '.'
		for isDigit(l.ch) {
			l.readChar()
		}
	}

	return l.input[position:l.position], tokenType
}

// readComment reads a single-line comment (BTW until end of line)
func (l *Lexer) readComment() error {
	// Skip "BTW"
	for i := 0; i < 3; i++ {
		l.readChar()
	}

	// Read until end of line or EOF
	for l.ch != '\n' && l.ch != '\r' && l.ch != 0 {
		l.readChar()
	}

	return nil
}

// NextToken returns the next token from the input
func (l *Lexer) NextToken() (Token, error) {
	var tok Token

	l.skipWhitespace()

	tok.Line = l.line
	tok.Column = l.column
	tok.Position = l.position

	switch l.ch {
	case '\n', '\r':
		tok = Token{Type: NEWLINE, Literal: "\\n", Line: l.line, Column: l.column, Position: l.position}
		l.readChar()
		if l.ch == '\n' && tok.Literal == "\\r" { // Handle CRLF
			l.readChar()
		}
	case '?':
		tok = Token{Type: QUESTION, Literal: string(l.ch), Line: l.line, Column: l.column, Position: l.position}
		l.readChar()
	case '(':
		tok = Token{Type: LPAREN, Literal: string(l.ch), Line: l.line, Column: l.column, Position: l.position}
		l.readChar()
	case ')':
		tok = Token{Type: RPAREN, Literal: string(l.ch), Line: l.line, Column: l.column, Position: l.position}
		l.readChar()
	case '"':
		str, err := l.readString()
		if err != nil {
			return tok, err
		}
		tok = Token{Type: STRING, Literal: str, Line: l.line, Column: l.column, Position: l.position}
		l.readChar() // consume closing quote
	case 0:
		tok = Token{Type: EOF, Literal: "", Line: l.line, Column: l.column, Position: l.position}
	default:
		if isLetter(l.ch) {
			// Check for BTW comment
			if l.ch == 'B' && l.position+2 < len(l.input) &&
				l.input[l.position:l.position+3] == "BTW" {
				err := l.readComment()
				if err != nil {
					return tok, err
				}
				return l.NextToken() // Return next token after comment
			}

			literal := l.readIdentifier()
			literal = strings.ToUpper(literal) // Case insensitive

			if tokenType, exists := keywords[literal]; exists {
				tok.Type = tokenType
			} else {
				tok.Type = IDENTIFIER
			}
			tok.Literal = literal
			return tok, nil
		} else if isDigit(l.ch) {
			literal, tokenType := l.readNumber()
			tok.Type = tokenType
			tok.Literal = literal
			return tok, nil
		} else if l.ch == '-' && isDigit(l.peekChar()) {
			// Negative number: '-' followed by digit
			literal, tokenType := l.readNumber()
			tok.Type = tokenType
			tok.Literal = literal
			return tok, nil
		} else {
			tok = Token{Type: ILLEGAL, Literal: string(l.ch), Line: l.line, Column: l.column, Position: l.position}
			l.readChar()
		}
	}

	return tok, nil
}

// GetAllTokens returns all tokens from the input
func (l *Lexer) GetAllTokens() ([]Token, error) {
	var tokens []Token

	for {
		tok, err := l.NextToken()
		if err != nil {
			return nil, err
		}

		tokens = append(tokens, tok)

		if tok.Type == EOF {
			break
		}
	}

	return tokens, nil
}

// Helper functions

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isHexDigit(ch byte) bool {
	return isDigit(ch) || ('A' <= ch && ch <= 'F') || ('a' <= ch && ch <= 'f')
}

// ConvertValue converts a token literal to a Go value
func ConvertValue(tok Token) (interface{}, error) {
	switch tok.Type {
	case STRING:
		// Handle escape sequences
		str := tok.Literal
		str = strings.ReplaceAll(str, "\\\"", "\"")
		str = strings.ReplaceAll(str, "\\n", "\n")
		str = strings.ReplaceAll(str, "\\t", "\t")
		str = strings.ReplaceAll(str, "\\r", "\r")
		str = strings.ReplaceAll(str, "\\\\", "\\")
		return str, nil

	case INTEGER:
		// Handle hex numbers
		literal := tok.Literal
		if strings.HasPrefix(strings.ToUpper(literal), "0X") {
			val, err := strconv.ParseInt(literal[2:], 16, 64)
			return val, err
		}
		// Handle negative hex numbers
		if strings.HasPrefix(strings.ToUpper(literal), "-0X") {
			val, err := strconv.ParseInt(literal[3:], 16, 64)
			if err != nil {
				return nil, err
			}
			return -val, nil
		}
		val, err := strconv.ParseInt(literal, 10, 64)
		return val, err

	case DOUBLE:
		val, err := strconv.ParseFloat(tok.Literal, 64)
		return val, err

	case YEZ:
		return true, nil

	case NO:
		return false, nil

	case NOTHIN:
		return nil, nil

	default:
		return tok.Literal, nil
	}
}
