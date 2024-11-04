package parsing

import (
	"strconv"
	"zen/lang/common"
	"zen/lang/lexing"
	"zen/lang/parsing/ast"
	"zen/lang/parsing/expression"
	"zen/lang/parsing/statement"
)

// parseType parses a type annotation, which can be either a basic type or a parametric type
func (p *Parser) parseType() ast.Expression {
	// Parse the base type name, which can be either a keyword (primitive) or identifier (user type)
	var typeToken lexing.Token
	if p.check(lexing.KEYWORD) {
		// Handle primitive types (string, int, float64, etc.)
		typeToken = p.advance()
	} else if p.check(lexing.IDENTIFIER) {
		// Handle user-defined types (Array, MyClass, etc.)
		typeToken = p.consume(lexing.IDENTIFIER, "Expected type name")
	} else {
		p.errorAtToken(p.peek(), "Expected KEYWORD or IDENTIFIER for type name")
		return nil
	}

	// If there's no angle bracket, it's a basic type
	if !p.match(lexing.LESS) {
		return expression.NewBasicType(typeToken.Literal, typeToken.Location)
	}

	// Parse parametric type parameters
	var params []expression.Parameter

	// Must have at least one parameter
	if p.check(lexing.GREATER) {
		p.error("Expected at least one type parameter")
		return nil
	}

	// Parse first parameter
	param := p.parseTypeParameter()
	if param == nil {
		return nil
	}
	params = append(params, *param)

	// Parse additional parameters
	for p.match(lexing.COMMA) {
		// Don't allow trailing comma
		if p.check(lexing.GREATER) {
			p.error("Unexpected trailing comma")
			return nil
		}

		param = p.parseTypeParameter()
		if param == nil {
			return nil
		}
		params = append(params, *param)
	}

	// Expect closing angle bracket
	if !p.match(lexing.GREATER) {
		p.error("Expected '>' after type parameters")
		return nil
	}

	return expression.NewParametricType(typeToken.Literal, params, typeToken.Location)
}

// parseTypeParameter parses a single type parameter, which can be:
// - A type name (keyword or identifier)
// - A nested parametric type
// - An integer literal
func (p *Parser) parseTypeParameter() *expression.Parameter {
	var location *common.SourceLocation

	// Try to parse an integer parameter first
	if p.check(lexing.INT) {
		token := p.advance()
		location = token.Location
		value, err := strconv.ParseInt(token.Literal, 10, 64)
		if err != nil {
			p.errorAtToken(token, "Invalid integer literal")
			return nil
		}
		return &expression.Parameter{
			Value:    value,
			IsType:   false,
			Location: location,
		}
	}

	// Try to parse a type parameter (keyword, identifier, or nested type)
	if p.check(lexing.KEYWORD) || p.check(lexing.IDENTIFIER) {
		token := p.advance()
		location = token.Location

		// If it's followed by a less-than, it's a nested parametric type
		if p.check(lexing.LESS) {
			// Parse the nested type
			if !p.match(lexing.LESS) {
				p.error("Expected '<' after type name")
				return nil
			}

			var params []expression.Parameter

			// Parse first parameter
			param := p.parseTypeParameter()
			if param == nil {
				return nil
			}
			params = append(params, *param)

			// Parse additional parameters
			for p.match(lexing.COMMA) {
				if p.check(lexing.GREATER) {
					p.error("Unexpected trailing comma")
					return nil
				}

				param = p.parseTypeParameter()
				if param == nil {
					return nil
				}
				params = append(params, *param)
			}

			// Expect closing angle bracket
			if !p.match(lexing.GREATER) {
				p.error("Expected '>' after type parameters")
				return nil
			}

			return &expression.Parameter{
				Value:    expression.NewParametricType(token.Literal, params, location),
				IsType:   true,
				Location: location,
			}
		}

		// Otherwise it's a basic type
		return &expression.Parameter{
			Value:    token.Literal,
			IsType:   true,
			Location: location,
		}
	}

	p.error("Expected type name or integer")
	return nil
}

// parseVarDeclaration: Parses a variable declaration
func (p *Parser) parseVarDeclaration() ast.Statement {
	isConstant := p.previous().Literal == "const"
	startToken := p.previous() // Save the 'var' or 'const' token for error reporting

	// Parse variable name
	name := p.consume(lexing.IDENTIFIER, "Expected variable name")
	if len(p.errors) > 0 {
		return nil
	}

	isNullable := false
	var varType ast.Expression

	// Parse optional type annotation
	if p.match(lexing.COLON) {
		varType = p.parseType()
		if varType == nil {
			return nil
		}

		// Handle nullable type
		if p.match(lexing.QMARK) {
			isNullable = true
		}
	}

	// Parse optional initializer
	var initializer ast.Expression
	if p.match(lexing.ASSIGN) {
		initializer = p.parseExpression()
		if initializer == nil {
			return nil
		}
	}

	return statement.NewVarDeclarationNode(
		name.Literal,
		varType,
		initializer,
		isConstant,
		isNullable,
		startToken.Location,
	)
}
