package parsing

import (
	"strconv"
	"zen/lang/lexing"
	"zen/lang/parsing/ast"
	"zen/lang/parsing/expression"
)

// parseExpression: Entry point for expression parsing
func (p *Parser) parseExpression() ast.Expression {
	return p.parseAssignment()
}

// parseAssignment parses assignment expressions including compound assignments (+=, -=, etc.)
func (p *Parser) parseAssignment() ast.Expression {
	expr := p.parseLogicalOr()

	if p.match(lexing.ASSIGN, lexing.PLUS_ASSIGN, lexing.MINUS_ASSIGN, lexing.MULTIPLY_ASSIGN, lexing.DIVIDE_ASSIGN) {
		operator := p.previous()
		right := p.parseAssignment()
		if right == nil {
			p.error("Expected expression after assignment operator")
			return nil
		}

		// For compound assignments (+=, -=, etc.), transform them into regular assignments
		// e.g., x += 1 becomes x = x + 1
		if operator.Type != lexing.ASSIGN {
			// Get the base operator from the compound operator
			var baseOp string
			switch operator.Type {
			case lexing.PLUS_ASSIGN:
				baseOp = "+"
			case lexing.MINUS_ASSIGN:
				baseOp = "-"
			case lexing.MULTIPLY_ASSIGN:
				baseOp = "*"
			case lexing.DIVIDE_ASSIGN:
				baseOp = "/"
			}

			// Create a binary expression for the right side
			// e.g., for x += 1, create (x + 1)
			right = expression.NewBinaryExpression(expr, baseOp, right, operator.Location)
		}

		// Create the assignment expression
		return expression.NewBinaryExpression(expr, "=", right, operator.Location)
	}

	return expr
}

// parseLogicalOr parses logical OR expressions
func (p *Parser) parseLogicalOr() ast.Expression {
	expr := p.parseLogicalAnd()

	for p.matchKeyword("or") {
		operator := p.previous().Literal
		right := p.parseLogicalAnd()
		if right == nil {
			p.error("Expected expression after 'or'")
			return nil
		}
		expr = expression.NewBinaryExpression(expr, operator, right, p.previous().Location)
	}

	return expr
}

// parseLogicalAnd parses logical AND expressions
func (p *Parser) parseLogicalAnd() ast.Expression {
	expr := p.parseEquality()

	for p.matchKeyword("and") {
		operator := p.previous().Literal
		right := p.parseEquality()
		if right == nil {
			p.error("Expected expression after 'and'")
			return nil
		}
		expr = expression.NewBinaryExpression(expr, operator, right, p.previous().Location)
	}

	return expr
}

// parseEquality parses equality expressions
func (p *Parser) parseEquality() ast.Expression {
	expr := p.parseComparison()

	for p.match(lexing.EQUALS, lexing.NOT_EQUALS) {
		operator := p.previous().Literal
		right := p.parseComparison()
		if right == nil {
			p.error("Expected expression after comparison operator")
			return nil
		}
		expr = expression.NewBinaryExpression(expr, operator, right, p.previous().Location)
	}

	return expr
}

// parseComparison parses comparison expressions
func (p *Parser) parseComparison() ast.Expression {
	expr := p.parseAdditive()

	for p.match(lexing.LESS, lexing.LESS_EQUALS, lexing.GREATER, lexing.GREATER_EQUALS) {
		operator := p.previous().Literal
		right := p.parseAdditive()
		if right == nil {
			p.error("Expected expression after comparison operator")
			return nil
		}
		expr = expression.NewBinaryExpression(expr, operator, right, p.previous().Location)
	}

	return expr
}

// parseAdditive: Parses addition and subtraction
func (p *Parser) parseAdditive() ast.Expression {
	expr := p.parseMultiplicative()

	for p.match(lexing.PLUS, lexing.MINUS) {
		operator := p.previous().Literal
		right := p.parseMultiplicative()
		if right == nil {
			if p.isAtEnd() {
				p.errorAtToken(p.peek(), "Expected expression after operator")
			} else if p.peek().Type == lexing.PLUS {
				p.errorAtToken(p.peek(), "Expected expression between operators")
			} else {
				p.errorAtToken(p.peek(), "Expected expression after operator")
			}
			return nil
		}
		expr = expression.NewBinaryExpression(expr, operator, right, p.previous().Location)
	}

	return expr
}

// parseMultiplicative: Parses multiplication and division
func (p *Parser) parseMultiplicative() ast.Expression {
	if p.check(lexing.MULTIPLY) || p.check(lexing.DIVIDE) {
		p.errorAtToken(p.peek(), "Expected expression before operator")
		p.advance() // Skip the operator
		return nil
	}

	expr := p.parseUnary()

	for p.match(lexing.MULTIPLY, lexing.DIVIDE) {
		operator := p.previous().Literal
		right := p.parseUnary()
		if right == nil {
			p.errorAtToken(p.peek(), "Expected expression after operator")
			return nil
		}
		expr = expression.NewBinaryExpression(expr, operator, right, p.previous().Location)
	}

	return expr
}

// parseUnary: Parses unary operators
func (p *Parser) parseUnary() ast.Expression {
	if p.match(lexing.MINUS) || p.matchKeyword("not") || p.matchKeyword("await") {
		operator := p.previous()
		expr := p.parseUnary()
		if expr == nil {
			p.errorAtToken(p.peek(), "Expected expression after unary operator")
			return nil
		}
		if operator.Literal == "await" {
			return expression.NewAwaitExpression(expr, operator.Location)
		}
		return expression.NewUnaryExpression(operator.Literal, expr, operator.Location)
	}

	return p.parsePostfix()
}

// parsePostfix parses postfix operators (++, --)
func (p *Parser) parsePostfix() ast.Expression {
	expr := p.parseCall()

	if p.match(lexing.INCREMENT, lexing.DECREMENT) {
		operator := p.previous()
		var baseOp string
		if operator.Type == lexing.INCREMENT {
			baseOp = "+"
		} else {
			baseOp = "-"
		}

		// Create a literal 1
		oneExpr := expression.NewLiteralExpression(int64(1), operator.Location)

		// Create the binary expression for the operation (i + 1 or i - 1)
		addExpr := expression.NewBinaryExpression(expr, baseOp, oneExpr, operator.Location)

		// Create the assignment (i = i + 1 or i = i - 1)
		return expression.NewBinaryExpression(expr, "=", addExpr, operator.Location)
	}

	return expr
}

// parseCall: Parses function calls and member access
func (p *Parser) parseCall() ast.Expression {
	expr := p.parsePrimary()
	if expr == nil {
		return nil
	}

	for {
		if p.match(lexing.LEFT_PAREN) {
			expr = p.finishCall(expr)
		} else if p.match(lexing.DOT) {
			// Handle member access (obj.prop)
			if !p.check(lexing.IDENTIFIER) {
				p.error("Expected property name after '.'")
				return nil
			}
			name := p.advance()
			// Build member access from left to right
			expr = expression.NewMemberAccessExpression(expr, name.Literal, name.Location)
		} else {
			break
		}
	}

	return expr
}

// finishCall handles the parsing of function call arguments after '(' has been matched
func (p *Parser) finishCall(callee ast.Expression) ast.Expression {
	args := make([]ast.Expression, 0)
	if !p.check(lexing.RIGHT_PAREN) {
		for {
			arg := p.parseExpression()
			if arg == nil {
				p.error("Expected argument in function call")
				return nil
			}
			args = append(args, arg)

			if !p.match(lexing.COMMA) {
				break
			}
		}
	}

	if !p.match(lexing.RIGHT_PAREN) {
		p.error("Expected ')' after function arguments")
		return nil
	}

	return &expression.CallExpression{
		Callee:    callee,
		Arguments: args,
		Location:  p.previous().Location,
	}
}

// parsePrimary: Parses primary expressions (literals, identifiers, and parentheses)
func (p *Parser) parsePrimary() ast.Expression {
	token := p.peek()

	switch token.Type {
	case lexing.LEFT_BRACKET:
		return p.parseArrayLiteral()

	case lexing.LEFT_BRACE:
		return p.parseMapLiteral()

	case lexing.STRING:
		p.advance()
		return expression.NewLiteralExpression(token.Literal, token.Location)

	case lexing.INT:
		p.advance()
		// Convert string to int
		value, err := strconv.ParseInt(token.Literal, 10, 64)
		if err != nil {
			p.errorAtToken(token, "Invalid integer literal")
			return nil
		}
		return expression.NewLiteralExpression(value, token.Location)

	case lexing.FLOAT:
		p.advance()
		// Convert string to float
		value, err := strconv.ParseFloat(token.Literal, 64)
		if err != nil {
			p.errorAtToken(token, "Invalid float literal")
			return nil
		}
		return expression.NewLiteralExpression(value, token.Location)

	case lexing.IDENTIFIER:
		p.advance()
		return expression.NewIdentifierExpression(token.Literal, token.Location)

	case lexing.KEYWORD:
		if token.Literal == "true" || token.Literal == "false" {
			p.advance()
			return expression.NewLiteralExpression(token.Literal == "true", token.Location)
		} else if token.Literal == "null" {
			p.advance()
			return expression.NewLiteralExpression(nil, token.Location)
		}

	case lexing.LEFT_PAREN:
		p.advance()
		expr := p.parseExpression()
		if expr == nil {
			return nil
		}
		if !p.match(lexing.RIGHT_PAREN) {
			p.errorAtToken(p.peek(), "Expected closing parenthesis")
			return nil
		}
		return expr
	}

	p.errorAtToken(token, "Expected expression")
	return nil
}
