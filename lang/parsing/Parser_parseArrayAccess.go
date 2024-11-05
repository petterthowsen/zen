package parsing

import (
	"strconv"
	"zen/lang/lexing"
	"zen/lang/parsing/ast"
	"zen/lang/parsing/expression"
)

// parseArrayAccessExpression parses array access expressions like array[idx] or array[2]
func (p *Parser) parseArrayAccessExpression(array ast.Expression) ast.Expression {
	// We've already consumed the '['
	var index ast.Expression

	// Parse index - only allow literals and identifiers for now
	token := p.peek()
	switch token.Type {
	case lexing.INT:
		p.advance()
		value, err := strconv.ParseInt(token.Literal, 10, 64)
		if err != nil {
			p.errorAtToken(token, "Invalid integer literal")
			return nil
		}
		index = expression.NewLiteralExpression(value, token.Location)
	case lexing.IDENTIFIER:
		p.advance()
		index = expression.NewIdentifierExpression(token.Literal, token.Location)
	default:
		p.error("Expected integer or identifier for array index")
		return nil
	}

	if !p.match(lexing.RIGHT_BRACKET) {
		p.error("Expected ']' after array index")
		return nil
	}

	return expression.NewArrayAccessExpression(array, index, p.previous().Location)
}
