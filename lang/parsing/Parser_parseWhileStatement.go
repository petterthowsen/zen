package parsing

import (
	"zen/lang/lexing"
	"zen/lang/parsing/ast"
	"zen/lang/parsing/statement"
)

// parseWhileStatement parses a while loop
// Syntax: while condition { body }
func (p *Parser) parseWhileStatement() ast.Statement {
	startToken := p.previous() // The 'while' token

	// Parse condition
	condition := p.parseExpression()
	if condition == nil {
		p.error("Expected condition after 'while'")
		return nil
	}

	// Parse body
	if !p.match(lexing.LEFT_BRACE) {
		p.error("Expected '{' before while loop body")
		return nil
	}

	body := p.parseBlock()

	if !p.match(lexing.RIGHT_BRACE) {
		p.error("Expected '}' after while loop body")
		return nil
	}

	return statement.NewWhileStatement(
		condition,
		body,
		startToken.Location,
	)
}
