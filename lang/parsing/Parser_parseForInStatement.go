package parsing

import (
	"zen/lang/lexing"
	"zen/lang/parsing/ast"
	"zen/lang/parsing/statement"
)

// parseForInStatement parses a for-in loop
// Syntax:
//
//	for key, value in container { body }
//	for value in container { body }
func (p *Parser) parseForInStatement() ast.Statement {
	startToken := p.previous() // The 'for' token

	// Parse key and value identifiers
	var key string
	var value string

	// First identifier is required (either key or value)
	if !p.check(lexing.IDENTIFIER) {
		p.error("Expected identifier after 'for'")
		return nil
	}
	firstIdent := p.advance()

	// Check if we have a comma (indicating key, value syntax)
	if p.match(lexing.COMMA) {
		// First identifier was the key
		key = firstIdent.Literal

		// Parse value identifier
		if !p.check(lexing.IDENTIFIER) {
			p.error("Expected value identifier after comma")
			return nil
		}
		value = p.advance().Literal
	} else {
		// No comma, so first identifier was the value
		value = firstIdent.Literal
	}

	// Parse 'in' keyword
	if !p.matchKeyword("in") {
		p.error("Expected 'in' after for loop variables")
		return nil
	}

	// Parse container expression
	container := p.parseExpression()
	if container == nil {
		p.error("Expected expression after 'in'")
		return nil
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

	return statement.NewForInStatement(
		key,
		value,
		container,
		body,
		startToken.Location,
	)
}
