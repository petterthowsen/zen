package parsing

import (
	"zen/lang/parsing/ast"
	"zen/lang/parsing/statement"
)

// parseBreakStatement parses a break statement
func (p *Parser) parseBreakStatement() ast.Statement {
	startToken := p.previous() // The 'break' token
	return statement.NewBreakStatement(startToken.Location)
}

// parseContinueStatement parses a continue statement
func (p *Parser) parseContinueStatement() ast.Statement {
	startToken := p.previous() // The 'continue' token
	return statement.NewContinueStatement(startToken.Location)
}
