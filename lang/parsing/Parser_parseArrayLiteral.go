package parsing

import (
	"zen/lang/lexing"
	"zen/lang/parsing/ast"
	"zen/lang/parsing/expression"
)

// parseArrayLiteral parses array literals like [1, 2, 3]
func (p *Parser) parseArrayLiteral() ast.Expression {
	openBracket := p.advance() // consume [
	elements := make([]ast.Expression, 0)

	// Handle empty array
	if p.match(lexing.RIGHT_BRACKET) {
		return expression.NewArrayLiteralExpression(elements, openBracket.Location)
	}

	// Parse first element
	expr := p.parseExpression()
	if expr == nil {
		return nil
	}
	elements = append(elements, expr)

	// Parse remaining elements
	for p.match(lexing.COMMA) {
		// Allow trailing comma
		if p.check(lexing.RIGHT_BRACKET) {
			break
		}

		expr = p.parseExpression()
		if expr == nil {
			return nil
		}
		elements = append(elements, expr)
	}

	if !p.match(lexing.RIGHT_BRACKET) {
		p.error("Expected ']' after array elements")
		return nil
	}

	return expression.NewArrayLiteralExpression(elements, openBracket.Location)
}
