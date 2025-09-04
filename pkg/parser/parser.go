package parser

import (
	"fmt"
	"slices"

	"github.com/bjia56/objective-lol/pkg/ast"
	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/types"
)

// Parser implements a recursive descent parser for Objective-LOL
type Parser struct {
	lexer        *Lexer
	currentToken Token
	peekToken    Token
	errors       []string
}

// NewParser creates a new parser instance
func NewParser(lexer *Lexer) *Parser {
	p := &Parser{
		lexer:  lexer,
		errors: []string{},
	}

	// Read two tokens, so currentToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

// nextToken advances both currentToken and peekToken
func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	var err error
	p.peekToken, err = p.lexer.NextToken()
	if err != nil {
		p.addError(err.Error())
	}
}

// addError adds an error message to the parser's error list
func (p *Parser) addError(msg string) {
	p.errors = append(p.errors, msg)
}

// Errors returns all parsing errors
func (p *Parser) Errors() []string {
	return p.errors
}

// convertPosition converts parser.PositionInfo to ast.PositionInfo
func (p *Parser) convertPosition(pos PositionInfo) ast.PositionInfo {
	return ast.PositionInfo{
		Line:   pos.Line,
		Column: pos.Column,
	}
}

// collectPrecedingComments collects documentation comments immediately preceding a declaration
func (p *Parser) collectPrecedingComments() []string {
	comments := p.lexer.GetRecentComments()
	if len(comments) == 0 {
		return nil
	}

	var documentation []string
	
	// Find the most recent contiguous block of comments
	// Comments must be immediately before the current token (allowing only whitespace/newlines between)
	currentLine := p.currentToken.Position.Line
	
	// Work backwards through comments to find the contiguous block
	var relevantComments []Token
	for i := len(comments) - 1; i >= 0; i-- {
		comment := comments[i]
		
		// Check if this comment is part of a contiguous block preceding the declaration
		expectedLine := currentLine - (len(relevantComments) + 1)
		if comment.Position.Line == expectedLine || 
		   (len(relevantComments) == 0 && comment.Position.Line < currentLine && comment.Position.Line >= currentLine-5) {
			relevantComments = append([]Token{comment}, relevantComments...)
		} else {
			break
		}
	}

	// Convert comment tokens to documentation strings
	for _, comment := range relevantComments {
		// Include empty comments as empty strings for proper formatting
		documentation = append(documentation, comment.Literal)
	}

	return documentation
}

// currentTokenIs checks if the current token matches the expected type
func (p *Parser) currentTokenIs(t ...TokenType) bool {
	return slices.Contains(t, p.currentToken.Type)
}

// peekTokenIs checks if the peek token matches the expected type
func (p *Parser) peekTokenIs(t TokenType) bool {
	return p.peekToken.Type == t
}

// expectPeek checks the peek token type and advances if it matches
func (p *Parser) expectPeek(t TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.addError(fmt.Sprintf("expected next token to be %v, got %v instead at line %d", t, p.peekToken.Type, p.peekToken.Position.Line))
	return false
}

// ParseProgram parses the entire program and returns the AST
func (p *Parser) ParseProgram() *ast.ProgramNode {
	program := &ast.ProgramNode{}
	program.Declarations = []ast.Node{}

	p.skipNewlines() // skip any leading newlines

	for !p.currentTokenIs(EOF) {
		var node ast.Node

		// Check if this is an import statement or a regular declaration
		if p.currentTokenIs(I) && p.peekTokenIs(CAN) {
			// This is an import statement
			node = p.parseImportStatement()
		} else if p.currentTokenIs(HAI) {
			// This is a regular declaration
			node = p.parseDeclaration()
		} else {
			p.addError(fmt.Sprintf("expected 'HAI' for declaration or 'I CAN HAS' for import, got %v at line %d", p.currentToken.Type, p.currentToken.Position.Line))
			p.nextToken()
			continue
		}

		if node != nil {
			program.Declarations = append(program.Declarations, node)
		}
		p.nextToken()
		p.skipNewlines() // skip newlines between declarations
	}

	return program
}

// parseDeclaration parses a top-level declaration
func (p *Parser) parseDeclaration() ast.Node {
	if !p.currentTokenIs(HAI) {
		p.addError(fmt.Sprintf("expected 'HAI' at start of declaration, got %v at line %d", p.currentToken.Type, p.currentToken.Position.Line))
		return nil
	}

	if !p.expectPeek(ME) {
		p.addError(fmt.Sprintf("expected 'ME' after 'HAI', got %v at line %d", p.peekToken.Type, p.peekToken.Position.Line))
		return nil
	}

	if !p.expectPeek(TEH) {
		p.addError(fmt.Sprintf("expected 'TEH' after 'HAI ME', got %v at line %d", p.peekToken.Type, p.peekToken.Position.Line))
		return nil
	}

	switch p.peekToken.Type {
	case VARIABLE, LOCKD:
		return p.parseVariableDeclaration()
	case FUNCSHUN:
		return p.parseFunctionDeclaration()
	case CLAS:
		return p.parseClassDeclaration()
	default:
		p.addError(fmt.Sprintf("unexpected token after 'HAI ME TEH': %v at line %d", p.peekToken.Type, p.peekToken.Position.Line))
		return nil
	}
}

// parseVariableDeclaration parses variable declarations
func (p *Parser) parseVariableDeclaration() *ast.VariableDeclarationNode {
	node := &ast.VariableDeclarationNode{}

	// Check for LOCKD
	if p.peekTokenIs(LOCKD) {
		node.IsLocked = true
		p.nextToken()
	}

	if !p.expectPeek(VARIABLE) {
		p.addError(fmt.Sprintf("expected 'VARIABLE', got %v at line %d", p.peekToken.Type, p.peekToken.Position.Line))
		return nil
	}

	if !p.expectPeek(IDENTIFIER) {
		p.addError(fmt.Sprintf("expected identifier, got %v at line %d", p.peekToken.Type, p.peekToken.Position.Line))
		return nil
	}
	node.Name = p.currentToken.Literal

	// Set position from the identifier token
	node.Position = p.convertPosition(p.currentToken.Position)

	if !p.expectPeek(TEH) {
		p.addError(fmt.Sprintf("expected 'TEH' after '%s', got %v at line %d", p.currentToken.Literal, p.peekToken.Type, p.peekToken.Position.Line))
		return nil
	}

	p.nextToken()
	if !p.currentTokenIs(IDENTIFIER) && !p.isTypeToken(p.currentToken.Type) {
		p.addError(fmt.Sprintf("expected type after 'TEH', got %v at line %d", p.currentToken.Type, p.currentToken.Position.Line))
		return nil
	}
	node.Type = p.currentToken.Literal

	// Check for initialization
	if p.peekTokenIs(ITZ) {
		p.nextToken() // consume ITZ
		p.nextToken() // move to expression
		node.Value = p.parseExpression()
	}

	return node
}

// parseIHasAVariableDeclaration parses "I HAS A" variable declarations
func (p *Parser) parseIHasAVariableDeclaration() *ast.VariableDeclarationNode {
	node := &ast.VariableDeclarationNode{}

	// Expect "I HAS A"
	if !p.expectPeek(HAS) {
		p.addError(fmt.Sprintf("expected 'HAS', got %v at line %d", p.peekToken.Type, p.peekToken.Position.Line))
		return nil
	}
	if !p.expectPeek(A) {
		p.addError(fmt.Sprintf("expected 'A', got %v at line %d", p.peekToken.Type, p.peekToken.Position.Line))
		return nil
	}

	// Check for LOCKD
	if p.peekTokenIs(LOCKD) {
		node.IsLocked = true
		p.nextToken()
	}

	if !p.expectPeek(VARIABLE) {
		p.addError(fmt.Sprintf("expected 'VARIABLE', got %v at line %d", p.peekToken.Type, p.peekToken.Position.Line))
		return nil
	}

	if !p.expectPeek(IDENTIFIER) {
		p.addError(fmt.Sprintf("expected identifier, got %v at line %d", p.peekToken.Type, p.peekToken.Position.Line))
		return nil
	}
	node.Name = p.currentToken.Literal

	// Set position from the identifier token
	node.Position = p.convertPosition(p.currentToken.Position)

	if !p.expectPeek(TEH) {
		p.addError(fmt.Sprintf("expected 'TEH' after '%s', got %v at line %d", p.currentToken.Literal, p.peekToken.Type, p.peekToken.Position.Line))
		return nil
	}

	p.nextToken()
	if !p.currentTokenIs(IDENTIFIER) && !p.isTypeToken(p.currentToken.Type) {
		p.addError(fmt.Sprintf("expected type after 'TEH', got %v at line %d", p.currentToken.Type, p.currentToken.Position.Line))
		return nil
	}
	node.Type = p.currentToken.Literal

	// Check for initialization
	if p.peekTokenIs(ITZ) {
		p.nextToken() // consume ITZ
		p.nextToken() // move to expression
		node.Value = p.parseExpression()
	}

	return node
}

// parseImportStatement parses both "I CAN HAS module?" and "I CAN HAS decl1 AN decl2 FROM module?" import statements
func (p *Parser) parseImportStatement() *ast.ImportStatementNode {
	node := &ast.ImportStatementNode{}

	// Expect "I CAN HAS"
	if !p.expectPeek(CAN) {
		p.addError(fmt.Sprintf("expected 'CAN', got %v at line %d", p.peekToken.Type, p.peekToken.Position.Line))
		return nil
	}
	if !p.expectPeek(HAS) {
		p.addError(fmt.Sprintf("expected 'HAS', got %v at line %d", p.peekToken.Type, p.peekToken.Position.Line))
		return nil
	}

	// Check if first token is STRING (file import) or IDENTIFIER (built-in module or selective import)
	if p.peekTokenIs(STRING) {
		// File import: I CAN HAS "filepath"?
		p.nextToken() // consume STRING
		node.IsFileImport = true
		node.ModuleName = p.currentToken.Literal

		// Set position from the string token
		node.Position = p.convertPosition(p.currentToken.Position)
	} else if p.peekTokenIs(IDENTIFIER) {
		// Either built-in module import or selective import
		p.nextToken() // consume IDENTIFIER
		firstIdentifier := p.currentToken.Literal

		// Check if this is selective import (next token is AN or FROM)
		if p.peekToken.Type == AN || p.peekToken.Type == FROM {
			// Selective import: I CAN HAS decl1 [AN decl2 ...] FROM module?
			node.Declarations = []string{firstIdentifier}

			// Parse additional declarations separated by AN
			for p.peekToken.Type == AN {
				p.nextToken() // consume AN
				if !p.expectPeek(IDENTIFIER) {
					p.addError(fmt.Sprintf("expected declaration name after 'AN', got %v at line %d", p.peekToken.Type, p.peekToken.Position.Line))
					return nil
				}
				node.Declarations = append(node.Declarations, p.currentToken.Literal)
			}

			// Expect FROM
			if !p.expectPeek(FROM) {
				p.addError(fmt.Sprintf("expected 'FROM', got %v at line %d", p.peekToken.Type, p.peekToken.Position.Line))
				return nil
			}

			// Module name can be either IDENTIFIER (built-in) or STRING (file)
			if p.peekTokenIs(STRING) {
				// Selective import from file: I CAN HAS FUNC FROM "file"?
				p.nextToken() // consume STRING
				node.IsFileImport = true
				node.ModuleName = p.currentToken.Literal
			} else if p.peekTokenIs(IDENTIFIER) {
				// Selective import from built-in: I CAN HAS FUNC FROM STDIO?
				p.nextToken() // consume IDENTIFIER
				node.IsFileImport = false
				node.ModuleName = p.currentToken.Literal
			} else {
				p.addError(fmt.Sprintf("expected module name (identifier or string) after 'FROM', got %v at line %d", p.peekToken.Type, p.peekToken.Position.Line))
				return nil
			}

			// Set position from the string or identifier token
			node.Position = p.convertPosition(p.currentToken.Position)
		} else {
			// Traditional built-in import: I CAN HAS STDIO?
			node.IsFileImport = false
			node.ModuleName = firstIdentifier

			// Set position from the identifier token
			node.Position = p.convertPosition(p.currentToken.Position)
		}
	} else {
		p.addError(fmt.Sprintf("expected module name (identifier) or file path (string) after 'HAS', got %v at line %d", p.peekToken.Type, p.peekToken.Position.Line))
		return nil
	}

	// Expect '?' at the end
	if !p.expectPeek(QUESTION) {
		p.addError(fmt.Sprintf("expected '?' at end of import statement, got %v at line %d", p.peekToken.Type, p.peekToken.Position.Line))
		return nil
	}

	return node
}

// parseFunctionDeclaration parses function declarations
func (p *Parser) parseFunctionDeclaration() *ast.FunctionDeclarationNode {
	node := &ast.FunctionDeclarationNode{}

	// Collect documentation comments before parsing the function declaration
	node.Documentation = p.collectPrecedingComments()
	
	// Clear the comments buffer to avoid reusing them for subsequent declarations
	p.lexer.ClearRecentComments()

	if !p.expectPeek(FUNCSHUN) {
		p.addError(fmt.Sprintf("expected 'FUNCSHUN', got %v at line %d", p.peekToken.Type, p.peekToken.Position.Line))
		return nil
	}

	if !p.expectPeek(IDENTIFIER) {
		p.addError(fmt.Sprintf("expected identifier, got %v at line %d", p.peekToken.Type, p.peekToken.Position.Line))
		return nil
	}
	node.Name = p.currentToken.Literal

	// Set position from the identifier token
	node.Position = p.convertPosition(p.currentToken.Position)

	// Check for return type (TEH type)
	if p.peekTokenIs(TEH) {
		p.nextToken() // consume TEH
		p.nextToken() // move to the type token
		if !p.currentTokenIs(IDENTIFIER) && !p.isTypeToken(p.currentToken.Type) {
			p.addError(fmt.Sprintf("expected type after TEH, got %v at line %d", p.currentToken.Type, p.currentToken.Position.Line))
			return nil
		}
		node.ReturnType = p.currentToken.Literal
	}

	// Check for parameters (WIT ...)
	if p.peekTokenIs(WIT) {
		p.nextToken() // consume WIT
		node.Parameters = p.parseParameterList()
	}

	// Parse function body
	node.Body = p.parseStatementBlock(KTHXBAI)

	if !p.currentTokenIs(KTHXBAI) {
		p.addError(fmt.Sprintf("expected 'KTHXBAI' to close function, got %v at line %d", p.currentToken.Type, p.currentToken.Position.Line))
		return nil
	}

	return node
}

// parseClassDeclaration parses class declarations
func (p *Parser) parseClassDeclaration() *ast.ClassDeclarationNode {
	node := &ast.ClassDeclarationNode{}

	if !p.expectPeek(CLAS) {
		p.addError(fmt.Sprintf("expected 'CLAS', got %v at line %d", p.peekToken.Type, p.peekToken.Position.Line))
		return nil
	}

	if !p.expectPeek(IDENTIFIER) {
		p.addError(fmt.Sprintf("expected identifier after 'CLAS', got %v at line %d", p.peekToken.Type, p.peekToken.Position.Line))
		return nil
	}
	node.Name = p.currentToken.Literal

	// Set position from the identifier token
	node.Position = p.convertPosition(p.currentToken.Position)

	// Check for inheritance (KITTEH OF parent)
	if p.peekTokenIs(KITTEH) {
		p.nextToken() // consume KITTEH
		if !p.expectPeek(OF) {
			p.addError(fmt.Sprintf("expected 'OF' after 'KITTEH', got %v at line %d", p.peekToken.Type, p.peekToken.Position.Line))
			return nil
		}
		if !p.expectPeek(IDENTIFIER) {
			p.addError(fmt.Sprintf("expected identifier after 'OF', got %v at line %d", p.peekToken.Type, p.peekToken.Position.Line))
			return nil
		}
		node.ParentClass = p.currentToken.Literal
	}

	// Parse class body
	p.nextToken()
	node.Members = p.parseClassMembers()

	if !p.currentTokenIs(KTHXBAI) {
		p.addError(fmt.Sprintf("expected 'KTHXBAI' to close class, got %v at line %d", p.currentToken.Type, p.currentToken.Position.Line))
		return nil
	}

	return node
}

// parseClassMembers parses class member declarations
func (p *Parser) parseClassMembers() []*ast.ClassMemberNode {
	var members []*ast.ClassMemberNode
	isPublic := true // default visibility

	for !p.currentTokenIs(KTHXBAI) && !p.currentTokenIs(EOF) {
		p.skipNewlines() // skip newlines in class body

		// Skip empty lines
		if p.currentTokenIs(NEWLINE) || p.currentTokenIs(KTHXBAI) || p.currentTokenIs(EOF) {
			break
		}

		// Check for visibility modifiers
		if p.currentTokenIs(EVRYONE) {
			isPublic = true
			p.nextToken()
			continue
		} else if p.currentTokenIs(MAHSELF) {
			isPublic = false
			p.nextToken()
			continue
		}

		if p.currentTokenIs(DIS) {
			member := p.parseClassMember()
			if member != nil {
				member.IsPublic = isPublic
				members = append(members, member)
			}
		}

		p.nextToken()
	}

	return members
}

// parseClassMember parses a single class member
func (p *Parser) parseClassMember() *ast.ClassMemberNode {
	if !p.currentTokenIs(DIS) {
		p.addError(fmt.Sprintf("expected 'DIS', got %v at line %d", p.currentToken.Type, p.currentToken.Position.Line))
		return nil
	}

	if !p.expectPeek(TEH) {
		p.addError(fmt.Sprintf("expected 'TEH', got %v at line %d", p.peekToken.Type, p.peekToken.Position.Line))
		return nil
	}

	member := &ast.ClassMemberNode{}

	// Check for SHARD
	if p.peekTokenIs(SHARD) {
		member.IsShared = true
		p.nextToken()
	}

	// Check for LOCKD (only valid for variables)
	isLocked := false
	if p.peekTokenIs(LOCKD) {
		isLocked = true
		p.nextToken()
	}

	if p.peekTokenIs(VARIABLE) {
		// Member variable
		p.nextToken() // consume VARIABLE
		member.IsVariable = true

		if !p.expectPeek(IDENTIFIER) {
			p.addError(fmt.Sprintf("expected identifier after 'VARIABLE', got %v at line %d", p.peekToken.Type, p.peekToken.Position.Line))
			return nil
		}
		name := p.currentToken.Literal

		if !p.expectPeek(TEH) {
			p.addError(fmt.Sprintf("expected 'TEH', got %v at line %d", p.peekToken.Type, p.peekToken.Position.Line))
			return nil
		}

		p.nextToken()
		if !p.currentTokenIs(IDENTIFIER) && !p.isTypeToken(p.currentToken.Type) {
			p.addError(fmt.Sprintf("expected type after 'TEH', got %v at line %d", p.currentToken.Type, p.currentToken.Position.Line))
			return nil
		}
		varType := p.currentToken.Literal

		varDecl := &ast.VariableDeclarationNode{
			Name:     name,
			Type:     varType,
			IsLocked: isLocked,
		}

		// Check for initialization
		if p.peekTokenIs(ITZ) {
			p.nextToken() // consume ITZ
			p.nextToken() // move to expression
			varDecl.Value = p.parseExpression()
		}

		member.Variable = varDecl

	} else if p.peekTokenIs(FUNCSHUN) {
		// Member function
		p.nextToken() // consume FUNCSHUN
		member.IsVariable = false

		if !p.expectPeek(IDENTIFIER) {
			p.addError(fmt.Sprintf("expected identifier after 'FUNCSHUN', got %v at line %d", p.peekToken.Type, p.peekToken.Position.Line))
			return nil
		}
		name := p.currentToken.Literal

		funcDecl := &ast.FunctionDeclarationNode{
			Name:     name,
			IsShared: &member.IsShared,
		}

		// Check for return type
		if p.peekTokenIs(TEH) {
			p.nextToken() // consume TEH
			p.nextToken()
			if !p.currentTokenIs(IDENTIFIER) && !p.isTypeToken(p.currentToken.Type) {
				p.addError(fmt.Sprintf("expected type after 'TEH', got %v at line %d", p.currentToken.Type, p.currentToken.Position.Line))
				return nil
			}
			funcDecl.ReturnType = p.currentToken.Literal
		}

		// Check for parameters
		if p.peekTokenIs(WIT) {
			p.nextToken() // consume WIT
			funcDecl.Parameters = p.parseParameterList()
		}

		// Parse function body
		funcDecl.Body = p.parseStatementBlock(KTHX)

		member.Function = funcDecl
	}

	return member
}

// parseParameterList parses function parameters
func (p *Parser) parseParameterList() []environment.Parameter {
	var params []environment.Parameter

	if !p.expectPeek(IDENTIFIER) {
		p.addError(fmt.Sprintf("expected parameter name, got %v at line %d", p.peekToken.Type, p.peekToken.Position.Line))
		return params
	}

	param := environment.Parameter{
		Name: p.currentToken.Literal,
	}

	if !p.expectPeek(TEH) {
		p.addError(fmt.Sprintf("expected 'TEH' after parameter name, got %v at line %d", p.peekToken.Type, p.peekToken.Position.Line))
		return params
	}

	p.nextToken() // move to the type token
	if !p.currentTokenIs(IDENTIFIER) && !p.isTypeToken(p.currentToken.Type) {
		p.addError(fmt.Sprintf("expected type, got %v at line %d", p.currentToken.Type, p.currentToken.Position.Line))
		return params
	}

	param.Type = p.currentToken.Literal
	params = append(params, param)

	// Parse additional parameters (AN WIT ...)
	for p.peekTokenIs(AN) {
		p.nextToken() // consume AN
		if !p.expectPeek(WIT) {
			p.addError(fmt.Sprintf("expected 'WIT' after 'AN', got %v at line %d", p.peekToken.Type, p.peekToken.Position.Line))
			break
		}
		if !p.expectPeek(IDENTIFIER) {
			p.addError(fmt.Sprintf("expected parameter name, got %v at line %d", p.peekToken.Type, p.peekToken.Position.Line))
			break
		}

		param := environment.Parameter{
			Name: p.currentToken.Literal,
		}

		if !p.expectPeek(TEH) {
			p.addError(fmt.Sprintf("expected 'TEH' after parameter name, got %v at line %d", p.peekToken.Type, p.peekToken.Position.Line))
			break
		}

		p.nextToken()
		if !p.currentTokenIs(IDENTIFIER) && !p.isTypeToken(p.currentToken.Type) {
			p.addError(fmt.Sprintf("expected type, got %v at line %d", p.currentToken.Type, p.currentToken.Position.Line))
			break
		}

		param.Type = p.currentToken.Literal
		params = append(params, param)
	}

	return params
}

// parseStatementBlock parses a block of statements until the end token
func (p *Parser) parseStatementBlock(endTokens ...TokenType) *ast.StatementBlockNode {
	block := &ast.StatementBlockNode{}
	block.Statements = []ast.Node{}
	block.Position = p.convertPosition(p.currentToken.Position)

	p.nextToken()    // move past opening token
	p.skipNewlines() // skip any leading newlines

	for !p.currentTokenIs(endTokens...) && !p.currentTokenIs(EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
		p.skipNewlines() // skip newlines between statements
	}

	return block
}

// skipNewlines skips over NEWLINE tokens
func (p *Parser) skipNewlines() {
	for p.currentTokenIs(NEWLINE) {
		p.nextToken()
	}
}

// parseStatement parses individual statements
func (p *Parser) parseStatement() ast.Node {
	p.skipNewlines() // skip any leading newlines

	switch p.currentToken.Type {
	case NEWLINE:
		return nil // empty statement
	case I:
		// Check if this is "I CAN HAS" (import) or "I HAS A" (variable)
		if p.peekTokenIs(CAN) {
			return p.parseImportStatement()
		} else {
			return p.parseIHasAVariableDeclaration()
		}
	case IZ:
		return p.parseIfStatement()
	case WHILE:
		return p.parseWhileStatement()
	case GIVEZ:
		return p.parseReturnStatement()
	case MAYB:
		return p.parseTryStatement()
	case OOPS:
		return p.parseThrowStatement()
	case IDENTIFIER:
		// Check for new member syntax: IDENTIFIER IDENTIFIER or IDENTIFIER DO IDENTIFIER
		if p.peekTokenIs(IDENTIFIER) {
			// IDENTIFIER IDENTIFIER - member variable access or assignment
			objName := p.currentToken.Literal
			objPos := p.currentToken.Position
			p.nextToken() // move to member name
			memberName := p.currentToken.Literal

			if p.peekTokenIs(ITZ) {
				// Member variable assignment: OBJECT MEMBER ITZ value
				p.nextToken() // consume ITZ
				p.nextToken() // move to value
				return &ast.AssignmentNode{
					Target: &ast.MemberAccessNode{
						Object:   &ast.IdentifierNode{Name: objName, Position: p.convertPosition(objPos)},
						Member:   memberName,
						Position: p.convertPosition(objPos),
					},
					Value:    p.parseExpression(),
					Position: p.convertPosition(objPos),
				}
			} else {
				// Member variable access in statement context (should not happen often)
				return &ast.MemberAccessNode{
					Object:   &ast.IdentifierNode{Name: objName, Position: p.convertPosition(objPos)},
					Member:   memberName,
					Position: p.convertPosition(objPos),
				}
			}
		} else if p.peekTokenIs(DO) {
			// IDENTIFIER DO IDENTIFIER - member function call
			objName := p.currentToken.Literal
			objPos := p.currentToken.Position
			p.nextToken() // consume DO
			if !p.expectPeek(IDENTIFIER) {
				p.addError(fmt.Sprintf("expected identifier after 'DO', got %v at line %d", p.peekToken.Type, p.peekToken.Position.Line))
				return nil
			}
			methodName := p.currentToken.Literal

			if p.peekTokenIs(WIT) {
				// Member function call with arguments: OBJECT DO METHOD WIT args
				p.nextToken() // consume WIT
				args := p.parseArgumentList()
				return &ast.FunctionCallNode{
					Function: &ast.MemberAccessNode{
						Object:   &ast.IdentifierNode{Name: objName, Position: p.convertPosition(objPos)},
						Member:   methodName,
						Position: p.convertPosition(objPos),
					},
					Arguments: args,
					Position:  p.convertPosition(objPos),
				}
			} else {
				// Member function call without arguments: OBJECT DO METHOD
				return &ast.FunctionCallNode{
					Function: &ast.MemberAccessNode{
						Object:   &ast.IdentifierNode{Name: objName, Position: p.convertPosition(objPos)},
						Member:   methodName,
						Position: p.convertPosition(objPos),
					},
					Arguments: []ast.Node{},
					Position:  p.convertPosition(objPos),
				}
			}
		} else if p.peekTokenIs(ITZ) {
			return p.parseAssignment()
		}
		return p.parseExpressionStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// parseIfStatement parses if statements
func (p *Parser) parseIfStatement() *ast.IfStatementNode {
	node := &ast.IfStatementNode{}

	// Set position from the IZ token
	node.Position = p.convertPosition(p.currentToken.Position)

	p.nextToken() // move past IZ
	node.Condition = p.parseExpression()

	if !p.expectPeek(QUESTION) {
		p.addError(fmt.Sprintf("expected '?' after if condition, got %v at line %d", p.peekToken.Type, p.peekToken.Position.Line))
		return nil
	}

	node.ThenBlock = p.parseStatementBlock(KTHX, NOPE)

	// Check for NOPE (else)
	if p.currentTokenIs(NOPE) {
		node.ElseBlock = p.parseStatementBlock(KTHX)
	}

	if !p.currentTokenIs(KTHX) {
		p.addError(fmt.Sprintf("expected 'KTHX' to close if statement, got %v at line %d", p.currentToken.Type, p.currentToken.Position.Line))
		return nil
	}

	return node
}

// parseWhileStatement parses while loops
func (p *Parser) parseWhileStatement() *ast.WhileStatementNode {
	node := &ast.WhileStatementNode{}

	// Set position from the WHILE token
	node.Position = p.convertPosition(p.currentToken.Position)

	p.nextToken() // move past WHILE
	node.Condition = p.parseExpression()

	node.Body = p.parseStatementBlock(KTHX)

	if !p.currentTokenIs(KTHX) {
		p.addError(fmt.Sprintf("expected KTHX to close while statement, got %v at line %d", p.currentToken.Type, p.currentToken.Position.Line))
		return nil
	}

	return node
}

// parseReturnStatement parses return statements
func (p *Parser) parseReturnStatement() *ast.ReturnStatementNode {
	node := &ast.ReturnStatementNode{}

	// Set position from the GIVEZ token
	node.Position = p.convertPosition(p.currentToken.Position)

	if p.peekTokenIs(UP) {
		p.nextToken() // consume UP
		// GIVEZ UP - return nothing
		node.Value = nil
	} else {
		p.nextToken() // move past GIVEZ
		node.Value = p.parseExpression()
	}

	return node
}

// parseAssignment parses assignment statements
func (p *Parser) parseAssignment() *ast.AssignmentNode {
	node := &ast.AssignmentNode{}

	// Set position from the identifier token
	node.Position = p.convertPosition(p.currentToken.Position)

	// Simple assignment: IDENTIFIER ITZ expression
	if p.peekTokenIs(ITZ) {
		node.Target = &ast.IdentifierNode{Name: p.currentToken.Literal, Position: p.convertPosition(p.currentToken.Position)}
		p.nextToken() // consume ITZ
		p.nextToken() // move to expression
		node.Value = p.parseExpression()
		return node
	}

	// Member assignment: IDENTIFIER IDENTIFIER ITZ expression
	if p.peekTokenIs(IDENTIFIER) {
		objName := p.currentToken.Literal
		objPos := p.currentToken.Position
		p.nextToken() // consume first IDENTIFIER
		memberName := p.currentToken.Literal

		node.Target = &ast.MemberAccessNode{
			Object:   &ast.IdentifierNode{Name: objName, Position: p.convertPosition(objPos)},
			Member:   memberName,
			Position: p.convertPosition(objPos),
		}

		if !p.peekTokenIs(ITZ) {
			p.addError(fmt.Sprintf("expected 'ITZ' after member name, got %v at line %d", p.peekToken.Type, p.peekToken.Position.Line))
			return nil
		}
		p.nextToken() // move to expression
		node.Value = p.parseExpression()
		return node
	}

	p.addError(fmt.Sprintf("expected 'ITZ' or member name after '%s', got %v at line %d", p.currentToken.Literal, p.peekToken.Type, p.peekToken.Position.Line))
	return nil
}

// parseExpressionStatement parses expression statements
func (p *Parser) parseExpressionStatement() ast.Node {
	return p.parseExpression()
}

// parseExpression parses expressions with precedence
func (p *Parser) parseExpression() ast.Node {
	return p.parseOrExpression()
}

// parseOrExpression parses OR expressions (lowest precedence)
func (p *Parser) parseOrExpression() ast.Node {
	left := p.parseAndExpression()

	for p.peekTokenIs(OR) {
		op := p.peekToken.Literal
		p.nextToken()
		p.nextToken()
		right := p.parseAndExpression()
		left = &ast.BinaryOpNode{
			Left:     left,
			Operator: op,
			Right:    right,
			Position: left.GetPosition(),
		}
	}

	return left
}

// parseAndExpression parses AN (AND) expressions
func (p *Parser) parseAndExpression() ast.Node {
	left := p.parseEqualityExpression()

	for p.peekTokenIs(AN) && !p.isTokenSequence(AN, WIT) {
		op := p.peekToken.Literal
		p.nextToken()
		p.nextToken()
		right := p.parseEqualityExpression()
		left = &ast.BinaryOpNode{
			Left:     left,
			Operator: op,
			Right:    right,
			Position: left.GetPosition(),
		}
	}

	return left
}

// parseEqualityExpression parses equality expressions (SAEM AS)
func (p *Parser) parseEqualityExpression() ast.Node {
	left := p.parseComparisonExpression()

	if p.peekTokenIs(SAEM) {
		p.nextToken() // consume SAEM
		if !p.expectPeek(AS) {
			return left
		}
		p.nextToken()
		right := p.parseComparisonExpression()
		return &ast.BinaryOpNode{
			Left:     left,
			Operator: "SAEM AS",
			Right:    right,
			Position: left.GetPosition(),
		}
	}

	return left
}

// parseComparisonExpression parses comparison expressions
func (p *Parser) parseComparisonExpression() ast.Node {
	left := p.parseArithmeticExpression()

	if p.peekTokenIs(BIGGR) || p.peekTokenIs(SMALLR) {
		op := p.peekToken.Literal
		p.nextToken()
		if !p.expectPeek(THAN) {
			return left
		}
		p.nextToken()
		right := p.parseArithmeticExpression()
		return &ast.BinaryOpNode{
			Left:     left,
			Operator: op + " THAN",
			Right:    right,
			Position: left.GetPosition(),
		}
	}

	return left
}

// parseArithmeticExpression parses arithmetic expressions
func (p *Parser) parseArithmeticExpression() ast.Node {
	left := p.parseTermExpression()

	for p.peekTokenIs(MOAR) || p.peekTokenIs(LES) {
		op := p.peekToken.Literal
		p.nextToken()
		p.nextToken()
		right := p.parseTermExpression()
		left = &ast.BinaryOpNode{
			Left:     left,
			Operator: op,
			Right:    right,
			Position: left.GetPosition(),
		}
	}

	return left
}

// parseTermExpression parses multiplication/division
func (p *Parser) parseTermExpression() ast.Node {
	left := p.parseCastExpression()

	for p.peekTokenIs(TIEMZ) || p.peekTokenIs(DIVIDEZ) {
		op := p.peekToken.Literal
		p.nextToken()
		p.nextToken()
		right := p.parseCastExpression()
		left = &ast.BinaryOpNode{
			Left:     left,
			Operator: op,
			Right:    right,
			Position: left.GetPosition(),
		}
	}

	return left
}

// parseCastExpression parses cast expressions
func (p *Parser) parseCastExpression() ast.Node {
	left := p.parsePrimaryExpression()

	if p.peekTokenIs(AS) {
		p.nextToken() // consume AS
		p.nextToken() // move to type
		if !p.currentTokenIs(IDENTIFIER) && !p.isTypeToken(p.currentToken.Type) {
			return left
		}
		return &ast.CastNode{
			Expression: left,
			TargetType: p.currentToken.Literal,
			Position:   left.GetPosition(),
		}
	}

	return left
}

// parsePrimaryExpression parses primary expressions
func (p *Parser) parsePrimaryExpression() ast.Node {
	switch p.currentToken.Type {
	case IDENTIFIER:
		// Check for function call, member access, or member function call
		if p.peekTokenIs(WIT) {
			// Regular function call: IDENTIFIER WIT args
			name := p.currentToken.Literal
			pos := p.currentToken.Position
			p.nextToken() // consume WIT
			args := p.parseArgumentList()
			return &ast.FunctionCallNode{
				Function:  &ast.IdentifierNode{Name: name, Position: p.convertPosition(pos)},
				Arguments: args,
				Position:  p.convertPosition(pos),
			}
		} else if p.peekTokenIs(IDENTIFIER) {
			// IDENTIFIER IDENTIFIER - member variable access
			objName := p.currentToken.Literal
			objPos := p.currentToken.Position
			p.nextToken() // move to member name
			memberName := p.currentToken.Literal
			return &ast.MemberAccessNode{
				Object:   &ast.IdentifierNode{Name: objName, Position: p.convertPosition(objPos)},
				Member:   memberName,
				Position: p.convertPosition(objPos),
			}
		} else if p.peekTokenIs(DO) {
			// IDENTIFIER DO IDENTIFIER - member function call
			objName := p.currentToken.Literal
			objPos := p.currentToken.Position
			p.nextToken() // consume DO
			p.nextToken() // move to method name
			if !p.currentTokenIs(IDENTIFIER) {
				p.addError(fmt.Sprintf("expected method name after 'DO', got %v at line %d", p.currentToken.Type, p.currentToken.Position.Line))
				return nil
			}
			methodName := p.currentToken.Literal

			if p.peekTokenIs(WIT) {
				// Member function call with arguments: OBJECT DO METHOD WIT args
				p.nextToken() // consume WIT
				args := p.parseArgumentList()
				return &ast.FunctionCallNode{
					Function: &ast.MemberAccessNode{
						Object:   &ast.IdentifierNode{Name: objName, Position: p.convertPosition(objPos)},
						Member:   methodName,
						Position: p.convertPosition(objPos),
					},
					Arguments: args,
					Position:  p.convertPosition(objPos),
				}
			} else {
				// Member function call without arguments: OBJECT DO METHOD
				return &ast.FunctionCallNode{
					Function: &ast.MemberAccessNode{
						Object:   &ast.IdentifierNode{Name: objName, Position: p.convertPosition(objPos)},
						Member:   methodName,
						Position: p.convertPosition(objPos),
					},
					Arguments: []ast.Node{},
					Position:  p.convertPosition(objPos),
				}
			}
		}
		return &ast.IdentifierNode{Name: p.currentToken.Literal, Position: p.convertPosition(p.currentToken.Position)}

	case STRING, INTEGER, DOUBLE, YEZ, NO, NOTHIN:
		return p.parseLiteral()

	case NEW:
		return p.parseObjectInstantiation()

	case LPAREN:
		// Parse parenthesized expression
		p.nextToken() // consume '('
		expr := p.parseExpression()
		if !p.expectPeek(RPAREN) {
			p.addError(fmt.Sprintf("expected ')' after expression, got %v at line %d", p.peekToken.Type, p.peekToken.Position.Line))
			return nil
		}
		return expr

	default:
		p.addError(fmt.Sprintf("unexpected token in expression: %v at line %d", p.currentToken.Type, p.currentToken.Position.Line))
		return nil
	}
}

// parseArgumentList parses function arguments
func (p *Parser) parseArgumentList() []ast.Node {
	var args []ast.Node

	p.nextToken() // move to first argument
	args = append(args, p.parseExpression())

	// Parse additional arguments (AN WIT ...)
	for p.peekTokenIs(AN) && p.isTokenSequence(AN, WIT) {
		p.nextToken() // consume AN
		p.nextToken() // consume WIT
		p.nextToken() // move to expression
		args = append(args, p.parseExpression())
	}

	return args
}

// parseLiteral parses literal values
func (p *Parser) parseLiteral() *ast.LiteralNode {
	value, err := ConvertValue(p.currentToken)
	if err != nil {
		p.addError(fmt.Sprintf("error converting literal: %v", err))
		return nil
	}

	return &ast.LiteralNode{
		Value:    types.ValueOf(value),
		Position: p.convertPosition(p.currentToken.Position),
	}
}

// parseObjectInstantiation parses object instantiation
func (p *Parser) parseObjectInstantiation() *ast.ObjectInstantiationNode {
	node := &ast.ObjectInstantiationNode{}

	// Set position from NEW token
	node.Position = p.convertPosition(p.currentToken.Position)

	p.nextToken() // move past NEW
	if !p.currentTokenIs(IDENTIFIER) && !p.isTypeToken(p.currentToken.Type) {
		p.addError(fmt.Sprintf("expected class name after 'NEW', got %v at line %d", p.currentToken.Type, p.currentToken.Position.Line))
		return nil
	}

	node.ClassName = p.currentToken.Literal

	// Check for constructor arguments (WIT ...)
	if p.peekTokenIs(WIT) {
		p.nextToken() // consume WIT
		node.ConstructorArgs = p.parseArgumentList()
	}

	return node
}

// parseTryStatement parses try-catch-finally blocks
func (p *Parser) parseTryStatement() *ast.TryStatementNode {
	node := &ast.TryStatementNode{}

	// Set position from MAYB token
	node.Position = p.convertPosition(p.currentToken.Position)

	// Parse try block: MAYB ... statements ... until OOPSIE WIT
	// Now we can use parseStatementBlock and look for OOPSIE WIT pattern
	block := &ast.StatementBlockNode{Statements: []ast.Node{}, Position: p.convertPosition(p.currentToken.Position)}

	p.nextToken() // move past MAYB
	p.skipNewlines()

	// Parse statements until we find OOPSIE IDENTIFIER (catch clause)
	for !p.currentTokenIs(EOF) {
		if p.currentTokenIs(OOPSIE) && p.peekTokenIs(IDENTIFIER) {
			// Found catch clause
			break
		}

		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
		p.skipNewlines()
	}

	node.TryBody = block

	if !p.currentTokenIs(OOPSIE) || !p.peekTokenIs(IDENTIFIER) {
		p.addError(fmt.Sprintf("expected 'OOPSIE identifier' to start catch block, got %v %v at line %d", p.currentToken.Type, p.peekToken.Type, p.currentToken.Position.Line))
		return nil
	}

	// Parse catch clause: OOPSIE variable_name
	p.nextToken() // consume OOPSIE

	if !p.currentTokenIs(IDENTIFIER) {
		p.addError(fmt.Sprintf("expected identifier after 'OOPSIE', got %v at line %d", p.currentToken.Type, p.currentToken.Position.Line))
		return nil
	}
	node.CatchVar = p.currentToken.Literal

	// Parse catch body until KTHX or ALWAYZ
	node.CatchBody = p.parseStatementBlock(KTHX, ALWAYZ)

	// Check for optional finally block: ALWAYZ ... statements ... KTHX
	if p.currentTokenIs(ALWAYZ) {
		node.FinallyBody = p.parseStatementBlock(KTHX)
	}

	if !p.currentTokenIs(KTHX) {
		p.addError(fmt.Sprintf("expected 'KTHX' to close try statement, got %v at line %d", p.currentToken.Type, p.currentToken.Position.Line))
		return nil
	}

	return node
}

// parseThrowStatement parses throw statements
func (p *Parser) parseThrowStatement() *ast.ThrowStatementNode {
	node := &ast.ThrowStatementNode{}

	// Set position from OOPS token
	node.Position = p.convertPosition(p.currentToken.Position)

	// Parse: OOPSIE expression (no WIT for throw)
	p.nextToken() // move to expression
	node.Expression = p.parseExpression()

	return node
}

// Helper methods

func (p *Parser) isTypeToken(tokenType TokenType) bool {
	return tokenType == BOOL_TYPE || tokenType == INTEGR_TYPE ||
		tokenType == DUBBLE_TYPE || tokenType == STRIN_TYPE
}

func (p *Parser) isTokenSequence(first, second TokenType) bool {
	return p.peekToken.Type == first &&
		len(p.lexer.input) > p.lexer.readPosition &&
		p.getTokenAfterPeek() == second
}

// getTokenAfterPeek is a helper to look two tokens ahead
func (p *Parser) getTokenAfterPeek() TokenType {
	// This is a simplified version - in a real implementation,
	// we'd need a proper lookahead buffer
	oldPos := p.lexer.position
	oldReadPos := p.lexer.readPosition
	oldCh := p.lexer.ch
	oldLine := p.lexer.line
	oldCol := p.lexer.column

	// Skip current peek token
	tempLexer := &Lexer{
		input:        p.lexer.input,
		position:     p.lexer.readPosition,
		readPosition: p.lexer.readPosition + 1,
		ch:           0,
		line:         p.lexer.line,
		column:       p.lexer.column,
	}
	if tempLexer.position < len(tempLexer.input) {
		tempLexer.ch = tempLexer.input[tempLexer.position]
	}

	// Get next token
	nextTok, _ := tempLexer.NextToken()

	// Restore lexer state
	p.lexer.position = oldPos
	p.lexer.readPosition = oldReadPos
	p.lexer.ch = oldCh
	p.lexer.line = oldLine
	p.lexer.column = oldCol

	return nextTok.Type
}
