package parsing

import (
	"zen/lang/lexing"
	"zen/lang/parsing/ast"
	"zen/lang/parsing/statement"
)

func (p *Parser) parseReturnStatement() *statement.ReturnStatmenet {
	startToken := p.previous() // The 'return' token

	// after 'return' could be a closing } or an expression
	var exp ast.Expression = nil

	if !p.check(lexing.RIGHT_BRACE) {
		exp = p.parseExpression()
	}

	return &statement.ReturnStatmenet{
		Location:   startToken.Location,
		Expression: exp,
	}
}
