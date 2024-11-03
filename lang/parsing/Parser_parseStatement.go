package parsing

import (
	"zen/lang/lexing"
	"zen/lang/parsing/ast"
	"zen/lang/parsing/statement"
)

// parseStatement parses a statement by delegating to other methods like parseVarDeclaration, parseIfStatement etc.
func (p *Parser) parseStatement() ast.Statement {
	// var/const declaration
	if p.matchKeyword("var", "const") {
		return p.parseVarDeclaration()
	}

	// func declaration
	if p.matchKeyword("func") {
		return p.parseFuncDeclaration()
	}

	// If Statement
	if p.matchKeyword("if") {
		return p.parseIfStatement()
	}

	// For Statement
	if p.matchKeyword("for") {
		return p.parseForStatement()
	}

	// Return statement
	if p.matchKeyword("return") {
		return p.parseReturnStatement()
	}

	// Try parsing an expression statement
	expr := p.parseExpression()
	if expr != nil {
		stmt := &statement.ExpressionStatement{
			Location:   p.previous().Location,
			Expression: expr,
		}
		return stmt
	}

	token := p.peek()
	if token.Type == lexing.EOF {
		return nil // End of file is not an error
	}

	p.errorAtToken(token, "Expected statement")
	return nil
}
