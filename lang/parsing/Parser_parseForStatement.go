package parsing

import (
	"zen/lang/lexing"
	"zen/lang/parsing/ast"
	"zen/lang/parsing/statement"
)

// parseForStatement parses a for loop
// Syntax: for initialization; condition; update { body }
func (p *Parser) parseForStatement() ast.Statement {
	startToken := p.previous() // The 'for' token

	// Parse initialization
	var init ast.Statement
	if !p.check(lexing.SEMICOLON) {
		expr := p.parseExpression()
		if expr != nil {
			init = &statement.ExpressionStatement{
				Location:   expr.GetLocation(),
				Expression: expr,
			}
		}
	}

	// Consume the semicolon after initialization
	if !p.match(lexing.SEMICOLON) {
		p.error("Expected ';' after for loop initialization")
		return nil
	}

	// Parse condition
	var condition ast.Expression
	if !p.check(lexing.SEMICOLON) {
		condition = p.parseExpression()
		if condition == nil {
			p.error("Invalid condition in for loop")
			return nil
		}
	}

	// Consume the semicolon after condition
	if !p.match(lexing.SEMICOLON) {
		p.error("Expected ';' after for loop condition")
		return nil
	}

	// Parse update
	var update ast.Statement
	if !p.check(lexing.LEFT_BRACE) {
		expr := p.parseExpression()
		if expr != nil {
			update = &statement.ExpressionStatement{
				Location:   expr.GetLocation(),
				Expression: expr,
			}
		}
	}

	// Parse body
	if !p.match(lexing.LEFT_BRACE) {
		p.error("Expected '{' before for loop body")
		return nil
	}

	body := p.parseBlock()

	if !p.match(lexing.RIGHT_BRACE) {
		p.error("Expected '}' after for loop body")
		return nil
	}

	return statement.NewForStatement(
		init,
		condition,
		update,
		body,
		startToken.Location,
	)
}
