package parsing

import (
	"zen/lang/lexing"
	"zen/lang/parsing/ast"
	"zen/lang/parsing/expression"
)

// parseMapAccessExpression parses map access expressions like map["name"] or map[keyName]
func (p *Parser) parseMapAccessExpression(mapExpr ast.Expression) ast.Expression {
	// We've already consumed the '{'
	var key ast.Expression

	// Parse key - only allow string literals and identifiers for now
	token := p.peek()
	switch token.Type {
	case lexing.STRING:
		p.advance()
		key = expression.NewLiteralExpression(token.Literal, token.Location)
	case lexing.IDENTIFIER:
		p.advance()
		key = expression.NewIdentifierExpression(token.Literal, token.Location)
	default:
		p.error("Expected string or identifier for map key")
		return nil
	}

	if !p.match(lexing.RIGHT_BRACE) {
		p.error("Expected '}' after map key")
		return nil
	}

	return expression.NewMapAccessExpression(mapExpr, key, p.previous().Location)
}
