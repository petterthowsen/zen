package parsing

import (
	"zen/lang/lexing"
	"zen/lang/parsing/ast"
	"zen/lang/parsing/statement"
)

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

	// Parse optional type annotation
	var varType string
	if p.match(lexing.COLON) {
		typeToken := p.consume(lexing.IDENTIFIER, "Expected type name")
		if len(p.errors) > 0 {
			return nil
		}
		varType = typeToken.Literal

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
