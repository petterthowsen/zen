package parsing

import (
	"zen/lang/lexing"
	"zen/lang/parsing/ast"
	"zen/lang/parsing/statement"
)

// parseIfStatement parses an if statement
func (p *Parser) parseIfStatement() ast.Statement {
	startToken := p.previous() // The 'if' token

	// parse the first IfConditionBlock explicitly since it is required
	primaryCondition := p.parseExpression()
	if primaryCondition == nil {
		p.errorAtToken(startToken, "Expected expression after 'if'")
		return nil
	}

	// parse the primary block
	if !p.match(lexing.LEFT_BRACE) {
		p.error("Expected '{' after 'if'")
		return nil
	}

	primaryBlock := p.parseBlock()

	if !p.match(lexing.RIGHT_BRACE) {
		p.error("Expected '}' after if block")
		return nil
	}

	// parse optional 'elif' condition blocks

	elifBlocks := make([]*statement.IfConditionBlock, 0)

	for {
		tok := p.peek()
		if tok.Type == lexing.EOF {
			break
		}

		if p.matchKeyword("elif") {
			block := p.parseIfConditionBlock()

			if !p.match(lexing.RIGHT_BRACE) {
				p.error("Expected '}' after elif block")
				return nil
			}

			elifBlocks = append(elifBlocks, block)
		} else {
			break
		}
	}

	// parse optional 'else' block
	elseBlock := make([]ast.Statement, 0)

	if p.matchKeyword("else") {
		if !p.match(lexing.LEFT_BRACE) {
			p.error("Expected '{' after else")
			return nil
		}
		elseBlock = p.parseBlock()

		if !p.match(lexing.RIGHT_BRACE) {
			p.error("Expected '}' after else block")
			return nil
		}
	}

	return &statement.IfStatement{
		Location:         startToken.Location,
		PrimaryCondition: primaryCondition,
		PrimaryBlock:     primaryBlock,
		ElseIfBlocks:     elifBlocks,
		ElseBlock:        elseBlock,
	}
}

// parseIfConditionBlock  parses a condition followed by a block
// I.E:
// 1+1 == 2 {
// }
func (p *Parser) parseIfConditionBlock() *statement.IfConditionBlock {
	condition := p.parseExpression()
	if condition == nil {
		return nil
	}

	// Parse body
	if !p.match(lexing.LEFT_BRACE) {
		p.error("Expected '{' after if condition")
		return nil
	}

	body := p.parseBlock()

	return &statement.IfConditionBlock{
		Location:  condition.GetLocation(),
		Condition: condition,
		Body:      body,
	}
}
