package parsing

import (
	"zen/lang/lexing"
	"zen/lang/parsing/ast"
	"zen/lang/parsing/expression"
	"zen/lang/parsing/statement"
)

// parseFuncDeclaration: Parses a function declaration
func (p *Parser) parseFuncDeclaration() ast.Statement {
	startToken := p.previous() // Save the 'func' token for error reporting

	// Parse function name
	name := p.consume(lexing.IDENTIFIER, "Expected function name")
	if len(p.errors) > 0 {
		return nil
	}

	// Parse function parameters
	if !p.match(lexing.LEFT_PAREN) {
		p.error("Expected '(' after function name")
		return nil
	}

	// Parse function parameters with comma between each parameter
	parameters := make([]expression.FuncParameterExpression, 0)

	// Handle parameters
	if !p.check(lexing.RIGHT_PAREN) {
		for {
			// Parse parameter
			param := p.parseFuncParameter()
			if param == nil {
				// Error already reported by parseFuncParameter
				return nil
			}
			parameters = append(parameters, *param)

			// Check for comma or end of parameters
			if p.check(lexing.RIGHT_PAREN) {
				break
			}

			if !p.match(lexing.COMMA) {
				p.error("Expected ',' or ')' after function parameter")
				return nil
			}
		}
	}

	// Consume the closing parenthesis
	if !p.match(lexing.RIGHT_PAREN) {
		p.error("Expected ')' after function parameters")
		return nil
	}

	// Parse optional return type (defaults to "void" if not specified)
	returnType := "void"
	if p.match(lexing.COLON) {
		tok := p.peek()

		if tok.Type == lexing.KEYWORD || tok.Type == lexing.IDENTIFIER {
			returnType = tok.Literal
		} else {
			p.error("Invalid function return type " + tok.Literal)
		}

		p.advance()
	}

	// Parse function body
	if !p.match(lexing.LEFT_BRACE) {
		p.error("Expected '{' after function declaration")
		return nil
	}

	// Parse the function body
	body := p.parseBlock()
	if body == nil {
		// Error already reported by parseBlock
		return nil
	}

	// Consume the closing brace
	if !p.match(lexing.RIGHT_BRACE) {
		p.error("Expected '}' after function body")
		return nil
	}

	return statement.NewFuncDeclaration(
		name.Literal,
		parameters,
		returnType,
		body,
		startToken.Location,
	)
}

// parseFuncParameter: Parses a function parameter
func (p *Parser) parseFuncParameter() *expression.FuncParameterExpression {
	name := p.consume(lexing.IDENTIFIER, "Expected parameter name")
	if len(p.errors) > 0 {
		return nil
	}

	// type
	if !p.match(lexing.COLON) {
		p.error("Expected ':' after parameter name")
		return nil
	}

	typeToken := p.consume(lexing.IDENTIFIER, "Expected parameter type")
	if len(p.errors) > 0 {
		return nil
	}

	// Check for nullable type marker
	isNullable := p.match(lexing.QMARK)

	// optional default value
	var defaultValue ast.Expression
	if p.match(lexing.ASSIGN) {
		defaultValue = p.parseExpression()
		if defaultValue == nil {
			p.error("Expected expression after '='")
			return nil
		}
	}

	return expression.NewFuncParameterExpression(
		name.Literal,
		typeToken.Literal,
		isNullable,
		name.Location,
		defaultValue,
	)
}
