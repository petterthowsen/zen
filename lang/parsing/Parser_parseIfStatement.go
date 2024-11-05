package parsing

import (
	"zen/lang/lexing"
	"zen/lang/parsing/ast"
	"zen/lang/parsing/expression"
	"zen/lang/parsing/statement"
)

// parseIfStatement parses an if statement with optional else/elif blocks
func (p *Parser) parseIfStatement() ast.Statement {
	startToken := p.previous() // The 'if' token

	// Save state before parsing condition
	current := p.current
	errorCount := len(p.errors)

	// Try parsing condition with map access enabled
	condition := p.parseExpression()
	if condition == nil {
		p.error("Expected condition after 'if'")
		return nil
	}

	// If we got a map access but no LEFT_BRACE follows, try again without map access
	if _, isMapAccess := condition.(*expression.MapAccessExpression); isMapAccess {
		if !p.check(lexing.LEFT_BRACE) {
			// Restore state and try again with map access disabled
			p.current = current
			p.errors = p.errors[:errorCount]
			p.DisableMapAccess()
			condition = p.parseExpression()
			p.EnableMapAccess()
		}
	}

	if !p.match(lexing.LEFT_BRACE) {
		p.error("Expected '{' after 'if' condition")
		return nil
	}

	body := p.parseBlock()

	if !p.match(lexing.RIGHT_BRACE) {
		p.error("Expected '}' after if body")
		return nil
	}

	var elseIfBlocks []*statement.IfConditionBlock

	// Parse elif blocks
	for p.matchKeyword("elif") {
		elifToken := p.previous()

		// Save state before parsing elif condition
		current = p.current
		errorCount = len(p.errors)

		// Try parsing condition with map access enabled
		elifCondition := p.parseExpression()
		if elifCondition == nil {
			p.error("Expected condition after 'elif'")
			return nil
		}

		// If we got a map access but no LEFT_BRACE follows, try again without map access
		if _, isMapAccess := elifCondition.(*expression.MapAccessExpression); isMapAccess {
			if !p.check(lexing.LEFT_BRACE) {
				// Restore state and try again with map access disabled
				p.current = current
				p.errors = p.errors[:errorCount]
				p.DisableMapAccess()
				elifCondition = p.parseExpression()
				p.EnableMapAccess()
			}
		}

		if !p.match(lexing.LEFT_BRACE) {
			p.error("Expected '{' after 'elif' condition")
			return nil
		}

		elifBody := p.parseBlock()

		if !p.match(lexing.RIGHT_BRACE) {
			p.error("Expected '}' after elif body")
			return nil
		}

		elseIfBlocks = append(elseIfBlocks, statement.NewIfConditionBlock(elifCondition, elifBody, elifToken.Location))
	}

	// Parse optional else block
	var elseBody []ast.Statement
	if p.matchKeyword("else") {
		if !p.match(lexing.LEFT_BRACE) {
			p.error("Expected '{' after 'else'")
			return nil
		}

		elseBody = p.parseBlock()

		if !p.match(lexing.RIGHT_BRACE) {
			p.error("Expected '}' after else body")
			return nil
		}
	}

	return statement.NewIfStatement(condition, body, elseIfBlocks, elseBody, startToken.Location)
}
